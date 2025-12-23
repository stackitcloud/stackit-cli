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
		Short: "Lists all snapshots of a resource pool",
		Long:  "Lists all snapshots of a resource pool.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`List all snapshots of a resource pool with ID "xxx"`,
				"$ stackit beta sfs snapshot list --resource-pool-id xxx",
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
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("list snapshot: %w", err)
			}

			// Truncate output
			items := utils.GetSliceFromPointer(resp.ResourcePoolSnapshots)
			if model.Limit != nil && len(items) > int(*model.Limit) {
				items = items[:*model.Limit]
			}

			return outputResult(params.Printer, model.OutputFormat, items)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), resourcePoolIdFlag, "The resource pool from which the snapshot should be created")
	cmd.Flags().Int64(limitFlag, 0, "Number of snapshots to list")

	err := flags.MarkFlagsRequired(cmd, resourcePoolIdFlag)
	cobra.CheckErr(err)
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *sfs.APIClient) sfs.ApiListResourcePoolSnapshotsRequest {
	req := apiClient.ListResourcePoolSnapshots(ctx, model.ProjectId, model.Region, model.ResourcePoolId)
	return req
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
		ResourcePoolId:  flags.FlagToStringValue(p, cmd, resourcePoolIdFlag),
		Limit:           limit,
	}

	p.DebugInputModel(model)
	return &model, nil
}

func outputResult(p *print.Printer, outputFormat string, resp []sfs.ResourcePoolSnapshot) error {
	return p.OutputResult(outputFormat, resp, func() error {
		if len(resp) == 0 {
			p.Outputln("No snapshots found")
			return nil
		}
		table := tables.NewTable()
		table.SetHeader(
			"NAME",
			"COMMENT",
			"RESOURCE POOL ID",
			"SIZE (GB)",
			"LOGICAL SIZE (GB)",
			"CREATED AT",
		)

		for _, snap := range resp {
			var comment string
			if snap.Comment != nil {
				comment = utils.PtrString(snap.Comment.Get())
			}
			table.AddRow(
				utils.PtrString(snap.SnapshotName),
				comment,
				utils.PtrString(snap.ResourcePoolId),
				utils.PtrString(snap.SizeGigabytes),
				utils.PtrString(snap.LogicalSizeGigabytes),
				utils.ConvertTimePToDateTimeString(snap.CreatedAt),
			)
		}

		p.Outputln(table.Render())
		return nil
	})
}
