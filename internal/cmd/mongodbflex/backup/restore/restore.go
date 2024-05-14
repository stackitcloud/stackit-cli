package restore

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	cliErr "github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/mongodbflex/client"
	mongodbUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/mongodbflex/utils"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/mongodbflex"
)

const (
	instanceIdFlag       = "instance-id"
	backupInstanceIdFlag = "backup-instance-id"
	backupIdFlag         = "backup-id"
	timestampFlag        = "timestamp"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	InstanceId       string
	BackupInstanceId string
	BackupId         string
	Timestamp        string
}

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "restore",
		Short: "Restores a MongoDB Flex instance from a backup",
		Long: fmt.Sprintf("%s\n%s\n%s",
			"Restores a MongoDB Flex instance from a backup of an instance.",
			"The backup can be specified by either a backup id or a timestamp.",
			"The target instance can be specified, otherwise the target instance will be the same as the backup.",
		),
		Args: args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Restores a MongoDB Flex instance with id "yyy" using backup with id "zzz"`,
				`$ stackit mongodbflex backup restore --instance-id yyy --backup-id zzz`),
			examples.NewExample(
				`Restores a MongoDB Flex instance with id "yyy" using backup with timestamp "zzz"`,
				`$ stackit mongodbflex backup restore --instance-id yyy --timestamp zzz`),
			examples.NewExample(
				`Restores a MongoDB Flex instance with id "yyy" using backup from instance with id "zzz" with backup id "aaa"`,
				`$ stackit mongodbflex backup restore --instance-id yyy --backup-instance-id zzz --backup-id aaa`),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()

			model, err := parseInput(p, cmd)
			if err != nil {
				return err
			}

			// Configure API client
			apiClient, err := client.ConfigureClient(p)
			if err != nil {
				return err
			}

			instanceLabel, err := mongodbUtils.GetInstanceName(ctx, apiClient, model.ProjectId, model.InstanceId)
			if err != nil {
				p.Debug(print.ErrorLevel, "get instance name: %v", err)
				instanceLabel = model.ProjectId
			}

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to restore MongoDB Flex instance %q?", instanceLabel)
				err = p.PromptForConfirmation(prompt)
				if err != nil {
					return err
				}
			}

			// If backupInstanceId is not provided, the target is the same instance as the backup
			if model.BackupInstanceId == "" {
				model.BackupInstanceId = model.InstanceId
			}

			isRestoreOperation := getIsRestoreOperation(model.BackupId, model.Timestamp)

			// If backupId is provided, restore the instance from the backup with the backupId
			if isRestoreOperation {
				req := buildRestoreRequest(ctx, model, apiClient)
				_, err = req.Execute()
				if err != nil {
					return fmt.Errorf("restore MongoDB Flex instance: %w", err)
				}

				p.Outputf("Restored instance %q with backup %q\n", model.InstanceId, model.BackupId)
				return nil
			}

			// Else, if timestamp is provided, clone the instance from the backup with the timestep
			req := buildCloneRequest(ctx, model, apiClient)
			_, err = req.Execute()
			if err != nil {
				return fmt.Errorf("clone MongoDB Flex instance: %w", err)
			}
			p.Outputf("Cloned instance %q from backup with timestamp %q\n", model.InstanceId, model.Timestamp)
			return nil
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), instanceIdFlag, "Instance id")
	cmd.Flags().Var(flags.UUIDFlag(), backupInstanceIdFlag, "Backup instance id")
	cmd.Flags().String(backupIdFlag, "", "Backup id")
	cmd.Flags().String(timestampFlag, "", "Timestamp of the backup")

	err := flags.MarkFlagsRequired(cmd, instanceIdFlag)
	cobra.CheckErr(err)

	cmd.MarkFlagsOneRequired(backupIdFlag, timestampFlag)
	cmd.MarkFlagsMutuallyExclusive(backupIdFlag, timestampFlag)
}

func parseInput(p *print.Printer, cmd *cobra.Command) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &cliErr.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel:  globalFlags,
		InstanceId:       flags.FlagToStringValue(p, cmd, instanceIdFlag),
		BackupInstanceId: flags.FlagToStringValue(p, cmd, backupInstanceIdFlag),
		BackupId:         flags.FlagToStringValue(p, cmd, backupIdFlag),
		Timestamp:        flags.FlagToStringValue(p, cmd, timestampFlag),
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

func buildRestoreRequest(ctx context.Context, model *inputModel, apiClient *mongodbflex.APIClient) mongodbflex.ApiRestoreInstanceRequest {
	req := apiClient.RestoreInstance(ctx, model.ProjectId, model.InstanceId)
	req = req.RestoreInstancePayload(mongodbflex.RestoreInstancePayload{
		BackupId:   &model.BackupId,
		InstanceId: &model.BackupInstanceId,
	})
	return req
}

func buildCloneRequest(ctx context.Context, model *inputModel, apiClient *mongodbflex.APIClient) mongodbflex.ApiCloneInstanceRequest {
	req := apiClient.CloneInstance(ctx, model.ProjectId, model.InstanceId)
	req = req.CloneInstancePayload(mongodbflex.CloneInstancePayload{
		Timestamp:  &model.Timestamp,
		InstanceId: &model.BackupInstanceId,
	})
	return req
}

func getIsRestoreOperation(backupId, timestamp string) bool {
	return backupId != "" && timestamp == ""
}
