package utils

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	scf "github.com/stackitcloud/stackit-sdk-go/services/scf/v1api"

	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

var (
	testProjectId      = uuid.NewString()
	testOrganizationId = uuid.NewString()
)

const (
	testOrganizationName = "organization"
	testRegion           = "eu01"
)

type mockSettings struct {
	getOrganizationFails bool
	getOrganizationResp  *scf.Organization
}

func newApiMock(s *mockSettings) scf.DefaultAPI {
	return &scf.DefaultAPIServiceMock{
		GetOrganizationExecuteMock: utils.Ptr(func(_ scf.ApiGetOrganizationRequest) (*scf.Organization, error) {
			if s.getOrganizationFails {
				return nil, fmt.Errorf("could not get organization")
			}
			return s.getOrganizationResp, nil
		}),
	}
}

func TestGetOrganizationName(t *testing.T) {
	tests := []struct {
		description          string
		getOrganizationFails bool
		getOrganizationResp  *scf.Organization
		isValid              bool
		expectedOutput       string
	}{
		{
			description: "base",
			getOrganizationResp: &scf.Organization{
				Name: testOrganizationName,
			},
			isValid:        true,
			expectedOutput: testOrganizationName,
		},
		{
			description:          "get organization fails",
			getOrganizationFails: true,
			isValid:              false,
		},
		{
			description:         "nil response",
			getOrganizationResp: nil,
			isValid:             false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			settings := &mockSettings{
				getOrganizationFails: tt.getOrganizationFails,
				getOrganizationResp:  tt.getOrganizationResp,
			}

			output, err := GetOrganizationName(context.Background(), newApiMock(settings), testProjectId, testRegion, testOrganizationId)

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
