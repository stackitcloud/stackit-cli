package utils

import (
	"context"
	"fmt"
	"slices"
	"sort"

	"github.com/stackitcloud/stackit-sdk-go/services/loadbalancer"
)

const (
	OP_FILTER_NOP = iota
	OP_FILTER_USED
	OP_FILTER_UNUSED
)

// enforce implementation of interfaces
var (
	_ LoadBalancerClient = &loadbalancer.APIClient{}
)

type LoadBalancerClient interface {
	GetCredentialsExecute(ctx context.Context, projectId, region, credentialsRef string) (*loadbalancer.GetCredentialsResponse, error)
	GetLoadBalancerExecute(ctx context.Context, projectId, region, name string) (*loadbalancer.LoadBalancer, error)
	UpdateTargetPool(ctx context.Context, projectId, region, loadBalancerName, targetPoolName string) loadbalancer.ApiUpdateTargetPoolRequest
	ListLoadBalancersExecute(ctx context.Context, projectId, region string) (*loadbalancer.ListLoadBalancersResponse, error)
}

func GetCredentialsDisplayName(ctx context.Context, apiClient LoadBalancerClient, projectId, region, credentialsRef string) (string, error) {
	resp, err := apiClient.GetCredentialsExecute(ctx, projectId, region, credentialsRef)
	if err != nil {
		return "", fmt.Errorf("get Load Balancer credentials: %w", err)
	}
	return *resp.Credential.DisplayName, nil
}

func GetLoadBalancerTargetPool(ctx context.Context, apiClient LoadBalancerClient, projectId, region, loadBalancerName, targetPoolName string) (*loadbalancer.TargetPool, error) {
	resp, err := apiClient.GetLoadBalancerExecute(ctx, projectId, region, loadBalancerName)
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

func GetTargetName(ctx context.Context, apiClient LoadBalancerClient, projectId, region, loadBalancerName, targetPoolName, targetIp string) (string, error) {
	targetPool, err := GetLoadBalancerTargetPool(ctx, apiClient, projectId, region, loadBalancerName, targetPoolName)
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
// It goes through all load balancers and checks what observability credentials are being used, then returns a list of those credentials.
func GetUsedObsCredentials(ctx context.Context, apiClient LoadBalancerClient, allCredentials []loadbalancer.CredentialsResponse, projectId, region string) ([]loadbalancer.CredentialsResponse, error) {
	var usedCredentialsSlice []loadbalancer.CredentialsResponse

	loadBalancers, err := apiClient.ListLoadBalancersExecute(ctx, projectId, region)
	if err != nil {
		return nil, fmt.Errorf("list load balancers: %w", err)
	}
	if loadBalancers == nil || loadBalancers.LoadBalancers == nil {
		return usedCredentialsSlice, nil
	}

	var usedCredentialsRefs []string
	for i := range *loadBalancers.LoadBalancers {
		loadBalancer := &(*loadBalancers.LoadBalancers)[i]

		if loadBalancer.Options == nil || loadBalancer.Options.Observability == nil {
			continue
		}

		if loadBalancer.Options != nil && loadBalancer.Options.Observability != nil && loadBalancer.Options.Observability.Logs != nil && loadBalancer.Options.Observability.Logs.CredentialsRef != nil {
			usedCredentialsRefs = append(usedCredentialsRefs, *loadBalancer.Options.Observability.Logs.CredentialsRef)
		}
		if loadBalancer.Options != nil && loadBalancer.Options.Observability != nil && loadBalancer.Options.Observability.Metrics != nil && loadBalancer.Options.Observability.Metrics.CredentialsRef != nil {
			usedCredentialsRefs = append(usedCredentialsRefs, *loadBalancer.Options.Observability.Metrics.CredentialsRef)
		}
	}

	usedCredentialsMap := make(map[string]loadbalancer.CredentialsResponse)
	for _, credential := range allCredentials {
		if credential.CredentialsRef == nil {
			continue
		}
		ref := *credential.CredentialsRef
		if slices.Contains(usedCredentialsRefs, ref) {
			usedCredentialsMap[ref] = credential
		}
	}

	for _, credential := range usedCredentialsMap {
		usedCredentialsSlice = append(usedCredentialsSlice, credential)
	}

	// sort credentials by reference to make output deterministic
	sort.Slice(usedCredentialsSlice, func(i, j int) bool {
		return *usedCredentialsSlice[i].CredentialsRef < *usedCredentialsSlice[j].CredentialsRef
	})

	return usedCredentialsSlice, nil
}

// GetUnusedObsCredentials returns a list of credentials that are not used by any load balancer for observability metrics or logs.
// It compares the list of all credentials with the list of used credentials and returns a list of credentials that are not used.
func GetUnusedObsCredentials(usedCredentials, allCredentials []loadbalancer.CredentialsResponse) []loadbalancer.CredentialsResponse {
	var unusedCredentials []loadbalancer.CredentialsResponse
	usedCredentialsRefs := make(map[string]bool)
	for _, credential := range usedCredentials {
		if credential.CredentialsRef != nil {
			usedCredentialsRefs[*credential.CredentialsRef] = true
		}
	}

	for _, credential := range allCredentials {
		if credential.CredentialsRef == nil {
			continue
		}
		if !usedCredentialsRefs[*credential.CredentialsRef] {
			unusedCredentials = append(unusedCredentials, credential)
		}
	}

	return unusedCredentials
}

// FilterCredentials filters a list of credentials based on the used and unused flags.
// If used is true, it returns only the credentials that are used by load balancers for observability metrics or logs.
// If unused is true, it returns only the credentials that are not used by any load balancer for observability metrics or logs.
// If both used and unused are true, it returns an error.
// If both used and unused are false, it returns the original list of credentials.
func FilterCredentials(ctx context.Context, client LoadBalancerClient, allCredentials []loadbalancer.CredentialsResponse, projectId, region string, filterOp int) ([]loadbalancer.CredentialsResponse, error) {
	// check that filter OP is valid
	if filterOp != OP_FILTER_USED && filterOp != OP_FILTER_UNUSED && filterOp != OP_FILTER_NOP {
		return nil, fmt.Errorf("invalid filter operation")
	}

	if filterOp == OP_FILTER_NOP {
		return allCredentials, nil
	}

	usedCredentials, err := GetUsedObsCredentials(ctx, client, allCredentials, projectId, region)
	if err != nil {
		return nil, fmt.Errorf("get used observability credentials: %w", err)
	}

	if filterOp == OP_FILTER_UNUSED {
		return GetUnusedObsCredentials(usedCredentials, allCredentials), nil
	}
	return usedCredentials, nil
}
