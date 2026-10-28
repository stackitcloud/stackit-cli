package utils

import (
	"context"
	"fmt"

	scf "github.com/stackitcloud/stackit-sdk-go/services/scf/v1api"
)

func GetOrganizationName(ctx context.Context, apiClient scf.DefaultAPI, projectId, region, organizationId string) (string, error) {
	resp, err := apiClient.GetOrganization(ctx, projectId, region, organizationId).Execute()
	if err != nil {
		return "", fmt.Errorf("get STACKIT Cloud Foundry organization: %w", err)
	}
	if resp == nil {
		return "", fmt.Errorf("get STACKIT Cloud Foundry organization: empty response")
	}
	return resp.Name, nil
}
