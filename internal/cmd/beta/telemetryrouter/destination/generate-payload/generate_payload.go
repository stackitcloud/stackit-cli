package generatepayload

import (
	"context"
	"encoding/json"
	"fmt"

	telemetryrouter "github.com/stackitcloud/stackit-sdk-go/services/telemetryrouter/v1api"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	cliErr "github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/fileutils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/telemetryrouter/client"
	telemetryrouterUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/telemetryrouter/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"

	"github.com/spf13/cobra"
)

const (
	instanceIdFlag    = "instance-id"
	destinationIdFlag = "destination-id"
	configTypeFlag    = "config-type"
	filePathFlag      = "file-path"

	configTypeOpenTelemetry = "opentelemetry"
	configTypeS3            = "s3"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	InstanceId    string
	DestinationId *string
	ConfigType    string
	FilePath      *string
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "generate-payload",
		Short: "Generates a payload to create/update TelemetryRouter destinations",
		Long: fmt.Sprintf("%s\n%s\n%s\n%s\n%s",
			"Generates a JSON payload with values to be used as --payload input for destination creation or update.",
			"This command can be used to generate a payload to update an existing destination or to create a new destination.",
			"To update an existing destination, provide the destination ID and the instance ID of the TelemetryRouter instance.",
			`To obtain a default payload to create a new destination, run the command with the "--config-type" flag set to either "opentelemetry" (default) or "s3".`,
			"Note that the default values provided, such as the URI, bucket name or credentials, should be adapted to your use case.",
		),
		Args: args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Generate a create payload with default values for an OpenTelemetry destination, and adapt it with custom values`,
				`$ stackit beta telemetryrouter destination generate-payload --file-path ./payload.json`,
				`<Modify payload in file, if needed>`,
				`$ stackit beta telemetryrouter destination create --instance-id xxx --payload @./payload.json`),
			examples.NewExample(
				`Generate a create payload with default values for an S3 destination`,
				`$ stackit beta telemetryrouter destination generate-payload --config-type s3 --file-path ./payload.json`),
			examples.NewExample(
				`Generate an update payload with the values of an existing destination "yyy" for TelemetryRouter instance "xxx", and adapt it with custom values`,
				`$ stackit beta telemetryrouter destination generate-payload --destination-id yyy --instance-id xxx --file-path ./payload.json`,
				`<Modify payload in file>`,
				`$ stackit beta telemetryrouter destination update yyy --instance-id xxx --payload @./payload.json`),
			examples.NewExample(
				`Generate an update payload with the values of an existing destination "yyy" for TelemetryRouter instance "xxx", and preview it in the terminal`,
				`$ stackit beta telemetryrouter destination generate-payload --destination-id yyy --instance-id xxx`),
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

			if model.DestinationId == nil {
				var createPayload telemetryrouter.CreateDestinationPayload
				switch model.ConfigType {
				case configTypeS3:
					createPayload = telemetryrouterUtils.DefaultCreateDestinationPayloadS3
				default:
					createPayload = telemetryrouterUtils.DefaultCreateDestinationPayloadOpenTelemetry
				}
				return outputCreateResult(params.Printer, model.FilePath, &createPayload)
			}

			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("read TelemetryRouter destination: %w", err)
			}

			payload, err := telemetryrouterUtils.MapToUpdateDestinationPayload(resp)
			if err != nil {
				return fmt.Errorf("map update destination payload: %w", err)
			}

			return outputUpdateResult(params.Printer, model.FilePath, payload)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), instanceIdFlag, "ID of the TelemetryRouter instance")
	cmd.Flags().Var(flags.UUIDFlag(), destinationIdFlag, "If set, generates an update payload with the current state of the given destination. If unset, generates a create payload with default values")
	cmd.Flags().String(configTypeFlag, configTypeOpenTelemetry, fmt.Sprintf(`Type of the default create payload to generate, one of %q or %q. Only relevant if "--destination-id" is unset`, configTypeOpenTelemetry, configTypeS3))
	cmd.Flags().StringP(filePathFlag, "f", "", "If set, writes the payload to the given file. If unset, writes the payload to the standard output")
}

func parseInput(p *print.Printer, cmd *cobra.Command, _ []string) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)

	destinationId := flags.FlagToStringPointer(p, cmd, destinationIdFlag)
	instanceId := flags.FlagToStringValue(p, cmd, instanceIdFlag)

	if destinationId != nil {
		if globalFlags.ProjectId == "" {
			return nil, &cliErr.ProjectIdError{}
		}
		if instanceId == "" {
			return nil, &cliErr.FlagValidationError{
				Flag:    instanceIdFlag,
				Details: fmt.Sprintf("must be set when %q is provided", destinationIdFlag),
			}
		}
	}

	configType := flags.FlagWithDefaultToStringValue(p, cmd, configTypeFlag)
	if configType != configTypeOpenTelemetry && configType != configTypeS3 {
		return nil, &cliErr.FlagValidationError{
			Flag:    configTypeFlag,
			Details: fmt.Sprintf("must be one of %q or %q", configTypeOpenTelemetry, configTypeS3),
		}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		InstanceId:      instanceId,
		DestinationId:   destinationId,
		ConfigType:      configType,
		FilePath:        flags.FlagToStringPointer(p, cmd, filePathFlag),
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *telemetryrouter.APIClient) telemetryrouter.ApiGetDestinationRequest {
	return apiClient.DefaultAPI.GetDestination(ctx, model.ProjectId, model.Region, model.InstanceId, *model.DestinationId)
}

func outputCreateResult(p *print.Printer, filePath *string, payload *telemetryrouter.CreateDestinationPayload) error {
	if payload == nil {
		return fmt.Errorf("payload is nil")
	}

	payloadBytes, err := json.MarshalIndent(*payload, "", "  ") // nolint:gosec // false positive
	if err != nil {
		return fmt.Errorf("marshal payload: %w", err)
	}

	if filePath != nil {
		err = fileutils.WriteToFile(*filePath, string(payloadBytes))
		if err != nil {
			return fmt.Errorf("write payload to the file: %w", err)
		}
	} else {
		p.Outputln(string(payloadBytes))
	}

	return nil
}

func outputUpdateResult(p *print.Printer, filePath *string, payload *telemetryrouter.UpdateDestinationPayload) error {
	if payload == nil {
		return fmt.Errorf("payload is nil")
	}

	payloadBytes, err := json.MarshalIndent(*payload, "", "  ") // nolint:gosec // false positive
	if err != nil {
		return fmt.Errorf("marshal payload: %w", err)
	}

	if filePath != nil {
		err = fileutils.WriteToFile(*filePath, string(payloadBytes))
		if err != nil {
			return fmt.Errorf("write payload to the file: %w", err)
		}
	} else {
		p.Outputln(string(payloadBytes))
	}

	return nil
}
