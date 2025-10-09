package disable

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
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

func NewCmd(params *params.CmdParams) *cobra.Command {
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
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := context.Background()
			model, err := parseInput(params.Printer, cmd)
			if err != nil {
				return err
			}

			// Configure API client
			apiClient, err := client.ConfigureClient(params.Printer, params.CliVersion)
			if err != nil {
				return err
			}

			projectLabel, err := projectname.GetProjectName(ctx, params.Printer, params.CliVersion, cmd)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get project name: %v", err)
				projectLabel = model.ProjectId
			}

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to disable SKE for project %q? (This will delete all associated clusters)", projectLabel)
				err = params.Printer.PromptForConfirmation(prompt)
				if err != nil {
					return err
				}
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			err = req.Execute()
			if err != nil {
				return fmt.Errorf("disable SKE: %w", err)
			}

			// Wait for async operation, if async mode not enabled
			if !model.Async {
				s := spinner.New(params.Printer)
				s.Start("Disabling SKE")
				_, err = wait.DisableServiceWaitHandler(ctx, apiClient, model.Region, model.ProjectId, utils.SKEServiceId).WaitWithContext(ctx)
				if err != nil {
					return fmt.Errorf("wait for SKE disabling: %w", err)
				}
				s.Stop()
			}

			operationState := "Disabled"
			if model.Async {
				operationState = "Triggered disablement of"
			}
			params.Printer.Info("%s SKE for project %q\n", operationState, projectLabel)
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

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *serviceenablement.APIClient) serviceenablement.ApiDisableServiceRegionalRequest {
	req := apiClient.DisableServiceRegional(ctx, model.Region, model.ProjectId, utils.SKEServiceId)
	return req
}
