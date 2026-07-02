package utils

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	secretsmanager "github.com/stackitcloud/stackit-sdk-go/services/secretsmanager/v1api"

	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

var (
	testProjectId  = uuid.NewString()
	testInstanceId = uuid.NewString()
	testUserId     = uuid.NewString()
)

const (
	testInstanceName = "instance"
	testUserName     = "user"
	testDescription  = "sample description"
)

type apiClientMockOptions struct {
	getInstanceFails bool
	getInstanceResp  *secretsmanager.Instance
	getUserFails     bool
	getUserResp      *secretsmanager.User
}

func newAPIClient(options apiClientMockOptions) secretsmanager.APIClient {
	return secretsmanager.APIClient{
		DefaultAPI: secretsmanager.DefaultAPIServiceMock{
			GetUserExecuteMock: utils.Ptr(func(_ secretsmanager.ApiGetUserRequest) (*secretsmanager.User, error) {
				if options.getUserFails {
					return nil, fmt.Errorf("could not get user")
				}
				return options.getUserResp, nil
			}),
			GetInstanceExecuteMock: utils.Ptr(func(_ secretsmanager.ApiGetInstanceRequest) (*secretsmanager.Instance, error) {
				if options.getInstanceFails {
					return nil, fmt.Errorf("could not get instance")
				}
				return options.getInstanceResp, nil
			}),
		},
	}
}

func TestGetInstanceName(t *testing.T) {
	tests := []struct {
		description      string
		getInstanceFails bool
		getInstanceResp  *secretsmanager.Instance
		isValid          bool
		expectedOutput   string
	}{
		{
			description: "base",
			getInstanceResp: &secretsmanager.Instance{
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
			mockOptions := &apiClientMockOptions{
				getInstanceFails: tt.getInstanceFails,
				getInstanceResp:  tt.getInstanceResp,
			}

			output, err := GetInstanceName(context.Background(), newAPIClient(*mockOptions).DefaultAPI, testProjectId, testInstanceId)

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

func TestGetUserDetails(t *testing.T) {
	tests := []struct {
		description    string
		getUserFails   bool
		GetUserResp    *secretsmanager.User
		isValid        bool
		expectedOutput string
	}{
		{
			description: "base",
			GetUserResp: &secretsmanager.User{
				Username:    testUserName,
				Description: testDescription,
			},
			isValid:        true,
			expectedOutput: fmt.Sprintf("%q (%s)", testUserName, testDescription),
		},
		{
			description: "user has no description",
			GetUserResp: &secretsmanager.User{
				Username: testUserName,
			},
			isValid:        true,
			expectedOutput: fmt.Sprintf("%q", testUserName),
		},
		{
			description: "user has empty description",
			GetUserResp: &secretsmanager.User{
				Username:    testUserName,
				Description: "",
			},
			isValid:        true,
			expectedOutput: fmt.Sprintf("%q", testUserName),
		},
		{
			description: "user has empty username",
			GetUserResp: &secretsmanager.User{
				Username: "",
			},
			isValid: false,
		},
		{
			description: "user has no username",
			GetUserResp: &secretsmanager.User{
				Username: "",
			},
			isValid: false,
		},
		{
			description:  "get user fails",
			getUserFails: true,
			isValid:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			options := &apiClientMockOptions{
				getUserFails: tt.getUserFails,
				getUserResp:  tt.GetUserResp,
			}

			userLabel, err := GetUserLabel(context.Background(), newAPIClient(*options).DefaultAPI, testProjectId, testInstanceId, testUserId)

			if tt.isValid && err != nil {
				t.Errorf("failed on valid input")
			}
			if !tt.isValid && err == nil {
				t.Errorf("did not fail on invalid input")
			}
			if !tt.isValid {
				return
			}
			if userLabel != tt.expectedOutput {
				t.Errorf("expected user label to be %s, got %s", tt.expectedOutput, userLabel)
			}
		})
	}
}
