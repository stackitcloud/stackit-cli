package list

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	cliErr "github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/modelexperiments/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"

	modelexperiments "github.com/stackitcloud/stackit-sdk-go/services/modelexperiments/v1api"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Lists AI Model Experiments instances",
		Long:  "Lists AI Model Experiments instances in a STACKIT project.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`List AI Model Experiments instances in a project`,
				`$ stackit ai-model-experiments instance list --region eu01`,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()

			model, err := parseInput(params.Printer, cmd, args)
			if err != nil {
				return err
			}

			apiClient, err := client.ConfigureClient(params.Printer, params.CliVersion)
			if err != nil {
				return err
			}

			req := buildListInstancesRequest(ctx, model, apiClient)

			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("list AI Model Experiments instances: %w", err)
			}

			return outputResult(params.Printer, model.OutputFormat, resp)
		},
	}

	return cmd
}

func parseInput(
	p *print.Printer,
	cmd *cobra.Command,
	_ []string,
) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)

	if globalFlags.ProjectId == "" {
		return nil, &cliErr.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
	}

	p.DebugInputModel(model)

	return &model, nil
}

func buildListInstancesRequest(
	ctx context.Context,
	model *inputModel,
	apiClient *modelexperiments.APIClient,
) modelexperiments.ApiListInstancesRequest {
	return apiClient.DefaultAPI.ListInstances(
		ctx,
		model.ProjectId,
		model.Region,
	)
}

func outputResult(
	p *print.Printer,
	outputFormat string,
	resp *modelexperiments.ListInstancesResponse,
) error {
	if resp == nil {
		return fmt.Errorf("response is nil")
	}

	return p.OutputResult(outputFormat, resp.Instances, func() error {
		if len(resp.Instances) == 0 {
			p.Outputf("No instances found\n")
			return nil
		}

		table := tables.NewTable()
		table.SetHeader("ID", "NAME", "REGION", "STATUS")

		for _, instance := range resp.Instances {
			table.AddRow(
				instance.Id,
				instance.Name,
				instance.Region,
				instance.State,
			)
		}

		return table.Display(p)
	})
}
