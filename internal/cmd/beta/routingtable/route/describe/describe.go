package describe

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/goccy/go-yaml"
	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/iaas/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/iaasalpha"
)

const (
	organizationIdFlag = "organization-id"
	networkAreaIdFlag  = "network-area-id"
	routingTableIdFlag = "routing-table-id"
	routeIdArg         = "ROUTE_ID_ARG"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	OrganizationId *string
	NetworkAreaId  *string
	RoutingTableId *string
	RouteID        *string
}

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("describe %s", routeIdArg),
		Short: "Describe a route within a routing-table",
		Long:  "Describe a route within a routing-table",
		Args:  args.SingleArg(routeIdArg, nil),
		Example: examples.Build(
			examples.NewExample(
				`Describe a route within a routing-table`,
				`$ stackit beta routing-table route describe xxxx-xxxx-xxxx-xxxx --routing-table-id xxx --organization-id yyy --network-area-id zzz`,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			model, err := parseInput(params.Printer, cmd, args)
			if err != nil {
				return err
			}

			// Configure API client
			apiClient, err := client.ConfigureAlphaClient(params.Printer, params.CliVersion)
			if err != nil {
				return err
			}

			// Call API
			request := apiClient.GetRouteOfRoutingTable(
				ctx,
				*model.OrganizationId,
				*model.NetworkAreaId,
				model.Region,
				*model.RoutingTableId,
				*model.RouteID,
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
	cmd.Flags().Var(flags.UUIDFlag(), organizationIdFlag, "Organization ID")
	cmd.Flags().Var(flags.UUIDFlag(), networkAreaIdFlag, "Network-Area ID")
	cmd.Flags().Var(flags.UUIDFlag(), routingTableIdFlag, "Routing-Table ID")

	err := flags.MarkFlagsRequired(cmd, organizationIdFlag, networkAreaIdFlag, routingTableIdFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command, args []string) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)

	if len(args) == 0 {
		return nil, fmt.Errorf("at least one argument is required")
	}
	routeId := args[0]

	model := inputModel{
		GlobalFlagModel: globalFlags,
		NetworkAreaId:   flags.FlagToStringPointer(p, cmd, networkAreaIdFlag),
		OrganizationId:  flags.FlagToStringPointer(p, cmd, organizationIdFlag),
		RoutingTableId:  flags.FlagToStringPointer(p, cmd, routingTableIdFlag),
		RouteID:         &routeId,
	}

	p.DebugInputModel(model)
	return &model, nil
}

func outputResult(p *print.Printer, outputFormat string, routingTable *iaasalpha.Route) error {
	switch outputFormat {
	case print.JSONOutputFormat:
		details, err := json.MarshalIndent(routingTable, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal route describe: %w", err)
		}
		p.Outputln(string(details))

		return nil
	case print.YAMLOutputFormat:
		details, err := yaml.MarshalWithOptions(routingTable, yaml.IndentSequence(true), yaml.UseJSONMarshaler())
		if err != nil {
			return fmt.Errorf("marshal route describe: %w", err)
		}
		p.Outputln(string(details))

		return nil
	default:
		var labels []string
		for key, value := range *routingTable.Labels {
			labels = append(labels, fmt.Sprintf("%s: %s", key, value))
		}

		destinationType := ""
		destinationValue := ""
		if dest := routingTable.Destination.DestinationCIDRv4; dest != nil {
			if dest.Type != nil {
				destinationType = *dest.Type
			}
			if dest.Value != nil {
				destinationValue = *dest.Value
			}
		}
		if dest := routingTable.Destination.DestinationCIDRv6; dest != nil {
			if dest.Type != nil {
				destinationType = *dest.Type
			}
			if dest.Value != nil {
				destinationValue = *dest.Value
			}
		}

		nextHopType := ""
		nextHopValue := ""
		if nextHop := routingTable.Destination.DestinationCIDRv4; nextHop != nil {
			if nextHop.Type != nil {
				nextHopType = *nextHop.Type
			}
			if nextHop.Value != nil {
				nextHopValue = *nextHop.Value
			}
		}
		if nextHop := routingTable.Destination.DestinationCIDRv6; nextHop != nil {
			if nextHop.Type != nil {
				nextHopType = *nextHop.Type
			}
			if nextHop.Value != nil {
				nextHopValue = *nextHop.Value
			}
		}

		table := tables.NewTable()
		table.SetHeader("ID", "CREATED_AT", "UPDATED_AT", "DESTINATION TYPE", "DESTINATION VALUE", "NEXTHOP TYPE", "NEXTHOP VALUE", "LABELS")
		table.AddRow(
			utils.PtrString(routingTable.Id),
			routingTable.CreatedAt.String(),
			routingTable.UpdatedAt.String(),
			destinationType,
			destinationValue,
			nextHopType,
			nextHopValue,
			strings.Join(labels, "\n"),
		)

		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}

		return nil
	}
}
