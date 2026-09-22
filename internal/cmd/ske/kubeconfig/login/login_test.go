package login

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	ske "github.com/stackitcloud/stackit-sdk-go/services/ske/v2api"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	clientauthenticationv1 "k8s.io/client-go/pkg/apis/clientauthentication/v1"
	"k8s.io/client-go/rest"

	"github.com/stackitcloud/stackit-cli/internal/pkg/auth"
	"github.com/stackitcloud/stackit-cli/internal/pkg/config"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testparams"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

type testCtxKey struct{}

var testCtx = context.WithValue(context.Background(), testCtxKey{}, "foo")
var testClient = &ske.APIClient{DefaultAPI: &ske.DefaultAPIService{}}
var testProjectId = uuid.NewString()
var testClusterName = "cluster"
var testOrganization = uuid.NewString()

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
	request := testClient.DefaultAPI.CreateKubeconfig(testCtx, testProjectId, testRegion, testClusterName)
	request = request.CreateKubeconfigPayload(ske.CreateKubeconfigPayload{})
	for _, mod := range mods {
		mod(&request)
	}
	return request
}

func setExecCredentialEnv(t *testing.T, cluster *clientauthenticationv1.Cluster) {
	t.Helper()
	execCredential := clientauthenticationv1.ExecCredential{
		TypeMeta: v1.TypeMeta{
			APIVersion: clientauthenticationv1.SchemeGroupVersion.String(),
			Kind:       "ExecCredential",
		},
		Spec: clientauthenticationv1.ExecCredentialSpec{
			Cluster: cluster,
		},
	}
	execCredentialJSON, err := json.Marshal(execCredential)
	if err != nil {
		t.Fatalf("marshal ExecCredential: %v", err)
	}
	t.Setenv("KUBERNETES_EXEC_INFO", string(execCredentialJSON))
}

func TestParseClusterConfigWithoutExecClusterInfo(t *testing.T) {
	viper.Reset()
	t.Cleanup(viper.Reset)
	viper.Set(config.ProjectIdKey, testProjectId)
	viper.Set(config.RegionKey, testRegion)
	t.Setenv(envServiceAccountEmail, "workload@sa.stackit.cloud")
	setExecCredentialEnv(t, nil)

	params := testparams.NewTestParams()
	cmd := &cobra.Command{}
	configureFlags(cmd)
	if err := cmd.Flags().Set(clusterNameFlag, testClusterName); err != nil {
		t.Fatalf("set cluster name flag: %v", err)
	}
	if err := cmd.Flags().Set(organizationFlag, testOrganization); err != nil {
		t.Fatalf("set organization flag: %v", err)
	}

	actual, err := parseClusterConfig(params.Printer, cmd, true, true, true)
	if err != nil {
		t.Fatalf("parse cluster config: %v", err)
	}
	expected := fixtureClusterConfig()
	if diff := cmp.Diff(actual, expected, cmpopts.IgnoreFields(clusterConfig{}, "cacheKey")); diff != "" {
		t.Fatalf("Data does not match: %s", diff)
	}
	if actual.cacheKey == "" {
		t.Fatal("cache key is empty")
	}
}

func TestParseClusterConfigPrecedence(t *testing.T) {
	globalProjectID := uuid.NewString()
	flagProjectID := uuid.NewString()
	globalRegion := "global-region"
	flagRegion := "flag-region"

	tests := []struct {
		name           string
		clusterConfig  bool
		explicitFlags  bool
		expectedID     string
		expectedRegion string
	}{
		{
			name:           "global config without exec cluster config",
			expectedID:     globalProjectID,
			expectedRegion: globalRegion,
		},
		{
			name:           "explicit flags without exec cluster config",
			explicitFlags:  true,
			expectedID:     flagProjectID,
			expectedRegion: flagRegion,
		},
		{
			name:           "exec cluster config overrides global config",
			clusterConfig:  true,
			expectedID:     testProjectId,
			expectedRegion: testRegion,
		},
		{
			name:           "explicit flags override exec cluster config",
			clusterConfig:  true,
			explicitFlags:  true,
			expectedID:     flagProjectID,
			expectedRegion: flagRegion,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			viper.Reset()
			t.Cleanup(viper.Reset)
			t.Setenv(envServiceAccountEmail, "workload@sa.stackit.cloud")

			globalConfig := fmt.Sprintf(`{"project_id":%q,"region":%q}`, globalProjectID, globalRegion)
			viper.SetConfigType("json")
			if err := viper.ReadConfig(strings.NewReader(globalConfig)); err != nil {
				t.Fatalf("read global config: %v", err)
			}

			execConfig := fixtureClusterConfig()
			if !tt.clusterConfig {
				execConfig.STACKITProjectID = ""
				execConfig.Region = ""
			}
			configJSON, err := json.Marshal(execConfig)
			if err != nil {
				t.Fatalf("marshal cluster config: %v", err)
			}
			setExecCredentialEnv(t, &clientauthenticationv1.Cluster{
				Server: "https://api.example.stackit.cloud",
				Config: runtime.RawExtension{Raw: configJSON},
			})

			params := testparams.NewTestParams()
			cmd := &cobra.Command{}
			if err := globalflags.Configure(cmd.Flags()); err != nil {
				t.Fatalf("configure global flags: %v", err)
			}
			configureFlags(cmd)
			if tt.explicitFlags {
				if err := cmd.Flags().Set(globalflags.ProjectIdFlag, flagProjectID); err != nil {
					t.Fatalf("set project ID flag: %v", err)
				}
				if err := cmd.Flags().Set(globalflags.RegionFlag, flagRegion); err != nil {
					t.Fatalf("set region flag: %v", err)
				}
			}

			actual, err := parseClusterConfig(params.Printer, cmd, true, true, true)
			if err != nil {
				t.Fatalf("parse cluster config: %v", err)
			}
			expected := fixtureClusterConfig(func(config *clusterConfig) {
				config.STACKITProjectID = tt.expectedID
				config.Region = tt.expectedRegion
			})
			if diff := cmp.Diff(actual, expected, cmpopts.IgnoreFields(clusterConfig{}, "cacheKey")); diff != "" {
				t.Fatalf("Data does not match: %s", diff)
			}
		})
	}
}

func TestParseClusterConfigPrefersExplicitClusterIdentity(t *testing.T) {
	viper.Reset()
	t.Cleanup(viper.Reset)
	t.Setenv(envServiceAccountEmail, "workload@sa.stackit.cloud")

	configJSON, err := json.Marshal(fixtureClusterConfig())
	if err != nil {
		t.Fatalf("marshal cluster config: %v", err)
	}
	setExecCredentialEnv(t, &clientauthenticationv1.Cluster{
		Server: "https://api.example.stackit.cloud",
		Config: runtime.RawExtension{Raw: configJSON},
	})

	explicitClusterName := "explicit-cluster"
	explicitOrganizationID := uuid.NewString()
	params := testparams.NewTestParams()
	cmd := &cobra.Command{}
	configureFlags(cmd)
	if err := cmd.Flags().Set(clusterNameFlag, explicitClusterName); err != nil {
		t.Fatalf("set cluster name flag: %v", err)
	}
	if err := cmd.Flags().Set(organizationFlag, explicitOrganizationID); err != nil {
		t.Fatalf("set organization flag: %v", err)
	}

	actual, err := parseClusterConfig(params.Printer, cmd, true, true, true)
	if err != nil {
		t.Fatalf("parse cluster config: %v", err)
	}
	expected := fixtureClusterConfig(func(config *clusterConfig) {
		config.ClusterName = explicitClusterName
		config.OrganizationID = explicitOrganizationID
	})
	if diff := cmp.Diff(actual, expected, cmpopts.IgnoreFields(clusterConfig{}, "cacheKey")); diff != "" {
		t.Fatalf("Data does not match: %s", diff)
	}
}

func TestParseClusterConfigUsesExecClusterInfoAsFallback(t *testing.T) {
	viper.Reset()
	t.Cleanup(viper.Reset)
	t.Setenv(envServiceAccountEmail, "workload@sa.stackit.cloud")

	configJSON, err := json.Marshal(fixtureClusterConfig())
	if err != nil {
		t.Fatalf("marshal cluster config: %v", err)
	}
	setExecCredentialEnv(t, &clientauthenticationv1.Cluster{
		Server: "https://api.example.stackit.cloud",
		Config: runtime.RawExtension{Raw: configJSON},
	})

	params := testparams.NewTestParams()
	cmd := &cobra.Command{}
	configureFlags(cmd)
	actual, err := parseClusterConfig(params.Printer, cmd, true, true, true)
	if err != nil {
		t.Fatalf("parse cluster config: %v", err)
	}
	expected := fixtureClusterConfig()
	if diff := cmp.Diff(actual, expected, cmpopts.IgnoreFields(clusterConfig{}, "cacheKey")); diff != "" {
		t.Fatalf("Data does not match: %s", diff)
	}
}

func TestParseClusterConfigReportsMissingExplicitFields(t *testing.T) {
	viper.Reset()
	t.Cleanup(viper.Reset)
	t.Setenv(envServiceAccountEmail, "workload@sa.stackit.cloud")
	setExecCredentialEnv(t, nil)

	params := testparams.NewTestParams()
	cmd := &cobra.Command{}
	configureFlags(cmd)
	_, err := parseClusterConfig(params.Printer, cmd, true, true, true)
	if err == nil {
		t.Fatal("Expected error but no error was returned")
	}
	for _, expectedFlag := range []string{"--cluster-name", "--project-id", "--region", "--organization-id"} {
		if !strings.Contains(err.Error(), expectedFlag) {
			t.Errorf("Expected error to mention %s, got %q", expectedFlag, err)
		}
	}
}

func TestParseClusterConfigWithAccessTokenWithoutStoredAuth(t *testing.T) {
	viper.Reset()
	t.Cleanup(viper.Reset)
	viper.Set(config.ProjectIdKey, testProjectId)
	viper.Set(config.RegionKey, testRegion)
	t.Setenv(envAccessToken, "environment-access-token")
	setExecCredentialEnv(t, nil)

	params := testparams.NewTestParams()
	cmd := &cobra.Command{}
	configureFlags(cmd)
	if err := cmd.Flags().Set(clusterNameFlag, testClusterName); err != nil {
		t.Fatalf("set cluster name flag: %v", err)
	}
	if err := cmd.Flags().Set(organizationFlag, testOrganization); err != nil {
		t.Fatalf("set organization flag: %v", err)
	}

	actual, err := parseClusterConfig(params.Printer, cmd, true, false, true)
	if err != nil {
		t.Fatalf("parse cluster config: %v", err)
	}
	expected := fixtureClusterConfig()
	if diff := cmp.Diff(actual, expected, cmpopts.IgnoreFields(clusterConfig{}, "cacheKey")); diff != "" {
		t.Fatalf("Data does not match: %s", diff)
	}
}

func TestGetAccessTokenFromEnvironmentWithoutStoredSession(t *testing.T) {
	const accessToken = "environment-access-token"
	t.Setenv(envAccessToken, accessToken)

	params := testparams.NewTestParams()
	actual, err := getAccessToken(params.CmdParams, false)
	if err != nil {
		t.Fatalf("get access token: %v", err)
	}
	if actual != accessToken {
		t.Fatalf("Expected access token %q, got %q", accessToken, actual)
	}
}

func TestWorkloadIdentityConfigured(t *testing.T) {
	tests := []struct {
		name      string
		configure func(t *testing.T)
		expected  bool
	}{
		{
			name: "explicit OIDC",
			configure: func(t *testing.T) {
				t.Setenv(auth.EnvUseOIDC, "1")
			},
			expected: true,
		},
		{
			name: "static federated token",
			configure: func(t *testing.T) {
				t.Setenv(auth.EnvServiceAccountFederatedToken, "federated-token")
			},
			expected: true,
		},
		{
			name: "federated token file",
			configure: func(t *testing.T) {
				t.Setenv(auth.EnvFederatedTokenFile, "/projected/token")
			},
			expected: true,
		},
		{
			name: "GitHub Actions",
			configure: func(t *testing.T) {
				t.Setenv(auth.EnvGitHubRequestURL, "https://github.example.test/oidc")
				t.Setenv(auth.EnvGitHubRequestToken, "request-token")
			},
			expected: true,
		},
		{
			name: "Azure DevOps",
			configure: func(t *testing.T) {
				t.Setenv(auth.EnvAzureOIDCRequestURI, "https://azure.example.test/oidc")
				t.Setenv(auth.EnvAzureAccessToken, "access-token")
			},
			expected: true,
		},
		{
			name:      "missing token source",
			configure: func(_ *testing.T) {},
			expected:  false,
		},
		{
			name: "missing service account email",
			configure: func(t *testing.T) {
				t.Setenv(envServiceAccountEmail, "")
				t.Setenv(auth.EnvServiceAccountFederatedToken, "federated-token")
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clearOIDCTokenSourceEnv(t)
			t.Setenv(envServiceAccountEmail, "workload@sa.stackit.cloud")
			tt.configure(t)
			if actual := workloadIdentityConfigured(); actual != tt.expected {
				t.Fatalf("workloadIdentityConfigured() = %t, want %t", actual, tt.expected)
			}
		})
	}
}

func TestGetWorkloadIdentityAccessToken(t *testing.T) {
	const serviceAccountEmail = "workload@sa.stackit.cloud"
	federatedToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(10 * time.Minute)),
	}).SignedString([]byte("federated-token-signing-key"))
	if err != nil {
		t.Fatalf("sign federated token: %v", err)
	}
	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
	}).SignedString([]byte("test-signing-key"))
	if err != nil {
		t.Fatalf("sign access token: %v", err)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		req.Body = http.MaxBytesReader(w, req.Body, 1<<20)
		if err := req.ParseForm(); err != nil {
			t.Errorf("parse form: %v", err)
		}
		expectedForm := url.Values{
			"grant_type":            {"client_credentials"},
			"client_assertion_type": {"urn:schwarz:params:oauth:client-assertion-type:workload-jwt"},
			"client_assertion":      {federatedToken},
			"client_id":             {serviceAccountEmail},
		}
		if diff := cmp.Diff(req.Form, expectedForm); diff != "" {
			t.Errorf("Token request does not match: %s", diff)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprintf(w, `{"access_token":%q,"expires_in":3600,"token_type":"Bearer"}`, accessToken)
	}))
	defer server.Close()

	tests := []struct {
		name      string
		configure func(t *testing.T)
	}{
		{
			name: "static federated token",
			configure: func(t *testing.T) {
				t.Setenv(auth.EnvServiceAccountFederatedToken, federatedToken)
			},
		},
		{
			name: "federated token file",
			configure: func(t *testing.T) {
				tokenPath := t.TempDir() + "/token"
				if err := os.WriteFile(tokenPath, []byte(federatedToken), 0o600); err != nil {
					t.Fatalf("write federated token: %v", err)
				}
				t.Setenv(auth.EnvFederatedTokenFile, tokenPath)
			},
		},
		{
			name: "GitHub Actions",
			configure: func(t *testing.T) {
				oidcServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
					w.Header().Set("Content-Type", "application/json")
					_ = json.NewEncoder(w).Encode(map[string]string{"value": federatedToken})
				}))
				t.Cleanup(oidcServer.Close)
				t.Setenv(auth.EnvGitHubRequestURL, oidcServer.URL)
				t.Setenv(auth.EnvGitHubRequestToken, "request-token")
			},
		},
		{
			name: "Azure DevOps",
			configure: func(t *testing.T) {
				oidcServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
					w.Header().Set("Content-Type", "application/json")
					_ = json.NewEncoder(w).Encode(map[string]string{"oidcToken": federatedToken})
				}))
				t.Cleanup(oidcServer.Close)
				t.Setenv(auth.EnvAzureOIDCRequestURI, oidcServer.URL)
				t.Setenv(auth.EnvAzureAccessToken, "access-token")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clearOIDCTokenSourceEnv(t)
			t.Setenv(envServiceAccountEmail, serviceAccountEmail)
			t.Setenv("STACKIT_IDP_TOKEN_ENDPOINT", server.URL)
			tt.configure(t)

			actual, err := getWorkloadIdentityAccessToken()
			if err != nil {
				t.Fatalf("get workload identity access token: %v", err)
			}
			if actual != accessToken {
				t.Fatalf("Expected access token %q, got %q", accessToken, actual)
			}
		})
	}
}

func clearOIDCTokenSourceEnv(t *testing.T) {
	t.Helper()
	for _, env := range []string{
		auth.EnvUseOIDC,
		auth.EnvServiceAccountFederatedToken,
		auth.EnvFederatedTokenFile,
		auth.EnvGitHubRequestURL,
		auth.EnvGitHubRequestToken,
		auth.EnvAzureOIDCRequestURI,
		auth.EnvAzureAccessToken,
	} {
		t.Setenv(env, "")
	}
}

func TestRetrieveTokenFromIDPWithoutCache(t *testing.T) {
	clusterCredential := uuid.NewString()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprintf(w, `{"access_token":%q}`, clusterCredential)
	}))
	defer server.Close()

	actual, err := retrieveTokenFromIDP(testCtx, server.Client(), server.URL, "access-token", fixtureClusterConfig(), false)
	if err != nil {
		t.Fatalf("retrieve token without cache: %v", err)
	}
	if actual != clusterCredential {
		t.Fatalf("Expected cluster credential %q, got %q", clusterCredential, actual)
	}
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
				cmpopts.EquateComparable(testClient.DefaultAPI),
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

func TestResourceForCluster(t *testing.T) {
	cc := fixtureClusterConfig()
	resource := resourceForCluster(cc)
	// somewhat redundant, but the resource string must not change unexpectedly
	expectedResource := "resource://organizations/" + testOrganization + "/projects/" + testProjectId + "/regions/" + testRegion + "/ske/" + testClusterName
	if resource != expectedResource {
		t.Fatalf("unexpected resource, got %v expected %v", resource, expectedResource)
	}
}
