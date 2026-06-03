package utils

import (
	"context"
	"fmt"
	"strings"

	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"

	sqlserverflex "github.com/stackitcloud/stackit-sdk-go/services/sqlserverflex/v2api"
)

const (
	ServiceCmd = "beta sqlserverflex"
)

// enforce implementation of interfaces
var (
	_ SQLServerFlexClient = sqlserverflex.APIClient{}.DefaultAPI
)

type SQLServerFlexClient interface {
	ListVersions(ctx context.Context, projectId string, region string) sqlserverflex.ApiListVersionsRequest
	GetInstance(ctx context.Context, projectId, instanceId string, region string) sqlserverflex.ApiGetInstanceRequest
	GetUser(ctx context.Context, projectId, instanceId, userId string, region string) sqlserverflex.ApiGetUserRequest
}

func ValidateFlavorId(flavorId string, flavors []sqlserverflex.InstanceFlavorEntry) error {
	if flavors == nil {
		return fmt.Errorf("nil flavors")
	}

	for _, f := range flavors {
		if f.Id != nil && strings.EqualFold(*f.Id, flavorId) {
			return nil
		}
	}

	return &errors.DatabaseInvalidFlavorError{
		Service: ServiceCmd,
		Details: fmt.Sprintf("You provided flavor ID '%s', which is invalid.", flavorId),
	}
}

func ValidateStorage(storageClass *string, storageSize *int64, storages *sqlserverflex.ListStoragesResponse, flavorId string) error {
	if storages == nil {
		return fmt.Errorf("nil storages")
	}

	if storageSize != nil {
		if *storageSize < *storages.StorageRange.Min || *storageSize > *storages.StorageRange.Max {
			return fmt.Errorf("%s", fmt.Sprintf("You provided storage size '%d', which is invalid. The valid range is %d-%d.", *storageSize, *storages.StorageRange.Min, *storages.StorageRange.Max))
		}
	}

	if storageClass == nil {
		return nil
	}

	for _, sc := range storages.StorageClasses {
		if strings.EqualFold(*storageClass, sc) {
			return nil
		}
	}
	return &errors.DatabaseInvalidStorageError{
		Service:  ServiceCmd,
		Details:  fmt.Sprintf("You provided storage class '%s', which is invalid.", *storageClass),
		FlavorId: flavorId,
	}
}

func LoadFlavorId(cpu, ram int32, flavors []sqlserverflex.InstanceFlavorEntry) (*string, error) {
	if flavors == nil {
		return nil, fmt.Errorf("nil flavors")
	}

	availableFlavors := ""
	for _, f := range flavors {
		if f.Id == nil || f.Cpu == nil || f.Memory == nil {
			continue
		}
		if *f.Cpu == cpu && *f.Memory == ram {
			return f.Id, nil
		}
		availableFlavors = fmt.Sprintf("%s\n- %d CPU, %d GB RAM", availableFlavors, *f.Cpu, *f.Cpu)
	}
	return nil, &errors.DatabaseInvalidFlavorError{
		Service: ServiceCmd,
		Details: "You provided an invalid combination for CPU and RAM.",
	}
}

func GetInstanceName(ctx context.Context, apiClient SQLServerFlexClient, projectId, instanceId, region string) (string, error) {
	resp, err := apiClient.GetInstance(ctx, projectId, instanceId, region).Execute()
	if err != nil {
		return "", fmt.Errorf("get SQLServer Flex instance: %w", err)
	}
	return *resp.Item.Name, nil
}

func GetUserName(ctx context.Context, apiClient SQLServerFlexClient, projectId, instanceId, userId, region string) (string, error) {
	resp, err := apiClient.GetUser(ctx, projectId, instanceId, userId, region).Execute()
	if err != nil {
		return "", fmt.Errorf("get SQLServer Flex user: %w", err)
	}
	return *resp.Item.Username, nil
}
