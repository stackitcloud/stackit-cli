package list

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/core/experimental/paginate"
	automation "github.com/stackitcloud/stackit-sdk-go/services/automation/v1betaapi"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	cliErr "github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/projectname"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/automation/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

const (
	automationIdFlag = "automation-id"
	limitFlag        = "limit"

	// maxPageSize is the maximum number of items the API returns per page.
	maxPageSize = 100
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	AutomationId string
	Limit        *int64
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all Volume Automation Executions",
		Long:  "List all Volume Automation Executions.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`List all Volume Automation Executions for Automation with ID "xxx"`,
				"$ stackit beta volume automation execution list --automation-id xxx"),
			examples.NewExample(
				`List up to 10 Volume Automation Executions for Automation with ID "xxx"`,
				"$ stackit beta volume automation execution list --automation-id xxx --limit 10"),
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
			resp, err := fetchExecutions(ctx, model, apiClient)
			if err != nil {
				return fmt.Errorf("list volume automation executions: %w", err)
			}

			// Get projectLabel
			projectLabel, err := projectname.GetProjectName(ctx, params.Printer, params.CliVersion, cmd)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get project name: %v", err)
				projectLabel = model.ProjectId
			} else if projectLabel == "" {
				projectLabel = model.ProjectId
			}

			return outputResult(params.Printer, model.OutputFormat, model.AutomationId, projectLabel, resp)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), automationIdFlag, "Automation ID")
	cmd.Flags().Int64(limitFlag, 0, "Maximum number of entries to list")

	err := flags.MarkFlagsRequired(cmd, automationIdFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command, _ []string) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &cliErr.ProjectIdError{}
	}

	limit := flags.FlagToInt64Pointer(p, cmd, limitFlag)
	if limit != nil && *limit < 1 {
		return nil, &cliErr.FlagValidationError{
			Flag:    limitFlag,
			Details: "must be greater than 0",
		}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		AutomationId:    flags.FlagToStringValue(p, cmd, automationIdFlag),
		Limit:           limit,
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *automation.APIClient) automation.ApiListVolumeExecutionsRequest {
	return apiClient.DefaultAPI.ListVolumeExecutions(ctx, model.ProjectId, model.Region, model.AutomationId)
}

func fetchExecutions(ctx context.Context, model *inputModel, apiClient *automation.APIClient) ([]automation.ListExecutionsItem, error) {
	req := buildRequest(ctx, model, apiClient)
	opts := []paginate.Option{paginate.WithPageSize(maxPageSize)}
	if model.Limit != nil {
		opts = append(opts, paginate.WithLimit(int(*model.Limit)))
	}

	items, err := paginate.All[automation.ListExecutionsItem](req, opts...)
	if err != nil {
		return nil, fmt.Errorf("list volume automation executions: %w", err)
	}
	if items == nil {
		items = []automation.ListExecutionsItem{}
	}

	return items, nil
}

func outputResult(p *print.Printer, outputFormat, automationId, projectLabel string, executions []automation.ListExecutionsItem) error {
	return p.OutputResult(outputFormat, executions, func() error {
		if len(executions) == 0 {
			p.Outputf("No volume automation executions found in project %q for automation %q\n", projectLabel, automationId)
			return nil
		}

		table := tables.NewTable()

		table.SetTitle("Volume Automation Executions")
		table.SetHeader("ID", "STATUS", "CREATE TIME", "END TIME")
		table.AddSeparator()
		for _, execution := range executions {
			table.AddRow(
				execution.Id,
				execution.Status,
				execution.CreateTime.Format(time.DateTime),
				utils.ConvertTimePToDateTimeString(execution.EndTime),
			)
			table.AddSeparator()
		}

		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("render tables: %w", err)
		}
		return nil
	})
}
