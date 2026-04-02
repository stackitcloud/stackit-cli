package utils

import (
	"context"
	"fmt"

	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stackitcloud/stackit-sdk-go/core/oapierror"
	objectstorage "github.com/stackitcloud/stackit-sdk-go/services/objectstorage/v2api"

	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

var (
	testProjectId          = uuid.NewString()
	testCredentialsGroupId = uuid.NewString()
)

const (
	testCredentialsGroupName = "testGroup"
	testCredentialsName      = "testCredential"
	testCredentialsId        = "credentialsID" //nolint:gosec // linter false positive
	testRegion               = "eu01"
)

type mockSettings struct {
	serviceDisabled       bool
	getServiceStatusFails bool

	listCredentialsGroupsFails bool
	listCredentialsGroupsResp  *objectstorage.ListCredentialsGroupsResponse

	listAccessKeysFails bool
	listAccessKeysResp  *objectstorage.ListAccessKeysResponse
}

func newAPIMock(settings *mockSettings) objectstorage.DefaultAPI {
	return &objectstorage.DefaultAPIServiceMock{
		GetServiceStatusExecuteMock: utils.Ptr(func(_ objectstorage.ApiGetServiceStatusRequest) (*objectstorage.ProjectStatus, error) {
			if settings.getServiceStatusFails {
				return nil, fmt.Errorf("could not get service status")
			}

			if settings.serviceDisabled {
				return nil, &oapierror.GenericOpenAPIError{StatusCode: http.StatusNotFound}
			}

			return &objectstorage.ProjectStatus{}, nil
		}),
		ListCredentialsGroupsExecuteMock: utils.Ptr(func(_ objectstorage.ApiListCredentialsGroupsRequest) (*objectstorage.ListCredentialsGroupsResponse, error) {
			if settings.listCredentialsGroupsFails {
				return nil, fmt.Errorf("could not list credentials groups")
			}

			return settings.listCredentialsGroupsResp, nil
		}),
		ListAccessKeysExecuteMock: utils.Ptr(func(_ objectstorage.ApiListAccessKeysRequest) (*objectstorage.ListAccessKeysResponse, error) {
			if settings.listAccessKeysFails {
				return nil, &oapierror.GenericOpenAPIError{StatusCode: http.StatusBadGateway}
			}

			return settings.listAccessKeysResp, nil
		}),
	}
}

func TestProjectEnabled(t *testing.T) {
	tests := []struct {
		description     string
		serviceDisabled bool
		getProjectFails bool
		isValid         bool
		expectedOutput  bool
	}{
		{
			description:    "project enabled",
			isValid:        true,
			expectedOutput: true,
		},
		{
			description:     "project disabled (404)",
			serviceDisabled: true,
			isValid:         true,
			expectedOutput:  false,
		},
		{
			description:     "get project fails",
			getProjectFails: true,
			isValid:         false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			client := newAPIMock(&mockSettings{
				serviceDisabled:       tt.serviceDisabled,
				getServiceStatusFails: tt.getProjectFails,
			})

			output, err := ProjectEnabled(context.Background(), client, testProjectId, testRegion)

			if tt.isValid && err != nil {
				fmt.Printf("failed on valid input: %v", err)
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
				CredentialsGroups: []objectstorage.CredentialsGroup{
					{
						CredentialsGroupId: testCredentialsGroupId,
						DisplayName:        testCredentialsGroupName,
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
				CredentialsGroups: []objectstorage.CredentialsGroup{
					{
						CredentialsGroupId: "test-UUID",
						DisplayName:        "test-name",
					},
					{
						CredentialsGroupId: testCredentialsGroupId,
						DisplayName:        testCredentialsGroupName,
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
			description: "empty credentials group id",
			listCredentialsGroupsResp: &objectstorage.ListCredentialsGroupsResponse{
				CredentialsGroups: []objectstorage.CredentialsGroup{
					{
						CredentialsGroupId: "",
					},
				},
			},
			isValid: false,
		},
		{
			description: "empty credentials group name",
			listCredentialsGroupsResp: &objectstorage.ListCredentialsGroupsResponse{
				CredentialsGroups: []objectstorage.CredentialsGroup{
					{
						CredentialsGroupId: testCredentialsGroupId,
						DisplayName:        "",
					},
				},
			},
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			client := newAPIMock(&mockSettings{
				listCredentialsGroupsFails: tt.listCredentialsGroupsFails,
				listCredentialsGroupsResp:  tt.listCredentialsGroupsResp,
			})

			output, err := GetCredentialsGroupName(context.Background(), client, testProjectId, testCredentialsGroupId, testRegion)

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

func TestGetCredentialsName(t *testing.T) {
	tests := []struct {
		description             string
		listAccessKeysResp      *objectstorage.ListAccessKeysResponse
		expectedOutput          string
		getCredentialsNameFails bool
		isValid                 bool
	}{
		{
			description: "base",
			listAccessKeysResp: &objectstorage.ListAccessKeysResponse{
				AccessKeys: []objectstorage.AccessKey{
					{
						KeyId:       testCredentialsId,
						DisplayName: testCredentialsName,
					},
				},
			},
			isValid:        true,
			expectedOutput: testCredentialsName,
		},
		{
			description:             "get credentials name fails",
			getCredentialsNameFails: true,
			isValid:                 false,
		},
		{
			description: "multiple credentials",
			listAccessKeysResp: &objectstorage.ListAccessKeysResponse{
				AccessKeys: []objectstorage.AccessKey{
					{
						KeyId:       "test-UUID",
						DisplayName: "test-name",
					},
					{
						KeyId:       testCredentialsId,
						DisplayName: testCredentialsName,
					},
				},
			},
			isValid:        true,
			expectedOutput: testCredentialsName,
		},
		{
			description: "nil credentials",
			listAccessKeysResp: &objectstorage.ListAccessKeysResponse{
				AccessKeys: nil,
			},
			isValid: false,
		},
		{
			description: "empty credentials id",
			listAccessKeysResp: &objectstorage.ListAccessKeysResponse{
				AccessKeys: []objectstorage.AccessKey{
					{
						KeyId: "",
					},
				},
			},
			isValid: false,
		},
		{
			description: "empty credentials name",
			listAccessKeysResp: &objectstorage.ListAccessKeysResponse{
				AccessKeys: []objectstorage.AccessKey{
					{
						KeyId:       testCredentialsId,
						DisplayName: "",
					},
				},
			},
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			client := newAPIMock(&mockSettings{
				listAccessKeysFails: tt.getCredentialsNameFails,
				listAccessKeysResp:  tt.listAccessKeysResp,
			})

			output, err := GetCredentialsName(context.Background(), client, testProjectId, testCredentialsGroupId, testCredentialsId, testRegion)

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
