package create

import (
	"context"
	"encoding/json"
	"errors"
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
	routeUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/routing-table/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/iaasalpha"
)

const (
	organizationIdFlag   = "organization-id"
	networkAreaIdFlag    = "network-area-id"
	routingTableIdFlag   = "routing-table-id"
	destinationTypeFlag  = "destination-type"
	destinationValueFlag = "destination-value"
	nextHopTypeFlag      = "nexthop-type"
	nextHopValueFlag     = "nexthop-value"
	labelFlag            = "labels"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	OrganizationId   *string
	NetworkAreaId    *string
	RoutingTableId   *string
	DestinationType  *string
	DestinationValue *string
	NextHopType      *string
	NextHopValue     *string
	Labels           *map[string]string
}

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Creates a route in a routing-table",
		Long:  "Creates a route in a routing-table.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample("Create a route with CIDRv4 destination and IPv4 nexthop",
				`stackit beta routing-tables route create  \ 
--routing-table-id xxx --organization-id yyy --network-area-id zzz \
--destination-type cidrv4 --destination-value <ipv4-cidr> \
--nexthop-type ipv4 --nexthop-value <ipv4-address>`),

			examples.NewExample("Create a route with CIDRv6 destination and IPv6 nexthop",
				`stackit beta routing-tables route create \
--routing-table-id xxx --organization-id yyy --network-area-id zzz \
--destination-type cidrv6 --destination-value <ipv6-cidr> \
--nexthop-type ipv6 --nexthop-value <ipv6-address>`),

			examples.NewExample("Create a route with CIDRv6 destination and Nexthop Internet",
				`stackit beta routing-tables route create \
--routing-table-id xxx --organization-id yyy --network-area-id zzz \
--destination-type cidrv6 --destination-value <ipv6-cidr> \
--nexthop-type internet`),
		),
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := context.Background()
			model, err := parseInput(params.Printer, cmd)
			if err != nil {
				return err
			}

			apiClient, err := client.ConfigureAlphaClient(params.Printer, params.CliVersion)
			if err != nil {
				return err
			}

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to create a route for routing-table with id %q?", *model.RoutingTableId)
				if err := params.Printer.PromptForConfirmation(prompt); err != nil {
					return err
				}
			}

			req, err := buildRequest(ctx, model, apiClient)
			if err != nil {
				return err
			}

			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("create route request failed: %w", err)
			}

			return outputResult(params.Printer, model.OutputFormat, *resp.Items)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), organizationIdFlag, "Organization ID")
	cmd.Flags().Var(flags.UUIDFlag(), networkAreaIdFlag, "Network-Area ID")
	cmd.Flags().Var(flags.UUIDFlag(), routingTableIdFlag, "Routing-Table ID")
	cmd.Flags().Var(flags.CIDRFlag(), destinationValueFlag, "Destination value")
	cmd.Flags().String(nextHopValueFlag, "", "NextHop value")

	cmd.Flags().Var(
		flags.EnumFlag(true, "", "cidrv4", "cidrv6"),
		destinationTypeFlag,
		"Destination type")

	cmd.Flags().Var(
		flags.EnumFlag(true, "", "ipv4", "ipv6", "internet", "blackhole"),
		nextHopTypeFlag,
		"Next hop type")

	cmd.Flags().StringToString(labelFlag, nil, "Key=value labels")

	err := flags.MarkFlagsRequired(cmd, organizationIdFlag, networkAreaIdFlag, routingTableIdFlag, destinationTypeFlag, destinationValueFlag, nextHopTypeFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)

	model := &inputModel{
		GlobalFlagModel:  globalFlags,
		OrganizationId:   flags.FlagToStringPointer(p, cmd, organizationIdFlag),
		NetworkAreaId:    flags.FlagToStringPointer(p, cmd, networkAreaIdFlag),
		RoutingTableId:   flags.FlagToStringPointer(p, cmd, routingTableIdFlag),
		DestinationType:  flags.FlagToStringPointer(p, cmd, destinationTypeFlag),
		DestinationValue: flags.FlagToStringPointer(p, cmd, destinationValueFlag),
		NextHopType:      flags.FlagToStringPointer(p, cmd, nextHopTypeFlag),
		NextHopValue:     flags.FlagToStringPointer(p, cmd, nextHopValueFlag),
		Labels:           flags.FlagToStringToStringPointer(p, cmd, labelFlag),
	}

	// Next Hop validation logic
	switch strings.ToLower(*model.NextHopType) {
	case "internet", "blackhole":
		if model.NextHopValue != nil && *model.NextHopValue != "" {
			return nil, errors.New("--nexthop-value is not allowed when --nexthop-type is 'internet' or 'blackhole'")
		}
	case "ipv4", "ipv6":
		if model.NextHopValue == nil || *model.NextHopValue == "" {
			return nil, errors.New("--nexthop-value is required when --nexthop-type is 'ipv4' or 'ipv6'")
		}
	default:
		return nil, fmt.Errorf("invalid nexthop-type: %q", *model.NextHopType)
	}

	p.DebugInputModel(model)
	return model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *iaasalpha.APIClient) (iaasalpha.ApiAddRoutesToRoutingTableRequest, error) {
	destination := buildDestination(model)
	nextHop := buildNextHop(model)

	if destination != nil && nextHop != nil {
		payload := iaasalpha.AddRoutesToRoutingTablePayload{
			Items: &[]iaasalpha.Route{
				{
					Destination: destination,
					Nexthop:     nextHop,
					Labels:      utils.ConvertStringMapToInterfaceMap(model.Labels),
				},
			},
		}

		return apiClient.AddRoutesToRoutingTable(
			ctx,
			*model.OrganizationId,
			*model.NetworkAreaId,
			model.Region,
			*model.RoutingTableId,
		).AddRoutesToRoutingTablePayload(payload), nil
	}

	return nil, fmt.Errorf("invalid input")
}

func buildDestination(model *inputModel) *iaasalpha.RouteDestination {
	if model.DestinationValue == nil {
		return nil
	}

	destinationType := strings.ToLower(*model.DestinationType)
	switch destinationType {
	case "cidrv4":
		return &iaasalpha.RouteDestination{
			DestinationCIDRv4: &iaasalpha.DestinationCIDRv4{
				Type:  model.DestinationType,
				Value: model.DestinationValue,
			},
		}
	case "cidrv6":
		return &iaasalpha.RouteDestination{
			DestinationCIDRv6: &iaasalpha.DestinationCIDRv6{
				Type:  model.DestinationType,
				Value: model.DestinationValue,
			},
		}
	default:
		return nil
	}
}

func buildNextHop(model *inputModel) *iaasalpha.RouteNexthop {
	nextHopType := strings.ToLower(*model.NextHopType)
	switch nextHopType {
	case "ipv4":
		return &iaasalpha.RouteNexthop{
			NexthopIPv4: &iaasalpha.NexthopIPv4{
				Type:  model.NextHopType,
				Value: model.NextHopValue,
			},
		}
	case "ipv6":
		return &iaasalpha.RouteNexthop{
			NexthopIPv6: &iaasalpha.NexthopIPv6{
				Type:  model.NextHopType,
				Value: model.NextHopValue,
			},
		}
	case "internet":
		return &iaasalpha.RouteNexthop{
			NexthopInternet: &iaasalpha.NexthopInternet{
				Type: model.NextHopType,
			},
		}
	case "blackhole":
		return &iaasalpha.RouteNexthop{
			NexthopBlackhole: &iaasalpha.NexthopBlackhole{
				Type: model.NextHopType,
			},
		}
	default:
		return nil
	}
}

func outputResult(p *print.Printer, outputFormat string, items []iaasalpha.Route) error {
	if len(items) == 0 {
		return fmt.Errorf("create routes response is empty")
	}

	switch outputFormat {
	case print.JSONOutputFormat:
		data, err := json.MarshalIndent(items, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal routes: %w", err)
		}
		p.Outputln(string(data))
	case print.YAMLOutputFormat:
		data, err := yaml.MarshalWithOptions(items, yaml.IndentSequence(true), yaml.UseJSONMarshaler())
		if err != nil {
			return fmt.Errorf("marshal routes: %w", err)
		}
		p.Outputln(string(data))
	default:
		table := tables.NewTable()
		table.SetHeader("ID", "DEST. TYPE", "DEST. VALUE", "NEXTHOP TYPE", "NEXTHOP VALUE", "LABELS", "CREATED", "UPDATED")
		for _, item := range items {
			destType, destValue, hopType, hopValue, labels := routeUtils.ExtractRouteDetails(item)

			table.AddRow(
				utils.PtrString(item.Id),
				destType,
				destValue,
				hopType,
				hopValue,
				labels,
				item.CreatedAt.String(),
				item.UpdatedAt.String(),
			)
		}
		return table.Display(p)
	}
	return nil
}
