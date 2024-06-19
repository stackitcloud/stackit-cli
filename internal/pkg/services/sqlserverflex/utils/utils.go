package utils

import (
	"context"
	"fmt"
	"strings"

	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"

	"github.com/stackitcloud/stackit-sdk-go/services/sqlserverflex"
)

const (
	ServiceCmd = "beta sqlserverflex"
)

type SQLServerFlexClient interface {
	ListVersionsExecute(ctx context.Context, projectId string) (*sqlserverflex.ListVersionsResponse, error)
	GetInstanceExecute(ctx context.Context, projectId, instanceId string) (*sqlserverflex.GetInstanceResponse, error)
	GetUserExecute(ctx context.Context, projectId, instanceId, userId string) (*sqlserverflex.GetUserResponse, error)
}

func ValidateFlavorId(flavorId string, flavors *[]sqlserverflex.InstanceFlavorEntry) error {
	if flavors == nil {
		return fmt.Errorf("nil flavors")
	}

	for _, f := range *flavors {
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

	for _, sc := range *storages.StorageClasses {
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

func LoadFlavorId(cpu, ram int64, flavors *[]sqlserverflex.InstanceFlavorEntry) (*string, error) {
	if flavors == nil {
		return nil, fmt.Errorf("nil flavors")
	}

	availableFlavors := ""
	for _, f := range *flavors {
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

func GetInstanceName(ctx context.Context, apiClient SQLServerFlexClient, projectId, instanceId string) (string, error) {
	resp, err := apiClient.GetInstanceExecute(ctx, projectId, instanceId)
	if err != nil {
		return "", fmt.Errorf("get SQLServer Flex instance: %w", err)
	}
	return *resp.Item.Name, nil
}

func GetUserName(ctx context.Context, apiClient SQLServerFlexClient, projectId, instanceId, userId string) (string, error) {
	resp, err := apiClient.GetUserExecute(ctx, projectId, instanceId, userId)
	if err != nil {
		return "", fmt.Errorf("get SQLServer Flex user: %w", err)
	}
	return *resp.Item.Username, nil
}
