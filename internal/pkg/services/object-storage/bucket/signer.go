package bucket

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	signingAlgorithm      = "AWS4-HMAC-SHA256"
	signingService        = "s3"
	terminalType          = "aws4_request"
	credentialScopePrefix = "AWS4"

	emptyPayloadHash = "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"

	unsignedPayload = "UNSIGNED-PAYLOAD"
)

type Signer struct {
	AccessKeyID     string
	SecretAccessKey string
	Region          string
	Service         string
}

func NewSigner(accessKeyID, secretAccessKey, region string) *Signer {
	return &Signer{
		AccessKeyID:     accessKeyID,
		SecretAccessKey: secretAccessKey,
		Region:          region,
		Service:         signingService,
	}
}

func (s *Signer) credentialScope(dateStamp string) string {
	return dateStamp + "/" + s.Region + "/" + s.Service + "/" + terminalType
}

func (s *Signer) signingKey(dateStamp string) []byte {
	kDate := hmacSHA256([]byte(credentialScopePrefix+s.SecretAccessKey), dateStamp)
	kRegion := hmacSHA256(kDate, s.Region)
	kService := hmacSHA256(kRegion, s.Service)
	return hmacSHA256(kService, terminalType)
}

func canonicalPath(req *http.Request) string {
	path := req.URL.Path
	if path == "" {
		return "/"
	}
	segments := strings.Split(path, "/")
	for i, segment := range segments {
		segments[i] = queryEncode(segment)
	}
	return strings.Join(segments, "/")
}

func (s *Signer) Sign(req *http.Request, now time.Time) error {
	now = now.UTC()
	amzDate := now.Format("20060102T150405Z")
	dateStamp := now.Format("20060102")
	req.Header.Set("X-Amz-Date", amzDate)

	canonicalHeaders, signedHeaders := canonicalHeaders(req)
	canonicalRequest := buildCanonicalRequest(
		req.Method,
		canonicalPath(req),
		canonicalQueryString(req.URL.Query()),
		canonicalHeaders,
		signedHeaders,
		payloadHash(req),
	)
	signature := s.signature(canonicalRequest, amzDate, dateStamp)
	req.Header.Set("Authorization", fmt.Sprintf("%s Credential=%s/%s, SignedHeaders=%s, Signature=%s",
		signingAlgorithm, s.AccessKeyID, s.credentialScope(dateStamp), signedHeaders, signature))
	return nil
}

func (s *Signer) Presign(req *http.Request, expires time.Duration, now time.Time, payloadHash string) (string, error) {
	now = now.UTC()
	amzDate := now.Format("20060102T150405Z")
	dateStamp := now.Format("20060102")

	query := req.URL.Query()
	query.Set("X-Amz-Algorithm", signingAlgorithm)
	query.Set("X-Amz-Credential", s.AccessKeyID+"/"+s.credentialScope(dateStamp))
	query.Set("X-Amz-Date", amzDate)
	query.Set("X-Amz-Expires", strconv.FormatInt(int64(expires.Seconds()), 10))
	query.Set("X-Amz-SignedHeaders", "host")

	host := hostOf(req)
	canonicalRequest := buildCanonicalRequest(
		req.Method,
		canonicalPath(req),
		canonicalQueryString(query),
		"host:"+host+"\n",
		"host",
		payloadHash,
	)
	signature := s.signature(canonicalRequest, amzDate, dateStamp)
	// X-Amz-Signature is not part of the canonical query string; append it last
	// to match the SDK wire format
	req.URL.RawQuery = canonicalQueryString(query) + "&X-Amz-Signature=" + signature
	return req.URL.String(), nil
}

func (s *Signer) signature(canonicalRequest, amzDate, dateStamp string) string {
	stringToSign := strings.Join([]string{
		signingAlgorithm,
		amzDate,
		s.credentialScope(dateStamp),
		sha256Hex([]byte(canonicalRequest)),
	}, "\n")
	return hex.EncodeToString(hmacSHA256(s.signingKey(dateStamp), stringToSign))
}

func buildCanonicalRequest(method, canonicalURI, canonicalQuery, canonicalHeaders, signedHeaders, payloadHash string) string {
	return strings.Join([]string{method, canonicalURI, canonicalQuery, canonicalHeaders, signedHeaders, payloadHash}, "\n")
}

func canonicalHeaders(req *http.Request) (headersBlock, signedHeaders string) {
	names := make(map[string]string)
	if host := hostOf(req); host != "" {
		names["host"] = host
	}
	for name, values := range req.Header {
		lower := strings.ToLower(name)
		if !strings.HasPrefix(lower, "x-amz-") && lower != "content-type" {
			continue
		}
		names[lower] = strings.Join(values, ",")
	}
	sorted := make([]string, 0, len(names))
	for name := range names {
		sorted = append(sorted, name)
	}
	sort.Strings(sorted)

	var builder strings.Builder
	for _, name := range sorted {
		builder.WriteString(name)
		builder.WriteByte(':')
		builder.WriteString(strings.TrimSpace(names[name]))
		builder.WriteByte('\n')
	}
	return builder.String(), strings.Join(sorted, ";")
}

func canonicalQueryString(values url.Values) string {
	type kv struct{ key, value string }
	pairs := make([]kv, 0, len(values))
	for key, vals := range values {
		for _, value := range vals {
			pairs = append(pairs, kv{key: key, value: value})
		}
	}
	sort.Slice(pairs, func(i, j int) bool {
		if pairs[i].key != pairs[j].key {
			return pairs[i].key < pairs[j].key
		}
		return pairs[i].value < pairs[j].value
	})
	parts := make([]string, 0, len(pairs))
	for _, pair := range pairs {
		parts = append(parts, queryEncode(pair.key)+"="+queryEncode(pair.value))
	}
	return strings.Join(parts, "&")
}

func payloadHash(req *http.Request) string {
	if v := req.Header.Get("X-Amz-Content-Sha256"); v != "" {
		return v
	}
	return emptyPayloadHash
}

func hostOf(req *http.Request) string {
	if req.Host != "" {
		return req.Host
	}
	return req.URL.Host
}

func queryEncode(s string) string {
	var builder strings.Builder
	for i := 0; i < len(s); i++ {
		c := s[i]
		switch {
		case (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z') || (c >= '0' && c <= '9'):
			builder.WriteByte(c)
		case c == '-' || c == '_' || c == '.' || c == '~':
			builder.WriteByte(c)
		default:
			fmt.Fprintf(&builder, "%%%02X", c)
		}
	}
	return builder.String()
}

func escapeKey(key string) string {
	segments := strings.Split(key, "/")
	for i, segment := range segments {
		segments[i] = queryEncode(segment)
	}
	return strings.Join(segments, "/")
}

func hmacSHA256(key []byte, data string) []byte {
	mac := hmac.New(sha256.New, key)
	mac.Write([]byte(data))
	return mac.Sum(nil)
}

func sha256Hex(data []byte) string {
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}
