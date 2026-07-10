package utils

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	opensearch "github.com/stackitcloud/stackit-sdk-go/services/opensearch/v2api"

	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

var (
	testProjectId     = uuid.NewString()
	testRegion        = "eu01"
	testInstanceId    = uuid.NewString()
	testCredentialsId = uuid.NewString()
)

const (
	testInstanceName        = "instance"
	testCredentialsUsername = "username"
)

type mockSettings struct {
	getInstanceFails    bool
	getInstanceResp     *opensearch.Instance
	getCredentialsFails bool
	getCredentialsResp  *opensearch.CredentialsResponse
}

func newAPIClientMock(m mockSettings) opensearch.DefaultAPI {
	return opensearch.DefaultAPIServiceMock{
		GetInstanceExecuteMock: utils.Ptr(func(_ opensearch.ApiGetInstanceRequest) (*opensearch.Instance, error) {
			if m.getInstanceFails {
				return nil, fmt.Errorf("could not get instance")
			}
			return m.getInstanceResp, nil
		}),
		GetCredentialsExecuteMock: utils.Ptr(func(_ opensearch.ApiGetCredentialsRequest) (*opensearch.CredentialsResponse, error) {
			if m.getCredentialsFails {
				return nil, fmt.Errorf("could not get user")
			}
			return m.getCredentialsResp, nil
		}),
	}
}

func TestGetInstanceName(t *testing.T) {
	tests := []struct {
		description        string
		mockClientSettings mockSettings
		isValid            bool
		expectedOutput     string
	}{
		{
			description: "base",
			mockClientSettings: mockSettings{
				getInstanceResp: &opensearch.Instance{
					Name: testInstanceName,
				},
			},
			isValid:        true,
			expectedOutput: testInstanceName,
		},
		{
			description: "get instance fails",
			mockClientSettings: mockSettings{
				getInstanceFails: true,
			},
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			output, err := GetInstanceName(context.Background(), newAPIClientMock(tt.mockClientSettings), testProjectId, testRegion, testInstanceId)

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
		description        string
		mockClientSettings mockSettings
		isValid            bool
		expectedOutput     string
	}{
		{
			description: "base",
			mockClientSettings: mockSettings{
				getCredentialsResp: &opensearch.CredentialsResponse{
					Raw: &opensearch.RawCredentials{
						Credentials: opensearch.Credentials{
							Username: testCredentialsUsername,
						},
					},
				},
			},
			isValid:        true,
			expectedOutput: testCredentialsUsername,
		},
		{
			description: "get credentials fails",
			mockClientSettings: mockSettings{
				getCredentialsFails: true,
			},
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			output, err := GetCredentialsUsername(context.Background(), newAPIClientMock(tt.mockClientSettings), testProjectId, testRegion, testInstanceId, testCredentialsId)

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
