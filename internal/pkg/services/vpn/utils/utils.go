package utils

import (
	"context"
	"fmt"

	vpn "github.com/stackitcloud/stackit-sdk-go/services/vpn/v1api"
)

func GetGatewayName(ctx context.Context, client vpn.DefaultAPI, projectId, region, gatewayId string) (string, error) {
	resp, err := client.GetGateway(ctx, projectId, region, gatewayId).Execute()
	if err != nil {
		return "", fmt.Errorf("get gateway: %w", err)
	}
	if resp != nil {
		return resp.DisplayName, nil
	}
	return "", nil
}
