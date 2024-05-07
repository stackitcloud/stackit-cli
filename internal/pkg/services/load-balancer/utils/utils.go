package utils

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-sdk-go/services/loadbalancer"
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

	if resp == nil {
		return nil, fmt.Errorf("no load balancer found")
	}
	if resp.TargetPools == nil {
		return nil, fmt.Errorf("no target pools found")
	}

	targetPool := FindLoadBalancerTargetPoolByName(resp.TargetPools, targetPoolName)
	if targetPool == nil {
		return nil, fmt.Errorf("target pool not found")
	}
	return targetPool, nil
}

func FindLoadBalancerTargetPoolByName(targetPools *[]loadbalancer.TargetPool, targetPoolName string) *loadbalancer.TargetPool {
	if targetPools == nil {
		return nil
	}
	for _, targetPool := range *targetPools {
		if targetPool.Name != nil && *targetPool.Name == targetPoolName {
			return &targetPool
		}
	}
	return nil
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

func RemoveTargetFromTargetPool(targetPool *loadbalancer.TargetPool, ip string) error {
	if targetPool == nil {
		return fmt.Errorf("target pool is nil")
	}
	if targetPool.Targets == nil {
		return fmt.Errorf("no targets found")
	}
	targets := *targetPool.Targets
	for i, target := range targets {
		if target.Ip != nil && *target.Ip == ip {
			newTargets := targets[:i]
			newTargets = append(newTargets, targets[i+1:]...)
			*targetPool.Targets = newTargets
			return nil
		}
	}
	return fmt.Errorf("target not found")
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

func GetTargetName(ctx context.Context, apiClient LoadBalancerClient, projectId, loadBalancerName, targetPoolName, targetIp string) (string, error) {
	targetPool, err := GetLoadBalancerTargetPool(ctx, apiClient, projectId, loadBalancerName, targetPoolName)
	if err != nil {
		return "", fmt.Errorf("get target pool: %w", err)
	}
	if targetPool.Targets == nil {
		return "", fmt.Errorf("no targets found")
	}
	for _, target := range *targetPool.Targets {
		if target.Ip != nil && *target.Ip == targetIp {
			if target.DisplayName == nil {
				return "", fmt.Errorf("nil target display name")
			}
			return *target.DisplayName, nil
		}
	}
	return "", fmt.Errorf("target not found")
}
