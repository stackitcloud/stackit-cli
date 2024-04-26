package utils

import (
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/loadbalancer"
)

type LoadBalancerClient interface {
}

func GetDefaultPayload() *loadbalancer.CreateLoadBalancerPayload {
	payloadListener := getDefaultPayloadListener()
	payloadNetwork := getDefaultPayloadNetwork()
	payloadTargetPool := getDefaultPayloadTargetPool()

	payload := &loadbalancer.CreateLoadBalancerPayload{
		ExternalAddress: utils.Ptr(""),
		Listeners: &[]loadbalancer.Listener{
			*payloadListener,
		},
		Name: utils.Ptr(""),
		Networks: &[]loadbalancer.Network{
			*payloadNetwork,
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
			*payloadTargetPool,
		},
	}
	return payload
}

func getDefaultPayloadListener() *loadbalancer.Listener {
	output := &loadbalancer.Listener{
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
	}
	return output
}

func getDefaultPayloadNetwork() *loadbalancer.Network {
	output := &loadbalancer.Network{
		NetworkId: utils.Ptr(""),
		Role:      utils.Ptr(""),
	}
	return output
}

func getDefaultPayloadTargetPool() *loadbalancer.TargetPool {
	output := &loadbalancer.TargetPool{
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
	}
	return output
}
