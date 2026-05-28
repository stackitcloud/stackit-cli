package utils

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	iaas "github.com/stackitcloud/stackit-sdk-go/services/iaas/v2api"

	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

type IaaSClientMocked struct {
	GetSecurityGroupRuleFails  bool
	GetSecurityGroupRuleResp   *iaas.SecurityGroupRule
	GetSecurityGroupFails      bool
	GetSecurityGroupResp       *iaas.SecurityGroup
	GetPublicIpFails           bool
	GetPublicIpResp            *iaas.PublicIp
	GetServerFails             bool
	GetServerResp              *iaas.Server
	GetVolumeFails             bool
	GetVolumeResp              *iaas.Volume
	GetNetworkFails            bool
	GetNetworkResp             *iaas.Network
	GetRoutingTableOfAreaFails bool
	GetRoutingTableOfAreaResp  *iaas.RoutingTable
	GetNetworkAreaFails        bool
	GetNetworkAreaResp         *iaas.NetworkArea
	GetAttachedProjectsFails   bool
	GetAttachedProjectsResp    *iaas.ProjectListResponse
	GetNetworkAreaRangeFails   bool
	GetNetworkAreaRangeResp    *iaas.NetworkRange
	GetImageFails              bool
	GetImageResp               *iaas.Image
	GetAffinityGroupsFails     bool
	GetAffinityGroupResp       *iaas.AffinityGroup
	GetBackupFails             bool
	GetBackupResp              *iaas.Backup
	GetSnapshotFails           bool
	GetSnapshotResp            *iaas.Snapshot
}

func newMock(m *IaaSClientMocked) iaas.DefaultAPI {
	return iaas.DefaultAPIServiceMock{
		GetAffinityGroupExecuteMock: utils.Ptr(func(_ iaas.ApiGetAffinityGroupRequest) (*iaas.AffinityGroup, error) {
			if m.GetAffinityGroupsFails {
				return nil, fmt.Errorf("could not get affinity groups")
			}
			return m.GetAffinityGroupResp, nil
		}),
		GetSecurityGroupRuleExecuteMock: utils.Ptr(func(_ iaas.ApiGetSecurityGroupRuleRequest) (*iaas.SecurityGroupRule, error) {
			if m.GetSecurityGroupRuleFails {
				return nil, fmt.Errorf("could not get security group rule")
			}
			return m.GetSecurityGroupRuleResp, nil
		}),
		GetSecurityGroupExecuteMock: utils.Ptr(func(_ iaas.ApiGetSecurityGroupRequest) (*iaas.SecurityGroup, error) {
			if m.GetSecurityGroupFails {
				return nil, fmt.Errorf("could not get security group")
			}
			return m.GetSecurityGroupResp, nil
		}),
		GetPublicIPExecuteMock: utils.Ptr(func(_ iaas.ApiGetPublicIPRequest) (*iaas.PublicIp, error) {
			if m.GetPublicIpFails {
				return nil, fmt.Errorf("could not get public ip")
			}
			return m.GetPublicIpResp, nil
		}),
		GetServerExecuteMock: utils.Ptr(func(_ iaas.ApiGetServerRequest) (*iaas.Server, error) {
			if m.GetServerFails {
				return nil, fmt.Errorf("could not get server")
			}
			return m.GetServerResp, nil
		}),
		GetVolumeExecuteMock: utils.Ptr(func(_ iaas.ApiGetVolumeRequest) (*iaas.Volume, error) {
			if m.GetVolumeFails {
				return nil, fmt.Errorf("could not get volume")
			}
			return m.GetVolumeResp, nil
		}),
		GetNetworkExecuteMock: utils.Ptr(func(_ iaas.ApiGetNetworkRequest) (*iaas.Network, error) {
			if m.GetNetworkFails {
				return nil, fmt.Errorf("could not get network")
			}
			return m.GetNetworkResp, nil
		}),
		GetRoutingTableOfAreaExecuteMock: utils.Ptr(func(_ iaas.ApiGetRoutingTableOfAreaRequest) (*iaas.RoutingTable, error) {
			if m.GetRoutingTableOfAreaFails {
				return nil, fmt.Errorf("could not get routing table")
			}
			return m.GetRoutingTableOfAreaResp, nil
		}),
		GetNetworkAreaExecuteMock: utils.Ptr(func(_ iaas.ApiGetNetworkAreaRequest) (*iaas.NetworkArea, error) {
			if m.GetNetworkAreaFails {
				return nil, fmt.Errorf("could not get network area")
			}
			return m.GetNetworkAreaResp, nil
		}),
		ListNetworkAreaProjectsExecuteMock: utils.Ptr(func(_ iaas.ApiListNetworkAreaProjectsRequest) (*iaas.ProjectListResponse, error) {
			if m.GetAttachedProjectsFails {
				return nil, fmt.Errorf("could not get attached projects")
			}
			return m.GetAttachedProjectsResp, nil
		}),
		GetNetworkAreaRangeExecuteMock: utils.Ptr(func(_ iaas.ApiGetNetworkAreaRangeRequest) (*iaas.NetworkRange, error) {
			if m.GetNetworkAreaRangeFails {
				return nil, fmt.Errorf("could not get network range")
			}
			return m.GetNetworkAreaRangeResp, nil
		}),
		GetImageExecuteMock: utils.Ptr(func(_ iaas.ApiGetImageRequest) (*iaas.Image, error) {
			if m.GetImageFails {
				return nil, fmt.Errorf("could not get image")
			}
			return m.GetImageResp, nil
		}),
		GetBackupExecuteMock: utils.Ptr(func(_ iaas.ApiGetBackupRequest) (*iaas.Backup, error) {
			if m.GetBackupFails {
				return nil, fmt.Errorf("could not get backup")
			}
			return m.GetBackupResp, nil
		}),
		GetSnapshotExecuteMock: utils.Ptr(func(_ iaas.ApiGetSnapshotRequest) (*iaas.Snapshot, error) {
			if m.GetSnapshotFails {
				return nil, fmt.Errorf("could not get snapshot")
			}
			return m.GetSnapshotResp, nil
		}),
	}
}

func TestGetSecurityGroupRuleName(t *testing.T) {
	type args struct {
		getInstanceFails bool
		getInstanceResp  *iaas.SecurityGroupRule
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
				getInstanceResp: &iaas.SecurityGroupRule{
					Ethertype: utils.Ptr("IPv6"),
					Direction: "ingress",
				},
			},
			want: "IPv6, ingress",
		},
		{
			name: "get security group rule fails",
			args: args{
				getInstanceFails: true,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &IaaSClientMocked{
				GetSecurityGroupRuleFails: tt.args.getInstanceFails,
				GetSecurityGroupRuleResp:  tt.args.getInstanceResp,
			}
			got, err := GetSecurityGroupRuleName(context.Background(), newMock(m), "", "", "", "")
			if (err != nil) != tt.wantErr {
				t.Errorf("GetSecurityGroupRuleName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetSecurityGroupRuleName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetSecurityGroupName(t *testing.T) {
	type args struct {
		getInstanceFails bool
		getInstanceResp  *iaas.SecurityGroup
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
				getInstanceResp: &iaas.SecurityGroup{
					Name: "test",
				},
			},
			want: "test",
		},
		{
			name: "get security group fails",
			args: args{
				getInstanceFails: true,
			},
			wantErr: true,
		},
		{
			name: "response is nil",
			args: args{
				getInstanceResp:  nil,
				getInstanceFails: false,
			},
			wantErr: true,
			want:    "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &IaaSClientMocked{
				GetSecurityGroupFails: tt.args.getInstanceFails,
				GetSecurityGroupResp:  tt.args.getInstanceResp,
			}
			got, err := GetSecurityGroupName(context.Background(), newMock(m), "", "", "")
			if (err != nil) != tt.wantErr {
				t.Errorf("GetSecurityGroupName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetSecurityGroupName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetPublicIp(t *testing.T) {
	type args struct {
		getPublicIpFails bool
		getPublicIpResp  *iaas.PublicIp
	}
	tests := []struct {
		name                   string
		args                   args
		wantPublicIp           string
		wantAssociatedResource string
		wantErr                bool
	}{
		{
			name: "base",
			args: args{
				getPublicIpResp: &iaas.PublicIp{
					Ip:               utils.Ptr("1.2.3.4"),
					NetworkInterface: *iaas.NewNullableString(utils.Ptr("5.6.7.8")),
				},
			},
			wantPublicIp:           "1.2.3.4",
			wantAssociatedResource: "5.6.7.8",
		},
		{
			name: "get public ip fails",
			args: args{
				getPublicIpFails: true,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &IaaSClientMocked{
				GetPublicIpFails: tt.args.getPublicIpFails,
				GetPublicIpResp:  tt.args.getPublicIpResp,
			}
			gotPublicIP, gotAssociatedResource, err := GetPublicIP(context.Background(), newMock(m), "", "", "")
			if (err != nil) != tt.wantErr {
				t.Errorf("GetPublicIP() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotPublicIP != tt.wantPublicIp {
				t.Errorf("GetPublicIP() = %v, want public IP %v", gotPublicIP, tt.wantPublicIp)
			}
			if gotAssociatedResource != tt.wantAssociatedResource {
				t.Errorf("GetPublicIP() = %v, want associated resource %v", gotAssociatedResource, tt.wantAssociatedResource)
			}
		})
	}
}

func TestGetServerName(t *testing.T) {
	type args struct {
		getInstanceFails bool
		getInstanceResp  *iaas.Server
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
				getInstanceResp: &iaas.Server{
					Name: "test",
				},
			},
			want: "test",
		},
		{
			name: "get server fails",
			args: args{
				getInstanceFails: true,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &IaaSClientMocked{
				GetServerFails: tt.args.getInstanceFails,
				GetServerResp:  tt.args.getInstanceResp,
			}
			got, err := GetServerName(context.Background(), newMock(m), "", "", "")
			if (err != nil) != tt.wantErr {
				t.Errorf("GetServerName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetServerName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetVolumeName(t *testing.T) {
	type args struct {
		getInstanceFails bool
		getInstanceResp  *iaas.Volume
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
				getInstanceResp: &iaas.Volume{
					Name: utils.Ptr("test"),
				},
			},
			want: "test",
		},
		{
			name: "get volume fails",
			args: args{
				getInstanceFails: true,
			},
			wantErr: true,
		},
		{
			name: "response is nil",
			args: args{
				getInstanceResp:  nil,
				getInstanceFails: false,
			},
			wantErr: true,
			want:    "",
		},
		{
			name: "name in response is empty",
			args: args{
				getInstanceResp: &iaas.Volume{
					Name: nil,
				},
				getInstanceFails: false,
			},
			wantErr: true,
			want:    "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &IaaSClientMocked{
				GetVolumeFails: tt.args.getInstanceFails,
				GetVolumeResp:  tt.args.getInstanceResp,
			}
			got, err := GetVolumeName(context.Background(), newMock(m), "", "", "")
			if (err != nil) != tt.wantErr {
				t.Errorf("GetVolumeName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetVolumeName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetNetworkName(t *testing.T) {
	type args struct {
		getInstanceFails bool
		getInstanceResp  *iaas.Network
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
				getInstanceResp: &iaas.Network{
					Name: "test",
				},
			},
			want: "test",
		},
		{
			name: "get network fails",
			args: args{
				getInstanceFails: true,
			},
			wantErr: true,
		},
		{
			name: "response is nil",
			args: args{
				getInstanceResp:  nil,
				getInstanceFails: false,
			},
			wantErr: true,
			want:    "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &IaaSClientMocked{
				GetNetworkFails: tt.args.getInstanceFails,
				GetNetworkResp:  tt.args.getInstanceResp,
			}
			got, err := GetNetworkName(context.Background(), newMock(m), "", "", "")
			if (err != nil) != tt.wantErr {
				t.Errorf("GetNetworkName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetNetworkName() = %v, want %v", got, tt.want)
			}
		})
	}
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
					Name: "test",
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
			want:    "",
		},
		{
			name: "response is nil",
			args: args{
				getInstanceResp:  nil,
				getInstanceFails: false,
			},
			wantErr: true,
			want:    "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &IaaSClientMocked{
				GetNetworkAreaFails: tt.args.getInstanceFails,
				GetNetworkAreaResp:  tt.args.getInstanceResp,
			}
			got, err := GetNetworkAreaName(context.Background(), newMock(m), "", "")
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
					Items: []string{"test"},
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
			got, err := ListAttachedProjects(context.Background(), newMock(m), "", "")
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
					Prefix: "test",
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
			got, err := GetNetworkRangePrefix(context.Background(), newMock(m), "", "", "", "")
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

func TestGetRouteFromAPIResponse(t *testing.T) {
	type args struct {
		prefix  string
		nexthop string
		routes  []iaas.Route
	}
	tests := []struct {
		name    string
		args    args
		want    iaas.Route
		wantErr bool
	}{
		{
			name: "base",
			args: args{
				prefix:  "1.1.1.0/24",
				nexthop: "1.1.1.1",
				routes: []iaas.Route{
					{
						Destination: iaas.RouteDestination{
							DestinationCIDRv4: &iaas.DestinationCIDRv4{
								Type:  "cidrv4",
								Value: "1.1.1.0/24",
							},
						},
						Nexthop: iaas.RouteNexthop{
							NexthopIPv4: &iaas.NexthopIPv4{
								Type:  "ipv4",
								Value: "1.1.1.1",
							},
						},
					},
					{
						Destination: iaas.RouteDestination{
							DestinationCIDRv4: &iaas.DestinationCIDRv4{
								Type:  "cidrv4",
								Value: "2.2.2.0/24",
							},
						},
						Nexthop: iaas.RouteNexthop{
							NexthopIPv4: &iaas.NexthopIPv4{
								Type:  "ipv4",
								Value: "2.2.2.2",
							},
						},
					},
					{
						Destination: iaas.RouteDestination{
							DestinationCIDRv4: &iaas.DestinationCIDRv4{
								Value: "3.3.3.0/24",
							},
						},
						Nexthop: iaas.RouteNexthop{
							NexthopBlackhole: &iaas.NexthopBlackhole{
								Type: "blackhole",
							},
						},
					},
					{
						Destination: iaas.RouteDestination{
							DestinationCIDRv4: &iaas.DestinationCIDRv4{
								Value: "4.4.4.0/24",
							},
						},
						Nexthop: iaas.RouteNexthop{
							NexthopInternet: &iaas.NexthopInternet{
								Type: "internet",
							},
						},
					},
				},
			},
			want: iaas.Route{
				Destination: iaas.RouteDestination{
					DestinationCIDRv4: &iaas.DestinationCIDRv4{
						Type:  "cidrv4",
						Value: "1.1.1.0/24",
					},
				},
				Nexthop: iaas.RouteNexthop{
					NexthopIPv4: &iaas.NexthopIPv4{
						Type:  "ipv4",
						Value: "1.1.1.1",
					},
				},
			},
		},
		{
			name: "nexthop internet",
			args: args{
				prefix:  "4.4.4.0/24",
				nexthop: "internet",
				routes: []iaas.Route{
					{
						Destination: iaas.RouteDestination{
							DestinationCIDRv4: &iaas.DestinationCIDRv4{
								Value: "1.1.1.0/24",
							},
						},
						Nexthop: iaas.RouteNexthop{
							NexthopIPv4: &iaas.NexthopIPv4{
								Value: "1.1.1.1",
							},
						},
					},
					{
						Destination: iaas.RouteDestination{
							DestinationCIDRv4: &iaas.DestinationCIDRv4{
								Value: "2.2.2.0/24",
							},
						},
						Nexthop: iaas.RouteNexthop{
							NexthopIPv4: &iaas.NexthopIPv4{
								Value: "2.2.2.2",
							},
						},
					},
					{
						Destination: iaas.RouteDestination{
							DestinationCIDRv4: &iaas.DestinationCIDRv4{
								Value: "3.3.3.0/24",
							},
						},
						Nexthop: iaas.RouteNexthop{
							NexthopBlackhole: &iaas.NexthopBlackhole{
								Type: "blackhole",
							},
						},
					},
					{
						Destination: iaas.RouteDestination{
							DestinationCIDRv4: &iaas.DestinationCIDRv4{
								Value: "4.4.4.0/24",
							},
						},
						Nexthop: iaas.RouteNexthop{
							NexthopInternet: &iaas.NexthopInternet{
								Type: "internet",
							},
						},
					},
				},
			},
			want: iaas.Route{
				Destination: iaas.RouteDestination{
					DestinationCIDRv4: &iaas.DestinationCIDRv4{
						Value: "4.4.4.0/24",
					},
				},
				Nexthop: iaas.RouteNexthop{
					NexthopInternet: &iaas.NexthopInternet{
						Type: "internet",
					},
				},
			},
		},
		{
			name: "nexthop blackhole",
			args: args{
				prefix:  "3.3.3.0/24",
				nexthop: "blackhole",
				routes: []iaas.Route{
					{
						Destination: iaas.RouteDestination{
							DestinationCIDRv4: &iaas.DestinationCIDRv4{
								Value: "1.1.1.0/24",
							},
						},
						Nexthop: iaas.RouteNexthop{
							NexthopIPv4: &iaas.NexthopIPv4{
								Value: "1.1.1.1",
							},
						},
					},
					{
						Destination: iaas.RouteDestination{
							DestinationCIDRv4: &iaas.DestinationCIDRv4{
								Value: "2.2.2.0/24",
							},
						},
						Nexthop: iaas.RouteNexthop{
							NexthopIPv4: &iaas.NexthopIPv4{
								Value: "2.2.2.2",
							},
						},
					},
					{
						Destination: iaas.RouteDestination{
							DestinationCIDRv4: &iaas.DestinationCIDRv4{
								Value: "3.3.3.0/24",
							},
						},
						Nexthop: iaas.RouteNexthop{
							NexthopBlackhole: &iaas.NexthopBlackhole{
								Type: "blackhole",
							},
						},
					},
					{
						Destination: iaas.RouteDestination{
							DestinationCIDRv4: &iaas.DestinationCIDRv4{
								Value: "4.4.4.0/24",
							},
						},
						Nexthop: iaas.RouteNexthop{
							NexthopInternet: &iaas.NexthopInternet{
								Type: "internet",
							},
						},
					},
				},
			},
			want: iaas.Route{
				Destination: iaas.RouteDestination{
					DestinationCIDRv4: &iaas.DestinationCIDRv4{
						Value: "3.3.3.0/24",
					},
				},
				Nexthop: iaas.RouteNexthop{
					NexthopBlackhole: &iaas.NexthopBlackhole{
						Type: "blackhole",
					},
				},
			},
		},
		{
			name: "not found",
			args: args{
				prefix:  "1.1.1.0/24",
				nexthop: "1.1.1.1",
				routes: []iaas.Route{
					{
						Destination: iaas.RouteDestination{
							DestinationCIDRv4: &iaas.DestinationCIDRv4{
								Value: "2.2.2.0/24",
							},
						},
						Nexthop: iaas.RouteNexthop{
							NexthopIPv4: &iaas.NexthopIPv4{
								Value: "2.2.2.2",
							},
						},
					},
					{
						Destination: iaas.RouteDestination{
							DestinationCIDRv4: &iaas.DestinationCIDRv4{
								Value: "3.3.3.0/24",
							},
						},
						Nexthop: iaas.RouteNexthop{
							NexthopIPv4: &iaas.NexthopIPv4{
								Value: "3.3.3.3",
							},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "empty",
			args: args{
				prefix:  "1.1.1.0/24",
				nexthop: "1.1.1.1",
				routes:  []iaas.Route{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetRouteFromAPIResponse(tt.args.prefix, tt.args.nexthop, tt.args.routes)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetRouteFromAPIResponse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetRouteFromAPIResponse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetNetworkRangeFromAPIResponse(t *testing.T) {
	type args struct {
		prefix        string
		networkRanges []iaas.NetworkRange
	}
	tests := []struct {
		name    string
		args    args
		want    iaas.NetworkRange
		wantErr bool
	}{
		{
			name: "base",
			args: args{
				prefix: "1.1.1.0/24",
				networkRanges: []iaas.NetworkRange{
					{
						Prefix: "1.1.1.0/24",
					},
					{
						Prefix: "2.2.2.0/24",
					},
					{
						Prefix: "3.3.3.0/24",
					},
				},
			},
			want: iaas.NetworkRange{
				Prefix: "1.1.1.0/24",
			},
		},
		{
			name: "not found",
			args: args{
				prefix: "1.1.1.0/24",
				networkRanges: []iaas.NetworkRange{
					{
						Prefix: "2.2.2.0/24",
					},
					{
						Prefix: "3.3.3.0/24",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "empty",
			args: args{
				prefix:        "1.1.1.0/24",
				networkRanges: []iaas.NetworkRange{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetNetworkRangeFromAPIResponse(tt.args.prefix, tt.args.networkRanges)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetNetworkRangeFromAPIResponse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetNetworkRangeFromAPIResponse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetImageName(t *testing.T) {
	tests := []struct {
		name      string
		imageResp *iaas.Image
		imageErr  bool
		want      string
		wantErr   bool
	}{
		{
			name:      "successful retrieval",
			imageResp: &iaas.Image{Name: "test-image"},
			want:      "test-image",
			wantErr:   false,
		},
		{
			name:     "error on retrieval",
			imageErr: true,
			wantErr:  true,
		},
		{
			name:      "response is nil",
			imageErr:  false,
			imageResp: nil,
			want:      "",
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &IaaSClientMocked{
				GetImageFails: tt.imageErr,
				GetImageResp:  tt.imageResp,
			}
			got, err := GetImageName(context.Background(), newMock(client), "", "", "")
			if (err != nil) != tt.wantErr {
				t.Errorf("GetImageName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetImageName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetAffinityGroupName(t *testing.T) {
	tests := []struct {
		name         string
		affinityResp *iaas.AffinityGroup
		affinityErr  bool
		want         string
		wantErr      bool
	}{
		{
			name:         "successful retrieval",
			affinityResp: &iaas.AffinityGroup{Name: "test-affinity"},
			want:         "test-affinity",
			wantErr:      false,
		},
		{
			name:        "error on retrieval",
			affinityErr: true,
			wantErr:     true,
		},
		{
			name:         "response is nil",
			affinityErr:  false,
			affinityResp: nil,
			want:         "",
			wantErr:      true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			client := &IaaSClientMocked{
				GetAffinityGroupsFails: tt.affinityErr,
				GetAffinityGroupResp:   tt.affinityResp,
			}
			got, err := GetAffinityGroupName(ctx, newMock(client), "", "", "")
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAffinityGroupName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetAffinityGroupName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetRoutingTableOfAreaName(t *testing.T) {
	type args struct {
		getInstanceFails bool
		getInstanceResp  *iaas.RoutingTable
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
				getInstanceResp: &iaas.RoutingTable{
					Name: "test",
				},
			},
			want: "test",
		},
		{
			name: "get routing table fails",
			args: args{
				getInstanceFails: true,
			},
			wantErr: true,
			want:    "",
		},
		{
			name: "response is nil",
			args: args{
				getInstanceResp:  nil,
				getInstanceFails: false,
			},
			wantErr: true,
			want:    "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &IaaSClientMocked{
				GetRoutingTableOfAreaFails: tt.args.getInstanceFails,
				GetRoutingTableOfAreaResp:  tt.args.getInstanceResp,
			}

			got, err := GetRoutingTableOfAreaName(context.Background(), newMock(m), "", "", "", "")
			if (err != nil) != tt.wantErr {
				t.Errorf("GetRoutingTableOfAreaName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetRoutingTableOfAreaName() = %v, want %v", got, tt.want)
			}
		})
	}
}
