package utils

import (
	"context"
	"fmt"
	"time"

	"github.com/stackitcloud/stackit-sdk-go/services/kms"
)

type KMSClient interface {
	GetKeyExecute(ctx context.Context, projectId string, regionId string, keyRingId string, keyId string) (*kms.Key, error)
	GetKeyRingExecute(ctx context.Context, projectId string, regionId string, keyRingId string) (*kms.KeyRing, error)
	GetWrappingKeyExecute(ctx context.Context, projectId string, regionId string, keyRingId string, wrappingKeyId string) (*kms.WrappingKey, error)
}

func GetKeyName(ctx context.Context, apiClient KMSClient, projectId, region, keyRingId, keyId string) (string, error) {
	resp, err := apiClient.GetKeyExecute(ctx, projectId, region, keyRingId, keyId)
	if err != nil {
		return "", fmt.Errorf("get KMS Key: %w", err)
	}

	if resp == nil || resp.DisplayName == nil {
		return "", fmt.Errorf("response is nil / empty")
	}

	return *resp.DisplayName, nil
}

func GetKeyDeletionDate(ctx context.Context, apiClient KMSClient, projectId, region, keyRingId, keyId string) (time.Time, error) {
	resp, err := apiClient.GetKeyExecute(ctx, projectId, region, keyRingId, keyId)
	if err != nil {
		return time.Now(), fmt.Errorf("get KMS Key: %w", err)
	}

	if resp == nil || resp.DeletionDate == nil {
		return time.Time{}, fmt.Errorf("response is nil / empty")
	}

	return *resp.DeletionDate, nil
}

func GetKeyRingName(ctx context.Context, apiClient KMSClient, projectId, id, region string) (string, error) {
	resp, err := apiClient.GetKeyRingExecute(ctx, projectId, region, id)
	if err != nil {
		return "", fmt.Errorf("get KMS key ring: %w", err)
	}

	if resp == nil || resp.DisplayName == nil {
		return "", fmt.Errorf("response is nil / empty")
	}

	return *resp.DisplayName, nil
}

func GetWrappingKeyName(ctx context.Context, apiClient KMSClient, projectId, region, keyRingId, wrappingKeyId string) (string, error) {
	resp, err := apiClient.GetWrappingKeyExecute(ctx, projectId, region, keyRingId, wrappingKeyId)
	if err != nil {
		return "", fmt.Errorf("get KMS Wrapping Key: %w", err)
	}

	if resp == nil || resp.DisplayName == nil {
		return "", fmt.Errorf("response is nil / empty")
	}

	return *resp.DisplayName, nil
}
