package update

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	vpn "github.com/stackitcloud/stackit-sdk-go/services/vpn/v1api"
	"github.com/stackitcloud/stackit-sdk-go/services/vpn/v1api/wait"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/projectname"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/vpn/client"
	vpnUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/vpn/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/spinner"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

const (
	gatewayIdArg = "GATEWAY_ID"

	availabilityZoneTunnel1Flag     = "availability-zone-tunnel-1"
	availabilityZoneTunnel2Flag     = "availability-zone-tunnel-2"
	bgpLocalAsnFlag                 = "bgp-local-asn"
	bgpOverrideAdvertisedRoutesFlag = "bgp-override-advertised-routes"
	nameFlag                        = "name"
	labelsFlag                      = "labels"
	planIdFlag                      = "plan-id"
	routingTypeFlag                 = "routing-type"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	GatewayId        string
	AvailabilityZone vpn.UpdateGatewayPayloadAvailabilityZones
	Bgp              *vpn.BGPGatewayConfig
	Name             string
	Labels           *map[string]string
	PlanId           string
	RoutingType      vpn.RoutingType
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("update %s", gatewayIdArg),
		Short: "Updates a vpn gateway",
		Long:  "Updates a vpn gateway.",
		Args:  args.SingleArg(gatewayIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Update vpn gateway with ID "xxx"`,
				"$ stackit beta vpn gateway update xxx",
			),
		),
		RunE: func(cmd *cobra.Command, inputArgs []string) error {
			ctx := context.Background()
			model, err := parseInput(params.Printer, cmd, inputArgs)
			if err != nil {
				return fmt.Errorf("unable to parse input: %w", err)
			}

			// Configure API client
			apiClient, err := client.ConfigureClient(params.Printer, params.CliVersion)
			if err != nil {
				return err
			}

			gatewayLabel, err := vpnUtils.GetGatewayName(ctx, apiClient.DefaultAPI, model.ProjectId, model.Region, model.GatewayId)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get gateway name: %v", err)
				gatewayLabel = model.GatewayId
			} else if gatewayLabel == "" {
				gatewayLabel = model.GatewayId
			}

			projectLabel, err := projectname.GetProjectName(ctx, params.Printer, params.CliVersion, cmd)
			if err != nil || projectLabel == "" {
				projectLabel = model.ProjectId
			}

			prompt := fmt.Sprintf("Are you sure you want to update vpn gateway %q for the project %q?", gatewayLabel, projectLabel)
			err = params.Printer.PromptForConfirmation(prompt)
			if err != nil {
				return err
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("update vpn gateway: %w", err)
			}

			// Wait for async operation, if async mode not enabled
			if !model.Async {
				err := spinner.Run(params.Printer, "Updating gateway", func() error {
					_, err = wait.UpdateGatewayWaitHandler(ctx, apiClient.DefaultAPI, model.ProjectId, vpn.Region(model.Region), model.GatewayId).WaitWithContext(ctx)
					return err
				})
				if err != nil {
					return fmt.Errorf("waiting for gateway update: %w", err)
				}
			}

			return outputResult(params.Printer, model.OutputFormat, model.Async, projectLabel, resp)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().String(availabilityZoneTunnel1Flag, "", "Availability Zone of Tunnel 1")
	cmd.Flags().String(availabilityZoneTunnel2Flag, "", "Availability Zone of Tunnel 2")
	cmd.Flags().Int64(bgpLocalAsnFlag, 0, "ASN for private use (reserved by IANA), both 16Bit and 32Bit ranges are valid (RFC 6996)")
	cmd.Flags().StringArray(bgpOverrideAdvertisedRoutesFlag, nil, "A list of IPv4 Prefixes to advertise via BGP")
	cmd.Flags().String(nameFlag, "", "Gateway name")
	cmd.Flags().StringToString(labelsFlag, nil, "Labels in key=value format, separated by commas")
	cmd.Flags().String(planIdFlag, "", "Plan ID")
	cmd.Flags().String(routingTypeFlag, "", "Routing Type: \"POLICY_BASED\", \"ROUTE_BASED\" or \"BGP_ROUTE_BASED\"")

	err := flags.MarkFlagsRequired(cmd, availabilityZoneTunnel1Flag, availabilityZoneTunnel2Flag, nameFlag, planIdFlag, routingTypeFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	gatewayId := inputArgs[0]
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	bgpLocalAsn := flags.FlagToInt64Pointer(p, cmd, bgpLocalAsnFlag)
	if bgpLocalAsn != nil {
		if *bgpLocalAsn < 0 {
			return nil, &errors.FlagValidationError{
				Flag:    bgpLocalAsnFlag,
				Details: "must be a positive integer",
			}
		}
	}
	bgpOverrideAdvertisedRoutes := flags.FlagToStringArrayValue(p, cmd, bgpOverrideAdvertisedRoutesFlag)

	var bgp *vpn.BGPGatewayConfig
	if bgpLocalAsn != nil || bgpOverrideAdvertisedRoutes != nil {
		bgp = &vpn.BGPGatewayConfig{
			LocalAsn:                 bgpLocalAsn,
			OverrideAdvertisedRoutes: flags.FlagToStringArrayValue(p, cmd, bgpOverrideAdvertisedRoutesFlag),
		}
	}

	model := inputModel{
		GatewayId:       gatewayId,
		GlobalFlagModel: globalFlags,
		AvailabilityZone: vpn.UpdateGatewayPayloadAvailabilityZones{
			Tunnel1: flags.FlagToStringValue(p, cmd, availabilityZoneTunnel1Flag),
			Tunnel2: flags.FlagToStringValue(p, cmd, availabilityZoneTunnel2Flag),
		},
		Bgp:         bgp,
		Name:        flags.FlagToStringValue(p, cmd, nameFlag),
		Labels:      flags.FlagToStringToStringPointer(p, cmd, labelsFlag),
		PlanId:      flags.FlagToStringValue(p, cmd, planIdFlag),
		RoutingType: vpn.RoutingType(flags.FlagToStringValue(p, cmd, routingTypeFlag)),
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *vpn.APIClient) vpn.ApiUpdateGatewayRequest {
	req := apiClient.DefaultAPI.UpdateGateway(ctx, model.ProjectId, vpn.Region(model.Region), model.GatewayId)
	req = req.UpdateGatewayPayload(vpn.UpdateGatewayPayload{
		AvailabilityZones: model.AvailabilityZone,
		Bgp:               model.Bgp,
		DisplayName:       model.Name,
		Labels:            model.Labels,
		PlanId:            model.PlanId,
		RoutingType:       model.RoutingType,
	})
	return req
}

func outputResult(p *print.Printer, outputFormat string, async bool, projectLabel string, item *vpn.GatewayResponse) error {
	return p.OutputResult(outputFormat, item, func() error {
		if item == nil {
			p.Outputln("vpn gateway response is empty")
			return nil
		}

		operation := "Updated"
		if async {
			operation = "Triggered update of"
		}
		p.Outputf(
			"%s vpn gateway %q in project %q.\n",
			operation,
			item.DisplayName,
			projectLabel,
		)
		return nil
	})
}
