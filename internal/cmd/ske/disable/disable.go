package disable

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/config"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/ske/client"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stackitcloud/stackit-sdk-go/services/ske"
	"github.com/stackitcloud/stackit-sdk-go/services/ske/wait"
)

const (
	projectIdFlag = "project-id"
)

type FlagModel struct {
	ProjectId string
}

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "disable",
		Short:   "Disables SKE for a project",
		Long:    "Disables SKE for a project",
		Example: `$ stackit ske disable --project-id xxx`,
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
			_, err = req.Execute()
			if err != nil {
				return fmt.Errorf("disable SKE: %w", err)
			}

			// Wait for async operation
			_, err = wait.DeleteProjectWaitHandler(ctx, apiClient, model.ProjectId).WaitWithContext(ctx)
			if err != nil {
				return fmt.Errorf("wait for SKE disabling: %w", err)
			}

			cmd.Printf("SKE disabled")
			return nil
		},
	}
	return cmd
}

func parseFlags() (*FlagModel, error) {
	projectId := viper.GetString(config.ProjectIdKey)
	if projectId == "" {
		return nil, fmt.Errorf("project ID not set")
	}

	return &FlagModel{
		ProjectId: projectId,
	}, nil
}

func buildRequest(ctx context.Context, model *FlagModel, apiClient *ske.APIClient) ske.ApiDeleteProjectRequest {
	req := apiClient.DeleteProject(ctx, model.ProjectId)
	return req
}
