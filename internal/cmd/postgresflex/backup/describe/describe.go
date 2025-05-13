package describe

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/goccy/go-yaml"
	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/postgresflex/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/postgresflex"
)

const (
	backupIdArg = "BACKUP_ID"

	instanceIdFlag = "instance-id"

	backupExpireYearOffset  = 0
	backupExpireMonthOffset = 0
	backupExpireDayOffset   = 30
)

type inputModel struct {
	*globalflags.GlobalFlagModel

	InstanceId string
	BackupId   string
}

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("describe %s", backupIdArg),
		Short: "Shows details of a backup for a PostgreSQL Flex instance",
		Long:  "Shows details of a backup for a PostgreSQL Flex instance.",
		Example: examples.Build(
			examples.NewExample(
				`Get details of a backup with ID "xxx" for a PostgreSQL Flex instance with ID "yyy"`,
				"$ stackit postgresflex backup describe xxx --instance-id yyy"),
			examples.NewExample(
				`Get details of a backup with ID "xxx" for a PostgreSQL Flex instance with ID "yyy" in JSON format`,
				"$ stackit postgresflex backup describe xxx --instance-id yyy --output-format json"),
		),
		Args: args.SingleArg(backupIdArg, nil),
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
				return fmt.Errorf("describe backup for PostgreSQL Flex instance: %w", err)
			}

			return outputResult(params.Printer, model.OutputFormat, *resp.Item)
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

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	backupId := inputArgs[0]

	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	return &inputModel{
		GlobalFlagModel: globalFlags,
		InstanceId:      flags.FlagToStringValue(p, cmd, instanceIdFlag),
		BackupId:        backupId,
	}, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *postgresflex.APIClient) postgresflex.ApiGetBackupRequest {
	req := apiClient.GetBackup(ctx, model.ProjectId, model.Region, model.InstanceId, model.BackupId)
	return req
}

func outputResult(p *print.Printer, outputFormat string, backup postgresflex.Backup) error {
	if backup.StartTime == nil || *backup.StartTime == "" {
		return fmt.Errorf("start time not defined")
	}
	backupStartTime, err := time.Parse(time.RFC3339, utils.PtrString(backup.StartTime))
	if err != nil {
		return fmt.Errorf("parse backup start time: %w", err)
	}
	backupExpireDate := backupStartTime.AddDate(backupExpireYearOffset, backupExpireMonthOffset, backupExpireDayOffset).Format(time.DateOnly)

	switch outputFormat {
	case print.JSONOutputFormat:
		details, err := json.MarshalIndent(backup, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal backup for PostgreSQL Flex backup: %w", err)
		}
		p.Outputln(string(details))

		return nil
	case print.YAMLOutputFormat:
		details, err := yaml.MarshalWithOptions(backup, yaml.IndentSequence(true), yaml.UseJSONMarshaler())
		if err != nil {
			return fmt.Errorf("marshal backup for PostgreSQL Flex backup: %w", err)
		}
		p.Outputln(string(details))

		return nil
	default:
		table := tables.NewTable()
		table.AddRow("ID", utils.PtrString(backup.Id))
		table.AddSeparator()
		table.AddRow("CREATED AT", utils.PtrString(backup.StartTime))
		table.AddSeparator()
		table.AddRow("EXPIRES AT", backupExpireDate)
		table.AddSeparator()

		backupSize := utils.PtrByteSizeDefault(backup.Size, "n/a")
		table.AddRow("BACKUP SIZE", backupSize)

		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}

		return nil
	}
}
