package utils

import (
	"context"
	"fmt"
	"testing"

	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/google/uuid"
	rabbitmq "github.com/stackitcloud/stackit-sdk-go/services/rabbitmq/v2api"
)

var (
	testProjectId     = uuid.NewString()
	testInstanceId    = uuid.NewString()
	testCredentialsId = uuid.NewString()
)

const (
	testInstanceName        = "instance"
	testRegion              = "eu01"
	testCredentialsUsername = "username"
)

type rabbitMQClientMockSettings struct {
	getInstanceFails    bool
	getInstanceResp     *rabbitmq.Instance
	getCredentialsFails bool
	getCredentialsResp  *rabbitmq.CredentialsResponse
}

func newApiMock(s *rabbitMQClientMockSettings) rabbitmq.DefaultAPI {
	return &rabbitmq.DefaultAPIServiceMock{
		GetInstanceExecuteMock: utils.Ptr(func(_ rabbitmq.ApiGetInstanceRequest) (*rabbitmq.Instance, error) {
			if s.getInstanceFails {
				return nil, fmt.Errorf("could not get instance")
			}
			return s.getInstanceResp, nil
		}),
		GetCredentialsExecuteMock: utils.Ptr(func(_ rabbitmq.ApiGetCredentialsRequest) (*rabbitmq.CredentialsResponse, error) {
			if s.getCredentialsFails {
				return nil, fmt.Errorf("could not get user")
			}
			return s.getCredentialsResp, nil
		}),
	}
}

func TestGetInstanceName(t *testing.T) {
	tests := []struct {
		description      string
		getInstanceFails bool
		getInstanceResp  *rabbitmq.Instance
		isValid          bool
		expectedOutput   string
	}{
		{
			description: "base",
			getInstanceResp: &rabbitmq.Instance{
				Name: testInstanceName,
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
			settings := &rabbitMQClientMockSettings{
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

func TestGetCredentialsUsername(t *testing.T) {
	tests := []struct {
		description         string
		getCredentialsFails bool
		getCredentialsResp  *rabbitmq.CredentialsResponse
		isValid             bool
		expectedOutput      string
	}{
		{
			description: "base",
			getCredentialsResp: &rabbitmq.CredentialsResponse{
				Raw: &rabbitmq.RawCredentials{
					Credentials: rabbitmq.Credentials{
						Username: testCredentialsUsername,
					},
				},
			},
			isValid:        true,
			expectedOutput: testCredentialsUsername,
		},
		{
			description:         "get credentials fails",
			getCredentialsFails: true,
			isValid:             false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			settings := &rabbitMQClientMockSettings{
				getCredentialsFails: tt.getCredentialsFails,
				getCredentialsResp:  tt.getCredentialsResp,
			}

			output, err := GetCredentialsUsername(context.Background(), newApiMock(settings), testProjectId, testRegion, testInstanceId, testCredentialsId)

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
