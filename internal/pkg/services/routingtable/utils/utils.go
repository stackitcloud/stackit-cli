package utils

import (
	"fmt"
	"time"

	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/iaas"
)

type RouteDetails struct {
	DestType  string
	DestValue string
	HopType   string
	HopValue  string
	CreatedAt string
	UpdatedAt string
	Labels    string
}

func ExtractRouteDetails(route iaas.Route) RouteDetails {
	var routeDetails RouteDetails

	if route.Destination != nil {
		if route.Destination.DestinationCIDRv4 != nil {
			routeDetails.DestType = utils.PtrString(route.Destination.DestinationCIDRv4.Type)
			routeDetails.DestValue = utils.PtrString(route.Destination.DestinationCIDRv4.Value)
		} else if route.Destination.DestinationCIDRv6 != nil {
			routeDetails.DestType = utils.PtrString(route.Destination.DestinationCIDRv6.Type)
			routeDetails.DestValue = utils.PtrString(route.Destination.DestinationCIDRv6.Value)
		}
	}

	if route.Nexthop != nil {
		if route.Nexthop.NexthopIPv4 != nil {
			routeDetails.HopType = utils.PtrString(route.Nexthop.NexthopIPv4.Type)
			routeDetails.HopValue = utils.PtrString(route.Nexthop.NexthopIPv4.Value)
		} else if route.Nexthop.NexthopIPv6 != nil {
			routeDetails.HopType = utils.PtrString(route.Nexthop.NexthopIPv6.Type)
			routeDetails.HopValue = utils.PtrString(route.Nexthop.NexthopIPv6.Value)
		} else if route.Nexthop.NexthopInternet != nil {
			routeDetails.HopType = utils.PtrString(route.Nexthop.NexthopInternet.Type)
		} else if route.Nexthop.NexthopBlackhole != nil {
			routeDetails.HopType = utils.PtrString(route.Nexthop.NexthopBlackhole.Type)
		}
	}

	if route.Labels != nil && len(*route.Labels) > 0 {
		stringMap := make(map[string]string)
		for key, value := range *route.Labels {
			stringMap[key] = fmt.Sprintf("%v", value)
		}
		routeDetails.Labels = utils.JoinStringMap(stringMap, ": ", "\n")
	}

	if route.CreatedAt != nil {
		routeDetails.CreatedAt = route.CreatedAt.Format(time.RFC3339)
	}

	if route.UpdatedAt != nil {
		routeDetails.UpdatedAt = route.UpdatedAt.Format(time.RFC3339)
	}

	return routeDetails
}
