package utils

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"

	modelexperiments "github.com/stackitcloud/stackit-sdk-go/services/modelexperiments/v1api"

	cliutils "github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

const (
	testRegion       = "eu01"
	testInstanceName = "instance"
	testTokenName    = "token"
)

var (
	testProjectID  = uuid.NewString()
	testInstanceID = uuid.NewString()
	testTokenID    = uuid.NewString()
)

type mockSettings struct {
	getInstanceFails bool
	getInstanceResp  *modelexperiments.GetInstanceResponse
	getTokenFails    bool
	getTokenResp     *modelexperiments.GetInstanceTokenResponse
}

func newAPIMock(settings mockSettings) modelexperiments.DefaultAPI {
	return &modelexperiments.DefaultAPIServiceMock{
		GetInstanceExecuteMock: cliutils.Ptr(func(modelexperiments.ApiGetInstanceRequest) (*modelexperiments.GetInstanceResponse, error) {
			if settings.getInstanceFails {
				return nil, fmt.Errorf("could not get instance")
			}
			return settings.getInstanceResp, nil
		}),
		GetInstanceTokenExecuteMock: cliutils.Ptr(func(modelexperiments.ApiGetInstanceTokenRequest) (*modelexperiments.GetInstanceTokenResponse, error) {
			if settings.getTokenFails {
				return nil, fmt.Errorf("could not get token")
			}
			return settings.getTokenResp, nil
		}),
	}
}

func TestGetInstanceName(t *testing.T) {
	tests := []struct {
		name     string
		settings mockSettings
		want     string
		wantErr  bool
	}{
		{
			name: "success",
			settings: mockSettings{
				getInstanceResp: &modelexperiments.GetInstanceResponse{
					Instance: modelexperiments.Instance{Name: testInstanceName},
				},
			},
			want: testInstanceName,
		},
		{
			name:     "request error",
			settings: mockSettings{getInstanceFails: true},
			wantErr:  true,
		},
		{
			name:     "nil response",
			settings: mockSettings{},
			wantErr:  true,
		},
		{
			name:     "nil instance",
			settings: mockSettings{getInstanceResp: &modelexperiments.GetInstanceResponse{}},
			wantErr:  true,
		},
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
	tests := []struct {
		name     string
		settings mockSettings
		want     string
		wantErr  bool
	}{
		{
			name: "success",
			settings: mockSettings{
				getTokenResp: &modelexperiments.GetInstanceTokenResponse{
					Token: modelexperiments.TokenMetadata{Name: testTokenName},
				},
			},
			want: testTokenName,
		},
		{
			name:     "request error",
			settings: mockSettings{getTokenFails: true},
			wantErr:  true,
		},
		{
			name:     "nil response",
			settings: mockSettings{},
			wantErr:  true,
		},
		{
			name:     "empty token name",
			settings: mockSettings{getTokenResp: &modelexperiments.GetInstanceTokenResponse{}},
			wantErr:  true,
		},
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
