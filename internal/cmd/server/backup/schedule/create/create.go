package create

import (
	"context"
	"encoding/json"
	"fmt"

	iaasClient "github.com/stackitcloud/stackit-cli/internal/pkg/services/iaas/client"

	"github.com/goccy/go-yaml"
	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	cliErr "github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	iaasUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/iaas/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/serverbackup/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/serverbackup"
)

const (
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
	BackupScheduleName    string
	Enabled               bool
	Rrule                 string
	BackupName            string
	BackupRetentionPeriod int64
	BackupVolumeIds       []string
}

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Creates a Server Backup Schedule",
		Long:  "Creates a Server Backup Schedule.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Create a Server Backup Schedule with name "myschedule" and backup name "mybackup"`,
				`$ stackit server backup schedule create --server-id xxx --backup-name=mybackup --backup-schedule-name=myschedule`),
			examples.NewExample(
				`Create a Server Backup Schedule with name "myschedule", backup name "mybackup" and retention period of 5 days`,
				`$ stackit server backup schedule create --server-id xxx --backup-name=mybackup --backup-schedule-name=myschedule --backup-retention-period=5`),
		),
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := context.Background()

			model, err := parseInput(params.Printer, cmd)
			if err != nil {
				return err
			}

			// Configure API client
			apiClient, err := client.ConfigureClient(params.Printer, params.CliVersion)
			if err != nil {
				return err
			}

			serverLabel := model.ServerId
			// Get server name
			if iaasApiClient, err := iaasClient.ConfigureClient(params.Printer, params.CliVersion); err == nil {
				serverName, err := iaasUtils.GetServerName(ctx, iaasApiClient, model.ProjectId, model.ServerId)
				if err != nil {
					params.Printer.Debug(print.ErrorLevel, "get server name: %v", err)
				} else if serverName != "" {
					serverLabel = serverName
				}
			}

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to create a Backup Schedule for server %s?", serverLabel)
				err = params.Printer.PromptForConfirmation(prompt)
				if err != nil {
					return err
				}
			}

			// Call API
			req, err := buildRequest(ctx, model, apiClient)
			if err != nil {
				return err
			}
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("create Server Backup Schedule: %w", err)
			}

			return outputResult(params.Printer, model.OutputFormat, serverLabel, *resp)
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

	err := flags.MarkFlagsRequired(cmd, serverIdFlag, backupScheduleNameFlag, backupNameFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &cliErr.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel:       globalFlags,
		ServerId:              flags.FlagToStringValue(p, cmd, serverIdFlag),
		BackupRetentionPeriod: flags.FlagWithDefaultToInt64Value(p, cmd, backupRetentionPeriodFlag),
		BackupScheduleName:    flags.FlagToStringValue(p, cmd, backupScheduleNameFlag),
		BackupName:            flags.FlagToStringValue(p, cmd, backupNameFlag),
		Rrule:                 flags.FlagWithDefaultToStringValue(p, cmd, rruleFlag),
		Enabled:               flags.FlagToBoolValue(p, cmd, enabledFlag),
		BackupVolumeIds:       flags.FlagToStringSliceValue(p, cmd, backupVolumeIdsFlag),
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

func buildRequest(ctx context.Context, model *inputModel, apiClient *serverbackup.APIClient) (serverbackup.ApiCreateBackupScheduleRequest, error) {
	req := apiClient.CreateBackupSchedule(ctx, model.ProjectId, model.ServerId, model.Region)
	backupProperties := serverbackup.BackupProperties{
		Name:            &model.BackupName,
		RetentionPeriod: &model.BackupRetentionPeriod,
		VolumeIds:       &model.BackupVolumeIds,
	}
	if model.BackupVolumeIds == nil {
		backupProperties.VolumeIds = nil
	}
	req = req.CreateBackupSchedulePayload(serverbackup.CreateBackupSchedulePayload{
		Enabled:          &model.Enabled,
		Name:             &model.BackupScheduleName,
		Rrule:            &model.Rrule,
		BackupProperties: &backupProperties,
	})
	return req, nil
}

func outputResult(p *print.Printer, outputFormat, serverLabel string, resp serverbackup.BackupSchedule) error {
	switch outputFormat {
	case print.JSONOutputFormat:
		details, err := json.MarshalIndent(resp, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal server backup schedule: %w", err)
		}
		p.Outputln(string(details))

		return nil
	case print.YAMLOutputFormat:
		details, err := yaml.MarshalWithOptions(resp, yaml.IndentSequence(true), yaml.UseJSONMarshaler())
		if err != nil {
			return fmt.Errorf("marshal server backup schedule: %w", err)
		}
		p.Outputln(string(details))

		return nil
	default:
		p.Outputf("Created server backup schedule for server %s. Backup Schedule ID: %s\n", serverLabel, utils.PtrString(resp.Id))
		return nil
	}
}
