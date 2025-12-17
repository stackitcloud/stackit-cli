package update

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/types"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/load-balancer/client"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/loadbalancer"
)

const (
	loadBalancerNameArg = "LOAD_BALANCER_NAME"
	payloadFlag         = "payload"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	LoadBalancerName string
	Payload          loadbalancer.UpdateLoadBalancerPayload
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("update %s", loadBalancerNameArg),
		Short: "Updates a Load Balancer",
		Long: fmt.Sprintf("%s\n%s\n%s",
			"Updates a load balancer.",
			"The payload can be provided as a JSON string or a file path prefixed with \"@\".",
			"See https://docs.api.stackit.cloud/documentation/load-balancer/version/v1#tag/Load-Balancer/operation/APIService_UpdateLoadBalancer for information regarding the payload structure.",
		),
		Args: args.SingleArg(loadBalancerNameArg, nil),
		Example: examples.Build(
			examples.NewExample(
				`Update a load balancer with name "xxx", using an API payload sourced from the file "./payload.json"`,
				"$ stackit load-balancer update xxx --payload @./payload.json"),
			examples.NewExample(
				`Update a load balancer with name "xxx", using an API payload provided as a JSON string`,
				`$ stackit load-balancer update xxx --payload "{...}"`),
			examples.NewExample(
				`Generate a payload with the current values of an existing load balancer xxx, and adapt it with custom values for the different configuration options`,
				`$ stackit load-balancer generate-payload --lb-name xxx > ./payload.json`,
				`<Modify payload in file>`,
				`$ stackit load-balancer update xxx --payload @./payload.json`),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			model, err := parseInput(params.Printer, cmd, args)
			if err != nil {
				return err
			}

			// Configure API client
			apiClient, err := client.ConfigureClient(params.Printer, params.CliVersion)
			if err != nil {
				return err
			}

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to update load balancer %q?", model.LoadBalancerName)
				err = params.Printer.PromptForConfirmation(prompt)
				if err != nil {
					return err
				}
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			_, err = req.Execute()
			if err != nil {
				return fmt.Errorf("update load balancer: %w", err)
			}

			// The API has no status to wait on, so async mode is default
			params.Printer.Info("Updated load balancer with name %q\n", model.LoadBalancerName)
			return nil
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.ReadFromFileFlag(), payloadFlag, `Request payload (JSON). Can be a string or a file path, if prefixed with "@". Example: @./payload.json`)

	err := flags.MarkFlagsRequired(cmd, payloadFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	loadBalancerName := inputArgs[0]

	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	payloadString := flags.FlagToStringValue(p, cmd, payloadFlag)
	var payload loadbalancer.UpdateLoadBalancerPayload
	err := json.Unmarshal([]byte(payloadString), &payload)
	if err != nil {
		return nil, fmt.Errorf("encode payload: %w", err)
	}

	model := inputModel{
		GlobalFlagModel:  globalFlags,
		LoadBalancerName: loadBalancerName,
		Payload:          payload,
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *loadbalancer.APIClient) loadbalancer.ApiUpdateLoadBalancerRequest {
	req := apiClient.UpdateLoadBalancer(ctx, model.ProjectId, model.Region, model.LoadBalancerName)

	req = req.UpdateLoadBalancerPayload(model.Payload)
	return req
}
