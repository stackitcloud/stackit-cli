package utils

import (
	"context"
	"fmt"
	"testing"

	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"github.com/stackitcloud/stackit-sdk-go/services/postgresflex"
)

var (
	testProjectId  = uuid.NewString()
	testInstanceId = uuid.NewString()
)

const (
	testInstanceName = "instance"
)

type postgresFlexClientMocked struct {
	listVersionsFails bool
	listVersionsResp  *postgresflex.ListVersionsResponse
	getInstanceFails  bool
	getInstanceResp   *postgresflex.InstanceResponse
}

func (m *postgresFlexClientMocked) ListVersionsExecute(_ context.Context, _ string) (*postgresflex.ListVersionsResponse, error) {
	if m.listVersionsFails {
		return nil, fmt.Errorf("could not list versions")
	}
	return m.listVersionsResp, nil
}

func (m *postgresFlexClientMocked) GetInstanceExecute(_ context.Context, _, _ string) (*postgresflex.InstanceResponse, error) {
	if m.getInstanceFails {
		return nil, fmt.Errorf("could not get instance")
	}
	return m.getInstanceResp, nil
}

func TestValidateStorage(t *testing.T) {
	tests := []struct {
		description  string
		storageClass *string
		storageSize  *int64
		storages     *postgresflex.ListStoragesResponse
		isValid      bool
	}{
		{
			description:  "base",
			storageClass: utils.Ptr("foo"),
			storageSize:  utils.Ptr(int64(10)),
			storages: &postgresflex.ListStoragesResponse{
				StorageClasses: &[]string{"bar-1", "bar-2", "foo"},
				StorageRange: &postgresflex.StorageRange{
					Min: utils.Ptr(int64(5)),
					Max: utils.Ptr(int64(20)),
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
			storages: &postgresflex.ListStoragesResponse{
				StorageClasses: &[]string{"bar-1", "bar-2", "foo"},
				StorageRange: &postgresflex.StorageRange{
					Min: utils.Ptr(int64(5)),
					Max: utils.Ptr(int64(20)),
				},
			},
			isValid: false,
		},
		{
			description:  "storage size out of range 2",
			storageClass: utils.Ptr("foo"),
			storageSize:  utils.Ptr(int64(200)),
			storages: &postgresflex.ListStoragesResponse{
				StorageClasses: &[]string{"bar-1", "bar-2", "foo"},
				StorageRange: &postgresflex.StorageRange{
					Min: utils.Ptr(int64(5)),
					Max: utils.Ptr(int64(20)),
				},
			},
			isValid: false,
		},
		{
			description:  "storage size in range limit 1",
			storageClass: utils.Ptr("foo"),
			storageSize:  utils.Ptr(int64(5)),
			storages: &postgresflex.ListStoragesResponse{
				StorageClasses: &[]string{"bar-1", "bar-2", "foo"},
				StorageRange: &postgresflex.StorageRange{
					Min: utils.Ptr(int64(5)),
					Max: utils.Ptr(int64(20)),
				},
			},
			isValid: true,
		},
		{
			description:  "storage size in range limit 2",
			storageClass: utils.Ptr("foo"),
			storageSize:  utils.Ptr(int64(20)),
			storages: &postgresflex.ListStoragesResponse{
				StorageClasses: &[]string{"bar-1", "bar-2", "foo"},
				StorageRange: &postgresflex.StorageRange{
					Min: utils.Ptr(int64(5)),
					Max: utils.Ptr(int64(20)),
				},
			},
			isValid: true,
		},
		{
			description:  "invalid storage",
			storageClass: utils.Ptr("foo"),
			storageSize:  utils.Ptr(int64(10)),
			storages: &postgresflex.ListStoragesResponse{
				StorageClasses: &[]string{"bar-1", "bar-2", "bar-3"},
				StorageRange: &postgresflex.StorageRange{
					Min: utils.Ptr(int64(5)),
					Max: utils.Ptr(int64(20)),
				},
			},
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			err := ValidateStorage(tt.storageClass, tt.storageSize, tt.storages, "flavor-id")
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
		flavors     *[]postgresflex.Flavor
		isValid     bool
	}{
		{
			description: "base",
			flavorId:    "foo",
			flavors: &[]postgresflex.Flavor{
				{Id: utils.Ptr("bar-1")},
				{Id: utils.Ptr("bar-2")},
				{Id: utils.Ptr("foo")},
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
			flavors:     &[]postgresflex.Flavor{},
			isValid:     false,
		},
		{
			description: "nil flavor id",
			flavorId:    "foo",
			flavors: &[]postgresflex.Flavor{
				{Id: utils.Ptr("bar-1")},
				{Id: nil},
				{Id: utils.Ptr("foo")},
			},
			isValid: true,
		},
		{
			description: "invalid flavor",
			flavorId:    "foo",
			flavors: &[]postgresflex.Flavor{
				{Id: utils.Ptr("bar-1")},
				{Id: utils.Ptr("bar-2")},
				{Id: utils.Ptr("bar-3")},
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
		flavors        *[]postgresflex.Flavor
		isValid        bool
		expectedOutput *string
	}{
		{
			description: "base",
			cpu:         2,
			ram:         4,
			flavors: &[]postgresflex.Flavor{
				{
					Id:     utils.Ptr("bar-1"),
					Cpu:    utils.Ptr(int64(2)),
					Memory: utils.Ptr(int64(2)),
				},
				{
					Id:     utils.Ptr("bar-2"),
					Cpu:    utils.Ptr(int64(4)),
					Memory: utils.Ptr(int64(4)),
				},
				{
					Id:     utils.Ptr("foo"),
					Cpu:    utils.Ptr(int64(2)),
					Memory: utils.Ptr(int64(4)),
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
			flavors:     &[]postgresflex.Flavor{},
			isValid:     false,
		},
		{
			description: "flavors with details missing",
			cpu:         2,
			ram:         4,
			flavors: &[]postgresflex.Flavor{
				{
					Id:     utils.Ptr("bar-1"),
					Cpu:    nil,
					Memory: nil,
				},
				{
					Id:     utils.Ptr("bar-2"),
					Cpu:    utils.Ptr(int64(4)),
					Memory: utils.Ptr(int64(4)),
				},
				{
					Id:     utils.Ptr("foo"),
					Cpu:    utils.Ptr(int64(2)),
					Memory: utils.Ptr(int64(4)),
				},
			},
			isValid:        true,
			expectedOutput: utils.Ptr("foo"),
		},
		{
			description: "match with nil id",
			cpu:         2,
			ram:         4,
			flavors: &[]postgresflex.Flavor{
				{
					Id:     utils.Ptr("bar-1"),
					Cpu:    utils.Ptr(int64(2)),
					Memory: utils.Ptr(int64(2)),
				},
				{
					Id:     utils.Ptr("bar-2"),
					Cpu:    utils.Ptr(int64(4)),
					Memory: utils.Ptr(int64(4)),
				},
				{
					Id:     nil,
					Cpu:    utils.Ptr(int64(2)),
					Memory: utils.Ptr(int64(4)),
				},
			},
			isValid: false,
		},
		{
			description: "invalid settings",
			cpu:         2,
			ram:         4,
			flavors: &[]postgresflex.Flavor{
				{
					Id:     utils.Ptr("bar-1"),
					Cpu:    utils.Ptr(int64(2)),
					Memory: utils.Ptr(int64(2)),
				},
				{
					Id:     utils.Ptr("bar-2"),
					Cpu:    utils.Ptr(int64(4)),
					Memory: utils.Ptr(int64(4)),
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
		description       string
		listVersionsFails bool
		listVersionsResp  *postgresflex.ListVersionsResponse
		isValid           bool
		expectedOutput    string
	}{
		{
			description: "base",
			listVersionsResp: &postgresflex.ListVersionsResponse{
				Versions: &[]string{"8", "10", "9"},
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
			listVersionsResp: &postgresflex.ListVersionsResponse{
				Versions: &[]string{},
			},
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			client := &postgresFlexClientMocked{
				listVersionsFails: tt.listVersionsFails,
				listVersionsResp:  tt.listVersionsResp,
			}

			output, err := GetLatestPostgreSQLVersion(context.Background(), client, testProjectId)

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
		getInstanceResp  *postgresflex.InstanceResponse
		isValid          bool
		expectedOutput   string
	}{
		{
			description: "base",
			getInstanceResp: &postgresflex.InstanceResponse{
				Item: &postgresflex.Instance{
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
			client := &postgresFlexClientMocked{
				getInstanceFails: tt.getInstanceFails,
				getInstanceResp:  tt.getInstanceResp,
			}

			output, err := GetInstanceName(context.Background(), client, testProjectId, testInstanceId)

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
