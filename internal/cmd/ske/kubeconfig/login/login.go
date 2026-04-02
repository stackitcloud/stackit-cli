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
	"time"

	"github.com/spf13/cobra"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	clientauthenticationv1 "k8s.io/client-go/pkg/apis/clientauthentication/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/auth/exec"
	"k8s.io/client-go/tools/clientcmd"

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

	idpFlag = "idp"
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
		),
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := context.Background()

			if err := cache.Init(); err != nil {
				return fmt.Errorf("cache init failed: %w", err)
			}

			env := os.Getenv("KUBERNETES_EXEC_INFO")
			if env == "" {
				return fmt.Errorf("%s\n%s\n%s", "KUBERNETES_EXEC_INFO env var is unset or empty.",
					"The command probably was not called from a Kubernetes client application!",
					"See `stackit ske kubeconfig login --help` for detailed usage instructions.")
			}

			idpMode := flags.FlagToBoolValue(params.Printer, cmd, idpFlag)
			clusterConfig, err := parseClusterConfig(params.Printer, cmd, idpMode)
			if err != nil {
				return fmt.Errorf("parseClusterConfig: %w", err)
			}

			if idpMode {
				accessToken, err := getAccessToken(params)
				if err != nil {
					return err
				}
				idpClient := &http.Client{}
				token, err := retrieveTokenFromIDP(ctx, idpClient, accessToken, clusterConfig)
				if err != nil {
					return err
				}
				return outputTokenKubeconfig(params.Printer, clusterConfig.cacheKey, token)
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
}

type clusterConfig struct {
	STACKITProjectID string `json:"stackitProjectID"`
	ClusterName      string `json:"clusterName"`
	Region           string `json:"region"`
	OrganizationID   string `json:"organizationID"`

	cacheKey string
}

func parseClusterConfig(p *print.Printer, cmd *cobra.Command, idpMode bool) (*clusterConfig, error) {
	obj, _, err := exec.LoadExecCredentialFromEnv()
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
	if execCredential == nil || execCredential.Spec.Cluster == nil {
		return nil, fmt.Errorf("ExecCredential contains not all needed fields")
	}
	clusterConfig := &clusterConfig{}
	err = json.Unmarshal(execCredential.Spec.Cluster.Config.Raw, clusterConfig)
	if err != nil {
		return nil, fmt.Errorf("unmarshal: %w", err)
	}

	authEmail, err := auth.GetAuthEmail()
	if err != nil {
		return nil, fmt.Errorf("error getting auth email: %w", err)
	}
	idpSuffix := ""
	if idpMode {
		idpSuffix = "\x00idp"
	}
	clusterConfig.cacheKey = fmt.Sprintf("ske-login-%x", sha256.Sum256([]byte(execCredential.Spec.Cluster.Server+"\x00"+authEmail+idpSuffix)))

	// NOTE: Fallback if region is not set in the kubeconfig (this was the case in the past)
	if clusterConfig.Region == "" {
		clusterConfig.Region = globalflags.Parse(p, cmd).Region
	}

	return clusterConfig, nil
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

func getAccessToken(params *types.CmdParams) (string, error) {
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

	err = auth.EnsureIDPTokenEndpoint(params.Printer)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}

func retrieveTokenFromIDP(ctx context.Context, idpClient *http.Client, accessToken string, clusterConfig *clusterConfig) (string, error) {
	resource := resourceForCluster(clusterConfig)

	cachedToken := getCachedToken(clusterConfig.cacheKey)
	if cachedToken == "" {
		return exchangeAndCacheToken(ctx, idpClient, accessToken, resource, clusterConfig.cacheKey)
	}

	expiry, err := auth.TokenExpirationTime(cachedToken)
	if err != nil {
		// token is expired or invalid, request new
		_ = cache.DeleteObject(clusterConfig.cacheKey)
		return exchangeAndCacheToken(ctx, idpClient, accessToken, resource, clusterConfig.cacheKey)
	} else if time.Now().Add(refreshTokenBeforeDuration).After(expiry) {
		// token expires soon -> refresh
		token, err := exchangeAndCacheToken(ctx, idpClient, accessToken, resource, clusterConfig.cacheKey)
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

func exchangeAndCacheToken(ctx context.Context, idpClient *http.Client, accessToken, resource, cacheKey string) (string, error) {
	clusterToken, err := auth.ExchangeToken(ctx, idpClient, accessToken, resource)
	if err != nil {
		return "", err
	}
	if err = cache.PutObject(cacheKey, []byte(clusterToken)); err != nil {
		return "", fmt.Errorf("cache token: %w", err)
	}
	return clusterToken, err
}

func outputTokenKubeconfig(p *print.Printer, cacheKey, token string) error {
	output, err := parseTokenToExecCredential(token)
	if err != nil {
		_ = cache.DeleteObject(cacheKey)
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
