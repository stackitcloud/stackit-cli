package utils

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/resourcemanager"
)

var (
	testOrgId = uuid.NewString()
)

const (
	testOrgName = "organization"
)

type resourceManagerClientMocked struct {
	getOrganizationFails bool
	getOrganizationResp  *resourcemanager.OrganizationResponse
	getProjectFails      bool
	getProjectResp       *resourcemanager.GetProjectResponse
}

func (s *resourceManagerClientMocked) GetOrganizationExecute(_ context.Context, _ string) (*resourcemanager.OrganizationResponse, error) {
	if s.getOrganizationFails {
		return nil, fmt.Errorf("could not get organization")
	}
	return s.getOrganizationResp, nil
}

func (s *resourceManagerClientMocked) GetProjectExecute(_ context.Context, _ string) (*resourcemanager.GetProjectResponse, error) {
	if s.getProjectFails {
		return nil, fmt.Errorf("could not get project")
	}
	return s.getProjectResp, nil
}

func TestGetOrganizationName(t *testing.T) {
	tests := []struct {
		description          string
		getOrganizationFails bool
		getOrganizationResp  *resourcemanager.OrganizationResponse
		isValid              bool
		expectedOutput       string
	}{
		{
			description: "base",
			getOrganizationResp: &resourcemanager.OrganizationResponse{
				Name: utils.Ptr(testOrgName),
			},
			isValid:        true,
			expectedOutput: testOrgName,
		},
		{
			description:          "get organization fails",
			getOrganizationFails: true,
			isValid:              false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			client := &resourceManagerClientMocked{
				getOrganizationFails: tt.getOrganizationFails,
				getOrganizationResp:  tt.getOrganizationResp,
			}

			output, err := GetOrganizationName(context.Background(), client, testOrgId)

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

func TestGetProjectName(t *testing.T) {
	tests := []struct {
		description     string
		getProjectFails bool
		getProjectResp  *resourcemanager.GetProjectResponse
		isValid         bool
		expectedOutput  string
	}{
		{
			description: "base",
			getProjectResp: &resourcemanager.GetProjectResponse{
				Name: utils.Ptr("project"),
			},
			isValid:        true,
			expectedOutput: "project",
		},
		{
			description:     "get project fails",
			getProjectFails: true,
			isValid:         false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			client := &resourceManagerClientMocked{
				getProjectFails: tt.getProjectFails,
				getProjectResp:  tt.getProjectResp,
			}

			output, err := GetProjectName(context.Background(), client, testOrgId)

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
