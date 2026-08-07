package utils

import (
	"context"
	"fmt"

	edge "github.com/stackitcloud/stackit-sdk-go/services/edge/v1beta1api"
)

func GetInstanceName(ctx context.Context, apiClient edge.DefaultAPI, projectId, region, instanceId string) (string, error) {
	resp, err := apiClient.GetInstance(ctx, projectId, region, instanceId).Execute()
	if err != nil {
		return "", fmt.Errorf("get Edge Cloud instance: %w", err)
	}
	return resp.DisplayName, nil
}
