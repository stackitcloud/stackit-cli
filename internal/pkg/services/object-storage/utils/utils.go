package utils

import (
	"context"
	"fmt"
	"net/http"

	"github.com/stackitcloud/stackit-sdk-go/core/oapierror"
	"github.com/stackitcloud/stackit-sdk-go/services/objectstorage"
)

type ObjectStorageClient interface {
	GetServiceStatusExecute(ctx context.Context, projectId, region string) (*objectstorage.ProjectStatus, error)
	ListCredentialsGroupsExecute(ctx context.Context, projectId, region string) (*objectstorage.ListCredentialsGroupsResponse, error)
	ListAccessKeys(ctx context.Context, projectId, region string) objectstorage.ApiListAccessKeysRequest
}

func ProjectEnabled(ctx context.Context, apiClient ObjectStorageClient, projectId, region string) (bool, error) {
	_, err := apiClient.GetServiceStatusExecute(ctx, projectId, region)
	if err != nil {
		oapiErr, ok := err.(*oapierror.GenericOpenAPIError) //nolint:errorlint //complaining that error.As should be used to catch wrapped errors, but this error should not be wrapped
		if !ok {
			return false, err
		}
		if oapiErr.StatusCode == http.StatusNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func GetCredentialsGroupName(ctx context.Context, apiClient ObjectStorageClient, projectId, credentialsGroupId, region string) (string, error) {
	resp, err := apiClient.ListCredentialsGroupsExecute(ctx, projectId, region)
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

func GetCredentialsName(ctx context.Context, apiClient ObjectStorageClient, projectId, credentialsGroupId, keyId, region string) (string, error) {
	req := apiClient.ListAccessKeys(ctx, projectId, region)
	req = req.CredentialsGroup(credentialsGroupId)
	resp, err := req.Execute()

	if err != nil {
		return "", fmt.Errorf("list Object Storage credentials: %w", err)
	}

	credentials := resp.AccessKeys
	if credentials == nil {
		return "", fmt.Errorf("nil Object Storage credentials list")
	}

	for _, credential := range *credentials {
		if credential.KeyId != nil && *credential.KeyId == keyId && credential.DisplayName != nil && *credential.DisplayName != "" {
			return *credential.DisplayName, nil
		}
	}

	return "", fmt.Errorf("could not find Object Storage credentials name")
}
