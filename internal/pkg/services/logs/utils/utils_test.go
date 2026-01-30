package utils

import (
	"context"
	"fmt"
	"testing"

	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/logs"

	"github.com/google/uuid"
)

var (
	testProjectId  = uuid.NewString()
	testInstanceId = uuid.NewString()
)

const (
	testInstanceName = "instance"
	testRegion       = "eu01"
)

type logsClientMocked struct {
	getInstanceFails   bool
	getInstanceResp    *logs.LogsInstance
	getAccessTokenResp *logs.AccessToken
}

func (m *logsClientMocked) GetLogsInstanceExecute(_ context.Context, _, _, _ string) (*logs.LogsInstance, error) {
	if m.getInstanceFails {
		return nil, fmt.Errorf("could not get instance")
	}
	return m.getInstanceResp, nil
}

func (m *logsClientMocked) GetAccessTokenExecute(_ context.Context, _, _, _, _ string) (*logs.AccessToken, error) {
	if m.getInstanceFails {
		return nil, fmt.Errorf("could not get instance")
	}
	return m.getAccessTokenResp, nil
}

func TestGetInstanceName(t *testing.T) {
	tests := []struct {
		description      string
		getInstanceFails bool
		getInstanceResp  *logs.LogsInstance
		isValid          bool
		expectedOutput   string
	}{
		{
			description: "base",
			getInstanceResp: &logs.LogsInstance{
				DisplayName: utils.Ptr(testInstanceName),
			},
			isValid:        true,
			expectedOutput: testInstanceName,
		},
		{
			description:      "get instance fails",
			getInstanceFails: true,
			isValid:          false,
		},
		{
			description:      "response is nil",
			getInstanceFails: false,
			getInstanceResp:  nil,
			isValid:          false,
		},
		{
			description:      "name in response is nil",
			getInstanceFails: false,
			getInstanceResp: &logs.LogsInstance{
				DisplayName: nil,
			},
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			client := &logsClientMocked{
				getInstanceFails: tt.getInstanceFails,
				getInstanceResp:  tt.getInstanceResp,
			}

			output, err := GetInstanceName(context.Background(), client, testProjectId, testRegion, testInstanceId)

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
