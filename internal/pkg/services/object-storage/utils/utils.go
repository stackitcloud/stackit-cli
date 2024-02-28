package utils

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-sdk-go/services/objectstorage"
)

type ObjectStorageClient interface {
	ListCredentialsGroupsExecute(ctx context.Context, projectId string) (*objectstorage.ListCredentialsGroupsResponse, error)
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

	var name string
	for _, group := range *credentialsGroups {
		if group.CredentialsGroupId != nil && *group.CredentialsGroupId == credentialsGroupId && group.DisplayName != nil {
			name = *group.DisplayName
		}
	}

	if name == "" {
		return "", fmt.Errorf("could not find Object Storage credentials group name")
	}
	return name, nil
}
