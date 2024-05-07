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

func fixtureLoadBalancer(mods ...func(*loadbalancer.LoadBalancer)) *loadbalancer.LoadBalancer {
	lb := loadbalancer.LoadBalancer{
		Name: utils.Ptr(testLoadBalancerName),
		TargetPools: &[]loadbalancer.TargetPool{
			{
				Name: utils.Ptr("target-pool-1"),
				Targets: &[]loadbalancer.Target{
					{
						DisplayName: utils.Ptr("target-1"),
						Ip:          utils.Ptr("1.2.3.4"),
					},
					{
						DisplayName: utils.Ptr("target-2"),
						Ip:          utils.Ptr("4.3.2.1"),
					},
				},
			},
			{
				Name: utils.Ptr("target-pool-2"),
				Targets: &[]loadbalancer.Target{
					{
						DisplayName: utils.Ptr("target-1"),
						Ip:          utils.Ptr("6.7.8.9"),
					},
					{
						DisplayName: utils.Ptr("target-2"),
						Ip:          utils.Ptr("9.8.7.6"),
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

func fixtureTargets(mod ...func(*[]loadbalancer.Target)) *[]loadbalancer.Target {
	targets := &[]loadbalancer.Target{
		{
			DisplayName: utils.Ptr("target-1"),
			Ip:          utils.Ptr("1.2.3.4"),
		},
		{
			DisplayName: utils.Ptr("target-2"),
			Ip:          utils.Ptr("2.2.2.2"),
		},
		{
			DisplayName: utils.Ptr("target-3"),
			Ip:          utils.Ptr("6.6.6.6"),
		},
	}

	for _, m := range mod {
		m(targets)
	}

	return targets
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
						Ip:          utils.Ptr("1.2.3.4"),
					},
					{
						DisplayName: utils.Ptr("target-2"),
						Ip:          utils.Ptr("4.3.2.1"),
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
				lb.TargetPools = &[]loadbalancer.TargetPool{}
			}),
			isValid: false,
		},
		{
			description: "nil target pools",
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

func TestFindLoadBalancerTargetPoolByName(t *testing.T) {
	tests := []struct {
		description        string
		targetPools        []loadbalancer.TargetPool
		targetPoolName     string
		expectedTargetPool *loadbalancer.TargetPool
	}{
		{
			description: "base",
			targetPools: []loadbalancer.TargetPool{
				{
					Name: utils.Ptr("target-pool-1"),
				},
				{
					Name: utils.Ptr("target-pool-2"),
				},
			},
			targetPoolName: "target-pool-1",
			expectedTargetPool: &loadbalancer.TargetPool{
				Name: utils.Ptr("target-pool-1"),
			},
		},
		{
			description: "target pool not found",
			targetPools: []loadbalancer.TargetPool{
				{
					Name: utils.Ptr("target-pool-1"),
				},
				{
					Name: utils.Ptr("target-pool-2"),
				},
			},
			targetPoolName:     "target-pool-3",
			expectedTargetPool: nil,
		},
		{
			description:        "nil target pools",
			targetPools:        nil,
			targetPoolName:     "target-pool-1",
			expectedTargetPool: nil,
		},
		{
			description:        "no target pools",
			targetPools:        []loadbalancer.TargetPool{},
			targetPoolName:     "target-pool-1",
			expectedTargetPool: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			output := FindLoadBalancerTargetPoolByName(tt.targetPools, tt.targetPoolName)

			diff := cmp.Diff(output, tt.expectedTargetPool)
			if diff != "" {
				t.Fatalf("Data does not match: %s", diff)
			}
		})
	}
}

func TestFindLoadBalancerListenerByTargetPool(t *testing.T) {
	tests := []struct {
		description    string
		listeners      []loadbalancer.Listener
		targetPoolName string
		expected       *loadbalancer.Listener
	}{
		{
			description: "base",
			listeners: []loadbalancer.Listener{
				{
					TargetPool: utils.Ptr("target-pool-1"),
				},
				{
					TargetPool: utils.Ptr("target-pool-2"),
				},
			},
			targetPoolName: "target-pool-1",
			expected: &loadbalancer.Listener{
				TargetPool: utils.Ptr("target-pool-1"),
			},
		},
		{
			description: "listener not found",
			listeners: []loadbalancer.Listener{
				{
					TargetPool: utils.Ptr("target-pool-1"),
				},
				{
					TargetPool: utils.Ptr("target-pool-2"),
				},
			},
			targetPoolName: "target-pool-3",
			expected:       nil,
		},
		{
			description:    "nil listeners",
			listeners:      nil,
			targetPoolName: "target-pool-1",
			expected:       nil,
		},
		{
			description:    "no listeners",
			listeners:      []loadbalancer.Listener{},
			targetPoolName: "target-pool-1",
			expected:       nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			output := FindLoadBalancerListenerByTargetPool(tt.listeners, tt.targetPoolName)

			diff := cmp.Diff(output, tt.expected)
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
						Ip:          utils.Ptr("1.2.3.4"),
					},
				},
			},
			target: &loadbalancer.Target{
				DisplayName: utils.Ptr("target-2"),
				Ip:          utils.Ptr("6.6.6.6"),
			},
			isValid: true,
			expectedTargetPool: &loadbalancer.TargetPool{
				Name: utils.Ptr("target-pool-1"),
				Targets: &[]loadbalancer.Target{
					{
						DisplayName: utils.Ptr("target-1"),
						Ip:          utils.Ptr("1.2.3.4"),
					},
					{
						DisplayName: utils.Ptr("target-2"),
						Ip:          utils.Ptr("6.6.6.6"),
					},
				},
			},
		},
		{
			description: "no target pool targets",
			targetPool: &loadbalancer.TargetPool{
				Name:    utils.Ptr("target-pool-1"),
				Targets: &[]loadbalancer.Target{},
			},
			target: &loadbalancer.Target{
				DisplayName: utils.Ptr("target-3"),
				Ip:          utils.Ptr("2.2.2.2"),
			},
			isValid: true,
			expectedTargetPool: &loadbalancer.TargetPool{
				Name: utils.Ptr("target-pool-1"),
				Targets: &[]loadbalancer.Target{
					{
						DisplayName: utils.Ptr("target-3"),
						Ip:          utils.Ptr("2.2.2.2"),
					},
				},
			},
		},
		{
			description: "nil target pool targets",
			targetPool: &loadbalancer.TargetPool{
				Name:    utils.Ptr("target-pool-1"),
				Targets: nil,
			},
			target: &loadbalancer.Target{
				DisplayName: utils.Ptr("target-3"),
				Ip:          utils.Ptr("2.2.2.2"),
			},
			isValid: true,
			expectedTargetPool: &loadbalancer.TargetPool{
				Name: utils.Ptr("target-pool-1"),
				Targets: &[]loadbalancer.Target{
					{
						DisplayName: utils.Ptr("target-3"),
						Ip:          utils.Ptr("2.2.2.2"),
					},
				},
			},
		},
		{
			description: "nil target pool",
			targetPool:  nil,
			target: &loadbalancer.Target{
				DisplayName: utils.Ptr("target-3"),
				Ip:          utils.Ptr("2.2.2.2"),
			},
			expectedTargetPool: nil,
		},
		{
			description: "nil new target",
			targetPool: &loadbalancer.TargetPool{
				Name: utils.Ptr("target-pool-1"),
				Targets: &[]loadbalancer.Target{
					{
						DisplayName: utils.Ptr("target-1"),
						Ip:          utils.Ptr("1.2.3.4"),
					},
				},
			},
			target:  nil,
			isValid: false,
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

func TestRemoveTargetFromTargetPool(t *testing.T) {
	tests := []struct {
		description        string
		targetPool         *loadbalancer.TargetPool
		targetIp           string
		isValid            bool
		expectedTargetPool *loadbalancer.TargetPool
	}{
		{
			description: "remove first target",
			targetPool: &loadbalancer.TargetPool{
				Name:    utils.Ptr("target-pool-1"),
				Targets: fixtureTargets(),
			},
			targetIp: "1.2.3.4",
			isValid:  true,
			expectedTargetPool: &loadbalancer.TargetPool{
				Name: utils.Ptr("target-pool-1"),
				Targets: &[]loadbalancer.Target{
					{
						DisplayName: utils.Ptr("target-2"),
						Ip:          utils.Ptr("2.2.2.2"),
					},
					{
						DisplayName: utils.Ptr("target-3"),
						Ip:          utils.Ptr("6.6.6.6"),
					},
				},
			},
		},
		{
			description: "remove last target",
			targetPool: &loadbalancer.TargetPool{
				Name:    utils.Ptr("target-pool-1"),
				Targets: fixtureTargets(),
			},
			targetIp: "6.6.6.6",
			isValid:  true,
			expectedTargetPool: &loadbalancer.TargetPool{
				Name: utils.Ptr("target-pool-1"),
				Targets: &[]loadbalancer.Target{
					{
						DisplayName: utils.Ptr("target-1"),
						Ip:          utils.Ptr("1.2.3.4"),
					},
					{
						DisplayName: utils.Ptr("target-2"),
						Ip:          utils.Ptr("2.2.2.2"),
					},
				},
			},
		},
		{
			description: "remove middle target",
			targetPool: &loadbalancer.TargetPool{
				Name:    utils.Ptr("target-pool-1"),
				Targets: fixtureTargets(),
			},
			targetIp: "2.2.2.2",
			isValid:  true,
			expectedTargetPool: &loadbalancer.TargetPool{
				Name: utils.Ptr("target-pool-1"),
				Targets: &[]loadbalancer.Target{
					{
						DisplayName: utils.Ptr("target-1"),
						Ip:          utils.Ptr("1.2.3.4"),
					},
					{
						DisplayName: utils.Ptr("target-3"),
						Ip:          utils.Ptr("6.6.6.6"),
					},
				},
			},
		},
		{
			description: "remove only target",
			targetPool: &loadbalancer.TargetPool{
				Name: utils.Ptr("target-pool-1"),
				Targets: &[]loadbalancer.Target{
					{
						DisplayName: utils.Ptr("target-1"),
						Ip:          utils.Ptr("1.2.3.4"),
					},
				},
			},
			targetIp: "1.2.3.4",
			isValid:  true,
			expectedTargetPool: &loadbalancer.TargetPool{
				Name:    utils.Ptr("target-pool-1"),
				Targets: &[]loadbalancer.Target{},
			},
		},
		{
			description: "no target pool targets",
			targetPool: &loadbalancer.TargetPool{
				Name:    utils.Ptr("target-pool-1"),
				Targets: &[]loadbalancer.Target{},
			},
			targetIp: "2.2.2.2",
			isValid:  false,
		},
		{
			description: "nil target pool targets",
			targetPool: &loadbalancer.TargetPool{
				Name:    utils.Ptr("target-pool-1"),
				Targets: nil,
			},
			targetIp: "2.2.2.2",
			isValid:  false,
		},
		{
			description:        "nil target pool",
			targetPool:         nil,
			targetIp:           "2.2.2.2",
			expectedTargetPool: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			err := RemoveTargetFromTargetPool(tt.targetPool, tt.targetIp)

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
						Ip:          utils.Ptr("1.2.3.4"),
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
						Ip:          utils.Ptr("1.2.3.4"),
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

func TestGetTargetName(t *testing.T) {
	tests := []struct {
		description          string
		targetPoolName       string
		targetIp             string
		getLoadBalancerFails bool
		getLoadBalancerResp  *loadbalancer.LoadBalancer
		isValid              bool
		expectedOutput       string
	}{
		{
			description:         "base",
			targetPoolName:      "target-pool-1",
			targetIp:            "1.2.3.4",
			getLoadBalancerResp: fixtureLoadBalancer(),
			isValid:             true,
			expectedOutput:      "target-1",
		},
		{
			description:         "target not found",
			targetPoolName:      "target-pool-1",
			targetIp:            "9.9.9.9",
			getLoadBalancerResp: fixtureLoadBalancer(),
			isValid:             false,
		},
		{
			description:    "no targets",
			targetPoolName: "target-pool-1",
			targetIp:       "1.2.3.4",
			getLoadBalancerResp: fixtureLoadBalancer(func(lb *loadbalancer.LoadBalancer) {
				lb.TargetPools = &[]loadbalancer.TargetPool{
					{
						Name:    utils.Ptr("target-pool-1"),
						Targets: &[]loadbalancer.Target{},
					},
				}
			}),
			isValid: false,
		},
		{
			description:    "nil targets",
			targetPoolName: "target-pool-1",
			targetIp:       "1.2.3.4",
			getLoadBalancerResp: fixtureLoadBalancer(func(lb *loadbalancer.LoadBalancer) {
				lb.TargetPools = &[]loadbalancer.TargetPool{
					{
						Name:    utils.Ptr("target-pool-1"),
						Targets: nil,
					},
				}
			}),
			isValid: false,
		},
		{
			description:    "nil target name",
			targetPoolName: "target-pool-1",
			targetIp:       "1.2.3.4",
			getLoadBalancerResp: fixtureLoadBalancer(
				func(lb *loadbalancer.LoadBalancer) {
					lb.TargetPools = &[]loadbalancer.TargetPool{
						{
							Name: utils.Ptr("target-pool-1"),
							Targets: &[]loadbalancer.Target{
								{
									DisplayName: nil,
									Ip:          utils.Ptr("1.2.3.4"),
								},
							},
						},
					}
				}),
			isValid: false,
		},
		{
			description:          "get target pool fails",
			targetPoolName:       "target-pool-1",
			targetIp:             "1.2.3.4",
			getLoadBalancerFails: true,
			isValid:              false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			client := &loadBalancerClientMocked{
				getLoadBalancerResp: tt.getLoadBalancerResp,
			}

			output, err := GetTargetName(context.Background(), client, testProjectId, testLoadBalancerName, tt.targetPoolName, tt.targetIp)

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
