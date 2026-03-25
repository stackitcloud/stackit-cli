package describe

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/object-storage/client"
	objectStorageUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/object-storage/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"

	"github.com/spf13/cobra"
	objectstorage "github.com/stackitcloud/stackit-sdk-go/services/objectstorage/v2api"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "describe",
		Short: "Describe object storage compliance lock",
		Long:  "Describe object storage compliance lock.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Describe object storage compliance lock`,
				"$ stackit object-storage compliance-lock describe"),
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

			// Check if the project is enabled before trying to describe
			enabled, err := objectStorageUtils.ProjectEnabled(ctx, apiClient.DefaultAPI, model.ProjectId, model.Region)
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
				return fmt.Errorf("get object storage compliance lock: %w", err)
			}

			return outputResult(params.Printer, model.OutputFormat, resp)
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

func buildRequest(ctx context.Context, model *inputModel, apiClient *objectstorage.APIClient) objectstorage.ApiGetComplianceLockRequest {
	req := apiClient.DefaultAPI.GetComplianceLock(ctx, model.ProjectId, model.Region)
	return req
}

func outputResult(p *print.Printer, outputFormat string, resp *objectstorage.ComplianceLockResponse) error {
	return p.OutputResult(outputFormat, resp, func() error {
		if resp == nil {
			return fmt.Errorf("response is empty")
		}

		table := tables.NewTable()
		table.AddRow("PROJECT ID", resp.Project)
		table.AddSeparator()
		table.AddRow("MAX RETENTION DAYS", resp.MaxRetentionDays)
		table.AddSeparator()

		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}

		return nil
	})
}
