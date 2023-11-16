package describe

import (
	"context"
	"encoding/json"
	"fmt"

	"stackit/internal/pkg/args"
	"stackit/internal/pkg/errors"
	"stackit/internal/pkg/examples"
	"stackit/internal/pkg/globalflags"
	"stackit/internal/pkg/services/ske/client"
	"stackit/internal/pkg/tables"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/ske"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
}

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "describe",
		Short: "Get overall details regarding SKE",
		Long:  "Get overall details regarding STACKIT Kubernetes Engine (SKE)",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Get details regarding SKE functionality on your project`,
				"$ stackit ske describe"),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			model, err := parseInput(cmd)
			if err != nil {
				return err
			}
			// Configure API client
			apiClient, err := client.ConfigureClient(cmd)
			if err != nil {
				return err
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

func parseInput(cmd *cobra.Command) (*inputModel, error) {
	globalFlags := globalflags.Parse(cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	return &inputModel{
		GlobalFlagModel: globalFlags,
	}, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *ske.APIClient) ske.ApiGetServiceStatusRequest {
	req := apiClient.GetServiceStatus(ctx, model.ProjectId)
	return req
}

func outputResult(cmd *cobra.Command, outputFormat string, project *ske.ProjectResponse) error {
	switch outputFormat {
	case globalflags.PrettyOutputFormat:
		table := tables.NewTable()
		table.AddRow("ID", *project.ProjectId)
		table.AddSeparator()
		table.AddRow("STATE", *project.State)
		err := table.Display(cmd)
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
