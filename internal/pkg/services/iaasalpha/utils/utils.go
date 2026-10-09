package utils

import (
	"context"
	"errors"
	"fmt"

	iaas "github.com/stackitcloud/stackit-sdk-go/services/iaas/v2alpha1api"
)

var (
	ErrResponseNil = errors.New("response is nil")
)

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
