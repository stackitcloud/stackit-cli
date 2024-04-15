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
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/ske/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/spinner"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/ske"
	"github.com/stackitcloud/stackit-sdk-go/services/ske/wait"
)

type InputModel struct {
	*globalflags.GlobalFlagModel
}

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "disable",
		Short: "Disables SKE for a project",
		Long:  "Disables STACKIT Kubernetes Engine (SKE) for a project. It will delete all associated clusters.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Disable SKE functionality for your project, deleting all associated clusters`,
				"$ stackit ske disable"),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			model, err := parseInput(cmd)
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
				prompt := fmt.Sprintf("Are you sure you want to disable SKE for project %q? (This will delete all associated clusters)", projectLabel)
				err = p.PromptForConfirmation(prompt)
				if err != nil {
					return err
				}
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			_, err = req.Execute()
			if err != nil {
				return fmt.Errorf("disable SKE: %w", err)
			}

			// Wait for async operation, if async mode not enabled
			if !model.Async {
				s := spinner.New(p)
				s.Start("Disabling SKE")
				_, err = wait.DisableServiceWaitHandler(ctx, apiClient, model.ProjectId).WaitWithContext(ctx)
				if err != nil {
					return fmt.Errorf("wait for SKE disabling: %w", err)
				}
				s.Stop()
			}

			operationState := "Disabled"
			if model.Async {
				operationState = "Triggered disablement of"
			}
			p.Info("%s SKE for project %q\n", operationState, projectLabel)
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

func buildRequest(ctx context.Context, model *InputModel, apiClient *ske.APIClient) ske.ApiDisableServiceRequest {
	req := apiClient.DisableService(ctx, model.ProjectId)
	return req
}
