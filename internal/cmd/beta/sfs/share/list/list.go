package list

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/sfs/client"
	sfsUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/sfs/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/sfs"
)

const (
	resourcePoolIdFlag = "resource-pool-id"
	limitFlag          = "limit"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	ResourcePoolId string
	Limit          *int64
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Lists all shares of a resource pool",
		Long:  "Lists all shares of a resource pool.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`List all shares from resource pool with ID "xxx"`,
				"$ stackit beta sfs export-policy list --resource-pool-id xxx",
			),
			examples.NewExample(
				`List up to 10 shares from resource pool with ID "xxx"`,
				"$ stackit beta sfs export-policy list --resource-pool-id xxx --limit 10",
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
				return fmt.Errorf("list SFS share: %w", err)
			}

			resourcePoolLabel, err := sfsUtils.GetResourcePoolName(ctx, apiClient, model.ProjectId, model.Region, model.ResourcePoolId)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get resource pool name: %v", err)
				resourcePoolLabel = model.ResourcePoolId
			} else if resourcePoolLabel == "" {
				resourcePoolLabel = model.ResourcePoolId
			}

			// Truncate output
			items := utils.GetSliceFromPointer(resp.Shares)
			if model.Limit != nil && len(items) > int(*model.Limit) {
				items = items[:*model.Limit]
			}

			return outputResult(params.Printer, model.OutputFormat, resourcePoolLabel, items)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), resourcePoolIdFlag, "The resource pool the share is assigned to")
	cmd.Flags().Int64(limitFlag, 0, "Maximum number of entries to list")

	err := flags.MarkFlagsRequired(cmd, resourcePoolIdFlag)
	cobra.CheckErr(err)
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
			Details: "must be grater than 0",
		}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		ResourcePoolId:  flags.FlagToStringValue(p, cmd, resourcePoolIdFlag),
		Limit:           limit,
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *sfs.APIClient) sfs.ApiListSharesRequest {
	return apiClient.ListShares(ctx, model.ProjectId, model.Region, model.ResourcePoolId)
}

func outputResult(p *print.Printer, outputFormat, resourcePoolLabel string, shares []sfs.Share) error {
	return p.OutputResult(outputFormat, shares, func() error {
		if len(shares) == 0 {
			p.Info("No shares found for resource pool %q\n", resourcePoolLabel)
			return nil
		}

		table := tables.NewTable()
		table.SetHeader("ID", "NAME", "STATE", "EXPORT POLICY", "MOUNT PATH", "HARD LIMIT (GB)", "CREATED AT")

		for _, share := range shares {
			var policy string
			if share.ExportPolicy != nil {
				if name, ok := share.ExportPolicy.Get().GetNameOk(); ok {
					policy = name
				} else if id, ok := share.ExportPolicy.Get().GetIdOk(); ok {
					policy = id
				}
			}
			table.AddRow(
				utils.PtrString(share.Id),
				utils.PtrString(share.Name),
				utils.PtrString(share.State),
				policy,
				utils.PtrString(share.MountPath),
				utils.PtrString(share.SpaceHardLimitGigabytes),
				utils.ConvertTimePToDateTimeString(share.CreatedAt),
			)
		}
		p.Outputln(table.Render())
		return nil
	})
}
