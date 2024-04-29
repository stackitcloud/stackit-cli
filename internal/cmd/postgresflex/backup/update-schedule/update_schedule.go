package updateschedule

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	cliErr "github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/postgresflex/client"
	postgresflexUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/postgresflex/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/postgresflex"
)

const (
	instanceIdFlag = "instance-id"
	scheduleFlag   = "schedule"
)

type inputModel struct {
	*globalflags.GlobalFlagModel

	InstanceId     *string
	BackupSchedule *string
}

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-schedule",
		Short: "Updates backup schedule for a PostgreSQL Flex instance",
		Long:  "Updates backup schedule for a PostgreSQL Flex instance. The current backup schedule can be seen in the output of instance describe command.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Update the backup schedule of a PostgreSQL Flex instance with ID "xxx"`,
				"$ stackit postgresflex backup update-schedule --instance-id xxx --schedule '6 6 * * *'"),
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

			instanceLabel, err := postgresflexUtils.GetInstanceName(ctx, apiClient, model.ProjectId, *model.InstanceId)
			if err != nil {
				p.Debug(print.ErrorLevel, "get instance name: %v", err)
				instanceLabel = *model.InstanceId
			}

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to update backup schedule of instance %q?", instanceLabel)
				err = p.PromptForConfirmation(prompt)
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
	cmd.Flags().Var(flags.UUIDFlag(), instanceIdFlag, "Instance ID")
	cmd.Flags().String(scheduleFlag, "", "Backup schedule, in the cron scheduling system format e.g. '0 0 * * *'")

	err := flags.MarkFlagsRequired(cmd, instanceIdFlag, scheduleFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &cliErr.ProjectIdError{}
	}

	return &inputModel{
		GlobalFlagModel: globalFlags,
		InstanceId:      flags.FlagToStringPointer(p, cmd, instanceIdFlag),
		BackupSchedule:  flags.FlagToStringPointer(p, cmd, scheduleFlag),
	}, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *postgresflex.APIClient) postgresflex.ApiUpdateBackupScheduleRequest {
	req := apiClient.UpdateBackupSchedule(ctx, model.ProjectId, *model.InstanceId)
	req = req.UpdateBackupSchedulePayload(postgresflex.UpdateBackupSchedulePayload{
		BackupSchedule: model.BackupSchedule,
	})
	return req
}
