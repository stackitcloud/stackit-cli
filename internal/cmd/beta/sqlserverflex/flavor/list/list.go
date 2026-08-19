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
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/sqlserverflex/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
)

const (
	limitFlag = "limit"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	Limit *int64
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Lists SQLServer Flex flavors",
		Long:  "Lists SQLServer Flex flavors.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`List SQLServer Flex flavors`,
				"$ stackit beta sqlserverflex flavor list"),
			examples.NewExample(
				`List up to 10 SQLServer Flex flavors`,
				"$ stackit beta sqlserverflex flavor list --limit 10"),
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
				return fmt.Errorf("get SQLServer Flex flavors: %w", err)
			}

			return outputResult(params.Printer, model.OutputFormat, flavors.Flavors)
		},
	}

	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Int64(limitFlag, 0, "Maximum number of entries to list")
}

func parseInput(p *print.Printer, cmd *cobra.Command, _ []string) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &cliErr.ProjectIdError{}
	}

	limit := flags.FlagToInt64Pointer(p, cmd, limitFlag)
	if limit != nil && *limit < 1 {
		return nil, &cliErr.FlagValidationError{
			Flag:    limitFlag,
			Details: "must be greater than 0",
		}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		Limit:           limit,
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient sqlserverflex.DefaultAPI) sqlserverflex.ApiListFlavorsRequest {
	req := apiClient.ListFlavors(ctx, model.ProjectId, model.Region)

	if model.Limit != nil {
		req = req.Size(*model.Limit)
	} else {
		// the default page size is only 10
		req = req.Size(100)
	}

	return req
}

func outputResult(p *print.Printer, outputFormat string, flavors []sqlserverflex.ListFlavors) error {
	return p.OutputResult(outputFormat, flavors, func() error {
		if len(flavors) == 0 {
			p.Outputf("No SQLServer Flex flavors found.")
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

		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}
		return nil
	})
}
