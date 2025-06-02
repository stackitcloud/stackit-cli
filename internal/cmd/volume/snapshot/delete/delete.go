package delete

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/iaas/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/spinner"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/stackitcloud/stackit-sdk-go/services/iaas"
	"github.com/stackitcloud/stackit-sdk-go/services/iaas/wait"
)

const (
	snapshotIdArg = "SNAPSHOT_ID"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	SnapshotId string
}

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("delete %s", snapshotIdArg),
		Short: "Deletes a snapshot",
		Long:  "Deletes a snapshot by its ID.",
		Args:  args.SingleArg(snapshotIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Delete a snapshot`,
				"$ stackit volume snapshot delete xxx-xxx-xxx"),
			examples.NewExample(
				`Delete a snapshot and wait for deletion to be completed`,
				"$ stackit volume snapshot delete xxx-xxx-xxx --async=false"),
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

			// Get snapshot name for label
			snapshotLabel := model.SnapshotId
			snapshot, err := apiClient.GetSnapshot(ctx, model.ProjectId, model.SnapshotId).Execute()
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get snapshot name: %v", err)
			} else if snapshot != nil && snapshot.Name != nil {
				snapshotLabel = *snapshot.Name
			}

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to delete snapshot %q? (This cannot be undone)", snapshotLabel)
				err = params.Printer.PromptForConfirmation(prompt)
				if err != nil {
					return err
				}
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			err = req.Execute()
			if err != nil {
				return fmt.Errorf("delete snapshot: %w", err)
			}

			// Wait for async operation, if async mode not enabled
			if !model.Async {
				s := spinner.New(params.Printer)
				s.Start("Deleting snapshot")
				_, err = wait.DeleteSnapshotWaitHandler(ctx, apiClient, model.ProjectId, model.SnapshotId).WaitWithContext(ctx)
				if err != nil {
					return fmt.Errorf("wait for snapshot deletion: %w", err)
				}
				s.Stop()
			}

			if model.Async {
				params.Printer.Info("Triggered deletion of snapshot %q\n", snapshotLabel)
			} else {
				params.Printer.Info("Deleted snapshot %q\n", snapshotLabel)
			}
			return nil
		},
	}
	return cmd
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	snapshotId := inputArgs[0]

	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		SnapshotId:      snapshotId,
	}

	if p.IsVerbosityDebug() {
		modelStr, err := print.BuildDebugStrFromInputModel(model)
		if err != nil {
			p.Debug(print.ErrorLevel, "convert model to string for debugging: %v", err)
		} else {
			p.Debug(print.DebugLevel, "parsed input values: %s", modelStr)
		}
	}

	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *iaas.APIClient) iaas.ApiDeleteSnapshotRequest {
	return apiClient.DeleteSnapshot(ctx, model.ProjectId, model.SnapshotId)
}
