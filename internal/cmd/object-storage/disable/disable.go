package disable

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/projectname"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/object-storage/client"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/objectstorage"
)

type InputModel struct {
	*globalflags.GlobalFlagModel
}

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "disable",
		Short: "Disables Object Storage for a project",
		Long:  "Disables Object Storage for a project. All buckets must be deleted beforehand.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Disable Object Storage functionality for your project.`,
				"$ stackit object-storage disable"),
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

			projectLabel, err := projectname.GetProjectName(ctx, cmd)
			if err != nil {
				projectLabel = model.ProjectId
			}

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to disable Object Storage for project %q?", projectLabel)
				err = p.PromptForConfirmation(prompt)
				if err != nil {
					return err
				}
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			_, err = req.Execute()
			if err != nil {
				return fmt.Errorf("disable Object Storage: %w", err)
			}

			operationState := "Disabled"
			if model.Async {
				operationState = "Triggered disablement of"
			}
			p.Info("%s Object Storage for project %q\n", operationState, projectLabel)
			return nil
		},
	}
	return cmd
}

func parseInput(cmd *cobra.Command) (*InputModel, error) {
	globalFlags := globalflags.Parse(cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	return &InputModel{
		GlobalFlagModel: globalFlags,
	}, nil
}

func buildRequest(ctx context.Context, model *InputModel, apiClient *objectstorage.APIClient) objectstorage.ApiDisableServiceRequest {
	req := apiClient.DisableService(ctx, model.ProjectId)
	return req
}
