package auth

import (
	"context"
	"fmt"
	"os"

	"github.com/stackitcloud/stackit-sdk-go/core/oidcadapters"
)

const (
	EnvUseOIDC                      = "STACKIT_USE_OIDC"
	EnvServiceAccountEmail          = "STACKIT_SERVICE_ACCOUNT_EMAIL"
	EnvServiceAccountFederatedToken = "STACKIT_SERVICE_ACCOUNT_FEDERATED_TOKEN" //nolint:gosec // linter false positive
	EnvFederatedTokenFile           = "STACKIT_FEDERATED_TOKEN_FILE"            //nolint:gosec // linter false positive
	EnvGitHubRequestURL             = "ACTIONS_ID_TOKEN_REQUEST_URL"
	EnvGitHubRequestToken           = "ACTIONS_ID_TOKEN_REQUEST_TOKEN" //nolint:gosec // linter false positive
	EnvAzureOIDCRequestURI          = "SYSTEM_OIDCREQUESTURI"
	EnvAzureAccessToken             = "SYSTEM_ACCESSTOKEN" //nolint:gosec // linter false positive
)

func IsOIDCEnabled() bool {
	return os.Getenv(EnvUseOIDC) == "1"
}

// IsOIDCEnabledWithOverride resolves OIDC mode using explicit input first and env fallback.
// If useOIDC is not nil, its value is used directly; otherwise STACKIT_USE_OIDC is evaluated.
func IsOIDCEnabledWithOverride(useOIDC *bool) bool {
	if useOIDC != nil {
		return *useOIDC
	}

	return IsOIDCEnabled()
}

func OIDCServiceAccountEmail() string {
	return os.Getenv(EnvServiceAccountEmail)
}

// TokenFunc returns the OIDCTokenFunc to use for Workload Identity Federation.
// It checks the following token sources in order: STACKIT_SERVICE_ACCOUNT_FEDERATED_TOKEN,
// STACKIT_FEDERATED_TOKEN_FILE, GitHub Actions (ACTIONS_ID_TOKEN_REQUEST_URL +
// ACTIONS_ID_TOKEN_REQUEST_TOKEN), and Azure DevOps (SYSTEM_OIDCREQUESTURI + SYSTEM_ACCESSTOKEN).
// Returns an error if no source is detected.
func OIDCTokenFunc() (oidcadapters.OIDCTokenFunc, error) {
	// static token provided directly via env var
	if token := os.Getenv(EnvServiceAccountFederatedToken); token != "" {
		return func(_ context.Context) (string, error) {
			return token, nil
		}, nil
	}

	// token read from filesystem path via env var
	if tokenFilePath := os.Getenv(EnvFederatedTokenFile); tokenFilePath != "" {
		return oidcadapters.ReadJWTFromFileSystem(tokenFilePath), nil
	}

	// GitHub Actions
	if ghURL := os.Getenv(EnvGitHubRequestURL); ghURL != "" {
		if ghToken := os.Getenv(EnvGitHubRequestToken); ghToken != "" {
			return oidcadapters.RequestGHOIDCToken(ghURL, ghToken), nil
		}
	}

	// Azure DevOps
	if adoURL := os.Getenv(EnvAzureOIDCRequestURI); adoURL != "" {
		if adoToken := os.Getenv(EnvAzureAccessToken); adoToken != "" {
			return oidcadapters.RequestAzureDevOpsOIDCToken(adoURL, adoToken, ""), nil
		}
	}

	return nil, fmt.Errorf(
		"%s is enabled but no OIDC token source was detected\n"+
			"Provide the token via %s or %s, or run in a supported CI environment:\n"+
			"  - GitHub Actions: grant 'id-token: write' permission; %s and %s are set automatically by the runner\n"+
			"  - Azure DevOps:   pass 'SYSTEM_ACCESSTOKEN: $(System.AccessToken)' in your pipeline step",
		EnvUseOIDC, EnvServiceAccountFederatedToken, EnvFederatedTokenFile,
		EnvGitHubRequestURL, EnvGitHubRequestToken,
	)
}
