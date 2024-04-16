package delete

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/projectname"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/resourcemanager/client"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/resourcemanager"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
}

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Deletes a STACKIT project",
		Long:  "Deletes a STACKIT project.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Delete the configured STACKIT project`,
				"$ stackit project delete"),
			examples.NewExample(
				`Delete a STACKIT project by explicitly providing the project ID`,
				"$ stackit project delete --project-id xxx"),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			model, err := parseInput(cmd, p)
			if err != nil {
				return err
			}

			// Configure API client
			apiClient, err := client.ConfigureClient(p)
			if err != nil {
				return err
			}

			projectLabel, err := projectname.GetProjectName(ctx, cmd, p)
			if err != nil {
				p.Debug(print.ErrorLevel, "get project name: %v", err)
				projectLabel = model.ProjectId
			}

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to delete the project %q?", projectLabel)
				err = p.PromptForConfirmation(prompt)
				if err != nil {
					return err
				}
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			err = req.Execute()
			if err != nil {
				return fmt.Errorf("delete project: %w", err)
			}

			p.Info("Deleted project %q\n", projectLabel)
			p.Warn(fmt.Sprintf("%s\n%s\n",
				"If this was your default project, consider configuring a new project ID by running:",
				"  $ stackit config set --project-id <PROJECT_ID>",
			))
			return nil
		},
	}
	return cmd
}

func parseInput(cmd *cobra.Command, p *print.Printer) (*inputModel, error) {
	globalFlags := globalflags.Parse(cmd, p)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	return &inputModel{
		GlobalFlagModel: globalFlags,
	}, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *resourcemanager.APIClient) resourcemanager.ApiDeleteProjectRequest {
	req := apiClient.DeleteProject(ctx, model.ProjectId)
	return req
}
