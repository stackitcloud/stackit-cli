package utils

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/google/uuid"
	"github.com/stackitcloud/stackit-sdk-go/services/argus"
)

var (
	testProjectId  = uuid.NewString()
	testInstanceId = uuid.NewString()
	testPlanId     = uuid.NewString()
)

const (
	testInstanceName = "instance"
	testPlanName     = "Plan-Name-01"
)

var testPlansResponse = argus.PlansResponse{
	Plans: &[]argus.Plan{
		{
			Id:   utils.Ptr(testPlanId),
			Name: utils.Ptr(testPlanName),
		},
	},
}

type argusClientMocked struct {
	getInstanceFails bool
	getInstanceResp  *argus.Instance
}

func (m *argusClientMocked) GetInstanceExecute(_ context.Context, _, _ string) (*argus.Instance, error) {
	if m.getInstanceFails {
		return nil, fmt.Errorf("could not get instance")
	}
	return m.getInstanceResp, nil
}
func TestGetInstanceName(t *testing.T) {
	tests := []struct {
		description      string
		getInstanceFails bool
		getInstanceResp  *argus.Instance
		isValid          bool
		expectedOutput   string
	}{
		{
			description: "base",
			getInstanceResp: &argus.Instance{
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
			client := &argusClientMocked{
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

func TestLoadPlanId(t *testing.T) {
	tests := []struct {
		description    string
		planName       string
		isValid        bool
		expectedOutput string
	}{
		{
			description:    "base case",
			planName:       testPlanName,
			expectedOutput: testPlanId,
			isValid:        true,
		},

		{
			description:    "different casing",
			planName:       strings.ToLower(testPlanName),
			expectedOutput: testPlanId,
			isValid:        true,
		},
		{
			description: "empty plan name",
			planName:    "",
			isValid:     false,
		},
		{
			description: "unexisting plan name",
			planName:    "another plan name",
			isValid:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			output, err := LoadPlanId(tt.planName, &testPlansResponse)

			if tt.isValid && err != nil {
				t.Errorf("failed on valid input")
			}
			if !tt.isValid && err == nil {
				t.Errorf("did not fail on invalid input")
			}
			if !tt.isValid {
				return
			}
			if *output != tt.expectedOutput {
				t.Errorf("expected output to be %s, got %s", tt.expectedOutput, *output)
			}
		})
	}
}
