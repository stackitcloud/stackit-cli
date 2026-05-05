package lock

import (
	"context"
	sysErrors "errors"
	"fmt"
	"net/http"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/core/oapierror"
	sfs "github.com/stackitcloud/stackit-sdk-go/services/sfs/v1api"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/projectname"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/sfs/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "lock",
		Short: "Enables lock for a project",
		Long:  "Enables lock for a project. Necessary for immutable snapshots and to prevent accidental deletion of resources.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Enable lock for project`,
				"$ stackit beta sfs project-lock lock",
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			model, err := parseInput(params.Printer, cmd, args)
			if err != nil {
				return fmt.Errorf("unable to parse input: %w", err)
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

			prompt := fmt.Sprintf("Are you sure you want to enable SFS lock for project %s?", projectLabel)
			err = params.Printer.PromptForConfirmation(prompt)
			if err != nil {
				return err
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				var oApiErr *oapierror.GenericOpenAPIError
				if sysErrors.As(err, &oApiErr) {
					if oApiErr.StatusCode == http.StatusConflict {
						params.Printer.Info("There is already an active lock for project %s\n", projectLabel)
						return err
					}
				}

				return fmt.Errorf("enable SFS project lock: %w", err)
			}

			return outputResult(params.Printer, model.OutputFormat, projectLabel, resp)
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

func buildRequest(ctx context.Context, model *inputModel, apiClient *sfs.APIClient) sfs.ApiEnableLockRequest {
	req := apiClient.DefaultAPI.EnableLock(ctx, model.Region, model.ProjectId)
	return req
}

func outputResult(p *print.Printer, outputFormat, projectLabel string, resp *sfs.EnableLockResponse) error {
	return p.OutputResult(outputFormat, resp, func() error {
		if resp == nil {
			return fmt.Errorf("enable project lock response is empty")
		}

		p.Outputf("Project %q is successfully locked.\n", projectLabel)
		return nil
	})
}
