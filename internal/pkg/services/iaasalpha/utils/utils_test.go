package utils

import (
	"context"
	"fmt"
	"testing"

	iaas "github.com/stackitcloud/stackit-sdk-go/services/iaas/v2alpha1api"
)

type IaaSClientMocked struct {
	GetVPCFails             bool
	GetVPCResp              *iaas.VPC
	GetVPCNetworkRangeFails bool
	GetVPCNetworkRangeResp  *iaas.VPCNetworkRange
}

func newMock(m *IaaSClientMocked) iaas.DefaultAPI {
	return iaas.DefaultAPIServiceMock{
		GetVPCExecuteMock: new(func(_ iaas.ApiGetVPCRequest) (*iaas.VPC, error) {
			if m.GetVPCFails {
				return nil, fmt.Errorf("could not get vpc")
			}
			return m.GetVPCResp, nil
		}),
		GetVPCNetworkRangeExecuteMock: new(func(_ iaas.ApiGetVPCNetworkRangeRequest) (*iaas.VPCNetworkRange, error) {
			if m.GetVPCNetworkRangeFails {
				return nil, fmt.Errorf("could not get vpc network range")
			}
			return m.GetVPCNetworkRangeResp, nil
		}),
	}
}

func TestGetVPCName(t *testing.T) {
	type args struct {
		getVpcFails bool
		getVpcResp  *iaas.VPC
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
				getVpcResp: &iaas.VPC{
					Name: "test-vpc",
				},
			},
			want: "test-vpc",
		},
		{
			name: "nil response",
			args: args{
				getVpcResp: nil,
			},
			wantErr: true,
		},
		{
			name: "get vpc fails",
			args: args{
				getVpcFails: true,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &IaaSClientMocked{
				GetVPCFails: tt.args.getVpcFails,
				GetVPCResp:  tt.args.getVpcResp,
			}
			got, err := GetVPCName(context.Background(), newMock(m), "", "")
			if (err != nil) != tt.wantErr {
				t.Errorf("GetVPCName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetVPCName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetVPCNetworkRangePrefix(t *testing.T) {
	type args struct {
		getVpcNetworkRangeFails bool
		getVpcNetworkRangeResp  *iaas.VPCNetworkRange
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "base ipv4",
			args: args{
				getVpcNetworkRangeResp: &iaas.VPCNetworkRange{
					VPCNetworkRangeIPv4: &iaas.VPCNetworkRangeIPv4{
						Prefix: "test",
					},
				},
			},
			want: "test",
		},
		{
			name: "base ipv6",
			args: args{
				getVpcNetworkRangeResp: &iaas.VPCNetworkRange{
					VPCNetworkRangeIPv6: &iaas.VPCNetworkRangeIPv6{
						Prefix: "test",
					},
				},
			},
			want: "test",
		},
		{
			name: "nil response",
			args: args{
				getVpcNetworkRangeResp: nil,
			},
			wantErr: true,
		},
		{
			name: "get network range fails",
			args: args{
				getVpcNetworkRangeFails: true,
			},
			wantErr: true,
		},
		{
			name: "nil ipv4 and ipv6",
			args: args{
				getVpcNetworkRangeResp: &iaas.VPCNetworkRange{
					VPCNetworkRangeIPv4: nil,
					VPCNetworkRangeIPv6: nil,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &IaaSClientMocked{
				GetVPCNetworkRangeFails: tt.args.getVpcNetworkRangeFails,
				GetVPCNetworkRangeResp:  tt.args.getVpcNetworkRangeResp,
			}
			got, err := GetVPCNetworkRangePrefix(context.Background(), newMock(m), "", "", "", "")
			if (err != nil) != tt.wantErr {
				t.Errorf("GetVPCNetworkRangePrefix() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetVPCNetworkRangePrefix() = %v, want %v", got, tt.want)
			}
		})
	}
}
