package utils

import (
	"cmp"
	"context"
	"fmt"
	"slices"

	"golang.org/x/mod/semver"

	mongodbflex "github.com/stackitcloud/stackit-sdk-go/services/mongodbflex/v2api"

	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
)

var instanceTypeToReplicas = map[string]int32{
	"Single":  1,
	"Replica": 3,
	"Sharded": 9,
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

func GetInstanceReplicas(instanceType string) (int32, error) {
	numReplicas, ok := instanceTypeToReplicas[instanceType]
	if !ok {
		return 0, fmt.Errorf("invalid instance type: %v", instanceType)
	}
	return numReplicas, nil
}

// Deprecated: Will be removed after 2027-03-07
func LoadFlavorId(cpu, ram int32, flavors *[]mongodbflex.InstanceFlavor) (*string, error) {
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
		Service: "mongodbflex",
		Details: "You provided an invalid combination for CPU and RAM.",
	}
}

// Deprecated: Will be removed after 2027-03-07.
func GetLatestMongoDBVersion(ctx context.Context, apiClient mongodbflex.DefaultAPI, projectId, region string) (string, error) {
	resp, err := apiClient.ListVersions(ctx, projectId, region).Execute()
	if err != nil {
		return "", fmt.Errorf("get MongoDB versions: %w", err)
	}

	latestVersion := "0"
	for i := range resp.Versions {
		oldSemVer := fmt.Sprintf("v%s", latestVersion)
		newSemVer := fmt.Sprintf("v%s", resp.Versions[i])
		if semver.Compare(newSemVer, oldSemVer) != 1 {
			continue
		}
		latestVersion = resp.Versions[i]
	}
	if latestVersion == "0" {
		return "", fmt.Errorf("no MongoDB versions found")
	}
	return latestVersion, nil
}

func GetInstanceName(ctx context.Context, apiClient mongodbflex.DefaultAPI, projectId, instanceId, region string) (string, error) {
	resp, err := apiClient.GetInstance(ctx, projectId, instanceId, region).Execute()
	if err != nil {
		return "", fmt.Errorf("get MongoDB Flex instance: %w", err)
	}
	return *resp.Item.Name, nil
}

func GetUserName(ctx context.Context, apiClient mongodbflex.DefaultAPI, projectId, instanceId, userId, region string) (string, error) {
	resp, err := apiClient.GetUser(ctx, projectId, instanceId, userId, region).Execute()
	if err != nil {
		return "", fmt.Errorf("get MongoDB Flex user: %w", err)
	}
	return *resp.Item.Username, nil
}

func GetRestoreStatus(backupId string, restoreJobs *mongodbflex.ListRestoreJobsResponse) string {
	state := "-"
	if restoreJobs.Items == nil {
		return state
	}

	restoreJobsSlice := restoreJobs.Items

	// sort array by descending date
	slices.SortFunc(restoreJobsSlice, func(i, j mongodbflex.RestoreInstanceStatus) int {
		// swap elements to sort by descending order
		return cmp.Compare(*j.Date, *i.Date)
	})

	for _, restoreJob := range restoreJobs.Items {
		if *restoreJob.BackupID == backupId {
			state = *restoreJob.Status
			break
		}
	}
	return state
}
