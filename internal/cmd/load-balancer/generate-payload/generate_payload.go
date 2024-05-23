package generatepayload

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/fileutils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/load-balancer/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/loadbalancer"

	"github.com/spf13/cobra"
)

const (
	loadBalancerNameFlag = "lb-name"
	filePathFlag         = "file-path"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	LoadBalancerName *string
	FilePath         *string
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
				`$ stackit load-balancer generate-payload --file-path ./payload.json`,
				`<Modify payload in file, if needed>`,
				`$ stackit load-balancer create --payload @./payload.json`),
			examples.NewExample(
				`Generate a payload with values of an existing load balancer, and adapt it with custom values for the different configuration options`,
				`$ stackit load-balancer generate-payload --lb-name xxx --file-path ./payload.json`,
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

			if model.LoadBalancerName == nil {
				createPayload := DefaultCreateLoadBalancerPayload
				return outputCreateResult(p, model, &createPayload)
			}

			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("read load balancer: %w", err)
			}

			listeners := modifyListener(resp)

			updatePayload := &loadbalancer.UpdateLoadBalancerPayload{
				ExternalAddress: resp.ExternalAddress,
				Listeners:       listeners,
				Name:            resp.Name,
				Networks:        resp.Networks,
				Options:         resp.Options,
				TargetPools:     resp.TargetPools,
				Version:         resp.Version,
			}
			return outputUpdateResult(p, model, updatePayload)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().StringP(loadBalancerNameFlag, "n", "", "If set, generates the payload with the current values of the given load balancer. If unset, generates the payload with empty values")
	cmd.Flags().StringP(filePathFlag, "f", "", "If set, writes the payload in the given file. If unset, writes the payload to the standard output")
}

func parseInput(p *print.Printer, cmd *cobra.Command) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)

	loadBalancerName := flags.FlagToStringPointer(p, cmd, loadBalancerNameFlag)
	// If load balancer name is provided, projectId is needed as well
	if loadBalancerName != nil && globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel:  globalFlags,
		LoadBalancerName: loadBalancerName,
		FilePath:         flags.FlagToStringPointer(p, cmd, filePathFlag),
	}

	if p.IsVerbosityDebug() {
		modelStr, err := print.BuildDebugStrFromInputModel(model)
		if err != nil {
			p.Debug(print.ErrorLevel, "convert model to string for debugging: %v", err)
		} else {
			p.Debug(print.DebugLevel, "parsed input values: %s", modelStr)
		}
	}

	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *loadbalancer.APIClient) loadbalancer.ApiGetLoadBalancerRequest {
	req := apiClient.GetLoadBalancer(ctx, model.ProjectId, *model.LoadBalancerName)
	return req
}

func outputCreateResult(p *print.Printer, model *inputModel, payload *loadbalancer.CreateLoadBalancerPayload) error {
	payloadBytes, err := json.MarshalIndent(*payload, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal create load balancer payload: %w", err)
	}

	if model.FilePath != nil {
		err = fileutils.FileOutput(*model.FilePath, string(payloadBytes))
		if err != nil {
			return fmt.Errorf("write create load balancer payload to the file: %w", err)
		}
	} else {
		p.Outputln(string(payloadBytes))
	}

	return nil
}

func outputUpdateResult(p *print.Printer, model *inputModel, payload *loadbalancer.UpdateLoadBalancerPayload) error {
	payloadBytes, err := json.MarshalIndent(*payload, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal update load balancer payload: %w", err)
	}

	if model.FilePath != nil {
		err = fileutils.FileOutput(*model.FilePath, string(payloadBytes))
		if err != nil {
			return fmt.Errorf("write update load balancer payload to the file: %w", err)
		}
	} else {
		p.Outputln(string(payloadBytes))
	}

	return nil
}

func modifyListener(resp *loadbalancer.LoadBalancer) *[]loadbalancer.Listener {
	listeners := *resp.Listeners

	for i := range listeners {
		listeners[i].Name = nil
	}

	return &listeners
}
