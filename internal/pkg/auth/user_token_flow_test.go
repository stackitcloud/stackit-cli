package auth

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/zalando/go-keyring"
)

type clientTransport struct {
	t                               *testing.T // May write test errors
	requestURL                      string
	refreshTokensFails              bool
	refreshTokensReturnsNonOKStatus bool
	requestSent                     *bool
	tokensRefreshed                 *bool
}

func (rt *clientTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	reqURL := req.Host + req.URL.Path
	if reqURL == rt.requestURL {
		return rt.roundTripRequest()
	}
	if reqURL == fmt.Sprintf("%s/token", authDomain) {
		return rt.roundTripRefreshTokens()
	}
	rt.t.Fatalf("unexpected request to \"%s\"", reqURL)
	return nil, fmt.Errorf("unexpected request to \"%s\"", reqURL)
}

func (rt *clientTransport) roundTripRequest() (*http.Response, error) {
	if *rt.requestSent {
		rt.t.Errorf("request executed multiple times")
	}
	*rt.requestSent = true

	resp := &http.Response{
		Status:     http.StatusText(http.StatusOK),
		StatusCode: http.StatusOK,
	}
	return resp, nil
}

func (rt *clientTransport) roundTripRefreshTokens() (*http.Response, error) {
	if rt.refreshTokensFails {
		return nil, fmt.Errorf("failed")
	}
	if rt.refreshTokensReturnsNonOKStatus {
		resp := &http.Response{
			Status:     http.StatusText(http.StatusInternalServerError),
			StatusCode: http.StatusInternalServerError,
		}
		return resp, nil
	}

	if *rt.tokensRefreshed {
		rt.t.Errorf("tokens refreshed more than once")
	}
	*rt.tokensRefreshed = true
	expirationTimestamp := time.Now().Add(time.Hour)
	accessToken, refreshToken, err := createTokens(expirationTimestamp, expirationTimestamp)
	if err != nil {
		rt.t.Fatalf("refresh token API: failed to create tokens: %v", err)
	}
	respBody := fmt.Sprintf(
		`{
			"access_token": "%s",
			"refresh_token": "%s"
		}`,
		accessToken,
		refreshToken,
	)
	resp := &http.Response{
		Status:     http.StatusText(http.StatusOK),
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewReader([]byte(respBody))),
	}
	return resp, nil
}

type authorizeUserContext struct {
	t                   *testing.T // May write test errors
	authorizeUserFails  bool
	authorizeUserCalled *bool
	tokensRefreshed     *bool
}

func reauthorizeUser(auCtx *authorizeUserContext) error {
	if *auCtx.authorizeUserCalled {
		auCtx.t.Errorf("user authenticated more than once")
	}
	*auCtx.authorizeUserCalled = true

	if auCtx.authorizeUserFails {
		return fmt.Errorf("failed")
	}

	if *auCtx.tokensRefreshed {
		auCtx.t.Errorf("tokens refreshed more than once")
	}
	*auCtx.tokensRefreshed = true
	err := setAuthStorage(
		time.Now().Add(time.Hour),
		time.Now().Add(time.Hour),
		true,
		true,
	)
	if err != nil {
		auCtx.t.Fatalf("failed to set auth vars in authorize user: %v", err)
	}
	return nil
}

func TestRoundTrip(t *testing.T) {
	tests := []struct {
		desc string
		// Test settings
		accessTokenExpiresAt            time.Time
		refreshTokenExpiresAt           time.Time
		authStorageFails                bool
		accessTokenInvalid              bool
		refreshTokenInvalid             bool
		authorizeUserFails              bool
		refreshTokensFails              bool
		refreshTokensReturnsNonOKStatus bool
		// Expected outcome settings
		isValid                      bool
		expectedReautorizeUserCalled bool
		expectedTokensRefreshed      bool
	}{
		{
			desc:                         "happy path",
			accessTokenExpiresAt:         time.Now().Add(time.Hour),
			refreshTokenExpiresAt:        time.Now().Add(time.Hour),
			isValid:                      true,
			expectedReautorizeUserCalled: false,
			expectedTokensRefreshed:      false,
		},
		{
			desc:                         "use access token",
			accessTokenExpiresAt:         time.Now().Add(time.Hour),
			refreshTokenExpiresAt:        time.Now().Add(-time.Hour),
			isValid:                      true,
			expectedReautorizeUserCalled: false,
			expectedTokensRefreshed:      false,
		},
		{
			desc:                         "use refresh token",
			accessTokenExpiresAt:         time.Now().Add(-time.Hour),
			refreshTokenExpiresAt:        time.Now().Add(time.Hour),
			isValid:                      true,
			expectedReautorizeUserCalled: false,
			expectedTokensRefreshed:      true,
		},
		{
			desc:                         "tokens expired",
			accessTokenExpiresAt:         time.Now().Add(-time.Hour),
			refreshTokenExpiresAt:        time.Now().Add(-time.Hour),
			isValid:                      true,
			expectedReautorizeUserCalled: true,
			expectedTokensRefreshed:      true,
		},
		{
			desc:                         "auth storage fails",
			accessTokenExpiresAt:         time.Now().Add(time.Hour),
			refreshTokenExpiresAt:        time.Now().Add(time.Hour),
			authStorageFails:             true,
			isValid:                      false,
			expectedReautorizeUserCalled: false,
			expectedTokensRefreshed:      false,
		},
		{
			desc:                         "access token invalid",
			accessTokenExpiresAt:         time.Now().Add(time.Hour),
			refreshTokenExpiresAt:        time.Now().Add(time.Hour),
			accessTokenInvalid:           true,
			isValid:                      false,
			expectedReautorizeUserCalled: false,
			expectedTokensRefreshed:      false,
		},
		{
			desc:                         "refresh token invalid",
			accessTokenExpiresAt:         time.Now().Add(-time.Hour),
			refreshTokenExpiresAt:        time.Now().Add(time.Hour),
			refreshTokenInvalid:          true,
			isValid:                      false,
			expectedReautorizeUserCalled: false,
			expectedTokensRefreshed:      false,
		},
		{
			desc:                         "refresh token invalid but unused",
			accessTokenExpiresAt:         time.Now().Add(time.Hour),
			refreshTokenExpiresAt:        time.Now().Add(time.Hour),
			refreshTokenInvalid:          true,
			isValid:                      true,
			expectedReautorizeUserCalled: false,
			expectedTokensRefreshed:      false,
		},
		{
			desc:                         "authorize user fails",
			accessTokenExpiresAt:         time.Now().Add(-time.Hour),
			refreshTokenExpiresAt:        time.Now().Add(-time.Hour),
			authorizeUserFails:           true,
			isValid:                      false,
			expectedReautorizeUserCalled: true,
			expectedTokensRefreshed:      false,
		},
		{
			desc:                         "refresh tokens fails",
			accessTokenExpiresAt:         time.Now().Add(-time.Hour),
			refreshTokenExpiresAt:        time.Now().Add(time.Hour),
			refreshTokensFails:           true,
			isValid:                      true,
			expectedReautorizeUserCalled: true,
			expectedTokensRefreshed:      true,
		},
		{
			desc:                            "refresh tokens non OK",
			accessTokenExpiresAt:            time.Now().Add(-time.Hour),
			refreshTokenExpiresAt:           time.Now().Add(time.Hour),
			refreshTokensReturnsNonOKStatus: true,
			isValid:                         true,
			expectedReautorizeUserCalled:    true,
			expectedTokensRefreshed:         true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			// Setup auth storage
			if tt.authStorageFails {
				keyring.MockInitWithError(fmt.Errorf("failed"))
			} else {
				keyring.MockInit()
				err := setAuthStorage(
					tt.accessTokenExpiresAt,
					tt.refreshTokenExpiresAt,
					tt.accessTokenInvalid,
					tt.refreshTokenInvalid,
				)
				if err != nil {
					t.Fatalf("failed to set auth storage: %v", err)
				}
			}

			// Setup transport and authorizeUser
			requestSent := false
			authorizeUserCalled := false
			tokensRefreshed := false
			refreshTokensTransport := &clientTransport{
				t:                               t,
				requestURL:                      "request/url",
				refreshTokensFails:              tt.refreshTokensFails,
				refreshTokensReturnsNonOKStatus: tt.refreshTokensReturnsNonOKStatus,
				requestSent:                     &requestSent,
				tokensRefreshed:                 &tokensRefreshed,
			}
			client := &http.Client{
				Transport: refreshTokensTransport,
			}
			authorizeUserContext := &authorizeUserContext{
				t:                   t,
				authorizeUserFails:  tt.authorizeUserFails,
				authorizeUserCalled: &authorizeUserCalled,
				tokensRefreshed:     &tokensRefreshed,
			}
			authorizeUserRoutine := func() error {
				return reauthorizeUser(authorizeUserContext)
			}

			cmd := &cobra.Command{}
			cmd.SetOut(io.Discard) // Suppresses console prints

			p := &print.Printer{Cmd: cmd}

			// Test RoundTripper
			rt := userTokenFlow{
				printer:                p,
				reauthorizeUserRoutine: authorizeUserRoutine,
				client:                 client,
			}
			req, err := http.NewRequest(http.MethodGet, "request/url", http.NoBody)
			if err != nil {
				t.Fatalf("failed to create request: %v", err)
			}

			resp, err := rt.RoundTrip(req)
			if err == nil {
				defer func() {
					tempErr := resp.Body.Close()
					if tempErr != nil {
						t.Fatalf("failed to close response body: %v", tempErr)
					}
				}()
			}

			if !tt.isValid && err == nil {
				if err == nil {
					t.Errorf("should have failed")
				}
				if requestSent {
					t.Errorf("request was sent")
				}
			}
			if tt.isValid && err != nil {
				if err != nil {
					t.Errorf("shouldn't have failed: %v", err)
				}
				if !requestSent {
					t.Errorf("request wasn't sent")
				}
			}
			if authorizeUserCalled && !tt.expectedReautorizeUserCalled {
				t.Errorf("reauthorizeUser was called")
			}
			if !authorizeUserCalled && tt.expectedReautorizeUserCalled {
				t.Errorf("reauthorizeUser wasn't called")
			}
			if tokensRefreshed && !tt.expectedTokensRefreshed {
				t.Errorf("tokens were refreshed")
			}
			if !tokensRefreshed && tt.expectedTokensRefreshed {
				t.Errorf("tokens weren't refreshed")
			}
		})
	}
}

// Generates access and refresh tokens with the expiration timestamp provided, then sets the auth fields in storage appropriately
func setAuthStorage(accessTokenExpiresAt, refreshTokenExpiresAt time.Time, accessTokenInvalid, refreshTokenInvalid bool) error {
	accessToken, refreshToken, err := createTokens(accessTokenExpiresAt, refreshTokenExpiresAt)
	if err != nil {
		return fmt.Errorf("create tokens: %w", err)
	}
	if accessTokenInvalid {
		accessToken = "foo.bar.baz" //nolint:gosec // Hardcoded bad credentials
	}
	if refreshTokenInvalid {
		refreshToken = "foo.bar.baz" //nolint:gosec // Hardcoded bad credentials
	}

	err = SetAuthFlow(AUTH_FLOW_USER_TOKEN)
	if err != nil {
		return fmt.Errorf("set auth flow type: %w", err)
	}
	err = SetAuthFieldMap(map[authFieldKey]string{
		ACCESS_TOKEN:  accessToken,
		REFRESH_TOKEN: refreshToken,
	})
	if err != nil {
		return fmt.Errorf("set refreshed tokens in auth storage: %w", err)
	}
	return nil
}

func createTokens(accessTokenExpiresAt, refreshTokenExpiresAt time.Time) (accessToken, refreshToken string, err error) {
	accessToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(accessTokenExpiresAt),
	}).SignedString([]byte("test"))
	if err != nil {
		return "", "", fmt.Errorf("create access token: %w", err)
	}

	refreshToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(refreshTokenExpiresAt),
	}).SignedString([]byte("test"))
	if err != nil {
		return "", "", fmt.Errorf("create refresh token: %w", err)
	}

	return accessToken, refreshToken, nil
}
