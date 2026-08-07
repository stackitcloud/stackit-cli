package login

import (
	"context"
	"crypto/sha256"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"
	clientauthenticationv1 "k8s.io/client-go/pkg/apis/clientauthentication/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/auth/exec"
	"k8s.io/client-go/tools/clientcmd"

	sdkAuth "github.com/stackitcloud/stackit-sdk-go/core/auth"
	"github.com/stackitcloud/stackit-sdk-go/core/clients"
	sdkConfig "github.com/stackitcloud/stackit-sdk-go/core/config"
	ske "github.com/stackitcloud/stackit-sdk-go/services/ske/v2api"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/auth"
	"github.com/stackitcloud/stackit-cli/internal/pkg/cache"
	cliErr "github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/ske/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
)

const (
	expirationSeconds          = 30 * 60          // 30 min
	refreshBeforeDuration      = 15 * time.Minute // 15 min
	refreshTokenBeforeDuration = 5 * time.Minute  // 5 min

	idpFlag          = "idp"
	clusterNameFlag  = "cluster-name"
	organizationFlag = "organization-id"

	envAccessToken            = "STACKIT_ACCESS_TOKEN"
	envServiceAccountEmail    = "STACKIT_SERVICE_ACCOUNT_EMAIL"
	defaultFederatedTokenPath = "/var/run/secrets/stackit.cloud/serviceaccount/token" //nolint:gosec // Public path, not a credential.
)

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "login",
		Short: "Login plugin for kubernetes clients",
		Long: fmt.Sprintf("%s\n%s\n%s",
			"Login plugin for kubernetes clients, that creates short-lived credentials to authenticate against a STACKIT Kubernetes Engine (SKE) cluster.",
			"First you need to obtain a kubeconfig for use with the login command (first or second example).",
			"Secondly you use the kubeconfig with your chosen Kubernetes client (third example), the client will automatically retrieve the credentials via the STACKIT CLI.",
		),
		Args: args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Get an admin, login kubeconfig for the SKE cluster with name "my-cluster". `+
					"This kubeconfig does not contain any credentials and instead obtains valid admin credentials via the `stackit ske kubeconfig login` command.",
				"$ stackit ske kubeconfig create my-cluster --login"),
			examples.NewExample(
				`Get an IDP kubeconfig for the SKE cluster with name "my-cluster". `+
					"This kubeconfig does not contain any credentials and instead obtains valid credentials via the `stackit ske kubeconfig login` command.",
				"$ stackit ske kubeconfig create my-cluster --idp"),
			examples.NewExample(
				"Use the previously saved kubeconfig to authenticate to the SKE cluster, in this case with kubectl.",
				"$ kubectl cluster-info",
				"$ kubectl get pods"),
			examples.NewExample(
				"Configure an IdP exec provider for a Kubernetes client that does not provide cluster information. In an SKE workload, the CLI automatically uses the projected workload identity token.",
				"$ stackit ske kubeconfig login --idp --cluster-name my-cluster --organization-id my-organization-id --project-id my-project-id --region eu01"),
		),
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := context.Background()

			env := os.Getenv("KUBERNETES_EXEC_INFO")
			if env == "" {
				return fmt.Errorf("%s\n%s\n%s", "KUBERNETES_EXEC_INFO env var is unset or empty.",
					"The command probably was not called from a Kubernetes client application!",
					"See `stackit ske kubeconfig login --help` for detailed usage instructions.")
			}

			idpMode := flags.FlagToBoolValue(params.Printer, cmd, idpFlag)
			workloadIdentityMode := idpMode && os.Getenv(envAccessToken) == "" && workloadIdentityConfigured()
			statelessIDPMode := idpMode && (workloadIdentityMode || os.Getenv(envAccessToken) != "")
			if !statelessIDPMode {
				if err := cache.Init(); err != nil {
					return fmt.Errorf("cache init failed: %w", err)
				}
			}
			clusterConfig, err := parseClusterConfig(params.Printer, cmd, idpMode, workloadIdentityMode, statelessIDPMode)
			if err != nil {
				return fmt.Errorf("parseClusterConfig: %w", err)
			}

			if idpMode {
				accessToken, err := getAccessToken(params, workloadIdentityMode)
				if err != nil {
					return err
				}
				tokenEndpoint, err := auth.GetIDPTokenEndpoint(params.Printer)
				if err != nil {
					return fmt.Errorf("get IDP token endpoint: %w", err)
				}
				idpClient := &http.Client{}
				token, err := retrieveTokenFromIDP(ctx, idpClient, tokenEndpoint, accessToken, clusterConfig, !statelessIDPMode)
				if err != nil {
					return err
				}
				return outputTokenKubeconfig(params.Printer, clusterConfig.cacheKey, token, !statelessIDPMode)
			}

			// Configure API client
			apiClient, err := client.ConfigureClient(params.Printer, params.CliVersion)
			if err != nil {
				return err
			}
			kubeconfig, err := retrieveLoginKubeconfig(ctx, apiClient, clusterConfig)
			if err != nil {
				return err
			}
			return outputLoginKubeconfig(params.Printer, clusterConfig.cacheKey, kubeconfig)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Bool(idpFlag, false, "Use the STACKIT IdP for authentication to the cluster.")
	cmd.Flags().String(clusterNameFlag, "", "SKE cluster name. Used when the Kubernetes exec request does not provide cluster information.")
	cmd.Flags().String(organizationFlag, "", "Organization ID. Used in IdP mode when the Kubernetes exec request does not provide cluster information.")
}

type clusterConfig struct {
	STACKITProjectID string `json:"stackitProjectID"`
	ClusterName      string `json:"clusterName"`
	Region           string `json:"region"`
	OrganizationID   string `json:"organizationID"`

	cacheKey string
}

func parseClusterConfig(p *print.Printer, cmd *cobra.Command, idpMode, workloadIdentityMode, statelessIDPMode bool) (*clusterConfig, error) {
	obj, err := loadExecCredentialFromEnv()
	if err != nil {
		return nil, fmt.Errorf("LoadExecCredentialFromEnv: %w", err)
	}

	if err := clientauthenticationv1.AddToScheme(scheme.Scheme); err != nil {
		return nil, err
	}

	obj, err = scheme.Scheme.ConvertToVersion(obj, clientauthenticationv1.SchemeGroupVersion)
	if err != nil {
		return nil, fmt.Errorf("ConvertToVersion: %w", err)
	}

	execCredential, ok := obj.(*clientauthenticationv1.ExecCredential)
	if !ok {
		return nil, fmt.Errorf("conversion to ExecCredential failed")
	}
	if execCredential == nil {
		return nil, fmt.Errorf("ExecCredential is empty")
	}
	clusterConfig := &clusterConfig{}
	clusterServer := ""
	if execCredential.Spec.Cluster != nil {
		clusterServer = execCredential.Spec.Cluster.Server
		if len(execCredential.Spec.Cluster.Config.Raw) > 0 {
			err = json.Unmarshal(execCredential.Spec.Cluster.Config.Raw, clusterConfig)
			if err != nil {
				return nil, fmt.Errorf("unmarshal: %w", err)
			}
		}
	}

	if clusterName := flags.FlagToStringValue(p, cmd, clusterNameFlag); clusterName != "" {
		clusterConfig.ClusterName = clusterName
	}
	if organizationID := flags.FlagToStringValue(p, cmd, organizationFlag); organizationID != "" {
		clusterConfig.OrganizationID = organizationID
	}
	globalFlags := globalflags.Parse(p, cmd)
	if clusterConfig.STACKITProjectID == "" {
		clusterConfig.STACKITProjectID = globalFlags.ProjectId
	}
	if clusterConfig.Region == "" {
		clusterConfig.Region = globalFlags.Region
	}

	missingFields := missingClusterConfigFields(clusterConfig, idpMode)
	if len(missingFields) > 0 {
		return nil, fmt.Errorf("ExecCredential cluster configuration is incomplete; provide %s", strings.Join(missingFields, ", "))
	}

	authIdentity, err := getAuthIdentity(workloadIdentityMode, statelessIDPMode)
	if err != nil {
		return nil, fmt.Errorf("error getting auth identity: %w", err)
	}
	idpSuffix := ""
	if idpMode {
		idpSuffix = "\x00idp"
	}
	clusterIdentity := clusterServer
	if clusterIdentity == "" {
		clusterIdentity = strings.Join([]string{
			clusterConfig.OrganizationID,
			clusterConfig.STACKITProjectID,
			clusterConfig.Region,
			clusterConfig.ClusterName,
		}, "\x00")
	}
	clusterConfig.cacheKey = fmt.Sprintf("ske-login-%x", sha256.Sum256([]byte(clusterIdentity+"\x00"+authIdentity+idpSuffix)))

	return clusterConfig, nil
}

func loadExecCredentialFromEnv() (runtime.Object, error) {
	execInfo := os.Getenv("KUBERNETES_EXEC_INFO")
	if execInfo == "" {
		return nil, errors.New("KUBERNETES_EXEC_INFO env var is unset or empty")
	}

	// client-go's loader rejects requests without cluster information. Decode the
	// v1 request directly when the Kubernetes client omits it.
	var execCredential clientauthenticationv1.ExecCredential
	if err := json.Unmarshal([]byte(execInfo), &execCredential); err == nil &&
		execCredential.APIVersion == clientauthenticationv1.SchemeGroupVersion.String() &&
		execCredential.Kind == "ExecCredential" && execCredential.Spec.Cluster == nil {
		return &execCredential, nil
	}

	obj, _, err := exec.LoadExecCredential([]byte(execInfo))
	if err != nil {
		return nil, err
	}
	return obj, nil
}

func missingClusterConfigFields(config *clusterConfig, idpMode bool) []string {
	missingFields := make([]string, 0, 4)
	if config.ClusterName == "" {
		missingFields = append(missingFields, "--cluster-name")
	}
	if config.STACKITProjectID == "" {
		missingFields = append(missingFields, "--project-id")
	}
	if config.Region == "" {
		missingFields = append(missingFields, "--region")
	}
	if idpMode && config.OrganizationID == "" {
		missingFields = append(missingFields, "--organization-id")
	}
	return missingFields
}

func getAuthIdentity(workloadIdentityMode, statelessIDPMode bool) (string, error) {
	if workloadIdentityMode {
		email := os.Getenv(envServiceAccountEmail)
		if email == "" {
			return "", fmt.Errorf("%s is not set", envServiceAccountEmail)
		}
		return email, nil
	}
	if statelessIDPMode {
		return "stateless-idp", nil
	}
	return auth.GetAuthEmail()
}

func retrieveLoginKubeconfig(ctx context.Context, apiClient *ske.APIClient, clusterConfig *clusterConfig) (*rest.Config, error) {
	cachedKubeconfig := getCachedKubeConfig(clusterConfig.cacheKey)
	if cachedKubeconfig == nil {
		return requestNewLoginKubeconfig(ctx, apiClient, clusterConfig)
	}

	isValid, notAfter := checkKubeconfigExpiry(cachedKubeconfig.CertData)
	if !isValid {
		// cert is expired or invalid, request new
		_ = cache.DeleteObject(clusterConfig.cacheKey)
		return requestNewLoginKubeconfig(ctx, apiClient, clusterConfig)
	} else if time.Now().Add(refreshBeforeDuration).After(notAfter.UTC()) {
		// cert expires within the next 15min -> refresh
		kubeconfig, err := requestNewLoginKubeconfig(ctx, apiClient, clusterConfig)
		// try to get a new one but use cache on failure
		if err != nil {
			return cachedKubeconfig, nil
		}
		return kubeconfig, nil
	}
	// cert not expired, nor will it expire in the next 15min; therefore, use the cached kubeconfig
	return cachedKubeconfig, nil
}

func getCachedKubeConfig(key string) *rest.Config {
	cachedKubeconfig, err := cache.GetObject(key)
	if err != nil {
		return nil
	}

	restConfig, err := clientcmd.RESTConfigFromKubeConfig(cachedKubeconfig)
	if err != nil {
		return nil
	}

	return restConfig
}

func checkKubeconfigExpiry(certData []byte) (bool, time.Time) {
	certPem, _ := pem.Decode(certData)
	if certPem == nil {
		return false, time.Time{}
	}

	certificate, err := x509.ParseCertificate(certPem.Bytes)
	if err != nil {
		return false, time.Time{}
	}

	// cert is expired
	if time.Now().After(certificate.NotAfter.UTC()) {
		return false, time.Time{}
	}
	return true, certificate.NotAfter.UTC()
}

func requestNewLoginKubeconfig(ctx context.Context, apiClient *ske.APIClient, clusterConfig *clusterConfig) (*rest.Config, error) {
	req := buildLoginKubeconfigRequest(ctx, apiClient, clusterConfig)
	kubeconfigResponse, err := req.Execute()
	if err != nil {
		return nil, fmt.Errorf("request kubeconfig: %w", err)
	}
	kubeconfig, err := clientcmd.RESTConfigFromKubeConfig([]byte(*kubeconfigResponse.Kubeconfig))
	if err != nil {
		return nil, fmt.Errorf("parse kubeconfig: %w", err)
	}
	if err = cache.PutObject(clusterConfig.cacheKey, []byte(*kubeconfigResponse.Kubeconfig)); err != nil {
		return nil, fmt.Errorf("cache kubeconfig: %w", err)
	}

	return kubeconfig, nil
}

func buildLoginKubeconfigRequest(ctx context.Context, apiClient *ske.APIClient, clusterConfig *clusterConfig) ske.ApiCreateKubeconfigRequest {
	req := apiClient.DefaultAPI.CreateKubeconfig(ctx, clusterConfig.STACKITProjectID, clusterConfig.Region, clusterConfig.ClusterName)
	expirationSeconds := strconv.Itoa(expirationSeconds)

	return req.CreateKubeconfigPayload(ske.CreateKubeconfigPayload{ExpirationSeconds: &expirationSeconds})
}

func outputLoginKubeconfig(p *print.Printer, cacheKey string, kubeconfig *rest.Config) error {
	output, err := parseLoginKubeConfigToExecCredential(kubeconfig)
	if err != nil {
		_ = cache.DeleteObject(cacheKey)
		return fmt.Errorf("convert to ExecCredential: %w", err)
	}

	p.Outputf("%s", string(output))
	return nil
}

func parseLoginKubeConfigToExecCredential(kubeconfig *rest.Config) ([]byte, error) {
	if kubeconfig == nil {
		return nil, errors.New("kubeconfig is nil")
	}

	certPem, _ := pem.Decode(kubeconfig.CertData)
	if certPem == nil {
		return nil, fmt.Errorf("decoded pem is nil")
	}

	certificate, err := x509.ParseCertificate(certPem.Bytes)
	if err != nil {
		return nil, fmt.Errorf("parse certificate: %w", err)
	}

	outputExecCredential := clientauthenticationv1.ExecCredential{
		TypeMeta: v1.TypeMeta{
			APIVersion: clientauthenticationv1.SchemeGroupVersion.String(),
			Kind:       "ExecCredential",
		},
		Status: &clientauthenticationv1.ExecCredentialStatus{
			ExpirationTimestamp:   &v1.Time{Time: certificate.NotAfter.Add(-refreshBeforeDuration)},
			ClientCertificateData: string(kubeconfig.CertData),
			ClientKeyData:         string(kubeconfig.KeyData),
		},
	}

	output, err := json.Marshal(outputExecCredential)
	if err != nil {
		return nil, fmt.Errorf("marshal: %w", err)
	}
	return output, nil
}

func getAccessToken(params *types.CmdParams, workloadIdentityMode bool) (string, error) {
	if workloadIdentityMode {
		accessToken, err := getWorkloadIdentityAccessToken()
		if err != nil {
			params.Printer.Debug(print.ErrorLevel, "get workload identity access token: %v", err)
			return "", fmt.Errorf("get workload identity access token: %w", err)
		}
		return accessToken, nil
	}

	if accessToken := os.Getenv(envAccessToken); accessToken != "" {
		return accessToken, nil
	}

	userSessionExpired, err := auth.UserSessionExpired()
	if err != nil {
		return "", err
	}
	if userSessionExpired {
		return "", &cliErr.SessionExpiredError{}
	}

	accessToken, err := auth.GetValidAccessToken(params.Printer)
	if err != nil {
		params.Printer.Debug(print.ErrorLevel, "get valid access token: %v", err)
		return "", &cliErr.SessionExpiredError{}
	}

	return accessToken, nil
}

func workloadIdentityConfigured() bool {
	if os.Getenv(envServiceAccountEmail) == "" {
		return false
	}
	if os.Getenv(clients.FederatedTokenFileEnv) != "" {
		return true
	}
	fileInfo, err := os.Stat(defaultFederatedTokenPath)
	return err == nil && !fileInfo.IsDir()
}

func getWorkloadIdentityAccessToken() (string, error) {
	roundTripper, err := sdkAuth.SetupAuth(&sdkConfig.Configuration{WorkloadIdentityFederation: true})
	if err != nil {
		return "", fmt.Errorf("configure workload identity federation: %w", err)
	}
	flow, ok := roundTripper.(interface {
		GetAccessToken() (string, error)
	})
	if !ok {
		return "", errors.New("configured authentication flow does not provide access tokens")
	}
	return flow.GetAccessToken()
}

func retrieveTokenFromIDP(ctx context.Context, idpClient *http.Client, tokenEndpoint, accessToken string, clusterConfig *clusterConfig, cacheEnabled bool) (string, error) {
	resource := resourceForCluster(clusterConfig)
	if !cacheEnabled {
		return auth.ExchangeTokenWithEndpoint(ctx, idpClient, tokenEndpoint, accessToken, resource)
	}

	cachedToken := getCachedToken(clusterConfig.cacheKey)
	if cachedToken == "" {
		return exchangeAndCacheToken(ctx, idpClient, tokenEndpoint, accessToken, resource, clusterConfig.cacheKey)
	}

	expiry, err := auth.TokenExpirationTime(cachedToken)
	if err != nil {
		// token is expired or invalid, request new
		_ = cache.DeleteObject(clusterConfig.cacheKey)
		return exchangeAndCacheToken(ctx, idpClient, tokenEndpoint, accessToken, resource, clusterConfig.cacheKey)
	} else if time.Now().Add(refreshTokenBeforeDuration).After(expiry) {
		// token expires soon -> refresh
		token, err := exchangeAndCacheToken(ctx, idpClient, tokenEndpoint, accessToken, resource, clusterConfig.cacheKey)
		// try to get a new one but use cache on failure
		if err != nil {
			return cachedToken, nil
		}
		return token, nil
	}
	// cached token is valid and won't expire soon
	return cachedToken, nil
}

func resourceForCluster(config *clusterConfig) string {
	return fmt.Sprintf(
		"resource://organizations/%s/projects/%s/regions/%s/ske/%s",
		config.OrganizationID,
		config.STACKITProjectID,
		config.Region,
		config.ClusterName,
	)
}

func getCachedToken(key string) string {
	token, err := cache.GetObject(key)
	if err != nil {
		return ""
	}
	return string(token)
}

func exchangeAndCacheToken(ctx context.Context, idpClient *http.Client, tokenEndpoint, accessToken, resource, cacheKey string) (string, error) {
	clusterToken, err := auth.ExchangeTokenWithEndpoint(ctx, idpClient, tokenEndpoint, accessToken, resource)
	if err != nil {
		return "", err
	}
	if err = cache.PutObject(cacheKey, []byte(clusterToken)); err != nil {
		return "", fmt.Errorf("cache token: %w", err)
	}
	return clusterToken, err
}

func outputTokenKubeconfig(p *print.Printer, cacheKey, token string, cacheEnabled bool) error {
	output, err := parseTokenToExecCredential(token)
	if err != nil {
		if cacheEnabled {
			_ = cache.DeleteObject(cacheKey)
		}
		return fmt.Errorf("convert to ExecCredential: %w", err)
	}

	p.Outputf("%s", string(output))
	return nil
}

func parseTokenToExecCredential(clusterToken string) ([]byte, error) {
	expiry, err := auth.TokenExpirationTime(clusterToken)
	if err != nil {
		return nil, fmt.Errorf("parse auth token for cluster: %w", err)
	}

	outputExecCredential := clientauthenticationv1.ExecCredential{
		TypeMeta: v1.TypeMeta{
			APIVersion: clientauthenticationv1.SchemeGroupVersion.String(),
			Kind:       "ExecCredential",
		},
		Status: &clientauthenticationv1.ExecCredentialStatus{
			ExpirationTimestamp: &v1.Time{Time: expiry.Add(-refreshTokenBeforeDuration)},
			Token:               clusterToken,
		},
	}
	output, err := json.Marshal(&outputExecCredential)
	if err != nil {
		return nil, fmt.Errorf("marshal: %w", err)
	}
	return output, nil
}
