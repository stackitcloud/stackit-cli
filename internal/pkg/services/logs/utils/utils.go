package utils

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-sdk-go/services/logs"
)

type LogsClient interface {
	GetLogsInstanceExecute(ctx context.Context, projectId, regionId, instanceId string) (*logs.LogsInstance, error)
}

func GetInstanceName(ctx context.Context, apiClient LogsClient, projectId, regionId, instanceId string) (string, error) {
	resp, err := apiClient.GetLogsInstanceExecute(ctx, projectId, regionId, instanceId)
	if err != nil {
		return "", fmt.Errorf("get Logs instance: %w", err)
	}
	return *resp.DisplayName, nil
}
