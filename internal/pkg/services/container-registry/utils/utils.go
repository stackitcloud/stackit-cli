package utils

import (
	"context"
	"fmt"

	registry "github.com/stackitcloud/stackit-sdk-go/services/registry/v1api"
)

func GetArtifactoryName(ctx context.Context, apiClient registry.DefaultAPI, projectId, region, artifactoryId string) (string, error) {
	resp, err := apiClient.GetArtifactory(ctx, projectId, region, artifactoryId).Execute()
	if err != nil {
		return "", fmt.Errorf("get artifactory: %w", err)
	}
	if resp == nil {
		return "", fmt.Errorf("get artifactory: empty response")
	}
	return resp.Name, nil
}
