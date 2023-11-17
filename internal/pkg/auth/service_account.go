package auth

import (
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stackitcloud/stackit-sdk-go/core/clients"
)

type tokenClaims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

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
func AuthenticateServiceAccount(rt http.RoundTripper) (email string, err error) {
	authFields := make(map[authFieldKey]string)
	var authFlowType authFlow
	switch flow := rt.(type) {
	case keyFlowInterface:
		authFlowType = AUTH_FLOW_SERVICE_ACCOUNT_KEY

		accessToken, err := flow.GetAccessToken()
		if err != nil {
			return "", fmt.Errorf("get access token: %w", err)
		}
		authFields[ACCESS_TOKEN] = accessToken
		authFields[REFRESH_TOKEN] = flow.GetToken().RefreshToken
		authFields[SERVICE_ACCOUNT_KEY] = flow.GetConfig().ServiceAccountKey
		authFields[PRIVATE_KEY] = flow.GetConfig().PrivateKey
	case tokenFlowInterface:
		authFlowType = AUTH_FLOW_SERVICE_ACCOUNT_TOKEN

		authFields[ACCESS_TOKEN] = flow.GetConfig().ServiceAccountToken
	default:
		return "", fmt.Errorf("could not authenticate using any of the supported authentication flows (key and token): please report this issue")
	}

	// We can safely use ParseUnverified because we are not authenticating the user at this point,
	// We are parsing the token just to get the service account e-mail
	parsedAccessToken, _, err := jwt.NewParser().ParseUnverified(authFields[ACCESS_TOKEN], &tokenClaims{})
	if err != nil {
		return "", fmt.Errorf("parse access token to read service account email: %w", err)
	}
	claims, ok := parsedAccessToken.Claims.(*tokenClaims)
	if !ok {
		return "", fmt.Errorf("get claims from parsed access token: unknown claims type, please report this issue")
	}
	authFields[SERVICE_ACCOUNT_EMAIL] = claims.Email

	sessionExpiresAtUnix, err := getStartingSessionExpiresAtUnix()
	if err != nil {
		return "", fmt.Errorf("compute session expiration timestamp: %w", err)
	}
	authFields[SESSION_EXPIRES_AT_UNIX] = sessionExpiresAtUnix

	err = SetAuthFlow(authFlowType)
	if err != nil {
		return "", fmt.Errorf("set auth flow type: %w", err)
	}
	err = SetAuthFieldMap(authFields)
	if err != nil {
		return "", fmt.Errorf("set in auth storage: %w", err)
	}

	return authFields[SERVICE_ACCOUNT_EMAIL], nil
}

// initKeyFlowWithStorage initializes the keyFlow from the SDK and creates a keyFlowWithStorage struct that uses that keyFlow
func initKeyFlowWithStorage() (*keyFlowWithStorage, error) {
	authFields := map[authFieldKey]string{
		ACCESS_TOKEN:          "",
		REFRESH_TOKEN:         "",
		SERVICE_ACCOUNT_KEY:   "",
		PRIVATE_KEY:           "",
		TOKEN_CUSTOM_ENDPOINT: "",
		JWKS_CUSTOM_ENDPOINT:  "",
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

	cfg := &clients.KeyFlowConfig{
		ServiceAccountKey: authFields[SERVICE_ACCOUNT_KEY],
		PrivateKey:        authFields[PRIVATE_KEY],
		ClientRetry:       clients.NewRetryConfig(),
		TokenUrl:          authFields[TOKEN_CUSTOM_ENDPOINT],
		JWKSUrl:           authFields[JWKS_CUSTOM_ENDPOINT],
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
