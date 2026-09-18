package update

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
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/telemetryrouter/client"
	telemetryrouterUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/telemetryrouter/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/spinner"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

const (
	destinationIdArg = "DESTINATION_ID"

	instanceIdFlag = "instance-id"
	payloadFlag    = "payload"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	InstanceId    string
	DestinationId string
	Payload       telemetryrouter.UpdateDestinationPayload
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("update %s", destinationIdArg),
		Short: "Updates a TelemetryRouter destination",
		Long: fmt.Sprintf("%s\n%s\n%s",
			"Updates a destination of a TelemetryRouter instance.",
			"The payload can be provided as a JSON string or a file path prefixed with \"@\".",
			"The \"config\" field is a discriminated union: exactly one of \"openTelemetry\" or \"s3\" must be set, matching the \"configType\".",
		),
		Args: args.SingleArg(destinationIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Update a destination with ID "xxx" from TelemetryRouter instance "yyy", using an API payload sourced from the file "./payload.json"`,
				"$ stackit beta telemetryrouter destination update xxx --instance-id yyy --payload @./payload.json"),
			examples.NewExample(
				`Update a destination with ID "xxx" from TelemetryRouter instance "yyy", using an API payload provided as a JSON string`,
				`$ stackit beta telemetryrouter destination update xxx --instance-id yyy --payload "{...}"`),
			examples.NewExample(
				`Generate a payload with the current values of a destination, and adapt it with custom values for the different configuration options`,
				`$ stackit beta telemetryrouter destination generate-payload --destination-id xxx --instance-id yyy > ./payload.json`,
				`<Modify payload in file>`,
				`$ stackit beta telemetryrouter destination update xxx --instance-id yyy --payload @./payload.json`),
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

			instanceLabel, err := telemetryrouterUtils.GetInstanceName(ctx, apiClient.DefaultAPI, model.ProjectId, model.Region, model.InstanceId)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get instance name: %v", err)
			}
			if instanceLabel == "" {
				instanceLabel = model.InstanceId
			}

			destinationLabel, err := telemetryrouterUtils.GetDestinationName(ctx, apiClient.DefaultAPI, model.ProjectId, model.Region, model.InstanceId, model.DestinationId)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get destination name: %v", err)
			}
			if destinationLabel == "" {
				destinationLabel = model.DestinationId
			}

			prompt := fmt.Sprintf("Are you sure you want to update destination %q for TelemetryRouter instance %q?", destinationLabel, instanceLabel)
			err = params.Printer.PromptForConfirmation(prompt)
			if err != nil {
				return err
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("update TelemetryRouter destination: %w", err)
			}
			if resp == nil {
				return fmt.Errorf("update TelemetryRouter destination: empty response from API")
			}

			// Wait for async operation, if async mode not enabled
			if !model.Async {
				err := spinner.Run(params.Printer, "Updating destination", func() error {
					_, err = wait.UpdateDestinationWaitHandler(ctx, apiClient.DefaultAPI, model.ProjectId, model.Region, model.InstanceId, model.DestinationId).WaitWithContext(ctx)
					return err
				})
				if err != nil {
					return fmt.Errorf("wait for TelemetryRouter destination update: %w", err)
				}
			}

			operationState := "Updated"
			if model.Async {
				operationState = "Triggered update of"
			}
			params.Printer.Outputf("%s destination %q for TelemetryRouter instance %q\n", operationState, destinationLabel, instanceLabel)
			return nil
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), instanceIdFlag, "ID of the TelemetryRouter instance")
	cmd.Flags().Var(flags.ReadFromFileFlag(), payloadFlag, `Request payload (JSON). Can be a string or a file path, if prefixed with "@". Example: @./payload.json`)

	err := flags.MarkFlagsRequired(cmd, instanceIdFlag, payloadFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	destinationId := inputArgs[0]

	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &cliErr.ProjectIdError{}
	}

	payloadString := flags.FlagToStringValue(p, cmd, payloadFlag)
	var payload telemetryrouter.UpdateDestinationPayload
	err := json.Unmarshal([]byte(payloadString), &payload)
	if err != nil {
		return nil, fmt.Errorf("decode payload: %w", err)
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		InstanceId:      flags.FlagToStringValue(p, cmd, instanceIdFlag),
		DestinationId:   destinationId,
		Payload:         payload,
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *telemetryrouter.APIClient) telemetryrouter.ApiUpdateDestinationRequest {
	req := apiClient.DefaultAPI.UpdateDestination(ctx, model.ProjectId, model.Region, model.InstanceId, model.DestinationId)
	req = req.UpdateDestinationPayload(model.Payload)
	return req
}
