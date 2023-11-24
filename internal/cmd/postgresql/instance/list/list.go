package list

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/postgresql/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/postgresql"
)

const (
	limitFlag = "limit"
)

type flagModel struct {
	*globalflags.GlobalFlagModel
	Limit *int64
}

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List all PostgreSQL instances",
		Long:    "List all PostgreSQL instances",
		Example: `$ stackit postgresql instance list --project-id xxx`,
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
				return fmt.Errorf("get PostgreSQL instances: %w", err)
			}
			instances := *resp.Instances
			if len(instances) == 0 {
				cmd.Printf("No instances found for product with ID %s\n", model.ProjectId)
				return nil
			}

			// Truncate output
			if model.Limit != nil && len(instances) > int(*model.Limit) {
				instances = instances[:*model.Limit]
			}

			// Show output as table
			table := tables.NewTable()
			table.SetHeader("ID", "NAME", "LAST_OPERATION.TYPE", "LAST_OPERATION.STATE")
			for i := range instances {
				instance := instances[i]
				table.AddRow(*instance.InstanceId, *instance.Name, *instance.LastOperation.Type, *instance.LastOperation.State)
			}
			err = table.Render(cmd)
			if err != nil {
				return fmt.Errorf("render table: %w", err)
			}

			return nil
		},
	}

	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Int64(limitFlag, 0, "Maximum number of entries to list")
}

func parseFlags(cmd *cobra.Command) (*flagModel, error) {
	globalFlags := globalflags.Parse()
	if globalFlags.ProjectId == "" {
		return nil, fmt.Errorf("project ID not set")
	}

	limit := utils.FlagToInt64Pointer(cmd, limitFlag)
	if limit != nil && *limit < 1 {
		return nil, fmt.Errorf("limit must be greater than 0")
	}

	return &flagModel{
		GlobalFlagModel: globalFlags,
		Limit:           utils.FlagToInt64Pointer(cmd, limitFlag),
	}, nil
}

func buildRequest(ctx context.Context, model *flagModel, apiClient *postgresql.APIClient) postgresql.ApiGetInstancesRequest {
	req := apiClient.GetInstances(ctx, model.ProjectId)
	return req
}
