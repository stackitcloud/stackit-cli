package restore

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
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/iaas/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/spinner"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	iaasutils "github.com/stackitcloud/stackit-cli/internal/pkg/services/iaas/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/iaas"
	"github.com/stackitcloud/stackit-sdk-go/services/iaas/wait"
)

const (
	backupIdArg = "BACKUP_ID"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	BackupId string
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("restore %s", backupIdArg),
		Short: "Restores a backup",
		Long:  "Restores a backup by its ID.",
		Args:  args.SingleArg(backupIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Restore a backup with ID "xxx"`, "$ stackit volume backup restore xxx"),
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

			backupLabel, err := iaasutils.GetBackupName(ctx, apiClient, model.ProjectId, model.Region, model.BackupId)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get backup details: %v", err)
			}

			// Get source details for labels
			var sourceLabel string
			backup, err := apiClient.GetBackup(ctx, model.ProjectId, model.Region, model.BackupId).Execute()
			if err == nil && backup != nil && backup.VolumeId != nil {
				sourceLabel = *backup.VolumeId
				name, err := iaasutils.GetVolumeName(ctx, apiClient, model.ProjectId, model.Region, *backup.VolumeId)
				if err != nil {
					params.Printer.Debug(print.ErrorLevel, "get volume details: %v", err)
				} else if name != "" {
					sourceLabel = name
				}
			}

			prompt := fmt.Sprintf("Are you sure you want to restore %q with backup %q? (This cannot be undone)", sourceLabel, backupLabel)
			err = params.Printer.PromptForConfirmation(prompt)
			if err != nil {
				return err
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			err = req.Execute()
			if err != nil {
				return fmt.Errorf("restore backup: %w", err)
			}

			// Wait for async operation, if async mode not enabled
			if !model.Async {
				s := spinner.New(params.Printer)
				s.Start("Restoring backup")
				_, err = wait.RestoreBackupWaitHandler(ctx, apiClient, model.ProjectId, model.Region, model.BackupId).WaitWithContext(ctx)
				if err != nil {
					return fmt.Errorf("wait for backup restore: %w", err)
				}
				s.Stop()
			}

			if model.Async {
				params.Printer.Outputf("Triggered restore of %q with %q in %q\n", sourceLabel, backupLabel, model.ProjectId)
			} else {
				params.Printer.Outputf("Restored %q with %q in %q\n", sourceLabel, backupLabel, model.ProjectId)
			}
			return nil
		},
	}
	return cmd
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	backupId := inputArgs[0]

	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		BackupId:        backupId,
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *iaas.APIClient) iaas.ApiRestoreBackupRequest {
	req := apiClient.RestoreBackup(ctx, model.ProjectId, model.Region, model.BackupId)
	return req
}
