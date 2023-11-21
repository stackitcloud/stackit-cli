package list

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/config"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/ske/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stackitcloud/stackit-sdk-go/services/ske"
)

const (
	projectIdFlag = "project-id"
)

type flagModel struct {
	ProjectId string
}

var Cmd = &cobra.Command{
	Use:     "list",
	Short:   "List all SKE clusters",
	Long:    "List all SKE clusters",
	Example: `$ stackit ske cluster list --project-id xxx`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		model, err := parseFlags()
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
			return fmt.Errorf("get SKE clusters: %w", err)
		}
		clusters := *resp.Items
		if len(clusters) == 0 {
			cmd.Printf("No clusters found for project with ID %s\n", model.ProjectId)
			return nil
		}

		// Show output as table
		table := tables.NewTable()
		table.SetHeader("NAME", "STATE")
		for i := range clusters {
			c := clusters[i]
			table.AddRow(*c.Name, *c.Status.Aggregated)
		}
		table.Render(cmd)

		return nil
	},
}

func parseFlags() (*flagModel, error) {
	projectId := viper.GetString(config.ProjectIdKey)
	if projectId == "" {
		return nil, fmt.Errorf("project ID not set")
	}

	return &flagModel{
		ProjectId: projectId,
	}, nil
}

func buildRequest(ctx context.Context, model *flagModel, apiClient *ske.APIClient) ske.ApiGetClustersRequest {
	req := apiClient.GetClusters(ctx, model.ProjectId)
	return req
}
