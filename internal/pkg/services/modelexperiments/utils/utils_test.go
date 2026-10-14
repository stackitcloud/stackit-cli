package utils

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"

	modelexperiments "github.com/stackitcloud/stackit-sdk-go/services/modelexperiments/v1api"

	cliutils "github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

var (
	testProjectID  = uuid.NewString()
	testRegion     = "eu01"
	testInstanceID = uuid.NewString()
	testTokenID    = uuid.NewString()
)

type mockSettings struct {
	instanceFails bool
	instanceResp  *modelexperiments.GetInstanceResponse
	tokenFails    bool
	tokenResp     *modelexperiments.GetInstanceTokenResponse
}

func newAPIMock(settings mockSettings) modelexperiments.DefaultAPI {
	return &modelexperiments.DefaultAPIServiceMock{
		GetInstanceExecuteMock: cliutils.Ptr(func(modelexperiments.ApiGetInstanceRequest) (*modelexperiments.GetInstanceResponse, error) {
			if settings.instanceFails {
				return nil, fmt.Errorf("could not get instance")
			}
			return settings.instanceResp, nil
		}),
		GetInstanceTokenExecuteMock: cliutils.Ptr(func(modelexperiments.ApiGetInstanceTokenRequest) (*modelexperiments.GetInstanceTokenResponse, error) {
			if settings.tokenFails {
				return nil, fmt.Errorf("could not get token")
			}
			return settings.tokenResp, nil
		}),
	}
}

func TestGetInstanceName(t *testing.T) {
	name := "instance"
	tests := []struct {
		name     string
		settings mockSettings
		want     string
		wantErr  bool
	}{
		{name: "success", settings: mockSettings{instanceResp: &modelexperiments.GetInstanceResponse{Instance: modelexperiments.Instance{Name: name}}}, want: name},
		{name: "request error", settings: mockSettings{instanceFails: true}, wantErr: true},
		{name: "nil response", settings: mockSettings{}, wantErr: true},
		{name: "nil instance", settings: mockSettings{instanceResp: &modelexperiments.GetInstanceResponse{}}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetInstanceName(context.Background(), newAPIMock(tt.settings), testProjectID, testRegion, testInstanceID)
			if (err != nil) != tt.wantErr {
				t.Fatalf("error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && got != tt.want {
				t.Fatalf("name = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestGetTokenName(t *testing.T) {
	name := "token"
	tests := []struct {
		name     string
		settings mockSettings
		want     string
		wantErr  bool
	}{
		{name: "success", settings: mockSettings{tokenResp: &modelexperiments.GetInstanceTokenResponse{Token: modelexperiments.TokenMetadata{Name: name}}}, want: name},
		{name: "request error", settings: mockSettings{tokenFails: true}, wantErr: true},
		{name: "nil response", settings: mockSettings{}, wantErr: true},
		{name: "zero token", settings: mockSettings{tokenResp: &modelexperiments.GetInstanceTokenResponse{}}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetTokenName(context.Background(), newAPIMock(tt.settings), testProjectID, testRegion, testInstanceID, testTokenID)
			if (err != nil) != tt.wantErr {
				t.Fatalf("error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && got != tt.want {
				t.Fatalf("name = %q, want %q", got, tt.want)
			}
		})
	}
}
