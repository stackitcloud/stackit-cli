package describe

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/config"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/ske/client"

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

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "describe",
		Short:   "Get overall details regarding SKE",
		Long:    "Get overall details regarding SKE",
		Example: `$ stackit ske describe --project-id xxx`,
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
				return fmt.Errorf("read SKE project details: %w", err)
			}

			// Show details
			details, err := json.MarshalIndent(resp, "", "  ")
			if err != nil {
				return fmt.Errorf("marshal SKE project details: %w", err)
			}
			cmd.Println(string(details))

			return nil
		},
	}
	return cmd
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

func buildRequest(ctx context.Context, model *flagModel, apiClient *ske.APIClient) ske.ApiGetProjectRequest {
	req := apiClient.GetProject(ctx, model.ProjectId)
	return req
}
