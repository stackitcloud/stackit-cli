package utils

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	sfs "github.com/stackitcloud/stackit-sdk-go/services/sfs/v1api"

	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

const (
	testShareName        = "share-name"
	testResourcePoolName = "resource-pool-name"
	testExportPolicyName = "export-policy-name"
	testRegion           = "eu01"
)

var (
	testPolicyId  = uuid.NewString()
	testProjectId = uuid.NewString()
)

type mockSettings struct {
	getShareFails bool
	getShareResp  *sfs.GetShareResponse

	getResourcePoolFails bool
	getResourcePoolResp  *sfs.GetResourcePoolResponse

	getExportPolicyFails bool
	getExportPolicyResp  *sfs.GetShareExportPolicyResponse
}

func newAPIMock(settings *mockSettings) sfs.DefaultAPI {
	return &sfs.DefaultAPIServiceMock{
		GetShareExecuteMock: utils.Ptr(func(_ sfs.ApiGetShareRequest) (*sfs.GetShareResponse, error) {
			if settings.getShareFails {
				return nil, fmt.Errorf("could not get share details")
			}

			return settings.getShareResp, nil
		}),
		GetShareExportPolicyExecuteMock: utils.Ptr(func(_ sfs.ApiGetShareExportPolicyRequest) (*sfs.GetShareExportPolicyResponse, error) {
			if settings.getExportPolicyFails {
				return nil, fmt.Errorf("could not get export policy details")
			}

			return settings.getExportPolicyResp, nil
		}),
		GetResourcePoolExecuteMock: utils.Ptr(func(_ sfs.ApiGetResourcePoolRequest) (*sfs.GetResourcePoolResponse, error) {
			if settings.getResourcePoolFails {
				return nil, fmt.Errorf("could not get resource pool details")
			}

			return settings.getResourcePoolResp, nil
		}),
	}
}

func TestGetExportPolicyName(t *testing.T) {
	tests := []struct {
		description          string
		getExportPolicyResp  *sfs.GetShareExportPolicyResponse
		getExportPolicyFails bool
		isValid              bool
		expectedOutput       string
	}{
		{
			description: "base",
			getExportPolicyResp: &sfs.GetShareExportPolicyResponse{
				ShareExportPolicy: &sfs.ShareExportPolicy{
					Name: utils.Ptr(testExportPolicyName),
				},
			},
			isValid:        true,
			expectedOutput: testExportPolicyName,
		},
		{
			description:          "get export policy fails",
			getExportPolicyFails: true,
			isValid:              false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			client := newAPIMock(&mockSettings{
				getExportPolicyFails: tt.getExportPolicyFails,
				getExportPolicyResp:  tt.getExportPolicyResp,
			})

			output, err := GetExportPolicyName(context.Background(), client, testProjectId, testRegion, testPolicyId)

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
				t.Errorf("expected output to be %s, got %s", tt.expectedOutput, output)
			}
		})
	}
}

func TestGetShareName(t *testing.T) {
	tests := []struct {
		description    string
		getShareResp   *sfs.GetShareResponse
		getShareFails  bool
		isValid        bool
		expectedOutput string
	}{
		{
			description: "base",
			getShareResp: &sfs.GetShareResponse{
				Share: &sfs.Share{
					Name: utils.Ptr(testShareName),
				},
			},
			isValid:        true,
			expectedOutput: testShareName,
		},
		{
			description:    "get share fails",
			getShareFails:  true,
			isValid:        false,
			expectedOutput: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			client := newAPIMock(&mockSettings{
				getShareFails: tt.getShareFails,
				getShareResp:  tt.getShareResp,
			})

			output, err := GetShareName(context.Background(), client, testProjectId, testRegion, "", "")

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
				t.Errorf("expected output to be %s, got %s", tt.expectedOutput, output)
			}
		})
	}
}

func TestGetResourcePoolName(t *testing.T) {
	tests := []struct {
		description          string
		getResourcePoolResp  *sfs.GetResourcePoolResponse
		getResourcePoolFails bool
		isValid              bool
		expectedOutput       string
	}{
		{
			description: "base",
			getResourcePoolResp: &sfs.GetResourcePoolResponse{
				ResourcePool: &sfs.ResourcePool{
					Name: utils.Ptr(testResourcePoolName),
				},
			},
			isValid:        true,
			expectedOutput: testResourcePoolName,
		},
		{
			description:          "get resource pool fails",
			getResourcePoolFails: true,
			isValid:              false,
			expectedOutput:       "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			client := newAPIMock(&mockSettings{
				getResourcePoolResp:  tt.getResourcePoolResp,
				getResourcePoolFails: tt.getResourcePoolFails,
			})

			output, err := GetResourcePoolName(context.Background(), client, testProjectId, testRegion, "")

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
				t.Errorf("expected output to be %s, got %s", tt.expectedOutput, output)
			}
		})
	}
}
