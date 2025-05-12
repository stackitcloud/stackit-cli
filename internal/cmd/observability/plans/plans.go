package plans

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/projectname"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/observability/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/goccy/go-yaml"
	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/observability"
)

const (
	limitFlag = "limit"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	Limit *int64
}

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "plans",
		Short: "Lists all Observability service plans",
		Long:  "Lists all Observability service plans.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`List all Observability service plans`,
				"$ stackit observability plans"),
			examples.NewExample(
				`List all Observability service plans in JSON format`,
				"$ stackit observability plans --output-format json"),
			examples.NewExample(
				`List up to 10 Observability service plans`,
				"$ stackit observability plans --limit 10"),
		),
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := context.Background()
			model, err := parseInput(params.Printer, cmd)
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
				return fmt.Errorf("get Observability service plans: %w", err)
			}
			plans := *resp.Plans
			if len(plans) == 0 {
				projectLabel, err := projectname.GetProjectName(ctx, params.Printer, params.CliVersion, cmd)
				if err != nil {
					params.Printer.Debug(print.ErrorLevel, "get project name: %v", err)
					projectLabel = model.ProjectId
				}
				params.Printer.Info("No plans found for project %q\n", projectLabel)
				return nil
			}

			// Truncate output
			if model.Limit != nil && len(plans) > int(*model.Limit) {
				plans = plans[:*model.Limit]
			}

			return outputResult(params.Printer, model.OutputFormat, plans)
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

func buildRequest(ctx context.Context, model *inputModel, apiClient *observability.APIClient) observability.ApiListPlansRequest {
	req := apiClient.ListPlans(ctx, model.ProjectId)
	return req
}

func outputResult(p *print.Printer, outputFormat string, plans []observability.Plan) error {
	switch outputFormat {
	case print.JSONOutputFormat:
		details, err := json.MarshalIndent(plans, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal Observability plans: %w", err)
		}
		p.Outputln(string(details))

		return nil
	case print.YAMLOutputFormat:
		details, err := yaml.MarshalWithOptions(plans, yaml.IndentSequence(true), yaml.UseJSONMarshaler())
		if err != nil {
			return fmt.Errorf("marshal Observability plans: %w", err)
		}
		p.Outputln(string(details))

		return nil
	default:
		table := tables.NewTable()
		table.SetHeader("ID", "PLAN NAME", "DESCRIPTION")
		for i := range plans {
			o := plans[i]
			table.AddRow(
				utils.PtrString(o.Id),
				utils.PtrString(o.Name),
				utils.PtrString(o.Description),
			)
			table.AddSeparator()
		}
		table.EnableAutoMergeOnColumns(1)
		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}

		return nil
	}
}
