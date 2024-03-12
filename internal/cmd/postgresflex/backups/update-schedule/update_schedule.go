package updateschedule

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/pkg/confirm"
	cliErr "github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/postgresflex/client"
	postgresflexUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/postgresflex/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/postgresflex"
)

const (
	instanceIdArg = "INSTANCE_ID"

	backupScheduleFlag = "backup-schedule"
)

type inputModel struct {
	*globalflags.GlobalFlagModel

	InstanceId     string
	BackupSchedule *string
}

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("update-schedule %s", instanceIdArg),
		Short: "Updates backup schedule for a specific PostgreSQL Flex instance",
		Long:  "Updates backup schedule for a specific PostgreSQL Flex instance.",
		Args:  args.SingleArg(instanceIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Update the backup schedule of a PostgreSQL Flex instance with ID "xxx"`,
				"$ stackit postgresflex backups update-schedule xxx --backup-schedule '6 6 * * *'"),
		),

		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()

			model, err := parseInput(cmd, args)
			if err != nil {
				return err
			}

			// Configure API client
			apiClient, err := client.ConfigureClient(cmd)
			if err != nil {
				return err
			}

			instanceLabel, err := postgresflexUtils.GetInstanceName(ctx, apiClient, model.ProjectId, model.InstanceId)
			if err != nil {
				instanceLabel = model.InstanceId
			}

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to update backup schedule of instance %q?", instanceLabel)
				err = confirm.PromptForConfirmation(cmd, prompt)
				if err != nil {
					return err
				}
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			err = req.Execute()
			if err != nil {
				return fmt.Errorf("update backup schedule of PostgreSQL Flex instance: %w", err)
			}

			cmd.Printf("Updated backup schedule of instance %q\n", instanceLabel)
			return nil
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().String(backupScheduleFlag, "", "Backup schedule")

	err := flags.MarkFlagsRequired(cmd, backupScheduleFlag)
	cobra.CheckErr(err)
}

func parseInput(cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	instanceId := inputArgs[0]

	globalFlags := globalflags.Parse(cmd)
	if globalFlags.ProjectId == "" {
		return nil, &cliErr.ProjectIdError{}
	}

	return &inputModel{
		GlobalFlagModel: globalFlags,
		InstanceId:      instanceId,
		BackupSchedule:  flags.FlagToStringPointer(cmd, backupScheduleFlag),
	}, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *postgresflex.APIClient) postgresflex.ApiUpdateBackupScheduleRequest {
	req := apiClient.UpdateBackupSchedule(ctx, model.ProjectId, model.InstanceId)
	req = req.UpdateBackupSchedulePayload(postgresflex.UpdateBackupSchedulePayload{
		BackupSchedule: model.BackupSchedule,
	})
	return req
}
