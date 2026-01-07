package describe

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
	snapshotNameArg = "SNAPSHOT_NAME"

	resourcePoolIdFlag = "resource-pool-id"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	ResourcePoolId string
	SnapshotName   string
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("describe %s", snapshotNameArg),
		Short: "Shows details of a snapshot",
		Long:  "Shows details of a snapshot.",
		Args:  args.SingleArg(snapshotNameArg, nil),
		Example: examples.Build(
			examples.NewExample(
				`Describe a snapshot with "SNAPSHOT_NAME" from resource pool with ID "yyy"`,
				"stackit beta sfs snapshot describe SNAPSHOT_NAME --resource-pool-id yyy",
			),
		),
		RunE: func(cmd *cobra.Command, inputArgs []string) error {
			ctx := context.Background()
			model, err := parseInput(params.Printer, cmd, inputArgs)
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
				return fmt.Errorf("create snapshot: %w", err)
			}

			return outputResult(params.Printer, model.OutputFormat, resp)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), resourcePoolIdFlag, "The resource pool from which the snapshot should be created")

	err := flags.MarkFlagsRequired(cmd, resourcePoolIdFlag)
	cobra.CheckErr(err)
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *sfs.APIClient) sfs.ApiGetResourcePoolSnapshotRequest {
	return apiClient.GetResourcePoolSnapshot(ctx, model.ProjectId, model.Region, model.ResourcePoolId, model.SnapshotName)
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	snapshotName := inputArgs[0]

	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		SnapshotName:    snapshotName,
		ResourcePoolId:  flags.FlagToStringValue(p, cmd, resourcePoolIdFlag),
	}

	p.DebugInputModel(model)
	return &model, nil
}

func outputResult(p *print.Printer, outputFormat string, resp *sfs.GetResourcePoolSnapshotResponse) error {
	return p.OutputResult(outputFormat, resp, func() error {
		if resp == nil || resp.ResourcePoolSnapshot == nil {
			p.Outputln("Resource pool snapshot response is empty")
			return nil
		}

		table := tables.NewTable()

		snap := *resp.ResourcePoolSnapshot
		table.AddRow("NAME", utils.PtrString(snap.SnapshotName))
		table.AddSeparator()
		if snap.Comment != nil {
			table.AddRow("COMMENT", utils.PtrString(snap.Comment.Get()))
			table.AddSeparator()
		}
		table.AddRow("RESOURCE POOL ID", utils.PtrString(snap.ResourcePoolId))
		table.AddSeparator()
		table.AddRow("SIZE (GB)", utils.PtrString(snap.SizeGigabytes))
		table.AddSeparator()
		table.AddRow("LOGICAL SIZE (GB)", utils.PtrString(snap.LogicalSizeGigabytes))
		table.AddSeparator()
		table.AddRow("CREATED AT", utils.ConvertTimePToDateTimeString(snap.CreatedAt))
		table.AddSeparator()

		p.Outputln(table.Render())
		return nil
	})
}
