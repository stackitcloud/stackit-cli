package utils

import (
	"context"
	"fmt"

	sdkSKE "github.com/stackitcloud/stackit-sdk-go/services/ske"
)

type SKEClient interface {
	GetClustersExecute(ctx context.Context, projectId string) (*sdkSKE.ClustersResponse, error)
}

func ClusterExists(ctx context.Context, apiClient SKEClient, projectId, clusterName string) (bool, error) {
	clusters, err := apiClient.GetClustersExecute(ctx, projectId)
	if err != nil {
		return false, fmt.Errorf("get SKE cluster: %w", err)
	}
	for _, cl := range *clusters.Items {
		if cl.Name != nil && *cl.Name == clusterName {
			return true, nil
		}
	}
	return false, nil
}
