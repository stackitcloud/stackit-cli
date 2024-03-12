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
	getInstanceResp  *argus.GetInstanceResponse
}

func (m *argusClientMocked) GetInstanceExecute(_ context.Context, _, _ string) (*argus.GetInstanceResponse, error) {
	if m.getInstanceFails {
		return nil, fmt.Errorf("could not get instance")
	}
	return m.getInstanceResp, nil
}
func TestGetInstanceName(t *testing.T) {
	tests := []struct {
		description      string
		getInstanceFails bool
		getInstanceResp  *argus.GetInstanceResponse
		isValid          bool
		expectedOutput   string
	}{
		{
			description: "base",
			getInstanceResp: &argus.GetInstanceResponse{
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
		plansResponse  *argus.PlansResponse
		isValid        bool
		expectedOutput string
	}{
		{
			description:    "base case",
			planName:       testPlanName,
			plansResponse:  utils.Ptr(testPlansResponse),
			expectedOutput: testPlanId,
			isValid:        true,
		},
		{
			description:    "different casing",
			planName:       strings.ToLower(testPlanName),
			plansResponse:  utils.Ptr(testPlansResponse),
			expectedOutput: testPlanId,
			isValid:        true,
		},
		{
			description:   "empty plan name",
			planName:      "",
			plansResponse: utils.Ptr(testPlansResponse),
			isValid:       false,
		},
		{
			description:   "unexisting plan name",
			planName:      "another plan name",
			plansResponse: utils.Ptr(testPlansResponse),
			isValid:       false,
		},
		{
			description: "unable to fetch plans",
			isValid:     false,
		},
		{
			description: "no available plans",
			planName:    testPlanName,
			plansResponse: &argus.PlansResponse{
				Plans: &[]argus.Plan{},
			},
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			output, err := LoadPlanId(tt.planName, tt.plansResponse)

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

func TestValidatePlanId(t *testing.T) {
	tests := []struct {
		description   string
		planId        string
		plansResponse *argus.PlansResponse
		isValid       bool
	}{
		{
			description:   "base case",
			planId:        testPlanId,
			plansResponse: utils.Ptr(testPlansResponse),
			isValid:       true,
		},
		{
			description:   "different casing",
			planId:        strings.ToLower(testPlanId),
			plansResponse: utils.Ptr(testPlansResponse),
			isValid:       true,
		},
		{
			description:   "empty plan id",
			planId:        "",
			plansResponse: utils.Ptr(testPlansResponse),
			isValid:       false,
		},
		{
			description:   "unexisting plan id",
			planId:        uuid.NewString(),
			plansResponse: utils.Ptr(testPlansResponse),
			isValid:       false,
		},
		{
			description: "unable to fetch plans",
			isValid:     false,
		},
		{
			description: "no available plans",
			planId:      testPlanId,
			plansResponse: &argus.PlansResponse{
				Plans: &[]argus.Plan{},
			},
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			err := ValidatePlanId(tt.planId, tt.plansResponse)

			if tt.isValid && err != nil {
				t.Errorf("failed on valid input")
			}
			if !tt.isValid && err == nil {
				t.Errorf("did not fail on invalid input")
			}
			if !tt.isValid {
				return
			}
		})
	}
}
