package list

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/iaas/client"
	rmClient "github.com/stackitcloud/stackit-cli/internal/pkg/services/resourcemanager/client"
	rmUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/resourcemanager/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/iaas"
)

const (
	labelSelectorFlag  = "label-selector"
	limitFlag          = "limit"
	networkAreaIdFlag  = "network-area-id"
	organizationIdFlag = "organization-id"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	LabelSelector  *string
	Limit          *int64
	NetworkAreaId  *string
	OrganizationId *string
}

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Lists all routing-tables",
		Long:  "Lists all routing-tables",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`List all routing-tables`,
				`$ stackit routing-table list --organization-id xxx --network-area-id yyy`,
			),
			examples.NewExample(
				`List all routing-tables with labels`,
				`$ stackit routing-table list --label-selector env=dev,env=rc --organization-id xxx --network-area-id yyy`,
			),
			examples.NewExample(
				`List all routing-tables with labels and set limit to 10`,
				`$ stackit routing-table list --label-selector env=dev,env=rc --limit 10 --organization-id xxx --network-area-id yyy`,
			),
		),
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := context.Background()
			model, err := parseInput(params.Printer, cmd, nil)
			if err != nil {
				return err
			}

			// Configure API client
			apiClient, err := client.ConfigureClient(params.Printer, params.CliVersion)
			if err != nil {
				return err
			}

			// Call API
			request := buildRequest(ctx, model, apiClient)

			response, err := request.Execute()
			if err != nil {
				return fmt.Errorf("list routing-tables: %w", err)
			}

			if items := response.Items; items == nil || len(*items) == 0 {
				var orgLabel string
				rmApiClient, err := rmClient.ConfigureClient(params.Printer, params.CliVersion)
				if err == nil {
					orgLabel, err = rmUtils.GetOrganizationName(ctx, rmApiClient, *model.OrganizationId)
					if err != nil {
						params.Printer.Debug(print.ErrorLevel, "get organization name: %v", err)
						orgLabel = *model.OrganizationId
					} else if orgLabel == "" {
						orgLabel = *model.OrganizationId
					}
				} else {
					params.Printer.Debug(print.ErrorLevel, "configure resource manager client: %v", err)
				}
				params.Printer.Info("No routing-tables found for organization %q\n", orgLabel)
				return nil
			}

			// Truncate output
			items := *response.Items
			if model.Limit != nil && len(items) > int(*model.Limit) {
				items = items[:*model.Limit]
			}

			return outputResult(params.Printer, model.OutputFormat, items)
		},
	}

	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Int64(limitFlag, 0, "Maximum number of entries to list")
	cmd.Flags().String(labelSelectorFlag, "", "Filter by label")
	cmd.Flags().Var(flags.UUIDFlag(), networkAreaIdFlag, "Network-Area ID")
	cmd.Flags().Var(flags.UUIDFlag(), organizationIdFlag, "Organization ID")

	err := flags.MarkFlagsRequired(cmd, organizationIdFlag, networkAreaIdFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command, _ []string) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)

	limit := flags.FlagToInt64Pointer(p, cmd, limitFlag)
	if limit != nil && *limit < 1 {
		return nil, &errors.FlagValidationError{
			Flag:    limitFlag,
			Details: "must be greater than 0",
		}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		LabelSelector:   flags.FlagToStringPointer(p, cmd, labelSelectorFlag),
		Limit:           limit,
		NetworkAreaId:   flags.FlagToStringPointer(p, cmd, networkAreaIdFlag),
		OrganizationId:  flags.FlagToStringPointer(p, cmd, organizationIdFlag),
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *iaas.APIClient) iaas.ApiListRoutingTablesOfAreaRequest {
	request := apiClient.ListRoutingTablesOfArea(ctx, *model.OrganizationId, *model.NetworkAreaId, model.Region)
	if model.LabelSelector != nil {
		request.LabelSelector(*model.LabelSelector)
	}

	return request
}
func outputResult(p *print.Printer, outputFormat string, items []iaas.RoutingTable) error {
	if len(items) == 0 {
		return fmt.Errorf("list routingtable response is empty")
	}

	return p.OutputResult(outputFormat, items, func() error {
		table := tables.NewTable()
		table.SetHeader("ID", "NAME", "DESCRIPTION", "CREATED_AT", "UPDATED_AT", "DEFAULT", "LABELS", "SYSTEM_ROUTES", "DYNAMIC_ROUTES")

		for _, item := range items {
			var labels []string
			for key, value := range *item.Labels {
				labels = append(labels, fmt.Sprintf("%s: %s", key, value))
			}

			createdAt := ""
			if item.CreatedAt != nil {
				createdAt = item.CreatedAt.Format(time.RFC3339)
			}

			updatedAt := ""
			if item.UpdatedAt != nil {
				updatedAt = item.UpdatedAt.Format(time.RFC3339)
			}

			table.AddRow(
				utils.PtrString(item.Id),
				utils.PtrString(item.Name),
				utils.PtrString(item.Description),
				createdAt,
				updatedAt,
				utils.PtrString(item.Default),
				strings.Join(labels, "\n"),
				utils.PtrString(item.SystemRoutes),
				utils.PtrString(item.DynamicRoutes),
			)
		}
		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}

		return nil
	})
}
