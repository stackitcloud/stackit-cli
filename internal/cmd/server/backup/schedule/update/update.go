package update

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/goccy/go-yaml"
	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	cliErr "github.com/stackitcloud/stackit-cli/internal/pkg/errors"
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
	scheduleIdArg = "SCHEDULE_ID"

	backupScheduleNameFlag    = "backup-schedule-name"
	enabledFlag               = "enabled"
	rruleFlag                 = "rrule"
	backupNameFlag            = "backup-name"
	backupVolumeIdsFlag       = "backup-volume-ids"
	backupRetentionPeriodFlag = "backup-retention-period"
	serverIdFlag              = "server-id"

	defaultRrule           = "DTSTART;TZID=Europe/Sofia:20200803T023000 RRULE:FREQ=DAILY;INTERVAL=1"
	defaultRetentionPeriod = 14
	defaultEnabled         = true
)

type inputModel struct {
	*globalflags.GlobalFlagModel

	ServerId              string
	BackupScheduleId      string
	BackupScheduleName    *string
	Enabled               *bool
	Rrule                 *string
	BackupName            *string
	BackupRetentionPeriod *int64
	BackupVolumeIds       []string
}

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("update %s", scheduleIdArg),
		Short: "Updates a Server Backup Schedule",
		Long:  "Updates a Server Backup Schedule.",
		Example: examples.Build(
			examples.NewExample(
				`Update the retention period of the backup schedule "zzz" of server "xxx"`,
				"$ stackit server backup schedule update zzz --server-id=xxx --backup-retention-period=20"),
			examples.NewExample(
				`Update the backup name of the backup schedule "zzz" of server "xxx"`,
				"$ stackit server backup schedule update zzz --server-id=xxx --backup-name=newname"),
		),
		Args: args.SingleArg(scheduleIdArg, nil),
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

			currentBackupSchedule, err := apiClient.GetBackupScheduleExecute(ctx, model.ProjectId, model.ServerId, model.Region, model.BackupScheduleId)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get current server backup schedule: %v", err)
				return err
			}

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to update Server Backup Schedule %q?", model.BackupScheduleId)
				err = params.Printer.PromptForConfirmation(prompt)
				if err != nil {
					return err
				}
			}

			// Call API
			req, err := buildRequest(ctx, model, apiClient, *currentBackupSchedule)
			if err != nil {
				return err
			}
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("update Server Backup Schedule: %w", err)
			}

			return outputResult(params.Printer, model.OutputFormat, *resp)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().VarP(flags.UUIDFlag(), serverIdFlag, "s", "Server ID")

	cmd.Flags().StringP(backupScheduleNameFlag, "n", "", "Backup schedule name")
	cmd.Flags().StringP(backupNameFlag, "b", "", "Backup name")
	cmd.Flags().Int64P(backupRetentionPeriodFlag, "d", defaultRetentionPeriod, "Backup retention period (in days)")
	cmd.Flags().BoolP(enabledFlag, "e", defaultEnabled, "Is the server backup schedule enabled")
	cmd.Flags().StringP(rruleFlag, "r", defaultRrule, "Backup RRULE (recurrence rule)")
	cmd.Flags().VarP(flags.UUIDSliceFlag(), backupVolumeIdsFlag, "i", "Backup volume IDs, as comma separated UUID values.")

	err := flags.MarkFlagsRequired(cmd, serverIdFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	scheduleId := inputArgs[0]

	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &cliErr.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel:       globalFlags,
		BackupScheduleId:      scheduleId,
		ServerId:              flags.FlagToStringValue(p, cmd, serverIdFlag),
		BackupRetentionPeriod: flags.FlagToInt64Pointer(p, cmd, backupRetentionPeriodFlag),
		BackupScheduleName:    flags.FlagToStringPointer(p, cmd, backupScheduleNameFlag),
		BackupName:            flags.FlagToStringPointer(p, cmd, backupNameFlag),
		Rrule:                 flags.FlagToStringPointer(p, cmd, rruleFlag),
		Enabled:               flags.FlagToBoolPointer(p, cmd, enabledFlag),
		BackupVolumeIds:       flags.FlagToStringSliceValue(p, cmd, backupVolumeIdsFlag),
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *serverbackup.APIClient, old serverbackup.BackupSchedule) (serverbackup.ApiUpdateBackupScheduleRequest, error) {
	req := apiClient.UpdateBackupSchedule(ctx, model.ProjectId, model.ServerId, model.Region, model.BackupScheduleId)

	if model.BackupName != nil {
		old.BackupProperties.Name = model.BackupName
	}
	if model.BackupRetentionPeriod != nil {
		old.BackupProperties.RetentionPeriod = model.BackupRetentionPeriod
	}
	if model.BackupVolumeIds != nil {
		old.BackupProperties.VolumeIds = &model.BackupVolumeIds
	}
	if model.Enabled != nil {
		old.Enabled = model.Enabled
	}
	if model.BackupScheduleName != nil {
		old.Name = model.BackupScheduleName
	}
	if model.Rrule != nil {
		old.Rrule = model.Rrule
	}

	req = req.UpdateBackupSchedulePayload(serverbackup.UpdateBackupSchedulePayload{
		Enabled:          old.Enabled,
		Name:             old.Name,
		Rrule:            old.Rrule,
		BackupProperties: old.BackupProperties,
	})
	return req, nil
}

func outputResult(p *print.Printer, outputFormat string, resp serverbackup.BackupSchedule) error {
	switch outputFormat {
	case print.JSONOutputFormat:
		details, err := json.MarshalIndent(resp, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal update server backup schedule: %w", err)
		}
		p.Outputln(string(details))

		return nil
	case print.YAMLOutputFormat:
		details, err := yaml.MarshalWithOptions(resp, yaml.IndentSequence(true), yaml.UseJSONMarshaler())
		if err != nil {
			return fmt.Errorf("marshal update server backup schedule: %w", err)
		}
		p.Outputln(string(details))

		return nil
	default:
		p.Info("Updated server backup schedule %s\n", utils.PtrString(resp.Id))
		return nil
	}
}
