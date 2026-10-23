package utils

import (
	"context"
	"errors"
	"fmt"
	"time"

	iaas "github.com/stackitcloud/stackit-sdk-go/services/iaas/v2alpha1api"

	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

var (
	ErrResponseNil = errors.New("response is nil")
)

type RouteDetails struct {
	DestType     string
	DestValue    string
	NextHopType  string
	NextHopValue string
	CreatedAt    string
	UpdatedAt    string
	Labels       string
}

func GetVPCName(ctx context.Context, apiClient iaas.DefaultAPI, projectId, vpcId string) (string, error) {
	resp, err := apiClient.GetVPC(ctx, projectId, vpcId).Execute()
	if err != nil {
		return "", fmt.Errorf("get vpc: %w", err)
	}

	if resp == nil {
		return "", ErrResponseNil
	}

	return resp.Name, nil
}

func GetVPCRoutingTableName(ctx context.Context, apiClient iaas.DefaultAPI, projectId, vpcId, region, routingTableId string) (string, error) {
	resp, err := apiClient.GetVPCRoutingTable(ctx, projectId, vpcId, region, routingTableId).Execute()
	if err != nil {
		return "", fmt.Errorf("get vpc routing-table: %w", err)
	}

	if resp == nil {
		return "", ErrResponseNil
	}

	return resp.Name, nil
}

func GetVPCNetworkRangePrefix(ctx context.Context, apiClient iaas.DefaultAPI, projectId, vpcId, region, networkRangeId string) (string, error) {
	resp, err := apiClient.GetVPCNetworkRange(ctx, projectId, vpcId, region, networkRangeId).Execute()
	if err != nil {
		return "", fmt.Errorf("get vpc network range: %w", err)
	}

	if resp != nil {
		if resp.VPCNetworkRangeIPv4 != nil {
			return resp.VPCNetworkRangeIPv4.Prefix, nil
		}

		if resp.VPCNetworkRangeIPv6 != nil {
			return resp.VPCNetworkRangeIPv6.Prefix, nil
		}
	}

	return "", ErrResponseNil
}

func ExtractRouteDetails(route *iaas.Route) RouteDetails {
	var routeDetails RouteDetails

	if route == nil {
		return routeDetails
	}

	if route.Destination.DestinationCIDRv4 != nil {
		routeDetails.DestType = route.Destination.DestinationCIDRv4.Type
		routeDetails.DestValue = route.Destination.DestinationCIDRv4.Value
	} else if route.Destination.DestinationCIDRv6 != nil {
		routeDetails.DestType = route.Destination.DestinationCIDRv6.Type
		routeDetails.DestValue = route.Destination.DestinationCIDRv6.Value
	}

	if route.Nexthop.NexthopIPv4 != nil {
		routeDetails.NextHopType = route.Nexthop.NexthopIPv4.Type
		routeDetails.NextHopValue = route.Nexthop.NexthopIPv4.Value
	} else if route.Nexthop.NexthopIPv6 != nil {
		routeDetails.NextHopType = route.Nexthop.NexthopIPv6.Type
		routeDetails.NextHopValue = route.Nexthop.NexthopIPv6.Value
	} else if route.Nexthop.NexthopInternet != nil {
		routeDetails.NextHopType = route.Nexthop.NexthopInternet.Type
	} else if route.Nexthop.NexthopBlackhole != nil {
		routeDetails.NextHopType = route.Nexthop.NexthopBlackhole.Type
	}

	if len(route.Labels) > 0 {
		stringMap := make(map[string]string)
		for key, value := range route.Labels {
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
