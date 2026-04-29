package utils

import (
	"context"
	"fmt"

	sfs "github.com/stackitcloud/stackit-sdk-go/services/sfs/v1api"
)

func GetShareName(ctx context.Context, client sfs.DefaultAPI, projectId, region, resourcePoolId, shareId string) (string, error) {
	resp, err := client.GetShare(ctx, projectId, region, resourcePoolId, shareId).Execute()
	if err != nil {
		return "", fmt.Errorf("get share: %w", err)
	}
	if resp != nil && resp.Share != nil && resp.Share.Name != nil {
		return *resp.Share.Name, nil
	}
	return "", nil
}

func GetExportPolicyName(ctx context.Context, apiClient sfs.DefaultAPI, projectId, region, policyId string) (string, error) {
	resp, err := apiClient.GetShareExportPolicy(ctx, projectId, region, policyId).Execute()
	if err != nil {
		return "", fmt.Errorf("get share export policy: %w", err)
	}
	if resp != nil && resp.ShareExportPolicy != nil && resp.ShareExportPolicy.Name != nil {
		return *resp.ShareExportPolicy.Name, nil
	}
	return "", nil
}

func GetResourcePoolName(ctx context.Context, client sfs.DefaultAPI, projectId, region, resourcePoolId string) (string, error) {
	resp, err := client.GetResourcePool(ctx, projectId, region, resourcePoolId).Execute()
	if err != nil {
		return "", fmt.Errorf("get resource pool: %w", err)
	}
	if resp != nil && resp.ResourcePool != nil && resp.ResourcePool.Name != nil {
		return *resp.ResourcePool.Name, nil
	}
	return "", nil
}
