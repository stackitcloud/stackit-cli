package list

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/projectname"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/sfs/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/sfs"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Lists all performances classes available",
		Long:  "Lists all performances classes available.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`List all performances classes`,
				"$ stackit beta sfs performance-class list",
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			model, err := parseInput(params.Printer, cmd, args)
			if err != nil {
				return fmt.Errorf("unable to parse input: %w", err)
			}

			// Configure API client
			apiClient, err := client.ConfigureClient(params.Printer, params.CliVersion)
			if err != nil {
				return err
			}

			// Call API
			resp, err := buildRequest(ctx, apiClient).Execute()
			if err != nil {
				return fmt.Errorf("list performance-class: %w", err)
			}

			// Get projectLabel
			projectLabel, err := projectname.GetProjectName(ctx, params.Printer, params.CliVersion, cmd)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get project name: %v", err)
				projectLabel = model.ProjectId
			} else if projectLabel == "" {
				projectLabel = model.ProjectId
			}

			performanceClasses := utils.GetSliceFromPointer(resp.PerformanceClasses)

			return outputResult(params.Printer, model.OutputFormat, projectLabel, performanceClasses)
		},
	}
	return cmd
}

func parseInput(p *print.Printer, cmd *cobra.Command, _ []string) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, apiClient *sfs.APIClient) sfs.ApiListPerformanceClassesRequest {
	return apiClient.ListPerformanceClasses(ctx)
}

func outputResult(p *print.Printer, outputFormat, projectLabel string, performanceClasses []sfs.PerformanceClass) error {
	return p.OutputResult(outputFormat, performanceClasses, func() error {
		if len(performanceClasses) == 0 {
			p.Outputf("No performance classes found for project %q\n", projectLabel)
			return nil
		}

		table := tables.NewTable()
		table.SetHeader("NAME", "IOPS", "THROUGHPUT")
		for _, performanceClass := range performanceClasses {
			table.AddRow(
				utils.PtrString(performanceClass.Name),
				utils.PtrString(performanceClass.Iops),
				utils.PtrString(performanceClass.Throughput),
			)
		}
		p.Outputln(table.Render())
		return nil
	})
}
