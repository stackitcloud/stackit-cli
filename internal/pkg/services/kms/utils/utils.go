package utils

import (
	"context"
	"fmt"
	"time"

	kms "github.com/stackitcloud/stackit-sdk-go/services/kms/v1api"
)

func GetKeyName(ctx context.Context, apiClient kms.DefaultAPI, projectId, region, keyRingId, keyId string) (string, error) {
	resp, err := apiClient.GetKey(ctx, projectId, region, keyRingId, keyId).Execute()
	if err != nil {
		return "", fmt.Errorf("get KMS Key: %w", err)
	}

	if resp == nil {
		return "", fmt.Errorf("response is nil / empty")
	}

	return resp.DisplayName, nil
}

func GetKeyDeletionDate(ctx context.Context, apiClient kms.DefaultAPI, projectId, region, keyRingId, keyId string) (time.Time, error) {
	resp, err := apiClient.GetKey(ctx, projectId, region, keyRingId, keyId).Execute()
	if err != nil {
		return time.Now(), fmt.Errorf("get KMS Key: %w", err)
	}

	if resp == nil || resp.DeletionDate == nil {
		return time.Time{}, fmt.Errorf("response is nil / empty")
	}

	return *resp.DeletionDate, nil
}

func GetKeyRingName(ctx context.Context, apiClient kms.DefaultAPI, projectId, id, region string) (string, error) {
	resp, err := apiClient.GetKeyRing(ctx, projectId, region, id).Execute()
	if err != nil {
		return "", fmt.Errorf("get KMS key ring: %w", err)
	}

	if resp == nil {
		return "", fmt.Errorf("response is nil / empty")
	}

	return resp.DisplayName, nil
}

func GetWrappingKeyName(ctx context.Context, apiClient kms.DefaultAPI, projectId, region, keyRingId, wrappingKeyId string) (string, error) {
	resp, err := apiClient.GetWrappingKey(ctx, projectId, region, keyRingId, wrappingKeyId).Execute()
	if err != nil {
		return "", fmt.Errorf("get KMS Wrapping Key: %w", err)
	}

	if resp == nil {
		return "", fmt.Errorf("response is nil / empty")
	}

	return resp.DisplayName, nil
}
