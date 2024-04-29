package describe

import (
	"context"
	"encoding/json"
	"fmt"

	"time"

	"github.com/inhies/go-bytesize"
	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/postgresflex/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
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

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("describe %s", backupIdArg),
		Short: "Shows details of a backup for a PostgreSQL Flex instance",
		Long:  "Shows details of a backup for a PostgreSQL Flex instance.",
		Example: examples.Build(
			examples.NewExample(
				`Get details of a backup with ID "xxx" for a PostgreSQL Flex instance with ID "yyy"`,
				"$ stackit postgresflex backup describe xxx --instance-id yyy"),
			examples.NewExample(
				`Get details of a backup with ID "xxx" for a PostgreSQL Flex instance with ID "yyy" in a table format`,
				"$ stackit postgresflex backup describe xxx --instance-id yyy --output-format pretty"),
		),
		Args: args.SingleArg(backupIdArg, nil),
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
				return fmt.Errorf("describe backup for PostgreSQL Flex instance: %w", err)
			}

			return outputResult(p, cmd, model.OutputFormat, *resp.Item)
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
	req := apiClient.GetBackup(ctx, model.ProjectId, model.InstanceId, model.BackupId)
	return req
}

func outputResult(p *print.Printer, cmd *cobra.Command, outputFormat string, backup postgresflex.Backup) error {
	backupStartTime, err := time.Parse(time.RFC3339, *backup.StartTime)
	if err != nil {
		return fmt.Errorf("parse backup start time: %w", err)
	}
	backupExpireDate := backupStartTime.AddDate(backupExpireYearOffset, backupExpireMonthOffset, backupExpireDayOffset).Format(time.DateOnly)

	switch outputFormat {
	case print.PrettyOutputFormat:
		table := tables.NewTable()
		table.AddRow("ID", *backup.Id)
		table.AddSeparator()
		table.AddRow("START TIME", *backup.StartTime)
		table.AddSeparator()
		table.AddRow("EXPIRES AT", backupExpireDate)
		table.AddSeparator()
		table.AddRow("BACKUP SIZE", bytesize.New(float64(*backup.Size)))

		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}

		return nil
	default:
		details, err := json.MarshalIndent(backup, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal backup for PostgreSQL Flex instance: %w", err)
		}
		cmd.Println(string(details))

		return nil
	}
}
