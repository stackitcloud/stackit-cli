package list

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/ske/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/ske"
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
		Short:   "List all SKE clusters",
		Long:    "List all SKE clusters",
		Example: `$ stackit ske cluster list --project-id xxx`,
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
				return fmt.Errorf("get SKE clusters: %w", err)
			}
			clusters := *resp.Items
			if len(clusters) == 0 {
				cmd.Printf("No clusters found for project with ID %s\n", model.ProjectId)
				return nil
			}

			// Truncate output
			if model.Limit != nil && len(clusters) > int(*model.Limit) {
				clusters = clusters[:*model.Limit]
			}

			return outputResult(cmd, model.OutputFormat, clusters)
		},
	}

	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Int64(limitFlag, 0, "Maximum number of entries to list")
}

func parseFlags(cmd *cobra.Command) (*flagModel, error) {
	globalFlags := globalflags.Parse(cmd)
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

func buildRequest(ctx context.Context, model *flagModel, apiClient *ske.APIClient) ske.ApiGetClustersRequest {
	req := apiClient.GetClusters(ctx, model.ProjectId)
	return req
}

func outputResult(cmd *cobra.Command, outputFormat string, clusters []ske.ClusterResponse) error {
	switch outputFormat {
	case globalflags.JSONOutputFormat:
		details, err := json.MarshalIndent(clusters, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal SKE cluster list: %w", err)
		}
		cmd.Println(string(details))

		return nil
	default:
		table := tables.NewTable()
		table.SetHeader("NAME", "STATE")
		for i := range clusters {
			c := clusters[i]
			table.AddRow(*c.Name, *c.Status.Aggregated)
		}
		err := table.Render(cmd)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}

		return nil
	}
}
