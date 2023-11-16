package auth

import (
	"fmt"
	"strconv"
	"time"

	"stackit/internal/pkg/config"

	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/cobra"
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
func AuthenticationConfig(cmd *cobra.Command, reauthorizeUserRoutine func() error) (authCfgOption sdkConfig.ConfigurationOption, err error) {
	flow, err := GetAuthFlow()
	if err != nil {
		return nil, fmt.Errorf("get authentication flow: %w", err)
	}
	if flow == "" {
		return nil, fmt.Errorf("authentication flow not set")
	}

	userSessionExpired, err := userSessionExpired()
	if err != nil {
		return nil, fmt.Errorf("check if user session expired: %w", err)
	}

	switch flow {
	case AUTH_FLOW_SERVICE_ACCOUNT_TOKEN:
		if userSessionExpired {
			return nil, fmt.Errorf("session expired")
		}
		accessToken, err := getAccessToken()
		if err != nil {
			return nil, fmt.Errorf("get service account access token: %w", err)
		}
		authCfgOption = sdkConfig.WithToken(accessToken)
	case AUTH_FLOW_SERVICE_ACCOUNT_KEY:
		if userSessionExpired {
			return nil, fmt.Errorf("session expired")
		}
		keyFlow, err := initKeyFlowWithStorage()
		if err != nil {
			return nil, fmt.Errorf("initialize service account key flow: %w", err)
		}
		authCfgOption = sdkConfig.WithCustomAuth(keyFlow)
	case AUTH_FLOW_USER_TOKEN:
		if userSessionExpired {
			cmd.Println("Session expired, logging in again...")
			err = reauthorizeUserRoutine()
			if err != nil {
				return nil, fmt.Errorf("user login: %w", err)
			}
		}
		userTokenFlow := UserTokenFlow(cmd)
		authCfgOption = sdkConfig.WithCustomAuth(userTokenFlow)
	default:
		return nil, fmt.Errorf("the provided authentication flow (%s) is not supported", flow)
	}
	return authCfgOption, nil
}

func userSessionExpired() (bool, error) {
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

func getAccessToken() (string, error) {
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
