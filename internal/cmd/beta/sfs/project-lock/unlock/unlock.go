package unlock

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/projectname"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"

	"github.com/spf13/cobra"
	sfs "github.com/stackitcloud/stackit-sdk-go/services/sfs/v1api"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/sfs/client"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "unlock",
		Short: "Clean up lock for a project",
		Long:  "Clean up lock for a project.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Disable lock for project`,
				"$ stackit beta sfs project-lock unlock"),
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
			} else if projectLabel == "" {
				projectLabel = model.ProjectId
			}

			prompt := fmt.Sprintf("Are you sure you want to disable lock for project %s?", projectLabel)
			err = params.Printer.PromptForConfirmation(prompt)
			if err != nil {
				return err
			}

			// Call API
			_, err = buildRequest(ctx, model, apiClient).Execute()
			if err != nil {
				return fmt.Errorf("disable project lock: %w", err)
			}

			params.Printer.Outputf("Project \"%s\" is successfully unlocked.\n", projectLabel)

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

func buildRequest(ctx context.Context, model *inputModel, apiClient *sfs.APIClient) sfs.ApiDisableLockRequest {
	req := apiClient.DefaultAPI.DisableLock(ctx, model.Region, model.ProjectId)
	return req
}
