package utils

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"github.com/stackitcloud/stackit-sdk-go/services/ske"
)

var (
	testProjectId = uuid.NewString()
)

const (
	testClusterName = "test-cluster"
)

type skeClientMocked struct {
	getServiceStatusFails    bool
	getServiceStatusResp     *ske.ProjectResponse
	listClustersFails        bool
	listClustersResp         *ske.ListClustersResponse
	listProviderOptionsFails bool
	listProviderOptionsResp  *ske.ProviderOptions
}

func (m *skeClientMocked) GetServiceStatusExecute(_ context.Context, _ string) (*ske.ProjectResponse, error) {
	if m.getServiceStatusFails {
		return nil, fmt.Errorf("could not get service status")
	}
	return m.getServiceStatusResp, nil
}

func (m *skeClientMocked) ListClustersExecute(_ context.Context, _ string) (*ske.ListClustersResponse, error) {
	if m.listClustersFails {
		return nil, fmt.Errorf("could not list clusters")
	}
	return m.listClustersResp, nil
}

func (m *skeClientMocked) ListProviderOptionsExecute(_ context.Context) (*ske.ProviderOptions, error) {
	if m.listProviderOptionsFails {
		return nil, fmt.Errorf("could not list provider options")
	}
	return m.listProviderOptionsResp, nil
}

func TestProjectEnabled(t *testing.T) {
	tests := []struct {
		description     string
		getProjectFails bool
		getProjectResp  *ske.ProjectResponse
		isValid         bool
		expectedOutput  bool
	}{
		{
			description:    "project enabled",
			getProjectResp: &ske.ProjectResponse{State: ske.PROJECTSTATE_CREATED.Ptr()},
			isValid:        true,
			expectedOutput: true,
		},
		{
			description:    "project disabled 1",
			getProjectResp: &ske.ProjectResponse{State: ske.PROJECTSTATE_CREATING.Ptr()},
			isValid:        true,
			expectedOutput: false,
		},
		{
			description:    "project disabled 2",
			getProjectResp: &ske.ProjectResponse{State: ske.PROJECTSTATE_DELETING.Ptr()},
			isValid:        true,
			expectedOutput: false,
		},
		{
			description:     "get clusters fails",
			getProjectFails: true,
			isValid:         false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			client := &skeClientMocked{
				getServiceStatusFails: tt.getProjectFails,
				getServiceStatusResp:  tt.getProjectResp,
			}

			output, err := ProjectEnabled(context.Background(), client, testProjectId)

			if tt.isValid && err != nil {
				t.Errorf("failed on valid input")
			}
			if !tt.isValid && err == nil {
				t.Errorf("did not fail on invalid input")
			}
			if !tt.isValid {
				return
			}
			if output != tt.expectedOutput {
				t.Errorf("expected output to be %t, got %t", tt.expectedOutput, output)
			}
		})
	}
}

func TestClusterExists(t *testing.T) {
	tests := []struct {
		description      string
		getClustersFails bool
		getClustersResp  *ske.ListClustersResponse
		isValid          bool
		expectedExists   bool
	}{
		{
			description:     "cluster exists",
			getClustersResp: &ske.ListClustersResponse{Items: &[]ske.Cluster{{Name: utils.Ptr(testClusterName)}}},
			isValid:         true,
			expectedExists:  true,
		},
		{
			description:     "cluster exists 2",
			getClustersResp: &ske.ListClustersResponse{Items: &[]ske.Cluster{{Name: utils.Ptr("some-cluster")}, {Name: utils.Ptr("some-other-cluster")}, {Name: utils.Ptr(testClusterName)}}},
			isValid:         true,
			expectedExists:  true,
		},
		{
			description:     "cluster does not exist",
			getClustersResp: &ske.ListClustersResponse{Items: &[]ske.Cluster{{Name: utils.Ptr("some-cluster")}, {Name: utils.Ptr("some-other-cluster")}}},
			isValid:         true,
			expectedExists:  false,
		},
		{
			description:      "get clusters fails",
			getClustersFails: true,
			isValid:          false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			client := &skeClientMocked{
				listClustersFails: tt.getClustersFails,
				listClustersResp:  tt.getClustersResp,
			}

			exists, err := ClusterExists(context.Background(), client, testProjectId, testClusterName)

			if tt.isValid && err != nil {
				t.Errorf("failed on valid input")
			}
			if !tt.isValid && err == nil {
				t.Errorf("did not fail on invalid input")
			}
			if !tt.isValid {
				return
			}
			if exists != tt.expectedExists {
				t.Errorf("expected exists to be %t, got %t", tt.expectedExists, exists)
			}
		})
	}
}

func fixtureProviderOptions(mods ...func(*ske.ProviderOptions)) *ske.ProviderOptions {
	providerOptions := &ske.ProviderOptions{
		KubernetesVersions: &[]ske.KubernetesVersion{
			{
				State:   utils.Ptr("supported"),
				Version: utils.Ptr("1.2.3"),
			},
			{
				State:   utils.Ptr("supported"),
				Version: utils.Ptr("3.2.1"),
			},
			{
				State:   utils.Ptr("not-supported"),
				Version: utils.Ptr("4.4.4"),
			},
		},
		MachineImages: &[]ske.MachineImage{
			{
				Name: utils.Ptr("flatcar"),
				Versions: &[]ske.MachineImageVersion{
					{
						State:   utils.Ptr("supported"),
						Version: utils.Ptr("1.2.3"),
						Cri: &[]ske.CRI{
							{
								Name: utils.Ptr("not-containerd"),
							},
							{
								Name: utils.Ptr("containerd"),
							},
						},
					},
					{
						State:   utils.Ptr("supported"),
						Version: utils.Ptr("3.2.1"),
						Cri: &[]ske.CRI{
							{
								Name: utils.Ptr("not-containerd"),
							},
							{
								Name: utils.Ptr("containerd"),
							},
						},
					},
				},
			},
			{
				Name: utils.Ptr("not-flatcar"),
				Versions: &[]ske.MachineImageVersion{
					{
						State:   utils.Ptr("supported"),
						Version: utils.Ptr("4.4.4"),
						Cri: &[]ske.CRI{
							{
								Name: utils.Ptr("containerd"),
							},
						},
					},
				},
			},
			{
				Name: utils.Ptr("flatcar"),
				Versions: &[]ske.MachineImageVersion{
					{
						State:   utils.Ptr("supported"),
						Version: utils.Ptr("4.4.4"),
					},
				},
			},
			{
				Name: utils.Ptr("flatcar"),
				Versions: &[]ske.MachineImageVersion{
					{
						State:   utils.Ptr("not-supported"),
						Version: utils.Ptr("4.4.4"),
						Cri: &[]ske.CRI{
							{
								Name: utils.Ptr("containerd"),
							},
						},
					},
				},
			},
			{
				Name: utils.Ptr("flatcar"),
				Versions: &[]ske.MachineImageVersion{
					{
						State:   utils.Ptr("supported"),
						Version: utils.Ptr("4.4.4"),
						Cri: &[]ske.CRI{
							{
								Name: utils.Ptr("not-containerd"),
							},
						},
					},
				},
			},
		},
	}
	for _, mod := range mods {
		mod(providerOptions)
	}
	return providerOptions
}

func fixtureGetDefaultPayload(mods ...func(*ske.CreateOrUpdateClusterPayload)) *ske.CreateOrUpdateClusterPayload {
	payload := &ske.CreateOrUpdateClusterPayload{
		Extensions: &ske.Extension{
			Acl: &ske.ACL{
				AllowedCidrs: &[]string{},
				Enabled:      utils.Ptr(false),
			},
		},
		Kubernetes: &ske.Kubernetes{
			Version: utils.Ptr("3.2.1"),
		},
		Nodepools: &[]ske.Nodepool{
			{
				AvailabilityZones: &[]string{
					"eu01-3",
				},
				Cri: &ske.CRI{
					Name: utils.Ptr("containerd"),
				},
				Machine: &ske.Machine{
					Type: utils.Ptr("b1.2"),
					Image: &ske.Image{
						Version: utils.Ptr("3.2.1"),
						Name:    utils.Ptr("flatcar"),
					},
				},
				MaxSurge: utils.Ptr(int64(1)),
				Maximum:  utils.Ptr(int64(2)),
				Minimum:  utils.Ptr(int64(1)),
				Name:     utils.Ptr("pool-default"),
				Volume: &ske.Volume{
					Type: utils.Ptr("storage_premium_perf2"),
					Size: utils.Ptr(int64(50)),
				},
			},
		},
	}
	for _, mod := range mods {
		mod(payload)
	}
	return payload
}

func TestGetDefaultPayload(t *testing.T) {
	tests := []struct {
		description              string
		listProviderOptionsFails bool
		listProviderOptionsResp  *ske.ProviderOptions
		isValid                  bool
		expectedOutput           *ske.CreateOrUpdateClusterPayload
	}{
		{
			description:             "base",
			listProviderOptionsResp: fixtureProviderOptions(),
			isValid:                 true,
			expectedOutput:          fixtureGetDefaultPayload(),
		},
		{
			description:              "get provider options fails",
			listProviderOptionsFails: true,
			isValid:                  false,
		},
		{
			description: "no Kubernetes versions 1",
			listProviderOptionsResp: fixtureProviderOptions(func(po *ske.ProviderOptions) {
				po.KubernetesVersions = nil
			}),
			isValid: false,
		},
		{
			description: "no Kubernetes versions 2",
			listProviderOptionsResp: fixtureProviderOptions(func(po *ske.ProviderOptions) {
				po.KubernetesVersions = &[]ske.KubernetesVersion{}
			}),
			isValid: false,
		},
		{
			description: "no supported Kubernetes versions",
			listProviderOptionsResp: fixtureProviderOptions(func(po *ske.ProviderOptions) {
				po.KubernetesVersions = &[]ske.KubernetesVersion{
					{
						State:   utils.Ptr("not-supported"),
						Version: utils.Ptr("1.2.3"),
					},
				}
			}),
			isValid: false,
		},
		{
			description: "no machine images 1",
			listProviderOptionsResp: fixtureProviderOptions(func(po *ske.ProviderOptions) {
				po.MachineImages = &[]ske.MachineImage{}
			}),
			isValid: false,
		},
		{
			description: "no machine images 2",
			listProviderOptionsResp: fixtureProviderOptions(func(po *ske.ProviderOptions) {
				po.MachineImages = nil
			}),
			isValid: false,
		},
		{
			description: "no machine image versions 1",
			listProviderOptionsResp: fixtureProviderOptions(func(po *ske.ProviderOptions) {
				po.MachineImages = &[]ske.MachineImage{
					{
						Name:     utils.Ptr("image-1"),
						Versions: nil,
					},
				}
			}),
			isValid: false,
		},
		{
			description: "no machine image versions 2",
			listProviderOptionsResp: fixtureProviderOptions(func(po *ske.ProviderOptions) {
				po.MachineImages = &[]ske.MachineImage{
					{
						Name:     utils.Ptr("image-1"),
						Versions: &[]ske.MachineImageVersion{},
					},
				}
			}),
			isValid: false,
		},
		{
			description: "no supported machine image versions",
			listProviderOptionsResp: fixtureProviderOptions(func(po *ske.ProviderOptions) {
				po.MachineImages = &[]ske.MachineImage{
					{
						Name: utils.Ptr("image-1"),
						Versions: &[]ske.MachineImageVersion{
							{
								State:   utils.Ptr("not-supported"),
								Version: utils.Ptr("1.2.3"),
							},
						},
					},
				}
			}),
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			client := &skeClientMocked{
				listProviderOptionsFails: tt.listProviderOptionsFails,
				listProviderOptionsResp:  tt.listProviderOptionsResp,
			}

			output, err := GetDefaultPayload(context.Background(), client)

			if tt.isValid && err != nil {
				t.Errorf("failed on valid input")
			}
			if !tt.isValid && err == nil {
				t.Errorf("did not fail on invalid input")
			}
			if !tt.isValid {
				return
			}
			diff := cmp.Diff(output, tt.expectedOutput)
			if diff != "" {
				t.Fatalf("Output is not as expected: %s", diff)
			}
		})
	}
}

func TestConvertToSeconds(t *testing.T) {
	tests := []struct {
		description    string
		expirationTime string
		isValid        bool
		expectedOutput string
	}{
		{
			description:    "seconds",
			expirationTime: "30s",
			isValid:        true,
			expectedOutput: "30",
		},
		{
			description:    "minutes",
			expirationTime: "30m",
			isValid:        true,
			expectedOutput: "1800",
		},
		{
			description:    "hours",
			expirationTime: "30h",
			isValid:        true,
			expectedOutput: "108000",
		},
		{
			description:    "days",
			expirationTime: "30d",
			isValid:        true,
			expectedOutput: "2592000",
		},
		{
			description:    "months",
			expirationTime: "30M",
			isValid:        true,
			expectedOutput: "77760000",
		},
		{
			description:    "leading zero",
			expirationTime: "0030M",
			isValid:        true,
			expectedOutput: "77760000",
		},
		{
			description:    "invalid unit",
			expirationTime: "30x",
			isValid:        false,
		},
		{
			description:    "invalid unit 2",
			expirationTime: "3000abcdef",
			isValid:        false,
		},
		{
			description:    "invalid unit 3",
			expirationTime: "3000abcdef000",
			isValid:        false,
		},
		{
			description:    "invalid time",
			expirationTime: "x",
			isValid:        false,
		},
		{
			description:    "empty",
			expirationTime: "",
			isValid:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			output, err := ConvertToSeconds(tt.expirationTime)

			if tt.isValid && err != nil {
				t.Errorf("failed on valid input")
			}
			if !tt.isValid && err == nil {
				t.Errorf("did not fail on invalid input")
			}
			if !tt.isValid {
				return
			}
			if *output != tt.expectedOutput {
				t.Errorf("expected output to be %s, got %s", tt.expectedOutput, *output)
			}
		})
	}
}

func TestWriteConfigFile(t *testing.T) {
	tests := []struct {
		description string
		location    string
		kubeconfig  string
		isValid     bool
		expectedErr string
	}{
		{
			description: "base",
			location:    "test_data/base/config",
			kubeconfig:  "kubeconfig",
			isValid:     true,
		},
		{
			description: "empty location",
			location:    "",
			kubeconfig:  "kubeconfig",
			isValid:     false,
		},
		{
			description: "no permission location",
			location:    "/root/config",
			kubeconfig:  "kubeconfig",
			isValid:     false,
		},
		{
			description: "path is only dir",
			location:    "test_data/only_dir/",
			kubeconfig:  "kubeconfig",
			isValid:     false,
		},
		{
			description: "empty kubeconfig",
			location:    "test_data/empty/config",
			kubeconfig:  "",
			isValid:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			err := WriteConfigFile(tt.location, tt.kubeconfig)

			if tt.isValid && err != nil {
				t.Errorf("failed on valid input")
			}
			if !tt.isValid && err == nil {
				t.Errorf("did not fail on invalid input")
			}

			if tt.isValid {
				data, err := os.ReadFile(tt.location)
				if err != nil {
					t.Errorf("could not read file: %s", tt.location)
				}
				if string(data) != tt.kubeconfig {
					t.Errorf("expected file content to be %s, got %s", tt.kubeconfig, string(data))
				}
			}
		})
	}
	// Cleanup
	err := os.RemoveAll("test_data/")
	if err != nil {
		t.Errorf("failed cleaning test data")
	}
}
