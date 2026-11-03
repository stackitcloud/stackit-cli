package utils

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	scf "github.com/stackitcloud/stackit-sdk-go/services/scf/v1api"

	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
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
	getSpaceFails        bool
	getSpaceResp         *scf.Space
}

func newApiMock(s *mockSettings) scf.DefaultAPI {
	return &scf.DefaultAPIServiceMock{
		GetSpaceExecuteMock: utils.Ptr(func(_ scf.ApiGetSpaceRequest) (*scf.Space, error) {
			if s.getSpaceFails {
				return nil, fmt.Errorf("could not get space")
			}
			return s.getSpaceResp, nil
		}),
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

func TestGetSpaceName(t *testing.T) {
	tests := []struct {
		description    string
		getSpaceFails  bool
		getSpaceResp   *scf.Space
		isValid        bool
		expectedOutput string
	}{
		{
			description:    "base",
			getSpaceResp:   &scf.Space{Name: "space"},
			isValid:        true,
			expectedOutput: "space",
		},
		{
			description:   "get space fails",
			getSpaceFails: true,
			isValid:       false,
		},
		{
			description:  "nil response",
			getSpaceResp: nil,
			isValid:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			settings := &mockSettings{
				getSpaceFails: tt.getSpaceFails,
				getSpaceResp:  tt.getSpaceResp,
			}

			output, err := GetSpaceName(context.Background(), newApiMock(settings), testProjectId, testRegion, testOrganizationId, uuid.NewString())

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

func TestGetOrganizationLabel(t *testing.T) {
	tests := []struct {
		description          string
		getOrganizationFails bool
		getOrganizationResp  *scf.Organization
		expectedOutput       string
	}{
		{
			description:         "base",
			getOrganizationResp: &scf.Organization{Name: testOrganizationName},
			expectedOutput:      testOrganizationName,
		},
		{
			description:          "get organization fails",
			getOrganizationFails: true,
			expectedOutput:       testOrganizationId,
		},
		{
			description:         "empty name",
			getOrganizationResp: &scf.Organization{},
			expectedOutput:      testOrganizationId,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			settings := &mockSettings{
				getOrganizationFails: tt.getOrganizationFails,
				getOrganizationResp:  tt.getOrganizationResp,
			}
			p := print.NewPrinter(nil, nil, nil)

			output := GetOrganizationLabel(context.Background(), p, newApiMock(settings), testProjectId, testRegion, testOrganizationId)
			if output != tt.expectedOutput {
				t.Errorf("expected output to be %s, got %s", tt.expectedOutput, output)
			}
		})
	}
}

func TestGetSpaceLabel(t *testing.T) {
	testSpaceId := uuid.NewString()
	tests := []struct {
		description    string
		getSpaceFails  bool
		getSpaceResp   *scf.Space
		expectedOutput string
	}{
		{
			description:    "base",
			getSpaceResp:   &scf.Space{Name: "space"},
			expectedOutput: "space",
		},
		{
			description:    "get space fails",
			getSpaceFails:  true,
			expectedOutput: testSpaceId,
		},
		{
			description:    "empty name",
			getSpaceResp:   &scf.Space{},
			expectedOutput: testSpaceId,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			settings := &mockSettings{
				getSpaceFails: tt.getSpaceFails,
				getSpaceResp:  tt.getSpaceResp,
			}
			p := print.NewPrinter(nil, nil, nil)

			output := GetSpaceLabel(context.Background(), p, newApiMock(settings), testProjectId, testRegion, testOrganizationId, testSpaceId)
			if output != tt.expectedOutput {
				t.Errorf("expected output to be %s, got %s", tt.expectedOutput, output)
			}
		})
	}
}

func TestUserLabel(t *testing.T) {
	tests := []struct {
		description    string
		userGuid       *string
		userName       *string
		expectedOutput string
	}{
		{
			description:    "guid",
			userGuid:       utils.Ptr("guid"),
			expectedOutput: "guid",
		},
		{
			description:    "name",
			userName:       utils.Ptr("name"),
			expectedOutput: "name",
		},
		{
			description:    "neither",
			expectedOutput: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			output := UserLabel(tt.userGuid, tt.userName)
			if output != tt.expectedOutput {
				t.Errorf("expected output to be %s, got %s", tt.expectedOutput, output)
			}
		})
	}
}
