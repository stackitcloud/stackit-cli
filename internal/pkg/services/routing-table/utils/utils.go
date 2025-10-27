package utils

import (
	"fmt"
	"strings"

	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/iaas"
)

func ExtractRouteDetails(item iaas.Route) (destType, destValue, hopType, hopValue, labels string) {
	if item.Destination.DestinationCIDRv4 != nil {
		destType = utils.PtrString(item.Destination.DestinationCIDRv4.Type)
		destValue = utils.PtrString(item.Destination.DestinationCIDRv4.Value)
	} else if item.Destination.DestinationCIDRv6 != nil {
		destType = utils.PtrString(item.Destination.DestinationCIDRv6.Type)
		destValue = utils.PtrString(item.Destination.DestinationCIDRv6.Value)
	}

	if item.Nexthop.NexthopIPv4 != nil {
		hopType = utils.PtrString(item.Nexthop.NexthopIPv4.Type)
		hopValue = utils.PtrString(item.Nexthop.NexthopIPv4.Value)
	} else if item.Nexthop.NexthopIPv6 != nil {
		hopType = utils.PtrString(item.Nexthop.NexthopIPv6.Type)
		hopValue = utils.PtrString(item.Nexthop.NexthopIPv6.Value)
	} else if item.Nexthop.NexthopInternet != nil {
		hopType = utils.PtrString(item.Nexthop.NexthopInternet.Type)
	} else if item.Nexthop.NexthopBlackhole != nil {
		hopType = utils.PtrString(item.Nexthop.NexthopBlackhole.Type)
	}

	var sortedLabels []string
	if item.Labels != nil && len(*item.Labels) > 0 {
		for key, value := range *item.Labels {
			sortedLabels = append(sortedLabels, fmt.Sprintf("%s: %s", key, value))
		}
	}

	return destType, destValue, hopType, hopValue, strings.Join(sortedLabels, ",")
}
