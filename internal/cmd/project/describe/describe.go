package describe

import (
	"context"
	"encoding/json"
	"fmt"

	"stackit/internal/pkg/args"
	"stackit/internal/pkg/errors"
	"stackit/internal/pkg/examples"
	"stackit/internal/pkg/flags"
	"stackit/internal/pkg/globalflags"
	"stackit/internal/pkg/services/resourcemanager/client"
	"stackit/internal/pkg/tables"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/resourcemanager"
)

const (
	includeParentsFlag = "include-parents"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	IncludeParents bool
}

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "describe",
		Short: "Get the details of a STACKIT project",
		Long:  "Get the details of a STACKIT project",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Get the details of the configured STACKIT project`,
				"$ stackit project describe"),
			examples.NewExample(
				`Get the details of a STACKIT project by explicitly providing the project ID`,
				"$ stackit project describe --project-id xxx"),
			examples.NewExample(
				`Get the details of the configured STACKIT project, including details of the parent resources`,
				"$ stackit project describe --include-parents"),
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
				return fmt.Errorf("read project details: %w", err)
			}

			return outputResult(cmd, model.OutputFormat, resp)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Bool(includeParentsFlag, false, "When true, the details of the parent resources will be included in the output")
}

func parseInput(cmd *cobra.Command) (*inputModel, error) {
	globalFlags := globalflags.Parse(cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	return &inputModel{
		GlobalFlagModel: globalFlags,
		IncludeParents:  flags.FlagToBoolValue(cmd, includeParentsFlag),
	}, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *resourcemanager.APIClient) resourcemanager.ApiGetProjectRequest {
	req := apiClient.GetProject(ctx, model.ProjectId)
	req.IncludeParents(model.IncludeParents)
	return req
}

func outputResult(cmd *cobra.Command, outputFormat string, project *resourcemanager.ProjectResponseWithParents) error {
	switch outputFormat {
	case globalflags.PrettyOutputFormat:
		table := tables.NewTable()
		table.AddRow("ID", *project.ProjectId)
		table.AddSeparator()
		table.AddRow("NAME", *project.Name)
		table.AddSeparator()
		table.AddRow("CREATION", *project.CreationTime)
		table.AddSeparator()
		table.AddRow("STATE", *project.LifecycleState)
		table.AddSeparator()
		table.AddRow("PARENT ID", *project.Parent.Id)
		err := table.Display(cmd)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}

		return nil
	default:
		details, err := json.MarshalIndent(project, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal project details: %w", err)
		}
		cmd.Println(string(details))

		return nil
	}
}
