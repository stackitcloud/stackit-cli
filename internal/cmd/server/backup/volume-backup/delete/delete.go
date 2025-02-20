package delete

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
	volumeBackupIdArg = "VOLUME_BACKUP_ID"
	serverIdFlag      = "server-id"
	backupIdFlag      = "backup-id"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	BackupId string
	VolumeId string
	ServerId string
}

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("delete %s", volumeBackupIdArg),
		Short: "Deletes a Server Volume Backup.",
		Long:  "Deletes a Server Volume Backup. Operation always is async.",
		Args:  args.SingleArg(volumeBackupIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Delete a Server Volume Backup with ID "xxx" for server "zzz" and backup "bbb"`,
				"$ stackit server backup volume-backup delete xxx --server-id=zzz --backup-id=bbb"),
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
				prompt := fmt.Sprintf("Are you sure you want to delete server volume backup %q? (This cannot be undone)", model.VolumeId)
				err = p.PromptForConfirmation(prompt)
				if err != nil {
					return err
				}
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			err = req.Execute()
			if err != nil {
				return fmt.Errorf("delete Server Volume Backup: %w", err)
			}

			p.Info("Triggered deletion of server volume backup %q\n", model.VolumeId)
			return nil
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().VarP(flags.UUIDFlag(), serverIdFlag, "s", "Server ID")
	cmd.Flags().VarP(flags.UUIDFlag(), backupIdFlag, "b", "Backup ID")

	err := flags.MarkFlagsRequired(cmd, serverIdFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	volumeId := inputArgs[0]

	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		VolumeId:        volumeId,
		BackupId:        flags.FlagToStringValue(p, cmd, backupIdFlag),
		ServerId:        flags.FlagToStringValue(p, cmd, serverIdFlag),
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

func buildRequest(ctx context.Context, model *inputModel, apiClient *serverbackup.APIClient) serverbackup.ApiDeleteVolumeBackupRequest {
	req := apiClient.DeleteVolumeBackup(ctx, model.ProjectId, model.ServerId, model.BackupId, model.VolumeId)
	return req
}
