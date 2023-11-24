package describe

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/postgresql/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/postgresql"
)

const (
	instanceIdFlag = "instance-id"
)

type flagModel struct {
	*globalflags.GlobalFlagModel
	InstanceId string
}

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "describe",
		Short:   "Get details of a PostgreSQL instance",
		Long:    "Get details of a PostgreSQL instance",
		Example: `$ stackit postgresql instance describe --project-id xxx --instance-id xxx`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			model, err := parseFlags(cmd)
			if err != nil {
				return err
			}
			// Configure API client
			apiClient, err := client.ConfigureClient(cmd)
			if err != nil {
				return fmt.Errorf("authentication failed, please run \"stackit auth login\" or \"stackit auth activate-service-account\"")
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("read PostgreSQL instance: %w", err)
			}

			return outputResult(cmd, model.OutputFormat, resp)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), instanceIdFlag, "Instance ID")

	err := utils.MarkFlagsRequired(cmd, instanceIdFlag)
	cobra.CheckErr(err)
}

func parseFlags(cmd *cobra.Command) (*flagModel, error) {
	globalFlags := globalflags.Parse(cmd)
	if globalFlags.ProjectId == "" {
		return nil, fmt.Errorf("project ID not set")
	}

	return &flagModel{
		GlobalFlagModel: globalFlags,
		InstanceId:      utils.FlagToStringValue(cmd, instanceIdFlag),
	}, nil
}

func buildRequest(ctx context.Context, model *flagModel, apiClient *postgresql.APIClient) postgresql.ApiGetInstanceRequest {
	req := apiClient.GetInstance(ctx, model.ProjectId, model.InstanceId)
	return req
}

func outputResult(cmd *cobra.Command, outputFormat string, instance *postgresql.Instance) error {
	switch outputFormat {
	case globalflags.TableOutputFormat:
		table := tables.NewTable()
		table.SetHeader("ID", "NAME", "LAST_OPERATION.TYPE", "LAST_OPERATION.STATE")
		table.AddRow(*instance.InstanceId, *instance.Name, *instance.LastOperation.Type, *instance.LastOperation.State)
		err := table.Render(cmd)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}

		return nil
	default:
		details, err := json.MarshalIndent(instance, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal PostgreSQL instance: %w", err)
		}
		cmd.Println(string(details))

		return nil
	}
}
