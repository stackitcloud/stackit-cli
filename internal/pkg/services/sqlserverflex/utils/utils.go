package utils

import (
	"context"
	"fmt"
	"math"
	"strings"

	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"

	sqlserverflex "github.com/stackitcloud/stackit-sdk-go/services/sqlserverflex/v3api"
)

const (
	ServiceCmd = "beta sqlserverflex"
)

func ValidateFlavorId(flavorId string, flavors []sqlserverflex.ListFlavors) error {
	if flavors == nil {
		return fmt.Errorf("nil flavors")
	}

	for _, f := range flavors {
		if strings.EqualFold(f.Id, flavorId) {
			return nil
		}
	}

	return &errors.DatabaseInvalidFlavorError{
		Service: ServiceCmd,
		Details: fmt.Sprintf("You provided flavor ID '%s', which is invalid.", flavorId),
	}
}

func ValidateStorage(storageClass string, storageSize *int64, storages *sqlserverflex.ListStoragesResponse, flavorId string) error {
	if storages == nil {
		return fmt.Errorf("nil storages")
	}

	if storageSize != nil {
		if *storageSize < int64(storages.StorageRange.Min) || *storageSize > int64(storages.StorageRange.Max) {
			return fmt.Errorf("%s", fmt.Sprintf("You provided storage size '%d', which is invalid. The valid range is %d-%d.", *storageSize, storages.StorageRange.Min, storages.StorageRange.Max))
		}
	}

	for _, sc := range storages.StorageClasses {
		if strings.EqualFold(storageClass, sc.Class) {
			return nil
		}
	}
	return &errors.DatabaseInvalidStorageError{
		Service:  ServiceCmd,
		Details:  fmt.Sprintf("You provided storage class '%s', which is invalid.", storageClass),
		FlavorId: flavorId,
	}
}

func LoadFlavorId(cpu, ram int64, flavors []sqlserverflex.ListFlavors) (string, error) {
	if flavors == nil {
		return "", fmt.Errorf("nil flavors")
	}

	availableFlavors := ""
	for _, f := range flavors {
		if f.Id == "" {
			continue
		}
		if f.Cpu == cpu && f.Memory == ram {
			return f.Id, nil
		}
		availableFlavors = fmt.Sprintf("%s\n- %d CPU, %d GB RAM", availableFlavors, f.Cpu, f.Cpu)
	}
	return "", &errors.DatabaseInvalidFlavorError{
		Service: ServiceCmd,
		Details: "You provided an invalid combination for CPU and RAM.",
	}
}

func ListAllFlavors(ctx context.Context, api sqlserverflex.DefaultAPI, projectId, region string) ([]sqlserverflex.ListFlavors, error) {
	const pageSize = 100
	const sort = sqlserverflex.FLAVORSORT_ID_ASC
	totalRows := int64(math.MaxInt64)
	var page int64 = 1
	var result []sqlserverflex.ListFlavors
	for ; int64(len(result)) < totalRows; page++ {
		flavors, err := api.ListFlavors(ctx, projectId, region).Page(page).Sort(sort).Size(pageSize).Execute()
		if err != nil {
			return result, fmt.Errorf("listing flavors: %w", err)
		}
		result = append(result, flavors.Flavors...)
		totalRows = flavors.Pagination.TotalRows
	}
	return result, nil
}

func GetInstanceName(ctx context.Context, apiClient sqlserverflex.DefaultAPI, projectId, instanceId, region string) (string, error) {
	resp, err := apiClient.GetInstance(ctx, projectId, region, instanceId).Execute()
	if err != nil {
		return "", fmt.Errorf("get SQLServer Flex instance: %w", err)
	}
	return resp.Name, nil
}

func GetUserName(ctx context.Context, apiClient sqlserverflex.DefaultAPI, projectId, instanceId string, userId int64, region string) (string, error) {
	resp, err := apiClient.GetUser(ctx, projectId, region, instanceId, userId).Execute()
	if err != nil {
		return "", fmt.Errorf("get SQLServer Flex user: %w", err)
	}
	return resp.Username, nil
}

func GetFlavor(ctx context.Context, client sqlserverflex.DefaultAPI, projectId, region, flavorId string) (*sqlserverflex.ListFlavors, error) {
	flavors, err := ListAllFlavors(ctx, client, projectId, region)
	if err != nil {
		return nil, fmt.Errorf("failed to list flavors: %w", err)
	}
	for _, flavor := range flavors {
		if flavor.Id == flavorId {
			return &flavor, nil
		}
	}
	return nil, fmt.Errorf("flavor with ID %q not found in project %q", flavorId, projectId)
}
