package auth_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stackitcloud/stackit-cli/internal/pkg/auth"
)

func boolPtr(v bool) *bool {
	return &v
}

func TestIsEnabled(t *testing.T) {
	tests := []struct {
		value    string
		expected bool
	}{
		{"1", true},
		{"0", false},
		{"", false},
		{"true", false},
		{"yes", false},
		{"random", false},
	}
	for _, tt := range tests {
		t.Run(tt.value, func(t *testing.T) {
			t.Setenv(auth.EnvUseOIDC, tt.value)
			got := auth.IsOIDCEnabled()
			if got != tt.expected {
				t.Errorf("IsOIDCEnabled() = %v, want %v (env=%q)", got, tt.expected, tt.value)
			}
		})
	}
}

func TestIsEnabled_Unset(t *testing.T) {
	// When the env var is not set at all IsEnabled must return false
	t.Setenv(auth.EnvUseOIDC, "")
	if auth.IsOIDCEnabled() {
		t.Error("IsOIDCEnabled() = true, want false when env var is empty")
	}
}

func TestIsOIDCEnabledWithOverride(t *testing.T) {
	tests := []struct {
		description string
		envUseOIDC  string
		override    *bool
		expected    bool
	}{
		{
			description: "uses env when override is nil",
			envUseOIDC:  "1",
			override:    nil,
			expected:    true,
		},
		{
			description: "override true wins over env false",
			envUseOIDC:  "0",
			override:    boolPtr(true),
			expected:    true,
		},
		{
			description: "override false wins over env true",
			envUseOIDC:  "1",
			override:    boolPtr(false),
			expected:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			t.Setenv(auth.EnvUseOIDC, tt.envUseOIDC)

			got := auth.IsOIDCEnabledWithOverride(tt.override)
			if got != tt.expected {
				t.Errorf("IsOIDCEnabledWithOverride() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestServiceAccountEmail(t *testing.T) {
	const want = "ci@sa.stackit.cloud"
	t.Setenv(auth.EnvServiceAccountEmail, want)
	if got := auth.OIDCServiceAccountEmail(); got != want {
		t.Errorf("OIDCServiceAccountEmail() = %q, want %q", got, want)
	}
}

func TestTokenFunc_StaticToken(t *testing.T) {
	const want = "my-static-oidc-token"
	t.Setenv(auth.EnvServiceAccountFederatedToken, want)
	// ensure GitHub / Azure vars are absent so we hit the static path first
	t.Setenv(auth.EnvGitHubRequestURL, "")
	t.Setenv(auth.EnvGitHubRequestToken, "")
	t.Setenv(auth.EnvAzureOIDCRequestURI, "")
	t.Setenv(auth.EnvAzureAccessToken, "")
	t.Setenv(auth.EnvFederatedTokenFile, "")

	fn, err := auth.OIDCTokenFunc()
	if err != nil {
		t.Fatalf("OIDCTokenFunc() unexpected error: %v", err)
	}
	got, err := fn(context.Background())
	if err != nil {
		t.Fatalf("fn() unexpected error: %v", err)
	}
	if got != want {
		t.Errorf("fn() = %q, want %q", got, want)
	}
}

func TestTokenFunc_GitHubActions(t *testing.T) {
	// Spin up a fake GitHub OIDC endpoint that matches the SDK's expected format.
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{"value": "gh-oidc-token"})
	}))
	defer srv.Close()

	t.Setenv(auth.EnvServiceAccountFederatedToken, "")
	t.Setenv(auth.EnvFederatedTokenFile, "")
	t.Setenv(auth.EnvGitHubRequestURL, srv.URL)
	t.Setenv(auth.EnvGitHubRequestToken, "gh-bearer-token")
	t.Setenv(auth.EnvAzureOIDCRequestURI, "")
	t.Setenv(auth.EnvAzureAccessToken, "")

	fn, err := auth.OIDCTokenFunc()
	if err != nil {
		t.Fatalf("OIDCTokenFunc() unexpected error: %v", err)
	}
	got, err := fn(context.Background())
	if err != nil {
		t.Fatalf("fn() unexpected error: %v", err)
	}
	if got != "gh-oidc-token" {
		t.Errorf("fn() = %q, want %q", got, "gh-oidc-token")
	}
}

func TestTokenFunc_AzureDevOps(t *testing.T) {
	// Spin up a fake Azure DevOps OIDC endpoint that matches the SDK's expected format.
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		tok := "ado-oidc-token"
		_ = json.NewEncoder(w).Encode(map[string]*string{"oidcToken": &tok})
	}))
	defer srv.Close()

	t.Setenv(auth.EnvServiceAccountFederatedToken, "")
	t.Setenv(auth.EnvFederatedTokenFile, "")
	t.Setenv(auth.EnvGitHubRequestURL, "")
	t.Setenv(auth.EnvGitHubRequestToken, "")
	t.Setenv(auth.EnvAzureOIDCRequestURI, srv.URL)
	t.Setenv(auth.EnvAzureAccessToken, "ado-access-token")

	fn, err := auth.OIDCTokenFunc()
	if err != nil {
		t.Fatalf("OIDCTokenFunc() unexpected error: %v", err)
	}
	got, err := fn(context.Background())
	if err != nil {
		t.Fatalf("fn() unexpected error: %v", err)
	}
	if got != "ado-oidc-token" {
		t.Errorf("fn() = %q, want %q", got, "ado-oidc-token")
	}
}

func TestTokenFunc_NoSource(t *testing.T) {
	// All env vars absent → must return an actionable error, no panic.
	t.Setenv(auth.EnvServiceAccountFederatedToken, "")
	t.Setenv(auth.EnvFederatedTokenFile, "")
	t.Setenv(auth.EnvGitHubRequestURL, "")
	t.Setenv(auth.EnvGitHubRequestToken, "")
	t.Setenv(auth.EnvAzureOIDCRequestURI, "")
	t.Setenv(auth.EnvAzureAccessToken, "")

	_, err := auth.OIDCTokenFunc()
	if err == nil {
		t.Fatal("OIDCTokenFunc() expected error when no OIDC source is available, got nil")
	}
}

func TestTokenFunc_GitHubURL_NoToken(t *testing.T) {
	// URL present but token absent → should fall through to Azure / error.
	t.Setenv(auth.EnvServiceAccountFederatedToken, "")
	t.Setenv(auth.EnvFederatedTokenFile, "")
	t.Setenv(auth.EnvGitHubRequestURL, "https://example.com")
	t.Setenv(auth.EnvGitHubRequestToken, "")
	t.Setenv(auth.EnvAzureOIDCRequestURI, "")
	t.Setenv(auth.EnvAzureAccessToken, "")

	_, err := auth.OIDCTokenFunc()
	if err == nil {
		t.Fatal("OIDCTokenFunc() expected error when GitHub token is missing, got nil")
	}
}
