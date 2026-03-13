package login

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	ske "github.com/stackitcloud/stackit-sdk-go/services/ske/v2api"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clientauthenticationv1 "k8s.io/client-go/pkg/apis/clientauthentication/v1"
	"k8s.io/client-go/rest"
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
