package bucket

import (
	"errors"
	"fmt"
	"testing"
)

func TestParseAPIError(t *testing.T) {
	t.Run("valid error body", func(t *testing.T) {
		body := []byte(`<Error><Code>NoSuchKey</Code><Message>The specified key does not exist.</Message><RequestId>ABC123</RequestId></Error>`)
		err := parseAPIError("GetObject", 404, body)
		var apiErr *APIError
		if !errors.As(err, &apiErr) {
			t.Fatalf("error type = %T, want *APIError", err)
		}
		if apiErr.StatusCode != 404 || apiErr.Code != "NoSuchKey" || apiErr.RequestID != "ABC123" {
			t.Fatalf("APIError = %+v", apiErr)
		}
		if want := "Object Storage API error 404: NoSuchKey: The specified key does not exist."; apiErr.Error() != want {
			t.Fatalf("Error() = %q, want %q", apiErr.Error(), want)
		}
	})

	t.Run("invalid error body falls back to generic message", func(t *testing.T) {
		err := parseAPIError("GetObject", 502, []byte("<html>bad gateway</html>"))
		var apiErr *APIError
		if !errors.As(err, &apiErr) {
			t.Fatalf("error type = %T, want *APIError", err)
		}
		if apiErr.Code != "" || apiErr.Message != "GetObject: unexpected response" {
			t.Fatalf("APIError = %+v", apiErr)
		}
	})
}

func TestIsNotFound(t *testing.T) {
	tests := []struct {
		description string
		err         error
		expected    bool
	}{
		{
			description: "NoSuchKey code",
			err:         &APIError{StatusCode: 404, Code: "NoSuchKey"},
			expected:    true,
		},
		{
			description: "NoSuchBucket code",
			err:         &APIError{StatusCode: 404, Code: "NoSuchBucket"},
			expected:    true,
		},
		{
			description: "404 without code",
			err:         &APIError{StatusCode: 404},
			expected:    true,
		},
		{
			description: "AccessDenied is not not-found",
			err:         &APIError{StatusCode: 403, Code: "AccessDenied"},
			expected:    false,
		},
		{
			description: "wrapped APIError",
			err:         fmt.Errorf("check object: %w", &APIError{StatusCode: 404, Code: "NoSuchKey"}),
			expected:    true,
		},
		{
			description: "error mentioning the code but not wrapping it",
			err:         errors.New("outer: Object Storage API error 404: NoSuchKey: The specified key does not exist."),
			expected:    false,
		},
		{
			description: "ErrNotFound sentinel",
			err:         ErrNotFound,
			expected:    true,
		},
		{
			description: "unrelated error",
			err:         errors.New("boom"),
			expected:    false,
		},
		{
			description: "nil",
			err:         nil,
			expected:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			if got := IsNotFound(tt.err); got != tt.expected {
				t.Fatalf("IsNotFound = %v, want %v", got, tt.expected)
			}
		})
	}
}
