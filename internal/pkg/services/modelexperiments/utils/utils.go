package utils

import (
	"context"
	"errors"
	"fmt"

	modelexperiments "github.com/stackitcloud/stackit-sdk-go/services/modelexperiments/v1api"
)

var ErrResponseNil = errors.New("response is nil")

func GetInstanceName(ctx context.Context, apiClient modelexperiments.DefaultAPI, projectID, regionID, instanceID string) (string, error) {
	resp, err := apiClient.GetInstance(ctx, projectID, regionID, instanceID).Execute()
	if err != nil {
		return "", fmt.Errorf("get AI Model Experiments instance: %w", err)
	}
	if resp == nil {
		return "", ErrResponseNil
	}
	if resp.Instance.Name == "" {
		return "", ErrResponseNil
	}
	return resp.Instance.Name, nil
}

func GetTokenName(ctx context.Context, apiClient modelexperiments.DefaultAPI, projectID, regionID, instanceID, tokenID string) (string, error) {
	resp, err := apiClient.GetInstanceToken(ctx, projectID, regionID, tokenID, instanceID).Execute()
	if err != nil {
		return "", fmt.Errorf("get AI Model Experiments token: %w", err)
	}
	if resp == nil || resp.Token.Name == "" {
		return "", ErrResponseNil
	}
	return resp.Token.Name, nil
}
