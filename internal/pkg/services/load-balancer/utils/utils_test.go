package utils

import (
	"testing"

	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/google/go-cmp/cmp"
	"github.com/stackitcloud/stackit-sdk-go/services/loadbalancer"
)

func fixtureGetDefaultPayload(mods ...func(payload *loadbalancer.CreateLoadBalancerPayload)) *loadbalancer.CreateLoadBalancerPayload {
	payload := &loadbalancer.CreateLoadBalancerPayload{
		ExternalAddress: utils.Ptr(""),

		Listeners: &[]loadbalancer.Listener{
			{
				DisplayName: utils.Ptr(""),
				Name:        utils.Ptr(""),
				Port:        utils.Ptr(int64(0)),
				Protocol:    utils.Ptr(""),
				ServerNameIndicators: &[]loadbalancer.ServerNameIndicator{
					{
						Name: utils.Ptr(""),
					},
				},
				TargetPool: utils.Ptr(""),
				Tcp: &loadbalancer.OptionsTCP{
					IdleTimeout: utils.Ptr(""),
				},
				Udp: &loadbalancer.OptionsUDP{
					IdleTimeout: utils.Ptr(""),
				},
			},
		},
		Name: utils.Ptr(""),
		Networks: &[]loadbalancer.Network{
			{
				NetworkId: utils.Ptr(""),
				Role:      utils.Ptr(""),
			},
		},
		Options: &loadbalancer.LoadBalancerOptions{
			AccessControl: &loadbalancer.LoadbalancerOptionAccessControl{
				AllowedSourceRanges: &[]string{
					"",
				},
			},
			EphemeralAddress: utils.Ptr(false),
			Observability: &loadbalancer.LoadbalancerOptionObservability{
				Logs: &loadbalancer.LoadbalancerOptionLogs{
					CredentialsRef: utils.Ptr(""),
					PushUrl:        utils.Ptr(""),
				},
				Metrics: &loadbalancer.LoadbalancerOptionMetrics{
					CredentialsRef: utils.Ptr(""),
					PushUrl:        utils.Ptr(""),
				},
			},
			PrivateNetworkOnly: utils.Ptr(false),
		},
		PrivateAddress: utils.Ptr(""),
		TargetPools: &[]loadbalancer.TargetPool{
			{
				ActiveHealthCheck: &loadbalancer.ActiveHealthCheck{
					HealthyThreshold:   utils.Ptr(int64(0)),
					Interval:           utils.Ptr(""),
					IntervalJitter:     utils.Ptr(""),
					Timeout:            utils.Ptr(""),
					UnhealthyThreshold: utils.Ptr(int64(0)),
				},
				Name: utils.Ptr(""),
				SessionPersistence: &loadbalancer.SessionPersistence{
					UseSourceIpAddress: utils.Ptr(false),
				},
				TargetPort: utils.Ptr(int64(0)),
				Targets: &[]loadbalancer.Target{
					{
						DisplayName: utils.Ptr(""),
						Ip:          utils.Ptr(""),
					},
				},
			},
		},
	}
	for _, mod := range mods {
		mod(payload)
	}
	return payload
}

func TestGetDefaultPayload(t *testing.T) {
	tests := []struct {
		description    string
		isValid        bool
		expectedOutput *loadbalancer.CreateLoadBalancerPayload
	}{
		{
			description:    "base",
			isValid:        true,
			expectedOutput: fixtureGetDefaultPayload(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			output := GetDefaultPayload()

			if !tt.isValid {
				t.Errorf("did not fail on invalid input")
			}
			if !tt.isValid {
				return
			}
			diff := cmp.Diff(output, tt.expectedOutput)
			if diff != "" {
				t.Fatalf("Output is not as expected: %s", diff)
			}
		})
	}
}
