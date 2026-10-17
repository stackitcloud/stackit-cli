package auth

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stackitcloud/stackit-cli/internal/pkg/config"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testparams"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stackitcloud/stackit-sdk-go/core/clients"
	"github.com/zalando/go-keyring"
)

const (
	tokenFlow = "token"
	keyFlow   = "key"
)

var accessTokenSigningKey = []byte("Test")

type keyFlowMocked struct {
	accessToken        jwt.Token
	config             clients.KeyFlowConfig
	tokenResponse      clients.TokenResponseBody
	getAccessTokenFail bool
	tokenInvalid       bool
}

func (f *keyFlowMocked) GetAccessToken() (string, error) {
	if f.getAccessTokenFail {
		return "", fmt.Errorf("error")
	}
	if f.tokenInvalid {
		return "", nil
	}
	raw, err := f.accessToken.SignedString(accessTokenSigningKey)
	if err != nil {
		return "", fmt.Errorf("sign string from token: %w", err)
	}
	return raw, nil
}

func (f *keyFlowMocked) GetConfig() clients.KeyFlowConfig {
	return f.config
}

func (f *keyFlowMocked) GetToken() clients.TokenResponseBody {
	return f.tokenResponse
}

func (f *keyFlowMocked) RoundTrip(*http.Request) (*http.Response, error) {
	return nil, nil
}

type tokenFlowMocked struct {
	config clients.TokenFlowConfig
}

func (f *tokenFlowMocked) GetConfig() clients.TokenFlowConfig {
	return f.config
}

func (f *tokenFlowMocked) RoundTrip(*http.Request) (*http.Response, error) {
	return nil, nil
}

type wifFlowMocked struct {
	accessToken        string
	getAccessTokenFail bool
}

func (f *wifFlowMocked) GetAccessToken() (string, error) {
	if f.getAccessTokenFail {
		return "", fmt.Errorf("mock WIF error")
	}
	return f.accessToken, nil
}

func (f *wifFlowMocked) RoundTrip(*http.Request) (*http.Response, error) {
	return nil, nil
}

func TestAuthenticateServiceAccount(t *testing.T) {
	tests := []struct {
		description        string
		flowType           string
		getAccessTokenFail bool
		tokenInvalid       bool
		accessToken        jwt.Token
		accessTokenRaw     string
		refreshToken       string
		expectedEmail      string
		isValid            bool
	}{
		{
			description: "base_key_flow",
			flowType:    keyFlow,
			accessToken: *jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
				Email:            "test_email",
				RegisteredClaims: jwt.RegisteredClaims{},
			}),
			refreshToken:  "refresh_token",
			expectedEmail: "test_email",
			isValid:       true,
		},
		{
			description: "base_token_flow",
			flowType:    tokenFlow,
			accessToken: *jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
				Email: "test_email",
			}),
			refreshToken:  "refresh_token",
			expectedEmail: "test_email",
			isValid:       true,
		},
		{
			description: "unsupported_flow",
			flowType:    "unsupported",
			isValid:     false,
		},
		{
			description:        "key_flow_failed_get_access_token",
			flowType:           keyFlow,
			getAccessTokenFail: true,
			isValid:            false,
		},
		{
			description:  "invalid_token",
			flowType:     keyFlow,
			tokenInvalid: true,
			isValid:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			keyring.MockInit()
			config.InitConfig() // AuthenticateServiceAccount accesses the config

			var flow http.RoundTripper
			switch tt.flowType {
			case keyFlow:
				flow = &keyFlowMocked{
					accessToken:        tt.accessToken,
					getAccessTokenFail: tt.getAccessTokenFail,
					tokenInvalid:       tt.tokenInvalid,
					config: clients.KeyFlowConfig{
						ServiceAccountKey: &clients.ServiceAccountKeyResponse{},
						PrivateKey:        "private_key",
					},
					tokenResponse: clients.TokenResponseBody{
						RefreshToken: tt.refreshToken,
					},
				}
			case tokenFlow:
				raw, err := tt.accessToken.SignedString(accessTokenSigningKey)
				if err != nil {
					t.Fatalf("signing string from token: %s", err)
				}
				flow = &tokenFlowMocked{
					config: clients.TokenFlowConfig{
						ServiceAccountToken: raw,
					},
				}
			default:
				flow = &http.Transport{}
			}

			params := testparams.NewTestParams()
			email, _, err := AuthenticateServiceAccount(params.Printer, flow, false)

			if !tt.isValid {
				if err == nil {
					t.Fatalf("Expected error but no error was returned")
				}
			} else {
				if err != nil {
					t.Fatalf("Expected no error but error was returned: %v", err)
				}
				if tt.expectedEmail != email {
					t.Fatalf("The returned email is wrong. Expected %s, got %s", tt.expectedEmail, email)
				}
			}
		})
	}
}

func TestAuthenticateServiceAccount_WIF(t *testing.T) {
	// Build a signed test JWT that getEmailFromToken can parse.
	testEmail := "ci@sa.stackit.cloud"
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		Email:            testEmail,
		RegisteredClaims: jwt.RegisteredClaims{},
	})
	raw, err := tok.SignedString(accessTokenSigningKey)
	if err != nil {
		t.Fatalf("sign test token: %v", err)
	}

	tests := []struct {
		description        string
		accessToken        string
		getAccessTokenFail bool
		disableWriting     bool
		isValid            bool
	}{
		{
			description:    "wif_success_credentials_written",
			accessToken:    raw,
			disableWriting: false,
			isValid:        true,
		},
		{
			description:    "wif_success_writing_disabled",
			accessToken:    raw,
			disableWriting: true,
			isValid:        true,
		},
		{
			description:        "wif_get_access_token_fails",
			getAccessTokenFail: true,
			isValid:            false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			keyring.MockInit()
			config.InitConfig()

			// leftovers from a previous key flow must not survive a WIF activation
			if err := SetAuthFieldMap(map[authFieldKey]string{
				SERVICE_ACCOUNT_KEY: "old_key",
				PRIVATE_KEY:         "old_private_key",
			}); err != nil {
				t.Fatalf("set previous credentials: %v", err)
			}
			if err := SetAuthFlow(AUTH_FLOW_SERVICE_ACCOUNT_KEY); err != nil {
				t.Fatalf("set previous auth flow: %v", err)
			}

			flow := &wifFlowMocked{
				accessToken:        tt.accessToken,
				getAccessTokenFail: tt.getAccessTokenFail,
			}

			params := testparams.NewTestParams()
			email, _, err := AuthenticateServiceAccount(params.Printer, flow, tt.disableWriting)

			if !tt.isValid {
				if err == nil {
					t.Fatal("Expected error but no error was returned")
				}
				return
			}
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			if email != testEmail {
				t.Fatalf("email = %q, want %q", email, testEmail)
			}

			storedFlow, _ := GetAuthFlow()
			storedEmail, _ := GetAuthField(SERVICE_ACCOUNT_EMAIL)
			storedKey, _ := GetAuthField(SERVICE_ACCOUNT_KEY)
			storedPrivKey, _ := GetAuthField(PRIVATE_KEY)

			if tt.disableWriting {
				if storedFlow != AUTH_FLOW_SERVICE_ACCOUNT_KEY {
					t.Errorf("auth flow = %q, want previous flow %q untouched when writing is disabled", storedFlow, AUTH_FLOW_SERVICE_ACCOUNT_KEY)
				}
				return
			}

			// Only the flow and the service account email are needed to rebuild the WIF flow,
			// the federated token is read from the environment on every command.
			if storedFlow != AUTH_FLOW_WORKLOAD_IDENTITY {
				t.Errorf("auth flow = %q, want %q", storedFlow, AUTH_FLOW_WORKLOAD_IDENTITY)
			}
			if storedEmail != testEmail {
				t.Errorf("SERVICE_ACCOUNT_EMAIL = %q, want %q", storedEmail, testEmail)
			}
			if storedKey != "" {
				t.Errorf("SERVICE_ACCOUNT_KEY = %q, want empty in WIF mode", storedKey)
			}
			if storedPrivKey != "" {
				t.Errorf("PRIVATE_KEY = %q, want empty in WIF mode", storedPrivKey)
			}
		})
	}
}
