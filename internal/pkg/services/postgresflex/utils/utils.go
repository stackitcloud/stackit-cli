package utils

import (
	"context"
	"fmt"
	"stackit/internal/pkg/errors"
	"strings"

	"github.com/stackitcloud/stackit-sdk-go/services/postgresflex"
)

func ValidateFlavorId(service, flavorId string, flavors *[]postgresflex.Flavor) error {
	for _, f := range *flavors {
		if f.Id != nil && strings.EqualFold(*f.Id, flavorId) {
			return nil
		}
	}

	return &errors.DatabaseInvalidFlavorError{
		Service: service,
		Details: fmt.Sprintf("You provided flavor ID '%s', which is invalid.", flavorId),
	}
}

func ValidateStorage(service string, storageClass *string, storageSize *int64, storages *postgresflex.ListStoragesResponse, flavorId string) error {
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
		Service:  service,
		Details:  fmt.Sprintf("You provided storage class '%s', which is invalid.", *storageClass),
		FlavorId: flavorId,
	}
}

func LoadFlavorId(service string, cpu, ram int64, flavors *[]postgresflex.Flavor) (*string, error) {
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
		Service: service,
		Details: "You provided an invalid combination for CPU and RAM.",
	}
}

type PostgreSQLFlexClient interface {
	GetInstanceExecute(ctx context.Context, projectId, instanceId string) (*postgresflex.InstanceResponse, error)
}

func GetInstanceName(ctx context.Context, apiClient PostgreSQLFlexClient, projectId, instanceId string) (string, error) {
	resp, err := apiClient.GetInstanceExecute(ctx, projectId, instanceId)
	if err != nil {
		return "", fmt.Errorf("get PostgreSQL instance: %w", err)
	}
	return *resp.Item.Name, nil
}
