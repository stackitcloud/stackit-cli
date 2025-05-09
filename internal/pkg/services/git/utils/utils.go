package utils

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-sdk-go/services/git"
)

type GitClient interface {
	GetInstanceExecute(ctx context.Context, projectId string, instanceId string) (*git.Instance, error)
}

func GetInstanceName(ctx context.Context, apiClient GitClient, projectId, instanceId string) (string, error) {
	resp, err := apiClient.GetInstanceExecute(ctx, projectId, instanceId)
	if err != nil {
		return "", fmt.Errorf("get instance: %w", err)
	}
	if resp.Name == nil {
		return "", nil
	}
	return *resp.Name, nil
}
