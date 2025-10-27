package utils

import (
	"testing"

	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/iaas"
)

func TestExtractRouteDetails(t *testing.T) {
	tests := []struct {
		description   string
		input         *iaas.Route
		wantDestType  string
		wantDestValue string
		wantHopType   string
		wantHopValue  string
		wantLabels    string
	}{
		{
			description: "CIDRv4 destination, IPv4 nexthop, with labels",
			input: &iaas.Route{
				Destination: &iaas.RouteDestination{
					DestinationCIDRv4: &iaas.DestinationCIDRv4{
						Type:  utils.Ptr("CIDRv4"),
						Value: utils.Ptr("10.0.0.0/24"),
					},
				},
				Nexthop: &iaas.RouteNexthop{
					NexthopIPv4: &iaas.NexthopIPv4{
						Type:  utils.Ptr("IPv4"),
						Value: utils.Ptr("10.0.0.1"),
					},
				},
				Labels: &map[string]interface{}{
					"key": "value",
				},
			},
			wantDestType:  "CIDRv4",
			wantDestValue: "10.0.0.0/24",
			wantHopType:   "IPv4",
			wantHopValue:  "10.0.0.1",
			wantLabels:    "key=value",
		},
		{
			description: "CIDRv6 destination, IPv6 nexthop, with no labels",
			input: &iaas.Route{
				Destination: &iaas.RouteDestination{
					DestinationCIDRv6: &iaas.DestinationCIDRv6{
						Type:  utils.Ptr("CIDRv6"),
						Value: utils.Ptr("2001:db8::/32"),
					},
				},
				Nexthop: &iaas.RouteNexthop{
					NexthopIPv4: &iaas.NexthopIPv4{
						Type:  utils.Ptr("IPv6"),
						Value: utils.Ptr("2001:db8::1"),
					},
				},
				Labels: nil,
			},
			wantDestType:  "CIDRv6",
			wantDestValue: "2001:db8::/32",
			wantHopType:   "IPv6",
			wantHopValue:  "2001:db8::1",
			wantLabels:    "",
		},
		{
			description: "Internet nexthop without value",
			input: &iaas.Route{
				Destination: &iaas.RouteDestination{
					DestinationCIDRv4: &iaas.DestinationCIDRv4{
						Type:  utils.Ptr("CIDRv4"),
						Value: utils.Ptr("0.0.0.0/0"),
					},
				},
				Nexthop: &iaas.RouteNexthop{
					NexthopInternet: &iaas.NexthopInternet{
						Type: utils.Ptr("Internet"),
					},
				},
				Labels: nil,
			},
			wantDestType:  "CIDRv4",
			wantDestValue: "0.0.0.0/0",
			wantHopType:   "Internet",
			wantHopValue:  "",
			wantLabels:    "",
		},
		{
			description: "Blackhole nexthop without value and nil labels map",
			input: &iaas.Route{
				Destination: &iaas.RouteDestination{
					DestinationCIDRv6: &iaas.DestinationCIDRv6{
						Type:  utils.Ptr("CIDRv6"),
						Value: utils.Ptr("::/0"),
					},
				},
				Nexthop: &iaas.RouteNexthop{
					NexthopBlackhole: &iaas.NexthopBlackhole{
						Type: utils.Ptr("Blackhole"),
					},
				},
				Labels: nil,
			},
			wantDestType:  "CIDRv6",
			wantDestValue: "::/0",
			wantHopType:   "Blackhole",
			wantHopValue:  "",
			wantLabels:    "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			destType, destValue, hopType, hopValue, labels := ExtractRouteDetails(*tt.input)

			if destType != tt.wantDestType {
				t.Errorf("destType = %v, want %v", destType, tt.wantDestType)
			}
			if destValue != tt.wantDestValue {
				t.Errorf("destValue = %v, want %v", destValue, tt.wantDestValue)
			}
			if hopType != tt.wantHopType {
				t.Errorf("hopType = %v, want %v", hopType, tt.wantHopType)
			}
			if hopValue != tt.wantHopValue {
				t.Errorf("hopValue = %v, want %v", hopValue, tt.wantHopValue)
			}
			if (tt.wantLabels != "" && labels == "") || (tt.wantLabels == "" && labels != "") {
				t.Errorf("labels mismatch: got %q, want %q", labels, tt.wantLabels)
			}
		})
	}
}
