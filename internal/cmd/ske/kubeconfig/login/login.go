package login

import (
	"context"
	"crypto/sha256"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"strconv"
	"time"

	"github.com/stackitcloud/stackit-cli/internal/pkg/cache"
	"k8s.io/client-go/rest"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/ske/client"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/ske"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	clientauthenticationv1 "k8s.io/client-go/pkg/apis/clientauthentication/v1"
	"k8s.io/client-go/tools/auth/exec"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	expirationSeconds     = 30 * 60          // 30 min
	refreshBeforeDuration = 15 * time.Minute // 15 min
)

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "login",
		Short: "login plugin for kubectl",
		Long:  "login plugin for kubectl to create a short-lived kubeconfig to authenticate against a STACKIT Kubernetes Engine (SKE) cluster. To obtain a kubeconfig for use with the login command, use the 'kubeconfig create' command.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				"login to a SKE cluster specified in the kubeconfig",
				"$ kubectl get pod"),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()

			model, err := parseInput()
			if err != nil {
				return fmt.Errorf("parseInput: %w", err)
			}

			// Configure API client
			apiClient, err := client.ConfigureClient(p)
			if err != nil {
				return err
			}

			cachedKubeconfig := getCachedKubeConfig(model.CacheKey)

			if cachedKubeconfig == nil {
				return GetAndOutputKubeconfig(ctx, cmd, apiClient, model, false, nil)
			}

			certPem, _ := pem.Decode(cachedKubeconfig.CertData)
			if certPem == nil {
				_ = cache.DeleteObject(model.CacheKey)
				return GetAndOutputKubeconfig(ctx, cmd, apiClient, model, false, nil)
			}

			certificate, err := x509.ParseCertificate(certPem.Bytes)
			if err != nil {
				_ = cache.DeleteObject(model.CacheKey)
				return GetAndOutputKubeconfig(ctx, cmd, apiClient, model, false, nil)
			}

			// cert is expired, request new
			if time.Now().After(certificate.NotAfter.UTC()) {
				_ = cache.DeleteObject(model.CacheKey)
				return GetAndOutputKubeconfig(ctx, cmd, apiClient, model, false, nil)
			}
			// cert expires within the next 15min, refresh (try to get a new, use cache on failure)
			if time.Now().Add(refreshBeforeDuration).After(certificate.NotAfter.UTC()) {
				return GetAndOutputKubeconfig(ctx, cmd, apiClient, model, true, cachedKubeconfig)
			}

			// cert not expired, nor will it expire in the next 15min; therefore, use the cached kubeconfig
			if err := output(cmd, model.CacheKey, cachedKubeconfig); err != nil {
				return err
			}
			return nil
		},
	}
	return cmd
}

type inputModel struct {
	ProjectId   string
	ClusterName string
	CacheKey    string
}

type SKEClusterConfig struct {
	STACKITProjectID string `json:"stackitProjectId"`
	ClusterName      string `json:"clusterName"`
}

func parseInput() (*inputModel, error) {
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
	config := &SKEClusterConfig{}
	err = json.Unmarshal(execCredential.Spec.Cluster.Config.Raw, config)
	if err != nil {
		return nil, fmt.Errorf("unmarshal: %w", err)
	}

	return &inputModel{
		ClusterName: config.ClusterName,
		ProjectId:   config.STACKITProjectID,
		CacheKey:    fmt.Sprintf("ske-login-%x", sha256.Sum256([]byte(execCredential.Spec.Cluster.Server))),
	}, nil
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

func GetAndOutputKubeconfig(ctx context.Context, cmd *cobra.Command, apiClient *ske.APIClient, model *inputModel, fallbackToCache bool, cachedKubeconfig *rest.Config) error {
	req := buildRequest(ctx, apiClient, model)
	kubeconfigResponse, err := req.Execute()
	if err != nil {
		if fallbackToCache {
			return output(cmd, model.CacheKey, cachedKubeconfig)
		}
		return fmt.Errorf("request kubeconfig: %w", err)
	}

	kubeconfig, err := clientcmd.RESTConfigFromKubeConfig([]byte(*kubeconfigResponse.Kubeconfig))
	if err != nil {
		if fallbackToCache {
			return output(cmd, model.CacheKey, cachedKubeconfig)
		}
		return fmt.Errorf("parse kubeconfig: %w", err)
	}
	if err = cache.PutObject(model.CacheKey, []byte(*kubeconfigResponse.Kubeconfig)); err != nil {
		if fallbackToCache {
			return output(cmd, model.CacheKey, cachedKubeconfig)
		}
		return fmt.Errorf("cache kubeconfig: %w", err)
	}

	return output(cmd, model.CacheKey, kubeconfig)
}

func buildRequest(ctx context.Context, apiClient *ske.APIClient, model *inputModel) ske.ApiCreateKubeconfigRequest {
	req := apiClient.CreateKubeconfig(ctx, model.ProjectId, model.ClusterName)
	expirationSeconds := strconv.Itoa(expirationSeconds)

	return req.CreateKubeconfigPayload(ske.CreateKubeconfigPayload{ExpirationSeconds: &expirationSeconds})
}

func output(cmd *cobra.Command, cacheKey string, kubeconfig *rest.Config) error {
	outputExecCredential, err := parseKubeConfigToExecCredential(kubeconfig)
	if err != nil {
		_ = cache.DeleteObject(cacheKey)
		return fmt.Errorf("convert to ExecCredential: %w", err)
	}

	output, err := json.Marshal(outputExecCredential)
	if err != nil {
		_ = cache.DeleteObject(cacheKey)
		return fmt.Errorf("marshal ExecCredential: %w", err)
	}

	cmd.Print(string(output))
	return nil
}

func parseKubeConfigToExecCredential(kubeconfig *rest.Config) (*clientauthenticationv1.ExecCredential, error) {
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
			ExpirationTimestamp:   &v1.Time{Time: certificate.NotAfter.Add(-time.Minute * 15)},
			ClientCertificateData: string(kubeconfig.CertData),
			ClientKeyData:         string(kubeconfig.KeyData),
		},
	}
	return &outputExecCredential, nil
}
