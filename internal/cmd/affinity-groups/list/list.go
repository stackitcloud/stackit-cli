package list

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/types"

	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/iaas"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/projectname"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/iaas/client"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	Limit *int64
}

const limitFlag = "limit"

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Lists affinity groups",
		Long:  `Lists affinity groups.`,
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				"Lists all affinity groups",
				"$ stackit affinity-group list",
			),
			examples.NewExample(
				"Lists up to 10 affinity groups",
				"$ stackit affinity-group list --limit=10",
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
			request := buildRequest(ctx, *model, apiClient)
			result, err := request.Execute()
			if err != nil {
				return fmt.Errorf("list affinity groups: %w", err)
			}
			items := result.GetItems()

			projectLabel, err := projectname.GetProjectName(ctx, params.Printer, params.CliVersion, cmd)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get project name: %v", err)
				projectLabel = model.ProjectId
			}

			// Truncate Output
			if model.Limit != nil && len(items) > int(*model.Limit) {
				items = items[:*model.Limit]
			}
			return outputResult(params.Printer, model.OutputFormat, projectLabel, items)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Int64(limitFlag, 0, "Limit the output to the first n elements")
}

func buildRequest(ctx context.Context, model inputModel, apiClient *iaas.APIClient) iaas.ApiListAffinityGroupsRequest {
	return apiClient.ListAffinityGroups(ctx, model.ProjectId, model.Region)
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

func outputResult(p *print.Printer, outputFormat, projectLabel string, items []iaas.AffinityGroup) error {
	return p.OutputResult(outputFormat, items, func() error {
		if len(items) == 0 {
			p.Outputf("No affinity groups found for project %q\n", projectLabel)
			return nil
		}
		table := tables.NewTable()
		table.SetHeader("ID", "NAME", "POLICY")
		for _, item := range items {
			table.AddRow(
				utils.PtrString(item.Id),
				utils.PtrString(item.Name),
				utils.PtrString(item.Policy),
			)
			table.AddSeparator()
		}

		if err := table.Display(p); err != nil {
			return fmt.Errorf("render table: %w", err)
		}

		return nil
	})
}
