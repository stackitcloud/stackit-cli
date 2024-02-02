package utils

import (
	"context"
	"fmt"
	"testing"

	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/google/uuid"
	"github.com/stackitcloud/stackit-sdk-go/services/opensearch"
)

var (
	testProjectId     = uuid.NewString()
	testInstanceId    = uuid.NewString()
	testCredentialsId = uuid.NewString()
)

const (
	testInstanceName        = "instance"
	testCredentialsUsername = "username"
)

type openSearchClientMocked struct {
	getInstanceFails    bool
	getInstanceResp     *opensearch.Instance
	getCredentialsFails bool
	getCredentialsResp  *opensearch.CredentialsResponse
}

func (m *openSearchClientMocked) GetInstanceExecute(_ context.Context, _, _ string) (*opensearch.Instance, error) {
	if m.getInstanceFails {
		return nil, fmt.Errorf("could not get instance")
	}
	return m.getInstanceResp, nil
}

func (m *openSearchClientMocked) GetCredentialsExecute(_ context.Context, _, _, _ string) (*opensearch.CredentialsResponse, error) {
	if m.getCredentialsFails {
		return nil, fmt.Errorf("could not get user")
	}
	return m.getCredentialsResp, nil
}

func TestGetInstanceName(t *testing.T) {
	tests := []struct {
		description      string
		getInstanceFails bool
		getInstanceResp  *opensearch.Instance
		isValid          bool
		expectedOutput   string
	}{
		{
			description: "base",
			getInstanceResp: &opensearch.Instance{
				Name: utils.Ptr(testInstanceName),
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
			client := &openSearchClientMocked{
				getInstanceFails: tt.getInstanceFails,
				getInstanceResp:  tt.getInstanceResp,
			}

			output, err := GetInstanceName(context.Background(), client, testProjectId, testInstanceId)

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
		getCredentialsResp  *opensearch.CredentialsResponse
		isValid             bool
		expectedOutput      string
	}{
		{
			description: "base",
			getCredentialsResp: &opensearch.CredentialsResponse{
				Raw: &opensearch.RawCredentials{
					Credentials: &opensearch.Credentials{
						Username: utils.Ptr(testCredentialsUsername),
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
			client := &openSearchClientMocked{
				getCredentialsFails: tt.getCredentialsFails,
				getCredentialsResp:  tt.getCredentialsResp,
			}

			output, err := GetCredentialsUsername(context.Background(), client, testProjectId, testInstanceId, testCredentialsId)

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
