package enable

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
		Use:   "enable",
		Short: "Enables Object Storage for a project",
		Long:  "Enables Object Storage for a project.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Enable Object Storage functionality for your project`,
				"$ stackit object-storage enable"),
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
				prompt := fmt.Sprintf("Are you sure you want to enable Object Storage for project %q?", projectLabel)
				err = p.PromptForConfirmation(cmd, prompt)
				if err != nil {
					return err
				}
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			_, err = req.Execute()
			if err != nil {
				return fmt.Errorf("enable Object Storage: %w", err)
			}

			operationState := "Enabled"
			if model.Async {
				operationState = "Triggered enablement of"
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

func buildRequest(ctx context.Context, model *InputModel, apiClient *objectstorage.APIClient) objectstorage.ApiEnableServiceRequest {
	req := apiClient.EnableService(ctx, model.ProjectId)
	return req
}
