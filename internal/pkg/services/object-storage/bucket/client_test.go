package bucket

import (
	"context"
	"crypto/md5" //nolint:gosec // test helper for Content-MD5, not used for security
	"encoding/base64"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

func newTestClient(t *testing.T, handler http.HandlerFunc) *Client {
	t.Helper()
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)
	client, err := NewClient(ClientConfig{
		Region:          "eu01",
		Endpoint:        server.URL,
		AccessKeyId:     "AKIDEXAMPLE",
		SecretAccessKey: "secret",
	})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	return client
}

func TestNewClientRequiresCredentials(t *testing.T) {
	_, err := NewClient(ClientConfig{Region: "eu01"})
	if err == nil {
		t.Fatal("expected an error when credentials are missing")
	}
}

func TestListBuckets(t *testing.T) {
	var gotAuth, gotPayloadHash string
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		gotPayloadHash = r.Header.Get("X-Amz-Content-Sha256")
		_, _ = w.Write([]byte(`<ListAllMyBucketsResult><Buckets>` +
			`<Bucket><Name>bucket-a</Name><CreationDate>2024-05-01T12:00:00Z</CreationDate></Bucket>` +
			`<Bucket><Name>bucket-b</Name><CreationDate>2024-06-01T12:00:00Z</CreationDate></Bucket>` +
			`</Buckets></ListAllMyBucketsResult>`))
	})

	buckets, err := client.ListBuckets(context.Background())
	if err != nil {
		t.Fatalf("ListBuckets: %v", err)
	}
	if len(buckets) != 2 {
		t.Fatalf("ListBuckets returned %d buckets, want 2", len(buckets))
	}
	if buckets[0].Name != "bucket-a" || buckets[1].Name != "bucket-b" {
		t.Fatalf("unexpected buckets: %+v", buckets)
	}
	if !strings.HasPrefix(gotAuth, "AWS4-HMAC-SHA256 Credential=AKIDEXAMPLE/") {
		t.Fatalf("unexpected Authorization header: %q", gotAuth)
	}
	if gotPayloadHash != emptyPayloadHash {
		t.Fatalf("X-Amz-Content-Sha256 = %q, want %q", gotPayloadHash, emptyPayloadHash)
	}
}

func TestPutObjectSendsBodyAndPayloadHash(t *testing.T) {
	var gotPayloadHash string
	var gotBody []byte
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPayloadHash = r.Header.Get("X-Amz-Content-Sha256")
		gotBody, _ = io.ReadAll(r.Body)
		w.WriteHeader(http.StatusOK)
	})

	err := client.PutObject(context.Background(), "my-bucket", "dir/file.txt", strings.NewReader("hello"))
	if err != nil {
		t.Fatalf("PutObject: %v", err)
	}
	if gotPayloadHash != unsignedPayload {
		t.Fatalf("X-Amz-Content-Sha256 = %q, want %q", gotPayloadHash, unsignedPayload)
	}
	if string(gotBody) != "hello" {
		t.Fatalf("server received body %q, want %q", gotBody, "hello")
	}
}

func TestDeleteObjectsSendsContentMD5AndCountsQuietResponse(t *testing.T) {
	var gotMD5 string
	var gotBody []byte
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMD5 = r.Header.Get("Content-MD5")
		gotBody, _ = io.ReadAll(r.Body)
		// Quiet mode response without per-object results
		_, _ = w.Write([]byte(`<DeleteResult></DeleteResult>`))
	})

	deleted, err := client.DeleteObjects(context.Background(), "my-bucket", []string{"a.txt", "b.txt"})
	if err != nil {
		t.Fatalf("DeleteObjects: %v", err)
	}
	if deleted != 2 {
		t.Fatalf("deleted = %d, want 2", deleted)
	}
	wantMD5 := base64.StdEncoding.EncodeToString(md5sum(gotBody))
	if gotMD5 != wantMD5 {
		t.Fatalf("Content-MD5 = %q, want %q for body %q", gotMD5, wantMD5, gotBody)
	}
}

func TestDeleteObjectsReportsPerObjectResults(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`<DeleteResult>` +
			`<Deleted><Key>a.txt</Key></Deleted>` +
			`<Error><Key>b.txt</Key><Code>AccessDenied</Code><Message>denied</Message></Error>` +
			`</DeleteResult>`))
	})

	deleted, err := client.DeleteObjects(context.Background(), "my-bucket", []string{"a.txt", "b.txt"})
	if err == nil {
		t.Fatal("expected an error when a delete fails")
	}
	if deleted != 1 {
		t.Fatalf("deleted = %d, want 1", deleted)
	}
	if !strings.Contains(err.Error(), "AccessDenied") {
		t.Fatalf("error = %v, want it to mention AccessDenied", err)
	}
}

func TestCreateBucketSendsLocationConstraint(t *testing.T) {
	var gotMethod, gotPath string
	var gotMD5, gotContentType string
	var gotBody []byte
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		gotMD5 = r.Header.Get("Content-MD5")
		gotContentType = r.Header.Get("Content-Type")
		gotBody, _ = io.ReadAll(r.Body)
		w.WriteHeader(http.StatusOK)
	})

	err := client.CreateBucket(context.Background(), "my-bucket", "eu02")
	if err != nil {
		t.Fatalf("CreateBucket: %v", err)
	}
	if gotMethod != http.MethodPut || gotPath != "/my-bucket" {
		t.Fatalf("request = %s %s, want PUT /my-bucket", gotMethod, gotPath)
	}
	wantBody := `<CreateBucketConfiguration><LocationConstraint>eu02</LocationConstraint></CreateBucketConfiguration>`
	if string(gotBody) != wantBody {
		t.Fatalf("body = %q, want %q", gotBody, wantBody)
	}
	if gotContentType != "application/xml" {
		t.Fatalf("Content-Type = %q, want %q", gotContentType, "application/xml")
	}
	wantMD5 := base64.StdEncoding.EncodeToString(md5sum(gotBody))
	if gotMD5 != wantMD5 {
		t.Fatalf("Content-MD5 = %q, want %q for body %q", gotMD5, wantMD5, gotBody)
	}
}

func TestCreateBucketWithoutConstraintSendsEmptyBody(t *testing.T) {
	var gotMethod, gotPath string
	var gotMD5, gotContentType string
	var gotBody []byte
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		gotMD5 = r.Header.Get("Content-MD5")
		gotContentType = r.Header.Get("Content-Type")
		gotBody, _ = io.ReadAll(r.Body)
		w.WriteHeader(http.StatusOK)
	})

	err := client.CreateBucket(context.Background(), "my-bucket", "")
	if err != nil {
		t.Fatalf("CreateBucket: %v", err)
	}
	if gotMethod != http.MethodPut || gotPath != "/my-bucket" {
		t.Fatalf("request = %s %s, want PUT /my-bucket", gotMethod, gotPath)
	}
	if len(gotBody) != 0 {
		t.Fatalf("body = %q, want empty", gotBody)
	}
	if gotMD5 != "" {
		t.Fatalf("Content-MD5 = %q, want empty", gotMD5)
	}
	if gotContentType != "" {
		t.Fatalf("Content-Type = %q, want empty", gotContentType)
	}
}

func TestRetryOnServerErrors(t *testing.T) {
	var attempts atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		if attempts.Add(1) <= 2 {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		_, _ = w.Write([]byte(`<ListAllMyBucketsResult><Buckets></Buckets></ListAllMyBucketsResult>`))
	}))
	t.Cleanup(server.Close)
	client, err := NewClient(ClientConfig{
		Region:          "eu01",
		Endpoint:        server.URL,
		AccessKeyId:     "AKIDEXAMPLE",
		SecretAccessKey: "secret",
	})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	client.retryDelay = time.Millisecond

	buckets, err := client.ListBuckets(context.Background())
	if err != nil {
		t.Fatalf("ListBuckets: %v", err)
	}
	if len(buckets) != 0 {
		t.Fatalf("ListBuckets returned %d buckets, want 0", len(buckets))
	}
	if got := attempts.Load(); got != 3 {
		t.Fatalf("attempts = %d, want 3", got)
	}
}

func TestNoRetryOnClientErrors(t *testing.T) {
	var attempts atomic.Int32
	client := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
		attempts.Add(1)
		w.WriteHeader(http.StatusForbidden)
		_, _ = w.Write([]byte(`<Error><Code>AccessDenied</Code><Message>Access Denied</Message></Error>`))
	})

	_, err := client.ListBuckets(context.Background())
	if err == nil {
		t.Fatal("expected an error on 403")
	}
	if got := attempts.Load(); got != 1 {
		t.Fatalf("attempts = %d, want 1", got)
	}
	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("error type = %T, want *APIError", err)
	}
	if apiErr.Code != "AccessDenied" || apiErr.StatusCode != http.StatusForbidden {
		t.Fatalf("APIError = %+v, want 403 AccessDenied", apiErr)
	}
}

func md5sum(data []byte) []byte {
	sum := md5.Sum(data) //nolint:gosec // test helper for Content-MD5, not used for security
	return sum[:]
}
