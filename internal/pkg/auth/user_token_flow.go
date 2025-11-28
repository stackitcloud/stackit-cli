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
	reauthorizeUserRoutine func(p *print.Printer, context StorageContext, isReauthentication bool) error // Called if the user needs to login again
	client                 *http.Client
	context                StorageContext
	authFlow               AuthFlow
	accessToken            string
	refreshToken           string
	tokenEndpoint          string
}

// Ensure the implementation satisfies the expected interface
var _ http.RoundTripper = &userTokenFlow{}

// Returns a round tripper that adds authentication according to the user token flow
// Uses the CLI storage context by default
func UserTokenFlow(p *print.Printer) *userTokenFlow {
	return UserTokenFlowWithContext(p, StorageContextCLI)
}

// Returns a round tripper that adds authentication according to the user token flow
// with the specified storage context
func UserTokenFlowWithContext(p *print.Printer, context StorageContext) *userTokenFlow {
	return &userTokenFlow{
		printer:                p,
		reauthorizeUserRoutine: AuthorizeUser,
		client:                 &http.Client{},
		context:                context,
	}
}

func (utf *userTokenFlow) RoundTrip(req *http.Request) (*http.Response, error) {
	// Set the storage printer so debug messages use the correct verbosity
	SetStoragePrinter(utf.printer)

	err := loadVarsFromStorage(utf)
	if err != nil {
		return nil, err
	}
	if utf.authFlow != AUTH_FLOW_USER_TOKEN {
		return nil, fmt.Errorf("auth flow is not user token")
	}

	accessTokenValid := false
	accessTokenExpired, err := TokenExpired(utf.accessToken)
	if err != nil {
		return nil, fmt.Errorf("check if access token has expired: %w", err)
	} else if !accessTokenExpired {
		accessTokenValid = true
	} else {
		utf.printer.Debug(print.DebugLevel, "access token expired, refreshing...")
		err = refreshTokens(utf)
		if err == nil {
			accessTokenValid = true
		} else {
			utf.printer.Debug(print.ErrorLevel, "refresh access token: %w", err)
		}
	}

	if !accessTokenValid {
		utf.printer.Debug(print.DebugLevel, "user access token is not valid, reauthenticating...")
		err = reauthenticateUser(utf)
		if err != nil {
			return nil, fmt.Errorf("reauthenticate user: %w", err)
		}
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", utf.accessToken))
	return utf.client.Do(req)
}

func loadVarsFromStorage(utf *userTokenFlow) error {
	authFlow, err := GetAuthFlowWithContext(utf.context)
	if err != nil {
		return fmt.Errorf("get auth flow type: %w", err)
	}
	authFields := map[authFieldKey]string{
		ACCESS_TOKEN:       "",
		REFRESH_TOKEN:      "",
		IDP_TOKEN_ENDPOINT: "",
	}
	err = GetAuthFieldMapWithContext(utf.context, authFields)
	if err != nil {
		return fmt.Errorf("get tokens from auth storage: %w", err)
	}

	utf.authFlow = authFlow
	utf.accessToken = authFields[ACCESS_TOKEN]
	utf.refreshToken = authFields[REFRESH_TOKEN]
	utf.tokenEndpoint = authFields[IDP_TOKEN_ENDPOINT]
	return nil
}

func reauthenticateUser(utf *userTokenFlow) error {
	err := utf.reauthorizeUserRoutine(utf.printer, utf.context, true)
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

func TokenExpired(token string) (bool, error) {
	// We can safely use ParseUnverified because we are not authenticating the user at this point.
	// We're just checking the expiration time
	tokenParsed, _, err := jwt.NewParser().ParseUnverified(token, &jwt.RegisteredClaims{})
	if err != nil {
		return false, fmt.Errorf("parse access token: %w", err)
	}
	expirationTimestampNumeric, err := tokenParsed.Claims.GetExpirationTime()
	if err != nil {
		return false, fmt.Errorf("get expiration timestamp from access token: %w", err)
	} else if expirationTimestampNumeric == nil {
		return false, nil
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

	// Debug log the request
	debugHTTPRequest(utf.printer, req)

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

	// Debug log the response
	debugHTTPResponse(utf.printer, resp)

	accessToken, refreshToken, err := parseRefreshTokensResponse(resp)
	if err != nil {
		return fmt.Errorf("parse API response: %w", err)
	}

	// Get the new access token's expiration time
	expiresAtUnix, err := getAccessTokenExpiresAtUnix(accessToken)
	if err != nil {
		return fmt.Errorf("get access token expiration: %w", err)
	}

	err = SetAuthFieldMapWithContext(utf.context, map[authFieldKey]string{
		ACCESS_TOKEN:            accessToken,
		REFRESH_TOKEN:           refreshToken,
		SESSION_EXPIRES_AT_UNIX: expiresAtUnix,
	})
	if err != nil {
		return fmt.Errorf("set refreshed tokens in auth storage: %w", err)
	}
	utf.accessToken = accessToken
	utf.refreshToken = refreshToken
	return nil
}

func buildRequestToRefreshTokens(utf *userTokenFlow) (*http.Request, error) {
	idpClientID, err := getIDPClientID()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(
		http.MethodPost,
		utf.tokenEndpoint,
		http.NoBody,
	)
	if err != nil {
		return nil, err
	}
	reqQuery := url.Values{}
	reqQuery.Set("grant_type", "refresh_token")
	reqQuery.Set("client_id", idpClientID)
	reqQuery.Set("refresh_token", utf.refreshToken)
	reqQuery.Set("token_format", "jwt")
	req.URL.RawQuery = reqQuery.Encode()

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
