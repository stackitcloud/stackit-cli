package plans

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/types"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/projectname"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/alb/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/alb"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "plans",
		Short: "Lists the application load balancer plans",
		Long:  "Lists the available application load balancer plans.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`List all application load balancer plans`,
				`$ stackit beta alb plans`,
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

			projectLabel, err := projectname.GetProjectName(ctx, params.Printer, params.CliVersion, cmd)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get project name: %v", err)
				projectLabel = model.ProjectId
			} else if projectLabel == "" {
				projectLabel = model.ProjectId
			}

			// Call API
			request := buildRequest(ctx, model, apiClient)

			response, err := request.Execute()
			if err != nil {
				return fmt.Errorf("list plans: %w", err)
			}
			items := response.GetValidPlans()

			return outputResult(params.Printer, model.OutputFormat, projectLabel, items)
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

func buildRequest(ctx context.Context, model *inputModel, apiClient *alb.APIClient) alb.ApiListPlansRequest {
	request := apiClient.ListPlans(ctx, model.Region)

	return request
}

func outputResult(p *print.Printer, outputFormat, projectLabel string, items []alb.PlanDetails) error {
	return p.OutputResult(outputFormat, items, func() error {
		if len(items) == 0 {
			p.Outputf("No plans found for project %q", projectLabel)
			return nil
		}

		table := tables.NewTable()
		table.SetHeader("PLAN ID", "NAME", "FLAVOR", "MAX CONNS", "DESCRIPTION")
		for _, item := range items {
			table.AddRow(utils.PtrString(item.PlanId),
				utils.PtrString(item.Name),
				utils.PtrString(item.FlavorName),
				utils.PtrString(item.MaxConnections),
				utils.Truncate(item.Description, 70),
			)
		}
		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}

		return nil
	})
}
