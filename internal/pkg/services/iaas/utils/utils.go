package utils

import (
	"context"
	"errors"
	"fmt"

	"github.com/stackitcloud/stackit-sdk-go/services/iaas"
)

var (
	ErrResponseNil = errors.New("response is nil")
	ErrNameNil     = errors.New("name is nil")
	ErrItemsNil    = errors.New("items is nil")
)

type IaaSClient interface {
	GetSecurityGroupRuleExecute(ctx context.Context, projectId, region, securityGroupRuleId, securityGroupId string) (*iaas.SecurityGroupRule, error)
	GetSecurityGroupExecute(ctx context.Context, projectId, region, securityGroupId string) (*iaas.SecurityGroup, error)
	GetPublicIPExecute(ctx context.Context, projectId, region, publicIpId string) (*iaas.PublicIp, error)
	GetServerExecute(ctx context.Context, projectId, region, serverId string) (*iaas.Server, error)
	GetVolumeExecute(ctx context.Context, projectId, region, volumeId string) (*iaas.Volume, error)
	GetNetworkExecute(ctx context.Context, projectId, region, networkId string) (*iaas.Network, error)
	GetNetworkAreaExecute(ctx context.Context, organizationId, areaId string) (*iaas.NetworkArea, error)
	ListNetworkAreaProjectsExecute(ctx context.Context, organizationId, areaId string) (*iaas.ProjectListResponse, error)
	GetNetworkAreaRangeExecute(ctx context.Context, organizationId, areaId, region, networkRangeId string) (*iaas.NetworkRange, error)
	GetImageExecute(ctx context.Context, projectId, region, imageId string) (*iaas.Image, error)
	GetAffinityGroupExecute(ctx context.Context, projectId, region, affinityGroupId string) (*iaas.AffinityGroup, error)
	GetSnapshotExecute(ctx context.Context, projectId, region, snapshotId string) (*iaas.Snapshot, error)
	GetBackupExecute(ctx context.Context, projectId, region, backupId string) (*iaas.Backup, error)
}

func GetSecurityGroupRuleName(ctx context.Context, apiClient IaaSClient, projectId, region, securityGroupRuleId, securityGroupId string) (string, error) {
	resp, err := apiClient.GetSecurityGroupRuleExecute(ctx, projectId, region, securityGroupRuleId, securityGroupId)
	if err != nil {
		return "", fmt.Errorf("get security group rule: %w", err)
	}
	securityGroupRuleName := *resp.Ethertype + ", " + *resp.Direction
	return securityGroupRuleName, nil
}

func GetSecurityGroupName(ctx context.Context, apiClient IaaSClient, projectId, region, securityGroupId string) (string, error) {
	resp, err := apiClient.GetSecurityGroupExecute(ctx, projectId, region, securityGroupId)
	if err != nil {
		return "", fmt.Errorf("get security group: %w", err)
	} else if resp == nil {
		return "", ErrResponseNil
	} else if resp.Name == nil {
		return "", ErrNameNil
	}
	return *resp.Name, nil
}

func GetPublicIP(ctx context.Context, apiClient IaaSClient, projectId, region, publicIpId string) (ip, associatedResource string, err error) {
	resp, err := apiClient.GetPublicIPExecute(ctx, projectId, region, publicIpId)
	if err != nil {
		return "", "", fmt.Errorf("get public ip: %w", err)
	}
	associatedResourceId := ""
	if resp.NetworkInterface != nil {
		associatedResourceId = *resp.NetworkInterface.Get()
	}
	return *resp.Ip, associatedResourceId, nil
}

func GetServerName(ctx context.Context, apiClient IaaSClient, projectId, region, serverId string) (string, error) {
	resp, err := apiClient.GetServerExecute(ctx, projectId, region, serverId)
	if err != nil {
		return "", fmt.Errorf("get server: %w", err)
	}
	return *resp.Name, nil
}

func GetVolumeName(ctx context.Context, apiClient IaaSClient, projectId, region, volumeId string) (string, error) {
	resp, err := apiClient.GetVolumeExecute(ctx, projectId, region, volumeId)
	if err != nil {
		return "", fmt.Errorf("get volume: %w", err)
	} else if resp == nil {
		return "", ErrResponseNil
	} else if resp.Name == nil {
		return "", ErrNameNil
	}
	return *resp.Name, nil
}

func GetNetworkName(ctx context.Context, apiClient IaaSClient, projectId, region, networkId string) (string, error) {
	resp, err := apiClient.GetNetworkExecute(ctx, projectId, region, networkId)
	if err != nil {
		return "", fmt.Errorf("get network: %w", err)
	} else if resp == nil {
		return "", ErrResponseNil
	} else if resp.Name == nil {
		return "", ErrNameNil
	}
	return *resp.Name, nil
}

func GetNetworkAreaName(ctx context.Context, apiClient IaaSClient, organizationId, areaId string) (string, error) {
	resp, err := apiClient.GetNetworkAreaExecute(ctx, organizationId, areaId)
	if err != nil {
		return "", fmt.Errorf("get network area: %w", err)
	} else if resp == nil {
		return "", ErrResponseNil
	} else if resp.Name == nil {
		return "", ErrNameNil
	}
	return *resp.Name, nil
}

func ListAttachedProjects(ctx context.Context, apiClient IaaSClient, organizationId, areaId string) ([]string, error) {
	resp, err := apiClient.ListNetworkAreaProjectsExecute(ctx, organizationId, areaId)
	if err != nil {
		return nil, fmt.Errorf("list network area attached projects: %w", err)
	} else if resp == nil {
		return nil, ErrResponseNil
	} else if resp.Items == nil {
		return nil, ErrItemsNil
	}
	return *resp.Items, nil
}

func GetNetworkRangePrefix(ctx context.Context, apiClient IaaSClient, organizationId, areaId, region, networkRangeId string) (string, error) {
	resp, err := apiClient.GetNetworkAreaRangeExecute(ctx, organizationId, areaId, region, networkRangeId)
	if err != nil {
		return "", fmt.Errorf("get network range: %w", err)
	}
	return *resp.Prefix, nil
}

// GetRouteFromAPIResponse returns the static route from the API response that matches the prefix and nexthop
// This works because static routes are unique by prefix and nexthop
func GetRouteFromAPIResponse(destination, nexthop string, routes *[]iaas.Route) (iaas.Route, error) {
	for _, route := range *routes {
		// Check if destination matches
		if dest := route.Destination; dest != nil {
			match := false
			if destV4 := dest.DestinationCIDRv4; destV4 != nil {
				if destV4.Value != nil && *destV4.Value == destination {
					match = true
				}
			} else if destV6 := dest.DestinationCIDRv6; destV6 != nil {
				if destV6.Value != nil && *destV6.Value == destination {
					match = true
				}
			}
			if !match {
				continue
			}
		}
		// Check if nexthop matches
		if routeNexthop := route.Nexthop; routeNexthop != nil {
			match := false
			if nexthopIPv4 := routeNexthop.NexthopIPv4; nexthopIPv4 != nil {
				if nexthopIPv4.Value != nil && *nexthopIPv4.Value == nexthop {
					match = true
				}
			} else if nexthopIPv6 := routeNexthop.NexthopIPv6; nexthopIPv6 != nil {
				if nexthopIPv6.Value != nil && *nexthopIPv6.Value == nexthop {
					match = true
				}
			} else if nexthopInternet := routeNexthop.NexthopInternet; nexthopInternet != nil {
				if nexthopInternet.Type != nil && *nexthopInternet.Type == nexthop {
					match = true
				}
			} else if nexthopBlackhole := routeNexthop.NexthopBlackhole; nexthopBlackhole != nil {
				if nexthopBlackhole.Type != nil && *nexthopBlackhole.Type == nexthop {
					match = true
				}
			}
			if match {
				return route, nil
			}
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

func GetImageName(ctx context.Context, apiClient IaaSClient, projectId, region, imageId string) (string, error) {
	resp, err := apiClient.GetImageExecute(ctx, projectId, region, imageId)
	if err != nil {
		return "", fmt.Errorf("get image: %w", err)
	} else if resp == nil {
		return "", ErrResponseNil
	} else if resp.Name == nil {
		return "", ErrNameNil
	}
	return *resp.Name, nil
}

func GetAffinityGroupName(ctx context.Context, apiClient IaaSClient, projectId, region, affinityGroupId string) (string, error) {
	resp, err := apiClient.GetAffinityGroupExecute(ctx, projectId, region, affinityGroupId)
	if err != nil {
		return "", fmt.Errorf("get affinity group: %w", err)
	} else if resp == nil {
		return "", ErrResponseNil
	} else if resp.Name == nil {
		return "", ErrNameNil
	}
	return *resp.Name, nil
}

func GetSnapshotName(ctx context.Context, apiClient IaaSClient, projectId, region, snapshotId string) (string, error) {
	resp, err := apiClient.GetSnapshotExecute(ctx, projectId, region, snapshotId)
	if err != nil {
		return "", fmt.Errorf("get snapshot: %w", err)
	} else if resp == nil {
		return "", ErrResponseNil
	} else if resp.Name == nil {
		return "", ErrNameNil
	}
	return *resp.Name, nil
}

func GetBackupName(ctx context.Context, apiClient IaaSClient, projectId, region, backupId string) (string, error) {
	resp, err := apiClient.GetBackupExecute(ctx, projectId, region, backupId)
	if err != nil {
		return backupId, fmt.Errorf("get backup: %w", err)
	}
	if resp != nil && resp.Name != nil {
		return *resp.Name, nil
	}
	return backupId, nil
}
