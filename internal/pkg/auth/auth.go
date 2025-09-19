package auth

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/stackitcloud/stackit-cli/internal/pkg/config"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"

	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
	sdkConfig "github.com/stackitcloud/stackit-sdk-go/core/config"
)

type tokenClaims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

// AuthenticationConfig reads the credentials from the storage and initializes the authentication flow.
// It returns the configuration option that can be used to create an authenticated SDK client.
//
// If the user was logged in and the user session expired, reauthorizeUserRoutine is called to reauthenticate the user again.
// If the environment variable STACKIT_ACCESS_TOKEN is set this token is used instead.
func AuthenticationConfig(p *print.Printer, reauthorizeUserRoutine func(p *print.Printer, _ bool) error) (authCfgOption sdkConfig.ConfigurationOption, err error) {
	// Get access token from env and use this if present
	accessToken := os.Getenv(envAccessTokenName)
	if accessToken != "" {
		authCfgOption = sdkConfig.WithToken(accessToken)
		return authCfgOption, nil
	}

	flow, err := GetAuthFlow()
	if err != nil {
		return nil, fmt.Errorf("get authentication flow: %w", err)
	}
	if flow == "" {
		return nil, fmt.Errorf("authentication flow not set")
	}

	userSessionExpired, err := UserSessionExpired()
	if err != nil {
		return nil, fmt.Errorf("check if user session expired: %w", err)
	}

	switch flow {
	case AUTH_FLOW_SERVICE_ACCOUNT_TOKEN:
		p.Debug(print.DebugLevel, "authenticating using service account token")
		if userSessionExpired {
			return nil, fmt.Errorf("session expired")
		}
		accessToken, err := GetAccessToken()
		if err != nil {
			return nil, fmt.Errorf("get service account access token: %w", err)
		}
		authCfgOption = sdkConfig.WithToken(accessToken)
	case AUTH_FLOW_SERVICE_ACCOUNT_KEY:
		p.Debug(print.DebugLevel, "authenticating using service account key")
		if userSessionExpired {
			return nil, fmt.Errorf("session expired")
		}
		keyFlow, err := initKeyFlowWithStorage()
		if err != nil {
			return nil, fmt.Errorf("initialize service account key flow: %w", err)
		}
		authCfgOption = sdkConfig.WithCustomAuth(keyFlow)
	case AUTH_FLOW_USER_TOKEN:
		p.Debug(print.DebugLevel, "authenticating using user token")
		if userSessionExpired {
			err = reauthorizeUserRoutine(p, true)
			if err != nil {
				return nil, fmt.Errorf("user login: %w", err)
			}
		}
		userTokenFlow := UserTokenFlow(p)
		authCfgOption = sdkConfig.WithCustomAuth(userTokenFlow)
	default:
		return nil, fmt.Errorf("the provided authentication flow (%s) is not supported", flow)
	}
	return authCfgOption, nil
}

func UserSessionExpired() (bool, error) {
	sessionExpiresAtString, err := GetAuthField(SESSION_EXPIRES_AT_UNIX)
	if err != nil {
		return false, fmt.Errorf("get %s: %w", SESSION_EXPIRES_AT_UNIX, err)
	}
	sessionExpiresAtInt, err := strconv.Atoi(sessionExpiresAtString)
	if err != nil {
		return false, fmt.Errorf("parse session expiration value \"%s\": %w", sessionExpiresAtString, err)
	}
	sessionExpiresAt := time.Unix(int64(sessionExpiresAtInt), 0)
	now := time.Now()
	return now.After(sessionExpiresAt), nil
}

func GetAccessToken() (string, error) {
	accessToken, err := GetAuthField(ACCESS_TOKEN)
	if err != nil {
		return "", fmt.Errorf("get %s: %w", ACCESS_TOKEN, err)
	}
	if accessToken == "" {
		return "", fmt.Errorf("%s not set", ACCESS_TOKEN)
	}
	return accessToken, nil
}

func getStartingSessionExpiresAtUnix() (string, error) {
	sessionStart := time.Now()
	sessionTimeLimitString := viper.GetString(config.SessionTimeLimitKey)
	sessionTimeLimit, err := time.ParseDuration(sessionTimeLimitString)
	if err != nil {
		return "", fmt.Errorf("parse session time limit \"%s\": %w", sessionTimeLimitString, err)
	}
	sessionExpiresAt := sessionStart.Add(sessionTimeLimit)
	return strconv.FormatInt(sessionExpiresAt.Unix(), 10), nil
}

func getEmailFromToken(token string) (string, error) {
	// We can safely use ParseUnverified because we are not authenticating the user at this point,
	// We are parsing the token just to get the service account e-mail
	parsedAccessToken, _, err := jwt.NewParser().ParseUnverified(token, &tokenClaims{})
	if err != nil {
		return "", fmt.Errorf("parse token: %w", err)
	}
	claims, ok := parsedAccessToken.Claims.(*tokenClaims)
	if !ok {
		return "", fmt.Errorf("get claims from parsed token: unknown claims type, please report this issue")
	}

	return claims.Email, nil
}

// GetValidAccessToken returns a valid access token for the current authentication flow.
// For user token flows, it refreshes the token if necessary.
// For service account flows, it returns the current access token.
func GetValidAccessToken(p *print.Printer) (string, error) {
	flow, err := GetAuthFlow()
	if err != nil {
		return "", fmt.Errorf("get authentication flow: %w", err)
	}

	// For service account flows, just return the current token
	if flow == AUTH_FLOW_SERVICE_ACCOUNT_TOKEN || flow == AUTH_FLOW_SERVICE_ACCOUNT_KEY {
		return GetAccessToken()
	}

	if flow != AUTH_FLOW_USER_TOKEN {
		return "", fmt.Errorf("unsupported authentication flow: %s", flow)
	}

	// Load tokens from storage
	authFields := map[authFieldKey]string{
		ACCESS_TOKEN:       "",
		REFRESH_TOKEN:      "",
		IDP_TOKEN_ENDPOINT: "",
	}
	err = GetAuthFieldMap(authFields)
	if err != nil {
		return "", fmt.Errorf("get tokens from auth storage: %w", err)
	}

	accessToken := authFields[ACCESS_TOKEN]
	refreshToken := authFields[REFRESH_TOKEN]
	tokenEndpoint := authFields[IDP_TOKEN_ENDPOINT]

	if accessToken == "" {
		return "", fmt.Errorf("access token not set")
	}
	if refreshToken == "" {
		return "", fmt.Errorf("refresh token not set")
	}
	if tokenEndpoint == "" {
		return "", fmt.Errorf("token endpoint not set")
	}

	// Check if access token is expired
	accessTokenExpired, err := TokenExpired(accessToken)
	if err != nil {
		return "", fmt.Errorf("check if access token has expired: %w", err)
	}
	if !accessTokenExpired {
		// Token is still valid, return it
		return accessToken, nil
	}

	p.Debug(print.DebugLevel, "access token expired, refreshing...")

	// Create a temporary userTokenFlow to reuse the refresh logic
	utf := &userTokenFlow{
		printer:       p,
		client:        &http.Client{},
		authFlow:      flow,
		accessToken:   accessToken,
		refreshToken:  refreshToken,
		tokenEndpoint: tokenEndpoint,
	}

	// Refresh the tokens
	err = refreshTokens(utf)
	if err != nil {
		return "", fmt.Errorf("access token and refresh token expired: %w", err)
	}

	// Return the new access token
	return utf.accessToken, nil
}
