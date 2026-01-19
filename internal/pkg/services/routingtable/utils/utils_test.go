package utils

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/iaas"
)

const ipv4 = "ipv4"
const ipv6 = "ipv6"
const cidrv4 = "cidrv4"
const cidrv6 = "cidrv6"

func TestExtractRouteDetails(t *testing.T) {
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
				Destination: &iaas.RouteDestination{
					DestinationCIDRv4: &iaas.DestinationCIDRv4{
						Type:  utils.Ptr(cidrv4),
						Value: utils.Ptr("10.0.0.0/24"),
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
				Nexthop: &iaas.RouteNexthop{
					NexthopIPv4: &iaas.NexthopIPv4{
						Type:  utils.Ptr(ipv4),
						Value: utils.Ptr("10.0.0.1"),
					},
				},
				Labels: &map[string]interface{}{}, // empty but non-nil
			},
			want: RouteDetails{
				HopType:  ipv4,
				HopValue: "10.0.0.1",
			},
		},
		{
			description: "destination present, nexthop struct nil, labels nil",
			input: &iaas.Route{
				Destination: &iaas.RouteDestination{
					DestinationCIDRv6: &iaas.DestinationCIDRv6{
						Type:  utils.Ptr(cidrv6),
						Value: utils.Ptr("2001:db8::/32"),
					},
				},
				Nexthop: nil,
				Labels:  nil,
			},
			want: RouteDetails{
				DestType:  cidrv6,
				DestValue: "2001:db8::/32",
			},
		},
		{
			description: "CIDRv4 destination, IPv4 nexthop, with labels",
			input: &iaas.Route{
				Destination: &iaas.RouteDestination{
					DestinationCIDRv4: &iaas.DestinationCIDRv4{
						Type:  utils.Ptr(cidrv4),
						Value: utils.Ptr("10.0.0.0/24"),
					},
				},
				Nexthop: &iaas.RouteNexthop{
					NexthopIPv4: &iaas.NexthopIPv4{
						Type:  utils.Ptr(ipv4),
						Value: utils.Ptr("10.0.0.1"),
					},
				},
				Labels: &map[string]interface{}{
					"key": "value",
				},
			},
			want: RouteDetails{
				DestType:  cidrv4,
				DestValue: "10.0.0.0/24",
				HopType:   ipv4,
				HopValue:  "10.0.0.1",
				Labels:    "key: value",
			},
		},
		{
			description: "CIDRv6 destination, IPv6 nexthop, with no labels",
			input: &iaas.Route{
				Destination: &iaas.RouteDestination{
					DestinationCIDRv6: &iaas.DestinationCIDRv6{
						Type:  utils.Ptr(cidrv6),
						Value: utils.Ptr("2001:db8::/32"),
					},
				},
				Nexthop: &iaas.RouteNexthop{
					NexthopIPv6: &iaas.NexthopIPv6{
						Type:  utils.Ptr(ipv6),
						Value: utils.Ptr("2001:db8::1"),
					},
				},
				Labels: nil,
			},
			want: RouteDetails{
				DestType:  cidrv6,
				DestValue: "2001:db8::/32",
				HopType:   ipv6,
				HopValue:  "2001:db8::1",
			},
		},
		{
			description: "Internet nexthop without value",
			input: &iaas.Route{
				Destination: &iaas.RouteDestination{
					DestinationCIDRv4: &iaas.DestinationCIDRv4{
						Type:  utils.Ptr(cidrv4),
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
			want: RouteDetails{
				DestType:  cidrv4,
				DestValue: "0.0.0.0/0",
				HopType:   "Internet",
				// HopValue empty
			},
		},
		{
			description: "Blackhole nexthop without value and nil labels map",
			input: &iaas.Route{
				Destination: &iaas.RouteDestination{
					DestinationCIDRv6: &iaas.DestinationCIDRv6{
						Type:  utils.Ptr(cidrv6),
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
			want: RouteDetails{
				DestType:  cidrv6,
				DestValue: "::/0",
				HopType:   "Blackhole",
			},
		},
		{
			description: "route with created and updated timestamps",
			input: &iaas.Route{
				Destination: &iaas.RouteDestination{
					DestinationCIDRv4: &iaas.DestinationCIDRv4{
						Type:  utils.Ptr(cidrv4),
						Value: utils.Ptr("10.0.0.0/24"),
					},
				},
				Nexthop: &iaas.RouteNexthop{
					NexthopIPv4: &iaas.NexthopIPv4{
						Type:  utils.Ptr(ipv4),
						Value: utils.Ptr("10.0.0.1"),
					},
				},
				CreatedAt: &created,
				UpdatedAt: &updated,
			},
			want: RouteDetails{
				DestType:  cidrv4,
				DestValue: "10.0.0.0/24",
				HopType:   ipv4,
				HopValue:  "10.0.0.1",
				CreatedAt: created.Format(time.RFC3339),
				UpdatedAt: updated.Format(time.RFC3339),
				Labels:    "",
			},
		},
		{
			description: "route with created timestamp only",
			input: &iaas.Route{
				Destination: &iaas.RouteDestination{
					DestinationCIDRv4: &iaas.DestinationCIDRv4{
						Type:  utils.Ptr(cidrv4),
						Value: utils.Ptr("10.0.0.0/24"),
					},
				},
				Nexthop: &iaas.RouteNexthop{
					NexthopIPv4: &iaas.NexthopIPv4{
						Type:  utils.Ptr(ipv4),
						Value: utils.Ptr("10.0.0.1"),
					},
				},
				CreatedAt: &created,
			},
			want: RouteDetails{
				DestType:  cidrv4,
				DestValue: "10.0.0.0/24",
				HopType:   ipv4,
				HopValue:  "10.0.0.1",
				CreatedAt: created.Format(time.RFC3339),
				UpdatedAt: "",
				Labels:    "",
			},
		},
		{
			description: "route with updated timestamp only",
			input: &iaas.Route{
				Destination: &iaas.RouteDestination{
					DestinationCIDRv4: &iaas.DestinationCIDRv4{
						Type:  utils.Ptr(cidrv4),
						Value: utils.Ptr("10.0.0.0/24"),
					},
				},
				Nexthop: &iaas.RouteNexthop{
					NexthopIPv4: &iaas.NexthopIPv4{
						Type:  utils.Ptr(ipv4),
						Value: utils.Ptr("10.0.0.1"),
					},
				},
				UpdatedAt: &updated,
			},
			want: RouteDetails{
				DestType:  cidrv4,
				DestValue: "10.0.0.0/24",
				HopType:   ipv4,
				HopValue:  "10.0.0.1",
				CreatedAt: "",
				UpdatedAt: updated.Format(time.RFC3339),
				Labels:    "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			got := ExtractRouteDetails(*tt.input)

			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Fatalf("ExtractRouteDetails() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
