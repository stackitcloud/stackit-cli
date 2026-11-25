package bucket

import (
	"bytes"
	"context"
	"crypto/md5" //nolint:gosec // Content-MD5 is a storage protocol requirement for request integrity, not used for security
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/stackitcloud/stackit-cli/internal/pkg/config"
)

const (
	clientTimeout  = 30 * time.Minute
	maxAttempts    = 3
	retryBaseDelay = time.Second

	contentTypeDefault = "application/octet-stream"
	contentTypeXML     = "application/xml"
)

const fallbackSigningRegion = "eu01"

type Client struct {
	endpoint   string
	region     string
	signer     *Signer
	httpClient *http.Client
	now        func() time.Time
	retryDelay time.Duration
}

func NewClient(cfg ClientConfig) (*Client, error) {
	if cfg.AccessKeyId == "" || cfg.SecretAccessKey == "" {
		return nil, fmt.Errorf("object Storage credentials are required: use the --%s and --%s flags, or set %s/%s in the config file or environment variables",
			AccessKeyIdFlagName, SecretAccessKeyFlagName, config.ObjectStorageAccessKeyIdKey, config.ObjectStorageSecretAccessKeyKey)
	}
	endpoint, err := cfg.ResolvedEndpoint()
	if err != nil {
		return nil, err
	}
	region := cfg.Region
	if region == "" {
		region = fallbackSigningRegion
	}
	endpoint = strings.TrimSuffix(endpoint, "/")
	return &Client{
		endpoint:   endpoint,
		region:     region,
		signer:     NewSigner(cfg.AccessKeyId, cfg.SecretAccessKey, region),
		httpClient: &http.Client{Timeout: clientTimeout},
		now:        time.Now,
		retryDelay: retryBaseDelay,
	}, nil
}

type ListObjectsV2Output struct {
	IsTruncated           bool
	NextContinuationToken string
	Objects               []ObjectInfo
	CommonPrefixes        []string
}

func objectPath(bucket, key string) string {
	if key == "" {
		return "/" + bucket
	}
	return "/" + bucket + "/" + escapeKey(key)
}

func isRetryableStatus(status int) bool {
	return status == http.StatusTooManyRequests || status >= 500
}

func (c *Client) do(ctx context.Context, method, path string, query url.Values, headers map[string]string, body io.Reader, rewind func() error) (*http.Response, error) {
	target := c.endpoint + path
	if len(query) > 0 {
		target += "?" + canonicalQueryString(query)
	}

	var lastErr error
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		if attempt > 1 {
			if rewind != nil {
				if err := rewind(); err != nil {
					return nil, fmt.Errorf("%s %s: rewind request body: %w", method, path, err)
				}
			}
			delay := c.retryDelay << uint(attempt-2)
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(delay):
			}
		}

		var bodyReader io.Reader = http.NoBody
		if body != nil {
			bodyReader = body
		}
		// Buffered bodies need a rewind for retries (their reader is consumed
		// by the first attempt)
		if rewind == nil {
			if reader, ok := body.(*bytes.Reader); ok {
				rewind = func() error {
					_, err := reader.Seek(0, io.SeekStart)
					return err
				}
			}
		}
		req, err := http.NewRequestWithContext(ctx, method, target, bodyReader)
		if err != nil {
			return nil, fmt.Errorf("%s %s: build request: %w", method, path, err)
		}
		// The wire path must be fully encoded so that the path the server
		// receives matches the canonical path used for signing
		req.URL.RawPath = canonicalPath(req)
		// The storage backends do not support chunked transfer encoding, so
		// the content length must be known for file bodies
		if size, ok := fileContentLength(body); ok {
			req.ContentLength = size
		}
		for name, value := range headers {
			req.Header.Set(name, value)
		}
		// The storage server expects the declared payload hash on every request: the
		// empty payload hash for bodyless requests, UNSIGNED-PAYLOAD otherwise
		if req.Header.Get("X-Amz-Content-Sha256") == "" {
			if body == nil {
				req.Header.Set("X-Amz-Content-Sha256", emptyPayloadHash)
			} else {
				req.Header.Set("X-Amz-Content-Sha256", unsignedPayload)
			}
		}
		if err := c.signer.Sign(req, c.now()); err != nil {
			return nil, fmt.Errorf("%s %s: sign request: %w", method, path, err)
		}

		resp, err := c.httpClient.Do(req)
		if err != nil {
			lastErr = err
			if attempt == maxAttempts {
				break
			}
			continue
		}
		if isRetryableStatus(resp.StatusCode) {
			bodyBytes, _ := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
			_ = resp.Body.Close()
			lastErr = parseAPIError(method, resp.StatusCode, bodyBytes)
			if attempt == maxAttempts {
				break
			}
			continue
		}
		return resp, nil
	}
	return nil, fmt.Errorf("%s %s: %w", method, path, lastErr)
}

func (c *Client) doAndCheck(ctx context.Context, method, path string, query url.Values, headers map[string]string, body io.Reader, rewind func() error, out io.Writer) error {
	resp, err := c.do(ctx, method, path, query, headers, body, rewind)
	if err != nil {
		return err
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			fmt.Fprintf(os.Stderr, "close response body: %v\n", cerr)
		}
	}()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
		return parseAPIError(method, resp.StatusCode, bodyBytes)
	}
	if out != nil {
		if _, err := io.Copy(out, resp.Body); err != nil {
			return fmt.Errorf("%s %s: read response body: %w", method, path, err)
		}
	}
	return nil
}

func (c *Client) ListBuckets(ctx context.Context) ([]BucketInfo, error) {
	var buf bytes.Buffer
	if err := c.doAndCheck(ctx, "GET", "/", nil, nil, nil, nil, &buf); err != nil {
		return nil, fmt.Errorf("list buckets: %w", err)
	}
	var result listAllMyBucketsResult
	if err := xml.Unmarshal(buf.Bytes(), &result); err != nil {
		return nil, fmt.Errorf("list buckets: parse response: %w", err)
	}
	buckets := make([]BucketInfo, 0, len(result.Buckets))
	for _, b := range result.Buckets {
		buckets = append(buckets, BucketInfo(b))
	}
	return buckets, nil
}

func (c *Client) ListObjectsV2(ctx context.Context, bucket, prefix, delimiter, continuationToken string) (*ListObjectsV2Output, error) {
	query := url.Values{}
	query.Set("list-type", "2")
	if prefix != "" {
		query.Set("prefix", prefix)
	}
	if delimiter != "" {
		query.Set("delimiter", delimiter)
	}
	if continuationToken != "" {
		query.Set("continuation-token", continuationToken)
	}
	var buf bytes.Buffer
	if err := c.doAndCheck(ctx, "GET", objectPath(bucket, ""), query, nil, nil, nil, &buf); err != nil {
		return nil, fmt.Errorf("list objects in bucket %q: %w", bucket, err)
	}
	var result listBucketResult
	if err := xml.Unmarshal(buf.Bytes(), &result); err != nil {
		return nil, fmt.Errorf("list objects in bucket %q: parse response: %w", bucket, err)
	}
	output := &ListObjectsV2Output{
		IsTruncated:           result.IsTruncated,
		NextContinuationToken: result.NextContinuationToken,
		Objects:               make([]ObjectInfo, 0, len(result.Contents)),
		CommonPrefixes:        make([]string, 0, len(result.CommonPrefixes)),
	}
	for _, obj := range result.Contents {
		output.Objects = append(output.Objects, ObjectInfo(obj))
	}
	for _, p := range result.CommonPrefixes {
		output.CommonPrefixes = append(output.CommonPrefixes, p.Prefix)
	}
	return output, nil
}

func (c *Client) HeadObject(ctx context.Context, bucket, key string) error {
	err := c.doAndCheck(ctx, "HEAD", objectPath(bucket, key), nil, nil, nil, nil, nil)
	if IsNotFound(err) {
		return ErrNotFound
	}
	return err
}

func (c *Client) GetObject(ctx context.Context, bucket, key string) (io.ReadCloser, error) {
	resp, err := c.do(ctx, "GET", objectPath(bucket, key), nil, nil, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("download object %s/%s: %w", bucket, key, err)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
		_ = resp.Body.Close()
		return nil, parseAPIError("GetObject", resp.StatusCode, bodyBytes)
	}
	return resp.Body, nil
}

func (c *Client) PutObject(ctx context.Context, bucket, key string, body io.Reader) error {
	var rewind func() error
	if file, ok := body.(*os.File); ok {
		rewind = func() error {
			_, err := file.Seek(0, io.SeekStart)
			return err
		}
	}
	headers := map[string]string{
		"Content-Type": contentTypeDefault,
	}
	// Wrap the body so the HTTP client does not close it (it closes request
	// bodies that implement io.Closer, which would break retry rewinds)
	err := c.doAndCheck(ctx, "PUT", objectPath(bucket, key), nil, headers, readOnlyReader{body}, rewind, nil)
	if err != nil {
		return fmt.Errorf("upload to %s/%s: %w", bucket, key, err)
	}
	return nil
}

type readOnlyReader struct {
	r io.Reader
}

func (r readOnlyReader) Read(p []byte) (int, error) { return r.r.Read(p) }

func bodyHeaders(contentType string, body []byte) map[string]string {
	sum := md5.Sum(body) //nolint:gosec // Content-MD5 is a storage protocol requirement, not used for security
	return map[string]string{
		"Content-MD5":  base64.StdEncoding.EncodeToString(sum[:]),
		"Content-Type": contentType,
	}
}

func fileContentLength(body io.Reader) (int64, bool) {
	if rr, ok := body.(readOnlyReader); ok {
		body = rr.r
	}
	if file, ok := body.(*os.File); ok {
		if stat, err := file.Stat(); err == nil {
			return stat.Size(), true
		}
	}
	return 0, false
}

func (c *Client) CopyObject(ctx context.Context, src, dst Reference) error {
	headers := map[string]string{
		"x-amz-copy-source": src.Bucket + "/" + escapeKey(src.Key),
	}
	err := c.doAndCheck(ctx, "PUT", objectPath(dst.Bucket, dst.Key), nil, headers, nil, nil, nil)
	if err != nil {
		return fmt.Errorf("copy %s/%s to %s/%s: %w", src.Bucket, src.Key, dst.Bucket, dst.Key, err)
	}
	return nil
}

func (c *Client) DeleteObjects(ctx context.Context, bucket string, keys []string) (int, error) {
	if len(keys) == 0 {
		return 0, nil
	}
	request := deleteRequest{Quiet: true}
	for _, key := range keys {
		request.Objects = append(request.Objects, deleteEntry{Key: key})
	}
	body, err := xml.Marshal(&request)
	if err != nil {
		return 0, fmt.Errorf("delete objects from bucket %q: encode request: %w", bucket, err)
	}
	var buf bytes.Buffer
	query := url.Values{"delete": {""}}
	headers := bodyHeaders(contentTypeXML, body)
	if err := c.doAndCheck(ctx, "POST", objectPath(bucket, ""), query, headers, bytes.NewReader(body), nil, &buf); err != nil {
		return 0, fmt.Errorf("delete objects from bucket %q: %w", bucket, err)
	}
	if len(buf.Bytes()) == 0 {
		// Quiet mode: no per-object result
		return len(keys), nil
	}
	var result deleteResult
	if err := xml.Unmarshal(buf.Bytes(), &result); err != nil {
		return 0, fmt.Errorf("delete objects from bucket %q: parse response: %w", bucket, err)
	}
	if len(result.Errors) > 0 {
		return len(result.Deleted), fmt.Errorf("delete objects from bucket %q: %d object(s) failed: %s: %s", bucket, len(result.Errors), result.Errors[0].Code, result.Errors[0].Message)
	}
	if len(result.Deleted) == 0 {
		// Quiet mode: the server reports no per-object results
		return len(keys), nil
	}
	return len(result.Deleted), nil
}

func (c *Client) CreateBucket(ctx context.Context, bucket, locationConstraint string) error {
	var body []byte
	var headers map[string]string
	if locationConstraint != "" {
		encoded, varErr := xml.Marshal(createBucketConfiguration{LocationConstraint: locationConstraint})
		if varErr != nil {
			return fmt.Errorf("create bucket %q: encode request: %w", bucket, varErr)
		}
		body = encoded
		headers = bodyHeaders(contentTypeXML, body)
	}
	var err error
	if body != nil {
		err = c.doAndCheck(ctx, "PUT", objectPath(bucket, ""), nil, headers, bytes.NewReader(body), nil, nil)
	} else {
		err = c.doAndCheck(ctx, "PUT", objectPath(bucket, ""), nil, headers, nil, nil, nil)
	}
	if err != nil {
		return fmt.Errorf("create bucket %q: %w", bucket, err)
	}
	return nil
}

func (c *Client) DeleteBucket(ctx context.Context, bucket string) error {
	if err := c.doAndCheck(ctx, "DELETE", objectPath(bucket, ""), nil, nil, nil, nil, nil); err != nil {
		return fmt.Errorf("delete bucket %q: %w", bucket, err)
	}
	return nil
}

func (c *Client) PresignGetObject(bucket, key string, expires time.Duration) (string, error) {
	req, err := http.NewRequest("GET", c.endpoint+objectPath(bucket, key), http.NoBody)
	if err != nil {
		return "", fmt.Errorf("presign URL for %s/%s: build request: %w", bucket, key, err)
	}
	req.URL.RawPath = canonicalPath(req)
	// Object Storage presigned downloads declare an unsigned payload, matching the reference implementation
	signedURL, err := c.signer.Presign(req, expires, c.now(), unsignedPayload)
	if err != nil {
		return "", fmt.Errorf("presign URL for %s/%s: %w", bucket, key, err)
	}
	return signedURL, nil
}
