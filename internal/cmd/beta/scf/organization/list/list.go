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
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
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
		Short: "Lists all STACKIT Cloud Foundry organizations",
		Long:  fmt.Sprintf("Lists all STACKIT Cloud Foundry organizations. The number of fetched organizations is limited to %d.", maxPageSize),
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`List all STACKIT Cloud Foundry organizations`,
				"$ stackit beta scf organization list"),
			examples.NewExample(
				`List up to 10 STACKIT Cloud Foundry organizations`,
				"$ stackit beta scf organization list --limit 10"),
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
			resp, err := fetchOrganizations(ctx, params.Printer, model, maxPageSize, apiClient.DefaultAPI)
			if err != nil {
				return fmt.Errorf("fetch STACKIT Cloud Foundry organizations: %w", err)
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

func buildRequest(ctx context.Context, model *inputModel, apiClient scf.DefaultAPI, page, perPage int64) scf.ApiListOrganizationsRequest {
	return apiClient.ListOrganizations(ctx, model.ProjectId, model.Region).WithDetails(true).Page(page).PerPage(perPage)
}

// fetchOrganizations fetches the organizations without using the non AIP-compliant pagination.
func fetchOrganizations(ctx context.Context, p *print.Printer, model *inputModel, pageSize int, apiClient scf.DefaultAPI) ([]scf.OrganizationsListItem, error) {
	if pageSize < 1 {
		return nil, fmt.Errorf("list STACKIT Cloud Foundry organizations: maxPageSize must not be smaller than 1")
	}
	perPage := int64(pageSize)
	if model.Limit != nil && *model.Limit < perPage {
		perPage = *model.Limit
	}
	resp, err := buildRequest(ctx, model, apiClient, 1, perPage).Execute()
	if err != nil {
		return nil, fmt.Errorf("list STACKIT Cloud Foundry organizations: %w", err)
	}
	organizations := resp.Resources
	if model.Limit != nil && int64(len(resp.Resources)) >= *model.Limit {
		organizations = resp.Resources[:*model.Limit]
	}

	if resp.Pagination.TotalPages != nil && *resp.Pagination.TotalPages > int64(1) {
		p.Warn("Server responded with more than 1 page. Output will not show all existing organizations.\n")
	}
	return organizations, nil
}

func outputResult(p *print.Printer, outputFormat, projectLabel string, organizations []scf.OrganizationsListItem) error {
	return p.OutputResult(outputFormat, organizations, func() error {
		if len(organizations) == 0 {
			p.Outputf("No STACKIT Cloud Foundry organizations found for project %q\n", projectLabel)
			return nil
		}

		table := tables.NewTable()

		table.SetTitle("STACKIT Cloud Foundry organizations")
		table.SetHeader("ID", "NAME", "STATUS", "PLATFORM ID", "CREATED")
		table.AddSeparator()
		for _, organization := range organizations {
			table.AddRow(
				organization.Guid,
				utils.PtrString(organization.Name),
				organization.Status,
				organization.PlatformId,
				utils.ConvertTimePToDateTimeString(organization.CreatedAt),
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
