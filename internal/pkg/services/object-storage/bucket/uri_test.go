package bucket

import (
	"testing"
)

func TestParseBucketURI(t *testing.T) {
	tests := []struct {
		description string
		uri         string
		expected    Reference
		isValid     bool
	}{
		{
			description: "bucket URI",
			uri:         "s3://my-bucket",
			expected:    Reference{Bucket: "my-bucket"},
			isValid:     true,
		},
		{
			description: "bucket URI with trailing slash",
			uri:         "s3://my-bucket/",
			expected:    Reference{Bucket: "my-bucket"},
			isValid:     true,
		},
		{
			description: "object URI",
			uri:         "s3://my-bucket/key/to/object.txt",
			expected:    Reference{Bucket: "my-bucket", Key: "key/to/object.txt"},
			isValid:     true,
		},
		{
			description: "prefix URI with trailing slash",
			uri:         "s3://my-bucket/prefix/",
			expected:    Reference{Bucket: "my-bucket", Key: "prefix/"},
			isValid:     true,
		},
		{
			description: "object URI with leading slash in key",
			uri:         "s3://my-bucket//key",
			expected:    Reference{Bucket: "my-bucket", Key: "key"},
			isValid:     true,
		},
		{
			description: "virtual hosted bucket URL",
			uri:         "https://my-bucket.object.storage.eu01.onstackit.cloud",
			expected: Reference{
				Bucket:   "my-bucket",
				Endpoint: "https://object.storage.eu01.onstackit.cloud",
				Region:   "eu01",
			},
			isValid: true,
		},
		{
			description: "virtual hosted bucket URL with trailing slash",
			uri:         "https://my-bucket.object.storage.eu01.onstackit.cloud/",
			expected: Reference{
				Bucket:   "my-bucket",
				Endpoint: "https://object.storage.eu01.onstackit.cloud",
				Region:   "eu01",
			},
			isValid: true,
		},
		{
			description: "virtual hosted object URL",
			uri:         "https://my-bucket.object.storage.eu01.onstackit.cloud/dir/file.txt",
			expected: Reference{
				Bucket:   "my-bucket",
				Key:      "dir/file.txt",
				Endpoint: "https://object.storage.eu01.onstackit.cloud",
				Region:   "eu01",
			},
			isValid: true,
		},
		{
			description: "virtual hosted bucket URL with dotted bucket name",
			uri:         "https://my.bucket.name.object.storage.de-fra.onstackit.cloud/key",
			expected: Reference{
				Bucket:   "my.bucket.name",
				Key:      "key",
				Endpoint: "https://object.storage.de-fra.onstackit.cloud",
				Region:   "de-fra",
			},
			isValid: true,
		},
		{
			description: "path style bucket URL",
			uri:         "https://object.storage.eu01.onstackit.cloud/my-bucket",
			expected: Reference{
				Bucket:   "my-bucket",
				Endpoint: "https://object.storage.eu01.onstackit.cloud",
				Region:   "eu01",
			},
			isValid: true,
		},
		{
			description: "path style object URL",
			uri:         "https://object.storage.eu01.onstackit.cloud/my-bucket/dir/file.txt",
			expected: Reference{
				Bucket:   "my-bucket",
				Key:      "dir/file.txt",
				Endpoint: "https://object.storage.eu01.onstackit.cloud",
				Region:   "eu01",
			},
			isValid: true,
		},
		{
			description: "wrong scheme",
			uri:         "ftp://my-bucket/key",
			isValid:     false,
		},
		{
			description: "unrecognized host",
			uri:         "https://example.com/my-bucket",
			isValid:     false,
		},
		{
			description: "path style URL without bucket",
			uri:         "https://object.storage.eu01.onstackit.cloud",
			isValid:     false,
		},
		{
			description: "missing bucket",
			uri:         "s3://",
			isValid:     false,
		},
		{
			description: "empty bucket",
			uri:         "s3:///key",
			isValid:     false,
		},
		{
			description: "invalid bucket name uppercase",
			uri:         "s3://My-Bucket",
			isValid:     false,
		},
		{
			description: "invalid bucket name too short",
			uri:         "s3://ab",
			isValid:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			result, err := ParseBucketURI(tt.uri)
			if (err != nil) != !tt.isValid {
				t.Fatalf("ParseBucketURI(%q) error = %v, wantErr %v", tt.uri, err, !tt.isValid)
			}
			if tt.isValid && result != tt.expected {
				t.Fatalf("ParseBucketURI(%q) = %+v, want %+v", tt.uri, result, tt.expected)
			}
		})
	}
}

func TestReferenceString(t *testing.T) {
	tests := []struct {
		description string
		reference   Reference
		expected    string
	}{
		{
			description: "bucket",
			reference:   Reference{Bucket: "my-bucket"},
			expected:    "s3://my-bucket",
		},
		{
			description: "object",
			reference:   Reference{Bucket: "my-bucket", Key: "key/to/object.txt"},
			expected:    "s3://my-bucket/key/to/object.txt",
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			if result := tt.reference.String(); result != tt.expected {
				t.Fatalf("String() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestIsBucketURI(t *testing.T) {
	tests := []struct {
		description string
		value       string
		expected    bool
	}{
		{description: "valid s3 uri", value: "s3://bucket", expected: true},
		{description: "local path", value: "./local/file.txt", expected: false},
		{description: "empty", value: "", expected: false},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			if result := IsBucketURI(tt.value); result != tt.expected {
				t.Fatalf("IsBucketURI(%q) = %v, want %v", tt.value, result, tt.expected)
			}
		})
	}
}

func TestValidateBucketName(t *testing.T) {
	tests := []struct {
		description string
		name        string
		isValid     bool
	}{
		{description: "valid", name: "my-bucket", isValid: true},
		{description: "valid with dots", name: "my.bucket.name", isValid: true},
		{description: "valid with digits", name: "bucket-123", isValid: true},
		{description: "too short", name: "ab", isValid: false},
		{description: "too long", name: "a234567890123456789012345678901234567890123456789012345678901234", isValid: false},
		{description: "uppercase", name: "My-Bucket", isValid: false},
		{description: "starts with hyphen", name: "-bucket", isValid: false},
		{description: "ends with dot", name: "bucket.", isValid: false},
		{description: "consecutive dots", name: "bu..cket", isValid: false},
		{description: "dot and hyphen adjacent", name: "buc-ket.", isValid: false},
		{description: "looks like ip", name: "192.168.1.1", isValid: false},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			err := ValidateBucketName(tt.name)
			if (err != nil) != !tt.isValid {
				t.Fatalf("ValidateBucketName(%q) error = %v, wantErr %v", tt.name, err, !tt.isValid)
			}
		})
	}
}
