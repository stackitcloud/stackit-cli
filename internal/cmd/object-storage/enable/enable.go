package enable

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
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/object-storage/client"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/objectstorage"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
}

func NewCmd(params *params.CmdParams) *cobra.Command {
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
			model, err := parseInput(params.Printer, cmd, args)
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
				prompt := fmt.Sprintf("Are you sure you want to enable Object Storage for project %q?", projectLabel)
				err = params.Printer.PromptForConfirmation(prompt)
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
			params.Printer.Info("%s Object Storage for project %q\n", operationState, projectLabel)
			return nil
		},
	}
	return cmd
}

func parseInput(p *print.Printer, cmd *cobra.Command, _ []string) (*inputModel, error) {
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

func buildRequest(ctx context.Context, model *inputModel, apiClient *objectstorage.APIClient) objectstorage.ApiEnableServiceRequest {
	req := apiClient.EnableService(ctx, model.ProjectId, model.Region)
	return req
}
