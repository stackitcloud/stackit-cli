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
