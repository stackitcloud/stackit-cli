package disable

import (
	"context"
	"fmt"

	"stackit/internal/pkg/args"
	"stackit/internal/pkg/confirm"
	"stackit/internal/pkg/errors"
	"stackit/internal/pkg/examples"
	"stackit/internal/pkg/globalflags"
	"stackit/internal/pkg/projectname"
	"stackit/internal/pkg/services/ske/client"
	"stackit/internal/pkg/spinner"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/ske"
	"github.com/stackitcloud/stackit-sdk-go/services/ske/wait"
)

type InputModel struct {
	*globalflags.GlobalFlagModel
}

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "disable",
		Short: "Disables SKE for a project",
		Long:  "Disables STACKIT Kubernetes Engine (SKE) for a project. It will delete all associated clusters",
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
			apiClient, err := client.ConfigureClient(cmd)
			if err != nil {
				return err
			}

			projectLabel, err := projectname.GetProjectName(ctx, cmd)
			if err != nil {
				projectLabel = model.ProjectId
			}

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to disable SKE for project %s? (This will delete all associated clusters)", projectLabel)
				err = confirm.PromptForConfirmation(cmd, prompt)
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
				s := spinner.New(cmd)
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
			cmd.Printf("%s SKE for project %s\n", operationState, projectLabel)
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
