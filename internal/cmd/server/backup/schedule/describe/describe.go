package describe

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/goccy/go-yaml"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/serverbackup/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/serverbackup"
)

const (
	backupScheduleIdArg = "BACKUP_SCHEDULE_ID"
	serverIdFlag        = "server-id"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	ServerId         string
	BackupScheduleId string
}

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("describe %s", backupScheduleIdArg),
		Short: "Shows details of a Server Backup Schedule",
		Long:  "Shows details of a Server Backup Schedule.",
		Args:  args.SingleArg(backupScheduleIdArg, nil),
		Example: examples.Build(
			examples.NewExample(
				`Get details of a Server Backup Schedule with id "my-schedule-id"`,
				"$ stackit server backup schedule describe my-schedule-id"),
			examples.NewExample(
				`Get details of a Server Backup Schedule with id "my-schedule-id" in JSON format`,
				"$ stackit server backup schedule describe my-schedule-id --output-format json"),
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

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("read server backup schedule: %w", err)
			}

			return outputResult(p, model.OutputFormat, *resp)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().VarP(flags.UUIDFlag(), serverIdFlag, "s", "Server ID")

	err := flags.MarkFlagsRequired(cmd, serverIdFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	backupScheduleId := inputArgs[0]

	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel:  globalFlags,
		ServerId:         flags.FlagToStringValue(p, cmd, serverIdFlag),
		BackupScheduleId: backupScheduleId,
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

func buildRequest(ctx context.Context, model *inputModel, apiClient *serverbackup.APIClient) serverbackup.ApiGetBackupScheduleRequest {
	req := apiClient.GetBackupSchedule(ctx, model.ProjectId, model.ServerId, model.BackupScheduleId)
	return req
}

func outputResult(p *print.Printer, outputFormat string, schedule serverbackup.BackupSchedule) error {
	switch outputFormat {
	case print.JSONOutputFormat:
		details, err := json.MarshalIndent(schedule, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal server backup schedule: %w", err)
		}
		p.Outputln(string(details))

		return nil
	case print.YAMLOutputFormat:
		details, err := yaml.MarshalWithOptions(schedule, yaml.IndentSequence(true), yaml.UseJSONMarshaler())
		if err != nil {
			return fmt.Errorf("marshal server backup schedule: %w", err)
		}
		p.Outputln(string(details))

		return nil
	default:
		table := tables.NewTable()
		table.AddRow("SCHEDULE ID", utils.PtrString(schedule.Id))
		table.AddSeparator()
		table.AddRow("SCHEDULE NAME", utils.PtrString(schedule.Name))
		table.AddSeparator()
		table.AddRow("ENABLED", utils.PtrString(schedule.Enabled))
		table.AddSeparator()
		table.AddRow("RRULE", utils.PtrString(schedule.Rrule))
		table.AddSeparator()
		if schedule.BackupProperties != nil {
			table.AddRow("BACKUP NAME", utils.PtrString(schedule.BackupProperties.Name))
			table.AddSeparator()
			table.AddRow("BACKUP RETENTION DAYS", utils.PtrString(schedule.BackupProperties.RetentionPeriod))
			table.AddSeparator()
			ids := schedule.BackupProperties.VolumeIds
			table.AddRow("BACKUP VOLUME IDS", utils.JoinStringPtr(ids, "\n"))
		}
		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}

		return nil
	}
}
