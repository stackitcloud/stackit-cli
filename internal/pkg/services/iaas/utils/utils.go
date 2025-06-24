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
	GetSecurityGroupRuleExecute(ctx context.Context, projectId, securityGroupRuleId, securityGroupId string) (*iaas.SecurityGroupRule, error)
	GetSecurityGroupExecute(ctx context.Context, projectId, securityGroupId string) (*iaas.SecurityGroup, error)
	GetPublicIPExecute(ctx context.Context, projectId, publicIpId string) (*iaas.PublicIp, error)
	GetServerExecute(ctx context.Context, projectId, serverId string) (*iaas.Server, error)
	GetVolumeExecute(ctx context.Context, projectId, volumeId string) (*iaas.Volume, error)
	GetNetworkExecute(ctx context.Context, projectId, networkId string) (*iaas.Network, error)
	GetNetworkAreaExecute(ctx context.Context, organizationId, areaId string) (*iaas.NetworkArea, error)
	ListNetworkAreaProjectsExecute(ctx context.Context, organizationId, areaId string) (*iaas.ProjectListResponse, error)
	GetNetworkAreaRangeExecute(ctx context.Context, organizationId, areaId, networkRangeId string) (*iaas.NetworkRange, error)
	GetImageExecute(ctx context.Context, projectId string, imageId string) (*iaas.Image, error)
	GetAffinityGroupExecute(ctx context.Context, projectId string, affinityGroupId string) (*iaas.AffinityGroup, error)
	GetSnapshotExecute(ctx context.Context, projectId, snapshotId string) (*iaas.Snapshot, error)
	GetBackupExecute(ctx context.Context, projectId, backupId string) (*iaas.Backup, error)
}

func GetSecurityGroupRuleName(ctx context.Context, apiClient IaaSClient, projectId, securityGroupRuleId, securityGroupId string) (string, error) {
	resp, err := apiClient.GetSecurityGroupRuleExecute(ctx, projectId, securityGroupRuleId, securityGroupId)
	if err != nil {
		return "", fmt.Errorf("get security group rule: %w", err)
	}
	securityGroupRuleName := *resp.Ethertype + ", " + *resp.Direction
	return securityGroupRuleName, nil
}

func GetSecurityGroupName(ctx context.Context, apiClient IaaSClient, projectId, securityGroupId string) (string, error) {
	resp, err := apiClient.GetSecurityGroupExecute(ctx, projectId, securityGroupId)
	if err != nil {
		return "", fmt.Errorf("get security group: %w", err)
	} else if resp == nil {
		return "", ErrResponseNil
	} else if resp.Name == nil {
		return "", ErrNameNil
	}
	return *resp.Name, nil
}

func GetPublicIP(ctx context.Context, apiClient IaaSClient, projectId, publicIpId string) (ip, associatedResource string, err error) {
	resp, err := apiClient.GetPublicIPExecute(ctx, projectId, publicIpId)
	if err != nil {
		return "", "", fmt.Errorf("get public ip: %w", err)
	}
	associatedResourceId := ""
	if resp.NetworkInterface != nil {
		associatedResourceId = *resp.NetworkInterface.Get()
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
	} else if resp == nil {
		return "", ErrResponseNil
	} else if resp.Name == nil {
		return "", ErrNameNil
	}
	return *resp.Name, nil
}

func GetNetworkName(ctx context.Context, apiClient IaaSClient, projectId, networkId string) (string, error) {
	resp, err := apiClient.GetNetworkExecute(ctx, projectId, networkId)
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

func GetImageName(ctx context.Context, apiClient IaaSClient, projectId, imageId string) (string, error) {
	resp, err := apiClient.GetImageExecute(ctx, projectId, imageId)
	if err != nil {
		return "", fmt.Errorf("get image: %w", err)
	} else if resp == nil {
		return "", ErrResponseNil
	} else if resp.Name == nil {
		return "", ErrNameNil
	}
	return *resp.Name, nil
}

func GetAffinityGroupName(ctx context.Context, apiClient IaaSClient, projectId, affinityGroupId string) (string, error) {
	resp, err := apiClient.GetAffinityGroupExecute(ctx, projectId, affinityGroupId)
	if err != nil {
		return "", fmt.Errorf("get affinity group: %w", err)
	} else if resp == nil {
		return "", ErrResponseNil
	} else if resp.Name == nil {
		return "", ErrNameNil
	}
	return *resp.Name, nil
}

func GetSnapshotName(ctx context.Context, apiClient IaaSClient, projectId, snapshotId string) (string, error) {
	resp, err := apiClient.GetSnapshotExecute(ctx, projectId, snapshotId)
	if err != nil {
		return "", fmt.Errorf("get snapshot: %w", err)
	} else if resp == nil {
		return "", ErrResponseNil
	} else if resp.Name == nil {
		return "", ErrNameNil
	}
	return *resp.Name, nil
}

func GetBackupName(ctx context.Context, apiClient IaaSClient, projectId, backupId string) (string, error) {
	resp, err := apiClient.GetBackupExecute(ctx, projectId, backupId)
	if err != nil {
		return backupId, fmt.Errorf("get backup: %w", err)
	}
	if resp != nil && resp.Name != nil {
		return *resp.Name, nil
	}
	return backupId, nil
}
