package utils

import (
	"context"
	"fmt"
	"testing"

	logs "github.com/stackitcloud/stackit-sdk-go/services/logs/v1api"

	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/google/uuid"
)

var (
	testProjectId     = uuid.NewString()
	testInstanceId    = uuid.NewString()
	testAccessTokenId = uuid.NewString()
)

const (
	testInstanceName = "instance"
	testRegion       = "eu01"
)

type mockSettings struct {
	getInstanceFails    bool
	getInstanceResp     *logs.LogsInstance
	getAccessTokenFails bool
	getAccessTokenResp  *logs.AccessToken
}

func newAPIMock(s mockSettings) logs.DefaultAPI {
	return &logs.DefaultAPIServiceMock{
		GetLogsInstanceExecuteMock: utils.Ptr(func(_ logs.ApiGetLogsInstanceRequest) (*logs.LogsInstance, error) {
			if s.getInstanceFails {
				return nil, fmt.Errorf("could not get instance")
			}
			return s.getInstanceResp, nil
		}),
		GetAccessTokenExecuteMock: utils.Ptr(func(_ logs.ApiGetAccessTokenRequest) (*logs.AccessToken, error) {
			if s.getAccessTokenFails {
				return nil, fmt.Errorf("could not get access token")
			}
			return s.getAccessTokenResp, nil
		}),
	}
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
		{
			description:      "response is nil",
			getInstanceFails: false,
			getInstanceResp:  nil,
			isValid:          false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			client := mockSettings{
				getInstanceFails: tt.getInstanceFails,
				getInstanceResp:  tt.getInstanceResp,
			}

			output, err := GetInstanceName(context.Background(), newAPIMock(client), testProjectId, testRegion, testInstanceId)

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

func TestGetAccessTokenName(t *testing.T) {
	tests := []struct {
		description         string
		getAccessTokenFails bool
		getAccessTokenResp  *logs.AccessToken
		isValid             bool
		expectedOutput      string
	}{
		{
			description: "base",
			getAccessTokenResp: &logs.AccessToken{
				DisplayName: testInstanceName,
			},
			isValid:        true,
			expectedOutput: testInstanceName,
		},
		{
			description:         "get instance fails",
			getAccessTokenFails: true,
			isValid:             false,
		},
		{
			description:         "response is nil",
			getAccessTokenFails: false,
			getAccessTokenResp:  nil,
			isValid:             false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			client := mockSettings{
				getAccessTokenFails: tt.getAccessTokenFails,
				getAccessTokenResp:  tt.getAccessTokenResp,
			}

			output, err := GetAccessTokenName(context.Background(), newAPIMock(client), testProjectId, testRegion, testInstanceId, testAccessTokenId)

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
