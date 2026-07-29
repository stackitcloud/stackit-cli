package list

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/types"

	"github.com/spf13/cobra"
	postgresflex "github.com/stackitcloud/stackit-sdk-go/services/postgresflex/v3api"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	cliErr "github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/postgresflex/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Lists PostgreSQL Flex flavors",
		Long:  "Lists PostgreSQL Flex flavors.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`List PostgreSQL Flex flavor`,
				"$ stackit postgresflex flavor list"),
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
			flavors, err := buildRequest(ctx, model, apiClient.DefaultAPI).Execute()
			if err != nil {
				return fmt.Errorf("get PostgreSQL Flex flavors: %w", err)
			}

			return outputResult(params.Printer, model.OutputFormat, flavors.Flavors)
		},
	}
	return cmd
}

func parseInput(p *print.Printer, cmd *cobra.Command, _ []string) (*inputModel, error) {
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

func buildRequest(ctx context.Context, model *inputModel, apiClient postgresflex.DefaultAPI) postgresflex.ApiListFlavorsRequest {
	return apiClient.ListFlavors(ctx, model.ProjectId, model.Region).Size(100)
}

func outputResult(p *print.Printer, outputFormat string, flavors []postgresflex.ListFlavors) error {
	return p.OutputResult(outputFormat, flavors, func() error {
		if len(flavors) == 0 {
			p.Outputf("No PostgreSQL flavors found.")
			return nil
		}

		table := tables.NewTable()
		table.SetTitle("Flavors")
		table.SetHeader("ID", "CPU", "MEMORY", "NODE TYPE", "DESCRIPTION")
		for _, f := range flavors {
			table.AddRow(
				f.Id,
				f.Cpu,
				f.Memory,
				f.NodeType,
				f.Description,
			)
			table.AddSeparator()
		}

		return table.Display(p)
	})
}
