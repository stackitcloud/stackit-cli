package auth

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"github.com/zalando/go-keyring"
)

type testCtxKey struct{}

var testCtx = context.WithValue(context.Background(), testCtxKey{}, "foo")

var testAccessToken = "access-token-test-" + uuid.NewString()
var testExchangedToken = "access-token-exchanged-" + uuid.NewString()
var testExchangeResource = "resource://for/token/exchange"

func fixtureTokenExchangeRequest(tokenEndpoint string) *http.Request {
	form := url.Values{}
	form.Set("grant_type", "urn:ietf:params:oauth:grant-type:token-exchange")
	form.Set("client_id", "stackit-cli-0000-0000-000000000001")
	form.Set("subject_token_type", "urn:ietf:params:oauth:token-type:access_token")
	form.Set("requested_token_type", "urn:ietf:params:oauth:token-type:id_token")
	form.Set("scope", "openid profile email groups")
	form.Set("subject_token", testAccessToken)
	form.Set("resource", testExchangeResource)

	req, _ := http.NewRequestWithContext(
		testCtx,
		http.MethodPost,
		tokenEndpoint,
		strings.NewReader(form.Encode()),
	)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	return req
}

func fixtureTokenExchangeResponse() string {
	type exchangeReponse struct {
		AccessToken    string `json:"access_token"`
		IssuedTokeType string `json:"issued_token_type"`
		TokenType      string `json:"token_type"`
	}
	response, _ := json.Marshal(exchangeReponse{ //nolint:gosec // just a testcase, no valid credentials
		AccessToken:    testExchangedToken,
		IssuedTokeType: "urn:ietf:params:oauth:token-type:id_token",
		TokenType:      "Bearer",
	})
	return string(response)
}

func TestBuildTokenExchangeRequest(t *testing.T) {
	expectedRequest := fixtureTokenExchangeRequest(testTokenEndpoint)
	req, err := buildRequestToExchangeTokens(testCtx, testTokenEndpoint, testAccessToken, testExchangeResource)
	if err != nil {
		t.Fatalf("func returned error: %s", err)
	}
	// directly using cmp.Diff is not possible, so dump the requests first
	expected, err := httputil.DumpRequest(expectedRequest, true)
	if err != nil {
		t.Fatalf("fail to dump expected: %s", err)
	}
	actual, err := httputil.DumpRequest(req, true)
	if err != nil {
		t.Fatalf("fail to dump actual: %s", err)
	}
	diff := cmp.Diff(actual, expected)
	if diff != "" {
		t.Fatalf("Data does not match: %s", diff)
	}
}

func TestParseTokenExchangeResponse(t *testing.T) {
	response := fixtureTokenExchangeResponse()

	tests := []struct {
		description string
		response    string
		status      int
		expectError bool
	}{
		{
			description: "valid response",
			response:    response,
			status:      http.StatusOK,
		},
		{
			description: "error status",
			response:    response, // valid response to make sure the status code is checked
			status:      http.StatusForbidden,
			expectError: true,
		},
		{
			description: "error content",
			response:    "{}",
			status:      http.StatusOK,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			w := httptest.NewRecorder()
			w.WriteHeader(tt.status)
			_, _ = w.WriteString(tt.response)
			resp := w.Result()

			defer func() {
				tempErr := resp.Body.Close()
				if tempErr != nil {
					t.Fatalf("failed to close response body: %v", tempErr)
				}
			}()
			accessToken, err := parseTokenExchangeResponse(resp)
			if tt.expectError {
				if err == nil {
					t.Fatal("expected error got nil")
				}
			} else {
				if err != nil {
					t.Fatalf("func returned error: %s", err)
				}
				diff := cmp.Diff(accessToken, testExchangedToken)
				if diff != "" {
					t.Fatalf("Token does not match: %s", diff)
				}
			}
		})
	}
}

func TestExchangeToken(t *testing.T) {
	var request *http.Request
	response := fixtureTokenExchangeResponse()

	handler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// only compare body as the headers will differ
		expected, err := io.ReadAll(request.Body)
		if err != nil {
			t.Errorf("fail to dump expected: %s", err)
		}
		actual, err := io.ReadAll(req.Body)
		if err != nil {
			t.Errorf("fail to dump actual: %s", err)
		}
		diff := cmp.Diff(actual, expected)
		if diff != "" {
			w.WriteHeader(http.StatusBadRequest)
			t.Errorf("request mismatch: %v", diff)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write([]byte(response))
		if err != nil {
			t.Errorf("Failed to write response: %v", err)
		}
	})
	server := httptest.NewServer(handler)
	defer server.Close()

	request = fixtureTokenExchangeRequest(server.URL)
	// use mock keyring to inject the token endpoint URL
	keyring.MockInit()
	err := SetAuthField(IDP_TOKEN_ENDPOINT, server.URL)
	if err != nil {
		t.Errorf("failed to inject idp token endpoint: %s", err)
	}

	idToken, err := ExchangeToken(testCtx, server.Client(), testAccessToken, testExchangeResource)
	if err != nil {
		t.Fatalf("func returned error: %s", err)
	}
	diff := cmp.Diff(idToken, testExchangedToken)
	if diff != "" {
		t.Fatalf("Exchanged token does not match: %s", diff)
	}
}
