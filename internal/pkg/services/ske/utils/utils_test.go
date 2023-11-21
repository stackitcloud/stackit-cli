package utils

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/ske"
)

var (
	testProjectId = uuid.NewString()
)

const (
	testClusterName = "test-cluster"
)

type skeClientMocked struct {
	getClustersFails bool
	getClustersResp  *ske.ClustersResponse
}

func (m *skeClientMocked) GetClustersExecute(_ context.Context, _ string) (*ske.ClustersResponse, error) {
	if m.getClustersFails {
		return nil, fmt.Errorf("could not get clusters")
	}
	return m.getClustersResp, nil
}

func TestClusterExists(t *testing.T) {
	tests := []struct {
		description      string
		getClustersFails bool
		getClustersResp  *ske.ClustersResponse
		isValid          bool
		expectedExists   bool
	}{
		{
			description:     "cluster exists",
			getClustersResp: &ske.ClustersResponse{Items: &[]ske.ClusterResponse{{Name: utils.Ptr(testClusterName)}}},
			isValid:         true,
			expectedExists:  true,
		},
		{
			description:     "cluster exists 2",
			getClustersResp: &ske.ClustersResponse{Items: &[]ske.ClusterResponse{{Name: utils.Ptr("some-cluster")}, {Name: utils.Ptr("some-other-cluster")}, {Name: utils.Ptr(testClusterName)}}},
			isValid:         true,
			expectedExists:  true,
		},
		{
			description:     "cluster does not exist",
			getClustersResp: &ske.ClustersResponse{Items: &[]ske.ClusterResponse{{Name: utils.Ptr("some-cluster")}, {Name: utils.Ptr("some-other-cluster")}}},
			isValid:         true,
			expectedExists:  false,
		},
		{
			description:      "get clusters fails",
			getClustersFails: true,
			isValid:          false,
		},
	}

	for _, tt := range tests {
		client := &skeClientMocked{
			getClustersFails: tt.getClustersFails,
			getClustersResp:  tt.getClustersResp,
		}

		exists, err := ClusterExists(context.Background(), client, testProjectId, testClusterName)

		if tt.isValid && err != nil {
			t.Errorf("failed on valid input")
		}
		if !tt.isValid && err == nil {
			t.Errorf("did not fail on invalid input")
		}
		if !tt.isValid {
			return
		}
		if exists != tt.expectedExists {
			t.Errorf("expected exists to be %t, got %t", tt.expectedExists, exists)
		}
	}
}
