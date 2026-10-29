package create

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	iaas "github.com/stackitcloud/stackit-sdk-go/services/iaas/v2alpha1api"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/iaasalpha/client"
	iaasUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/iaasalpha/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

const (
	destinationValueFlag = "destination-value"
	labelFlag            = "labels"
	vpcIdFlag            = "vpc-id"
	nextHopValueFlag     = "nexthop-value"
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

var (
	destinationTypeFlag = flags.StringEnumFlag(
		"destination-type",
		[]string{destTypeCIDRv4, destTypeCIDRv6},
		"Destination type",
		flags.StringEnumIgnoreCase[string](),
	)
	nextHopTypeFlag = flags.StringEnumFlag(
		"nexthop-type",
		[]string{nextHopTypeIPv4, nextHopTypeIPv6, nextHopTypeInternet, nextHopTypeBlackhole},
		"Next hop type",
		flags.StringEnumIgnoreCase[string](),
	)
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	DestinationType  string
	DestinationValue string
	Labels           map[string]any
	VpcId            string
	NextHopType      string
	NextHopValue     *string
	RoutingTableId   string
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Creates a static-route in a routing-table",
		Long:  "Creates a static-route in a routing-table.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample("Create a static-route with CIDRv4 destination and IPv4 nexthop",
				`$ stackit alpha vpc routing-table static-route create --routing-table-id xxx --vpc-id yyy --destination-type cidrv4 --destination-value <ipv4-cidr> --nexthop-type ipv4 --nexthop-value <ipv4-address>`),

			examples.NewExample("Create a route with CIDRv6 destination and IPv6 nexthop",
				`$ stackit alpha vpc routing-table static-route create --routing-table-id xxx --vpc-id yyy --destination-type cidrv6 --destination-value <ipv6-cidr> --nexthop-type ipv6 --nexthop-value <ipv6-address>`),

			examples.NewExample("Create a route with CIDRv6 destination and nexthop Internet",
				`$ stackit alpha vpc routing-table static-route create --routing-table-id xxx --vpc-id yyy --destination-type cidrv6 --destination-value <ipv6-cidr> --nexthop-type internet`),
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

			routingTableLabel, err := iaasUtils.GetVPCRoutingTableName(ctx, apiClient.DefaultAPI, model.ProjectId, model.VpcId, model.Region, model.RoutingTableId)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get vpc routing-table name: %v", err)
				routingTableLabel = model.RoutingTableId
			} else if routingTableLabel == "" {
				routingTableLabel = model.RoutingTableId
			}

			prompt := fmt.Sprintf("Are you sure you want to create a static-route for vpc routing-table %q?", routingTableLabel)
			err = params.Printer.PromptForConfirmation(prompt)
			if err != nil {
				return err
			}

			req, err := buildRequest(ctx, model, apiClient.DefaultAPI)
			if err != nil {
				return err
			}

			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("create route request failed: %w", err)
			}

			return outputResult(params.Printer, model.OutputFormat, resp)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.CIDRFlag(), destinationValueFlag, "Destination value")
	cmd.Flags().String(nextHopValueFlag, "", "NextHop value")
	cmd.Flags().Var(flags.UUIDFlag(), vpcIdFlag, "VPC ID")
	cmd.Flags().Var(flags.UUIDFlag(), routingTableIdFlag, "Routing-Table ID")
	destinationTypeFlag.Register(cmd.Flags())
	nextHopTypeFlag.Register(cmd.Flags())

	cmd.Flags().StringToString(labelFlag, nil, "Labels are key-value string pairs which can be attached to a static route. E.g. '--labels key1=value1,key2=value2,...'")

	err := flags.MarkFlagsRequired(cmd, vpcIdFlag, routingTableIdFlag, destinationTypeFlag.Name(), destinationValueFlag, nextHopTypeFlag.Name())
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command, _ []string) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)

	model := &inputModel{
		GlobalFlagModel:  globalFlags,
		DestinationType:  destinationTypeFlag.Get(),
		DestinationValue: flags.FlagToStringValue(p, cmd, destinationValueFlag),
		Labels:           flags.FlagToStringToAny(p, cmd, labelFlag),
		VpcId:            flags.FlagToStringValue(p, cmd, vpcIdFlag),
		NextHopType:      nextHopTypeFlag.Get(),
		NextHopValue:     flags.FlagToStringPointer(p, cmd, nextHopValueFlag),
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

func buildRequest(ctx context.Context, model *inputModel, apiClient iaas.DefaultAPI) (iaas.ApiAddVPCStaticRouteRequest, error) {
	destination := buildDestination(model)
	nextHop := buildNextHop(model)

	if destination != nil && nextHop != nil {
		payload := iaas.AddVPCStaticRoutePayload{
			Destination: *destination,
			Nexthop:     *nextHop,
			Labels:      model.Labels,
		}

		return apiClient.AddVPCStaticRoute(
			ctx,
			model.ProjectId,
			model.VpcId,
			model.Region,
			model.RoutingTableId,
		).AddVPCStaticRoutePayload(payload), nil
	}

	return iaas.ApiAddVPCStaticRouteRequest{}, fmt.Errorf("invalid input")
}

func buildDestination(model *inputModel) *iaas.AddVPCStaticRoutePayloadDestination {
	destinationType := strings.ToLower(model.DestinationType)
	switch destinationType {
	case destTypeCIDRv4:
		return &iaas.AddVPCStaticRoutePayloadDestination{
			DestinationCIDRv4: &iaas.DestinationCIDRv4{
				Type:  model.DestinationType,
				Value: model.DestinationValue,
			},
		}
	case destTypeCIDRv6:
		return &iaas.AddVPCStaticRoutePayloadDestination{
			DestinationCIDRv6: &iaas.DestinationCIDRv6{
				Type:  model.DestinationType,
				Value: model.DestinationValue,
			},
		}
	default:
		return nil
	}
}

func buildNextHop(model *inputModel) *iaas.AddVPCStaticRoutePayloadNexthop {
	nextHopType := strings.ToLower(model.NextHopType)
	switch nextHopType {
	case nextHopTypeIPv4:
		return &iaas.AddVPCStaticRoutePayloadNexthop{
			NexthopIPv4: &iaas.NexthopIPv4{
				Type:  model.NextHopType,
				Value: *model.NextHopValue,
			},
		}
	case nextHopTypeIPv6:
		return &iaas.AddVPCStaticRoutePayloadNexthop{
			NexthopIPv6: &iaas.NexthopIPv6{
				Type:  model.NextHopType,
				Value: *model.NextHopValue,
			},
		}
	case nextHopTypeInternet:
		return &iaas.AddVPCStaticRoutePayloadNexthop{
			NexthopInternet: &iaas.NexthopInternet{
				Type: model.NextHopType,
			},
		}
	case nextHopTypeBlackhole:
		return &iaas.AddVPCStaticRoutePayloadNexthop{
			NexthopBlackhole: &iaas.NexthopBlackhole{
				Type: model.NextHopType,
			},
		}
	default:
		return nil
	}
}

func outputResult(p *print.Printer, outputFormat string, route *iaas.Route) error {
	if route == nil {
		return fmt.Errorf("create route response is empty")
	}

	return p.OutputResult(outputFormat, route, func() error {
		p.Outputf("Created route with ID %q\n", utils.PtrString(route.Id))
		return nil
	})
}
