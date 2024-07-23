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
}

func (s *resourceManagerClientMocked) GetOrganizationExecute(_ context.Context, _ string) (*resourcemanager.OrganizationResponse, error) {
	if s.getOrganizationFails {
		return nil, fmt.Errorf("could not get organization")
	}
	return s.getOrganizationResp, nil
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
