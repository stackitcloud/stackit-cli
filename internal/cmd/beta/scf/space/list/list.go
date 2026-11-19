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
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/scf/client"
	scfUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/scf/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

const (
	organizationIdFlag = "organization-id"
	limitFlag          = "limit"

	// maxPageSize is the maximum number of items the API returns per page. 5000 is the maximum the API supports.
	maxPageSize = 5000
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	OrganizationId string
	Limit          *int64
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Lists all STACKIT Cloud Foundry spaces of an organization",
		Long:  fmt.Sprintf("Lists all STACKIT Cloud Foundry spaces of an organization. The number of fetched spaces is limited to %d.", maxPageSize),
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`List all STACKIT Cloud Foundry spaces of the organization with ID "xxx"`,
				"$ stackit beta scf space list --organization-id xxx"),
			examples.NewExample(
				`List up to 10 STACKIT Cloud Foundry spaces of the organization with ID "xxx"`,
				"$ stackit beta scf space list --organization-id xxx --limit 10"),
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
				return fmt.Errorf("fetch STACKIT Cloud Foundry spaces: %w", err)
			}

			organizationLabel := scfUtils.GetOrganizationLabel(ctx, params.Printer, apiClient.DefaultAPI, model.ProjectId, model.Region, model.OrganizationId)

			return outputResult(params.Printer, model.OutputFormat, model.Limit, organizationLabel, resp)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), organizationIdFlag, "The UUID of the organization")
	cmd.Flags().Int64(limitFlag, 0, "Maximum number of entries to list")

	err := flags.MarkFlagsRequired(cmd, organizationIdFlag)
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
		OrganizationId:  flags.FlagToStringValue(p, cmd, organizationIdFlag),
		Limit:           limit,
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient scf.DefaultAPI, page, perPage int64) scf.ApiListSpacesRequest {
	pageSize := perPage
	if model.Limit != nil && *model.Limit < pageSize {
		pageSize = *model.Limit
	}
	return apiClient.ListSpaces(ctx, model.ProjectId, model.Region, model.OrganizationId).Page(page).PerPage(pageSize)
}

func outputResult(p *print.Printer, outputFormat string, limit *int64, organizationLabel string, resp *scf.SpacesList) error {
	return p.OutputResult(outputFormat, resp, func() error {
		if resp == nil || len(resp.Resources) == 0 {
			p.Outputf("No STACKIT Cloud Foundry spaces found for organization %q\n", organizationLabel)
			return nil
		}
		limitReached := limit != nil && int64(len(resp.Resources)) >= *limit
		if resp.Pagination.TotalPages != nil && *resp.Pagination.TotalPages > int64(1) && !limitReached {
			p.Warn("Server responded with more than 1 page. Output will not show all existing spaces.\n")
		}

		table := tables.NewTable()

		table.SetTitle("STACKIT Cloud Foundry spaces")
		table.SetHeader("ID", "NAME", "PLATFORM ID", "CREATED")
		table.AddSeparator()
		for _, space := range resp.Resources {
			table.AddRow(
				space.Guid,
				space.Name,
				space.PlatformId,
				utils.ConvertTimePToDateTimeString(&space.CreatedAt),
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
