package utils

import (
	"context"
	"fmt"
	"testing"

	git "github.com/stackitcloud/stackit-sdk-go/services/git/v1betaapi"

	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

type mockSettings struct {
	GetInstanceFails bool
	GetInstanceResp  *git.Instance
}

func newAPIMock(settings *mockSettings) git.DefaultAPI {
	return &git.DefaultAPIServiceMock{
		GetInstanceExecuteMock: utils.Ptr(func(_ git.ApiGetInstanceRequest) (*git.Instance, error) {
			if settings.GetInstanceFails {
				return nil, fmt.Errorf("could not get instance details")
			}

			return settings.GetInstanceResp, nil
		}),
	}
}

func TestGetinstanceName(t *testing.T) {
	tests := []struct {
		name         string
		instanceResp *git.Instance
		instanceErr  bool
		want         string
		wantErr      bool
	}{
		{
			name: "successful retrieval",
			instanceResp: &git.Instance{
				Name: "test-instance",
			},
			want:    "test-instance",
			wantErr: false,
		},
		{
			name:        "error on retrieval",
			instanceErr: true,
			wantErr:     true,
		},
		{
			name:         "nil name",
			instanceErr:  false,
			instanceResp: &git.Instance{},
			want:         "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := newAPIMock(&mockSettings{
				GetInstanceFails: tt.instanceErr,
				GetInstanceResp:  tt.instanceResp,
			})
			got, err := GetInstanceName(context.Background(), client, "", "")
			if (err != nil) != tt.wantErr {
				t.Errorf("GetInstanceName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetInstanceName() = %v, want %v", got, tt.want)
			}
		})
	}
}
