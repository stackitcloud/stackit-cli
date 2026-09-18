package create

import (
	"context"
	"encoding/json"
	"fmt"

	telemetryrouter "github.com/stackitcloud/stackit-sdk-go/services/telemetryrouter/v1api"
	"github.com/stackitcloud/stackit-sdk-go/services/telemetryrouter/v1api/wait"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	cliErr "github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/projectname"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/telemetryrouter/client"
	telemetryrouterUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/telemetryrouter/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/spinner"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"

	"github.com/spf13/cobra"
)

const (
	instanceIdFlag = "instance-id"
	payloadFlag    = "payload"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	InstanceId string
	Payload    *telemetryrouter.CreateDestinationPayload
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Creates a TelemetryRouter destination",
		Long: fmt.Sprintf("%s\n%s\n%s\n%s",
			"Creates a destination for a TelemetryRouter instance.",
			"The payload can be provided as a JSON string or a file path prefixed with \"@\".",
			"If no payload is provided, a default payload will be used.",
			"The \"config\" field is a discriminated union: exactly one of \"openTelemetry\" or \"s3\" must be set, matching the \"configType\".",
		),
		Args: args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Create a destination on TelemetryRouter instance "xxx" using the default (OpenTelemetry) configuration`,
				"$ stackit beta telemetryrouter destination create --instance-id xxx"),
			examples.NewExample(
				`Create a destination on TelemetryRouter instance "xxx" using an API payload sourced from the file "./payload.json"`,
				"$ stackit beta telemetryrouter destination create --payload @./payload.json --instance-id xxx"),
			examples.NewExample(
				`Create a destination on TelemetryRouter instance "xxx" using an API payload provided as a JSON string`,
				`$ stackit beta telemetryrouter destination create --payload "{...}" --instance-id xxx`),
			examples.NewExample(
				`Generate a payload with default values, and adapt it with custom values for the different configuration options`,
				`$ stackit beta telemetryrouter destination generate-payload --config-type s3 > ./payload.json`,
				`<Modify payload in file, if needed>`,
				`$ stackit beta telemetryrouter destination create --payload @./payload.json --instance-id xxx`),
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

			projectLabel, err := projectname.GetProjectName(ctx, params.Printer, params.CliVersion, cmd)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get project name: %v", err)
				projectLabel = model.ProjectId
			}

			instanceLabel, err := telemetryrouterUtils.GetInstanceName(ctx, apiClient.DefaultAPI, model.ProjectId, model.Region, model.InstanceId)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get instance name: %v", err)
				instanceLabel = model.InstanceId
			}

			// Fill in default payload, if needed
			if model.Payload == nil {
				defaultPayload := telemetryrouterUtils.DefaultCreateDestinationPayloadOpenTelemetry
				model.Payload = &defaultPayload
			}

			prompt := fmt.Sprintf("Are you sure you want to create destination %q for TelemetryRouter instance %q in project %q?", model.Payload.DisplayName, instanceLabel, projectLabel)
			err = params.Printer.PromptForConfirmation(prompt)
			if err != nil {
				return err
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("create TelemetryRouter destination: %w", err)
			}
			if resp == nil {
				return fmt.Errorf("create TelemetryRouter destination: empty response from API")
			}

			// Wait for async operation, if async mode not enabled
			if !model.Async {
				err := spinner.Run(params.Printer, "Creating destination", func() error {
					_, err = wait.CreateDestinationWaitHandler(ctx, apiClient.DefaultAPI, model.ProjectId, model.Region, model.InstanceId, resp.Id).WaitWithContext(ctx)
					return err
				})
				if err != nil {
					return fmt.Errorf("wait for TelemetryRouter destination creation: %w", err)
				}
			}

			return outputResult(params.Printer, model, instanceLabel, resp)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), instanceIdFlag, "ID of the TelemetryRouter instance")
	cmd.Flags().Var(flags.ReadFromFileFlag(), payloadFlag, `Request payload (JSON). Can be a string or a file path, if prefixed with "@" (example: @./payload.json). If unset, will use a default payload (you can check it by running "stackit beta telemetryrouter destination generate-payload")`)

	err := flags.MarkFlagsRequired(cmd, instanceIdFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command, _ []string) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &cliErr.ProjectIdError{}
	}

	payloadValue := flags.FlagToStringPointer(p, cmd, payloadFlag)
	var payload *telemetryrouter.CreateDestinationPayload
	if payloadValue != nil {
		payload = &telemetryrouter.CreateDestinationPayload{}
		err := json.Unmarshal([]byte(*payloadValue), payload)
		if err != nil {
			return nil, fmt.Errorf("decode payload: %w", err)
		}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		InstanceId:      flags.FlagToStringValue(p, cmd, instanceIdFlag),
		Payload:         payload,
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *telemetryrouter.APIClient) telemetryrouter.ApiCreateDestinationRequest {
	req := apiClient.DefaultAPI.CreateDestination(ctx, model.ProjectId, model.Region, model.InstanceId)
	req = req.CreateDestinationPayload(*model.Payload)
	return req
}

func outputResult(p *print.Printer, model *inputModel, instanceLabel string, resp *telemetryrouter.DestinationResponse) error {
	if resp == nil {
		return fmt.Errorf("create TelemetryRouter destination response is empty")
	} else if model == nil || model.GlobalFlagModel == nil {
		return fmt.Errorf("input model is nil")
	}

	return p.OutputResult(model.OutputFormat, resp, func() error {
		operationState := "Created"
		if model.Async {
			operationState = "Triggered creation of"
		}
		p.Outputf("%s destination %q for TelemetryRouter instance %q. Destination ID: %s\n", operationState, resp.DisplayName, instanceLabel, resp.Id)
		return nil
	})
}
