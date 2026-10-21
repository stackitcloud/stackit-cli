package utils

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	iaas "github.com/stackitcloud/stackit-sdk-go/services/iaas/v2alpha1api"
)

type IaaSClientMocked struct {
	GetVPCFails             bool
	GetVPCResp              *iaas.VPC
	GetVPCNetworkRangeFails bool
	GetVPCNetworkRangeResp  *iaas.VPCNetworkRange
	GetVPCRoutingTableFails bool
	GetVPCRoutingTableResp  *iaas.VPCRoutingTable
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
		GetVPCRoutingTableExecuteMock: new(func(_ iaas.ApiGetVPCRoutingTableRequest) (*iaas.VPCRoutingTable, error) {
			if m.GetVPCRoutingTableFails {
				return nil, fmt.Errorf("could not get vpc routing table")
			}

			return m.GetVPCRoutingTableResp, nil
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

func TestGetVPCRoutingTable(t *testing.T) {
	type args struct {
		getVpcRoutingTableFails bool
		getVpcRoutingTableResp  *iaas.VPCRoutingTable
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
				getVpcRoutingTableResp: &iaas.VPCRoutingTable{
					Name: "test-vpc",
				},
			},
			want: "test-vpc",
		},
		{
			name: "nil response",
			args: args{
				getVpcRoutingTableResp: nil,
			},
			wantErr: true,
		},
		{
			name: "get fails",
			args: args{
				getVpcRoutingTableFails: true,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &IaaSClientMocked{
				GetVPCRoutingTableFails: tt.args.getVpcRoutingTableFails,
				GetVPCRoutingTableResp:  tt.args.getVpcRoutingTableResp,
			}
			got, err := GetVPCRoutingTableName(context.Background(), newMock(m), "", "", "", "")

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

func TestExtractRouteDetails(t *testing.T) {
	const (
		ipv4   = "ipv4"
		ipv6   = "ipv6"
		cidrv4 = "cidrv4"
		cidrv6 = "cidrv6"
	)

	created := time.Date(2024, 1, 2, 3, 4, 5, 0, time.UTC)
	updated := time.Date(2024, 1, 2, 4, 5, 6, 0, time.UTC)

	tests := []struct {
		description string
		input       *iaas.Route
		want        RouteDetails
	}{
		{
			description: "completely empty route (zero value)",
			input:       &iaas.Route{},
			want:        RouteDetails{},
		},
		{
			description: "destination only, no nexthop, no labels",
			input: &iaas.Route{
				Destination: iaas.RouteDestination{
					DestinationCIDRv4: &iaas.DestinationCIDRv4{
						Type:  cidrv4,
						Value: "10.0.0.0/24",
					},
				},
			},
			want: RouteDetails{
				DestType:  cidrv4,
				DestValue: "10.0.0.0/24",
			},
		},
		{
			description: "nexthop only, no destination, empty labels map",
			input: &iaas.Route{
				Nexthop: iaas.RouteNexthop{
					NexthopIPv4: &iaas.NexthopIPv4{
						Type:  ipv4,
						Value: "10.0.0.1",
					},
				},
				Labels: map[string]any{}, // empty but non-nil
			},
			want: RouteDetails{
				NextHopType:  ipv4,
				NextHopValue: "10.0.0.1",
			},
		},
		{
			description: "destination present, nexthop struct nil, labels nil",
			input: &iaas.Route{
				Destination: iaas.RouteDestination{
					DestinationCIDRv6: &iaas.DestinationCIDRv6{
						Type:  cidrv6,
						Value: "2001:db8::/32",
					},
				},
				Labels: nil,
			},
			want: RouteDetails{
				DestType:  cidrv6,
				DestValue: "2001:db8::/32",
			},
		},
		{
			description: "CIDRv4 destination, IPv4 nexthop, with labels",
			input: &iaas.Route{
				Destination: iaas.RouteDestination{
					DestinationCIDRv4: &iaas.DestinationCIDRv4{
						Type:  cidrv4,
						Value: "10.0.0.0/24",
					},
				},
				Nexthop: iaas.RouteNexthop{
					NexthopIPv4: &iaas.NexthopIPv4{
						Type:  ipv4,
						Value: "10.0.0.1",
					},
				},
				Labels: map[string]any{
					"key": "value",
				},
			},
			want: RouteDetails{
				DestType:     cidrv4,
				DestValue:    "10.0.0.0/24",
				NextHopType:  ipv4,
				NextHopValue: "10.0.0.1",
				Labels:       "key: value",
			},
		},
		{
			description: "CIDRv6 destination, IPv6 nexthop, with no labels",
			input: &iaas.Route{
				Destination: iaas.RouteDestination{
					DestinationCIDRv6: &iaas.DestinationCIDRv6{
						Type:  cidrv6,
						Value: "2001:db8::/32",
					},
				},
				Nexthop: iaas.RouteNexthop{
					NexthopIPv6: &iaas.NexthopIPv6{
						Type:  ipv6,
						Value: "2001:db8::1",
					},
				},
				Labels: nil,
			},
			want: RouteDetails{
				DestType:     cidrv6,
				DestValue:    "2001:db8::/32",
				NextHopType:  ipv6,
				NextHopValue: "2001:db8::1",
			},
		},
		{
			description: "Internet nexthop without value",
			input: &iaas.Route{
				Destination: iaas.RouteDestination{
					DestinationCIDRv4: &iaas.DestinationCIDRv4{
						Type:  cidrv4,
						Value: "0.0.0.0/0",
					},
				},
				Nexthop: iaas.RouteNexthop{
					NexthopInternet: &iaas.NexthopInternet{
						Type: "Internet",
					},
				},
				Labels: nil,
			},
			want: RouteDetails{
				DestType:    cidrv4,
				DestValue:   "0.0.0.0/0",
				NextHopType: "Internet",
				// HopValue empty
			},
		},
		{
			description: "Blackhole nexthop without value and nil labels map",
			input: &iaas.Route{
				Destination: iaas.RouteDestination{
					DestinationCIDRv6: &iaas.DestinationCIDRv6{
						Type:  cidrv6,
						Value: "::/0",
					},
				},
				Nexthop: iaas.RouteNexthop{
					NexthopBlackhole: &iaas.NexthopBlackhole{
						Type: "Blackhole",
					},
				},
				Labels: nil,
			},
			want: RouteDetails{
				DestType:    cidrv6,
				DestValue:   "::/0",
				NextHopType: "Blackhole",
			},
		},
		{
			description: "route with created and updated timestamps",
			input: &iaas.Route{
				Destination: iaas.RouteDestination{
					DestinationCIDRv4: &iaas.DestinationCIDRv4{
						Type:  cidrv4,
						Value: "10.0.0.0/24",
					},
				},
				Nexthop: iaas.RouteNexthop{
					NexthopIPv4: &iaas.NexthopIPv4{
						Type:  ipv4,
						Value: "10.0.0.1",
					},
				},
				CreatedAt: &created,
				UpdatedAt: &updated,
			},
			want: RouteDetails{
				DestType:     cidrv4,
				DestValue:    "10.0.0.0/24",
				NextHopType:  ipv4,
				NextHopValue: "10.0.0.1",
				CreatedAt:    created.Format(time.RFC3339),
				UpdatedAt:    updated.Format(time.RFC3339),
				Labels:       "",
			},
		},
		{
			description: "route with created timestamp only",
			input: &iaas.Route{
				Destination: iaas.RouteDestination{
					DestinationCIDRv4: &iaas.DestinationCIDRv4{
						Type:  cidrv4,
						Value: "10.0.0.0/24",
					},
				},
				Nexthop: iaas.RouteNexthop{
					NexthopIPv4: &iaas.NexthopIPv4{
						Type:  ipv4,
						Value: "10.0.0.1",
					},
				},
				CreatedAt: &created,
			},
			want: RouteDetails{
				DestType:     cidrv4,
				DestValue:    "10.0.0.0/24",
				NextHopType:  ipv4,
				NextHopValue: "10.0.0.1",
				CreatedAt:    created.Format(time.RFC3339),
				UpdatedAt:    "",
				Labels:       "",
			},
		},
		{
			description: "route with updated timestamp only",
			input: &iaas.Route{
				Destination: iaas.RouteDestination{
					DestinationCIDRv4: &iaas.DestinationCIDRv4{
						Type:  cidrv4,
						Value: "10.0.0.0/24",
					},
				},
				Nexthop: iaas.RouteNexthop{
					NexthopIPv4: &iaas.NexthopIPv4{
						Type:  ipv4,
						Value: "10.0.0.1",
					},
				},
				UpdatedAt: &updated,
			},
			want: RouteDetails{
				DestType:     cidrv4,
				DestValue:    "10.0.0.0/24",
				NextHopType:  ipv4,
				NextHopValue: "10.0.0.1",
				CreatedAt:    "",
				UpdatedAt:    updated.Format(time.RFC3339),
				Labels:       "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			got := ExtractRouteDetails(tt.input)

			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Fatalf("ExtractRouteDetails() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
