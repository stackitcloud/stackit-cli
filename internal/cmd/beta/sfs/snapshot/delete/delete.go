package delete

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
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
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
		Use:   fmt.Sprintf("delete %s", snapshotNameArg),
		Short: "Deletes a snapshot",
		Long:  "Deletes a snapshot.",
		Args:  args.SingleArg(snapshotNameArg, nil),
		Example: examples.Build(
			examples.NewExample(
				`Delete a snapshot with "SNAPSHOT_NAME" from resource pool with ID "yyy"`,
				"$ stackit beta sfs snapshot delete SNAPSHOT_NAME --resource-pool-id yyy"),
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

			resourcePoolLabel, err := sfsUtils.GetResourcePoolName(ctx, apiClient, model.ProjectId, model.Region, model.ResourcePoolId)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get resource pool name: %v", err)
				resourcePoolLabel = model.ResourcePoolId
			} else if resourcePoolLabel == "" {
				resourcePoolLabel = model.ResourcePoolId
			}

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to delete snapshot %q for resource pool %q?", model.SnapshotName, resourcePoolLabel)
				err = params.Printer.PromptForConfirmation(prompt)
				if err != nil {
					return err
				}
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			_, err = req.Execute()
			if err != nil {
				return fmt.Errorf("delete snapshot: %w", err)
			}

			params.Printer.Outputf("Deleted snapshot %q from resource pool %q.\n", model.SnapshotName, resourcePoolLabel)
			return nil
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

func buildRequest(ctx context.Context, model *inputModel, apiClient *sfs.APIClient) sfs.ApiDeleteResourcePoolSnapshotRequest {
	return apiClient.DeleteResourcePoolSnapshot(ctx, model.ProjectId, model.Region, model.ResourcePoolId, model.SnapshotName)
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
