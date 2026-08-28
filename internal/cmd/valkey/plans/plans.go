package plans

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	valkey "github.com/stackitcloud/stackit-sdk-go/services/valkey/v2api"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/projectname"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/valkey/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
)

const (
	limitFlag = "limit"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	Limit *int64
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "plans",
		Short: "Lists all Valkey service plans",
		Long:  "Lists all Valkey service plans.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				"Lists all Valkey service plans",
				"$ stackit valkey plans",
			),
			examples.NewExample(
				`List all Valkey service plans in JSON format`,
				"$ stackit valkey plans --output-format json"),
			examples.NewExample(
				`List up to 10 Valkey service plans`,
				"$ stackit valkey plans --limit 10"),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			model, err := parseInput(params.Printer, cmd, args)
			if err != nil {
				return err
			}
			apiclient, err := client.ConfigureClient(params.Printer, params.CliVersion)
			if err != nil {
				return err
			}

			// Call API
			req := buildRequest(ctx, model, apiclient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("get Valkey service plans: %w", err)
			}
			plans := resp.Offerings

			// Truncate output
			if model.Limit != nil && len(plans) > int(*model.Limit) {
				plans = plans[:*model.Limit]
			}

			projectLabel, err := projectname.GetProjectName(ctx, params.Printer, params.CliVersion, cmd)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get project name: %v", err)
				projectLabel = model.ProjectId
			}

			return outputResult(params.Printer, model.OutputFormat, projectLabel, plans)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Int64(limitFlag, 0, "Maximum number of entries to list")
}

func parseInput(p *print.Printer, cmd *cobra.Command, _ []string) (*inputModel, error) {
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

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *valkey.APIClient) valkey.ApiListOfferingsRequest {
	req := apiClient.DefaultAPI.ListOfferings(ctx, model.ProjectId, model.Region)
	return req
}

func outputResult(p *print.Printer, outputFormat, projectLabel string, plans []valkey.Offering) error {
	return p.OutputResult(outputFormat, plans, func() error {
		if len(plans) == 0 {
			p.Outputf("No plans found for project %q\n", projectLabel)
			return nil
		}
		table := tables.NewTable()
		table.SetHeader("OFFERING NAME", "VERSION", "ID", "NAME", "DESCRIPTION")
		for i := range plans {
			o := plans[i]
			if o.Plans != nil {
				for j := range o.Plans {
					plan := (o.Plans)[j]
					table.AddRow(
						o.Name,
						o.Version,
						plan.Id,
						plan.Name,
						plan.Description,
					)
				}
				table.AddSeparator()
			}
		}
		table.EnableAutoMergeOnColumns(1, 2)
		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}

		return nil
	})
}
