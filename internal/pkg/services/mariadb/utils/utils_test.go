package utils

import (
	"context"
	"fmt"
	"testing"

	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/google/uuid"
	mariadb "github.com/stackitcloud/stackit-sdk-go/services/mariadb/v1api"
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

type mockSettings struct {
	getInstanceFails    bool
	getInstanceResp     *mariadb.Instance
	getCredentialsFails bool
	getCredentialsResp  *mariadb.CredentialsResponse
}

func newAPIMock(m mockSettings) mariadb.DefaultAPI {
	return &mariadb.DefaultAPIServiceMock{
		GetInstanceExecuteMock: utils.Ptr(func(_ mariadb.ApiGetInstanceRequest) (*mariadb.Instance, error) {
			if m.getInstanceFails {
				return nil, fmt.Errorf("could not get instance")
			}
			return m.getInstanceResp, nil
		}),
		GetCredentialsExecuteMock: utils.Ptr(func(_ mariadb.ApiGetCredentialsRequest) (*mariadb.CredentialsResponse, error) {
			if m.getCredentialsFails {
				return nil, fmt.Errorf("could not get user")
			}
			return m.getCredentialsResp, nil
		}),
	}
}

func TestGetInstanceName(t *testing.T) {
	tests := []struct {
		description      string
		getInstanceFails bool
		getInstanceResp  *mariadb.Instance
		isValid          bool
		expectedOutput   string
	}{
		{
			description: "base",
			getInstanceResp: &mariadb.Instance{
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
			settings := mockSettings{
				getInstanceFails: tt.getInstanceFails,
				getInstanceResp:  tt.getInstanceResp,
			}

			output, err := GetInstanceName(context.Background(), newAPIMock(settings), testProjectId, testInstanceId)

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
		getCredentialsResp  *mariadb.CredentialsResponse
		isValid             bool
		expectedOutput      string
	}{
		{
			description: "base",
			getCredentialsResp: &mariadb.CredentialsResponse{
				Raw: &mariadb.RawCredentials{
					Credentials: mariadb.Credentials{
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
			settings := mockSettings{
				getCredentialsFails: tt.getCredentialsFails,
				getCredentialsResp:  tt.getCredentialsResp,
			}

			output, err := GetCredentialsUsername(context.Background(), newAPIMock(settings), testProjectId, testInstanceId, testCredentialsId)

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
