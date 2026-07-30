package plans

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	vpn "github.com/stackitcloud/stackit-sdk-go/services/vpn/v1api"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/vpn/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
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
		Short: "Lists all available plans",
		Long:  "Lists all available plans.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`List all available plans`,
				"$ stackit beta vpn plans",
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
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("list vpn plans: %w", err)
			}

			// Truncate output
			if model.Limit != nil && len(resp.Plans) > int(*model.Limit) {
				resp.Plans = resp.Plans[:*model.Limit]
			}

			return outputResult(params.Printer, model.OutputFormat, resp.Plans, model.Region)
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

func buildRequest(ctx context.Context, model *inputModel, apiClient *vpn.APIClient) vpn.ApiListPlansRequest {
	return apiClient.DefaultAPI.ListPlans(ctx, model.Region)
}

func outputResult(p *print.Printer, outputFormat string, plans []vpn.Plan, region string) error {
	return p.OutputResult(outputFormat, plans, func() error {
		if len(plans) == 0 {
			p.Info("No plans found for Region %q\n", region)
			return nil
		}

		table := tables.NewTable()
		table.SetHeader("ID", "NAME", "MAX BANDWIDTH", "MAX CONNECTIONS", "SKU")

		for _, plan := range plans {
			table.AddRow(
				utils.PtrString(plan.PlanId),
				utils.PtrString(plan.Name),
				utils.PtrString(plan.MaxBandwidth),
				utils.PtrString(plan.MaxConnections),
				utils.PtrString(plan.Sku),
			)
		}
		p.Outputln(table.Render())
		return nil
	})
}
