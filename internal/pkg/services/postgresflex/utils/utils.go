package utils

import (
	"context"
	"fmt"
	"slices"

	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"

	postgresflex "github.com/stackitcloud/stackit-sdk-go/services/postgresflex/v3api"
	"golang.org/x/mod/semver"
)

// The number of replicas is enforced by the API according to the instance type
var instanceTypeToReplicas = map[string]int32{
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

func LoadFlavorId(cpu, ram int64, flavors []postgresflex.ListFlavors) (*string, error) {
	if flavors == nil {
		return nil, fmt.Errorf("nil flavors")
	}

	for _, f := range flavors {
		if f.Cpu == cpu && f.Memory == ram {
			return &f.Id, nil
		}
	}

	return nil, &errors.DatabaseInvalidFlavorError{
		Service: "postgresflex",
		Details: "You provided an invalid combination for CPU and RAM.",
	}
}

func GetLatestPostgreSQLVersion(ctx context.Context, apiClient postgresflex.DefaultAPI, projectId, region string) (string, error) {
	resp, err := apiClient.ListVersions(ctx, projectId, region).Execute()
	if err != nil {
		return "", fmt.Errorf("get PostgreSQL versions: %w", err)
	}

	latestVersion := "0"
	for i := range resp.Versions {
		oldSemVer := fmt.Sprintf("v%s", latestVersion)
		newSemVer := fmt.Sprintf("v%s", resp.Versions[i].Version)
		if semver.Compare(newSemVer, oldSemVer) != 1 {
			continue
		}
		latestVersion = resp.Versions[i].Version
	}
	if latestVersion == "0" {
		return "", fmt.Errorf("no PostgreSQL versions found")
	}
	return latestVersion, nil
}

func GetInstanceName(ctx context.Context, apiClient postgresflex.DefaultAPI, projectId, region, instanceId string) (string, error) {
	resp, err := apiClient.GetInstance(ctx, projectId, region, instanceId).Execute()
	if err != nil {
		return "", fmt.Errorf("get PostgreSQL Flex instance: %w", err)
	}
	return resp.Name, nil
}

func GetUserName(ctx context.Context, apiClient postgresflex.DefaultAPI, projectId, region, instanceId string, userId int64) (string, error) {
	resp, err := apiClient.GetUser(ctx, projectId, region, instanceId, userId).Execute()
	if err != nil {
		return "", fmt.Errorf("get PostgreSQL Flex user: %w", err)
	}
	return resp.Name, nil
}
