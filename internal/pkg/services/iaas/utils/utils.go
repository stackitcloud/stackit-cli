package utils

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-sdk-go/services/iaas"
)

type IaaSClient interface {
	GetNetworkAreaExecute(ctx context.Context, organizationId, areaId string) (*iaas.NetworkArea, error)
}

func GetNetworkAreaName(ctx context.Context, apiClient IaaSClient, organizationId, areaId string) (string, error) {
	resp, err := apiClient.GetNetworkAreaExecute(ctx, organizationId, areaId)
	if err != nil {
		return "", fmt.Errorf("get network area: %w", err)
	}
	return *resp.Name, nil
}

func ListProjectsAttached(ctx context.Context, apiClient iaas.APIClient, organizationId, areaId string) ([]string, error) {
	resp, err := apiClient.
}