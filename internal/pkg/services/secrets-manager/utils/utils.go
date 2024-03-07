package utils

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-sdk-go/services/secretsmanager"
)

type SecretsManagerClient interface {
	GetInstanceExecute(ctx context.Context, projectId, instanceId string) (*secretsmanager.Instance, error)
	GetUserExecute(ctx context.Context, projectId string, instanceId string, userId string) (*secretsmanager.User, error)
}

func GetInstanceName(ctx context.Context, apiClient SecretsManagerClient, projectId, instanceId string) (string, error) {
	resp, err := apiClient.GetInstanceExecute(ctx, projectId, instanceId)
	if err != nil {
		return "", fmt.Errorf("get Secrets Manager instance: %w", err)
	}
	return *resp.Name, nil
}

func GetUserDetails(ctx context.Context, apiClient SecretsManagerClient, projectId, instanceId, userId string) (username, description string, err error) {
	resp, err := apiClient.GetUserExecute(ctx, projectId, instanceId, userId)
	if err != nil {
		return "", "", fmt.Errorf("get Secrets Manager user: %w", err)
	}
	return *resp.Username, *resp.Description, nil
}
