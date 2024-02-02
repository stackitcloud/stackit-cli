package utils

import (
	"context"
	"fmt"
	"strings"

	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"

	"github.com/stackitcloud/stackit-sdk-go/services/mongodbflex"
)

// The number of replicas is enforced by the API according to the instance type
var instanceTypeToReplicas = map[string]int64{
	"Single":  1,
	"Replica": 3,
	"Sharded": 9,
}

func ValidateFlavorId(service, flavorId string, flavors *[]mongodbflex.HandlersInfraFlavor) error {
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

func ValidateStorage(service string, storageClass *string, storageSize *int64, storages *mongodbflex.ListStoragesResponse, flavorId string) error {
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

func LoadFlavorId(service string, cpu, ram int64, flavors *[]mongodbflex.HandlersInfraFlavor) (*string, error) {
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

type MongoDBFlexClient interface {
	GetInstanceExecute(ctx context.Context, projectId, instanceId string) (*mongodbflex.GetInstanceResponse, error)
	GetUserExecute(ctx context.Context, projectId, instanceId, userId string) (*mongodbflex.GetUserResponse, error)
}

func GetInstanceName(ctx context.Context, apiClient MongoDBFlexClient, projectId, instanceId string) (string, error) {
	resp, err := apiClient.GetInstanceExecute(ctx, projectId, instanceId)
	if err != nil {
		return "", fmt.Errorf("get MongoDBFlex instance: %w", err)
	}
	return *resp.Item.Name, nil
}

func GetUserName(ctx context.Context, apiClient MongoDBFlexClient, projectId, instanceId, userId string) (string, error) {
	resp, err := apiClient.GetUserExecute(ctx, projectId, instanceId, userId)
	if err != nil {
		return "", fmt.Errorf("get MongoDBFlex user: %w", err)
	}
	return *resp.Item.Username, nil
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
