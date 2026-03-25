package create

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/iaas/client"
	iaasUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/iaas/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/iaas"
)

const (
	destinationTypeFlag  = "destination-type"
	destinationValueFlag = "destination-value"
	labelFlag            = "labels"
	networkAreaIdFlag    = "network-area-id"
	nextHopTypeFlag      = "nexthop-type"
	nextHopValueFlag     = "nexthop-value"
	organizationIdFlag   = "organization-id"
	routingTableIdFlag   = "routing-table-id"

	// Destination Type Constants
	destTypeCIDRv4 = "cidrv4"
	destTypeCIDRv6 = "cidrv6"

	// NextHop Type Constants
	nextHopTypeIPv4      = "ipv4"
	nextHopTypeIPv6      = "ipv6"
	nextHopTypeInternet  = "internet"
	nextHopTypeBlackhole = "blackhole"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	DestinationType  string
	DestinationValue *string
	Labels           *map[string]string
	NetworkAreaId    string
	NextHopType      string
	NextHopValue     *string
	OrganizationId   string
	RoutingTableId   string
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Creates a route in a routing-table",
		Long:  "Creates a route in a routing-table.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample("Create a route with CIDRv4 destination and IPv4 nexthop",
				`$ stackit network-area routing-table route create --routing-table-id xxx --organization-id yyy --network-area-id zzz --destination-type cidrv4 --destination-value <ipv4-cidr> --nexthop-type ipv4 --nexthop-value <ipv4-address>`),

			examples.NewExample("Create a route with CIDRv6 destination and IPv6 nexthop",
				`$ stackit network-area routing-table route create --routing-table-id xxx --organization-id yyy --network-area-id zzz --destination-type cidrv6 --destination-value <ipv6-cidr> --nexthop-type ipv6 --nexthop-value <ipv6-address>`),

			examples.NewExample("Create a route with CIDRv6 destination and Nexthop Internet",
				`$ stackit network-area routing-table route create --routing-table-id xxx --organization-id yyy --network-area-id zzz --destination-type cidrv6 --destination-value <ipv6-cidr> --nexthop-type internet`),
		),
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := context.Background()
			model, err := parseInput(params.Printer, cmd, nil)
			if err != nil {
				return err
			}

			apiClient, err := client.ConfigureClient(params.Printer, params.CliVersion)
			if err != nil {
				return err
			}

			routingTableLabel, err := iaasUtils.GetRoutingTableOfAreaName(ctx, apiClient, model.OrganizationId, model.NetworkAreaId, model.Region, model.RoutingTableId)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get routing-table name: %v", err)
				routingTableLabel = model.RoutingTableId
			} else if routingTableLabel == "" {
				routingTableLabel = model.RoutingTableId
			}

			prompt := fmt.Sprintf("Are you sure you want to create a route for routing-table with id %q?", routingTableLabel)
			err = params.Printer.PromptForConfirmation(prompt)
			if err != nil {
				return err
			}

			req, err := buildRequest(ctx, model, apiClient)
			if err != nil {
				return err
			}

			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("create route request failed: %w", err)
			}

			return outputResult(params.Printer, model.OutputFormat, resp.GetItems())
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.CIDRFlag(), destinationValueFlag, "Destination value")
	cmd.Flags().String(nextHopValueFlag, "", "NextHop value")
	cmd.Flags().Var(flags.UUIDFlag(), networkAreaIdFlag, "Network-Area ID")
	cmd.Flags().Var(flags.UUIDFlag(), organizationIdFlag, "Organization ID")
	cmd.Flags().Var(flags.UUIDFlag(), routingTableIdFlag, "Routing-Table ID")

	cmd.Flags().Var(
		flags.EnumFlag(true, "", destTypeCIDRv4, destTypeCIDRv6),
		destinationTypeFlag,
		"Destination type")

	cmd.Flags().Var(
		flags.EnumFlag(true, "", nextHopTypeIPv4, nextHopTypeIPv6, nextHopTypeInternet, nextHopTypeBlackhole),
		nextHopTypeFlag,
		"Next hop type")

	cmd.Flags().StringToString(labelFlag, nil, "Key=value labels")

	err := flags.MarkFlagsRequired(cmd, organizationIdFlag, networkAreaIdFlag, routingTableIdFlag, destinationTypeFlag, destinationValueFlag, nextHopTypeFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command, _ []string) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)

	model := &inputModel{
		GlobalFlagModel:  globalFlags,
		DestinationType:  flags.FlagToStringValue(p, cmd, destinationTypeFlag),
		DestinationValue: flags.FlagToStringPointer(p, cmd, destinationValueFlag),
		Labels:           flags.FlagToStringToStringPointer(p, cmd, labelFlag),
		NetworkAreaId:    flags.FlagToStringValue(p, cmd, networkAreaIdFlag),
		NextHopType:      flags.FlagToStringValue(p, cmd, nextHopTypeFlag),
		NextHopValue:     flags.FlagToStringPointer(p, cmd, nextHopValueFlag),
		OrganizationId:   flags.FlagToStringValue(p, cmd, organizationIdFlag),
		RoutingTableId:   flags.FlagToStringValue(p, cmd, routingTableIdFlag),
	}

	// Next Hop validation logic
	switch strings.ToLower(model.NextHopType) {
	case nextHopTypeInternet, nextHopTypeBlackhole:
		if model.NextHopValue != nil && *model.NextHopValue != "" {
			return nil, errors.New("--nexthop-value is not allowed when --nexthop-type is 'internet' or 'blackhole'")
		}
	case nextHopTypeIPv4, nextHopTypeIPv6:
		if model.NextHopValue == nil || *model.NextHopValue == "" {
			return nil, errors.New("--nexthop-value is required when --nexthop-type is 'ipv4' or 'ipv6'")
		}
	default:
		return nil, fmt.Errorf("invalid nexthop-type: %q", model.NextHopType)
	}

	p.DebugInputModel(model)
	return model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *iaas.APIClient) (iaas.ApiAddRoutesToRoutingTableRequest, error) {
	destination := buildDestination(model)
	nextHop := buildNextHop(model)

	if destination != nil && nextHop != nil {
		payload := iaas.AddRoutesToRoutingTablePayload{
			Items: &[]iaas.Route{
				{
					Destination: destination,
					Nexthop:     nextHop,
					Labels:      utils.ConvertStringMapToInterfaceMap(model.Labels),
				},
			},
		}

		return apiClient.AddRoutesToRoutingTable(
			ctx,
			model.OrganizationId,
			model.NetworkAreaId,
			model.Region,
			model.RoutingTableId,
		).AddRoutesToRoutingTablePayload(payload), nil
	}

	return nil, fmt.Errorf("invalid input")
}

func buildDestination(model *inputModel) *iaas.RouteDestination {
	if model.DestinationValue == nil {
		return nil
	}

	destinationType := strings.ToLower(model.DestinationType)
	switch destinationType {
	case destTypeCIDRv4:
		return &iaas.RouteDestination{
			DestinationCIDRv4: &iaas.DestinationCIDRv4{
				Type:  &model.DestinationType,
				Value: model.DestinationValue,
			},
		}
	case destTypeCIDRv6:
		return &iaas.RouteDestination{
			DestinationCIDRv6: &iaas.DestinationCIDRv6{
				Type:  &model.DestinationType,
				Value: model.DestinationValue,
			},
		}
	default:
		return nil
	}
}

func buildNextHop(model *inputModel) *iaas.RouteNexthop {
	nextHopType := strings.ToLower(model.NextHopType)
	switch nextHopType {
	case nextHopTypeIPv4:
		return &iaas.RouteNexthop{
			NexthopIPv4: &iaas.NexthopIPv4{
				Type:  &model.NextHopType,
				Value: model.NextHopValue,
			},
		}
	case nextHopTypeIPv6:
		return &iaas.RouteNexthop{
			NexthopIPv6: &iaas.NexthopIPv6{
				Type:  &model.NextHopType,
				Value: model.NextHopValue,
			},
		}
	case nextHopTypeInternet:
		return &iaas.RouteNexthop{
			NexthopInternet: &iaas.NexthopInternet{
				Type: &model.NextHopType,
			},
		}
	case nextHopTypeBlackhole:
		return &iaas.RouteNexthop{
			NexthopBlackhole: &iaas.NexthopBlackhole{
				Type: &model.NextHopType,
			},
		}
	default:
		return nil
	}
}

func outputResult(p *print.Printer, outputFormat string, routes []iaas.Route) error {
	if len(routes) == 0 {
		return fmt.Errorf("create routes response is empty")
	}

	return p.OutputResult(outputFormat, routes, func() error {
		for _, route := range routes {
			p.Outputf("Created route with ID %q\n", utils.PtrString(route.Id))
		}
		return nil
	})
}
