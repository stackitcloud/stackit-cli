package utils

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-sdk-go/services/sfs"
)

type SfsClient interface {
	GetShareExportPolicyExecute(ctx context.Context, projectId string, region string, policyId string) (*sfs.GetShareExportPolicyResponse, error)
	GetShareExecute(ctx context.Context, projectId, region, resourcePoolId, shareId string) (*sfs.GetShareResponse, error)
	GetResourcePoolExecute(ctx context.Context, projectId, region, resourcePoolId string) (*sfs.GetResourcePoolResponse, error)
}

func GetShareName(ctx context.Context, client SfsClient, projectId, region, resourcePoolId, shareId string) (string, error) {
	resp, err := client.GetShareExecute(ctx, projectId, region, resourcePoolId, shareId)
	if err != nil {
		return "", fmt.Errorf("get share: %w", err)
	}
	if resp != nil && resp.Share != nil && resp.Share.Name != nil {
		return *resp.Share.Name, nil
	}
	return "", nil
}

func GetExportPolicyName(ctx context.Context, apiClient SfsClient, projectId, region, policyId string) (string, error) {
	resp, err := apiClient.GetShareExportPolicyExecute(ctx, projectId, region, policyId)
	if err != nil {
		return "", fmt.Errorf("get share export policy: %w", err)
	}
	if resp != nil && resp.ShareExportPolicy != nil && resp.ShareExportPolicy.Name != nil {
		return *resp.ShareExportPolicy.Name, nil
	}
	return "", nil
}

func GetResourcePoolName(ctx context.Context, client SfsClient, projectId, region, resourcePoolId string) (string, error) {
	resp, err := client.GetResourcePoolExecute(ctx, projectId, region, resourcePoolId)
	if err != nil {
		return "", fmt.Errorf("get resource pool: %w", err)
	}
	if resp != nil && resp.ResourcePool != nil && resp.ResourcePool.Name != nil {
		return *resp.ResourcePool.Name, nil
	}
	return "", nil
}
