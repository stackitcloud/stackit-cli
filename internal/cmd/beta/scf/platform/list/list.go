package list

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	scf "github.com/stackitcloud/stackit-sdk-go/services/scf/v1api"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	cliErr "github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/projectname"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/scf/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
)

const (
	limitFlag = "limit"

	// maxPageSize is the maximum number of items the API returns per page. 5000 is the maximum the API supports.
	maxPageSize = 5000
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	Limit *int64
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Lists all STACKIT Cloud Foundry platforms",
		Long:  fmt.Sprintf("Lists all STACKIT Cloud Foundry platforms. The number of fetched platforms is limited to %d.", maxPageSize),
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`List all STACKIT Cloud Foundry platforms`,
				"$ stackit beta scf platform list"),
			examples.NewExample(
				`List up to 10 STACKIT Cloud Foundry platforms`,
				"$ stackit beta scf platform list --limit 10"),
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
			resp, err := buildRequest(ctx, model, apiClient.DefaultAPI, 1, maxPageSize).Execute()
			if err != nil {
				return fmt.Errorf("fetch STACKIT Cloud Foundry platforms: %w", err)
			}

			// Get projectLabel
			projectLabel, err := projectname.GetProjectName(ctx, params.Printer, params.CliVersion, cmd)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get project name: %v", err)
				projectLabel = model.ProjectId
			} else if projectLabel == "" {
				projectLabel = model.ProjectId
			}

			return outputResult(params.Printer, model.OutputFormat, projectLabel, resp)
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

func buildRequest(ctx context.Context, model *inputModel, apiClient scf.DefaultAPI, page, perPage int64) scf.ApiListPlatformsRequest {
	// the API has non AIP-compliant pagination, which cannot be adapted with out SDK pagination utils.
	pageSize := int64(perPage)
	if model.Limit != nil && *model.Limit < pageSize {
		pageSize = *model.Limit
	}
	return apiClient.ListPlatforms(ctx, model.ProjectId, model.Region).Page(page).PerPage(pageSize)
}

func outputResult(p *print.Printer, outputFormat, projectLabel string, resp *scf.PlatformList) error {
	return p.OutputResult(outputFormat, resp, func() error {
		if resp == nil || len(resp.Resources) == 0 {
			p.Outputf("No STACKIT Cloud Foundry platforms found for project %q\n", projectLabel)
			return nil
		}
		if resp.Pagination.TotalPages != nil && *resp.Pagination.TotalPages > int64(1) {
			p.Warn("Server responded with more than 1 page. Output will not show all existing platforms.\n")
		}

		table := tables.NewTable()

		table.SetTitle("STACKIT Cloud Foundry platforms")
		table.SetHeader("ID", "NAME", "REGION", "SYSTEM ID")
		table.AddSeparator()
		for _, platform := range resp.Resources {
			table.AddRow(
				platform.Guid,
				platform.DisplayName,
				platform.Region,
				platform.SystemId,
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
