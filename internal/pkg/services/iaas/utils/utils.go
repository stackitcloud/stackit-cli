package utils

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-sdk-go/services/iaas"
)

type IaaSClient interface {
	GetPublicIPExecute(ctx context.Context, projectId, publicIpId string) (*iaas.PublicIp, error)
	GetServerExecute(ctx context.Context, projectId, serverId string) (*iaas.Server, error)
	GetVolumeExecute(ctx context.Context, projectId, volumeId string) (*iaas.Volume, error)
	GetNetworkExecute(ctx context.Context, projectId, networkId string) (*iaas.Network, error)
	GetNetworkAreaExecute(ctx context.Context, organizationId, areaId string) (*iaas.NetworkArea, error)
	ListNetworkAreaProjectsExecute(ctx context.Context, organizationId, areaId string) (*iaas.ProjectListResponse, error)
	GetNetworkAreaRangeExecute(ctx context.Context, organizationId, areaId, networkRangeId string) (*iaas.NetworkRange, error)
}

func GetPublicIP(ctx context.Context, apiClient IaaSClient, projectId, publicIpId string) (ip, associatedResource string, err error) {
	resp, err := apiClient.GetPublicIPExecute(ctx, projectId, publicIpId)
	if err != nil {
		return "", "", fmt.Errorf("get public ip: %w", err)
	}
	associatedResourceId := ""
	if resp.GetNetworkInterface() != nil {
		associatedResourceId = *resp.GetNetworkInterface()
	}
	return *resp.Ip, associatedResourceId, nil
}

func GetServerName(ctx context.Context, apiClient IaaSClient, projectId, serverId string) (string, error) {
	resp, err := apiClient.GetServerExecute(ctx, projectId, serverId)
	if err != nil {
		return "", fmt.Errorf("get server: %w", err)
	}
	return *resp.Name, nil
}

func GetVolumeName(ctx context.Context, apiClient IaaSClient, projectId, volumeId string) (string, error) {
	resp, err := apiClient.GetVolumeExecute(ctx, projectId, volumeId)
	if err != nil {
		return "", fmt.Errorf("get volume: %w", err)
	}
	return *resp.Name, nil
}

func GetNetworkName(ctx context.Context, apiClient IaaSClient, projectId, networkId string) (string, error) {
	resp, err := apiClient.GetNetworkExecute(ctx, projectId, networkId)
	if err != nil {
		return "", fmt.Errorf("get network: %w", err)
	}
	return *resp.Name, nil
}

func GetNetworkAreaName(ctx context.Context, apiClient IaaSClient, organizationId, areaId string) (string, error) {
	resp, err := apiClient.GetNetworkAreaExecute(ctx, organizationId, areaId)
	if err != nil {
		return "", fmt.Errorf("get network area: %w", err)
	}
	return *resp.Name, nil
}

func ListAttachedProjects(ctx context.Context, apiClient IaaSClient, organizationId, areaId string) ([]string, error) {
	resp, err := apiClient.ListNetworkAreaProjectsExecute(ctx, organizationId, areaId)
	if err != nil {
		return nil, fmt.Errorf("list network area attached projects: %w", err)
	}
	return *resp.Items, nil
}

func GetNetworkRangePrefix(ctx context.Context, apiClient IaaSClient, organizationId, areaId, networkRangeId string) (string, error) {
	resp, err := apiClient.GetNetworkAreaRangeExecute(ctx, organizationId, areaId, networkRangeId)
	if err != nil {
		return "", fmt.Errorf("get network range: %w", err)
	}
	return *resp.Prefix, nil
}

// GetRouteFromAPIResponse returns the static route from the API response that matches the prefix and nexthop
// This works because static routes are unique by prefix and nexthop
func GetRouteFromAPIResponse(prefix, nexthop string, routes *[]iaas.Route) (iaas.Route, error) {
	for _, route := range *routes {
		if *route.Prefix == prefix && *route.Nexthop == nexthop {
			return route, nil
		}
	}
	return iaas.Route{}, fmt.Errorf("new static route not found in API response")
}

// GetNetworkRangeFromAPIResponse returns the network range from the API response that matches the given prefix
// This works because network range prefixes are unique in the same SNA
func GetNetworkRangeFromAPIResponse(prefix string, networkRanges *[]iaas.NetworkRange) (iaas.NetworkRange, error) {
	for _, networkRange := range *networkRanges {
		if *networkRange.Prefix == prefix {
			return networkRange, nil
		}
	}
	return iaas.NetworkRange{}, fmt.Errorf("new network range not found in API response")
}
