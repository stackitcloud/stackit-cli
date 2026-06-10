package utils

import (
	"context"
	"fmt"

	git "github.com/stackitcloud/stackit-sdk-go/services/git/v1betaapi"
)

func GetInstanceName(ctx context.Context, apiClient git.DefaultAPI, projectId, instanceId string) (string, error) {
	resp, err := apiClient.GetInstance(ctx, projectId, instanceId).Execute()
	if err != nil {
		return "", fmt.Errorf("get instance: %w", err)
	}
	return resp.Name, nil
}
