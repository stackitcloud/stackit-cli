package updateschedule

import (
	"context"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	cliErr "github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/mongodbflex/client"
	mongoDBflexUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/mongodbflex/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/mongodbflex"
)

const (
	instanceIdFlag                     = "instance-id"
	scheduleFlag                       = "schedule"
	snapshotRetentionDaysFlag          = "save-snapshot-days"
	dailySnapshotRetentionDaysFlag     = "save-daily-snapshot-days"
	weeklySnapshotRetentionWeeksFlag   = "save-weekly-snapshot-weeks"
	monthlySnapshotRetentionMonthsFlag = "save-monthly-snapshot-months"

	// Default values for the backup schedule options
	defaultBackupSchedule                       = "0 0/6 * * *"
	defaultSnapshotRetentionDays          int64 = 3
	defaultDailySnapshotRetentionDays     int64 = 0
	defaultWeeklySnapshotRetentionWeeks   int64 = 3
	defaultMonthlySnapshotRetentionMonths int64 = 1
	defaultPointInTimeWindowHours         int64 = 30
)

type inputModel struct {
	*globalflags.GlobalFlagModel

	InstanceId                     *string
	BackupSchedule                 *string
	SnapshotRetentionDays          *int64
	DailySnaphotRetentionDays      *int64
	WeeklySnapshotRetentionWeeks   *int64
	MonthlySnapshotRetentionMonths *int64
}

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-schedule",
		Short: "Updates the backup schedule and retention policy for a MongoDB Flex instance",
		Long: fmt.Sprintf("%s\n%s\n%s\n%s",
			"Updates the backup schedule and retention policy for a MongoDB Flex instance.",
			`The current backup schedule and retention policy can be seen in the output of the "stackit mongodbflex backup schedule" command.`,
			"The backup schedule is defined in the cron scheduling system format e.g. '0 0 * * *'.",
			"See below for more detail on the retention policy options.",
		),
		Args: args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Update the backup schedule of a MongoDB Flex instance with ID "xxx"`,
				"$ stackit mongodbflex backup update-schedule --instance-id xxx --schedule '6 6 * * *'"),
			examples.NewExample(
				`Update the retention days for snapshots of a MongoDB Flex instance with ID "xxx" to 5 days`,
				"$ stackit mongodbflex backup update-schedule --instance-id xxx --save-snapshot-days 5"),
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

			instanceLabel, err := mongoDBflexUtils.GetInstanceName(ctx, apiClient, model.ProjectId, *model.InstanceId)
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

			// Get current instance
			getReq := buildGetInstanceRequest(ctx, model, apiClient)
			getResp, err := getReq.Execute()
			if err != nil {
				return fmt.Errorf("get MongoDB Flex instance %q: %w", instanceLabel, err)
			}

			instance := getResp.Item

			// Call API
			req := buildUpdateBackupScheduleRequest(ctx, model, instance, apiClient)
			_, err = req.Execute()
			if err != nil {
				return fmt.Errorf("update backup schedule of MongoDB Flex instance: %w", err)
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
	cmd.Flags().Int64(snapshotRetentionDaysFlag, 0, "Number of days to retain snapshots. Should be less than or equal to the value of the daily backup.")
	cmd.Flags().Int64(dailySnapshotRetentionDaysFlag, 0, "Number of days to retain daily snapshots. Should be less than or equal to the number of days of the selected weekly or monthly value.")
	cmd.Flags().Int64(weeklySnapshotRetentionWeeksFlag, 0, "Number of weeks to retain weekly snapshots. Should be less than or equal to the number of weeks of the selected monthly value.")
	cmd.Flags().Int64(monthlySnapshotRetentionMonthsFlag, 0, "Number of months to retain monthly snapshots")

	err := flags.MarkFlagsRequired(cmd, instanceIdFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &cliErr.ProjectIdError{}
	}

	schedule := flags.FlagToStringPointer(p, cmd, scheduleFlag)
	snapshotRetentionDays := flags.FlagToInt64Pointer(p, cmd, snapshotRetentionDaysFlag)
	dailySnapshotRetentionDays := flags.FlagToInt64Pointer(p, cmd, dailySnapshotRetentionDaysFlag)
	weeklySnapshotRetentionWeeks := flags.FlagToInt64Pointer(p, cmd, weeklySnapshotRetentionWeeksFlag)
	monthlySnapshotRetentionMonths := flags.FlagToInt64Pointer(p, cmd, monthlySnapshotRetentionMonthsFlag)

	if schedule == nil && snapshotRetentionDays == nil && dailySnapshotRetentionDays == nil && weeklySnapshotRetentionWeeks == nil && monthlySnapshotRetentionMonths == nil {
		return nil, &cliErr.EmptyUpdateError{}
	}

	return &inputModel{
		GlobalFlagModel:                globalFlags,
		InstanceId:                     flags.FlagToStringPointer(p, cmd, instanceIdFlag),
		BackupSchedule:                 schedule,
		DailySnaphotRetentionDays:      dailySnapshotRetentionDays,
		MonthlySnapshotRetentionMonths: monthlySnapshotRetentionMonths,
		SnapshotRetentionDays:          snapshotRetentionDays,
		WeeklySnapshotRetentionWeeks:   weeklySnapshotRetentionWeeks,
	}, nil
}

func buildUpdateBackupScheduleRequest(ctx context.Context, model *inputModel, instance *mongodbflex.Instance, apiClient *mongodbflex.APIClient) mongodbflex.ApiUpdateBackupScheduleRequest {
	req := apiClient.UpdateBackupSchedule(ctx, model.ProjectId, *model.InstanceId)

	payload := getUpdateBackupSchedulePayload(instance)

	if model.BackupSchedule != nil {
		payload.BackupSchedule = model.BackupSchedule
	}
	if model.DailySnaphotRetentionDays != nil {
		payload.DailySnapshotRetentionDays = model.DailySnaphotRetentionDays
	}
	if model.MonthlySnapshotRetentionMonths != nil {
		payload.MonthlySnapshotRetentionMonths = model.MonthlySnapshotRetentionMonths
	}
	if model.SnapshotRetentionDays != nil {
		payload.SnapshotRetentionDays = model.SnapshotRetentionDays
	}
	if model.WeeklySnapshotRetentionWeeks != nil {
		payload.WeeklySnapshotRetentionWeeks = model.WeeklySnapshotRetentionWeeks
	}

	req = req.UpdateBackupSchedulePayload(payload)
	return req
}

// getUpdateBackupSchedulePayload creates a payload for the UpdateBackupSchedule API call
// it will use the values already set in the instance object
// falls back to default values if the values are not set
func getUpdateBackupSchedulePayload(instance *mongodbflex.Instance) mongodbflex.UpdateBackupSchedulePayload {
	options := make(map[string]string)
	if instance == nil || instance.Options != nil {
		options = *instance.Options
	}

	backupSchedule := instance.BackupSchedule
	if backupSchedule == nil {
		backupSchedule = utils.Ptr(defaultBackupSchedule)
	}
	dailySnapshotRetentionDays, err := strconv.ParseInt(options["dailySnapshotRetentionDays"], 10, 64)
	if err != nil {
		dailySnapshotRetentionDays = defaultDailySnapshotRetentionDays
	}
	weeklySnapshotRetentionWeeks, err := strconv.ParseInt(options["weeklySnapshotRetentionWeeks"], 10, 64)
	if err != nil {
		weeklySnapshotRetentionWeeks = defaultWeeklySnapshotRetentionWeeks
	}
	monthlySnapshotRetentionMonths, err := strconv.ParseInt(options["monthlySnapshotRetentionMonths"], 10, 64)
	if err != nil {
		monthlySnapshotRetentionMonths = defaultMonthlySnapshotRetentionMonths
	}
	pointInTimeWindowHours, err := strconv.ParseInt(options["pointInTimeWindowHours"], 10, 64)
	if err != nil {
		pointInTimeWindowHours = defaultPointInTimeWindowHours
	}
	snapshotRetentionDays, err := strconv.ParseInt(options["snapshotRetentionDays"], 10, 64)
	if err != nil {
		snapshotRetentionDays = defaultSnapshotRetentionDays
	}

	defaultPayload := mongodbflex.UpdateBackupSchedulePayload{
		BackupSchedule:                 backupSchedule,
		DailySnapshotRetentionDays:     &dailySnapshotRetentionDays,
		MonthlySnapshotRetentionMonths: &monthlySnapshotRetentionMonths,
		PointInTimeWindowHours:         &pointInTimeWindowHours,
		SnapshotRetentionDays:          &snapshotRetentionDays,
		WeeklySnapshotRetentionWeeks:   &weeklySnapshotRetentionWeeks,
	}
	return defaultPayload
}

func buildGetInstanceRequest(ctx context.Context, model *inputModel, apiClient *mongodbflex.APIClient) mongodbflex.ApiGetInstanceRequest {
	req := apiClient.GetInstance(ctx, model.ProjectId, *model.InstanceId)
	return req
}
