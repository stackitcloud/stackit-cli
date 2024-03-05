package utils

import (
	"context"
	"fmt"
	"testing"

	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/google/uuid"
	"github.com/stackitcloud/stackit-sdk-go/services/objectstorage"
)

var (
	testProjectId          = uuid.NewString()
	testCredentialsGroupId = uuid.NewString()
)

const (
	testCredentialsGroupName = "testGroup"
)

type objectStorageClientMocked struct {
	listCredentialsGroupsFails bool
	listCredentialsGroupsResp  *objectstorage.ListCredentialsGroupsResponse
	listAccessKeysReq          objectstorage.ApiListAccessKeysRequest
}

func (m *objectStorageClientMocked) ListCredentialsGroupsExecute(_ context.Context, _ string) (*objectstorage.ListCredentialsGroupsResponse, error) {
	if m.listCredentialsGroupsFails {
		return nil, fmt.Errorf("could not list credentials groups")
	}
	return m.listCredentialsGroupsResp, nil
}

func (m *objectStorageClientMocked) ListAccessKeys(_ context.Context, _ string) objectstorage.ApiListAccessKeysRequest {
	return m.listAccessKeysReq
}

func TestGetCredentialsGroupName(t *testing.T) {
	tests := []struct {
		description                string
		listCredentialsGroupsFails bool
		listCredentialsGroupsResp  *objectstorage.ListCredentialsGroupsResponse
		isValid                    bool
		expectedOutput             string
	}{
		{
			description: "base",
			listCredentialsGroupsResp: &objectstorage.ListCredentialsGroupsResponse{
				CredentialsGroups: &[]objectstorage.CredentialsGroup{
					{
						CredentialsGroupId: utils.Ptr(testCredentialsGroupId),
						DisplayName:        utils.Ptr(testCredentialsGroupName),
					},
				},
			},
			isValid:        true,
			expectedOutput: testCredentialsGroupName,
		},
		{
			description:                "list credentials groups fails",
			listCredentialsGroupsFails: true,
			isValid:                    false,
		},
		{
			description: "multiple credentials groups",
			listCredentialsGroupsResp: &objectstorage.ListCredentialsGroupsResponse{
				CredentialsGroups: &[]objectstorage.CredentialsGroup{
					{
						CredentialsGroupId: utils.Ptr("test-UUID"),
						DisplayName:        utils.Ptr("test-name"),
					},
					{
						CredentialsGroupId: utils.Ptr(testCredentialsGroupId),
						DisplayName:        utils.Ptr(testCredentialsGroupName),
					},
				},
			},
			isValid:        true,
			expectedOutput: testCredentialsGroupName,
		},
		{
			description: "nil credentials groups",
			listCredentialsGroupsResp: &objectstorage.ListCredentialsGroupsResponse{
				CredentialsGroups: nil,
			},
			isValid: false,
		},
		{
			description: "nil credentials group id",
			listCredentialsGroupsResp: &objectstorage.ListCredentialsGroupsResponse{
				CredentialsGroups: &[]objectstorage.CredentialsGroup{
					{
						CredentialsGroupId: nil,
					},
				},
			},
			isValid: false,
		},
		{
			description: "nil credentials group name",
			listCredentialsGroupsResp: &objectstorage.ListCredentialsGroupsResponse{
				CredentialsGroups: &[]objectstorage.CredentialsGroup{
					{
						CredentialsGroupId: utils.Ptr(testCredentialsGroupId),
						DisplayName:        nil,
					},
				},
			},
			isValid: false,
		},
		{
			description: "empty credentials group name",
			listCredentialsGroupsResp: &objectstorage.ListCredentialsGroupsResponse{
				CredentialsGroups: &[]objectstorage.CredentialsGroup{
					{
						CredentialsGroupId: utils.Ptr(testCredentialsGroupId),
						DisplayName:        utils.Ptr(""),
					},
				},
			},
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			client := &objectStorageClientMocked{
				listCredentialsGroupsFails: tt.listCredentialsGroupsFails,
				listCredentialsGroupsResp:  tt.listCredentialsGroupsResp,
			}

			output, err := GetCredentialsGroupName(context.Background(), client, testProjectId, testCredentialsGroupId)

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
