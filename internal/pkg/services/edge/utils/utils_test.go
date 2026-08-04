package utils

import (
	"context"
	"fmt"
	"testing"

	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/google/uuid"
	edge "github.com/stackitcloud/stackit-sdk-go/services/edge/v1beta1api"
)

var (
	testProjectId  = uuid.NewString()
	testInstanceId = uuid.NewString()
)

const (
	testInstanceName = "instance"
	testRegion       = "eu01"
)

type mockSettings struct {
	getInstanceFails bool
	getInstanceResp  *edge.Instance
}

func newApiMock(s *mockSettings) edge.DefaultAPI {
	return &edge.DefaultAPIServiceMock{
		GetInstanceExecuteMock: utils.Ptr(func(_ edge.ApiGetInstanceRequest) (*edge.Instance, error) {
			if s.getInstanceFails {
				return nil, fmt.Errorf("could not get instance")
			}
			return s.getInstanceResp, nil
		}),
	}
}

func TestGetInstanceName(t *testing.T) {
	tests := []struct {
		description      string
		getInstanceFails bool
		getInstanceResp  *edge.Instance
		isValid          bool
		expectedOutput   string
	}{
		{
			description: "base",
			getInstanceResp: &edge.Instance{
				DisplayName: testInstanceName,
			},
			isValid:        true,
			expectedOutput: testInstanceName,
		},
		{
			description:      "get instance fails",
			getInstanceFails: true,
			isValid:          false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			settings := &mockSettings{
				getInstanceFails: tt.getInstanceFails,
				getInstanceResp:  tt.getInstanceResp,
			}

			output, err := GetInstanceName(context.Background(), newApiMock(settings), testProjectId, testRegion, testInstanceId)

			if tt.isValid && err != nil {
				t.Errorf("failed on valid input")
			}
			if !tt.isValid && err == nil {
				t.Errorf("did not fail on invalid input")
			}
			if !tt.isValid {
				return
			}
			if output != tt.expectedOutput {
				t.Errorf("expected output to be %s, got %s", tt.expectedOutput, output)
			}
		})
	}
}
