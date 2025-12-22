package auth

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/zalando/go-keyring"
)

type apiClientMocked struct {
	getFails    bool
	getResponse string
}

func (a *apiClientMocked) Do(_ *http.Request) (*http.Response, error) {
	if a.getFails {
		return &http.Response{
			StatusCode: http.StatusNotFound,
		}, fmt.Errorf("not found")
	}
	return &http.Response{
		Status:     "200 OK",
		StatusCode: http.StatusAccepted,
		Body:       io.NopCloser(strings.NewReader(a.getResponse)),
	}, nil
}

func TestParseWellKnownConfig(t *testing.T) {
	tests := []struct {
		name        string
		getFails    bool
		getResponse string
		isValid     bool
		expected    *wellKnownConfig
	}{
		{
			name:        "success",
			getFails:    false,
			getResponse: `{"issuer":"issuer","authorization_endpoint":"auth","token_endpoint":"token"}`,
			isValid:     true,
			expected: &wellKnownConfig{
				Issuer:                "issuer",
				AuthorizationEndpoint: "auth",
				TokenEndpoint:         "token",
			},
		},
		{
			name:        "get_fails",
			getFails:    true,
			getResponse: "",
			isValid:     false,
			expected:    nil,
		},
		{
			name:        "empty_response",
			getFails:    true,
			getResponse: "",
			isValid:     false,
			expected:    nil,
		},
		{
			name:        "missing_issuer",
			getFails:    true,
			getResponse: `{"authorization_endpoint":"auth","token_endpoint":"token"}`,
			isValid:     false,
			expected:    nil,
		},
		{
			name:        "missing_authorization",
			getFails:    true,
			getResponse: `{"issuer":"issuer","token_endpoint":"token"}`,
			isValid:     false,
			expected:    nil,
		},
		{
			name:        "missing_token",
			getFails:    true,
			getResponse: `{"issuer":"issuer","authorization_endpoint":"auth"}`,
			isValid:     false,
			expected:    nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			keyring.MockInit()

			testClient := apiClientMocked{
				tt.getFails,
				tt.getResponse,
			}

			p := print.NewPrinter()

			got, err := parseWellKnownConfiguration(p, &testClient, "", StorageContextCLI)

			if tt.isValid && err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if !tt.isValid && err == nil {
				t.Fatalf("expected error, got none")
			}

			if tt.isValid && !cmp.Equal(*got, *tt.expected) {
				t.Fatalf("expected %v, got %v", tt.expected, got)
			}
		})
	}
}
