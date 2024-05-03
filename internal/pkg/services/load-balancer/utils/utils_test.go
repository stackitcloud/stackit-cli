package utils

import (
	"context"
	"fmt"
	"testing"

	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/stackitcloud/stackit-sdk-go/services/loadbalancer"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
)

var (
	testProjectId = uuid.NewString()
)

const (
	testCredentialsRef         = "credentials-ref"
	testCredentialsDisplayName = "credentials-name"
	testLoadBalancerName       = "my-load-balancer"
)

type loadBalancerClientMocked struct {
	getCredentialsFails  bool
	getCredentialsResp   *loadbalancer.GetCredentialsResponse
	getLoadBalancerFails bool
	getLoadBalancerResp  *loadbalancer.LoadBalancer
}

func (m *loadBalancerClientMocked) GetCredentialsExecute(_ context.Context, _, _ string) (*loadbalancer.GetCredentialsResponse, error) {
	if m.getCredentialsFails {
		return nil, fmt.Errorf("could not get credentials")
	}
	return m.getCredentialsResp, nil
}

func (m *loadBalancerClientMocked) GetLoadBalancerExecute(_ context.Context, _, _ string) (*loadbalancer.LoadBalancer, error) {
	if m.getLoadBalancerFails {
		return nil, fmt.Errorf("could not get load balancer")
	}
	return m.getLoadBalancerResp, nil
}

func (m *loadBalancerClientMocked) UpdateTargetPool(_ context.Context, _, _, _ string) loadbalancer.ApiUpdateTargetPoolRequest {
	return loadbalancer.ApiUpdateTargetPoolRequest{}
}

func TestGetCredentialsDisplayName(t *testing.T) {
	tests := []struct {
		description         string
		getCredentialsFails bool
		getCredentialsResp  *loadbalancer.GetCredentialsResponse
		isValid             bool
		expectedOutput      string
	}{
		{
			description: "base",
			getCredentialsResp: &loadbalancer.GetCredentialsResponse{
				Credential: &loadbalancer.CredentialsResponse{
					DisplayName: utils.Ptr(testCredentialsDisplayName),
				},
			},
			isValid:        true,
			expectedOutput: testCredentialsDisplayName,
		},
		{
			description:         "get credentials fails",
			getCredentialsFails: true,
			isValid:             false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			client := &loadBalancerClientMocked{
				getCredentialsFails: tt.getCredentialsFails,
				getCredentialsResp:  tt.getCredentialsResp,
			}

			output, err := GetCredentialsDisplayName(context.Background(), client, testProjectId, testCredentialsRef)

			if tt.isValid && err != nil {
				t.Errorf("failed on valid input")
			}
			if !tt.isValid && err == nil {
				t.Errorf("did not fail on invalid input")
			}
			if !tt.isValid {
				return
			}
			if output != tt.expectedOutput {
				t.Errorf("expected output to be %s, got %s", tt.expectedOutput, output)
			}
		})
	}
}

func fixtureLoadBalancer(mods ...func(*loadbalancer.LoadBalancer)) *loadbalancer.LoadBalancer {
	lb := loadbalancer.LoadBalancer{
		Name: utils.Ptr(testLoadBalancerName),
		TargetPools: &[]loadbalancer.TargetPool{
			{
				Name: utils.Ptr("target-pool-1"),
				Targets: &[]loadbalancer.Target{
					{
						DisplayName: utils.Ptr("target-1"),
						Ip:          utils.Ptr("1.2.3.4/32"),
					},
					{
						DisplayName: utils.Ptr("target-2"),
						Ip:          utils.Ptr("4.3.2.1/32"),
					},
				},
			},
			{
				Name: utils.Ptr("target-pool-2"),
				Targets: &[]loadbalancer.Target{
					{
						DisplayName: utils.Ptr("target-1"),
						Ip:          utils.Ptr("6.7.8.9/32"),
					},
					{
						DisplayName: utils.Ptr("target-2"),
						Ip:          utils.Ptr("9.8.7.6/32"),
					},
				},
			},
		},
	}

	for _, mod := range mods {
		mod(&lb)
	}
	return &lb
}

func TestGetLoadBalancerTargetPool(t *testing.T) {
	tests := []struct {
		description          string
		targetPoolName       string
		getLoadBalancerFails bool
		getLoadBalancerResp  *loadbalancer.LoadBalancer
		isValid              bool
		expectedOutput       *loadbalancer.TargetPool
	}{
		{
			description:         "base",
			targetPoolName:      "target-pool-1",
			getLoadBalancerResp: fixtureLoadBalancer(),
			isValid:             true,
			expectedOutput: &loadbalancer.TargetPool{
				Name: utils.Ptr("target-pool-1"),
				Targets: &[]loadbalancer.Target{
					{
						DisplayName: utils.Ptr("target-1"),
						Ip:          utils.Ptr("1.2.3.4/32"),
					},
					{
						DisplayName: utils.Ptr("target-2"),
						Ip:          utils.Ptr("4.3.2.1/32"),
					},
				},
			},
		},
		{
			description:         "target pool not found",
			targetPoolName:      "target-pool-non-existent",
			getLoadBalancerResp: fixtureLoadBalancer(),
			isValid:             false,
		},
		{
			description: "no target pools",
			getLoadBalancerResp: fixtureLoadBalancer(func(lb *loadbalancer.LoadBalancer) {
				lb.TargetPools = nil
			}),
			isValid: false,
		},
		{
			description:          "get load balancer fails",
			getLoadBalancerFails: true,
			isValid:              false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			client := &loadBalancerClientMocked{
				getLoadBalancerFails: tt.getLoadBalancerFails,
				getLoadBalancerResp:  tt.getLoadBalancerResp,
			}

			output, err := GetLoadBalancerTargetPool(context.Background(), client, testProjectId, testLoadBalancerName, tt.targetPoolName)

			if tt.isValid && err != nil {
				t.Errorf("failed on valid input")
			}
			if !tt.isValid && err == nil {
				t.Errorf("did not fail on invalid input")
			}
			if !tt.isValid {
				return
			}
			diff := cmp.Diff(output, tt.expectedOutput)
			if diff != "" {
				t.Fatalf("Data does not match: %s", diff)
			}
		})
	}
}

func TestAddTargetToTargetPool(t *testing.T) {
	tests := []struct {
		description        string
		targetPool         *loadbalancer.TargetPool
		target             *loadbalancer.Target
		isValid            bool
		expectedTargetPool *loadbalancer.TargetPool
	}{
		{
			description: "base",
			targetPool: &loadbalancer.TargetPool{
				Name: utils.Ptr("target-pool-1"),
				Targets: &[]loadbalancer.Target{
					{
						DisplayName: utils.Ptr("target-1"),
						Ip:          utils.Ptr("1.2.3.4/32"),
					},
				},
			},
			target: &loadbalancer.Target{
				DisplayName: utils.Ptr("target-2"),
				Ip:          utils.Ptr("6.6.6.6/32"),
			},
			isValid: true,
			expectedTargetPool: &loadbalancer.TargetPool{
				Name: utils.Ptr("target-pool-1"),
				Targets: &[]loadbalancer.Target{
					{
						DisplayName: utils.Ptr("target-1"),
						Ip:          utils.Ptr("1.2.3.4/32"),
					},
					{
						DisplayName: utils.Ptr("target-2"),
						Ip:          utils.Ptr("6.6.6.6/32"),
					},
				},
			},
		},
		{
			description: "nil target pool",
			targetPool:  nil,
			target: &loadbalancer.Target{
				DisplayName: utils.Ptr("target-3"),
				Ip:          utils.Ptr("2.2.2.2/32"),
			},
			expectedTargetPool: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			err := AddTargetToTargetPool(tt.targetPool, tt.target)

			if tt.isValid && err != nil {
				t.Errorf("failed on valid input")
			}
			if !tt.isValid && err == nil {
				t.Errorf("did not fail on invalid input")
			}
			if !tt.isValid {
				return
			}
			diff := cmp.Diff(tt.targetPool, tt.expectedTargetPool)
			if diff != "" {
				t.Fatalf("Data does not match: %s", diff)
			}
		})
	}
}

func TestToPayloadTargetPool(t *testing.T) {
	tests := []struct {
		description string
		input       *loadbalancer.TargetPool
		expected    *loadbalancer.UpdateTargetPoolPayload
	}{
		{
			description: "base",
			input: &loadbalancer.TargetPool{
				Name: utils.Ptr("target-pool-1"),
				ActiveHealthCheck: &loadbalancer.ActiveHealthCheck{
					UnhealthyThreshold: utils.Ptr(int64(3)),
				},
				SessionPersistence: &loadbalancer.SessionPersistence{
					UseSourceIpAddress: utils.Ptr(true),
				},
				TargetPort: utils.Ptr(int64(80)),
				Targets: &[]loadbalancer.Target{
					{
						DisplayName: utils.Ptr("target-1"),
						Ip:          utils.Ptr("1.2.3.4/32"),
					},
				},
			},
			expected: &loadbalancer.UpdateTargetPoolPayload{
				Name: utils.Ptr("target-pool-1"),
				ActiveHealthCheck: &loadbalancer.ActiveHealthCheck{
					UnhealthyThreshold: utils.Ptr(int64(3)),
				},
				SessionPersistence: &loadbalancer.SessionPersistence{
					UseSourceIpAddress: utils.Ptr(true),
				},
				TargetPort: utils.Ptr(int64(80)),
				Targets: &[]loadbalancer.Target{
					{
						DisplayName: utils.Ptr("target-1"),
						Ip:          utils.Ptr("1.2.3.4/32"),
					},
				},
			},
		},
		{
			description: "nil target pool",
			input:       nil,
			expected:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			output := ToPayloadTargetPool(tt.input)

			diff := cmp.Diff(output, tt.expected)
			if diff != "" {
				t.Errorf("expected output to be %+v, got %+v", tt.expected, output)
			}
		})
	}
}
