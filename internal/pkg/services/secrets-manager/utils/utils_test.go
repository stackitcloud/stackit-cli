package utils

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/secretsmanager"
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

type secretsManagerClientMocked struct {
	getInstanceFails bool
	getInstanceResp  *secretsmanager.Instance
	getUserFails     bool
	getUserResp      *secretsmanager.User
}

func (s *secretsManagerClientMocked) GetInstanceExecute(_ context.Context, _, _ string) (*secretsmanager.Instance, error) {
	if s.getInstanceFails {
		return nil, fmt.Errorf("could not get instance")
	}
	return s.getInstanceResp, nil
}

func (s *secretsManagerClientMocked) GetUserExecute(_ context.Context, _, _, _ string) (*secretsmanager.User, error) {
	if s.getUserFails {
		return nil, fmt.Errorf("could not get user")
	}
	return s.getUserResp, nil
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
			client := &secretsManagerClientMocked{
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

func TestGetUserDetails(t *testing.T) {
	tests := []struct {
		description         string
		getUserFails        bool
		GetUserResp         *secretsmanager.User
		isValid             bool
		expectedUserName    string
		expectedDescription string
	}{
		{
			description: "base",
			GetUserResp: &secretsmanager.User{
				Username:    utils.Ptr(testUserName),
				Description: utils.Ptr(testDescription),
			},
			isValid:             true,
			expectedUserName:    testUserName,
			expectedDescription: testDescription,
		},
		{
			description:  "get user fails",
			getUserFails: true,
			isValid:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			client := &secretsManagerClientMocked{
				getUserFails: tt.getUserFails,
				getUserResp:  tt.GetUserResp,
			}

			username, description, err := GetUserDetails(context.Background(), client, testProjectId, testInstanceId, testUserId)

			if tt.isValid && err != nil {
				t.Errorf("failed on valid input")
			}
			if !tt.isValid && err == nil {
				t.Errorf("did not fail on invalid input")
			}
			if !tt.isValid {
				return
			}
			if username != tt.expectedUserName {
				t.Errorf("expected username to be %s, got %s", tt.expectedUserName, username)
			}
			if description != tt.expectedDescription {
				t.Errorf("expected description to be %s, got %s", tt.expectedDescription, description)
			}
		})
	}
}
