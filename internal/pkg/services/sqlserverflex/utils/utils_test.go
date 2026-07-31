package utils

import (
	"context"
	"fmt"
	"testing"

	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	sqlserverflex "github.com/stackitcloud/stackit-sdk-go/services/sqlserverflex/v3api"
)

var (
	testProjectId  = uuid.NewString()
	testInstanceId = uuid.NewString()
	testUserId     = int64(123123)
)

const (
	testInstanceName = "instance"
	testUserName     = "user"
	testRegion       = "eu01"
)

type mockSettings struct {
	listFlavorsFails     bool
	listFlavorsResp      *sqlserverflex.ListFlavorsResponse
	listFlavorsResps     []*sqlserverflex.ListFlavorsResponse
	listFlavorsCallCount int
	listVersionsFails    bool
	listVersionsResp     *sqlserverflex.ListVersionsResponse
	getInstanceFails     bool
	getInstanceResp      *sqlserverflex.GetInstanceResponse
	getUserFails         bool
	getUserResp          *sqlserverflex.GetUserResponse
	listRestoreJobsFails bool
	listRestoreJobsResp  *sqlserverflex.ListCurrentRunningRestoreJobs
}

func newApiMock(s *mockSettings) sqlserverflex.DefaultAPI {
	return &sqlserverflex.DefaultAPIServiceMock{
		ListFlavorsExecuteMock: utils.Ptr(func(_ sqlserverflex.ApiListFlavorsRequest) (*sqlserverflex.ListFlavorsResponse, error) {
			if s.listFlavorsFails {
				return nil, fmt.Errorf("could not list flavors")
			}
			if len(s.listFlavorsResps) > 0 {
				if s.listFlavorsCallCount >= len(s.listFlavorsResps) {
					return nil, fmt.Errorf("no more mock responses")
				}
				resp := s.listFlavorsResps[s.listFlavorsCallCount]
				s.listFlavorsCallCount++
				return resp, nil
			}
			return s.listFlavorsResp, nil
		}),
		ListVersionsExecuteMock: utils.Ptr(func(_ sqlserverflex.ApiListVersionsRequest) (*sqlserverflex.ListVersionsResponse, error) {
			if s.listVersionsFails {
				return nil, fmt.Errorf("could not list versions")
			}
			return s.listVersionsResp, nil
		}),
		GetInstanceExecuteMock: utils.Ptr(func(_ sqlserverflex.ApiGetInstanceRequest) (*sqlserverflex.GetInstanceResponse, error) {
			if s.getInstanceFails {
				return nil, fmt.Errorf("could not get instance")
			}
			return s.getInstanceResp, nil
		}),
		GetUserExecuteMock: utils.Ptr(func(_ sqlserverflex.ApiGetUserRequest) (*sqlserverflex.GetUserResponse, error) {
			if s.getUserFails {
				return nil, fmt.Errorf("could not get user")
			}
			return s.getUserResp, nil
		}),
		ListCurrentRunningRestoreJobsExecuteMock: utils.Ptr(func(_ sqlserverflex.ApiListCurrentRunningRestoreJobsRequest) (*sqlserverflex.ListCurrentRunningRestoreJobs, error) {
			if s.listRestoreJobsFails {
				return nil, fmt.Errorf("could not list versions")
			}
			return s.listRestoreJobsResp, nil
		}),
	}
}

func TestListAllFlavors(t *testing.T) {
	tests := []struct {
		description      string
		listFlavorsFails bool
		listFlavorsResp  *sqlserverflex.ListFlavorsResponse
		listFlavorsResps []*sqlserverflex.ListFlavorsResponse
		isValid          bool
		expectedOutput   []sqlserverflex.ListFlavors
	}{
		{
			description: "base",
			listFlavorsResp: &sqlserverflex.ListFlavorsResponse{
				Flavors: []sqlserverflex.ListFlavors{
					{
						Id: "flavor-1",
					},
					{
						Id: "flavor-2",
					},
				},
				Pagination: sqlserverflex.Pagination{
					TotalRows: 2,
				},
			},
			isValid: true,
			expectedOutput: []sqlserverflex.ListFlavors{
				{
					Id: "flavor-1",
				},
				{
					Id: "flavor-2",
				},
			},
		},
		{
			description: "multiple pages",
			listFlavorsResps: []*sqlserverflex.ListFlavorsResponse{
				{
					Flavors: []sqlserverflex.ListFlavors{
						{
							Id: "flavor-1",
						},
						{
							Id: "flavor-2",
						},
					},
					Pagination: sqlserverflex.Pagination{
						TotalRows: 3,
					},
				},
				{
					Flavors: []sqlserverflex.ListFlavors{
						{
							Id: "flavor-3",
						},
					},
					Pagination: sqlserverflex.Pagination{
						TotalRows: 3,
					},
				},
			},
			isValid: true,
			expectedOutput: []sqlserverflex.ListFlavors{
				{
					Id: "flavor-1",
				},
				{
					Id: "flavor-2",
				},
				{
					Id: "flavor-3",
				},
			},
		},
		{
			description: "empty response",
			listFlavorsResp: &sqlserverflex.ListFlavorsResponse{
				Flavors: []sqlserverflex.ListFlavors{},
				Pagination: sqlserverflex.Pagination{
					TotalRows: 0,
				},
			},
			isValid:        true,
			expectedOutput: nil,
		},
		{
			description:      "list flavors fails",
			listFlavorsFails: true,
			isValid:          false,
		},
		{
			description: "page 2 fails",
			listFlavorsResps: []*sqlserverflex.ListFlavorsResponse{
				{
					Flavors: []sqlserverflex.ListFlavors{
						{
							Id: "flavor-1",
						},
					},
					Pagination: sqlserverflex.Pagination{
						TotalRows: 2,
					},
				},
			},
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			settings := &mockSettings{
				listFlavorsFails: tt.listFlavorsFails,
				listFlavorsResp:  tt.listFlavorsResp,
				listFlavorsResps: tt.listFlavorsResps,
			}

			output, err := ListAllFlavors(context.Background(), newApiMock(settings), testProjectId, testRegion)

			if tt.isValid && err != nil {
				t.Fatalf("failed on valid input: %v", err)
			}
			if !tt.isValid && err == nil {
				t.Fatalf("did not fail on invalid input")
			}
			if !tt.isValid {
				return
			}
			diff := cmp.Diff(output, tt.expectedOutput)
			if diff != "" {
				t.Fatalf("outputs do not match: %s", diff)
			}
		})
	}
}

func TestValidateStorage(t *testing.T) {
	tests := []struct {
		description  string
		storageClass *string
		storageSize  *int64
		storages     *sqlserverflex.ListStoragesResponse
		isValid      bool
	}{
		{
			description:  "base",
			storageClass: utils.Ptr("foo"),
			storageSize:  utils.Ptr(int64(10)),
			storages: &sqlserverflex.ListStoragesResponse{
				StorageClasses: []sqlserverflex.FlavorStorageClassesStorageClass{
					{Class: "bar-1"},
					{Class: "bar-2"},
					{Class: "foo"},
				},
				StorageRange: sqlserverflex.FlavorStorageRange{
					Min: int32(5),
					Max: int32(20),
				},
			},
			isValid: true,
		},
		{
			description:  "nil response",
			storageClass: utils.Ptr("foo"),
			storageSize:  utils.Ptr(int64(10)),
			storages:     nil,
			isValid:      false,
		},
		{
			description:  "storage size out of range 1",
			storageClass: utils.Ptr("foo"),
			storageSize:  utils.Ptr(int64(1)),
			storages: &sqlserverflex.ListStoragesResponse{
				StorageClasses: []sqlserverflex.FlavorStorageClassesStorageClass{
					{Class: "bar-1"},
					{Class: "bar-2"},
					{Class: "foo"},
				},
				StorageRange: sqlserverflex.FlavorStorageRange{
					Min: int32(5),
					Max: int32(20),
				},
			},
			isValid: false,
		},
		{
			description:  "storage size out of range 2",
			storageClass: utils.Ptr("foo"),
			storageSize:  utils.Ptr(int64(200)),
			storages: &sqlserverflex.ListStoragesResponse{
				StorageClasses: []sqlserverflex.FlavorStorageClassesStorageClass{
					{Class: "bar-1"},
					{Class: "bar-2"},
					{Class: "foo"},
				},
				StorageRange: sqlserverflex.FlavorStorageRange{
					Min: int32(5),
					Max: int32(20),
				},
			},
			isValid: false,
		},
		{
			description:  "storage size in range limit 1",
			storageClass: utils.Ptr("foo"),
			storageSize:  utils.Ptr(int64(5)),
			storages: &sqlserverflex.ListStoragesResponse{
				StorageClasses: []sqlserverflex.FlavorStorageClassesStorageClass{
					{Class: "bar-1"},
					{Class: "bar-2"},
					{Class: "foo"},
				},
				StorageRange: sqlserverflex.FlavorStorageRange{
					Min: int32(5),
					Max: int32(20),
				},
			},
			isValid: true,
		},
		{
			description:  "storage size in range limit 2",
			storageClass: utils.Ptr("foo"),
			storageSize:  utils.Ptr(int64(20)),
			storages: &sqlserverflex.ListStoragesResponse{
				StorageClasses: []sqlserverflex.FlavorStorageClassesStorageClass{
					{Class: "bar-1"},
					{Class: "bar-2"},
					{Class: "foo"},
				},
				StorageRange: sqlserverflex.FlavorStorageRange{
					Min: int32(5),
					Max: int32(20),
				},
			},
			isValid: true,
		},
		{
			description:  "invalid storage",
			storageClass: utils.Ptr("foo"),
			storageSize:  utils.Ptr(int64(10)),
			storages: &sqlserverflex.ListStoragesResponse{
				StorageClasses: []sqlserverflex.FlavorStorageClassesStorageClass{
					{Class: "bar-1"},
					{Class: "bar-2"},
					{Class: "bar-3"},
				},
				StorageRange: sqlserverflex.FlavorStorageRange{
					Min: int32(5),
					Max: int32(20),
				},
			},
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			err := ValidateStorage(*tt.storageClass, tt.storageSize, tt.storages, "flavor-id")
			if tt.isValid && err != nil {
				t.Fatalf("should not have failed: %v", err)
			}
			if !tt.isValid && err == nil {
				t.Fatalf("should have failed")
			}
		})
	}
}

func TestValidateFlavorId(t *testing.T) {
	tests := []struct {
		description string
		flavorId    string
		flavors     []sqlserverflex.ListFlavors
		isValid     bool
	}{
		{
			description: "base",
			flavorId:    "foo",
			flavors: []sqlserverflex.ListFlavors{
				{Id: "bar-1"},
				{Id: "bar-2"},
				{Id: "foo"},
			},
			isValid: true,
		},
		{
			description: "nil flavors",
			flavorId:    "foo",
			flavors:     nil,
			isValid:     false,
		},
		{
			description: "no flavors",
			flavorId:    "foo",
			flavors:     []sqlserverflex.ListFlavors{},
			isValid:     false,
		},
		{
			description: "empty flavor id",
			flavorId:    "foo",
			flavors: []sqlserverflex.ListFlavors{
				{Id: "bar-1"},
				{Id: ""},
				{Id: "foo"},
			},
			isValid: true,
		},
		{
			description: "invalid flavor",
			flavorId:    "foo",
			flavors: []sqlserverflex.ListFlavors{
				{Id: "bar-1"},
				{Id: "bar-2"},
				{Id: "bar-3"},
			},
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			err := ValidateFlavorId(tt.flavorId, tt.flavors)
			if tt.isValid && err != nil {
				t.Fatalf("should not have failed: %v", err)
			}
			if !tt.isValid && err == nil {
				t.Fatalf("should have failed")
			}
		})
	}
}

func TestLoadFlavorId(t *testing.T) {
	tests := []struct {
		description    string
		cpu            int64
		ram            int64
		flavors        []sqlserverflex.ListFlavors
		isValid        bool
		expectedOutput string
	}{
		{
			description: "base",
			cpu:         2,
			ram:         4,
			flavors: []sqlserverflex.ListFlavors{
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
			expectedOutput: "foo",
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
			flavors:     []sqlserverflex.ListFlavors{},
			isValid:     false,
		},
		{
			description: "flavors with details missing",
			cpu:         2,
			ram:         4,
			flavors: []sqlserverflex.ListFlavors{
				{
					Id:     "bar-1",
					Cpu:    0,
					Memory: 0,
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
			expectedOutput: "foo",
		},
		{
			description: "match with nil id",
			cpu:         2,
			ram:         4,
			flavors: []sqlserverflex.ListFlavors{
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
					Id:     "",
					Cpu:    int64(2),
					Memory: int64(4),
				},
			},
			isValid: false,
		},
		{
			description: "invalid settings",
			cpu:         2,
			ram:         4,
			flavors: []sqlserverflex.ListFlavors{
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
			if output == "" {
				t.Fatalf("returned nil output")
			}
			diff := cmp.Diff(output, tt.expectedOutput)
			if diff != "" {
				t.Fatalf("outputs do not match: %s", diff)
			}
		})
	}
}

func TestGetInstanceName(t *testing.T) {
	tests := []struct {
		description      string
		getInstanceFails bool
		getInstanceResp  *sqlserverflex.GetInstanceResponse
		isValid          bool
		expectedOutput   string
	}{
		{
			description: "base",
			getInstanceResp: &sqlserverflex.GetInstanceResponse{
				Name: testInstanceName,
			},
			isValid:        true,
			expectedOutput: testInstanceName,
		},
		{
			description:      "get instance fails",
			getInstanceFails: true,
			isValid:          false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			settings := &mockSettings{
				getInstanceFails: tt.getInstanceFails,
				getInstanceResp:  tt.getInstanceResp,
			}

			output, err := GetInstanceName(context.Background(), newApiMock(settings), testProjectId, testInstanceId, testRegion)

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
		description    string
		getUserFails   bool
		getUserResp    *sqlserverflex.GetUserResponse
		isValid        bool
		expectedOutput string
	}{
		{
			description: "base",
			getUserResp: &sqlserverflex.GetUserResponse{
				Username: testUserName,
			},
			isValid:        true,
			expectedOutput: testUserName,
		},
		{
			description:  "get user fails",
			getUserFails: true,
			isValid:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			settings := &mockSettings{
				getUserFails: tt.getUserFails,
				getUserResp:  tt.getUserResp,
			}

			output, err := GetUserName(context.Background(), newApiMock(settings), testProjectId, testInstanceId, testUserId, testRegion)

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
