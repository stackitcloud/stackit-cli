package utils

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/stackitcloud/stackit-sdk-go/services/ske"
	"golang.org/x/mod/semver"
)

const (
	defaultNodepoolAvailabilityZone = "eu01-3"
	defaultNodepoolCRI              = "containerd"
	defaultNodepoolMachineType      = "b1.2"
	defaultNodepoolMachineImageName = "flatcar"
	defaultNodepoolMaxSurge         = 1
	defaultNodepoolMaximum          = 2
	defaultNodepoolMinimum          = 1
	defaultNodepoolName             = "pool-default"
	defaultNodepoolVolumeType       = "storage_premium_perf2"
	defaultNodepoolVolumeSize       = 50

	supportedState = "supported"
)

type SKEClient interface {
	GetServiceStatusExecute(ctx context.Context, projectId string) (*ske.ProjectResponse, error)
	ListClustersExecute(ctx context.Context, projectId string) (*ske.ListClustersResponse, error)
	ListProviderOptionsExecute(ctx context.Context) (*ske.ProviderOptions, error)
}

func ProjectEnabled(ctx context.Context, apiClient SKEClient, projectId string) (bool, error) {
	project, err := apiClient.GetServiceStatusExecute(ctx, projectId)
	if err != nil {
		return false, fmt.Errorf("get SKE status: %w", err)
	}
	return *project.State == ske.PROJECTSTATE_CREATED, nil
}

func ClusterExists(ctx context.Context, apiClient SKEClient, projectId, clusterName string) (bool, error) {
	clusters, err := apiClient.ListClustersExecute(ctx, projectId)
	if err != nil {
		return false, fmt.Errorf("list SKE clusters: %w", err)
	}
	for _, cl := range *clusters.Items {
		if cl.Name != nil && *cl.Name == clusterName {
			return true, nil
		}
	}
	return false, nil
}

func GetDefaultPayload(ctx context.Context, apiClient SKEClient) (*ske.CreateOrUpdateClusterPayload, error) {
	resp, err := apiClient.ListProviderOptionsExecute(ctx)
	if err != nil {
		return nil, fmt.Errorf("get SKE provider options: %w", err)
	}

	payloadKubernetes, err := getDefaultPayloadKubernetes(resp)
	if err != nil {
		return nil, err
	}
	payloadNodepool, err := getDefaultPayloadNodepool(resp)
	if err != nil {
		return nil, err
	}

	payload := &ske.CreateOrUpdateClusterPayload{
		Extensions: &ske.Extension{
			Acl: &ske.ACL{
				AllowedCidrs: &[]string{},
				Enabled:      utils.Ptr(false),
			},
		},
		Kubernetes: payloadKubernetes,
		Nodepools: &[]ske.Nodepool{
			*payloadNodepool,
		},
	}
	return payload, nil
}

func getDefaultPayloadKubernetes(resp *ske.ProviderOptions) (*ske.Kubernetes, error) {
	output := &ske.Kubernetes{}

	if resp.KubernetesVersions == nil {
		return nil, fmt.Errorf("no supported Kubernetes version found")
	}
	foundKubernetesVersion := false
	versions := *resp.KubernetesVersions
	for i := range versions {
		version := versions[i]
		if *version.State != supportedState {
			continue
		}
		if output.Version != nil {
			oldSemVer := fmt.Sprintf("v%s", *output.Version)
			newSemVer := fmt.Sprintf("v%s", *version.Version)
			if semver.Compare(newSemVer, oldSemVer) != 1 {
				continue
			}
		}

		foundKubernetesVersion = true
		output.Version = version.Version
	}
	if !foundKubernetesVersion {
		return nil, fmt.Errorf("no supported Kubernetes version found")
	}
	return output, nil
}

func getDefaultPayloadNodepool(resp *ske.ProviderOptions) (*ske.Nodepool, error) {
	output := &ske.Nodepool{
		AvailabilityZones: &[]string{
			defaultNodepoolAvailabilityZone,
		},
		Cri: &ske.CRI{
			Name: utils.Ptr(defaultNodepoolCRI),
		},
		Machine: &ske.Machine{
			Type: utils.Ptr(defaultNodepoolMachineType),
			Image: &ske.Image{
				Name: utils.Ptr(defaultNodepoolMachineImageName),
			},
		},
		MaxSurge: utils.Ptr(int64(defaultNodepoolMaxSurge)),
		Maximum:  utils.Ptr(int64(defaultNodepoolMaximum)),
		Minimum:  utils.Ptr(int64(defaultNodepoolMinimum)),
		Name:     utils.Ptr(defaultNodepoolName),
		Volume: &ske.Volume{
			Type: utils.Ptr(defaultNodepoolVolumeType),
			Size: utils.Ptr(int64(defaultNodepoolVolumeSize)),
		},
	}

	// Fill in Cri and Machine.Image
	if resp.MachineImages == nil {
		return nil, fmt.Errorf("no supported image versions found")
	}
	foundImageVersion := false
	images := *resp.MachineImages
	for i := range images {
		image := images[i]
		if *image.Name != defaultNodepoolMachineImageName {
			continue
		}
		if image.Versions == nil {
			continue
		}
		versions := *image.Versions
		for j := range versions {
			version := versions[j]
			if *version.State != supportedState {
				continue
			}

			// Check if default CRI is supported
			if version.Cri == nil || len(*version.Cri) == 0 {
				continue
			}
			criSupported := false
			for k := range *version.Cri {
				cri := (*version.Cri)[k]
				if *cri.Name == defaultNodepoolCRI {
					criSupported = true
					break
				}
			}
			if !criSupported {
				continue
			}

			if output.Machine.Image.Version != nil {
				oldSemVer := fmt.Sprintf("v%s", *output.Machine.Image.Version)
				newSemVer := fmt.Sprintf("v%s", *version.Version)
				if semver.Compare(newSemVer, oldSemVer) != 1 {
					continue
				}
			}

			foundImageVersion = true
			output.Machine.Image.Version = version.Version
		}
	}
	if !foundImageVersion {
		return nil, fmt.Errorf("no supported images found")
	}

	return output, nil
}

func ConvertToSeconds(timeStr string) (*string, error) {
	if len(timeStr) < 2 {
		return nil, fmt.Errorf("invalid time format: %s", timeStr)
	}

	unit := timeStr[len(timeStr)-1:]

	valueStr := timeStr[:len(timeStr)-1]
	value, err := strconv.ParseUint(valueStr, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("unable to parse uint: %s", timeStr)
	}

	var multiplier uint64
	switch unit {
	// second
	case "s":
		multiplier = 1
	// minute
	case "m":
		multiplier = 60
	// hour
	case "h":
		multiplier = 60 * 60
	// day
	case "d":
		multiplier = 60 * 60 * 24
	// month, assume 30 days
	case "M":
		multiplier = 60 * 60 * 24 * 30
	default:
		return nil, fmt.Errorf("invalid time format: %s", timeStr)
	}

	result := uint64(value) * multiplier
	return utils.Ptr(strconv.FormatUint(result, 10)), nil
}

func WriteConfigFile(configPath string, data string) error {
	if data == "" {
		return fmt.Errorf("no data to write")
	}

	dir := filepath.Dir(configPath)

	err := os.MkdirAll(dir, 0o700)
	if err != nil {
		return fmt.Errorf("create config directory: %w", err)
	}

	err = os.WriteFile(configPath, []byte(data), 0o600)
	if err != nil {
		return fmt.Errorf("write file: %w", err)
	}
	return nil
}

func GetDefaultKubeconfigLocation() (string, error) {
	userHome, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("get user home directory: %w", err)
	}

	return filepath.Join(userHome, ".kube", "config"), nil
}
