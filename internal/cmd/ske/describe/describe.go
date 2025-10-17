package describe

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/service-enablement/client"
	skeUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/service-enablement/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/serviceenablement"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
}

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "describe",
		Short: "Shows overall details regarding SKE",
		Long:  "Shows overall details regarding STACKIT Kubernetes Engine (SKE).",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Get details regarding SKE functionality on your project`,
				"$ stackit ske describe"),
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

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("read SKE project details: %w", err)
			}

			return outputResult(params.Printer, model.OutputFormat, resp, model.ProjectId)
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

func buildRequest(ctx context.Context, model *inputModel, apiClient *serviceenablement.APIClient) serviceenablement.ApiGetServiceStatusRegionalRequest {
	req := apiClient.GetServiceStatusRegional(ctx, model.Region, model.ProjectId, skeUtils.SKEServiceId)
	return req
}

func outputResult(p *print.Printer, outputFormat string, project *serviceenablement.ServiceStatus, projectId string) error {
	if project == nil {
		return fmt.Errorf("project is nil")
	}

	return p.OutputResult(outputFormat, project, func() error {
		table := tables.NewTable()
		table.AddRow("ID", projectId)
		table.AddSeparator()
		if project.HasState() {
			table.AddRow("STATE", utils.PtrString(project.State))
		}
		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}

		return nil
	})
}
