package generatepayload

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/load-balancer/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/loadbalancer"

	"github.com/spf13/cobra"
)

const (
	instanceNameFlag = "instance-name"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	InstanceName *string
}

var (
	defaultPayloadListener = &loadbalancer.Listener{
		DisplayName: utils.Ptr(""),
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

	defaultPayloadNetwork = &loadbalancer.Network{
		NetworkId: utils.Ptr(""),
		Role:      utils.Ptr(""),
	}

	defaultPayloadTargetPool = &loadbalancer.TargetPool{
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

	DefaultCreateLoadBalancerPayload = loadbalancer.CreateLoadBalancerPayload{
		ExternalAddress: utils.Ptr(""),
		Listeners: &[]loadbalancer.Listener{
			*defaultPayloadListener,
		},
		Name: utils.Ptr(""),
		Networks: &[]loadbalancer.Network{
			*defaultPayloadNetwork,
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
		TargetPools: &[]loadbalancer.TargetPool{
			*defaultPayloadTargetPool,
		},
	}
)

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "generate-payload",
		Short: "Generates a payload to create/update a Load Balancer",
		Long: fmt.Sprintf("%s\n%s",
			"Generates a JSON payload with values to be used as --payload input for load balancer creation or update.",
			"See https://docs.api.stackit.cloud/documentation/load-balancer/version/v1#tag/Load-Balancer/operation/APIService_CreateLoadBalancer for information regarding the payload structure.",
		),
		Args: args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Generate a payload, and adapt it with custom values for the different configuration options`,
				`$ stackit load-balancer generate-payload > ./payload.json`,
				`<Modify payload in file, if needed>`,
				`$ stackit load-balancer create --payload @./payload.json`),
			examples.NewExample(
				`Generate a payload with values of an existing load balancer, and adapt it with custom values for the different configuration options`,
				`$ stackit load-balancer generate-payload --instance-name xxx > ./payload.json`,
				`<Modify payload in file>`,
				`$ stackit load-balancer update xxx --payload @./payload.json`),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			model, err := parseInput(p, cmd)
			if err != nil {
				return err
			}

			// Configure API client
			apiClient, err := client.ConfigureClient(p)
			if err != nil {
				return err
			}

			if model.InstanceName == nil {
				createPayload := DefaultCreateLoadBalancerPayload
				return outputCreateResult(p, &createPayload)
			}

			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("read load balancer: %w", err)
			}

			listeners := *resp.Listeners

			for i := range listeners {
				listener := listeners[i]
				listeners[i] = loadbalancer.Listener{
					DisplayName: listener.DisplayName,
					Port:        listener.Port,
					Protocol:    listener.Protocol,
					TargetPool:  listener.TargetPool,
				}

				if listener.ServerNameIndicators != nil {
					listeners[i].ServerNameIndicators = listener.ServerNameIndicators
				}

				if listener.Tcp != nil {
					listeners[i].Tcp = listener.Tcp
				}

				if listener.Udp != nil {
					listeners[i].Udp = listener.Udp
				}
			}

			updatePayload := &loadbalancer.UpdateLoadBalancerPayload{
				ExternalAddress: resp.ExternalAddress,
				Listeners:       &listeners,
				Name:            resp.Name,
				Networks:        resp.Networks,
				Options:         resp.Options,
				TargetPools:     resp.TargetPools,
				Version:         resp.Version,
			}
			return outputUpdateResult(p, updatePayload)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().StringP(instanceNameFlag, "n", "", "If set, generates the payload with the current values of the given load balancer. If unset, generates the payload with empty values")
}

func parseInput(p *print.Printer, cmd *cobra.Command) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)

	instanceName := flags.FlagToStringPointer(p, cmd, instanceNameFlag)
	// If instanceName is provided, projectId is needed as well
	if instanceName != nil && globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	return &inputModel{
		GlobalFlagModel: globalFlags,
		InstanceName:    instanceName,
	}, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *loadbalancer.APIClient) loadbalancer.ApiGetLoadBalancerRequest {
	req := apiClient.GetLoadBalancer(ctx, model.ProjectId, *model.InstanceName)
	return req
}

func outputCreateResult(p *print.Printer, payload *loadbalancer.CreateLoadBalancerPayload) error {
	payloadBytes, err := json.MarshalIndent(*payload, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal create load balancer payload: %w", err)
	}
	p.Outputln(string(payloadBytes))

	return nil
}

func outputUpdateResult(p *print.Printer, payload *loadbalancer.UpdateLoadBalancerPayload) error {
	payloadBytes, err := json.MarshalIndent(*payload, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal update load balancer payload: %w", err)
	}
	p.Outputln(string(payloadBytes))

	return nil
}
