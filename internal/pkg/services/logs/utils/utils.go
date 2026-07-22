package utils

import (
	"context"
	"errors"
	"fmt"

	logs "github.com/stackitcloud/stackit-sdk-go/services/logs/v1api"
)

var (
	ErrResponseNil = errors.New("response is nil")
)

func GetInstanceName(ctx context.Context, apiClient logs.DefaultAPI, projectId, regionId, instanceId string) (string, error) {
	resp, err := apiClient.GetLogsInstance(ctx, projectId, regionId, instanceId).Execute()
	if err != nil {
		return "", fmt.Errorf("get Logs instance: %w", err)
	} else if resp == nil {
		return "", ErrResponseNil
	}
	return resp.DisplayName, nil
}

func GetAccessTokenName(ctx context.Context, apiClient logs.DefaultAPI, projectId, regionId, instanceId, accessTokenId string) (string, error) {
	resp, err := apiClient.GetAccessToken(ctx, projectId, regionId, instanceId, accessTokenId).Execute()
	if err != nil {
		return "", fmt.Errorf("get Logs access token: %w", err)
	} else if resp == nil {
		return "", ErrResponseNil
	}
	return resp.DisplayName, nil
}
