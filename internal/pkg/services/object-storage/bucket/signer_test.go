package bucket

import (
	"net/http"
	"testing"
	"time"
)

var (
	testTime   = time.Date(2015, 8, 30, 12, 36, 0, 0, time.UTC)
	testSigner = Signer{ //nolint:gosec // public example credentials from the AWS documentation
		AccessKeyID:     "AKIDEXAMPLE",
		SecretAccessKey: "wJalrXUtnFEMI/K7MDENG+bPxRfiCYEXAMPLEKEY",
		Region:          "us-east-1",
		Service:         "service",
	}
)

func TestSign(t *testing.T) {
	// Reference computed with botocore (pinned timestamp): GET
	// https://example.amazonaws.com/?Host=example.amazonaws.com, service "service",
	// region us-east-1, date 20150830T123600Z
	req, err := http.NewRequest("GET", "https://example.amazonaws.com/?Host=example.amazonaws.com", http.NoBody)
	if err != nil {
		t.Fatalf("build request: %v", err)
	}
	if err := testSigner.Sign(req, testTime); err != nil {
		t.Fatalf("Sign: %v", err)
	}
	wantAuth := "AWS4-HMAC-SHA256 Credential=AKIDEXAMPLE/20150830/us-east-1/service/aws4_request, SignedHeaders=host;x-amz-date, Signature=ca564d2dbd30b429f5d637740c18b0039343ae449cf2bd0ba0b947298053f3c8"
	if got := req.Header.Get("Authorization"); got != wantAuth {
		t.Fatalf("Authorization = %q, want %q", got, wantAuth)
	}
	if got := req.Header.Get("X-Amz-Date"); got != "20150830T123600Z" {
		t.Fatalf("X-Amz-Date = %q, want 20150830T123600Z", got)
	}
}

func TestSignSortsQueryParameters(t *testing.T) {
	// Reference computed with botocore (pinned timestamp): the request query is
	// out of order and must be canonicalized
	req, err := http.NewRequest("GET", "https://example.amazonaws.com/?Param2=2&Param1=1", http.NoBody)
	if err != nil {
		t.Fatalf("build request: %v", err)
	}
	if err := testSigner.Sign(req, testTime); err != nil {
		t.Fatalf("Sign: %v", err)
	}
	wantAuth := "AWS4-HMAC-SHA256 Credential=AKIDEXAMPLE/20150830/us-east-1/service/aws4_request, SignedHeaders=host;x-amz-date, Signature=afce2c564917f25d63492b2cf7c08495be06c0a160aba74e53bacf4caffe854c"
	if got := req.Header.Get("Authorization"); got != wantAuth {
		t.Fatalf("Authorization = %q, want %q", got, wantAuth)
	}
}

func TestPresign(t *testing.T) {
	// Reference computed with botocore's S3SigV4QueryAuth (pinned timestamp):
	// presigned download with UNSIGNED-PAYLOAD, signing only the host header
	client := &Client{
		endpoint: "https://object.storage.eu01.onstackit.cloud",
		region:   "eu01",
		signer:   NewSigner("AKID", "SECRET", "eu01"),
		now:      func() time.Time { return time.Date(2024, 5, 1, 12, 0, 0, 0, time.UTC) },
	}
	got, err := client.PresignGetObject("my-bucket", "dir/file.txt", 15*time.Minute)
	if err != nil {
		t.Fatalf("PresignGetObject: %v", err)
	}
	want := "https://object.storage.eu01.onstackit.cloud/my-bucket/dir/file.txt" +
		"?X-Amz-Algorithm=AWS4-HMAC-SHA256" +
		"&X-Amz-Credential=AKID%2F20240501%2Feu01%2Fs3%2Faws4_request" +
		"&X-Amz-Date=20240501T120000Z" +
		"&X-Amz-Expires=900" +
		"&X-Amz-SignedHeaders=host" +
		"&X-Amz-Signature=c58f23addabe8c12ea43bc65ef7e8e51c8f5bae00e0db74d7adb283448f39d65"
	if got != want {
		t.Fatalf("presigned URL = %s, want %s", got, want)
	}
}

func TestPresignEncodesKeyFully(t *testing.T) {
	// Vector verified byte-for-byte against the AWS SDK presigner: a key with
	// spaces, & and parentheses, UNSIGNED-PAYLOAD, and an extra query parameter
	req, err := http.NewRequest("GET", "https://object.storage.eu01.onstackit.cloud/my-bucket/dir/file with space & (parens).txt", http.NoBody)
	if err != nil {
		t.Fatalf("build request: %v", err)
	}
	query := req.URL.Query()
	query.Set("x-id", "GetObject")
	req.URL.RawQuery = query.Encode()
	req.URL.RawPath = canonicalPath(req)

	signer := NewSigner("AKIDEXAMPLE", "wJalrXUtnFEMI/K7MDENG+bPxRfiCYEXAMPLEKEY", "eu01")
	got, err := signer.Presign(req, 15*time.Minute, time.Date(2026, 11, 16, 17, 6, 13, 0, time.UTC), unsignedPayload)
	if err != nil {
		t.Fatalf("Presign: %v", err)
	}
	want := "https://object.storage.eu01.onstackit.cloud/my-bucket/dir/file%20with%20space%20%26%20%28parens%29.txt" +
		"?X-Amz-Algorithm=AWS4-HMAC-SHA256" +
		"&X-Amz-Credential=AKIDEXAMPLE%2F20261116%2Feu01%2Fs3%2Faws4_request" +
		"&X-Amz-Date=20261116T170613Z" +
		"&X-Amz-Expires=900" +
		"&X-Amz-SignedHeaders=host" +
		"&x-id=GetObject" +
		"&X-Amz-Signature=a9e8de480dde798cb296f1f81433ac4bf6183f03c49711646c0f99ebdcea908c"
	if got != want {
		t.Fatalf("presigned URL = %s, want %s", got, want)
	}
}

func TestCanonicalPath(t *testing.T) {
	tests := []struct {
		description string
		rawURL      string
		expected    string
	}{
		{description: "root", rawURL: "https://example.com", expected: "/"},
		{description: "simple", rawURL: "https://example.com/bucket/key.txt", expected: "/bucket/key.txt"},
		{description: "reserved characters are fully encoded", rawURL: "https://example.com/bucket/a b&c(d).txt", expected: "/bucket/a%20b%26c%28d%29.txt"},
		{description: "slash stays unencoded", rawURL: "https://example.com/bucket/dir/file.txt", expected: "/bucket/dir/file.txt"},
		{description: "unreserved characters stay as-is", rawURL: "https://example.com/bucket/a-b_c.d~e.txt", expected: "/bucket/a-b_c.d~e.txt"},
	}
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			req, err := http.NewRequest("GET", tt.rawURL, http.NoBody)
			if err != nil {
				t.Fatalf("build request: %v", err)
			}
			if got := canonicalPath(req); got != tt.expected {
				t.Fatalf("canonicalPath = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestQueryEncode(t *testing.T) {
	tests := []struct {
		description string
		value       string
		expected    string
	}{
		{description: "unreserved", value: "AZaz09-_.~", expected: "AZaz09-_.~"},
		{description: "space", value: "a b", expected: "a%20b"},
		{description: "ampersand", value: "a&b", expected: "a%26b"},
		{description: "parens", value: "(a)", expected: "%28a%29"},
		{description: "slash is encoded in query values", value: "a/b", expected: "a%2Fb"},
		{description: "plus", value: "a+b", expected: "a%2Bb"},
		{description: "empty", value: "", expected: ""},
	}
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			if got := queryEncode(tt.value); got != tt.expected {
				t.Fatalf("queryEncode(%q) = %q, want %q", tt.value, got, tt.expected)
			}
		})
	}
}

func TestCanonicalQueryString(t *testing.T) {
	req, err := http.NewRequest("GET", "https://example.com/?list-type=2&prefix=dir/sub&prefix=a b&z=1&a=2", http.NoBody)
	if err != nil {
		t.Fatalf("build request: %v", err)
	}
	got := canonicalQueryString(req.URL.Query())
	want := "a=2&list-type=2&prefix=a%20b&prefix=dir%2Fsub&z=1"
	if got != want {
		t.Fatalf("canonicalQueryString = %q, want %q", got, want)
	}
}

func TestPayloadHash(t *testing.T) {
	req, _ := http.NewRequest("GET", "https://example.com/", http.NoBody)
	if got := payloadHash(req); got != emptyPayloadHash {
		t.Fatalf("payloadHash without header = %q, want %q", got, emptyPayloadHash)
	}
	req.Header.Set("X-Amz-Content-Sha256", unsignedPayload)
	if got := payloadHash(req); got != unsignedPayload {
		t.Fatalf("payloadHash with header = %q, want %q", got, unsignedPayload)
	}
}
