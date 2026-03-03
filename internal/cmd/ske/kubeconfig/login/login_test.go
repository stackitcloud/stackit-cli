package login

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
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"github.com/stackitcloud/stackit-cli/internal/pkg/cache"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/ske"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clientauthenticationv1 "k8s.io/client-go/pkg/apis/clientauthentication/v1"
	"k8s.io/client-go/rest"
)

type testCtxKey struct{}

var testCtx = context.WithValue(context.Background(), testCtxKey{}, "foo")
var testClient = &ske.APIClient{}
var testProjectId = uuid.NewString()
var testClusterName = "cluster"
var testOrganization = uuid.NewString()
var testAccessToken = "access-token-test-" + uuid.NewString()
var testExchangedToken = "access-token-exchanged-" + uuid.NewString()
var testTokenEndpoint = "https://accounts.stackit.cloud/test/endpoint" //nolint:gosec // Actually just a URL

const testRegion = "eu01"

func fixtureClusterConfig(mods ...func(clusterConfig *clusterConfig)) *clusterConfig {
	clusterConfig := &clusterConfig{
		STACKITProjectID: testProjectId,
		ClusterName:      testClusterName,
		cacheKey:         "",
		Region:           testRegion,
		OrganizationID:   testOrganization,
	}
	for _, mod := range mods {
		mod(clusterConfig)
	}
	return clusterConfig
}

func fixtureLoginRequest(mods ...func(request *ske.ApiCreateKubeconfigRequest)) ske.ApiCreateKubeconfigRequest {
	request := testClient.CreateKubeconfig(testCtx, testProjectId, testRegion, testClusterName)
	request = request.CreateKubeconfigPayload(ske.CreateKubeconfigPayload{})
	for _, mod := range mods {
		mod(&request)
	}
	return request
}

func fixtureTokenExchangeRequest(tokenEndpoint string) *http.Request {
	form := url.Values{}
	form.Set("grant_type", "urn:ietf:params:oauth:grant-type:token-exchange")
	form.Set("client_id", "stackit-cli-0000-0000-000000000001")
	form.Set("subject_token_type", "urn:ietf:params:oauth:token-type:access_token")
	form.Set("requested_token_type", "urn:ietf:params:oauth:token-type:id_token")
	form.Set("scope", "openid profile email groups")
	form.Set("subject_token", testAccessToken)
	form.Set("resource", "resource://organizations/"+testOrganization+"/projects/"+testProjectId+"/regions/"+testRegion+"/ske/"+testClusterName)

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
	response, _ := json.Marshal(exchangeReponse{
		AccessToken:    testExchangedToken,
		IssuedTokeType: "urn:ietf:params:oauth:token-type:id_token",
		TokenType:      "Bearer",
	})
	return string(response)
}

func TestBuildRequest(t *testing.T) {
	tests := []struct {
		description     string
		clusterConfig   *clusterConfig
		expectedRequest ske.ApiCreateKubeconfigRequest
	}{
		{
			description:   "expiration time",
			clusterConfig: fixtureClusterConfig(),
			expectedRequest: fixtureLoginRequest().CreateKubeconfigPayload(ske.CreateKubeconfigPayload{
				ExpirationSeconds: utils.Ptr("1800")}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			request := buildLoginKubeconfigRequest(testCtx, testClient, tt.clusterConfig)

			diff := cmp.Diff(request, tt.expectedRequest,
				cmp.AllowUnexported(tt.expectedRequest),
				cmpopts.EquateComparable(testCtx),
			)
			if diff != "" {
				t.Fatalf("Data does not match: %s", diff)
			}
		})
	}
}

func TestParseKubeConfigToExecCredential(t *testing.T) {
	expectedTime, _ := time.Parse(time.RFC3339, "2024-01-01T00:45:00Z")

	tests := []struct {
		description                   string
		kubeconfig                    *rest.Config
		expectedExecCredentialRequest *clientauthenticationv1.ExecCredential
	}{
		{
			description: "expiration time",
			kubeconfig: &rest.Config{
				TLSClientConfig: rest.TLSClientConfig{
					CertData: []byte(`-----BEGIN CERTIFICATE-----
MIIBhTCCASugAwIBAgIIF8+zRM8UalAwCgYIKoZIzj0EAwIwGDEWMBQGA1UEAxMN
Y2EtY2xpZW50LXh5ejAeFw0yNDAxMDEwMDAwMDBaFw0yNDAxMDEwMTAwMDBaMC8x
FzAVBgNVBAoTDnN5c3RlbTptYXN0ZXJzMRQwEgYDVQQDEwtza2U6Y2x1c3RlcjBZ
MBMGByqGSM49AgEGCCqGSM49AwEHA0IABJaxZ8G4wEZ1xf44hMV1pQWsti5SL6PH
QF0bRniQEJHSOcZMwc0OrVIfuSV1qSMyvYIaFtBj1j9f2v8oPux7V02jSDBGMA4G
A1UdDwEB/wQEAwIFoDATBgNVHSUEDDAKBggrBgEFBQcDAjAfBgNVHSMEGDAWgBQt
Pn1pNgfb8xcdRVxVnHDIvb8abzAKBggqhkjOPQQDAgNIADBFAiEA8gG2l0schbMu
zbRjZmli7cnenEnfnNoFIGbgkbjGXRUCIC5zFtWXFK7kA+B2vDxD0DlLcQodNwi4
2JKP8gT9ol16
-----END CERTIFICATE-----`),
					KeyData: []byte("keykeykey"),
				},
			},
			expectedExecCredentialRequest: &clientauthenticationv1.ExecCredential{
				TypeMeta: v1.TypeMeta{
					APIVersion: clientauthenticationv1.SchemeGroupVersion.String(),
					Kind:       "ExecCredential",
				},
				Status: &clientauthenticationv1.ExecCredentialStatus{
					ExpirationTimestamp: &v1.Time{Time: expectedTime},
					ClientCertificateData: `-----BEGIN CERTIFICATE-----
MIIBhTCCASugAwIBAgIIF8+zRM8UalAwCgYIKoZIzj0EAwIwGDEWMBQGA1UEAxMN
Y2EtY2xpZW50LXh5ejAeFw0yNDAxMDEwMDAwMDBaFw0yNDAxMDEwMTAwMDBaMC8x
FzAVBgNVBAoTDnN5c3RlbTptYXN0ZXJzMRQwEgYDVQQDEwtza2U6Y2x1c3RlcjBZ
MBMGByqGSM49AgEGCCqGSM49AwEHA0IABJaxZ8G4wEZ1xf44hMV1pQWsti5SL6PH
QF0bRniQEJHSOcZMwc0OrVIfuSV1qSMyvYIaFtBj1j9f2v8oPux7V02jSDBGMA4G
A1UdDwEB/wQEAwIFoDATBgNVHSUEDDAKBggrBgEFBQcDAjAfBgNVHSMEGDAWgBQt
Pn1pNgfb8xcdRVxVnHDIvb8abzAKBggqhkjOPQQDAgNIADBFAiEA8gG2l0schbMu
zbRjZmli7cnenEnfnNoFIGbgkbjGXRUCIC5zFtWXFK7kA+B2vDxD0DlLcQodNwi4
2JKP8gT9ol16
-----END CERTIFICATE-----`,
					ClientKeyData: "keykeykey",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			execCredential, err := parseLoginKubeConfigToExecCredential(tt.kubeconfig)
			if err != nil {
				t.Fatalf("func returned error: %s", err)
			}
			if execCredential == nil {
				t.Fatal("execCredential is nil")
			}
			expected, _ := json.Marshal(tt.expectedExecCredentialRequest)
			diff := cmp.Diff(execCredential, expected)
			if diff != "" {
				t.Fatalf("Data does not match: %s", diff)
			}
		})
	}
}

func TestBuildTokenExchangeRequest(t *testing.T) {
	cfg := fixtureClusterConfig()
	expectedRequest := fixtureTokenExchangeRequest(testTokenEndpoint)
	req, err := buildRequestToExchangeTokens(testCtx, testTokenEndpoint, testAccessToken, cfg)
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
	config := fixtureClusterConfig(func(clusterConfig *clusterConfig) {
		clusterConfig.cacheKey = "test-exchange-token-" + uuid.NewString()
	})
	var request *http.Request
	response := fixtureTokenExchangeResponse()
	defer cache.OverwriteCacheDir(t)()
	if err := cache.Init(); err != nil {
		t.Fatalf("cache init failed: %s", err)
	}

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
	idToken, err := exchangeToken(testCtx, server.Client(), server.URL, testAccessToken, config)
	if err != nil {
		t.Fatalf("func returned error: %s", err)
	}
	diff := cmp.Diff(idToken, testExchangedToken)
	if diff != "" {
		t.Fatalf("Exchanged token does not match: %s", diff)
	}
}

func TestParseTokenToExecCredential(t *testing.T) {
	expirationTime := time.Now().Add(30 * time.Minute)
	expectedTime := expirationTime.Add(-5 * time.Minute)
	token, err := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(expirationTime),
	}).SigningString()
	if err != nil {
		t.Fatalf("token generation failed: %v", err)
	}
	token += ".signatureAAA"

	tests := []struct {
		description                   string
		token                         string
		expectedExecCredentialRequest *clientauthenticationv1.ExecCredential
	}{
		{
			description: "expiration time",
			token:       token,
			expectedExecCredentialRequest: &clientauthenticationv1.ExecCredential{
				TypeMeta: v1.TypeMeta{
					APIVersion: clientauthenticationv1.SchemeGroupVersion.String(),
					Kind:       "ExecCredential",
				},
				Status: &clientauthenticationv1.ExecCredentialStatus{
					ExpirationTimestamp: &v1.Time{Time: expectedTime},
					Token:               token,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			execCredential, err := parseTokenToExecCredential(tt.token)
			if err != nil {
				t.Fatalf("func returned error: %s", err)
			}
			if execCredential == nil {
				t.Fatal("execCredential is nil")
			}
			expected, _ := json.Marshal(tt.expectedExecCredentialRequest)
			diff := cmp.Diff(execCredential, expected)
			if diff != "" {
				t.Fatalf("Data does not match: %s", diff)
			}
		})
	}
}
