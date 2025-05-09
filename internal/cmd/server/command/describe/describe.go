package describe

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/goccy/go-yaml"
	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/runcommand/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/runcommand"
)

const (
	commandIdArg = "COMMAND_ID"
	serverIdFlag = "server-id"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	ServerId  string
	CommandId string
}

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("describe %s", commandIdArg),
		Short: "Shows details of a Server Command",
		Long:  "Shows details of a Server Command.",
		Args:  args.SingleArg(commandIdArg, nil),
		Example: examples.Build(
			examples.NewExample(
				`Get details of a Server Command with ID "xxx" for server with ID "yyy"`,
				"$ stackit server command describe xxx --server-id=yyy"),
			examples.NewExample(
				`Get details of a Server Command with ID "xxx" for server with ID "yyy" in JSON format`,
				"$ stackit server command describe xxx --server-id=yyy --output-format json"),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			model, err := parseInput(params.Printer, cmd, args)
			if err != nil {
				return err
			}
			// Configure API client
			apiClient, err := client.ConfigureClient(params.Printer)
			if err != nil {
				return err
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("read server command: %w", err)
			}

			return outputResult(params.Printer, model.OutputFormat, *resp)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().VarP(flags.UUIDFlag(), serverIdFlag, "s", "Server ID")

	err := flags.MarkFlagsRequired(cmd, serverIdFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	commandId := inputArgs[0]

	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		ServerId:        flags.FlagToStringValue(p, cmd, serverIdFlag),
		CommandId:       commandId,
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

func buildRequest(ctx context.Context, model *inputModel, apiClient *runcommand.APIClient) runcommand.ApiGetCommandRequest {
	req := apiClient.GetCommand(ctx, model.ProjectId, model.Region, model.ServerId, model.CommandId)
	return req
}

func outputResult(p *print.Printer, outputFormat string, command runcommand.CommandDetails) error {
	switch outputFormat {
	case print.JSONOutputFormat:
		details, err := json.MarshalIndent(command, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal server command: %w", err)
		}
		p.Outputln(string(details))

		return nil
	case print.YAMLOutputFormat:
		details, err := yaml.MarshalWithOptions(command, yaml.IndentSequence(true), yaml.UseJSONMarshaler())
		if err != nil {
			return fmt.Errorf("marshal server command: %w", err)
		}
		p.Outputln(string(details))

		return nil
	default:
		table := tables.NewTable()
		table.AddRow("ID", utils.PtrString(command.Id))
		table.AddSeparator()
		table.AddRow("COMMAND TEMPLATE NAME", utils.PtrString(command.CommandTemplateName))
		table.AddSeparator()
		table.AddRow("COMMAND TEMPLATE TITLE", utils.PtrString(command.CommandTemplateTitle))
		table.AddSeparator()
		table.AddRow("STATUS", utils.PtrString(command.Status))
		table.AddSeparator()
		table.AddRow("STARTED AT", utils.PtrString(command.StartedAt))
		table.AddSeparator()
		table.AddRow("FINISHED AT", utils.PtrString(command.FinishedAt))
		table.AddSeparator()
		table.AddRow("EXIT CODE", utils.PtrString(command.ExitCode))
		table.AddSeparator()
		table.AddRow("COMMAND SCRIPT", utils.PtrString(command.Script))
		table.AddSeparator()
		table.AddRow("COMMAND OUTPUT", utils.PtrString(command.Output))
		table.AddSeparator()
		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}

		return nil
	}
}
