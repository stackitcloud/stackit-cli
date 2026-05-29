package oidc_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stackitcloud/stackit-cli/internal/pkg/auth/oidc"
)

func TestIsEnabled(t *testing.T) {
	tests := []struct {
		value    string
		expected bool
	}{
		{"1", true},
		{"true", true},
		{"True", true},
		{"TRUE", true},
		{"yes", true},
		{"YES", true},
		{"Yes", true},
		{"0", false},
		{"false", false},
		{"no", false},
		{"", false},
		{"random", false},
		{" 1 ", true}, // leading/trailing whitespace
		{" true", true},
	}
	for _, tt := range tests {
		t.Run(tt.value, func(t *testing.T) {
			t.Setenv(oidc.EnvUseOIDC, tt.value)
			got := oidc.IsEnabled()
			if got != tt.expected {
				t.Errorf("IsEnabled() = %v, want %v (env=%q)", got, tt.expected, tt.value)
			}
		})
	}
}

func TestIsEnabled_Unset(t *testing.T) {
	// When the env var is not set at all IsEnabled must return false
	t.Setenv(oidc.EnvUseOIDC, "")
	if oidc.IsEnabled() {
		t.Error("IsEnabled() = true, want false when env var is empty")
	}
}

func TestServiceAccountEmail(t *testing.T) {
	const want = "ci@sa.stackit.cloud"
	t.Setenv(oidc.EnvServiceAccountEmail, want)
	if got := oidc.ServiceAccountEmail(); got != want {
		t.Errorf("ServiceAccountEmail() = %q, want %q", got, want)
	}
}

func TestTokenFunc_StaticToken(t *testing.T) {
	const want = "my-static-oidc-token"
	t.Setenv(oidc.EnvServiceAccountFederatedToken, want)
	// ensure GitHub / Azure vars are absent so we hit the static path first
	t.Setenv(oidc.EnvGitHubRequestURL, "")
	t.Setenv(oidc.EnvGitHubRequestToken, "")
	t.Setenv(oidc.EnvAzureOIDCRequestURI, "")
	t.Setenv(oidc.EnvAzureAccessToken, "")

	fn, err := oidc.TokenFunc()
	if err != nil {
		t.Fatalf("TokenFunc() unexpected error: %v", err)
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

	t.Setenv(oidc.EnvServiceAccountFederatedToken, "")
	t.Setenv(oidc.EnvGitHubRequestURL, srv.URL)
	t.Setenv(oidc.EnvGitHubRequestToken, "gh-bearer-token")
	t.Setenv(oidc.EnvAzureOIDCRequestURI, "")
	t.Setenv(oidc.EnvAzureAccessToken, "")

	fn, err := oidc.TokenFunc()
	if err != nil {
		t.Fatalf("TokenFunc() unexpected error: %v", err)
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

	t.Setenv(oidc.EnvServiceAccountFederatedToken, "")
	t.Setenv(oidc.EnvGitHubRequestURL, "")
	t.Setenv(oidc.EnvGitHubRequestToken, "")
	t.Setenv(oidc.EnvAzureOIDCRequestURI, srv.URL)
	t.Setenv(oidc.EnvAzureAccessToken, "ado-access-token")

	fn, err := oidc.TokenFunc()
	if err != nil {
		t.Fatalf("TokenFunc() unexpected error: %v", err)
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
	t.Setenv(oidc.EnvServiceAccountFederatedToken, "")
	t.Setenv(oidc.EnvGitHubRequestURL, "")
	t.Setenv(oidc.EnvGitHubRequestToken, "")
	t.Setenv(oidc.EnvAzureOIDCRequestURI, "")
	t.Setenv(oidc.EnvAzureAccessToken, "")

	_, err := oidc.TokenFunc()
	if err == nil {
		t.Fatal("TokenFunc() expected error when no OIDC source is available, got nil")
	}
}

func TestTokenFunc_GitHubURL_NoToken(t *testing.T) {
	// URL present but token absent → should fall through to Azure / error.
	t.Setenv(oidc.EnvServiceAccountFederatedToken, "")
	t.Setenv(oidc.EnvGitHubRequestURL, "https://example.com")
	t.Setenv(oidc.EnvGitHubRequestToken, "")
	t.Setenv(oidc.EnvAzureOIDCRequestURI, "")
	t.Setenv(oidc.EnvAzureAccessToken, "")

	_, err := oidc.TokenFunc()
	if err == nil {
		t.Fatal("TokenFunc() expected error when GitHub token is missing, got nil")
	}
}
