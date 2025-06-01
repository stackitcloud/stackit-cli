package describe

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/goccy/go-yaml"
	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/iaas/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/stackitcloud/stackit-sdk-go/services/iaas"
)

const (
	backupIdArg = "BACKUP_ID"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	BackupId string
}

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("describe %s", backupIdArg),
		Short: "Describes a backup",
		Long:  "Describes a backup by its ID.",
		Args:  args.SingleArg(backupIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Get details of a backup`,
				"$ stackit volume backup describe xxx-xxx-xxx"),
			examples.NewExample(
				`Get details of a backup in JSON format`,
				"$ stackit volume backup describe xxx-xxx-xxx --output-format json"),
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
			backup, err := req.Execute()
			if err != nil {
				return fmt.Errorf("get backup details: %w", err)
			}

			return outputResult(params.Printer, model.OutputFormat, backup)
		},
	}
	return cmd
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	backupId := inputArgs[0]

	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
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

func buildRequest(ctx context.Context, model *inputModel, apiClient *iaas.APIClient) iaas.ApiGetBackupRequest {
	req := apiClient.GetBackup(ctx, model.ProjectId, model.BackupId)
	return req
}

func outputResult(p *print.Printer, outputFormat string, backup *iaas.Backup) error {
	if backup == nil {
		return fmt.Errorf("backup response is empty")
	}

	switch outputFormat {
	case print.JSONOutputFormat:
		details, err := json.MarshalIndent(backup, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal backup: %w", err)
		}
		p.Outputln(string(details))
		return nil

	case print.YAMLOutputFormat:
		details, err := yaml.MarshalWithOptions(backup, yaml.IndentSequence(true), yaml.UseJSONMarshaler())
		if err != nil {
			return fmt.Errorf("marshal backup: %w", err)
		}
		p.Outputln(string(details))
		return nil

	default:
		table := tables.NewTable()
		table.AddRow(
			"ID",
			utils.PtrString(backup.Id),
		)
		table.AddRow(
			"NAME",
			utils.PtrString(backup.Name),
		)
		table.AddRow(
			"SIZE",
			utils.PtrByteSizeDefault((*int64)(backup.Size), ""),
		)
		table.AddRow(
			"STATUS",
			utils.PtrString(backup.Status),
		)
		table.AddRow(
			"SNAPSHOT ID",
			utils.PtrString(backup.SnapshotId),
		)
		table.AddRow(
			"VOLUME ID",
			utils.PtrString(backup.VolumeId),
		)
		table.AddRow(
			"AVAILABILITY ZONE",
			utils.PtrString(backup.AvailabilityZone),
		)
		table.AddRow(
			"LABELS",
			utils.PtrStringDefault(backup.Labels, ""),
		)
		table.AddRow(
			"CREATED AT",
			utils.ConvertTimePToDateTimeString(backup.CreatedAt),
		)
		table.AddRow(
			"UPDATED AT",
			utils.ConvertTimePToDateTimeString(backup.UpdatedAt),
		)

		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}

		return nil
	}
}
