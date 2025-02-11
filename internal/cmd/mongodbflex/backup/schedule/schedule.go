package schedule

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/goccy/go-yaml"
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

func NewCmd(p *print.Printer) *cobra.Command {
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
		RunE: func(cmd *cobra.Command, _ []string) error {
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

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("read MongoDB Flex instance: %w", err)
			}

			return outputResult(p, model.OutputFormat, resp.Item)
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

func parseInput(p *print.Printer, cmd *cobra.Command) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		InstanceId:      *flags.FlagToStringPointer(p, cmd, instanceIdFlag),
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

func buildRequest(ctx context.Context, model *inputModel, apiClient *mongodbflex.APIClient) mongodbflex.ApiGetInstanceRequest {
	req := apiClient.GetInstance(ctx, model.ProjectId, model.InstanceId)
	return req
}

func outputResult(p *print.Printer, outputFormat string, instance *mongodbflex.Instance) error {
	output := struct {
		BackupSchedule                 string `json:"backup_schedule"`
		DailySnaphotRetentionDays      string `json:"daily_snapshot_retention_days"`
		MonthlySnapshotRetentionMonths string `json:"monthly_snapshot_retention_months"`
		PointInTimeWindowHours         string `json:"point_in_time_window_hours"`
		SnapshotRetentionDays          string `json:"snapshot_retention_days"`
		WeeklySnapshotRetentionWeeks   string `json:"weekly_snapshot_retention_weeks"`
	}{
		BackupSchedule:                 *instance.BackupSchedule,
		DailySnaphotRetentionDays:      (*instance.Options)["dailySnapshotRetentionDays"],
		MonthlySnapshotRetentionMonths: (*instance.Options)["monthlySnapshotRetentionDays"],
		PointInTimeWindowHours:         (*instance.Options)["pointInTimeWindowHours"],
		SnapshotRetentionDays:          (*instance.Options)["snapshotRetentionDays"],
		WeeklySnapshotRetentionWeeks:   (*instance.Options)["weeklySnapshotRetentionWeeks"],
	}

	switch outputFormat {
	case print.JSONOutputFormat:
		details, err := json.MarshalIndent(output, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal MongoDB Flex backup schedule: %w", err)
		}
		p.Outputln(string(details))

		return nil
	case print.YAMLOutputFormat:
		details, err := yaml.MarshalWithOptions(output, yaml.IndentSequence(true), yaml.UseJSONMarshaler())
		if err != nil {
			return fmt.Errorf("marshal MongoDB Flex backup schedule: %w", err)
		}
		p.Outputln(string(details))

		return nil
	default:
		table := tables.NewTable()
		table.AddRow("BACKUP SCHEDULE (UTC)", utils.PtrString(instance.BackupSchedule))
		table.AddSeparator()
		if instance.Options != nil {
			table.AddRow("DAILY SNAPSHOT RETENTION (DAYS)", (*instance.Options)["dailySnapshotRetentionDays"])
			table.AddSeparator()
			table.AddRow("MONTHLY SNAPSHOT RETENTION (MONTHS)", (*instance.Options)["monthlySnapshotRetentionMonths"])
			table.AddSeparator()
			table.AddRow("POINT IN TIME WINDOW (HOURS)", (*instance.Options)["pointInTimeWindowHours"])
			table.AddSeparator()
			table.AddRow("SNAPSHOT RETENTION (DAYS)", (*instance.Options)["snapshotRetentionDays"])
			table.AddSeparator()
			table.AddRow("WEEKLY SNAPSHOT RETENTION (WEEKS)", (*instance.Options)["weeklySnapshotRetentionWeeks"])
			table.AddSeparator()
		}
		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}

		return nil
	}
}
