package utils

import (
	"context"
	"fmt"
	"net/http"

	"github.com/stackitcloud/stackit-sdk-go/core/oapierror"
	objectstorage "github.com/stackitcloud/stackit-sdk-go/services/objectstorage/v2api"
)

func ProjectEnabled(ctx context.Context, apiClient objectstorage.DefaultAPI, projectId, region string) (bool, error) {
	_, err := apiClient.GetServiceStatus(ctx, projectId, region).Execute()
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

func GetCredentialsGroupName(ctx context.Context, apiClient objectstorage.DefaultAPI, projectId, credentialsGroupId, region string) (string, error) {
	resp, err := apiClient.ListCredentialsGroups(ctx, projectId, region).Execute()
	if err != nil {
		return "", fmt.Errorf("list Object Storage credentials groups: %w", err)
	}

	credentialsGroups := resp.CredentialsGroups
	if credentialsGroups == nil {
		return "", fmt.Errorf("nil Object Storage credentials group list: %w", err)
	}

	for _, group := range credentialsGroups {
		if group.CredentialsGroupId == credentialsGroupId && group.DisplayName != "" {
			return group.DisplayName, nil
		}
	}

	return "", fmt.Errorf("could not find Object Storage credentials group name")
}

func GetCredentialsName(ctx context.Context, apiClient objectstorage.DefaultAPI, projectId, credentialsGroupId, keyId, region string) (string, error) {
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

	for _, credential := range credentials {
		if credential.KeyId == keyId && credential.DisplayName != "" {
			return credential.DisplayName, nil
		}
	}

	return "", fmt.Errorf("could not find Object Storage credentials name")
}
