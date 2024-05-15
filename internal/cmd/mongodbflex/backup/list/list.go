package list

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/goccy/go-yaml"
	"github.com/inhies/go-bytesize"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/mongodbflex/client"
	mongodbflexUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/mongodbflex/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/mongodbflex"
)

const (
	instanceIdFlag = "instance-id"
	limitFlag      = "limit"
)

type inputModel struct {
	*globalflags.GlobalFlagModel

	InstanceId *string
	Limit      *int64
}

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Lists all backups which are available for a MongoDB Flex instance",
		Long:  "Lists all backups which are available for a MongoDB Flex instance.",
		Example: examples.Build(
			examples.NewExample(
				`List all backups of instance with ID "xxx"`,
				"$ stackit mongodbflex backup list --instance-id xxx"),
			examples.NewExample(
				`List all backups of instance with ID "xxx" in JSON format`,
				"$ stackit mongodbflex backup list --instance-id xxx --output-format json"),
			examples.NewExample(
				`List up to 10 backups of instance with ID "xxx"`,
				"$ stackit mongodbflex backup list --instance-id xxx --limit 10"),
		),
		Args: args.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
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

			instanceLabel, err := mongodbflexUtils.GetInstanceName(ctx, apiClient, model.ProjectId, *model.InstanceId)
			if err != nil {
				p.Debug(print.ErrorLevel, "get instance name: %v", err)
				instanceLabel = *model.InstanceId
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("get backups for MongoDB Flex instance %q: %w", instanceLabel, err)
			}
			if resp.Items == nil || len(*resp.Items) == 0 {
				cmd.Printf("No backups found for instance %q\n", instanceLabel)
				return nil
			}
			backups := *resp.Items

			restoreJobs, err := apiClient.ListRestoreJobs(ctx, model.ProjectId, *model.InstanceId).Execute()
			if err != nil {
				return fmt.Errorf("get restore jobs for MongoDB Flex instance %q: %w", instanceLabel, err)
			}

			// Truncate output
			if model.Limit != nil && len(backups) > int(*model.Limit) {
				backups = backups[:*model.Limit]
			}

			return outputResult(p, model.OutputFormat, backups, restoreJobs)
		},
	}

	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), instanceIdFlag, "Instance ID")
	cmd.Flags().Int64(limitFlag, 0, "Maximum number of entries to list")

	err := flags.MarkFlagsRequired(cmd, instanceIdFlag)
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
		InstanceId:      flags.FlagToStringPointer(p, cmd, instanceIdFlag),
		Limit:           flags.FlagToInt64Pointer(p, cmd, limitFlag),
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

func buildRequest(ctx context.Context, model *inputModel, apiClient *mongodbflex.APIClient) mongodbflex.ApiListBackupsRequest {
	req := apiClient.ListBackups(ctx, model.ProjectId, *model.InstanceId)
	return req
}

func outputResult(p *print.Printer, outputFormat string, backups []mongodbflex.Backup, restoreJobs *mongodbflex.ListRestoreJobsResponse) error {
	switch outputFormat {
	case print.JSONOutputFormat:
		details, err := json.MarshalIndent(backups, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal MongoDB Flex backups list: %w", err)
		}
		p.Outputln(string(details))

		return nil
	case print.YAMLOutputFormat:
		details, err := yaml.Marshal(backups)
		if err != nil {
			return fmt.Errorf("marshal MongoDB Flex backups list: %w", err)
		}
		p.Outputln(string(details))

		return nil
	default:
		table := tables.NewTable()
		table.SetHeader("ID", "CREATED AT", "EXPIRES AT", "BACKUP SIZE", "RESTORE STATUS")
		for i := range backups {
			backup := backups[i]
			restoreStatus := mongodbflexUtils.GetRestoreStatus(*backup.Id, restoreJobs)
			table.AddRow(*backup.Id, *backup.StartTime, *backup.EndTime, bytesize.New(float64(*backup.Size)), restoreStatus)
		}
		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}

		return nil
	}
}
