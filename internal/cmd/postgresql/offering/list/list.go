package list

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/commonflags"
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
	ProjectId string
	Limit     *int64
}

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List all PostgreSQL service offerings",
		Long:    "List all PostgreSQL service offerings",
		Example: `$ stackit postgresql offering list --project-id xxx`,
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
			if len(offerings) == 0 {
				cmd.Printf("No offerings found for project with ID %s\n", model.ProjectId)
				return nil
			}

			// Truncate output
			if model.Limit != nil && len(offerings) > int(*model.Limit) {
				offerings = offerings[:*model.Limit]
			}

			// Show output as table
			table := tables.NewTable()
			table.SetHeader("NAME", "PLAN.ID", "PLAN.NAME", "PLAN.DESCRIPTION")
			for i := range offerings {
				o := offerings[i]
				for j := range *o.Plans {
					p := (*o.Plans)[j]
					table.AddRow(*o.Name, *p.Id, *p.Name, *p.Description)
				}
				table.AddSeparator()
			}
			table.EnableAutoMergeOnColumns(1)
			table.Render(cmd)

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
	projectId := commonflags.GetString(commonflags.ProjectIdFlag)
	if projectId == "" {
		return nil, fmt.Errorf("project ID not set")
	}

	limit := utils.FlagToInt64Pointer(cmd, limitFlag)
	if limit != nil && *limit < 1 {
		return nil, fmt.Errorf("limit must be greater than 0")
	}

	return &flagModel{
		ProjectId: projectId,
		Limit:     limit,
	}, nil
}

func buildRequest(ctx context.Context, model *flagModel, apiClient *postgresql.APIClient) postgresql.ApiGetOfferingsRequest {
	req := apiClient.GetOfferings(ctx, model.ProjectId)
	return req
}
