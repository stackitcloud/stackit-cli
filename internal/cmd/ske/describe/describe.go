package describe

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/ske/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/ske"
)

type flagModel struct {
	*globalflags.GlobalFlagModel
}

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "describe",
		Short:   "Get overall details regarding SKE",
		Long:    "Get overall details regarding SKE",
		Example: `$ stackit ske describe --project-id xxx`,
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
				return fmt.Errorf("read SKE project details: %w", err)
			}

			return outputResult(cmd, model.OutputFormat, resp)
		},
	}
	return cmd
}

func parseFlags(cmd *cobra.Command) (*flagModel, error) {
	globalFlags := globalflags.Parse(cmd)
	if globalFlags.ProjectId == "" {
		return nil, fmt.Errorf("project ID not set")
	}

	return &flagModel{
		GlobalFlagModel: globalFlags,
	}, nil
}

func buildRequest(ctx context.Context, model *flagModel, apiClient *ske.APIClient) ske.ApiGetProjectRequest {
	req := apiClient.GetProject(ctx, model.ProjectId)
	return req
}

func outputResult(cmd *cobra.Command, outputFormat string, project *ske.ProjectResponse) error {
	switch outputFormat {
	case globalflags.TableOutputFormat:
		table := tables.NewTable()
		table.SetHeader("ID", "STATE")
		table.AddRow(*project.ProjectId, *project.State)
		err := table.Render(cmd)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}

		return nil
	default:
		details, err := json.MarshalIndent(project, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal SKE project details: %w", err)
		}
		cmd.Println(string(details))

		return nil
	}
}
