package list

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/iaas/client"
	routeUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/routingtable/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/iaas"
)

const (
	labelSelectorFlag  = "label-selector"
	limitFlag          = "limit"
	networkAreaIdFlag  = "network-area-id"
	organizationIdFlag = "organization-id"
	routingTableIdFlag = "routing-table-id"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	LabelSelector  *string
	Limit          *int64
	NetworkAreaId  string
	OrganizationId string
	RoutingTableId string
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Lists all routes within a routing-table",
		Long:  "Lists all routes within a routing-table",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`List all routes within a routing-table`,
				`$ stackit routing-table route list --routing-table-id xxx --organization-id yyy --network-area-id zzz`,
			),
			examples.NewExample(
				`List all routes within a routing-table with labels`,
				`$ stackit routing-table list --routing-table-id xxx --organization-id yyy --network-area-id zzz --label-selector env=dev,env=rc`,
			),
			examples.NewExample(
				`List all routes within a routing-tables with labels and limit to 10`,
				`$ stackit routing-table list --routing-table-id xxx --organization-id yyy --network-area-id zzz --label-selector env=dev,env=rc --limit 10`,
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
			request := apiClient.ListRoutesOfRoutingTable(
				ctx,
				model.OrganizationId,
				model.NetworkAreaId,
				model.Region,
				model.RoutingTableId,
			)

			if model.LabelSelector != nil {
				request.LabelSelector(*model.LabelSelector)
			}

			response, err := request.Execute()
			if err != nil {
				return fmt.Errorf("list routes: %w", err)
			}

			if items := response.Items; items == nil {
				params.Printer.Outputf("No routes found for routing-table %q\n", model.RoutingTableId)
				return nil
			}

			// Truncate output
			items := response.GetItems()
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
	cmd.Flags().Var(flags.UUIDFlag(), routingTableIdFlag, "Routing-Table ID")

	err := flags.MarkFlagsRequired(cmd, organizationIdFlag, networkAreaIdFlag, routingTableIdFlag)
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
		NetworkAreaId:   flags.FlagToStringValue(p, cmd, networkAreaIdFlag),
		OrganizationId:  flags.FlagToStringValue(p, cmd, organizationIdFlag),
		RoutingTableId:  flags.FlagToStringValue(p, cmd, routingTableIdFlag),
	}

	p.DebugInputModel(model)
	return &model, nil
}

func outputResult(p *print.Printer, outputFormat string, routes []iaas.Route) error {
	if routes == nil {
		return fmt.Errorf("list routes routes are nil")
	}

	return p.OutputResult(outputFormat, routes, func() error {
		table := tables.NewTable()
		table.SetHeader("ID", "DESTINATION TYPE", "DESTINATION VALUE", "NEXTHOP TYPE", "NEXTHOP VALUE", "LABELS", "CREATED AT", "UPDATED AT")
		for _, route := range routes {
			routeDetails := routeUtils.ExtractRouteDetails(route)
			table.AddRow(
				utils.PtrString(route.Id),
				routeDetails.DestType,
				routeDetails.DestValue,
				routeDetails.HopType,
				routeDetails.HopValue,
				routeDetails.Labels,
				routeDetails.CreatedAt,
				routeDetails.UpdatedAt,
			)
		}
		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}
		return nil
	})
}
