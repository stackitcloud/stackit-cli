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
	lbUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/load-balancer/utils"
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
				`$ stackit load-balancer generate-payload --instance-name my-lb > ./payload.json`,
				`<Modify payload in file>`,
				`$ stackit load-balancer update my-lb --payload @./payload.json`),
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

			var payload *loadbalancer.CreateLoadBalancerPayload
			if model.InstanceName == nil {
				payload = lbUtils.GetDefaultPayload()
			} else {
				req := buildRequest(ctx, model, apiClient)
				resp, err := req.Execute()
				if err != nil {
					return fmt.Errorf("read load balancer: %w", err)
				}
				payload = &loadbalancer.CreateLoadBalancerPayload{
					ExternalAddress: resp.ExternalAddress,
					Listeners:       resp.Listeners,
					Name:            resp.Name,
					Networks:        resp.Networks,
					Options:         resp.Options,
					PrivateAddress:  resp.PrivateAddress,
					TargetPools:     resp.TargetPools,
					Version:         resp.Version,
				}
			}

			return outputResult(p, payload)
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

func outputResult(p *print.Printer, payload *loadbalancer.CreateLoadBalancerPayload) error {
	payloadBytes, err := json.MarshalIndent(*payload, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal payload: %w", err)
	}
	p.Outputln(string(payloadBytes))

	return nil
}
