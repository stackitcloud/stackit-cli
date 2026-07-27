package utils

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stackitcloud/stackit-sdk-go/core/oapierror"
	serviceenablement "github.com/stackitcloud/stackit-sdk-go/services/serviceenablement/v2api"

	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

const testRegion = "eu01"

var (
	testProjectId = uuid.NewString()
)

type mockSettings struct {
	serviceDisabled       bool
	getServiceStatusFails bool
	getServiceStatusResp  *serviceenablement.ServiceStatus
}

func newServiceEnableClientMock(m mockSettings) serviceenablement.DefaultAPI {
	return serviceenablement.DefaultAPIServiceMock{
		GetServiceStatusRegionalExecuteMock: utils.Ptr(func(_ serviceenablement.ApiGetServiceStatusRegionalRequest) (*serviceenablement.ServiceStatus, error) {
			if m.getServiceStatusFails {
				return nil, fmt.Errorf("could not get service status")
			}
			if m.serviceDisabled {
				return nil, &oapierror.GenericOpenAPIError{StatusCode: 404}
			}
			return m.getServiceStatusResp, nil
		}),
	}
}

func TestProjectEnabled(t *testing.T) {
	tests := []struct {
		description    string
		mockSettings   mockSettings
		isValid        bool
		expectedOutput bool
	}{
		{
			description: "project enabled",
			mockSettings: mockSettings{
				getServiceStatusResp: &serviceenablement.ServiceStatus{State: serviceenablement.SERVICESTATUSSTATE_ENABLED.Ptr()},
			},
			isValid:        true,
			expectedOutput: true,
		},
		{
			description: "project disabled (404)",
			mockSettings: mockSettings{
				serviceDisabled: true,
			},
			isValid:        true,
			expectedOutput: false,
		},
		{
			description: "project disabled 1",
			mockSettings: mockSettings{
				getServiceStatusResp: &serviceenablement.ServiceStatus{State: serviceenablement.SERVICESTATUSSTATE_ENABLING.Ptr()},
			},
			isValid:        true,
			expectedOutput: false,
		},
		{
			description: "project disabled 2",
			mockSettings: mockSettings{
				getServiceStatusResp: &serviceenablement.ServiceStatus{State: serviceenablement.SERVICESTATUSSTATE_DISABLING.Ptr()},
			},
			isValid:        true,
			expectedOutput: false,
		},
		{
			description: "project disabled 3",
			mockSettings: mockSettings{
				getServiceStatusResp: &serviceenablement.ServiceStatus{State: serviceenablement.SERVICESTATUSSTATE_DISABLING.Ptr()},
			},
			isValid:        true,
			expectedOutput: false,
		},
		{
			description: "get clusters fails",
			mockSettings: mockSettings{
				getServiceStatusFails: true,
			},
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			client := newServiceEnableClientMock(tt.mockSettings)

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
