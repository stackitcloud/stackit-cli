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
	volumeBackupIdArg   = "VOLUME_BACKUP_ID"
	serverIdFlag        = "server-id"
	backupIdFlag        = "backup-id"
	restoreVolumeIdFlag = "restore-volume-id"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	VolumeBackupId  string
	BackupId        string
	ServerId        string
	RestoreVolumeId string
}

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("restore %s", volumeBackupIdArg),
		Short: "Restore a Server Volume Backup to a volume.",
		Long:  "Restore a Server Volume Backup to a volume. Operation always is async.",
		Args:  args.SingleArg(volumeBackupIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Restore a Server Volume Backup with ID "xxx" for server "zzz" and backup "bbb" to volume "rrr"`,
				"$ stackit server backup volume-backup restore xxx --server-id=zzz --backup-id=bbb --restore-volume-id=rrr"),
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
				prompt := fmt.Sprintf("Are you sure you want to restore volume backup %q? (This cannot be undone)", model.VolumeBackupId)
				err = p.PromptForConfirmation(prompt)
				if err != nil {
					return err
				}
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			err = req.Execute()
			if err != nil {
				return fmt.Errorf("restore Server Volume Backup: %w", err)
			}

			p.Info("Triggered restoring of server volume backup %q\n", model.VolumeBackupId)
			return nil
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().VarP(flags.UUIDFlag(), serverIdFlag, "s", "Server ID")
	cmd.Flags().VarP(flags.UUIDFlag(), backupIdFlag, "b", "Backup ID")
	cmd.Flags().VarP(flags.UUIDFlag(), restoreVolumeIdFlag, "r", "Restore Volume ID")

	err := flags.MarkFlagsRequired(cmd, serverIdFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	volumeBackupId := inputArgs[0]

	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		VolumeBackupId:  volumeBackupId,
		ServerId:        flags.FlagToStringValue(p, cmd, serverIdFlag),
		BackupId:        flags.FlagToStringValue(p, cmd, backupIdFlag),
		RestoreVolumeId: flags.FlagToStringValue(p, cmd, restoreVolumeIdFlag),
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

func buildRequest(ctx context.Context, model *inputModel, apiClient *serverbackup.APIClient) serverbackup.ApiRestoreVolumeBackupRequest {
	req := apiClient.RestoreVolumeBackup(ctx, model.ProjectId, model.ServerId, model.Region, model.BackupId, model.VolumeBackupId)
	payload := serverbackup.RestoreVolumeBackupPayload{
		RestoreVolumeId: &model.RestoreVolumeId,
	}
	req = req.RestoreVolumeBackupPayload(payload)
	return req
}
