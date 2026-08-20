package list

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/types"

	"github.com/spf13/cobra"
	sqlserverflex "github.com/stackitcloud/stackit-sdk-go/services/sqlserverflex/v3api"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	cliErr "github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/sqlserverflex/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Lists SQLServer Flex versions",
		Long:  "Lists SQLServer Flex versions.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`List SQLServer Flex version options`,
				"$ stackit sqlserverflex version list"),
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
			versions, err := buildRequest(ctx, model, apiClient.DefaultAPI).Execute()
			if err != nil {
				return fmt.Errorf("get SQLServer Flex versions: %w", err)
			}

			return outputResult(params.Printer, model.OutputFormat, versions.Versions)
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

func buildRequest(ctx context.Context, model *inputModel, apiClient sqlserverflex.DefaultAPI) sqlserverflex.ApiListVersionsRequest {
	return apiClient.ListVersions(ctx, model.ProjectId, model.Region)
}

func outputResult(p *print.Printer, outputFormat string, versions []sqlserverflex.Version) error {
	return p.OutputResult(outputFormat, versions, func() error {
		if len(versions) == 0 {
			p.Outputf("No SQLServer Flex versions found.")
			return nil
		}

		table := tables.NewTable()
		table.SetTitle("Versions")
		table.SetHeader("VERSION", "RECOMMEND", "BETA", "DEPRECATED")

		for _, v := range versions {
			table.AddRow(v.Version, v.Recommend, v.Beta, v.Deprecated)
			table.AddSeparator()
		}

		return table.Display(p)
	})
}
