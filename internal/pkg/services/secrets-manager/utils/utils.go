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

func GetUserLabel(ctx context.Context, apiClient SecretsManagerClient, projectId, instanceId, userId string) (string, error) {
	resp, err := apiClient.GetUserExecute(ctx, projectId, instanceId, userId)
	if err != nil {
		return "", fmt.Errorf("get Secrets Manager user: %w", err)
	}

	var userLabel string
	if resp.Username == nil || *resp.Username == "" {
		// Should never happen, username is auto-generated
		return "", fmt.Errorf("username is empty")
	}

	if resp.Description == nil || *resp.Description == "" {
		userLabel = fmt.Sprintf("%q", *resp.Username)
	} else {
		userLabel = fmt.Sprintf("%q (%s)", *resp.Username, *resp.Description)
	}
	return userLabel, nil
}
