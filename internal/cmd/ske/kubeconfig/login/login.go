package login

import (
	"context"
	"crypto/sha256"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
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

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/auth"
	"github.com/stackitcloud/stackit-cli/internal/pkg/cache"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/ske/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-sdk-go/services/ske"
)

const (
	expirationSeconds     = 30 * 60          // 30 min
	refreshBeforeDuration = 15 * time.Minute // 15 min
)

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "login",
		Short: "Login plugin for kubernetes clients",
		Long: fmt.Sprintf("%s\n%s\n%s",
			"Login plugin for kubernetes clients, that creates short-lived credentials to authenticate against a STACKIT Kubernetes Engine (SKE) cluster.",
			"First you need to obtain a kubeconfig for use with the login command (first example).",
			"Secondly you use the kubeconfig with your chosen Kubernetes client (second example), the client will automatically retrieve the credentials via the STACKIT CLI.",
		),
		Args: args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Get a login kubeconfig for the SKE cluster with name "my-cluster". `+
					"This kubeconfig does not contain any credentials and instead obtains valid credentials via the `stackit ske kubeconfig login` command.",
				"$ stackit ske kubeconfig create my-cluster --login"),
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

			clusterConfig, err := parseClusterConfig(params.Printer, cmd)
			if err != nil {
				return fmt.Errorf("parseClusterConfig: %w", err)
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
	return cmd
}

type clusterConfig struct {
	STACKITProjectID string `json:"stackitProjectID"`
	ClusterName      string `json:"clusterName"`
	Region           string `json:"region"`

	cacheKey string
}

func parseClusterConfig(p *print.Printer, cmd *cobra.Command) (*clusterConfig, error) {
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

	clusterConfig.cacheKey = fmt.Sprintf("ske-login-%x", sha256.Sum256([]byte(execCredential.Spec.Cluster.Server+"\x00"+authEmail)))

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
	req := apiClient.CreateKubeconfig(ctx, clusterConfig.STACKITProjectID, clusterConfig.Region, clusterConfig.ClusterName)
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
