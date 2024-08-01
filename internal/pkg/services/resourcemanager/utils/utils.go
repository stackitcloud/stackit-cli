package utils

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-sdk-go/services/resourcemanager"
)

type ResourceManagerClient interface {
	GetOrganizationExecute(ctx context.Context, organizationId string) (*resourcemanager.OrganizationResponse, error)
	GetProjectExecute(ctx context.Context, projectId string) (*resourcemanager.GetProjectResponse, error)
}

// GetOrganizationName returns the name of an organization by its ID.
func GetOrganizationName(ctx context.Context, apiClient ResourceManagerClient, orgId string) (string, error) {
	resp, err := apiClient.GetOrganizationExecute(ctx, orgId)
	if err != nil {
		return "", fmt.Errorf("get organization details: %w", err)
	}

	return *resp.Name, nil
}

func GetProjectName(ctx context.Context, apiClient ResourceManagerClient, projectId string) (string, error) {
	resp, err := apiClient.GetProjectExecute(ctx, projectId)
	if err != nil {
		return "", fmt.Errorf("get project details: %w", err)
	}

	return *resp.Name, nil
}
