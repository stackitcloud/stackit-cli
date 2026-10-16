package utils

import (
	"context"
	"fmt"
	"testing"

	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	mongodbflex "github.com/stackitcloud/stackit-sdk-go/services/mongodbflex/v2api"
)

var (
	testProjectId  = uuid.NewString()
	testInstanceId = uuid.NewString()
	testUserId     = uuid.NewString()
	testBackupId   = uuid.NewString()
)

const (
	testRegion       = "eu02"
	testInstanceName = "instance"
	testUserName     = "user"
)

type mockSettings struct {
	listVersionsFails    bool
	listVersionsResp     *mongodbflex.ListVersionsResponse
	getInstanceFails     bool
	getInstanceResp      *mongodbflex.InstanceResponse
	getUserFails         bool
	getUserResp          *mongodbflex.GetUserResponse
	listRestoreJobsFails bool
	listRestoreJobsResp  *mongodbflex.ListRestoreJobsResponse
}

func newAPIClientMock(m mockSettings) mongodbflex.DefaultAPI {
	return mongodbflex.DefaultAPIServiceMock{
		ListVersionsExecuteMock: utils.Ptr(func(_ mongodbflex.ApiListVersionsRequest) (*mongodbflex.ListVersionsResponse, error) {
			if m.listVersionsFails {
				return nil, fmt.Errorf("could not list versions")
			}
			return m.listVersionsResp, nil
		}),
		ListRestoreJobsExecuteMock: utils.Ptr(func(_ mongodbflex.ApiListRestoreJobsRequest) (*mongodbflex.ListRestoreJobsResponse, error) {
			if m.listRestoreJobsFails {
				return nil, fmt.Errorf("could not list versions")
			}
			return m.listRestoreJobsResp, nil
		}),
		GetInstanceExecuteMock: utils.Ptr(func(_ mongodbflex.ApiGetInstanceRequest) (*mongodbflex.InstanceResponse, error) {
			if m.getInstanceFails {
				return nil, fmt.Errorf("could not get instance")
			}
			return m.getInstanceResp, nil
		}),
		GetUserExecuteMock: utils.Ptr(func(_ mongodbflex.ApiGetUserRequest) (*mongodbflex.GetUserResponse, error) {
			if m.getUserFails {
				return nil, fmt.Errorf("could not get user")
			}
			return m.getUserResp, nil
		}),
	}
}

func TestLoadFlavorId(t *testing.T) {
	tests := []struct {
		description    string
		cpu            int32
		ram            int32
		flavors        *[]mongodbflex.InstanceFlavor
		isValid        bool
		expectedOutput *string
	}{
		{
			description: "base",
			cpu:         2,
			ram:         4,
			flavors: &[]mongodbflex.InstanceFlavor{
				{
					Id:     utils.Ptr("bar-1"),
					Cpu:    utils.Ptr(int32(2)),
					Memory: utils.Ptr(int32(2)),
				},
				{
					Id:     utils.Ptr("bar-2"),
					Cpu:    utils.Ptr(int32(4)),
					Memory: utils.Ptr(int32(4)),
				},
				{
					Id:     utils.Ptr("foo"),
					Cpu:    utils.Ptr(int32(2)),
					Memory: utils.Ptr(int32(4)),
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
			flavors:     &[]mongodbflex.InstanceFlavor{},
			isValid:     false,
		},
		{
			description: "flavors with details missing",
			cpu:         2,
			ram:         4,
			flavors: &[]mongodbflex.InstanceFlavor{
				{
					Id:     utils.Ptr("bar-1"),
					Cpu:    nil,
					Memory: nil,
				},
				{
					Id:     utils.Ptr("bar-2"),
					Cpu:    utils.Ptr(int32(4)),
					Memory: utils.Ptr(int32(4)),
				},
				{
					Id:     utils.Ptr("foo"),
					Cpu:    utils.Ptr(int32(2)),
					Memory: utils.Ptr(int32(4)),
				},
			},
			isValid:        true,
			expectedOutput: utils.Ptr("foo"),
		},
		{
			description: "match with nil id",
			cpu:         2,
			ram:         4,
			flavors: &[]mongodbflex.InstanceFlavor{
				{
					Id:     utils.Ptr("bar-1"),
					Cpu:    utils.Ptr(int32(2)),
					Memory: utils.Ptr(int32(2)),
				},
				{
					Id:     utils.Ptr("bar-2"),
					Cpu:    utils.Ptr(int32(4)),
					Memory: utils.Ptr(int32(4)),
				},
				{
					Id:     nil,
					Cpu:    utils.Ptr(int32(2)),
					Memory: utils.Ptr(int32(4)),
				},
			},
			isValid: false,
		},
		{
			description: "invalid settings",
			cpu:         2,
			ram:         4,
			flavors: &[]mongodbflex.InstanceFlavor{
				{
					Id:     utils.Ptr("bar-1"),
					Cpu:    utils.Ptr(int32(2)),
					Memory: utils.Ptr(int32(2)),
				},
				{
					Id:     utils.Ptr("bar-2"),
					Cpu:    utils.Ptr(int32(4)),
					Memory: utils.Ptr(int32(4)),
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

func TestGetLatestMongoDBFlexVersion(t *testing.T) {
	tests := []struct {
		description       string
		listVersionsFails bool
		listVersionsResp  *mongodbflex.ListVersionsResponse
		isValid           bool
		expectedOutput    string
	}{
		{
			description: "base",
			listVersionsResp: &mongodbflex.ListVersionsResponse{
				Versions: []string{"8", "10", "9"},
			},
			isValid:        true,
			expectedOutput: "10",
		},
		{
			description:       "get instance fails",
			listVersionsFails: true,
			isValid:           false,
		},
		{
			description: "no versions",
			listVersionsResp: &mongodbflex.ListVersionsResponse{
				Versions: []string{},
			},
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			settings := mockSettings{
				listVersionsFails: tt.listVersionsFails,
				listVersionsResp:  tt.listVersionsResp,
			}

			output, err := GetLatestMongoDBVersion(context.Background(), newAPIClientMock(settings), testProjectId, testRegion)

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
		description      string
		getInstanceFails bool
		getInstanceResp  *mongodbflex.InstanceResponse
		isValid          bool
		expectedOutput   string
	}{
		{
			description: "base",
			getInstanceResp: &mongodbflex.InstanceResponse{
				Item: &mongodbflex.Instance{
					Name: utils.Ptr(testInstanceName),
				},
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
			settings := mockSettings{
				getInstanceFails: tt.getInstanceFails,
				getInstanceResp:  tt.getInstanceResp,
			}

			output, err := GetInstanceName(context.Background(), newAPIClientMock(settings), testProjectId, testInstanceId, testRegion)

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
		getUserResp    *mongodbflex.GetUserResponse
		isValid        bool
		expectedOutput string
	}{
		{
			description: "base",
			getUserResp: &mongodbflex.GetUserResponse{
				Item: &mongodbflex.InstanceResponseUser{
					Username: utils.Ptr(testUserName),
				},
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
			settings := mockSettings{
				getUserFails: tt.getUserFails,
				getUserResp:  tt.getUserResp,
			}

			output, err := GetUserName(context.Background(), newAPIClientMock(settings), testProjectId, testInstanceId, testUserId, testRegion)

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

func TestGetRestoreStatus(t *testing.T) {
	tests := []struct {
		description         string
		listRestoreJobsResp *mongodbflex.ListRestoreJobsResponse
		expectedOutput      string
	}{
		{
			description: "base",
			listRestoreJobsResp: &mongodbflex.ListRestoreJobsResponse{
				Items: []mongodbflex.RestoreInstanceStatus{
					{
						BackupID: utils.Ptr(testBackupId),
						Date:     utils.Ptr("2024-05-14T12:01:11Z"),
						Status:   utils.Ptr("state"),
					},
					{
						BackupID: utils.Ptr("bar"),
						Date:     utils.Ptr("2024-05-14T12:01:11Z"),
						Status:   utils.Ptr("state 2"),
					},
				},
			},
			expectedOutput: "state",
		},
		{
			description: "get latest restore, ordered array",
			listRestoreJobsResp: &mongodbflex.ListRestoreJobsResponse{
				Items: []mongodbflex.RestoreInstanceStatus{
					{
						BackupID: utils.Ptr(testBackupId),
						Date:     utils.Ptr("2024-05-14T12:01:11Z"),
						Status:   utils.Ptr("in progress"),
					},
					{
						BackupID: utils.Ptr(testBackupId),
						Date:     utils.Ptr("2024-05-13T12:01:11Z"),
						Status:   utils.Ptr("finished"),
					},
				},
			},
			expectedOutput: "in progress",
		},
		{
			description: "get latest restore, unordered array",
			listRestoreJobsResp: &mongodbflex.ListRestoreJobsResponse{
				Items: []mongodbflex.RestoreInstanceStatus{
					{
						BackupID: utils.Ptr(testBackupId),
						Date:     utils.Ptr("2024-05-13T12:01:11Z"),
						Status:   utils.Ptr("finished"),
					},
					{
						BackupID: utils.Ptr(testBackupId),
						Date:     utils.Ptr("2024-05-14T12:01:11Z"),
						Status:   utils.Ptr("in progress"),
					},
				},
			},
			expectedOutput: "in progress",
		},
		{
			description: "get latest restore, another date format",
			listRestoreJobsResp: &mongodbflex.ListRestoreJobsResponse{
				Items: []mongodbflex.RestoreInstanceStatus{
					{
						BackupID: utils.Ptr(testBackupId),
						Date:     utils.Ptr("2009-11-10 23:00:00 +0000 UTC m=+0.000000001"),
						Status:   utils.Ptr("finished"),
					},
					{
						BackupID: utils.Ptr(testBackupId),
						Date:     utils.Ptr("2009-11-11 23:00:00 +0000 UTC m=+0.000000001"),
						Status:   utils.Ptr("in progress"),
					},
				},
			},
			expectedOutput: "in progress",
		},
		{
			description: "no restore job for that backup",
			listRestoreJobsResp: &mongodbflex.ListRestoreJobsResponse{
				Items: []mongodbflex.RestoreInstanceStatus{
					{
						BackupID: utils.Ptr("bar"),
						Date:     utils.Ptr("2024-05-13T12:01:11Z"),
						Status:   utils.Ptr("in progress"),
					},
					{
						BackupID: utils.Ptr("bar"),
						Date:     utils.Ptr("2024-05-13T12:01:11Z"),
						Status:   utils.Ptr("finished"),
					},
				},
			},
			expectedOutput: "-",
		},
		{
			description: "no restore jobs",
			listRestoreJobsResp: &mongodbflex.ListRestoreJobsResponse{
				Items: nil,
			},
			expectedOutput: "-",
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			output := GetRestoreStatus(testBackupId, tt.listRestoreJobsResp)

			if output != tt.expectedOutput {
				t.Errorf("expected output to be %s, got %s", tt.expectedOutput, output)
			}
		})
	}
}
