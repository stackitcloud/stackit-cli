package utils

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	vpn "github.com/stackitcloud/stackit-sdk-go/services/vpn/v1api"

	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

const (
	testGatewayName = "test-gateway-01"
	testRegion      = "eu01"
)

var (
	testProjectId = uuid.NewString()
	testGatewayId = uuid.NewString()
)

type mockSettings struct {
	getGatewayFails bool
	getGatewayResp  *vpn.GatewayResponse
}

func newAPIMock(settings *mockSettings) vpn.DefaultAPI {
	return &vpn.DefaultAPIServiceMock{
		GetGatewayExecuteMock: utils.Ptr(func(_ vpn.ApiGetGatewayRequest) (*vpn.GatewayResponse, error) {
			if settings.getGatewayFails {
				return nil, fmt.Errorf("could not get gateway details")
			}

			return settings.getGatewayResp, nil
		}),
	}
}

func TestGetGatewayName(t *testing.T) {
	tests := []struct {
		description     string
		getGatewayResp  *vpn.GatewayResponse
		getGatewayFails bool
		isValid         bool
		expectedOutput  string
	}{
		{
			description: "base",
			getGatewayResp: &vpn.GatewayResponse{
				DisplayName: testGatewayName,
			},
			isValid:        true,
			expectedOutput: testGatewayName,
		},
		{
			description:     "get gateway fails",
			getGatewayFails: true,
			isValid:         false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			client := newAPIMock(&mockSettings{
				getGatewayFails: tt.getGatewayFails,
				getGatewayResp:  tt.getGatewayResp,
			})

			output, err := GetGatewayName(context.Background(), client, testProjectId, testRegion, testGatewayId)

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
