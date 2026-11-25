package utils

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	registry "github.com/stackitcloud/stackit-sdk-go/services/registry/v1api"

	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

var (
	testProjectId     = uuid.NewString()
	testArtifactoryId = uuid.NewString()
)

const (
	testRegion          = "eu01"
	testArtifactoryName = "my-artifactory"
)

type registryClientMocked struct {
	getArtifactoryFails bool
	getArtifactoryResp  *registry.Artifactory
}

func (m *registryClientMocked) newMock() registry.DefaultAPI {
	return registry.DefaultAPIServiceMock{
		GetArtifactoryExecuteMock: utils.Ptr(func(_ registry.ApiGetArtifactoryRequest) (*registry.Artifactory, error) {
			if m.getArtifactoryFails {
				return nil, fmt.Errorf("could not get artifactory")
			}
			return m.getArtifactoryResp, nil
		}),
	}
}

func TestGetArtifactoryName(t *testing.T) {
	tests := []struct {
		description         string
		getArtifactoryFails bool
		getArtifactoryResp  *registry.Artifactory
		isValid             bool
		expectedOutput      string
	}{
		{
			description: "base",
			getArtifactoryResp: &registry.Artifactory{
				Name: testArtifactoryName,
			},
			isValid:        true,
			expectedOutput: testArtifactoryName,
		},
		{
			description:         "get artifactory fails",
			getArtifactoryFails: true,
			isValid:             false,
		},
		{
			description:        "empty response (nil)",
			getArtifactoryResp: nil,
			isValid:            false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			client := &registryClientMocked{
				getArtifactoryFails: tt.getArtifactoryFails,
				getArtifactoryResp:  tt.getArtifactoryResp,
			}

			output, err := GetArtifactoryName(context.Background(), client.newMock(), testProjectId, testRegion, testArtifactoryId)

			if tt.isValid && err != nil {
				t.Errorf("failed on valid input: %v", err)
			}
			if !tt.isValid && err == nil {
				t.Errorf("did not fail on invalid input")
			}
			if tt.isValid && output != tt.expectedOutput {
				t.Errorf("expected output %q, got %q", tt.expectedOutput, output)
			}
		})
	}
}
