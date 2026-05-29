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
			description:    "wif_success_no_credentials_written",
			accessToken:    raw,
			disableWriting: false, // even when false, WIF forces no-write internally
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

			// Verify no credentials were written to the keyring / file.
			// After a WIF authentication, the auth storage should not contain a
			// service account key or private key.
			storedKey, _ := GetAuthField(SERVICE_ACCOUNT_KEY)
			if storedKey != "" {
				t.Errorf("SERVICE_ACCOUNT_KEY was written to storage in WIF mode, want empty")
			}
			storedPrivKey, _ := GetAuthField(PRIVATE_KEY)
			if storedPrivKey != "" {
				t.Errorf("PRIVATE_KEY was written to storage in WIF mode, want empty")
			}
		})
	}
}
