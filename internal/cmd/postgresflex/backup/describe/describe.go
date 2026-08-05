package describe

import (
	"context"
	"fmt"
	"strconv"

	"github.com/stackitcloud/stackit-cli/internal/pkg/types"

	"github.com/spf13/cobra"
	postgresflex "github.com/stackitcloud/stackit-sdk-go/services/postgresflex/v3api"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/postgresflex/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
)

const (
	backupIdArg = "BACKUP_ID"

	instanceIdFlag = "instance-id"
)

type inputModel struct {
	*globalflags.GlobalFlagModel

	InstanceId string
	BackupId   int64
}

func NewCmd(params *types.CmdParams) *cobra.Command {
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

			return outputResult(params.Printer, model.OutputFormat, resp)
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
	backupIdStr := inputArgs[0]

	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	backupId, err := strconv.ParseInt(backupIdStr, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid backup id format, must be an integer: %w", err)
	}

	return &inputModel{
		GlobalFlagModel: globalFlags,
		InstanceId:      flags.FlagToStringValue(p, cmd, instanceIdFlag),
		BackupId:        backupId,
	}, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *postgresflex.APIClient) postgresflex.ApiGetBackupRequest {
	req := apiClient.DefaultAPI.GetBackup(ctx, model.ProjectId, model.Region, model.InstanceId, model.BackupId)
	return req
}

func outputResult(p *print.Printer, outputFormat string, backup *postgresflex.BackupData) error {
	return p.OutputResult(outputFormat, backup, func() error {
		if backup == nil {
			return fmt.Errorf("backup is nil")
		}

		table := tables.NewTable()
		table.AddRow("ID", backup.Id)
		table.AddSeparator()
		table.AddRow("COMPLETED AT", backup.CompletionTime)
		table.AddSeparator()
		table.AddRow("RETAINED UNTIL", backup.RetainedUntil)
		table.AddSeparator()
		table.AddRow("BACKUP SIZE", backup.Size)

		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}

		return nil
	})
}
