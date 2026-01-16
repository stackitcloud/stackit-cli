package list

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	cliErr "github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/projectname"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/sfs/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/sfs"
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
		Short: "Lists all SFS resource pools",
		Long:  "Lists all SFS resource pools.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`List all SFS resource pools`,
				"$ stackit beta sfs resource-pool list"),
			examples.NewExample(
				`List all SFS resource pools for another region than the default one`,
				"$ stackit beta sfs resource-pool list --region eu01"),
			examples.NewExample(
				`List up to 10 SFS resource pools`,
				"$ stackit beta sfs resource-pool list --limit 10"),
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
			resp, err := buildRequest(ctx, model, apiClient).Execute()
			if err != nil {
				return fmt.Errorf("list SFS resource pools: %w", err)
			}

			resourcePools := utils.GetSliceFromPointer(resp.ResourcePools)

			// Truncate output
			if model.Limit != nil && len(resourcePools) > int(*model.Limit) {
				resourcePools = resourcePools[:*model.Limit]
			}

			// Get projectLabel
			projectLabel, err := projectname.GetProjectName(ctx, params.Printer, params.CliVersion, cmd)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get project name: %v", err)
				projectLabel = model.ProjectId
			} else if projectLabel == "" {
				projectLabel = model.ProjectId
			}

			return outputResult(params.Printer, model.OutputFormat, projectLabel, resourcePools)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Int64(limitFlag, 0, "Maximum number of entries to list")
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *sfs.APIClient) sfs.ApiListResourcePoolsRequest {
	req := apiClient.ListResourcePools(ctx, model.ProjectId, model.Region)
	return req
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
func outputResult(p *print.Printer, outputFormat, projectLabel string, resourcePools []sfs.ResourcePool) error {
	return p.OutputResult(outputFormat, resourcePools, func() error {
		if len(resourcePools) == 0 {
			p.Outputf("No resource pools found for project %q\n", projectLabel)
			return nil
		}

		table := tables.NewTable()
		table.SetHeader("ID", "NAME", "AVAILABILITY ZONE", "STATE", "TOTAL SIZE (GB)", "USED SIZE (GB)")
		for _, resourcePool := range resourcePools {
			totalSizeGigabytes, usedSizeGigabytes := "", ""
			if resourcePool.HasSpace() {
				totalSizeGigabytes = utils.PtrString(resourcePool.Space.SizeGigabytes)
				usedSizeGigabytes = utils.PtrString(resourcePool.Space.UsedGigabytes)
			}
			table.AddRow(
				utils.PtrString(resourcePool.Id),
				utils.PtrString(resourcePool.Name),
				utils.PtrString(resourcePool.AvailabilityZone),
				utils.PtrString(resourcePool.State),
				totalSizeGigabytes,
				usedSizeGigabytes,
			)
		}
		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}

		return nil
	})
}
