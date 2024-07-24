package utils

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-sdk-go/services/iaas"
)

type IaaSClient interface {
	GetNetworkAreaExecute(ctx context.Context, organizationId, areaId string) (*iaas.NetworkArea, error)
	ListNetworkAreaProjectsExecute(ctx context.Context, organizationId, areaId string) (*iaas.ProjectListResponse, error)
	GetNetworkAreaRangeExecute(ctx context.Context, organizationId, areaId, networkRangeId string) (*iaas.NetworkRange, error)
}

func GetNetworkAreaName(ctx context.Context, apiClient IaaSClient, organizationId, areaId string) (string, error) {
	resp, err := apiClient.GetNetworkAreaExecute(ctx, organizationId, areaId)
	if err != nil {
		return "", fmt.Errorf("get network area: %w", err)
	}
	return *resp.Name, nil
}

func ListAttachedProjects(ctx context.Context, apiClient IaaSClient, organizationId, areaId string) ([]string, error) {
	resp, err := apiClient.ListNetworkAreaProjectsExecute(ctx, organizationId, areaId)
	if err != nil {
		return nil, fmt.Errorf("list network area attached projects: %w", err)
	}
	return *resp.Items, nil
}

func GetNetworkRangePrefix(ctx context.Context, apiClient IaaSClient, organizationId, areaId, networkRangeId string) (string, error) {
	resp, err := apiClient.GetNetworkAreaRangeExecute(ctx, organizationId, areaId, networkRangeId)
	if err != nil {
		return "", fmt.Errorf("get network range: %w", err)
	}
	return *resp.Prefix, nil
}
