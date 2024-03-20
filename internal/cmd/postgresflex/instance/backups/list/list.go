package list

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/depp/bytesize"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/postgresflex/client"
	postgresflexUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/postgresflex/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/postgresflex"
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

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Lists all backups which are available for a specific PostgreSQL Flex instance",
		Long:  "Lists all backups which are available for a specific PostgreSQL Flex instance.",
		Example: examples.Build(
			examples.NewExample(
				`List all backups of instance with ID "xxx"`,
				"$ stackit postgresflex instance backups list --instance-id xxx"),
			examples.NewExample(
				`List all backups of instance with ID "xxx" in JSON format`,
				"$ stackit postgresflex instance backups list --instance-id xxx --output-format json"),
			examples.NewExample(
				`List up to 10 backups of instance with ID "xxx"`,
				"$ stackit postgresflex instance backups list --instance-id xxx --limit 10"),
		),
		Args: args.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			model, err := parseInput(cmd)
			if err != nil {
				return err
			}

			// Configure API client
			apiClient, err := client.ConfigureClient(cmd)
			if err != nil {
				return err
			}

			instanceLabel, err := postgresflexUtils.GetInstanceName(ctx, apiClient, model.ProjectId, *model.InstanceId)
			if err != nil {
				instanceLabel = *model.InstanceId
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("get backups for PostgreSQL Flex instance %q: %w\n", instanceLabel, err)
			}
			if resp.Items == nil || len(*resp.Items) == 0 {
				cmd.Printf("No backups found for instance %q\n", instanceLabel)
				return nil
			}
			backups := *resp.Items

			// Truncate output
			if model.Limit != nil && len(backups) > int(*model.Limit) {
				backups = backups[:*model.Limit]
			}

			return outputResult(cmd, model.OutputFormat, backups)
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

func parseInput(cmd *cobra.Command) (*inputModel, error) {
	globalFlags := globalflags.Parse(cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	limit := flags.FlagToInt64Pointer(cmd, limitFlag)
	if limit != nil && *limit < 1 {
		return nil, &errors.FlagValidationError{
			Flag:    limitFlag,
			Details: "must be greater than 0",
		}
	}

	return &inputModel{
		GlobalFlagModel: globalFlags,
		InstanceId:      flags.FlagToStringPointer(cmd, instanceIdFlag),
		Limit:           flags.FlagToInt64Pointer(cmd, limitFlag),
	}, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *postgresflex.APIClient) postgresflex.ApiListBackupsRequest {
	req := apiClient.ListBackups(ctx, model.ProjectId, *model.InstanceId)
	return req
}

func outputResult(cmd *cobra.Command, outputFormat string, backups []postgresflex.Backup) error {
	switch outputFormat {
	case globalflags.JSONOutputFormat:
		details, err := json.MarshalIndent(backups, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal PostgreSQL Flex instance list: %w", err)
		}
		cmd.Println(string(details))

		return nil
	default:
		table := tables.NewTable()
		table.SetHeader("ID", "NAME", "START TIME", "END TIME", "BACKUP SIZE")
		for i := range backups {
			backup := backups[i]
			table.AddRow(*backup.Id, *backup.Name, *backup.StartTime, *backup.EndTime, bytesize.Format(uint64(*backup.Size)))
		}
		err := table.Display(cmd)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}

		return nil
	}
}