package utils

import (
	"context"
	"fmt"
	"testing"

	"stackit/internal/pkg/utils"

	"github.com/google/uuid"
	"github.com/stackitcloud/stackit-sdk-go/services/postgresflex"
)

var (
	testProjectId  = uuid.NewString()
	testInstanceId = uuid.NewString()
)

const (
	testInstanceName = "instance"
	testUserName     = "user"
)

type postgresFlexClientMocked struct {
	getInstanceFails bool
	getInstanceResp  *postgresflex.InstanceResponse
}

func (m *postgresFlexClientMocked) GetInstanceExecute(_ context.Context, _, _ string) (*postgresflex.InstanceResponse, error) {
	if m.getInstanceFails {
		return nil, fmt.Errorf("could not get instance")
	}
	return m.getInstanceResp, nil
}

func TestGetInstanceName(t *testing.T) {
	tests := []struct {
		description      string
		getInstanceFails bool
		getInstanceResp  *postgresflex.InstanceResponse
		isValid          bool
		expectedOutput   string
	}{
		{
			description: "base",
			getInstanceResp: &postgresflex.InstanceResponse{
				Item: &postgresflex.Instance{
					Name: utils.Ptr(testInstanceName),
				},
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
			client := &postgresFlexClientMocked{
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
