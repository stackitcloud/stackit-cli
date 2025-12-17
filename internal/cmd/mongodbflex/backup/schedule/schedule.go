package schedule

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/types"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/mongodbflex/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/mongodbflex"
)

const (
	instanceIdFlag = "instance-id"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	InstanceId string
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "schedule",
		Short: "Shows details of the backup schedule and retention policy of a MongoDB Flex instance",
		Long:  "Shows details of the backup schedule and retention policy of a MongoDB Flex instance.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Get details of the backup schedule of a MongoDB Flex instance with ID "xxx"`,
				"$ stackit mongodbflex backup schedule --instance-id xxx"),
			examples.NewExample(
				`Get details of the backup schedule of a MongoDB Flex instance with ID "xxx" in JSON format`,
				"$ stackit mongodbflex backup schedule --instance-id xxx --output-format json"),
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

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("read MongoDB Flex instance: %w", err)
			}

			return outputResult(params.Printer, model.OutputFormat, resp.Item)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), instanceIdFlag, "Instance ID")

	err := flags.MarkFlagsRequired(cmd, instanceIdFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command, _ []string) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		InstanceId:      *flags.FlagToStringPointer(p, cmd, instanceIdFlag),
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *mongodbflex.APIClient) mongodbflex.ApiGetInstanceRequest {
	req := apiClient.GetInstance(ctx, model.ProjectId, model.InstanceId, model.Region)
	return req
}

func outputResult(p *print.Printer, outputFormat string, instance *mongodbflex.Instance) error {
	if instance == nil {
		return fmt.Errorf("instance is nil")
	}

	output := struct {
		BackupSchedule                 string `json:"backup_schedule"`
		DailySnaphotRetentionDays      string `json:"daily_snapshot_retention_days"`
		MonthlySnapshotRetentionMonths string `json:"monthly_snapshot_retention_months"`
		PointInTimeWindowHours         string `json:"point_in_time_window_hours"`
		SnapshotRetentionDays          string `json:"snapshot_retention_days"`
		WeeklySnapshotRetentionWeeks   string `json:"weekly_snapshot_retention_weeks"`
	}{
		BackupSchedule: utils.PtrString(instance.BackupSchedule),
	}
	if instance.Options != nil {
		output.DailySnaphotRetentionDays = (*instance.Options)["dailySnapshotRetentionDays"]
		output.MonthlySnapshotRetentionMonths = (*instance.Options)["monthlySnapshotRetentionDays"]
		output.PointInTimeWindowHours = (*instance.Options)["pointInTimeWindowHours"]
		output.SnapshotRetentionDays = (*instance.Options)["snapshotRetentionDays"]
		output.WeeklySnapshotRetentionWeeks = (*instance.Options)["weeklySnapshotRetentionWeeks"]
	}

	return p.OutputResult(outputFormat, output, func() error {
		table := tables.NewTable()
		table.AddRow("BACKUP SCHEDULE (UTC)", output.BackupSchedule)
		table.AddSeparator()
		table.AddRow("DAILY SNAPSHOT RETENTION (DAYS)", output.DailySnaphotRetentionDays)
		table.AddSeparator()
		table.AddRow("MONTHLY SNAPSHOT RETENTION (MONTHS)", output.MonthlySnapshotRetentionMonths)
		table.AddSeparator()
		table.AddRow("POINT IN TIME WINDOW (HOURS)", output.PointInTimeWindowHours)
		table.AddSeparator()
		table.AddRow("SNAPSHOT RETENTION (DAYS)", output.SnapshotRetentionDays)
		table.AddSeparator()
		table.AddRow("WEEKLY SNAPSHOT RETENTION (WEEKS)", output.WeeklySnapshotRetentionWeeks)
		table.AddSeparator()

		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}

		return nil
	})
}
