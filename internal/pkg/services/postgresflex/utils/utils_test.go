package utils

import (
	"context"
	"fmt"
	"testing"

	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	postgresflex "github.com/stackitcloud/stackit-sdk-go/services/postgresflex/v3api"
)

var (
	testProjectId  = uuid.NewString()
	testInstanceId = uuid.NewString()
)

const (
	testInstanceName = "instance"
	testUserName     = "user"
	testRegion       = "eu01"
	testUserId       = int64(102)
)

type mockSettings struct {
	listVersionsFails bool
	listVersionsResp  *postgresflex.ListVersionsResponse
	getInstanceFails  bool
	getInstanceResp   *postgresflex.GetInstanceResponse
	getUserFails      bool
	getUserResp       *postgresflex.GetUserResponse
}

func newAPIMockClient(s mockSettings) postgresflex.DefaultAPI {
	return postgresflex.DefaultAPIServiceMock{
		ListVersionsExecuteMock: utils.Ptr(func(_ postgresflex.ApiListVersionsRequest) (*postgresflex.ListVersionsResponse, error) {
			if s.listVersionsFails {
				return nil, fmt.Errorf("could not list versions")
			}
			return s.listVersionsResp, nil
		}),
		GetInstanceExecuteMock: utils.Ptr(func(_ postgresflex.ApiGetInstanceRequest) (*postgresflex.GetInstanceResponse, error) {
			if s.getInstanceFails {
				return nil, fmt.Errorf("could not get instance")
			}
			return s.getInstanceResp, nil
		}),
		GetUserExecuteMock: utils.Ptr(func(_ postgresflex.ApiGetUserRequest) (*postgresflex.GetUserResponse, error) {
			if s.getUserFails {
				return nil, fmt.Errorf("could not get user")
			}
			return s.getUserResp, nil
		}),
	}
}

func TestLoadFlavorId(t *testing.T) {
	tests := []struct {
		description    string
		cpu            int64
		ram            int64
		flavors        []postgresflex.ListFlavors
		isValid        bool
		expectedOutput *string
	}{
		{
			description: "base",
			cpu:         2,
			ram:         4,
			flavors: []postgresflex.ListFlavors{
				{
					Id:     "bar-1",
					Cpu:    int64(2),
					Memory: int64(2),
				},
				{
					Id:     "bar-2",
					Cpu:    int64(4),
					Memory: int64(4),
				},
				{
					Id:     "foo",
					Cpu:    int64(2),
					Memory: int64(4),
				},
			},
			isValid:        true,
			expectedOutput: utils.Ptr("foo"),
		},
		{
			description: "nil flavors",
			cpu:         2,
			ram:         4,
			flavors:     nil,
			isValid:     false,
		},
		{
			description: "no flavors",
			cpu:         2,
			ram:         4,
			flavors:     []postgresflex.ListFlavors{},
			isValid:     false,
		},
		{
			description: "invalid settings",
			cpu:         2,
			ram:         4,
			flavors: []postgresflex.ListFlavors{
				{
					Id:     "bar-1",
					Cpu:    int64(2),
					Memory: int64(2),
				},
				{
					Id:     "bar-2",
					Cpu:    int64(4),
					Memory: int64(4),
				},
			},
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			output, err := LoadFlavorId(tt.cpu, tt.ram, tt.flavors)

			if !tt.isValid {
				if err == nil {
					t.Fatalf("should have failed")
				}
				return
			}

			if err != nil {
				t.Fatalf("should not have failed: %v", err)
			}
			if output == nil {
				t.Fatalf("returned nil output")
			}
			diff := cmp.Diff(output, tt.expectedOutput)
			if diff != "" {
				t.Fatalf("outputs do not match: %s", diff)
			}
		})
	}
}

func TestGetLatestPostgreSQLVersion(t *testing.T) {
	tests := []struct {
		description        string
		mockClientSettings mockSettings
		isValid            bool
		expectedOutput     string
	}{
		{
			description: "base",
			mockClientSettings: mockSettings{
				listVersionsResp: &postgresflex.ListVersionsResponse{
					Versions: []postgresflex.Version{
						{
							Version: "8",
						},
						{
							Version: "10",
						},
						{
							Version: "9",
						},
					},
				},
			},
			isValid:        true,
			expectedOutput: "10",
		},
		{
			description: "get instance fails",
			mockClientSettings: mockSettings{
				listVersionsFails: true,
			},
			isValid: false,
		},
		{
			description: "no versions",
			mockClientSettings: mockSettings{
				listVersionsResp: &postgresflex.ListVersionsResponse{
					Versions: []postgresflex.Version{},
				},
			},
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			output, err := GetLatestPostgreSQLVersion(context.Background(), newAPIMockClient(tt.mockClientSettings), testProjectId, testRegion)

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

func TestGetInstanceName(t *testing.T) {
	tests := []struct {
		description        string
		mockClientSettings mockSettings
		isValid            bool
		expectedOutput     string
	}{
		{
			description: "base",
			mockClientSettings: mockSettings{
				getInstanceResp: &postgresflex.GetInstanceResponse{
					Name: testInstanceName,
				},
			},
			isValid:        true,
			expectedOutput: testInstanceName,
		},
		{
			description: "get instance fails",
			mockClientSettings: mockSettings{
				getInstanceFails: true,
			},
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			output, err := GetInstanceName(context.Background(), newAPIMockClient(tt.mockClientSettings), testProjectId, testRegion, testInstanceId)

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

func TestGetUserName(t *testing.T) {
	tests := []struct {
		description        string
		mockClientSettings mockSettings
		isValid            bool
		expectedOutput     string
	}{
		{
			description: "base",
			mockClientSettings: mockSettings{
				getUserResp: &postgresflex.GetUserResponse{
					Name: testUserName,
				},
			},
			isValid:        true,
			expectedOutput: testUserName,
		},
		{
			description: "get user fails",
			mockClientSettings: mockSettings{
				getUserFails: true,
			},
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			output, err := GetUserName(context.Background(), newAPIMockClient(tt.mockClientSettings), testProjectId, testRegion, testInstanceId, testUserId)

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
