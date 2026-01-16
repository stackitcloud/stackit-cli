package utils

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/sfs"
)

const (
	testShareName        = "share-name"
	testResourcePoolName = "resource-pool-name"
	testExportPolicyName = "export-policy-name"
	testSnapshotName     = "snapshot-name"
	testRegion           = "eu01"
)

var (
	testPolicyId  = uuid.NewString()
	testProjectId = uuid.NewString()
)

type sfsClientMocked struct {
	getShareFails        bool
	getShareResp         *sfs.GetShareResponse
	getResourcePoolFails bool
	getResourcePoolResp  *sfs.GetResourcePoolResponse
	getExportPolicyFails bool
	getExportPolicyResp  *sfs.GetShareExportPolicyResponse
}

func (s *sfsClientMocked) GetShareExecute(_ context.Context, _, _, _, _ string) (*sfs.GetShareResponse, error) {
	if s.getShareFails {
		return nil, fmt.Errorf("could not get share")
	}
	return s.getShareResp, nil
}

func (s *sfsClientMocked) GetResourcePoolExecute(_ context.Context, _, _, _ string) (*sfs.GetResourcePoolResponse, error) {
	if s.getResourcePoolFails {
		return nil, fmt.Errorf("could not get resource pool")
	}
	return s.getResourcePoolResp, nil
}

func (s *sfsClientMocked) GetShareExportPolicyExecute(_ context.Context, _, _, _ string) (*sfs.GetShareExportPolicyResponse, error) {
	if s.getExportPolicyFails {
		return nil, fmt.Errorf("could not get export policy")
	}
	return s.getExportPolicyResp, nil
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
				ShareExportPolicy: &sfs.GetShareExportPolicyResponseShareExportPolicy{
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
			client := &sfsClientMocked{
				getExportPolicyFails: tt.getExportPolicyFails,
				getExportPolicyResp:  tt.getExportPolicyResp,
			}

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
				Share: &sfs.GetShareResponseShare{
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
			client := &sfsClientMocked{
				getShareFails: tt.getShareFails,
				getShareResp:  tt.getShareResp,
			}

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
				ResourcePool: &sfs.GetResourcePoolResponseResourcePool{
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
			client := &sfsClientMocked{
				getResourcePoolResp:  tt.getResourcePoolResp,
				getResourcePoolFails: tt.getResourcePoolFails,
			}

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
