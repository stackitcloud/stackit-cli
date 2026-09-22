package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/spf13/viper"
	"github.com/zalando/go-keyring"

	"github.com/stackitcloud/stackit-cli/internal/pkg/config"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testparams"
)

func TestGetWellKnownConfig(t *testing.T) {
	tests := []struct {
		name              string
		idpCustomEndpoint string
		allowedUrlDomain  string
		isValid           bool
		expected          string
	}{
		{
			name:              "custom endpoint specified",
			idpCustomEndpoint: "https://example.stackit.cloud",
			allowedUrlDomain:  "stackit.cloud",
			isValid:           true,
			expected:          "https://example.stackit.cloud",
		},
		{
			name:              "custom endpoint outside STACKIT",
			idpCustomEndpoint: "https://www.very-suspicious-website.com/",
			allowedUrlDomain:  "stackit.cloud",
			isValid:           false,
		},
		{
			name:              "non-STACKIT custom endpoint invalid",
			idpCustomEndpoint: "https://www.very-suspicious-website.com/",
			allowedUrlDomain:  "stackit.cloud",
			isValid:           false,
		},
		{
			name:              "non-STACKIT custom endpoint valid",
			idpCustomEndpoint: "https://www.test.example.com/",
			allowedUrlDomain:  "example.com",
			isValid:           true,
			expected:          "https://www.test.example.com/",
		},
		{
			name:              "every URL valid",
			idpCustomEndpoint: "https://www.test.example.com/",
			allowedUrlDomain:  "",
			isValid:           true,
			expected:          "https://www.test.example.com/",
		},
		{
			name:              "custom endpoint not specified",
			idpCustomEndpoint: "",
			allowedUrlDomain:  "",
			isValid:           true,
			expected:          defaultWellKnownConfig,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			viper.Reset()
			viper.Set(config.IdentityProviderCustomWellKnownConfigurationKey, tt.idpCustomEndpoint)
			viper.Set(config.AllowedUrlDomainKey, tt.allowedUrlDomain)

			got, err := getIDPWellKnownConfigURL()

			if tt.isValid && err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if !tt.isValid && err == nil {
				t.Fatalf("expected error, got none")
			}

			if got != tt.expected {
				t.Fatalf("expected idp endpoint %q, got %q", tt.expected, got)
			}
		})
	}
}

func TestGetIDPClientID(t *testing.T) {
	tests := []struct {
		name              string
		idpCustomClientID string
		isValid           bool
		expected          string
	}{
		{
			name:              "custom client ID specified",
			idpCustomClientID: "custom-client-id",
			isValid:           true,
			expected:          "custom-client-id",
		},
		{
			name:              "custom client ID not specified",
			idpCustomClientID: "",
			isValid:           true,
			expected:          defaultCLIClientID,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			viper.Reset()
			viper.Set(config.IdentityProviderCustomClientIdKey, tt.idpCustomClientID)

			got, err := getIDPClientID()

			if tt.isValid && err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if !tt.isValid && err == nil {
				t.Fatalf("expected error, got none")
			}

			if got != tt.expected {
				t.Fatalf("expected idp client ID %q, got %q", tt.expected, got)
			}
		})
	}
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
			getResponse: `{"issuer":"https://issuer.stackit.cloud/endpoint","authorization_endpoint":"https://auth.stackit.cloud/enpoint","token_endpoint":"https://token.stackit.cloud/endpoint"}`,
			isValid:     true,
			expected: &wellKnownConfig{ //nolint:gosec // just a testcase; no credentials
				Issuer:                "https://issuer.stackit.cloud/endpoint",
				AuthorizationEndpoint: "https://auth.stackit.cloud/enpoint",
				TokenEndpoint:         "https://token.stackit.cloud/endpoint",
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
			getResponse: `{"authorization_endpoint":"https://auth.stackit.cloud/enpoint","token_endpoint":"https://token.stackit.cloud/endpoint"}`,
			isValid:     false,
			expected:    nil,
		},
		{
			name:        "missing_authorization",
			getFails:    true,
			getResponse: `{"issuer":"https://issuer.stackit.cloud/endpoint","token_endpoint":"https://token.stackit.cloud/endpoint"}`,
			isValid:     false,
			expected:    nil,
		},
		{
			name:        "missing_token",
			getFails:    true,
			getResponse: `{"issuer":"https://issuer.stackit.cloud/endpoint","authorization_endpoint":"https://auth.stackit.cloud/enpoint"}`,
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

			got, err := parseWellKnownConfiguration(&testClient, "")

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

func TestParseWellKnownConfigWithoutStorage(t *testing.T) {
	keyring.MockInit()
	const existingEndpoint = "https://existing.stackit.cloud/oauth"
	if err := SetAuthField(IDP_TOKEN_ENDPOINT, existingEndpoint); err != nil {
		t.Fatalf("Set existing token endpoint: %v", err)
	}
	testClient := apiClientMocked{
		false,
		`{"issuer":"https://issuer.stackit.cloud/endpoint","authorization_endpoint":"https://auth.stackit.cloud/endpoint","token_endpoint":"https://token.stackit.cloud/endpoint"}`,
	}

	wellKnownConfig, err := parseWellKnownConfigurationWithStorage(&testClient, "", false)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if wellKnownConfig.TokenEndpoint != "https://token.stackit.cloud/endpoint" {
		t.Fatalf("Unexpected token endpoint %q", wellKnownConfig.TokenEndpoint)
	}

	storedTokenEndpoint, err := GetAuthField(IDP_TOKEN_ENDPOINT)
	if err != nil {
		t.Fatalf("Get stored token endpoint: %v", err)
	}
	if storedTokenEndpoint != existingEndpoint {
		t.Fatalf("Expected token endpoint not to be changed, got %q", storedTokenEndpoint)
	}
}

func TestGetIDPTokenEndpointStatelessIgnoresStorage(t *testing.T) {
	viper.Reset()
	t.Cleanup(viper.Reset)
	keyring.MockInit()

	const storedEndpoint = "https://stored.stackit.cloud/oauth"
	if err := SetAuthField(IDP_TOKEN_ENDPOINT, storedEndpoint); err != nil {
		t.Fatalf("Set stored token endpoint: %v", err)
	}

	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"issuer":"https://issuer.stackit.cloud/endpoint","authorization_endpoint":"https://auth.stackit.cloud/endpoint","token_endpoint":"https://discovered.stackit.cloud/oauth"}`))
	}))
	defer server.Close()

	originalTransport := http.DefaultTransport
	http.DefaultTransport = server.Client().Transport
	t.Cleanup(func() {
		http.DefaultTransport = originalTransport
	})
	viper.Set(config.IdentityProviderCustomWellKnownConfigurationKey, server.URL)
	viper.Set(config.AllowedUrlDomainKey, "")

	params := testparams.NewTestParams()
	params.Printer.AssumeYes = true
	actual, err := GetIDPTokenEndpointStateless(params.Printer)
	if err != nil {
		t.Fatalf("Get stateless token endpoint: %v", err)
	}
	const discoveredEndpoint = "https://discovered.stackit.cloud/oauth"
	if actual != discoveredEndpoint {
		t.Fatalf("Expected discovered endpoint %q, got %q", discoveredEndpoint, actual)
	}

	actualStoredEndpoint, err := GetAuthField(IDP_TOKEN_ENDPOINT)
	if err != nil {
		t.Fatalf("Get stored token endpoint: %v", err)
	}
	if actualStoredEndpoint != storedEndpoint {
		t.Fatalf("Expected stored endpoint to remain %q, got %q", storedEndpoint, actualStoredEndpoint)
	}
}
