package utils

import (
	"context"
	"errors"
	"fmt"

	"github.com/stackitcloud/stackit-sdk-go/services/logs"
)

var (
	ErrResponseNil = errors.New("response is nil")
	ErrNameNil     = errors.New("display name is nil")
)

type LogsClient interface {
	GetLogsInstanceExecute(ctx context.Context, projectId, regionId, instanceId string) (*logs.LogsInstance, error)
	GetAccessTokenExecute(ctx context.Context, projectId string, regionId string, instanceId string, tId string) (*logs.AccessToken, error)
}

func GetInstanceName(ctx context.Context, apiClient LogsClient, projectId, regionId, instanceId string) (string, error) {
	resp, err := apiClient.GetLogsInstanceExecute(ctx, projectId, regionId, instanceId)
	if err != nil {
		return "", fmt.Errorf("get Logs instance: %w", err)
	} else if resp == nil {
		return "", ErrResponseNil
	} else if resp.DisplayName == nil {
		return "", ErrNameNil
	}
	return *resp.DisplayName, nil
}

func GetAccessTokenName(ctx context.Context, apiClient LogsClient, projectId, regionId, instanceId, accessTokenId string) (string, error) {
	resp, err := apiClient.GetAccessTokenExecute(ctx, projectId, regionId, instanceId, accessTokenId)
	if err != nil {
		return "", fmt.Errorf("get Logs access token: %w", err)
	} else if resp == nil {
		return "", ErrResponseNil
	} else if resp.DisplayName == nil {
		return "", ErrNameNil
	}
	return *resp.DisplayName, nil
}
