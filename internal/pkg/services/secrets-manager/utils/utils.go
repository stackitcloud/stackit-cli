package utils

import (
	"context"
	"fmt"

	secretsmanager "github.com/stackitcloud/stackit-sdk-go/services/secretsmanager/v1api"
)

type SecretsManagerClient interface {
	GetInstance(ctx context.Context, projectId, instanceId string) secretsmanager.ApiGetInstanceRequest
	GetUser(ctx context.Context, projectId string, instanceId string, userId string) secretsmanager.ApiGetUserRequest
}

func GetInstanceName(ctx context.Context, apiClient SecretsManagerClient, projectId, instanceId string) (string, error) {
	resp, err := apiClient.GetInstance(ctx, projectId, instanceId).Execute()
	if err != nil {
		return "", fmt.Errorf("get Secrets Manager instance: %w", err)
	}
	return resp.Name, nil
}

func GetUserLabel(ctx context.Context, apiClient SecretsManagerClient, projectId, instanceId, userId string) (string, error) {
	resp, err := apiClient.GetUser(ctx, projectId, instanceId, userId).Execute()
	if err != nil {
		return "", fmt.Errorf("get Secrets Manager user: %w", err)
	}

	if resp.Username == "" {
		// Should never happen, username is auto-generated
		return "", fmt.Errorf("username is empty")
	}

	var userLabel string
	if resp.Description == "" {
		userLabel = fmt.Sprintf("%q", resp.Username)
	} else {
		userLabel = fmt.Sprintf("%q (%s)", resp.Username, resp.Description)
	}
	return userLabel, nil
}
