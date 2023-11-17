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
	Short:   "List all PostgreSQL service offerings",
	Long:    "List all PostgreSQL service offerings",
	Example: `$ stackit postgresql offerings list --project-id xxx`,
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
			return fmt.Errorf("get PostgreSQL service offerings: %w", err)
		}
		offerings := *resp.Offerings

		// Show output as table
		table := tables.NewTable()
		table.SetHeader("NAME", "PLAN.ID", "PLAN.NAME", "PLAN.DESCRIPTION")
		for _, o := range offerings {
			for _, p := range *o.Plans {
				table.AddRow(*o.Name, *p.Id, *p.Name, *p.Description)
			}
			table.AddSeparator()
		}
		table.EnableAutoMergeOnColumns(1)
		table.Render()

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

func buildRequest(ctx context.Context, model *flagModel, apiClient *postgresql.APIClient) postgresql.ApiGetOfferingsRequest {
	req := apiClient.GetOfferings(ctx, model.ProjectId)
	return req
}
