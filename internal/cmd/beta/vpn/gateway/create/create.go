package create

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
	"github.com/stackitcloud/stackit-cli/internal/pkg/spinner"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

const (
	availabilityZoneTunnel1Flag     = "availability-zone-tunnel-1"
	availabilityZoneTunnel2Flag     = "availability-zone-tunnel-2"
	bgpLocalAsnFlag                 = "bgp-local-asn"
	bgpOverrideAdvertisedRoutesFlag = "bgp-override-advertised-routes"
	nameFlag                        = "name"
	labelsFlag                      = "labels"
	planIdFlag                      = "plan-id"
)

var (
	routingTypeFlag = flags.StringEnumFlag(
		"routing-type",
		vpn.AllowedRoutingTypeEnumValues,
		"Routing Type of the VPN",
	)
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	AvailabilityZone vpn.CreateGatewayPayloadAvailabilityZones
	Bgp              *vpn.BGPGatewayConfig
	Name             string
	Labels           *map[string]string
	PlanId           string
	RoutingType      vpn.RoutingType
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Creates a vpn gateway",
		Long:  "Creates a vpn gateway.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Create a vpn gateway with name "xxx", plan "p500", policy based routing and both tunnels in availability-zone eu01-1`,
				"$ stackit beta vpn gateway create --name xxx --plan-id p500 --routing-type POLICY_BASED --availability-zone-tunnel-1 eu01-1 --availability-zone-tunnel-2 eu01-1",
			),
			examples.NewExample(
				`Create a vpn gateway with the labels foo=bar and x=y`,
				"$ stackit beta vpn gateway create --name xxx --plan-id p500 --routing-type POLICY_BASED --availability-zone-tunnel-1 eu01-1 --availability-zone-tunnel-2 eu01-1 --label foo=bar,x=y",
			),
			examples.NewExample(
				`Create a vpn gateway with bgp enabled, yyy as local asn and [aaa, bbb] as override advertised routes`,
				"$ stackit beta vpn gateway create --name xxx --plan-id p500 --routing-type POLICY_BASED --availability-zone-tunnel-1 eu01-1 --availability-zone-tunnel-2 eu01-1 --bgp-local-asn yyy --bgp-override-advertised-routes aaa,bbb",
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			model, err := parseInput(params.Printer, cmd, args)
			if err != nil {
				return fmt.Errorf("unable to parse input: %w", err)
			}

			// Configure API client
			apiClient, err := client.ConfigureClient(params.Printer, params.CliVersion)
			if err != nil {
				return err
			}

			projectLabel, err := projectname.GetProjectName(ctx, params.Printer, params.CliVersion, cmd)
			if err != nil || projectLabel == "" {
				projectLabel = model.ProjectId
			}

			prompt := fmt.Sprintf("Are you sure you want to create a vpn gateway for project %q?", projectLabel)
			err = params.Printer.PromptForConfirmation(prompt)
			if err != nil {
				return err
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("create vpn gateway: %w", err)
			}
			var gatewayId string
			if resp != nil && resp.HasId() {
				gatewayId = *resp.Id
			}

			// Wait for async operation, if async mode not enabled
			if !model.Async {
				err := spinner.Run(params.Printer, "Creating vpn gateway", func() error {
					_, err = wait.CreateGatewayWaitHandler(ctx, apiClient.DefaultAPI, model.ProjectId, model.Region, gatewayId).WaitWithContext(ctx)
					return err
				})
				if err != nil {
					return fmt.Errorf("waiting for vpn gateway creation: %w", err)
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
	routingTypeFlag.Register(cmd.Flags())

	err := flags.MarkFlagsRequired(cmd, availabilityZoneTunnel1Flag, availabilityZoneTunnel2Flag, nameFlag, planIdFlag, routingTypeFlag.Name())
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command, _ []string) (*inputModel, error) {
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
	if bgpOverrideAdvertisedRoutes != nil {
		if bgpLocalAsn == nil {
			return nil, &errors.DependingFlagIsMissing{
				MissingFlag: bgpLocalAsnFlag,
				SetFlag:     bgpOverrideAdvertisedRoutesFlag,
			}
		}
	}

	var bgp *vpn.BGPGatewayConfig
	if bgpLocalAsn != nil || bgpOverrideAdvertisedRoutes != nil {
		bgp = &vpn.BGPGatewayConfig{
			LocalAsn:                 *bgpLocalAsn,
			OverrideAdvertisedRoutes: flags.FlagToStringArrayValue(p, cmd, bgpOverrideAdvertisedRoutesFlag),
		}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		AvailabilityZone: vpn.CreateGatewayPayloadAvailabilityZones{
			Tunnel1: flags.FlagToStringValue(p, cmd, availabilityZoneTunnel1Flag),
			Tunnel2: flags.FlagToStringValue(p, cmd, availabilityZoneTunnel2Flag),
		},
		Bgp:         bgp,
		Name:        flags.FlagToStringValue(p, cmd, nameFlag),
		Labels:      flags.FlagToStringToStringPointer(p, cmd, labelsFlag),
		PlanId:      flags.FlagToStringValue(p, cmd, planIdFlag),
		RoutingType: vpn.RoutingType(routingTypeFlag.Get()),
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *vpn.APIClient) vpn.ApiCreateGatewayRequest {
	req := apiClient.DefaultAPI.CreateGateway(ctx, model.ProjectId, model.Region)
	req = req.CreateGatewayPayload(
		vpn.CreateGatewayPayload{
			AvailabilityZones: model.AvailabilityZone,
			Bgp:               model.Bgp,
			DisplayName:       model.Name,
			Labels:            model.Labels,
			PlanId:            model.PlanId,
			RoutingType:       model.RoutingType,
		},
	)
	return req
}

func outputResult(p *print.Printer, outputFormat string, async bool, projectLabel string, item *vpn.GatewayResponse) error {
	return p.OutputResult(outputFormat, item, func() error {
		if item == nil {
			p.Outputln("vpn gateway response is empty")
			return nil
		}
		operation := "Created"
		if async {
			operation = "Triggered creation of"
		}
		p.Outputf(
			"%s vpn gateway %q in project %q.\nGateway ID: %s\n",
			operation,
			item.DisplayName,
			projectLabel,
			utils.PtrString(item.Id),
		)
		return nil
	})
}
