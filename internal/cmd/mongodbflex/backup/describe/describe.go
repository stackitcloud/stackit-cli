package describe

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
	mongoUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/mongodbflex/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/mongodbflex"
)

const (
	backupIdArg = "BACKUP_ID"

	instanceIdFlag = "instance-id"
)

type inputModel struct {
	*globalflags.GlobalFlagModel

	InstanceId string
	BackupId   string
}

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("describe %s", backupIdArg),
		Short: "Shows details of a backup for a MongoDB Flex instance",
		Long:  "Shows details of a backup for a MongoDB Flex instance.",
		Example: examples.Build(
			examples.NewExample(
				`Get details of a backup with ID "xxx" for a MongoDB Flex instance with ID "yyy"`,
				"$ stackit mongodbflex backup describe xxx --instance-id yyy"),
			examples.NewExample(
				`Get details of a backup with ID "xxx" for a MongoDB Flex instance with ID "yyy" in JSON format`,
				"$ stackit mongodbflex backup describe xxx --instance-id yyy --output-format json"),
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

			instanceLabel, err := mongoUtils.GetInstanceName(ctx, apiClient, model.ProjectId, model.InstanceId)
			if err != nil {
				p.Debug(print.ErrorLevel, "get instance name: %v", err)
				instanceLabel = model.InstanceId
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()

			if err != nil {
				return fmt.Errorf("describe backup for MongoDB Flex instance: %w", err)
			}

			restoreJobs, err := apiClient.ListRestoreJobs(ctx, model.ProjectId, model.InstanceId).Execute()
			if err != nil {
				return fmt.Errorf("get restore jobs for MongoDB Flex instance %q: %w", instanceLabel, err)
			}

			restoreJobState := mongoUtils.GetRestoreStatus(model.BackupId, restoreJobs)
			return outputResult(p, model.OutputFormat, restoreJobState, *resp.Item)
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

	model := inputModel{
		GlobalFlagModel: globalFlags,
		InstanceId:      flags.FlagToStringValue(p, cmd, instanceIdFlag),
		BackupId:        backupId,
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

func buildRequest(ctx context.Context, model *inputModel, apiClient *mongodbflex.APIClient) mongodbflex.ApiGetBackupRequest {
	req := apiClient.GetBackup(ctx, model.ProjectId, model.InstanceId, model.BackupId)
	return req
}

func outputResult(p *print.Printer, outputFormat, restoreStatus string, backup mongodbflex.Backup) error {
	switch outputFormat {
	case print.JSONOutputFormat:
		details, err := json.MarshalIndent(backup, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal MongoDB Flex backup: %w", err)
		}
		p.Outputln(string(details))

		return nil
	case print.YAMLOutputFormat:
		details, err := yaml.MarshalWithOptions(backup, yaml.IndentSequence(true), yaml.UseJSONMarshaler())
		if err != nil {
			return fmt.Errorf("marshal MongoDB Flex backup: %w", err)
		}
		p.Outputln(string(details))

		return nil
	default:
		table := tables.NewTable()
		table.AddRow("ID", utils.PtrString(backup.Id))
		table.AddSeparator()
		table.AddRow("CREATED AT", utils.PtrString(backup.StartTime))
		table.AddSeparator()
		table.AddRow("EXPIRES AT", utils.PtrString(backup.EndTime))
		table.AddSeparator()
		backupSize := utils.PtrByteSizeDefault(backup.Size, "n/a")
		table.AddRow("BACKUP SIZE", backupSize)
		table.AddSeparator()
		table.AddRow("RESTORE STATUS", restoreStatus)

		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}

		return nil
	}
}
