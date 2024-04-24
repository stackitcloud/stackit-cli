package auth

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
)

type userTokenFlow struct {
	printer                *print.Printer
	reauthorizeUserRoutine func() error // Called if the user needs to login again
	client                 *http.Client
	authFlow               AuthFlow
	accessToken            string
	refreshToken           string
}

// Ensure the implementation satisfies the expected interface
var _ http.RoundTripper = &userTokenFlow{}

// Returns a round tripper that adds authentication according to the user token flow
func UserTokenFlow(p *print.Printer) *userTokenFlow {
	return &userTokenFlow{
		printer:                p,
		reauthorizeUserRoutine: AuthorizeUser,
		client:                 &http.Client{},
	}
}

func (utf *userTokenFlow) RoundTrip(req *http.Request) (*http.Response, error) {
	err := loadVarsFromStorage(utf)
	if err != nil {
		return nil, err
	}
	if utf.authFlow != AUTH_FLOW_USER_TOKEN {
		return nil, fmt.Errorf("auth flow is not user token")
	}

	accessTokenValid := false
	if accessTokenExpired, err := tokenExpired(utf.accessToken); err != nil {
		return nil, fmt.Errorf("check if access token has expired: %w", err)
	} else if !accessTokenExpired {
		accessTokenValid = true
	} else if refreshTokenExpired, err := tokenExpired(utf.refreshToken); err != nil {
		return nil, fmt.Errorf("check if refresh token has expired: %w", err)
	} else if !refreshTokenExpired {
		err = refreshTokens(utf)
		if err == nil {
			accessTokenValid = true
		} else {
			utf.printer.Debug(print.ErrorLevel, "refresh access token: %v", err)
		}
	}

	if !accessTokenValid {
		utf.printer.Warn("Session expired, logging in again...")
		err = reauthenticateUser(utf)
		if err != nil {
			return nil, fmt.Errorf("reauthenticate user: %w", err)
		}
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", utf.accessToken))
	return utf.client.Do(req)
}

func loadVarsFromStorage(utf *userTokenFlow) error {
	authFlow, err := GetAuthFlow()
	if err != nil {
		return fmt.Errorf("get auth flow type: %w", err)
	}
	authFields := map[authFieldKey]string{
		ACCESS_TOKEN:  "",
		REFRESH_TOKEN: "",
	}
	err = GetAuthFieldMap(authFields)
	if err != nil {
		return fmt.Errorf("get tokens from auth storage: %w", err)
	}

	utf.authFlow = authFlow
	utf.accessToken = authFields[ACCESS_TOKEN]
	utf.refreshToken = authFields[REFRESH_TOKEN]
	return nil
}

func reauthenticateUser(utf *userTokenFlow) error {
	err := utf.reauthorizeUserRoutine()
	if err != nil {
		return fmt.Errorf("authenticate user: %w", err)
	}
	err = loadVarsFromStorage(utf)
	if err != nil {
		return fmt.Errorf("load auth vars after user authentication: %w", err)
	}
	if utf.authFlow != AUTH_FLOW_USER_TOKEN {
		return fmt.Errorf("auth flow is not user token")
	}
	return nil
}

func tokenExpired(token string) (bool, error) {
	// We can safely use ParseUnverified because we are not authenticating the user at this point.
	// We're just checking the expiration time
	tokenParsed, _, err := jwt.NewParser().ParseUnverified(token, &jwt.RegisteredClaims{})
	if err != nil {
		return false, fmt.Errorf("parse access token: %w", err)
	}
	expirationTimestampNumeric, err := tokenParsed.Claims.GetExpirationTime()
	if err != nil {
		return false, fmt.Errorf("get expiration timestamp from access token: %w", err)
	}
	expirationTimestamp := expirationTimestampNumeric.Time
	now := time.Now()
	return now.After(expirationTimestamp), nil
}

// Refresh access and refresh tokens using a valid refresh token
func refreshTokens(utf *userTokenFlow) (err error) {
	req, err := buildRequestToRefreshTokens(utf)
	if err != nil {
		return fmt.Errorf("build request: %w", err)
	}

	resp, err := utf.client.Do(req)
	if err != nil {
		return fmt.Errorf("call API: %w", err)
	}
	defer func() {
		tempErr := resp.Body.Close()
		if tempErr != nil {
			err = fmt.Errorf("close response body: %w", tempErr)
		}
	}()

	accessToken, refreshToken, err := parseRefreshTokensResponse(resp)
	if err != nil {
		return fmt.Errorf("parse API response: %w", err)
	}
	err = SetAuthFieldMap(map[authFieldKey]string{
		ACCESS_TOKEN:  accessToken,
		REFRESH_TOKEN: refreshToken,
	})
	if err != nil {
		return fmt.Errorf("set refreshed tokens in auth storage: %w", err)
	}
	utf.accessToken = accessToken
	utf.refreshToken = refreshToken
	return nil
}

func buildRequestToRefreshTokens(utf *userTokenFlow) (*http.Request, error) {
	req, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("https://%s/token", authDomain),
		http.NoBody,
	)
	if err != nil {
		return nil, err
	}
	reqQuery := url.Values{}
	reqQuery.Set("grant_type", "refresh_token")
	reqQuery.Set("client_id", clientId)
	reqQuery.Set("refresh_token", utf.refreshToken)
	reqQuery.Set("token_format", "jwt")
	req.URL.RawQuery = reqQuery.Encode()

	// without this header, the API returns error "An Authentication object was not found in the SecurityContext"
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	return req, nil
}

func parseRefreshTokensResponse(resp *http.Response) (accessToken, refreshToken string, err error) {
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", fmt.Errorf("read body: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf("non-OK %d status: %s", resp.StatusCode, string(respBody))
	}

	respContent := struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}{}
	err = json.Unmarshal(respBody, &respContent)
	if err != nil {
		return "", "", fmt.Errorf("unmarshal body: %w", err)
	}
	if respContent.AccessToken == "" {
		return "", "", fmt.Errorf("no access token found")
	}
	if respContent.RefreshToken == "" {
		return "", "", fmt.Errorf("refresh token found")
	}
	return respContent.AccessToken, respContent.RefreshToken, nil
}
