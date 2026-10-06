package utils

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	resourcemanager "github.com/stackitcloud/stackit-sdk-go/services/resourcemanager/v0api"

	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

var (
	testOrgId     = uuid.NewString()
	testFolderId  = uuid.NewString()
	testProjectId = uuid.NewString()
)

const (
	testOrgName     = "organization"
	testFolderName  = "folder"
	testProjectName = "project"
)

type resourceManagerClientMocked struct {
	getOrganizationFails bool
	getOrganizationResp  *resourcemanager.OrganizationResponse
	getProjectFails      bool
	getProjectResp       *resourcemanager.GetProjectResponse
	getFolderFails       bool
	getFolderResp        *resourcemanager.GetFolderDetailsResponse
}

func (s *resourceManagerClientMocked) newMock() resourcemanager.DefaultAPI {
	return resourcemanager.DefaultAPIServiceMock{
		GetOrganizationExecuteMock: utils.Ptr(func(_ resourcemanager.ApiGetOrganizationRequest) (*resourcemanager.OrganizationResponse, error) {
			if s.getOrganizationFails {
				return nil, fmt.Errorf("could not get organization")
			}
			return s.getOrganizationResp, nil
		}),
		GetProjectExecuteMock: utils.Ptr(func(_ resourcemanager.ApiGetProjectRequest) (*resourcemanager.GetProjectResponse, error) {
			if s.getProjectFails {
				return nil, fmt.Errorf("could not get project")
			}
			return s.getProjectResp, nil
		}),
		GetFolderDetailsExecuteMock: utils.Ptr(func(_ resourcemanager.ApiGetFolderDetailsRequest) (*resourcemanager.GetFolderDetailsResponse, error) {
			if s.getFolderFails {
				return nil, fmt.Errorf("could not get folder")
			}
			return s.getFolderResp, nil
		}),
	}
}

func TestGetResourceLabel(t *testing.T) {
	tests := []struct {
		description              string
		isValid                  bool
		expectedOutput           string
		resourceType             string
		resourceId               string
		getProjectNameFails      bool
		getOrganizationNameFails bool
		getFolderNameFails       bool
		getOrganizationResp      *resourcemanager.OrganizationResponse
		getProjectResp           *resourcemanager.GetProjectResponse
		getFolderResp            *resourcemanager.GetFolderDetailsResponse
	}{
		{
			description:         "get project label",
			isValid:             true,
			resourceType:        "project",
			resourceId:          testProjectId,
			getProjectNameFails: false,
			expectedOutput:      testProjectName,
			getProjectResp: &resourcemanager.GetProjectResponse{
				Name: "project",
			},
		},
		{
			description:         "get project name fails",
			getProjectNameFails: true,
			isValid:             false,
		},
		{
			description:         "get folder label",
			isValid:             true,
			resourceType:        "folder",
			resourceId:          testFolderId,
			getProjectNameFails: false,
			expectedOutput:      testFolderName,
			getFolderResp: &resourcemanager.GetFolderDetailsResponse{
				Name: testFolderName,
			},
		},
		{
			description:        "get folder name fails",
			getFolderNameFails: true,
			isValid:            false,
		},
		{
			description:         "get organization label",
			isValid:             true,
			resourceType:        "organization",
			resourceId:          testOrgId,
			getProjectNameFails: false,
			expectedOutput:      testOrgName,
			getOrganizationResp: &resourcemanager.OrganizationResponse{
				Name: testOrgName,
			},
		},
		{
			description:              "get organization name fails",
			getOrganizationNameFails: true,
			isValid:                  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			client := &resourceManagerClientMocked{
				getFolderFails:       tt.getFolderNameFails,
				getOrganizationFails: tt.getOrganizationNameFails,
				getProjectFails:      tt.getProjectNameFails,
				getProjectResp:       tt.getProjectResp,
				getOrganizationResp:  tt.getOrganizationResp,
				getFolderResp:        tt.getFolderResp,
			}

			output, err := GetResourceLabel(context.Background(), client.newMock(), tt.resourceId, tt.resourceType)

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
