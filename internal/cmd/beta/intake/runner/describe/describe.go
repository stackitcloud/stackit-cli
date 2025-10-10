package describe

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/goccy/go-yaml"
	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/intake"

	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	cliErr "github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/intake/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

const (
	runnerIdArg = "RUNNER_ID"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	RunnerId string
}

func NewDescribeCmd(p *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("describe %s", runnerIdArg),
		Short: "Shows details of an Intake Runner",
		Long:  "Shows details of an Intake Runner.",
		Args:  args.SingleArg(runnerIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Get details of an Intake Runner with ID "xxx"`,
				`$ stackit beta intake runner describe xxx`),
			examples.NewExample(
				`Get details of an Intake Runner with ID "xxx" in JSON format`,
				`$ stackit beta intake runner describe xxx --output-format json`),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			model, err := parseInput(p.Printer, cmd, args)
			if err != nil {
				return err
			}

			// Configure API client
			apiClient, err := client.ConfigureClient(p.Printer, p.CliVersion)
			if err != nil {
				return err
			}

			// Call API to get a single runner
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("get Intake Runner: %w", err)
			}

			return outputResult(p.Printer, model.OutputFormat, resp)
		},
	}
	return cmd
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	runnerId := inputArgs[0]

	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &cliErr.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		RunnerId:        runnerId,
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

// buildRequest creates the API request to get a single Intake Runner
func buildRequest(ctx context.Context, model *inputModel, apiClient *intake.APIClient) intake.ApiGetIntakeRunnerRequest {
	req := apiClient.GetIntakeRunner(ctx, model.ProjectId, model.Region, model.RunnerId)
	return req
}

// outputResult formats the API response and prints it to the console
func outputResult(p *print.Printer, outputFormat string, runner *intake.IntakeRunnerResponse) error {
	if runner == nil {
		return fmt.Errorf("received nil runner, could not display details")
	}

	switch outputFormat {
	case print.JSONOutputFormat:
		details, err := json.MarshalIndent(runner, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal Intake Runner: %w", err)
		}
		p.Outputln(string(details))
		return nil

	case print.YAMLOutputFormat:
		details, err := yaml.MarshalWithOptions(runner, yaml.IndentSequence(true), yaml.UseJSONMarshaler())
		if err != nil {
			return fmt.Errorf("marshal Intake Runner: %w", err)
		}
		p.Outputln(string(details))
		return nil

	default:
		table := tables.NewTable()
		table.SetHeader("Attribute", "Value")
		table.AddRow("ID", runner.GetId())
		table.AddRow("Name", runner.GetDisplayName())
		table.AddRow("State", runner.GetState())
		table.AddRow("Created", runner.GetCreateTime())
		table.AddRow("Labels", runner.GetLabels())
		table.AddRow("Description", runner.GetDescription())
		table.AddRow("Max Message Size (KiB)", runner.GetMaxMessageSizeKiB())
		table.AddRow("Max Messages/Hour", runner.GetMaxMessagesPerHour())
		table.AddRow("Ingestion URI", runner.GetUri())

		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}
		return nil
	}
}
