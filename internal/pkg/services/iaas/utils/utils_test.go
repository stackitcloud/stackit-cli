package utils

import (
	"context"
	"fmt"
	"testing"

	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/iaas"
)

type IaaSClientMocked struct {
	GetNetworkAreaFails      bool
	GetNetworkAreaResp       *iaas.NetworkArea
	GetAttachedProjectsFails bool
	GetAttachedProjectsResp  *iaas.ProjectListResponse
	GetNetworkAreaRangeFails bool
	GetNetworkAreaRangeResp  *iaas.NetworkRange
}

func (m *IaaSClientMocked) GetNetworkAreaExecute(_ context.Context, _, _ string) (*iaas.NetworkArea, error) {
	if m.GetNetworkAreaFails {
		return nil, fmt.Errorf("could not get network area")
	}
	return m.GetNetworkAreaResp, nil
}

func (m *IaaSClientMocked) ListNetworkAreaProjectsExecute(_ context.Context, _, _ string) (*iaas.ProjectListResponse, error) {
	if m.GetAttachedProjectsFails {
		return nil, fmt.Errorf("could not get attached projects")
	}
	return m.GetAttachedProjectsResp, nil
}

func (m *IaaSClientMocked) GetNetworkAreaRangeExecute(_ context.Context, _, _, _ string) (*iaas.NetworkRange, error) {
	if m.GetNetworkAreaRangeFails {
		return nil, fmt.Errorf("could not get network range")
	}
	return m.GetNetworkAreaRangeResp, nil
}

func TestGetNetworkAreaName(t *testing.T) {
	type args struct {
		getInstanceFails bool
		getInstanceResp  *iaas.NetworkArea
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "base",
			args: args{
				getInstanceResp: &iaas.NetworkArea{
					Name: utils.Ptr("test"),
				},
			},
			want: "test",
		},
		{
			name: "get network area fails",
			args: args{
				getInstanceFails: true,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &IaaSClientMocked{
				GetNetworkAreaFails: tt.args.getInstanceFails,
				GetNetworkAreaResp:  tt.args.getInstanceResp,
			}
			got, err := GetNetworkAreaName(context.Background(), m, "", "")
			if (err != nil) != tt.wantErr {
				t.Errorf("GetNetworkAreaName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetNetworkAreaName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestListAttachedProjects(t *testing.T) {
	type args struct {
		getAttachedProjectsFails bool
		getAttachedProjectsResp  *iaas.ProjectListResponse
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "base",
			args: args{
				getAttachedProjectsResp: &iaas.ProjectListResponse{
					Items: &[]string{"test"},
				},
			},
			want: []string{"test"},
		},
		{
			name: "get attached projects fails",
			args: args{
				getAttachedProjectsFails: true,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &IaaSClientMocked{
				GetAttachedProjectsFails: tt.args.getAttachedProjectsFails,
				GetAttachedProjectsResp:  tt.args.getAttachedProjectsResp,
			}
			got, err := ListAttachedProjects(context.Background(), m, "", "")
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAttachedProjects() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if fmt.Sprintf("%v", got) != fmt.Sprintf("%v", tt.want) {
				t.Errorf("GetAttachedProjects() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetNetworkRangePrefix(t *testing.T) {
	type args struct {
		getNetworkAreaRangeFails bool
		getNetworkAreaRangeResp  *iaas.NetworkRange
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "base",
			args: args{
				getNetworkAreaRangeResp: &iaas.NetworkRange{
					Prefix: utils.Ptr("test"),
				},
			},
			want: "test",
		},
		{
			name: "get network area range fails",
			args: args{
				getNetworkAreaRangeFails: true,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &IaaSClientMocked{
				GetNetworkAreaRangeFails: tt.args.getNetworkAreaRangeFails,
				GetNetworkAreaRangeResp:  tt.args.getNetworkAreaRangeResp,
			}
			got, err := GetNetworkRangePrefix(context.Background(), m, "", "", "")
			if (err != nil) != tt.wantErr {
				t.Errorf("GetNetworkRangePrefix() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetNetworkRangePrefix() = %v, want %v", got, tt.want)
			}
		})
	}
}
