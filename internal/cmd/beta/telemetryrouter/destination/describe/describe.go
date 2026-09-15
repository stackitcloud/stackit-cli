package describe

import (
	"context"
	"fmt"

	telemetryrouter "github.com/stackitcloud/stackit-sdk-go/services/telemetryrouter/v1api"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	cliErr "github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/telemetryrouter/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

const (
	instanceIdFlag   = "instance-id"
	destinationIdArg = "DESTINATION_ID"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	InstanceId    string
	DestinationId string
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("describe %s", destinationIdArg),
		Short: "Shows details of a TelemetryRouter destination",
		Long:  "Shows details of a TelemetryRouter destination.",
		Args:  args.SingleArg(destinationIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Show details of a TelemetryRouter destination with ID "xxx"`,
				"$ stackit beta telemetryrouter destination describe xxx --instance-id yyy",
			),
			examples.NewExample(
				`Show details of a TelemetryRouter destination with ID "xxx" in JSON format`,
				"$ stackit beta telemetryrouter destination describe xxx --instance-id yyy --output-format json",
			),
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

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("read TelemetryRouter destination: %w", err)
			}

			return outputResult(params.Printer, model.OutputFormat, resp)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), instanceIdFlag, "ID of the TelemetryRouter instance")

	err := flags.MarkFlagsRequired(cmd, instanceIdFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	destinationId := inputArgs[0]

	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &cliErr.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		InstanceId:      flags.FlagToStringValue(p, cmd, instanceIdFlag),
		DestinationId:   destinationId,
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *telemetryrouter.APIClient) telemetryrouter.ApiGetDestinationRequest {
	return apiClient.DefaultAPI.GetDestination(ctx, model.ProjectId, model.Region, model.InstanceId, model.DestinationId)
}

func outputResult(p *print.Printer, outputFormat string, destination *telemetryrouter.DestinationResponse) error {
	if destination == nil {
		return fmt.Errorf("destination response is empty")
	}
	return p.OutputResult(outputFormat, destination, func() error {
		table := tables.NewTable()
		table.AddRow("ID", destination.Id)
		table.AddSeparator()
		table.AddRow("DISPLAY NAME", destination.DisplayName)
		table.AddSeparator()
		table.AddRow("DESCRIPTION", utils.PtrString(destination.Description))
		table.AddSeparator()
		table.AddRow("STATUS", destination.Status)
		table.AddSeparator()
		table.AddRow("CREDENTIAL TYPE", destination.CredentialType)
		table.AddSeparator()
		table.AddRow("CONFIG TYPE", destination.Config.ConfigType)
		table.AddSeparator()
		table.AddRow("CREATION TIME", destination.CreationTime)

		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("display table: %w", err)
		}
		return nil
	})
}
