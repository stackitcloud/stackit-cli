package utils

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-sdk-go/services/objectstorage"
)

type ObjectStorageClient interface {
	ListCredentialsGroupsExecute(ctx context.Context, projectId string) (*objectstorage.ListCredentialsGroupsResponse, error)
	ListAccessKeys(ctx context.Context, projectId string) objectstorage.ApiListAccessKeysRequest
}

func GetCredentialsGroupName(ctx context.Context, apiClient ObjectStorageClient, projectId, credentialsGroupId string) (string, error) {
	resp, err := apiClient.ListCredentialsGroupsExecute(ctx, projectId)
	if err != nil {
		return "", fmt.Errorf("list Object Storage credentials groups: %w", err)
	}

	credentialsGroups := resp.CredentialsGroups
	if credentialsGroups == nil {
		return "", fmt.Errorf("nil Object Storage credentials group list: %w", err)
	}

	for _, group := range *credentialsGroups {
		if group.CredentialsGroupId != nil && *group.CredentialsGroupId == credentialsGroupId && group.DisplayName != nil && *group.DisplayName != "" {
			return *group.DisplayName, nil
		}
	}

	return "", fmt.Errorf("could not find Object Storage credentials group name")
}

func GetCredentialsName(ctx context.Context, apiClient ObjectStorageClient, projectId, credentialsGroupId, keyId string) (string, error) {
	req := apiClient.ListAccessKeys(ctx, projectId)
	req = req.CredentialsGroup(credentialsGroupId)
	resp, err := req.Execute()

	if err != nil {
		return "", fmt.Errorf("list Object Storage credentials: %w", err)
	}

	credentials := resp.AccessKeys
	if credentials == nil {
		return "", fmt.Errorf("nil Object Storage credentials list: %w", err)
	}

	for _, credential := range *credentials {
		if credential.KeyId != nil && *credential.KeyId == keyId && credential.DisplayName != nil && *credential.DisplayName != "" {
			return *credential.DisplayName, nil
		}
	}

	return "", fmt.Errorf("could not find Object Storage credential name")
}
