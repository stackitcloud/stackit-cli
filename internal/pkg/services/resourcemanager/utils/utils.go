package utils

import (
	"context"
	"fmt"

	resourcemanager "github.com/stackitcloud/stackit-sdk-go/services/resourcemanager/v0api"
)

// GetOrganizationName returns the name of an organization by its ID.
func GetOrganizationName(ctx context.Context, apiClient resourcemanager.DefaultAPI, orgId string) (string, error) {
	resp, err := apiClient.GetOrganization(ctx, orgId).Execute()
	if err != nil {
		return "", fmt.Errorf("get organization details: %w", err)
	}

	return resp.Name, nil
}

func GetProjectName(ctx context.Context, apiClient resourcemanager.DefaultAPI, projectId string) (string, error) {
	resp, err := apiClient.GetProject(ctx, projectId).Execute()
	if err != nil {
		return "", fmt.Errorf("get project details: %w", err)
	}

	return resp.Name, nil
}
