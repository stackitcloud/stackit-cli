package list

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
	limitFlag = "limit"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	Limit *int64
}

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Lists all server command templates",
		Long:  "Lists all server command templates.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`List all command templates`,
				"$ stackit beta server command template list"),
			examples.NewExample(
				`List all commands templates in JSON format`,
				"$ stackit beta server command template list --output-format json"),
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

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("list server command templates: %w", err)
			}
			if templates := resp.Items; templates == nil || len(*templates) == 0 {
				p.Info("No commands templates found\n")
				return nil
			}
			templates := *resp.Items

			// Truncate output
			if model.Limit != nil && len(templates) > int(*model.Limit) {
				templates = templates[:*model.Limit]
			}
			return outputResult(p, model.OutputFormat, templates)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Int64(limitFlag, 0, "Maximum number of entries to list")
}

func parseInput(p *print.Printer, cmd *cobra.Command) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	limit := flags.FlagToInt64Pointer(p, cmd, limitFlag)
	if limit != nil && *limit < 1 {
		return nil, &errors.FlagValidationError{
			Flag:    limitFlag,
			Details: "must be greater than 0",
		}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		Limit:           limit,
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

func buildRequest(ctx context.Context, _ *inputModel, apiClient *runcommand.APIClient) runcommand.ApiListCommandTemplatesRequest {
	req := apiClient.ListCommandTemplates(ctx)
	return req
}

func outputResult(p *print.Printer, outputFormat string, templates []runcommand.CommandTemplate) error {
	switch outputFormat {
	case print.JSONOutputFormat:
		details, err := json.MarshalIndent(templates, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal server command template list: %w", err)
		}
		p.Outputln(string(details))

		return nil
	case print.YAMLOutputFormat:
		details, err := yaml.MarshalWithOptions(templates, yaml.IndentSequence(true), yaml.UseJSONMarshaler())
		if err != nil {
			return fmt.Errorf("marshal server command template list: %w", err)
		}
		p.Outputln(string(details))

		return nil
	default:
		table := tables.NewTable()
		table.SetHeader("NAME", "OS TYPE", "TITLE")
		for i := range templates {
			s := templates[i]

			var osType string
			if s.OsType != nil && len(*s.OsType) > 0 {
				osType = utils.JoinStringPtr(s.OsType, ",")
			}

			table.AddRow(
				utils.PtrString(s.Name),
				osType,
				utils.PtrString(s.Title),
			)
		}
		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}
		return nil
	}
}
