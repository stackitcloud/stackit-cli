package restore

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/serverbackup/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/serverbackup"
)

const (
	backupIdArg                 = "BACKUP_ID"
	serverIdFlag                = "server-id"
	startServerAfterRestoreFlag = "start-server-after-restore"
	backupVolumeIdsFlag         = "volume-ids"

	defaultStartServerAfterRestore = false
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	BackupId                string
	ServerId                string
	StartServerAfterRestore bool
	BackupVolumeIds         []string
}

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("restore %s", backupIdArg),
		Short: "Restores a Server Backup.",
		Long:  "Restores a Server Backup. Operation always is async.",
		Args:  args.SingleArg(backupIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Restore a Server Backup with ID "xxx" for server "zzz"`,
				"$ stackit beta server backup restore xxx --server-id=zzz"),
			examples.NewExample(
				`Restore a Server Backup with ID "xxx" for server "zzz" and start the server afterwards`,
				"$ stackit beta server backup restore xxx --server-id=zzz --start-server-after-restore"),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			model, err := parseInput(p, cmd, args)
			if err != nil {
				return err
			}

			// Configure API client
			apiClient, err := client.ConfigureClient(p)
			if err != nil {
				return err
			}

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to restore server backup %q? (This cannot be undone)", model.BackupId)
				err = p.PromptForConfirmation(prompt)
				if err != nil {
					return err
				}
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			err = req.Execute()
			if err != nil {
				return fmt.Errorf("restore Server Backup: %w", err)
			}

			p.Info("Triggered restoring of server backup %q\n", model.BackupId)
			return nil
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().VarP(flags.UUIDFlag(), serverIdFlag, "s", "Server ID")
	cmd.Flags().VarP(flags.UUIDSliceFlag(), backupVolumeIdsFlag, "i", "Backup volume IDs, as comma separated UUID values.")
	cmd.Flags().BoolP(startServerAfterRestoreFlag, "u", defaultStartServerAfterRestore, "Should the server start after the backup restoring.")

	err := flags.MarkFlagsRequired(cmd, serverIdFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	backupId := inputArgs[0]

	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel:         globalFlags,
		BackupId:                backupId,
		ServerId:                flags.FlagToStringValue(p, cmd, serverIdFlag),
		BackupVolumeIds:         flags.FlagToStringSliceValue(p, cmd, backupVolumeIdsFlag),
		StartServerAfterRestore: flags.FlagToBoolValue(p, cmd, startServerAfterRestoreFlag),
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

func buildRequest(ctx context.Context, model *inputModel, apiClient *serverbackup.APIClient) serverbackup.ApiRestoreBackupRequest {
	req := apiClient.RestoreBackup(ctx, model.ProjectId, model.ServerId, model.BackupId)
	payload := serverbackup.RestoreBackupPayload{
		StartServerAfterRestore: &model.StartServerAfterRestore,
		VolumeIds:               &model.BackupVolumeIds,
	}
	if model.BackupVolumeIds == nil {
		payload.VolumeIds = nil
	}
	req = req.RestoreBackupPayload(payload)
	return req
}
