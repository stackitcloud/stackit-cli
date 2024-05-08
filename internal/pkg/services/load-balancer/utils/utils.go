package utils

import (
	"context"
	"fmt"
	"slices"

	"github.com/stackitcloud/stackit-sdk-go/services/loadbalancer"
)

type LoadBalancerClient interface {
	GetCredentialsExecute(ctx context.Context, projectId, credentialsRef string) (*loadbalancer.GetCredentialsResponse, error)
	GetLoadBalancerExecute(ctx context.Context, projectId, name string) (*loadbalancer.LoadBalancer, error)
	UpdateTargetPool(ctx context.Context, projectId, loadBalancerName, targetPoolName string) loadbalancer.ApiUpdateTargetPoolRequest
	ListLoadBalancersExecute(ctx context.Context, projectId string) (*loadbalancer.ListLoadBalancersResponse, error)
	ListCredentialsExecute(ctx context.Context, projectId string) (*loadbalancer.ListCredentialsResponse, error)
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

	targetPool := FindLoadBalancerTargetPoolByName(*resp.TargetPools, targetPoolName)
	if targetPool == nil {
		return nil, fmt.Errorf("target pool not found")
	}
	return targetPool, nil
}

func FindLoadBalancerTargetPoolByName(targetPools []loadbalancer.TargetPool, targetPoolName string) *loadbalancer.TargetPool {
	if targetPools == nil {
		return nil
	}
	for _, targetPool := range targetPools {
		if targetPool.Name != nil && *targetPool.Name == targetPoolName {
			return &targetPool
		}
	}
	return nil
}

func FindLoadBalancerListenerByTargetPool(listeners []loadbalancer.Listener, targetPoolName string) *loadbalancer.Listener {
	if listeners == nil {
		return nil
	}
	for _, listener := range listeners {
		if listener.TargetPool != nil && *listener.TargetPool == targetPoolName {
			return &listener
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

// GetUsedObsCredentials returns a list of credentials that are used by load balancers for observability metrics or logs.
// It goes through all load balancers and checks what credentials are being used, then returns a list of those credentials.
func GetUsedObsCredentials(ctx context.Context, apiClient LoadBalancerClient, projectId string) (map[string]loadbalancer.CredentialsResponse, error) {
	loadBalancers, err := apiClient.ListLoadBalancersExecute(ctx, projectId)
	if err != nil {
		return nil, fmt.Errorf("list load balancers: %w", err)
	}

	if loadBalancers == nil || loadBalancers.LoadBalancers == nil {
		return nil, fmt.Errorf("no load balancers found")
	}

	var usedCredentials []string

	for _, loadBalancer := range *loadBalancers.LoadBalancers {
		if loadBalancer.Options == nil || loadBalancer.Options.Observability == nil {
			continue
		}

		if loadBalancer.Options != nil && loadBalancer.Options.Observability != nil && loadBalancer.Options.Observability.Logs != nil && loadBalancer.Options.Observability.Metrics != nil {
			usedCredentials = append(usedCredentials, *loadBalancer.Options.Observability.Logs.CredentialsRef)
		}
		if loadBalancer.Options != nil && loadBalancer.Options.Observability != nil && loadBalancer.Options.Observability.Metrics != nil && loadBalancer.Options.Observability.Logs == nil {
			usedCredentials = append(usedCredentials, *loadBalancer.Options.Observability.Metrics.CredentialsRef)
		}
	}

	credentials, err := apiClient.ListCredentialsExecute(ctx, projectId)
	if err != nil {
		return nil, fmt.Errorf("get credentials: %w", err)
	}

	if credentials == nil || credentials.Credentials == nil {
		return nil, fmt.Errorf("no credentials found")
	}
	var usedObsCredentials map[string]loadbalancer.CredentialsResponse

	for _, credential := range *credentials.Credentials {
		if credential.CredentialsRef == nil {
			continue
		}
		ref := *credential.CredentialsRef
		if slices.Contains(usedCredentials, ref) {
			usedObsCredentials[ref] = credential
		}
	}

	return usedObsCredentials, nil
}
