package lock

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/projectname"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/object-storage/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/object-storage/utils"

	"github.com/spf13/cobra"
	objectstorage "github.com/stackitcloud/stackit-sdk-go/services/objectstorage/v2api"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "lock",
		Short: "Create object storage compliance lock",
		Long:  "Create object storage compliance lock.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Create object storage compliance lock`,
				"$ stackit object-storage compliance-lock lock"),
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

			prompt := fmt.Sprintf("Are you sure you want to create object storage compliance-lock for project %s?", projectLabel)
			err = params.Printer.PromptForConfirmation(prompt)
			if err != nil {
				return err
			}

			// Check if the project is enabled before trying to create
			enabled, err := utils.ProjectEnabled(ctx, apiClient.DefaultAPI, model.ProjectId, model.Region)
			if err != nil {
				return fmt.Errorf("check if Object Storage is enabled: %w", err)
			}
			if !enabled {
				return &errors.ServiceDisabledError{
					Service: "object-storage",
				}
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("create object storage compliance lock: %w", err)
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

func buildRequest(ctx context.Context, model *inputModel, apiClient *objectstorage.APIClient) objectstorage.ApiCreateComplianceLockRequest {
	req := apiClient.DefaultAPI.CreateComplianceLock(ctx, model.ProjectId, model.Region)
	return req
}

func outputResult(p *print.Printer, outputFormat, projectLabel string, resp *objectstorage.ComplianceLockResponse) error {
	return p.OutputResult(outputFormat, resp, func() error {
		if resp == nil {
			return fmt.Errorf("create compliance lock response is empty")
		}

		p.Outputf("Created object storage compliance lock for project \"%s\" with maximum retention period of %d days.\n", projectLabel, resp.MaxRetentionDays)
		return nil
	})
}
