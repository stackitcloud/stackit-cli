package list

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/config"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/postgresql/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stackitcloud/stackit-sdk-go/services/postgresql"
)

type flagModel struct {
	ProjectId string
}

var Cmd = &cobra.Command{
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

		// Show output as table
		table := tables.NewTable()
		table.SetHeader("ID", "NAME", "LAST_OPERATION.TYPE", "LAST_OPERATION.STATE")
		for _, i := range instances {
			table.AddRow(*i.InstanceId, *i.Name, *i.LastOperation.Type, *i.LastOperation.State)
		}
		table.Render(cmd)

		return nil
	},
}

func parseFlags(_ *cobra.Command) (*flagModel, error) {
	projectId := viper.GetString(config.ProjectIdKey)
	if projectId == "" {
		return nil, fmt.Errorf("project ID not set")
	}

	return &flagModel{
		ProjectId: projectId,
	}, nil
}

func buildRequest(ctx context.Context, model *flagModel, apiClient *postgresql.APIClient) postgresql.ApiGetInstancesRequest {
	req := apiClient.GetInstances(ctx, model.ProjectId)
	return req
}
