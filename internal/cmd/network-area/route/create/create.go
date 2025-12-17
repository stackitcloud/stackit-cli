package create

import (
	"context"
	"fmt"
	"net"
	"os"

	"github.com/stackitcloud/stackit-cli/internal/pkg/types"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/iaas/client"
	iaasUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/iaas/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/iaas"

	"github.com/spf13/cobra"
)

const (
	organizationIdFlag = "organization-id"
	networkAreaIdFlag  = "network-area-id"
	// Deprecated: prefixFlag is deprecated and will be removed after April 2026. Use instead destinationFlag
	prefixFlag      = "prefix"
	destinationFlag = "destination"
	// Deprecated: nexthopFlag is deprecated and will be removed after April 2026. Use instead nexthopIPv4Flag or nexthopIPv6Flag
	nexthopFlag          = "next-hop"
	nexthopIPv4Flag      = "next-hop-ipv4"
	nexthopIPv6Flag      = "next-hop-ipv6"
	nexthopBlackholeFlag = "nexthop-blackhole"
	nexthopInternetFlag  = "nexthop-internet"
	labelFlag            = "labels"
)

const (
	destinationCIDRv4Type = "cidrv4"
	destinationCIDRv6Type = "cidrv6"

	nexthopBlackholeType = "blackhole"
	nexthopInternetType  = "internet"
	nexthopIPv4Type      = "ipv4"
	nexthopIPv6Type      = "ipv6"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	OrganizationId   *string
	NetworkAreaId    *string
	DestinationV4    *string
	DestinationV6    *string
	NexthopV4        *string
	NexthopV6        *string
	NexthopBlackhole *bool
	NexthopInternet  *bool
	Labels           *map[string]string
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Creates a static route in a STACKIT Network Area (SNA)",
		Long: fmt.Sprintf("%s\n%s\n",
			"Creates a static route in a STACKIT Network Area (SNA).",
			"This command is currently asynchonous only due to limitations in the waiting functionality of the SDK. This will be updated in a future release.",
		),
		Args: args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Create a static route with destination "1.1.1.0/24" and next hop "1.1.1.1" in a STACKIT Network Area with ID "xxx" in organization with ID "yyy"`,
				"$ stackit network-area route create --organization-id yyy --network-area-id xxx --destination 1.1.1.0/24 --next-hop 1.1.1.1",
			),
			examples.NewExample(
				`Create a static route with labels "key:value" and "foo:bar" with destination "1.1.1.0/24" and next hop "1.1.1.1" in a STACKIT Network Area with ID "xxx" in organization with ID "yyy"`,
				"$ stackit network-area route create --labels key=value,foo=bar --organization-id yyy --network-area-id xxx --destination 1.1.1.0/24 --next-hop 1.1.1.1",
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

			// Get network area label
			networkAreaLabel, err := iaasUtils.GetNetworkAreaName(ctx, apiClient, *model.OrganizationId, *model.NetworkAreaId)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get network area name: %v", err)
				networkAreaLabel = *model.NetworkAreaId
			}

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to create a static route for STACKIT Network Area (SNA) %q?", networkAreaLabel)
				err = params.Printer.PromptForConfirmation(prompt)
				if err != nil {
					return err
				}
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("create static route: %w", err)
			}

			if resp.Items == nil || len(*resp.Items) == 0 {
				return fmt.Errorf("empty response from API")
			}

			var destination string
			var nexthop string
			if model.DestinationV4 != nil {
				destination = *model.DestinationV4
			} else if model.DestinationV6 != nil {
				destination = *model.DestinationV6
			}

			if model.NexthopV4 != nil {
				nexthop = *model.NexthopV4
			} else if model.NexthopV6 != nil {
				nexthop = *model.NexthopV6
			} else if model.NexthopBlackhole != nil {
				// For nexthopBlackhole the type is assigned to nexthop, because it doesn't have any value
				nexthop = nexthopBlackholeType
			} else if model.NexthopInternet != nil {
				// For nexthopInternet the type is assigned to nexthop, because it doesn't have any value
				nexthop = nexthopInternetType
			}

			route, err := iaasUtils.GetRouteFromAPIResponse(destination, nexthop, resp.Items)
			if err != nil {
				return err
			}

			return outputResult(params.Printer, model.OutputFormat, networkAreaLabel, route)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), organizationIdFlag, "Organization ID")
	cmd.Flags().Var(flags.UUIDFlag(), networkAreaIdFlag, "STACKIT Network Area ID")
	cmd.Flags().Var(flags.CIDRFlag(), prefixFlag, "Static route prefix")
	cmd.Flags().Var(flags.CIDRFlag(), destinationFlag, "Destination route. Must be a valid IPv4 or IPv6 CIDR")

	cmd.Flags().StringToString(labelFlag, nil, "Labels are key-value string pairs which can be attached to a route. A label can be provided with the format key=value and the flag can be used multiple times to provide a list of labels")
	cmd.Flags().String(nexthopFlag, "", "Next hop IP address. Must be a valid IPv4")
	cmd.Flags().String(nexthopIPv4Flag, "", "Next hop IPv4 address")
	cmd.Flags().String(nexthopIPv6Flag, "", "Next hop IPv6 address")
	cmd.Flags().Bool(nexthopBlackholeFlag, false, "Sets next hop to black hole")
	cmd.Flags().Bool(nexthopInternetFlag, false, "Sets next hop to internet")

	cobra.CheckErr(cmd.Flags().MarkDeprecated(nexthopFlag, fmt.Sprintf("The flag %q is deprecated and will be removed after April 2026. Use instead %q to configure a IPv4 next hop.", nexthopFlag, nexthopBlackholeFlag)))
	cobra.CheckErr(cmd.Flags().MarkDeprecated(prefixFlag, fmt.Sprintf("The flag %q is deprecated and will be removed after April 2026. Use instead %q to configure a destination.", prefixFlag, destinationFlag)))
	// Set the output for deprecation warnings to stderr
	cmd.Flags().SetOutput(os.Stderr)

	destinationFlags := []string{prefixFlag, destinationFlag}
	nexthopFlags := []string{nexthopFlag, nexthopIPv4Flag, nexthopIPv6Flag, nexthopBlackholeFlag, nexthopInternetFlag}
	cmd.MarkFlagsMutuallyExclusive(destinationFlags...)
	cmd.MarkFlagsMutuallyExclusive(nexthopFlags...)

	cmd.MarkFlagsOneRequired(destinationFlags...)
	cmd.MarkFlagsOneRequired(nexthopFlags...)
	err := flags.MarkFlagsRequired(cmd, organizationIdFlag, networkAreaIdFlag)
	cobra.CheckErr(err)
}

func parseDestination(input string) (destinationV4, destinationV6 *string, err error) {
	ip, _, err := net.ParseCIDR(input)
	if err != nil {
		return nil, nil, fmt.Errorf("parse CIDR: %w", err)
	}
	if ip.To4() != nil { // CIDR is IPv4
		destinationV4 = utils.Ptr(input)
		return destinationV4, nil, nil
	}
	// CIDR is IPv6
	destinationV6 = utils.Ptr(input)
	return nil, destinationV6, nil
}

func parseInput(p *print.Printer, cmd *cobra.Command, _ []string) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)

	var destinationV4, destinationV6 *string
	if destination := flags.FlagToStringPointer(p, cmd, destinationFlag); destination != nil {
		var err error
		destinationV4, destinationV6, err = parseDestination(*destination)
		if err != nil {
			return nil, err
		}
	}
	if prefix := flags.FlagToStringPointer(p, cmd, prefixFlag); prefix != nil {
		var err error
		destinationV4, destinationV6, err = parseDestination(*prefix)
		if err != nil {
			return nil, err
		}
	}

	nexthopIPv4 := flags.FlagToStringPointer(p, cmd, nexthopIPv4Flag)
	nexthopIPv6 := flags.FlagToStringPointer(p, cmd, nexthopIPv6Flag)
	nexthopInternet := flags.FlagToBoolPointer(p, cmd, nexthopInternetFlag)
	nexthopBlackhole := flags.FlagToBoolPointer(p, cmd, nexthopBlackholeFlag)
	if nexthop := flags.FlagToStringPointer(p, cmd, nexthopFlag); nexthop != nil {
		nexthopIPv4 = nexthop
	}

	model := inputModel{
		GlobalFlagModel:  globalFlags,
		OrganizationId:   flags.FlagToStringPointer(p, cmd, organizationIdFlag),
		NetworkAreaId:    flags.FlagToStringPointer(p, cmd, networkAreaIdFlag),
		DestinationV4:    destinationV4,
		DestinationV6:    destinationV6,
		NexthopV4:        nexthopIPv4,
		NexthopV6:        nexthopIPv6,
		NexthopBlackhole: nexthopBlackhole,
		NexthopInternet:  nexthopInternet,
		Labels:           flags.FlagToStringToStringPointer(p, cmd, labelFlag),
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *iaas.APIClient) iaas.ApiCreateNetworkAreaRouteRequest {
	req := apiClient.CreateNetworkAreaRoute(ctx, *model.OrganizationId, *model.NetworkAreaId, model.Region)

	var destinationV4 *iaas.DestinationCIDRv4
	var destinationV6 *iaas.DestinationCIDRv6
	if model.DestinationV4 != nil {
		destinationV4 = &iaas.DestinationCIDRv4{
			Type:  utils.Ptr(destinationCIDRv4Type),
			Value: model.DestinationV4,
		}
	}
	if model.DestinationV6 != nil {
		destinationV6 = &iaas.DestinationCIDRv6{
			Type:  utils.Ptr(destinationCIDRv6Type),
			Value: model.DestinationV6,
		}
	}

	var nexthopIPv4 *iaas.NexthopIPv4
	var nexthopIPv6 *iaas.NexthopIPv6
	var nexthopBlackhole *iaas.NexthopBlackhole
	var nexthopInternet *iaas.NexthopInternet

	if model.NexthopV4 != nil {
		nexthopIPv4 = &iaas.NexthopIPv4{
			Type:  utils.Ptr(nexthopIPv4Type),
			Value: model.NexthopV4,
		}
	} else if model.NexthopV6 != nil {
		nexthopIPv6 = &iaas.NexthopIPv6{
			Type:  utils.Ptr(nexthopIPv6Type),
			Value: model.NexthopV6,
		}
	} else if model.NexthopBlackhole != nil {
		nexthopBlackhole = &iaas.NexthopBlackhole{
			Type: utils.Ptr(nexthopBlackholeType),
		}
	} else if model.NexthopInternet != nil {
		nexthopInternet = &iaas.NexthopInternet{
			Type: utils.Ptr(nexthopInternetType),
		}
	}

	payload := iaas.CreateNetworkAreaRoutePayload{
		Items: &[]iaas.Route{
			{
				Destination: &iaas.RouteDestination{
					DestinationCIDRv4: destinationV4,
					DestinationCIDRv6: destinationV6,
				},
				Nexthop: &iaas.RouteNexthop{
					NexthopIPv4:      nexthopIPv4,
					NexthopIPv6:      nexthopIPv6,
					NexthopBlackhole: nexthopBlackhole,
					NexthopInternet:  nexthopInternet,
				},
				Labels: utils.ConvertStringMapToInterfaceMap(model.Labels),
			},
		},
	}
	return req.CreateNetworkAreaRoutePayload(payload)
}

func outputResult(p *print.Printer, outputFormat, networkAreaLabel string, route iaas.Route) error {
	return p.OutputResult(outputFormat, route, func() error {
		p.Outputf("Created static route for SNA %q.\nStatic route ID: %s\n", networkAreaLabel, utils.PtrString(route.Id))
		return nil
	})
}
