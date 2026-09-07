package list

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
	mongodbflex "github.com/stackitcloud/stackit-sdk-go/services/mongodbflex/v2api"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	cliErr "github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/mongodbflex/client"
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
		Short: "Lists MongoDB Flex flavors",
		Long:  "Lists MongoDB Flex flavors.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`List MongoDB Flex flavor`,
				"$ stackit mongodbflex flavor list"),
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
				return fmt.Errorf("get MongoDB Flex flavors: %w", err)
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

func buildRequest(ctx context.Context, model *inputModel, apiClient mongodbflex.DefaultAPI) mongodbflex.ApiListFlavorsRequest {
	return apiClient.ListFlavors(ctx, model.ProjectId, model.Region)
}

func outputResult(p *print.Printer, outputFormat string, flavors []mongodbflex.InstanceFlavor) error {
	return p.OutputResult(outputFormat, flavors, func() error {
		if len(flavors) == 0 {
			p.Outputf("No MongoDB flavors found.")
			return nil
		}

		table := tables.NewTable()
		table.SetTitle("Flavors")
		table.SetHeader("ID", "CPU", "MEMORY", "DESCRIPTION", "VALID INSTANCE TYPES")
		for _, f := range flavors {
			table.AddRow(
				utils.PtrString(f.Id),
				utils.PtrString(f.Cpu),
				utils.PtrString(f.Memory),
				utils.PtrString(f.Description),
				f.Categories,
			)
			table.AddSeparator()
		}

		return table.Display(p)
	})
}
