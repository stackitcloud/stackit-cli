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
	flavorIdFlag = "flavor-id"
	limitFlag    = "limit"
)

type inputModel struct {
	*globalflags.GlobalFlagModel

	FlavorId *string
	Limit    *int64
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Lists MongoDB Flex storages for a certain flavor",
		Long:  "Lists MongoDB Flex storages for a certain flavor.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`List MongoDB Flex storages for flavor with ID "xxx"`,
				"$ stackit mongodbflex storage list --flavor-id xxx"),
			examples.NewExample(
				`List MongoDB Flex storages for flavor with ID "xxx" in JSON format`,
				"$ stackit mongodbflex storage list --flavor-id xxx --output-format json"),
			examples.NewExample(
				`List up to 10 storages for flavor with ID "xxx"`,
				"$ stackit mongodbflex storage list --flavor-id xxx --limit 10"),
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
			storages, err := buildRequest(ctx, model, apiClient.DefaultAPI).Execute()
			if err != nil {
				return fmt.Errorf("get MongoDB Flex storages: %w", err)
			}

			return outputResult(params.Printer, model.OutputFormat, storages)
		},
	}

	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().String(flavorIdFlag, "", "Flavor ID")
	cmd.Flags().Int64(limitFlag, 0, "Maximum number of entries to list")

	err := flags.MarkFlagsRequired(cmd, flavorIdFlag)
	cobra.CheckErr(err)
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
		FlavorId:        flags.FlagToStringPointer(p, cmd, flavorIdFlag),
		Limit:           limit,
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient mongodbflex.DefaultAPI) mongodbflex.ApiListStoragesRequest {
	return apiClient.ListStorages(ctx, model.ProjectId, *model.FlavorId, model.Region)
}

func outputResult(p *print.Printer, outputFormat string, storagesResp *mongodbflex.ListStoragesResponse) error {
	return p.OutputResult(outputFormat, storagesResp, func() error {
		if storagesResp == nil {
			return fmt.Errorf("storages resp is empty")
		}
		storages := storagesResp.StorageClasses
		if len(storages) == 0 {
			p.Outputf("No MongoDB Flex storages found.")
			return nil
		}

		table := tables.NewTable()
		table.SetTitle("Storages")
		table.SetHeader("MINIMUM", "MAXIMUM", "STORAGE CLASS")
		for _, storageClass := range storages {
			table.AddRow(
				utils.PtrString(storagesResp.StorageRange.Min),
				utils.PtrString(storagesResp.StorageRange.Max),
				storageClass,
			)
			table.AddSeparator()
		}

		return table.Display(p)
	})
}
