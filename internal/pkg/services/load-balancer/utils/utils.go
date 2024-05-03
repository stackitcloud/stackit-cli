package utils

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-sdk-go/services/loadbalancer"
)

type LoadBalancerClient interface {
	GetCredentialsExecute(ctx context.Context, projectId, credentialsRef string) (*loadbalancer.GetCredentialsResponse, error)
}

func GetCredentialsDisplayName(ctx context.Context, apiClient LoadBalancerClient, projectId, credentialsRef string) (string, error) {
	resp, err := apiClient.GetCredentialsExecute(ctx, projectId, credentialsRef)
	if err != nil {
		return "", fmt.Errorf("get Load Balancer credentials: %w", err)
	}
	return *resp.Credential.DisplayName, nil
}
