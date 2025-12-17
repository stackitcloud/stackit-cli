package list

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/types"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/iaas/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/iaas"
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
		Use:   "list",
		Short: "Lists all STACKIT public-ip ranges",
		Long:  "Lists all STACKIT public-ip ranges.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Lists all STACKIT public-ip ranges`,
				"$ stackit public-ip ranges list",
			),
			examples.NewExample(
				`Lists all STACKIT public-ip ranges, piping to a tool like fzf for interactive selection`,
				"$ stackit public-ip ranges list -o pretty | fzf",
			),
			examples.NewExample(
				`Lists up to 10 STACKIT public-ip ranges`,
				"$ stackit public-ip ranges list --limit 10",
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
			req := apiClient.ListPublicIPRanges(ctx)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("list public IP ranges: %w", err)
			}
			publicIpRanges := utils.GetSliceFromPointer(resp.Items)

			// Truncate output
			if model.Limit != nil && len(publicIpRanges) > int(*model.Limit) {
				publicIpRanges = publicIpRanges[:*model.Limit]
			}

			return outputResult(params.Printer, model.OutputFormat, publicIpRanges)
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

func outputResult(p *print.Printer, outputFormat string, publicIpRanges []iaas.PublicNetwork) error {
	return p.OutputResult(outputFormat, publicIpRanges, func() error {
		if len(publicIpRanges) == 0 {
			p.Outputln("No public IP ranges found")
			return nil
		}

		for _, item := range publicIpRanges {
			if item.Cidr != nil && *item.Cidr != "" {
				p.Outputln(*item.Cidr)
			}
		}

		return nil
	})
}
