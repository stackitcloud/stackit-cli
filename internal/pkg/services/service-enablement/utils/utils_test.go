package utils

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stackitcloud/stackit-sdk-go/core/oapierror"
	"github.com/stackitcloud/stackit-sdk-go/services/serviceenablement"
)

var (
	testProjectId = uuid.NewString()
	testRegion    = "eu01"
)

type serviceEnableClientMocked struct {
	serviceDisabled       bool
	getServiceStatusFails bool
	getServiceStatusResp  *serviceenablement.ServiceStatus
}

func (m *serviceEnableClientMocked) GetServiceStatusRegionalExecute(_ context.Context, _, _, _ string) (*serviceenablement.ServiceStatus, error) {
	if m.getServiceStatusFails {
		return nil, fmt.Errorf("could not get service status")
	}
	if m.serviceDisabled {
		return nil, &oapierror.GenericOpenAPIError{StatusCode: 404}
	}
	return m.getServiceStatusResp, nil
}

func TestProjectEnabled(t *testing.T) {
	tests := []struct {
		description     string
		serviceDisabled bool
		getProjectFails bool
		getProjectResp  *serviceenablement.ServiceStatus
		isValid         bool
		expectedOutput  bool
	}{
		{
			description:    "project enabled",
			getProjectResp: &serviceenablement.ServiceStatus{State: serviceenablement.SERVICESTATUSSTATE_ENABLED.Ptr()},
			isValid:        true,
			expectedOutput: true,
		},
		{
			description:     "project disabled (404)",
			serviceDisabled: true,
			isValid:         true,
			expectedOutput:  false,
		},
		{
			description:    "project disabled 1",
			getProjectResp: &serviceenablement.ServiceStatus{State: serviceenablement.SERVICESTATUSSTATE_ENABLING.Ptr()},
			isValid:        true,
			expectedOutput: false,
		},
		{
			description:    "project disabled 2",
			getProjectResp: &serviceenablement.ServiceStatus{State: serviceenablement.SERVICESTATUSSTATE_DISABLING.Ptr()},
			isValid:        true,
			expectedOutput: false,
		},
		{
			description:    "project disabled 3",
			getProjectResp: &serviceenablement.ServiceStatus{State: serviceenablement.SERVICESTATUSSTATE_DISABLING.Ptr()},
			isValid:        true,
			expectedOutput: false,
		},
		{
			description:     "get clusters fails",
			getProjectFails: true,
			isValid:         false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			client := &serviceEnableClientMocked{
				serviceDisabled:       tt.serviceDisabled,
				getServiceStatusFails: tt.getProjectFails,
				getServiceStatusResp:  tt.getProjectResp,
			}

			output, err := ProjectEnabled(context.Background(), client, testRegion, testProjectId)

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
				t.Errorf("expected output to be %t, got %t", tt.expectedOutput, output)
			}
		})
	}
}
