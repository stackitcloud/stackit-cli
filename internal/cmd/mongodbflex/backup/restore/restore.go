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
	"github.com/stackitcloud/stackit-cli/internal/pkg/spinner"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/mongodbflex"
	"github.com/stackitcloud/stackit-sdk-go/services/mongodbflex/wait"
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
			"Restores a MongoDB Flex instance from a backup of an instance or clones a MongoDB Flex instance from a point-in-time snapshot.",
			"The backup is specified by a backup ID and the point-in-time snapshot is specified by a timestamp.",
			"You can specify the instance to which the backup will be applied. If not specified, the backup will be applied to the same instance from which it was taken.",
		),
		Args: args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Restores a MongoDB Flex instance with ID "yyy" using backup with ID "zzz"`,
				`$ stackit mongodbflex backup restore --instance-id yyy --backup-id zzz`),
			examples.NewExample(
				`Clone a MongoDB Flex instance with ID "yyy" via point-in-time restore to timestamp "2024-05-14T14:31:48Z"`,
				`$ stackit mongodbflex backup restore --instance-id yyy --timestamp 2024-05-14T14:31:48Z`),
			examples.NewExample(
				`Restores a MongoDB Flex instance with ID "yyy", using backup from instance with ID "zzz" with backup ID "xxx"`,
				`$ stackit mongodbflex backup restore --instance-id zzz --backup-instance-id yyy --backup-id xxx`),
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

				if !model.Async {
					s := spinner.New(p)
					s.Start("Restoring instance")
					_, err = wait.RestoreInstanceWaitHandler(ctx, apiClient, model.ProjectId, model.InstanceId, model.BackupId).WaitWithContext(ctx)
					if err != nil {
						return fmt.Errorf("wait for MongoDB Flex instance restoration: %w", err)
					}
					s.Stop()
				}

				p.Outputf("Restored instance %q with backup %q\n", model.InstanceId, model.BackupId)
				return nil
			}

			// Else, if timestamp is provided, clone the instance from a point-in-time snapshot
			req := buildCloneRequest(ctx, model, apiClient)
			_, err = req.Execute()
			if err != nil {
				return fmt.Errorf("clone MongoDB Flex instance: %w", err)
			}

			if !model.Async {
				s := spinner.New(p)
				s.Start("Cloning instance")
				_, err = wait.CloneInstanceWaitHandler(ctx, apiClient, model.ProjectId, model.InstanceId).WaitWithContext(ctx)
				if err != nil {
					return fmt.Errorf("wait for MongoDB Flex instance cloning: %w", err)
				}
				s.Stop()
			}

			p.Outputf("Cloned instance %q from backup with timestamp %q\n", model.InstanceId, model.Timestamp)
			return nil
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), instanceIdFlag, "Instance ID")
	cmd.Flags().Var(flags.UUIDFlag(), backupInstanceIdFlag, "Instance ID of the target instance to restore the backup to")
	cmd.Flags().String(backupIdFlag, "", "Backup ID")
	cmd.Flags().String(timestampFlag, "", "Timestamp of the snapshot to use as a source for cloning the instance in a date-time with the RFC3339 layout format, e.g. 2024-01-01T00:00:00Z")

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
