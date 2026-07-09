package utils

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	redis "github.com/stackitcloud/stackit-sdk-go/services/redis/v2api"

	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

var (
	testProjectId     = uuid.NewString()
	testInstanceId    = uuid.NewString()
	testCredentialsId = uuid.NewString()
)

const (
	testInstanceName        = "instance"
	testCredentialsUsername = "username"
	testRegion              = "eu01"
)

type mockSettings struct {
	getInstanceFails    bool
	getInstanceResp     *redis.Instance
	getCredentialsFails bool
	getCredentialsResp  *redis.CredentialsResponse
}

func newAPIMock(settings *mockSettings) redis.DefaultAPI {
	return &redis.DefaultAPIServiceMock{
		GetInstanceExecuteMock: utils.Ptr(func(_ redis.ApiGetInstanceRequest) (*redis.Instance, error) {
			if settings.getInstanceFails {
				return nil, fmt.Errorf("could not get instance")
			}

			return settings.getInstanceResp, nil
		}),
		GetCredentialsExecuteMock: utils.Ptr(func(_ redis.ApiGetCredentialsRequest) (*redis.CredentialsResponse, error) {
			if settings.getCredentialsFails {
				return nil, fmt.Errorf("could not get user")
			}

			return settings.getCredentialsResp, nil
		}),
	}
}

func TestGetInstanceName(t *testing.T) {
	tests := []struct {
		description      string
		getInstanceFails bool
		getInstanceResp  *redis.Instance
		isValid          bool
		expectedOutput   string
	}{
		{
			description: "base",
			getInstanceResp: &redis.Instance{
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
			client := newAPIMock(&mockSettings{
				getInstanceFails: tt.getInstanceFails,
				getInstanceResp:  tt.getInstanceResp,
			})

			output, err := GetInstanceName(context.Background(), client, testProjectId, testInstanceId, testRegion)

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
		getCredentialsResp  *redis.CredentialsResponse
		isValid             bool
		expectedOutput      string
	}{
		{
			description: "base",
			getCredentialsResp: &redis.CredentialsResponse{
				Raw: &redis.RawCredentials{
					Credentials: redis.Credentials{
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
			client := newAPIMock(&mockSettings{
				getCredentialsFails: tt.getCredentialsFails,
				getCredentialsResp:  tt.getCredentialsResp,
			})

			output, err := GetCredentialsUsername(context.Background(), client, testProjectId, testInstanceId, testCredentialsId, testRegion)

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
