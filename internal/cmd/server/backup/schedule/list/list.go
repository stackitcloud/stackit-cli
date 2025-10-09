package list

import (
	"context"
	"encoding/json"
	"fmt"

	iaasClient "github.com/stackitcloud/stackit-cli/internal/pkg/services/iaas/client"

	"github.com/goccy/go-yaml"
	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	iaasUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/iaas/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/serverbackup/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/serverbackup"
)

const (
	limitFlag    = "limit"
	serverIdFlag = "server-id"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	ServerId string
	Limit    *int64
}

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Lists all server backup schedules",
		Long:  "Lists all server backup schedules.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`List all backup schedules for a server with ID "xxx"`,
				"$ stackit server backup schedule list --server-id xxx"),
			examples.NewExample(
				`List all backup schedules for a server with ID "xxx" in JSON format`,
				"$ stackit server backup schedule list --server-id xxx --output-format json"),
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

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("list server backup schedules: %w", err)
			}
			schedules := *resp.Items
			if len(schedules) == 0 {
				params.Printer.Info("No backup schedules found for server %s\n", serverLabel)
				return nil
			}

			// Truncate output
			if model.Limit != nil && len(schedules) > int(*model.Limit) {
				schedules = schedules[:*model.Limit]
			}
			return outputResult(params.Printer, model.OutputFormat, schedules)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Int64(limitFlag, 0, "Maximum number of entries to list")
	cmd.Flags().VarP(flags.UUIDFlag(), serverIdFlag, "s", "Server ID")

	err := flags.MarkFlagsRequired(cmd, serverIdFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	limit := flags.FlagToInt64Pointer(p, cmd, limitFlag)
	if limit != nil && *limit < 1 {
		return nil, &errors.FlagValidationError{
			Flag:    limitFlag,
			Details: "must be greater than 0",
		}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		ServerId:        flags.FlagToStringValue(p, cmd, serverIdFlag),
		Limit:           limit,
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *serverbackup.APIClient) serverbackup.ApiListBackupSchedulesRequest {
	req := apiClient.ListBackupSchedules(ctx, model.ProjectId, model.ServerId, model.Region)
	return req
}

func outputResult(p *print.Printer, outputFormat string, schedules []serverbackup.BackupSchedule) error {
	switch outputFormat {
	case print.JSONOutputFormat:
		details, err := json.MarshalIndent(schedules, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal Server Backup Schedules list: %w", err)
		}
		p.Outputln(string(details))

		return nil
	case print.YAMLOutputFormat:
		details, err := yaml.MarshalWithOptions(schedules, yaml.IndentSequence(true), yaml.UseJSONMarshaler())
		if err != nil {
			return fmt.Errorf("marshal Server Backup Schedules list: %w", err)
		}
		p.Outputln(string(details))

		return nil
	default:
		table := tables.NewTable()
		table.SetHeader("SCHEDULE ID", "SCHEDULE NAME", "ENABLED", "RRULE", "BACKUP NAME", "BACKUP RETENTION DAYS", "BACKUP VOLUME IDS")
		for i := range schedules {
			s := schedules[i]

			backupName := ""
			retentionPeriod := ""
			ids := ""
			if s.BackupProperties != nil {
				backupName = utils.PtrString(s.BackupProperties.Name)
				retentionPeriod = utils.PtrString(s.BackupProperties.RetentionPeriod)

				ids = utils.JoinStringPtr(s.BackupProperties.VolumeIds, ",")
			}
			table.AddRow(
				utils.PtrString(s.Id),
				utils.PtrString(s.Name),
				utils.PtrString(s.Enabled),
				utils.PtrString(s.Rrule),
				backupName,
				retentionPeriod,
				ids,
			)
		}
		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}
		return nil
	}
}
