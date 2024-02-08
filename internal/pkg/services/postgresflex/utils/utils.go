package utils

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"

	"github.com/stackitcloud/stackit-sdk-go/services/postgresflex"
	"golang.org/x/mod/semver"
)

// The number of replicas is enforced by the API according to the instance type
var instanceTypeToReplicas = map[string]int64{
	"Single":  1,
	"Replica": 3,
}

func AvailableInstanceTypes() []string {
	instanceTypes := make([]string, len(instanceTypeToReplicas))
	i := 0
	for k := range instanceTypeToReplicas {
		instanceTypes[i] = k
		i++
	}
	// Dict keys aren't iterated in a consistent order
	// So we sort the array to make the output consistent
	slices.Sort(instanceTypes)
	return instanceTypes
}

func GetInstanceReplicas(instanceType string) (int64, error) {
	numReplicas, ok := instanceTypeToReplicas[instanceType]
	if !ok {
		return 0, fmt.Errorf("invalid instance type: %v", instanceType)
	}
	return numReplicas, nil
}

func GetInstanceType(numReplicas int64) (string, error) {
	for k, v := range instanceTypeToReplicas {
		if v == numReplicas {
			return k, nil
		}
	}
	return "", fmt.Errorf("invalid number of replicas: %v", numReplicas)
}

func ValidateFlavorId(flavorId string, flavors *[]postgresflex.Flavor) error {
	if flavors == nil {
		return fmt.Errorf("nil flavors")
	}

	for _, f := range *flavors {
		if f.Id != nil && strings.EqualFold(*f.Id, flavorId) {
			return nil
		}
	}

	return &errors.DatabaseInvalidFlavorError{
		Service: "postgresflex",
		Details: fmt.Sprintf("You provided flavor ID '%s', which is invalid.", flavorId),
	}
}

func ValidateStorage(storageClass *string, storageSize *int64, storages *postgresflex.ListStoragesResponse, flavorId string) error {
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
		Service:  "postgresflex",
		Details:  fmt.Sprintf("You provided storage class '%s', which is invalid.", *storageClass),
		FlavorId: flavorId,
	}
}

func LoadFlavorId(cpu, ram int64, flavors *[]postgresflex.Flavor) (*string, error) {
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
		Service: "postgresflex",
		Details: "You provided an invalid combination for CPU and RAM.",
	}
}

type PostgresFlexClient interface {
	ListVersionsExecute(ctx context.Context, projectId string) (*postgresflex.ListVersionsResponse, error)
	GetInstanceExecute(ctx context.Context, projectId, instanceId string) (*postgresflex.InstanceResponse, error)
	GetUserExecute(ctx context.Context, projectId, instanceId, userId string) (*postgresflex.GetUserResponse, error)
}

func GetLatestPostgreSQLVersion(ctx context.Context, apiClient PostgresFlexClient, projectId string) (string, error) {
	resp, err := apiClient.ListVersionsExecute(ctx, projectId)
	if err != nil {
		return "", fmt.Errorf("get PostgreSQL versions: %w", err)
	}
	versions := *resp.Versions

	latestVersion := "0"
	for i := range versions {
		oldSemVer := fmt.Sprintf("v%s", latestVersion)
		newSemVer := fmt.Sprintf("v%s", versions[i])
		if semver.Compare(newSemVer, oldSemVer) != 1 {
			continue
		}
		latestVersion = versions[i]
	}
	if latestVersion == "0" {
		return "", fmt.Errorf("no PostgreSQL versions found")
	}
	return latestVersion, nil
}

func GetInstanceName(ctx context.Context, apiClient PostgresFlexClient, projectId, instanceId string) (string, error) {
	resp, err := apiClient.GetInstanceExecute(ctx, projectId, instanceId)
	if err != nil {
		return "", fmt.Errorf("get PostgreSQL Flex instance: %w", err)
	}
	return *resp.Item.Name, nil
}

func GetUserName(ctx context.Context, apiClient PostgresFlexClient, projectId, instanceId, userId string) (string, error) {
	resp, err := apiClient.GetUserExecute(ctx, projectId, instanceId, userId)
	if err != nil {
		return "", fmt.Errorf("get PostgreSQL Flex user: %w", err)
	}
	return *resp.Item.Username, nil
}
