package list

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	edge "github.com/stackitcloud/stackit-sdk-go/services/edge/v1beta1api"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	cliErr "github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/projectname"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/edge/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

// User input struct for the command
const (
	limitFlag = "limit"
)

// Struct to model user input (arguments and/or flags)
type inputModel struct {
	*globalflags.GlobalFlagModel
	Limit *int64
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Lists available Edge Cloud service plans",
		Long:  "Lists available STACKIT Edge Cloud (STEC) service plans of a project",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Lists all Edge Cloud plans for a given project`,
				`$ stackit beta edge-cloud plans list`),
			examples.NewExample(
				`Lists all Edge Cloud plans for a given project and limits the output to two plans`,
				`$ stackit beta edge-cloud plans list --limit=2`),
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

			projectLabel, err := projectname.GetProjectName(ctx, params.Printer, params.CliVersion, cmd)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get project name: %v", err)
				// If project label can't be determined, fall back to project ID
				projectLabel = model.ProjectId
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("get Edge Cloud plans: %w", err)
			}
			plans := resp.GetValidPlans()

			// Truncate output
			if model.Limit != nil && len(plans) > int(*model.Limit) {
				plans = (plans)[:*model.Limit]
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
		return nil, &cliErr.ProjectIdError{}
	}

	// Parse and validate user input then add it to the model
	limit := flags.FlagToInt64Pointer(p, cmd, limitFlag)
	if limit != nil && *limit < 1 {
		return nil, &cliErr.FlagValidationError{
			Flag:    limitFlag,
			Details: "must be greater than 0",
		}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		Limit:           limit,
	}

	// Log the parsed model if --verbosity is set to debug
	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *edge.APIClient) edge.ApiListPlansProjectRequest {
	return apiClient.DefaultAPI.ListPlansProject(ctx, model.ProjectId)
}

func outputResult(p *print.Printer, outputFormat, projectLabel string, plans []edge.Plan) error {
	return p.OutputResult(outputFormat, plans, func() error {
		// No plans found for project
		if len(plans) == 0 {
			p.Outputf("No plans found for project %q\n", projectLabel)
			return nil
		}

		// Display plans found for project in a table
		table := tables.NewTable()
		// List: only output the most important fields. Be sure to filter for any non-required fields.
		table.SetHeader("ID", "NAME", "DESCRIPTION", "MAX EDGE HOSTS")
		for i := range plans {
			plan := plans[i]
			table.AddRow(
				utils.PtrString(plan.Id),
				utils.PtrString(plan.Name),
				utils.PtrString(plan.Description),
				utils.PtrString(plan.MaxEdgeHosts))
		}
		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}
		return nil
	})
}
