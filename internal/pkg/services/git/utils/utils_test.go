package utils

import (
	"context"
	"fmt"
	"testing"

	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/git"
)

type GitClientMocked struct {
	GetInstanceFails bool
	GetInstanceResp  *git.Instance
}

func (m *GitClientMocked) GetInstanceExecute(_ context.Context, _, _ string) (*git.Instance, error) {
	if m.GetInstanceFails {
		return nil, fmt.Errorf("could not get instance")
	}
	return m.GetInstanceResp, nil
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
			name:         "successful retrieval",
			instanceResp: &git.Instance{Name: utils.Ptr("test-instance")},
			want:         "test-instance",
			wantErr:      false,
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
			client := &GitClientMocked{
				GetInstanceFails: tt.instanceErr,
				GetInstanceResp:  tt.instanceResp,
			}
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
