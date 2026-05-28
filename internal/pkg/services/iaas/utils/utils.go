package utils

import (
	"context"
	"errors"
	"fmt"

	iaas "github.com/stackitcloud/stackit-sdk-go/services/iaas/v2api"
)

var (
	ErrResponseNil = errors.New("response is nil")
	ErrNameNil     = errors.New("name is nil")
	ErrItemsNil    = errors.New("items is nil")
)

func GetSecurityGroupRuleName(ctx context.Context, apiClient iaas.DefaultAPI, projectId, region, securityGroupRuleId, securityGroupId string) (string, error) {
	resp, err := apiClient.GetSecurityGroupRule(ctx, projectId, region, securityGroupRuleId, securityGroupId).Execute()
	if err != nil {
		return "", fmt.Errorf("get security group rule: %w", err)
	}
	securityGroupRuleName := *resp.Ethertype + ", " + resp.Direction
	return securityGroupRuleName, nil
}

func GetSecurityGroupName(ctx context.Context, apiClient iaas.DefaultAPI, projectId, region, securityGroupId string) (string, error) {
	resp, err := apiClient.GetSecurityGroup(ctx, projectId, region, securityGroupId).Execute()
	if err != nil {
		return "", fmt.Errorf("get security group: %w", err)
	} else if resp == nil {
		return "", ErrResponseNil
	}
	return resp.Name, nil
}

func GetPublicIP(ctx context.Context, apiClient iaas.DefaultAPI, projectId, region, publicIpId string) (ip, associatedResource string, err error) {
	resp, err := apiClient.GetPublicIP(ctx, projectId, region, publicIpId).Execute()
	if err != nil {
		return "", "", fmt.Errorf("get public ip: %w", err)
	}
	associatedResourceId := ""
	if resp.NetworkInterface.IsSet() {
		associatedResourceId = *resp.NetworkInterface.Get()
	}
	return *resp.Ip, associatedResourceId, nil
}

func GetServerName(ctx context.Context, apiClient iaas.DefaultAPI, projectId, region, serverId string) (string, error) {
	resp, err := apiClient.GetServer(ctx, projectId, region, serverId).Execute()
	if err != nil {
		return "", fmt.Errorf("get server: %w", err)
	}
	return resp.Name, nil
}

func GetVolumeName(ctx context.Context, apiClient iaas.DefaultAPI, projectId, region, volumeId string) (string, error) {
	resp, err := apiClient.GetVolume(ctx, projectId, region, volumeId).Execute()
	if err != nil {
		return "", fmt.Errorf("get volume: %w", err)
	} else if resp == nil {
		return "", ErrResponseNil
	} else if resp.Name == nil {
		return "", ErrNameNil
	}
	return *resp.Name, nil
}

func GetNetworkName(ctx context.Context, apiClient iaas.DefaultAPI, projectId, region, networkId string) (string, error) {
	resp, err := apiClient.GetNetwork(ctx, projectId, region, networkId).Execute()
	if err != nil {
		return "", fmt.Errorf("get network: %w", err)
	} else if resp == nil {
		return "", ErrResponseNil
	}
	return resp.Name, nil
}

func GetRoutingTableOfAreaName(ctx context.Context, apiClient iaas.DefaultAPI, organizationId, areaId, region, routingTableId string) (string, error) {
	resp, err := apiClient.GetRoutingTableOfArea(ctx, organizationId, areaId, region, routingTableId).Execute()
	if err != nil {
		return "", fmt.Errorf("get routing-table: %w", err)
	} else if resp == nil {
		return "", ErrResponseNil
	}
	return resp.Name, nil
}

func GetNetworkAreaName(ctx context.Context, apiClient iaas.DefaultAPI, organizationId, areaId string) (string, error) {
	resp, err := apiClient.GetNetworkArea(ctx, organizationId, areaId).Execute()
	if err != nil {
		return "", fmt.Errorf("get network area: %w", err)
	} else if resp == nil {
		return "", ErrResponseNil
	}
	return resp.Name, nil
}

func ListAttachedProjects(ctx context.Context, apiClient iaas.DefaultAPI, organizationId, areaId string) ([]string, error) {
	resp, err := apiClient.ListNetworkAreaProjects(ctx, organizationId, areaId).Execute()
	if err != nil {
		return nil, fmt.Errorf("list network area attached projects: %w", err)
	} else if resp == nil {
		return nil, ErrResponseNil
	} else if resp.Items == nil {
		return nil, ErrItemsNil
	}
	return resp.Items, nil
}

func GetNetworkRangePrefix(ctx context.Context, apiClient iaas.DefaultAPI, organizationId, areaId, region, networkRangeId string) (string, error) {
	resp, err := apiClient.GetNetworkAreaRange(ctx, organizationId, areaId, region, networkRangeId).Execute()
	if err != nil {
		return "", fmt.Errorf("get network range: %w", err)
	}
	return resp.Prefix, nil
}

// GetRouteFromAPIResponse returns the static route from the API response that matches the prefix and nexthop
// This works because static routes are unique by prefix and nexthop
func GetRouteFromAPIResponse(destination, nexthop string, routes []iaas.Route) (iaas.Route, error) {
	for _, route := range routes {
		// Check if destination matches
		destV4 := route.Destination.DestinationCIDRv4
		destV4Matches := destV4 != nil && destV4.Value == destination
		destV6 := route.Destination.DestinationCIDRv6
		destV6Matches := destV6 != nil && destV6.Value == destination
		destMatches := destV4Matches || destV6Matches
		if !destMatches {
			continue
		}
		// Check if nexthop matches
		nextHopV4 := route.Nexthop.NexthopIPv4
		nextHopV4Matches := nextHopV4 != nil && nextHopV4.Value == nexthop
		nextHopV6 := route.Nexthop.NexthopIPv6
		nextHopV6Matches := nextHopV6 != nil && nextHopV6.Value == nexthop
		nextHopInet := route.Nexthop.NexthopInternet
		nextHopInetMatches := nextHopInet != nil && nextHopInet.Type == nexthop
		nextHopBlackhole := route.Nexthop.NexthopBlackhole
		nextHopBlackholeMatches := nextHopBlackhole != nil && nextHopBlackhole.Type == nexthop
		nextHopMatches := nextHopV4Matches || nextHopV6Matches || nextHopInetMatches || nextHopBlackholeMatches
		if nextHopMatches {
			return route, nil
		}
	}
	return iaas.Route{}, fmt.Errorf("new static route not found in API response")
}

// GetNetworkRangeFromAPIResponse returns the network range from the API response that matches the given prefix
// This works because network range prefixes are unique in the same SNA
func GetNetworkRangeFromAPIResponse(prefix string, networkRanges []iaas.NetworkRange) (iaas.NetworkRange, error) {
	for _, networkRange := range networkRanges {
		if networkRange.Prefix == prefix {
			return networkRange, nil
		}
	}
	return iaas.NetworkRange{}, fmt.Errorf("new network range not found in API response")
}

func GetImageName(ctx context.Context, apiClient iaas.DefaultAPI, projectId, region, imageId string) (string, error) {
	resp, err := apiClient.GetImage(ctx, projectId, region, imageId).Execute()
	if err != nil {
		return "", fmt.Errorf("get image: %w", err)
	} else if resp == nil {
		return "", ErrResponseNil
	}
	return resp.Name, nil
}

func GetAffinityGroupName(ctx context.Context, apiClient iaas.DefaultAPI, projectId, region, affinityGroupId string) (string, error) {
	resp, err := apiClient.GetAffinityGroup(ctx, projectId, region, affinityGroupId).Execute()
	if err != nil {
		return "", fmt.Errorf("get affinity group: %w", err)
	} else if resp == nil {
		return "", ErrResponseNil
	}
	return resp.Name, nil
}

func GetSnapshotName(ctx context.Context, apiClient iaas.DefaultAPI, projectId, region, snapshotId string) (string, error) {
	resp, err := apiClient.GetSnapshot(ctx, projectId, region, snapshotId).Execute()
	if err != nil {
		return "", fmt.Errorf("get snapshot: %w", err)
	} else if resp == nil {
		return "", ErrResponseNil
	} else if resp.Name == nil {
		return "", ErrNameNil
	}
	return *resp.Name, nil
}

func GetBackupName(ctx context.Context, apiClient iaas.DefaultAPI, projectId, region, backupId string) (string, error) {
	resp, err := apiClient.GetBackup(ctx, projectId, region, backupId).Execute()
	if err != nil {
		return backupId, fmt.Errorf("get backup: %w", err)
	}
	if resp != nil && resp.Name != nil {
		return *resp.Name, nil
	}
	return backupId, nil
}
