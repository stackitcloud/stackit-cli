package describe

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/goccy/go-yaml"
	"github.com/spf13/cobra"
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
	commandTemplateNameArg = "COMMAND_TEMPLATE_NAME"
	serverIdFlag           = "server-id"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	ServerId            string
	CommandTemplateName string
}

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("describe %s", commandTemplateNameArg),
		Short: "Shows details of a Server Command Template",
		Long:  "Shows details of a Server Command Template.",
		Args:  args.SingleArg(commandTemplateNameArg, nil),
		Example: examples.Build(
			examples.NewExample(
				`Get details of a Server Command Template with name "RunShellScript" for server with ID "xxx"`,
				"$ stackit server command template describe RunShellScript --server-id=xxx"),
			examples.NewExample(
				`Get details of a Server Command Template with name "RunShellScript" for server with ID "xxx" in JSON format`,
				"$ stackit server command template describe RunShellScript --server-id=xxx --output-format json"),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			model, err := parseInput(p, cmd, args)
			if err != nil {
				return err
			}
			// Configure API client
			apiClient, err := client.ConfigureClient(p)
			if err != nil {
				return err
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("read server command template: %w", err)
			}

			return outputResult(p, model.OutputFormat, *resp)
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
	commandTemplateName := inputArgs[0]

	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel:     globalFlags,
		ServerId:            flags.FlagToStringValue(p, cmd, serverIdFlag),
		CommandTemplateName: commandTemplateName,
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

func buildRequest(ctx context.Context, model *inputModel, apiClient *runcommand.APIClient) runcommand.ApiGetCommandTemplateRequest {
	req := apiClient.GetCommandTemplate(ctx, model.ProjectId, model.ServerId, model.CommandTemplateName, model.Region)
	return req
}

func outputResult(p *print.Printer, outputFormat string, commandTemplate runcommand.CommandTemplateSchema) error {
	switch outputFormat {
	case print.JSONOutputFormat:
		details, err := json.MarshalIndent(commandTemplate, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal server command template: %w", err)
		}
		p.Outputln(string(details))

		return nil
	case print.YAMLOutputFormat:
		details, err := yaml.MarshalWithOptions(commandTemplate, yaml.IndentSequence(true), yaml.UseJSONMarshaler())
		if err != nil {
			return fmt.Errorf("marshal server command template: %w", err)
		}
		p.Outputln(string(details))

		return nil
	default:
		table := tables.NewTable()
		table.AddRow("NAME", utils.PtrString(commandTemplate.Name))
		table.AddSeparator()
		table.AddRow("TITLE", utils.PtrString(commandTemplate.Title))
		table.AddSeparator()
		table.AddRow("DESCRIPTION", utils.PtrString(commandTemplate.Description))
		table.AddSeparator()
		if commandTemplate.OsType != nil {
			table.AddRow("OS TYPE", utils.JoinStringPtr(commandTemplate.OsType, "\n"))
			table.AddSeparator()
		}
		if commandTemplate.ParameterSchema != nil {
			table.AddRow("PARAMS", *commandTemplate.ParameterSchema)
		} else {
			table.AddRow("PARAMS", "")
		}
		table.AddSeparator()
		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}

		return nil
	}
}
