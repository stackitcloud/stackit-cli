package auth

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"strconv"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-sdk-go/core/clients"
	sdkConfig "github.com/stackitcloud/stackit-sdk-go/core/config"
	"github.com/zalando/go-keyring"
)

const saKeyStrPattern = `{
	"active": true,
	"createdAt": "2023-03-23T18:26:20.335Z",
	"credentials": {
	  "aud": "https://stackit-service-account-prod.apps.01.cf.eu01.stackit.cloud",
	  "iss": "stackit@sa.stackit.cloud",
	  "kid": "%s",
	  "sub": "%s"
	},
	"id": "%s",
	"keyAlgorithm": "RSA_2048",
	"keyOrigin": "USER_PROVIDED",
	"keyType": "USER_MANAGED",
	"publicKey": "...",
	"validUntil": "2024-03-22T18:05:41Z"
}`

var (
	testSigningKey        = []byte("Test")
	testServiceAccountKey = fmt.Sprintf(saKeyStrPattern, uuid.New().String(), uuid.New().String(), uuid.New().String())
)

func generatePrivateKey() ([]byte, error) {
	// Generate a new RSA key pair with a size of 2048 bits
	privKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}

	// Encode the private key in PEM format
	privKeyPEM := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privKey),
	}

	// Print the private and public keys
	return pem.EncodeToMemory(privKeyPEM), nil
}

func TestAuthenticationConfig(t *testing.T) {
	tests := []struct {
		description                   string
		flow                          AuthFlow
		sessionExpiresAt              time.Time
		accessTokenSet                bool
		refreshToken                  string
		saKey                         string
		privateKeySet                 bool
		tokenEndpoint                 string
		isValid                       bool
		expectedCustomAuthSet         bool
		expectedTokenSet              bool
		expectedReauthorizeUserCalled bool
	}{
		{
			description:      "base_service_account_token",
			flow:             AUTH_FLOW_SERVICE_ACCOUNT_TOKEN,
			sessionExpiresAt: time.Now().Add(time.Hour),
			accessTokenSet:   true,
			refreshToken:     "refresh_token",
			isValid:          true,
			expectedTokenSet: true,
		},
		{
			description:      "service_account_token_session_expired",
			flow:             AUTH_FLOW_SERVICE_ACCOUNT_TOKEN,
			sessionExpiresAt: time.Now().Add(-time.Hour),
			accessTokenSet:   true,
			refreshToken:     "refresh_token",
			isValid:          false,
		},
		{
			description:           "base_service_account_key",
			flow:                  AUTH_FLOW_SERVICE_ACCOUNT_KEY,
			sessionExpiresAt:      time.Now().Add(time.Hour),
			accessTokenSet:        true,
			refreshToken:          "refresh_token",
			saKey:                 testServiceAccountKey,
			privateKeySet:         true,
			tokenEndpoint:         "token_url",
			isValid:               true,
			expectedCustomAuthSet: true,
		},
		{
			description:      "service_account_key_session_expired",
			flow:             AUTH_FLOW_SERVICE_ACCOUNT_KEY,
			sessionExpiresAt: time.Now().Add(-time.Hour),
			accessTokenSet:   true,
			refreshToken:     "refresh_token",
			saKey:            testServiceAccountKey,
			privateKeySet:    true,
			tokenEndpoint:    "token_url",
			isValid:          false,
		},
		{
			description:      "base_user_token",
			flow:             AUTH_FLOW_USER_TOKEN,
			sessionExpiresAt: time.Now().Add(time.Hour),
			accessTokenSet:   true,
			refreshToken:     "refresh_token",
			isValid:          true,
		},
		{
			description:                   "user_token_session_expired",
			flow:                          AUTH_FLOW_USER_TOKEN,
			sessionExpiresAt:              time.Now().Add(-time.Hour),
			accessTokenSet:                true,
			refreshToken:                  "refresh_token",
			isValid:                       true,
			expectedReauthorizeUserCalled: true,
		},
		{
			description: "unsupported_flow",
			flow:        "test_flow",
			isValid:     false,
		},
		{
			description:    "unset_access_token",
			accessTokenSet: false,
			isValid:        false,
		},
		{
			description: "unset_flow",
			flow:        "",
			isValid:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			keyring.MockInit()
			timestamp := time.Now().Add(24 * time.Hour)
			authFields := make(map[authFieldKey]string)
			var accessToken string
			var err error
			if tt.accessTokenSet {
				accessTokenJWT := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
					ExpiresAt: jwt.NewNumericDate(timestamp)})
				accessToken, err = accessTokenJWT.SignedString(testSigningKey)
				if err != nil {
					t.Fatalf("Get test access token as string: %s", err)
				}
			}

			if tt.privateKeySet {
				privateKey, err := generatePrivateKey()
				if err != nil {
					t.Fatalf("Generate private key: %s", err)
				}
				authFields[PRIVATE_KEY] = string(privateKey)
			}
			authFields[SESSION_EXPIRES_AT_UNIX] = strconv.FormatInt(tt.sessionExpiresAt.Unix(), 10)
			authFields[ACCESS_TOKEN] = accessToken
			authFields[REFRESH_TOKEN] = tt.refreshToken
			authFields[SERVICE_ACCOUNT_KEY] = tt.saKey
			authFields[TOKEN_CUSTOM_ENDPOINT] = tt.tokenEndpoint

			err = SetAuthFlow(tt.flow)
			if err != nil {
				t.Fatalf("Failed to set auth flow: %s", err)
			}
			err = SetAuthFieldMap(authFields)
			if err != nil {
				t.Fatalf("Failed to set in auth storage: %v", err)
			}

			reauthorizeUserCalled := false
			reauthenticateUser := func(_ *print.Printer, _ StorageContext, _ bool) error {
				if reauthorizeUserCalled {
					t.Errorf("user reauthorized more than once")
				}
				reauthorizeUserCalled = true
				return nil
			}

			cmd := &cobra.Command{}
			cmd.SetOut(io.Discard) // Suppresses console prints
			p := &print.Printer{Cmd: cmd}

			authCfgOption, err := AuthenticationConfig(p, reauthenticateUser)

			if !tt.isValid {
				if err == nil {
					t.Fatalf("Expected error but no error was returned")
				}
			} else {
				if err != nil {
					t.Fatalf("Expected no error but error was returned: %v", err)
				}

				if reauthorizeUserCalled && !tt.expectedReauthorizeUserCalled {
					t.Errorf("Unexpected user reauthentication")
				} else if !reauthorizeUserCalled && tt.expectedReauthorizeUserCalled {
					t.Errorf("User wasn't reauthenticated when it should've been")
				}

				baseCfg := &sdkConfig.Configuration{}
				err := authCfgOption(baseCfg)
				if err != nil {
					t.Fatalf("Applying returned auth config option: %v", err)
				}
				if tt.expectedCustomAuthSet && baseCfg.CustomAuth == nil {
					t.Fatalf("The returned auth configuration option should set the CustomAuth field but it is nil")
				}
				if tt.expectedTokenSet && baseCfg.Token == "" {
					t.Fatalf("The returned auth configuration option should set the Token field but it is empty")
				}
			}
		})
	}
}

func TestInitKeyFlow(t *testing.T) {
	tests := []struct {
		description    string
		accessTokenSet bool
		refreshToken   string
		saKey          string
		privateKeySet  bool
		tokenEndpoint  string
		isValid        bool
	}{
		{
			description:    "base",
			accessTokenSet: true,
			refreshToken:   "refresh_token",
			saKey:          testServiceAccountKey,
			privateKeySet:  true,
			tokenEndpoint:  "token_url",
			isValid:        true,
		},
		{
			description:    "invalid_service_account_key",
			accessTokenSet: true,
			refreshToken:   "refresh_token",
			saKey:          "",
			privateKeySet:  true,
			tokenEndpoint:  "token_url",
			isValid:        false,
		},
		{
			description:    "invalid_private_key",
			accessTokenSet: true,
			refreshToken:   "refresh_token",
			saKey:          testServiceAccountKey,
			privateKeySet:  false,
			tokenEndpoint:  "token_url",
			isValid:        false,
		},
		{
			description:    "invalid_access_token",
			accessTokenSet: false,
			refreshToken:   "refresh_token",
			saKey:          testServiceAccountKey,
			privateKeySet:  true,
			tokenEndpoint:  "token_url",
			isValid:        false,
		},
		{
			description:    "empty_refresh_token",
			accessTokenSet: false,
			refreshToken:   "",
			saKey:          testServiceAccountKey,
			privateKeySet:  true,
			tokenEndpoint:  "token_url",
			isValid:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			keyring.MockInit()
			timestamp := time.Now().Add(24 * time.Hour)
			authFields := make(map[authFieldKey]string)
			var accessToken string
			var err error
			if tt.accessTokenSet {
				accessTokenJWT := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
					ExpiresAt: jwt.NewNumericDate(timestamp)})
				accessToken, err = accessTokenJWT.SignedString(testSigningKey)
				if err != nil {
					t.Fatalf("Get test access token as string: %s", err)
				}
			}
			if tt.privateKeySet {
				privateKey, err := generatePrivateKey()
				if err != nil {
					t.Fatalf("Generate private key: %s", err)
				}
				authFields[PRIVATE_KEY] = string(privateKey)
			}
			authFields[ACCESS_TOKEN] = accessToken
			authFields[REFRESH_TOKEN] = tt.refreshToken
			authFields[SERVICE_ACCOUNT_KEY] = tt.saKey
			authFields[TOKEN_CUSTOM_ENDPOINT] = tt.tokenEndpoint
			err = SetAuthFieldMap(authFields)
			if err != nil {
				t.Fatalf("Failed to set in auth storage: %v", err)
			}

			keyFlowWithStorage, err := initKeyFlowWithStorage()

			if !tt.isValid {
				if err == nil {
					t.Fatalf("Expected error but no error was returned")
				}
			} else {
				if err != nil {
					t.Fatalf("Expected no error but error was returned: %v", err)
				}
				expectedToken := &clients.TokenResponseBody{
					AccessToken:  accessToken,
					ExpiresIn:    int(timestamp.Unix()),
					RefreshToken: tt.refreshToken,
					Scope:        "",
					TokenType:    "Bearer",
				}
				if !cmp.Equal(*expectedToken, keyFlowWithStorage.keyFlow.GetToken()) {
					t.Errorf("The returned result is wrong. Expected %+v, got %+v", expectedToken, keyFlowWithStorage.keyFlow.GetToken())
				}
			}
		})
	}
}
