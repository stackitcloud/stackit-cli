package utils

import (
	"context"
	"fmt"

	resourcemanager "github.com/stackitcloud/stackit-sdk-go/services/resourcemanager/v0api"

	rmUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/resourcemanager/utils"
)

// GetResourceLabel returns the name of the given resource
func GetResourceLabel(ctx context.Context, apiClient resourcemanager.DefaultAPI, resourceId, resourceType string) (string, error) {
	var resourceLabel string
	var err error
	switch resourceType {
	case "project":
		resourceLabel, err = rmUtils.GetProjectName(ctx, apiClient, resourceId)
	case "organization":
		resourceLabel, err = rmUtils.GetOrganizationName(ctx, apiClient, resourceId)
	case "folder":
		resourceLabel, err = rmUtils.GetFolderName(ctx, apiClient, resourceId)
	default:
		return "", fmt.Errorf("unknown resource type: %v", resourceType)
	}
	if err != nil {
		return "", fmt.Errorf("get %v name: %w", resourceType, err)
	}

	if resourceLabel == "" {
		resourceLabel = resourceId
	}

	return resourceLabel, nil
}
