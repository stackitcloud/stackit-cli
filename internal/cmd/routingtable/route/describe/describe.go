package describe

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
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
	networkAreaIdFlag  = "network-area-id"
	organizationIdFlag = "organization-id"
	routeIdArg         = "ROUTE_ID"
	routingTableIdFlag = "routing-table-id"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	NetworkAreaId  string
	OrganizationId string
	RouteID        string
	RoutingTableId string
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("describe %s", routeIdArg),
		Short: "Describes a route within a routing-table",
		Long:  "Describes a route within a routing-table",
		Args:  args.SingleArg(routeIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Describe a route within a routing-table`,
				`$ stackit routing-table route describe xxx --routing-table-id xxx --organization-id yyy --network-area-id zzz`,
			),
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
			request := apiClient.GetRouteOfRoutingTable(
				ctx,
				model.OrganizationId,
				model.NetworkAreaId,
				model.Region,
				model.RoutingTableId,
				model.RouteID,
			)

			response, err := request.Execute()
			if err != nil {
				return fmt.Errorf("describe route: %w", err)
			}

			return outputResult(params.Printer, model.OutputFormat, response)
		},
	}

	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), networkAreaIdFlag, "Network-Area ID")
	cmd.Flags().Var(flags.UUIDFlag(), organizationIdFlag, "Organization ID")
	cmd.Flags().Var(flags.UUIDFlag(), routingTableIdFlag, "Routing-Table ID")

	err := flags.MarkFlagsRequired(cmd, organizationIdFlag, networkAreaIdFlag, routingTableIdFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	routeId := inputArgs[0]

	model := inputModel{
		GlobalFlagModel: globalFlags,
		NetworkAreaId:   flags.FlagToStringValue(p, cmd, networkAreaIdFlag),
		OrganizationId:  flags.FlagToStringValue(p, cmd, organizationIdFlag),
		RouteID:         routeId,
		RoutingTableId:  flags.FlagToStringValue(p, cmd, routingTableIdFlag),
	}

	p.DebugInputModel(model)
	return &model, nil
}

func outputResult(p *print.Printer, outputFormat string, route *iaas.Route) error {
	if route == nil {
		return fmt.Errorf("describe route response is empty")
	}

	return p.OutputResult(outputFormat, route, func() error {
		routeDetails := routeUtils.ExtractRouteDetails(*route)

		table := tables.NewTable()
		table.SetHeader("ID", "DESTINATION TYPE", "DESTINATION VALUE", "NEXTHOP TYPE", "NEXTHOP VALUE", "LABELS", "CREATED AT", "UPDATED AT")
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

		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}

		return nil
	})
}
