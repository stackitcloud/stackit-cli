package create

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"

	"github.com/stackitcloud/stackit-cli/internal/pkg/projectname"
	"github.com/stackitcloud/stackit-sdk-go/services/loadbalancer"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/load-balancer/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/spinner"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/loadbalancer/wait"
)

const (
	payloadFlag = "payload"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	Payload *loadbalancer.CreateLoadBalancerPayload
}

var (
	xRequestId = uuid.NewString()
)

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Creates a Load Balancer",
		Long: fmt.Sprintf("%s\n%s\n%s",
			"Creates a Load Balancer.",
			"The payload can be provided as a JSON string or a file path prefixed with \"@\".",
			"See https://docs.api.stackit.cloud/documentation/load-balancer/version/v1#tag/Load-Balancer/operation/APIService_CreateLoadBalancer for information regarding the payload structure.",
		),
		Args: args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Create a load balancer using an API payload sourced from the file "./payload.json"`,
				"$ stackit load-balancer create --payload @./payload.json"),
			examples.NewExample(
				`Create a load balancer using an API payload provided as a JSON string`,
				`$ stackit load-balancer create --payload "{...}"`),
			examples.NewExample(
				`Generate a payload with default values, and adapt it with custom values for the different configuration options`,
				`$ stackit load-balancer generate-payload > ./payload.json`,
				`<Modify payload in file>`,
				`$ stackit load-balancer create --payload @./payload.json`),
		),
		RunE: func(cmd *cobra.Command, _ []string) error {
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

			projectLabel, err := projectname.GetProjectName(ctx, p, cmd)
			if err != nil {
				p.Debug(print.ErrorLevel, "get project name: %v", err)
				projectLabel = model.ProjectId
			}

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to create a load balancer for project %q?", projectLabel)
				err = p.PromptForConfirmation(prompt)
				if err != nil {
					return err
				}
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			_, err = req.Execute()
			if err != nil {
				return fmt.Errorf("create load balancer: %w", err)
			}

			// Wait for async operation, if async mode not enabled
			if !model.Async {
				s := spinner.New(p)
				s.Start("Creating load balancer")
				_, err = wait.CreateLoadBalancerWaitHandler(ctx, apiClient, model.ProjectId, *model.Payload.Name).WaitWithContext(ctx)
				if err != nil {
					return fmt.Errorf("wait for load balancer creation: %w", err)
				}
				s.Stop()
			}

			operationState := "Created"
			if model.Async {
				operationState = "Triggered creation of"
			}
			p.Outputf("%s load balancer with name %q \n", operationState, *model.Payload.Name)
			return nil
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.ReadFromFileFlag(), payloadFlag, `Request payload (JSON). Can be a string or a file path, if prefixed with "@" (example: @./payload.json).`)

	err := flags.MarkFlagsRequired(cmd, payloadFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	payloadValue := flags.FlagToStringPointer(p, cmd, payloadFlag)
	var payload *loadbalancer.CreateLoadBalancerPayload
	if payloadValue != nil {
		payload = &loadbalancer.CreateLoadBalancerPayload{}
		err := json.Unmarshal([]byte(*payloadValue), payload)
		if err != nil {
			return nil, fmt.Errorf("encode payload: %w", err)
		}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		Payload:         payload,
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

func buildRequest(ctx context.Context, model *inputModel, apiClient *loadbalancer.APIClient) loadbalancer.ApiCreateLoadBalancerRequest {
	req := apiClient.CreateLoadBalancer(ctx, model.ProjectId)
	req = req.CreateLoadBalancerPayload(*model.Payload)
	req = req.XRequestID(xRequestId)
	return req
}
