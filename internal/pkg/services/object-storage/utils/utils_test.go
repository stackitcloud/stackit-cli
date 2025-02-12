package utils

import (
	"context"
	"encoding/json"
	"fmt"

	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/core/config"
	"github.com/stackitcloud/stackit-sdk-go/core/oapierror"
	"github.com/stackitcloud/stackit-sdk-go/services/objectstorage"
)

var (
	testProjectId          = uuid.NewString()
	testCredentialsGroupId = uuid.NewString()
	testCredentialsId      = "credentialsID" //nolint:gosec // linter false positive
	testRegion             = "eu01"
)

const (
	testCredentialsGroupName = "testGroup"
	testCredentialsName      = "testCredential"
)

type objectStorageClientMocked struct {
	serviceDisabled            bool
	getServiceStatusFails      bool
	listCredentialsGroupsFails bool
	listCredentialsGroupsResp  *objectstorage.ListCredentialsGroupsResponse
	listAccessKeysReq          objectstorage.ApiListAccessKeysRequest
}

func (m *objectStorageClientMocked) GetServiceStatusExecute(_ context.Context, _, _ string) (*objectstorage.ProjectStatus, error) {
	if m.getServiceStatusFails {
		return nil, fmt.Errorf("could not get service status")
	}
	if m.serviceDisabled {
		return nil, &oapierror.GenericOpenAPIError{StatusCode: 404}
	}
	return &objectstorage.ProjectStatus{}, nil
}

func (m *objectStorageClientMocked) ListCredentialsGroupsExecute(_ context.Context, _, _ string) (*objectstorage.ListCredentialsGroupsResponse, error) {
	if m.listCredentialsGroupsFails {
		return nil, fmt.Errorf("could not list credentials groups")
	}
	return m.listCredentialsGroupsResp, nil
}

func (m *objectStorageClientMocked) ListAccessKeys(_ context.Context, _, _ string) objectstorage.ApiListAccessKeysRequest {
	return m.listAccessKeysReq
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
			client := &objectStorageClientMocked{
				serviceDisabled:       tt.serviceDisabled,
				getServiceStatusFails: tt.getProjectFails,
			}

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
				AccessKeys: &[]objectstorage.AccessKey{
					{
						KeyId:       utils.Ptr(testCredentialsId),
						DisplayName: utils.Ptr(testCredentialsName),
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
				AccessKeys: &[]objectstorage.AccessKey{
					{
						KeyId:       utils.Ptr("test-UUID"),
						DisplayName: utils.Ptr("test-name"),
					},
					{
						KeyId:       utils.Ptr(testCredentialsId),
						DisplayName: utils.Ptr(testCredentialsName),
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
			description: "nil credentials id",
			listAccessKeysResp: &objectstorage.ListAccessKeysResponse{
				AccessKeys: &[]objectstorage.AccessKey{
					{
						KeyId: nil,
					},
				},
			},
			isValid: false,
		},
		{
			description: "nil credentials name",
			listAccessKeysResp: &objectstorage.ListAccessKeysResponse{
				AccessKeys: &[]objectstorage.AccessKey{
					{
						KeyId:       utils.Ptr(testCredentialsId),
						DisplayName: nil,
					},
				},
			},
			isValid: false,
		},
		{
			description: "empty credentials name",
			listAccessKeysResp: &objectstorage.ListAccessKeysResponse{
				AccessKeys: &[]objectstorage.AccessKey{
					{
						KeyId:       utils.Ptr(testCredentialsId),
						DisplayName: utils.Ptr(""),
					},
				},
			},
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			mockedRespBytes, err := json.Marshal(tt.listAccessKeysResp)
			if err != nil {
				t.Fatalf("Failed to marshal mocked response: %v", err)
			}

			handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				if tt.getCredentialsNameFails {
					w.WriteHeader(http.StatusBadGateway)
					w.Header().Set("Content-Type", "application/json")
					_, err := w.Write([]byte("{\"message\": \"Something bad happened\""))
					if err != nil {
						t.Errorf("Failed to write bad response: %v", err)
					}
					return
				}

				_, err := w.Write(mockedRespBytes)
				if err != nil {
					t.Errorf("Failed to write response: %v", err)
				}
			})
			mockedServer := httptest.NewServer(handler)
			defer mockedServer.Close()
			client, err := objectstorage.NewAPIClient(
				config.WithEndpoint(mockedServer.URL),
				config.WithoutAuthentication(),
			)
			if err != nil {
				t.Fatalf("Failed to initialize client: %v", err)
			}

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
