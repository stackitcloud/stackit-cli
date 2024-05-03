package utils

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-sdk-go/services/loadbalancer"
)

const (
// service = "loadbalancer"
)

type LoadBalancerClient interface {
	GetCredentialsExecute(ctx context.Context, projectId, credentialsRef string) (*loadbalancer.GetCredentialsResponse, error)
	GetLoadBalancerExecute(ctx context.Context, projectId, name string) (*loadbalancer.LoadBalancer, error)
	UpdateTargetPool(ctx context.Context, projectId, loadBalancerName, targetPoolName string) loadbalancer.ApiUpdateTargetPoolRequest
}

func GetCredentialsDisplayName(ctx context.Context, apiClient LoadBalancerClient, projectId, credentialsRef string) (string, error) {
	resp, err := apiClient.GetCredentialsExecute(ctx, projectId, credentialsRef)
	if err != nil {
		return "", fmt.Errorf("get Load Balancer credentials: %w", err)
	}
	return *resp.Credential.DisplayName, nil
}

func GetLoadBalancerTargetPool(ctx context.Context, apiClient LoadBalancerClient, projectId, loadBalancerName, targetPoolName string) (*loadbalancer.TargetPool, error) {
	resp, err := apiClient.GetLoadBalancerExecute(ctx, projectId, loadBalancerName)
	if err != nil {
		return nil, fmt.Errorf("get load balancer: %w", err)
	}

	if resp.TargetPools == nil {
		return nil, fmt.Errorf("no target pools found")
	}

	for _, targetPool := range *resp.TargetPools {
		if *targetPool.Name == targetPoolName {
			return &targetPool, nil
		}
	}

	return nil, fmt.Errorf("target pool not found")
}

func AddTargetToTargetPool(targetPool *loadbalancer.TargetPool, target *loadbalancer.Target) error {
	if targetPool == nil {
		return fmt.Errorf("target pool is nil")
	}
	if target == nil {
		return fmt.Errorf("target is nil")
	}
	if targetPool.Targets == nil {
		targetPool.Targets = &[]loadbalancer.Target{*target}
		return nil
	}
	*targetPool.Targets = append(*targetPool.Targets, *target)
	return nil
}

func ToPayloadTargetPool(targetPool *loadbalancer.TargetPool) *loadbalancer.UpdateTargetPoolPayload {
	if targetPool == nil {
		return nil
	}
	return &loadbalancer.UpdateTargetPoolPayload{
		Name:               targetPool.Name,
		ActiveHealthCheck:  targetPool.ActiveHealthCheck,
		SessionPersistence: targetPool.SessionPersistence,
		TargetPort:         targetPool.TargetPort,
		Targets:            targetPool.Targets,
	}
}
