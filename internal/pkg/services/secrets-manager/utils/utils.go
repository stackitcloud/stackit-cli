package utils

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-sdk-go/services/secretsmanager"
)

type SecretsManagerClient interface {
	GetInstanceExecute(ctx context.Context, projectId, instanceId string) (*secretsmanager.Instance, error)
}

func GetInstanceName(ctx context.Context, apiClient SecretsManagerClient, projectId, instanceId string) (string, error) {
	resp, err := apiClient.GetInstanceExecute(ctx, projectId, instanceId)
	if err != nil {
		return "", fmt.Errorf("get Secrets Manager instance: %w", err)
	}
	return *resp.Name, nil
}
