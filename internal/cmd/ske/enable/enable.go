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
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/service-enablement/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/service-enablement/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/spinner"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/serviceenablement"
	"github.com/stackitcloud/stackit-sdk-go/services/serviceenablement/wait"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
}

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "enable",
		Short: "Enables SKE for a project",
		Long:  "Enables STACKIT Kubernetes Engine (SKE) for a project.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Enable SKE functionality for your project`,
				"$ stackit ske enable"),
		),
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := context.Background()
			model, err := parseInput(p, cmd)
			if err != nil {
				return err
			}

			// Configure API client
			apiClient, err := client.ConfigureClient(p)
			if err != nil {
				return err
			}

			projectLabel, err := projectname.GetProjectName(ctx, p, cmd)
			if err != nil {
				p.Debug(print.ErrorLevel, "get project name: %v", err)
				projectLabel = model.ProjectId
			}

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to enable SKE for project %q?", projectLabel)
				err = p.PromptForConfirmation(prompt)
				if err != nil {
					return err
				}
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			err = req.Execute()
			if err != nil {
				return fmt.Errorf("enable SKE: %w", err)
			}

			// Wait for async operation, if async mode not enabled
			if !model.Async {
				s := spinner.New(p)
				s.Start("Enabling SKE")
				_, err = wait.EnableServiceWaitHandler(ctx, apiClient, model.Region, model.ProjectId, utils.SKEServiceId).WaitWithContext(ctx)
				if err != nil {
					return fmt.Errorf("wait for SKE enabling: %w", err)
				}
				s.Stop()
			}

			operationState := "Enabled"
			if model.Async {
				operationState = "Triggered enablement of"
			}
			p.Info("%s SKE for project %q\n", operationState, projectLabel)
			return nil
		},
	}
	return cmd
}

func parseInput(p *print.Printer, cmd *cobra.Command) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
	}

	if p.IsVerbosityDebug() {
		modelStr, err := print.BuildDebugStrFromInputModel(model)
		if err != nil {
			p.Debug(print.ErrorLevel, "convert model to string for debugging: %v", err)
		} else {
			p.Debug(print.DebugLevel, "parsed input values: %s", modelStr)
		}
	}

	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *serviceenablement.APIClient) serviceenablement.ApiEnableServiceRegionalRequest {
	req := apiClient.EnableServiceRegional(ctx, model.Region, model.ProjectId, utils.SKEServiceId)
	return req
}
