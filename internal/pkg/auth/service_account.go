package auth

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"

	"github.com/stackitcloud/stackit-sdk-go/core/clients"
)

type keyFlowInterface interface {
	GetAccessToken() (string, error)
	GetConfig() clients.KeyFlowConfig
	GetToken() clients.TokenResponseBody
	RoundTrip(*http.Request) (*http.Response, error)
}

type tokenFlowInterface interface {
	GetConfig() clients.TokenFlowConfig
	RoundTrip(*http.Request) (*http.Response, error)
}

type keyFlowWithStorage struct {
	keyFlow *clients.KeyFlow
}

// Ensure the implementation satisfies the expected interface
var _ http.RoundTripper = &keyFlowWithStorage{}

// AuthenticateServiceAccount checks the type of the provided roundtripper,
// authenticates the CLI accordingly and store the credentials.
// For the key flow, it fetches an access and refresh token from the Service Account API.
// For the token flow, it just stores the provided token and doesn't check if it is valid.
// It returns the email associated with the service account
// If disableWriting is set to true the credentials are not stored on disk (keyring, file).
func AuthenticateServiceAccount(p *print.Printer, rt http.RoundTripper, disableWriting bool) (email, accessToken string, err error) {
	// Set the storage printer so debug messages use the correct verbosity
	SetStoragePrinter(p)

	authFields := make(map[authFieldKey]string)
	var authFlowType AuthFlow
	switch flow := rt.(type) {
	case keyFlowInterface:
		p.Debug(print.DebugLevel, "authenticating using service account key")
		authFlowType = AUTH_FLOW_SERVICE_ACCOUNT_KEY

		accessToken, err := flow.GetAccessToken()
		if err != nil {
			p.Debug(print.ErrorLevel, "get access token: %v", err)
			return "", "", &errors.ActivateServiceAccountError{}
		}
		serviceAccountKey := flow.GetConfig().ServiceAccountKey
		saKeyBytes, err := json.Marshal(serviceAccountKey)
		if err != nil {
			return "", "", fmt.Errorf("marshal service account key: %w", err)
		}

		authFields[ACCESS_TOKEN] = accessToken
		authFields[REFRESH_TOKEN] = flow.GetToken().RefreshToken
		authFields[SERVICE_ACCOUNT_KEY] = string(saKeyBytes)
		authFields[PRIVATE_KEY] = flow.GetConfig().PrivateKey
	case tokenFlowInterface:
		p.Debug(print.DebugLevel, "authenticating using service account token")
		authFlowType = AUTH_FLOW_SERVICE_ACCOUNT_TOKEN

		authFields[ACCESS_TOKEN] = flow.GetConfig().ServiceAccountToken
	default:
		return "", "", fmt.Errorf("could not authenticate using any of the supported authentication flows (key and token): please report this issue")
	}

	email, err = getEmailFromToken(authFields[ACCESS_TOKEN])
	if err != nil {
		return "", "", fmt.Errorf("get email from access token: %w", err)
	}

	p.Debug(print.DebugLevel, "successfully authenticated service account %s", email)

	authFields[SERVICE_ACCOUNT_EMAIL] = email

	sessionExpiresAtUnix, err := getStartingSessionExpiresAtUnix()
	if err != nil {
		return "", "", fmt.Errorf("compute session expiration timestamp: %w", err)
	}
	authFields[SESSION_EXPIRES_AT_UNIX] = sessionExpiresAtUnix

	if !disableWriting {
		err = SetAuthFlow(authFlowType)
		if err != nil {
			return "", "", fmt.Errorf("set auth flow type: %w", err)
		}
		err = SetAuthFieldMap(authFields)
		if err != nil {
			return "", "", fmt.Errorf("set in auth storage: %w", err)
		}
	}

	return authFields[SERVICE_ACCOUNT_EMAIL], authFields[ACCESS_TOKEN], nil
}

// initKeyFlowWithStorage initializes the keyFlow from the SDK and creates a keyFlowWithStorage struct that uses that keyFlow
func initKeyFlowWithStorage() (*keyFlowWithStorage, error) {
	authFields := map[authFieldKey]string{
		ACCESS_TOKEN:          "",
		REFRESH_TOKEN:         "",
		SERVICE_ACCOUNT_KEY:   "",
		PRIVATE_KEY:           "",
		TOKEN_CUSTOM_ENDPOINT: "",
	}
	err := GetAuthFieldMap(authFields)
	if err != nil {
		return nil, fmt.Errorf("get from auth storage: %w", err)
	}
	if authFields[ACCESS_TOKEN] == "" {
		return nil, fmt.Errorf("access token not set")
	}
	if authFields[REFRESH_TOKEN] == "" {
		return nil, fmt.Errorf("refresh token not set")
	}

	var serviceAccountKey = &clients.ServiceAccountKeyResponse{}
	err = json.Unmarshal([]byte(authFields[SERVICE_ACCOUNT_KEY]), serviceAccountKey)
	if err != nil {
		return nil, fmt.Errorf("unmarshalling service account key: %w", err)
	}

	cfg := &clients.KeyFlowConfig{
		ServiceAccountKey: serviceAccountKey,
		PrivateKey:        authFields[PRIVATE_KEY],
		TokenUrl:          authFields[TOKEN_CUSTOM_ENDPOINT],
	}

	keyFlow := &clients.KeyFlow{}
	err = keyFlow.Init(cfg)
	if err != nil {
		return nil, fmt.Errorf("initialize key flow: %w", err)
	}
	err = keyFlow.SetToken(authFields[ACCESS_TOKEN], authFields[REFRESH_TOKEN])
	if err != nil {
		return nil, fmt.Errorf("set access and refresh token: %w", err)
	}

	// create keyFlowWithStorage roundtripper that stores the credentials after executing a request
	keyFlowWithStorage := &keyFlowWithStorage{
		keyFlow: keyFlow,
	}
	return keyFlowWithStorage, nil
}

// The keyFlowWithStorage Roundtrip executes the keyFlow roundtrip and then stores the access and refresh tokens
func (kf *keyFlowWithStorage) RoundTrip(req *http.Request) (*http.Response, error) {
	resp, err := kf.keyFlow.RoundTrip(req)

	token := kf.keyFlow.GetToken()
	accessToken := token.AccessToken
	refreshToken := token.RefreshToken
	tokenValues := map[authFieldKey]string{
		ACCESS_TOKEN:  accessToken,
		REFRESH_TOKEN: refreshToken,
	}

	storageErr := SetAuthFieldMap(tokenValues)
	if storageErr != nil {
		return nil, fmt.Errorf("set access and refresh token in the storage: %w", err)
	}

	return resp, err
}
