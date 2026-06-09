package create

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	cliErr "github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/projectname"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/vpn/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	vpn "github.com/stackitcloud/stackit-sdk-go/services/vpn/v1api"
)

const (
	gatewayIdFlag = "gateway-id"

	displayNameFlag   = "display-name"
	enabledFlag       = "enabled"
	labelsFlag        = "labels"
	localSubnetsFlag  = "local-subnets"
	remoteSubnetsFlag = "remote-subnets"
	staticRoutesFlag  = "static-routes"

	tunnel1BgpRemoteAsnFlag         = "tunnel1-bgp-remote-asn"
	tunnel1PeeringLocalAddressFlag  = "tunnel1-peering-local-address"
	tunnel1PeeringRemoteAddressFlag = "tunnel1-peering-remote-address"
	tunnel1Phase1RekeyTimeFlag      = "tunnel1-phase1-rekey-time"
	tunnel1Phase2RekeyTimeFlag      = "tunnel1-phase2-rekey-time"
	tunnel1PreSharedKeyFlag         = "tunnel1-pre-shared-key"
	tunnel1RemoteAddressFlag        = "tunnel1-remote-address"

	tunnel2BgpRemoteAsnFlag         = "tunnel2-bgp-remote-asn"
	tunnel2PeeringLocalAddressFlag  = "tunnel2-peering-local-address"
	tunnel2PeeringRemoteAddressFlag = "tunnel2-peering-remote-address"
	tunnel2Phase1RekeyTimeFlag      = "tunnel2-phase1-rekey-time"
	tunnel2Phase2RekeyTimeFlag      = "tunnel2-phase2-rekey-time"
	tunnel2PreSharedKeyFlag         = "tunnel2-pre-shared-key"
	tunnel2RemoteAddressFlag        = "tunnel2-remote-address"
)

var (
	// tunnel 1
	tunnel1Phase1DhGroupsFlag = flags.StringEnumSliceFlag(
		"tunnel1-phase1-dh-groups",
		vpn.AllowedPhaseDhGroupsInnerEnumValues,
		"Tunnel 1 Phase 1 DH Groups.\nThe Diffie-Hellman Group. Required, except if AEAD algorithms are selected.",
	)
	tunnel1Phase1EncryptionAlgorithmsFlag = flags.StringEnumSliceFlag(
		"tunnel1-phase1-encryption-algorithms",
		vpn.AllowedPhaseEncryptionAlgorithmsInnerEnumValues,
		"Required: Tunnel 1 Phase 1 Encryption Algorithms",
	)
	tunnel1Phase1IntegrityAlgorithmsFlag = flags.StringEnumSliceFlag(
		"tunnel1-phase1-integrity-algorithms",
		vpn.AllowedPhaseIntegrityAlgorithmsInnerEnumValues,
		"Required: Tunnel 1 Phase 1 Integrity Algorithms",
	)
	tunnel1Phase2DhGroupsFlag = flags.StringEnumSliceFlag(
		"tunnel1-phase2-dh-groups",
		vpn.AllowedPhaseDhGroupsInnerEnumValues,
		"Tunnel 1 Phase 2 DH Groups",
	)
	tunnel1Phase2EncryptionAlgorithmsFlag = flags.StringEnumSliceFlag(
		"tunnel1-phase2-encryption-algorithms",
		vpn.AllowedPhaseEncryptionAlgorithmsInnerEnumValues,
		"Required: Tunnel 1 Phase 2 Encryption Algorithms",
	)
	tunnel1Phase2IntegrityAlgorithmsFlag = flags.StringEnumSliceFlag(
		"tunnel1-phase2-integrity-algorithms",
		vpn.AllowedPhaseIntegrityAlgorithmsInnerEnumValues,
		"Required: Tunnel 1 Phase 2 Integrity Algorithms",
	)
	tunnel1Phase2DpdActionFlag = flags.StringEnumFlag(
		"tunnel1-phase2-dpd-action",
		vpn.AllowedTunnelConfigurationPhase2AllOfDpdActionEnumValues,
		"Tunnel 1 Phase 2 DPD Action.\nAction to perform for this CHILD_SA on DPD timeout. \"clear\": Closes the CHILD_SA and does not take further action. \"restart\": immediately tries to re-negotiate the CILD_SA under a fresh IKE_SA.",
	)
	tunnel1Phase2StartActionFlag = flags.StringEnumFlag(
		"tunnel1-phase2-start-action",
		vpn.AllowedTunnelConfigurationPhase2AllOfStartActionEnumValues,
		"Tunnel 1 Phase 2 Start Action.\nAction to perform after loading the connection configuration. \"none\": The connection will be loaded but needs to be manually initiated. \"start\": initiates the connection actively.",
	)
	// tunnel 2
	tunnel2Phase1DhGroupsFlag = flags.StringEnumSliceFlag(
		"tunnel2-phase1-dh-groups",
		vpn.AllowedPhaseDhGroupsInnerEnumValues,
		"Tunnel 2 Phase 1 DH Groups\nThe Diffie-Hellman Group. Required, except if AEAD algorithms are selected.",
	)
	tunnel2Phase1EncryptionAlgorithmsFlag = flags.StringEnumSliceFlag(
		"tunnel2-phase1-encryption-algorithms",
		vpn.AllowedPhaseEncryptionAlgorithmsInnerEnumValues,
		"Required: Tunnel 2 Phase 1 Encryption Algorithms",
	)
	tunnel2Phase1IntegrityAlgorithmsFlag = flags.StringEnumSliceFlag(
		"tunnel2-phase1-integrity-algorithms",
		vpn.AllowedPhaseIntegrityAlgorithmsInnerEnumValues,
		"Required: Tunnel 2 Phase 1 Integrity Algorithms",
	)
	tunnel2Phase2DhGroupsFlag = flags.StringEnumSliceFlag(
		"tunnel2-phase2-dh-groups",
		vpn.AllowedPhaseDhGroupsInnerEnumValues,
		"Tunnel 2 Phase 2 DH Groups",
	)
	tunnel2Phase2EncryptionAlgorithmsFlag = flags.StringEnumSliceFlag(
		"tunnel2-phase2-encryption-algorithms",
		vpn.AllowedPhaseEncryptionAlgorithmsInnerEnumValues,
		"Required: Tunnel 2 Phase 2 Encryption Algorithms",
	)
	tunnel2Phase2IntegrityAlgorithmsFlag = flags.StringEnumSliceFlag(
		"tunnel2-phase2-integrity-algorithms",
		vpn.AllowedPhaseIntegrityAlgorithmsInnerEnumValues,
		"Required: Tunnel 2 Phase 2 Integrity Algorithms",
	)
	tunnel2Phase2DpdActionFlag = flags.StringEnumFlag(
		"tunnel2-phase2-dpd-action",
		vpn.AllowedTunnelConfigurationPhase2AllOfDpdActionEnumValues,
		"Tunnel 2 Phase 2 DPD Action.\nAction to perform for this CHILD_SA on DPD timeout. \"clear\": Closes the CHILD_SA and does not take further action. \"restart\": immediately tries to re-negotiate the CILD_SA under a fresh IKE_SA.",
	)
	tunnel2Phase2StartActionFlag = flags.StringEnumFlag(
		"tunnel2-phase2-start-action",
		vpn.AllowedTunnelConfigurationPhase2AllOfStartActionEnumValues,
		"Tunnel 2 Phase 2 Start Action.\nDefault: \"start\"\nEnum: \"none\" \"start\"\nAction to perform after loading the connection configuration. \"none\": The connection will be loaded but needs to be manually initiated. \"start\": initiates the connection actively.",
	)
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	GatewayId string

	DisplayName   string
	Enabled       *bool
	Labels        *map[string]string
	LocalSubnets  []string
	RemoteSubnets []string
	StaticRoutes  []string

	Tunnel1BgpRemoteAsn               *int64
	Tunnel1PeeringLocalAddress        *string
	Tunnel1PeeringRemoteAddress       *string
	Tunnel1Phase1DhGroups             []vpn.PhaseDhGroupsInner
	Tunnel1Phase1EncryptionAlgorithms []vpn.PhaseEncryptionAlgorithmsInner
	Tunnel1Phase1IntegrityAlgorithms  []vpn.PhaseIntegrityAlgorithmsInner
	Tunnel1Phase1RekeyTime            *int32
	Tunnel1Phase2DhGroups             []vpn.PhaseDhGroupsInner
	Tunnel1Phase2EncryptionAlgorithms []vpn.PhaseEncryptionAlgorithmsInner
	Tunnel1Phase2IntegrityAlgorithms  []vpn.PhaseIntegrityAlgorithmsInner
	Tunnel1Phase2RekeyTime            *int32
	Tunnel1Phase2DpdAction            *vpn.TunnelConfigurationPhase2AllOfDpdAction
	Tunnel1Phase2StartAction          *vpn.TunnelConfigurationPhase2AllOfStartAction
	Tunnel1PreSharedKey               string
	Tunnel1RemoteAddress              string

	Tunnel2BgpRemoteAsn               *int64
	Tunnel2PeeringLocalAddress        *string
	Tunnel2PeeringRemoteAddress       *string
	Tunnel2Phase1DhGroups             []vpn.PhaseDhGroupsInner
	Tunnel2Phase1EncryptionAlgorithms []vpn.PhaseEncryptionAlgorithmsInner
	Tunnel2Phase1IntegrityAlgorithms  []vpn.PhaseIntegrityAlgorithmsInner
	Tunnel2Phase1RekeyTime            *int32
	Tunnel2Phase2DhGroups             []vpn.PhaseDhGroupsInner
	Tunnel2Phase2EncryptionAlgorithms []vpn.PhaseEncryptionAlgorithmsInner
	Tunnel2Phase2IntegrityAlgorithms  []vpn.PhaseIntegrityAlgorithmsInner
	Tunnel2Phase2RekeyTime            *int32
	Tunnel2Phase2DpdAction            *vpn.TunnelConfigurationPhase2AllOfDpdAction
	Tunnel2Phase2StartAction          *vpn.TunnelConfigurationPhase2AllOfStartAction
	Tunnel2PreSharedKey               string
	Tunnel2RemoteAddress              string
}

func NewCmd(p *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Creates a VPN connection",
		Long:  "Creates a VPN connection.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Create a VPN connection`,
				"$ stackit beta vpn connection create --gateway-id xxx --display-name my-connection --tunnel1-remote-address 1.2.3.4 --tunnel2-remote-address 5.6.7.8"),
		),
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := context.Background()
			model, err := parseInput(p.Printer, cmd)
			if err != nil {
				return err
			}

			// Configure API client
			apiClient, err := client.ConfigureClient(p.Printer, p.CliVersion)
			if err != nil {
				return err
			}

			projectLabel, err := projectname.GetProjectName(ctx, p.Printer, p.CliVersion, cmd)
			if err != nil {
				p.Printer.Debug(print.ErrorLevel, "get project name: %v", err)
				projectLabel = model.ProjectId
			}

			prompt := fmt.Sprintf("Are you sure you want to create a VPN connection for gateway %q?", model.GatewayId)
			err = p.Printer.PromptForConfirmation(prompt)
			if err != nil {
				return err
			}

			// Call API
			req, err := buildRequest(ctx, model, apiClient)
			if err != nil {
				return err
			}
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("create VPN connection: %w", err)
			}

			return outputResult(p.Printer, model, projectLabel, resp)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), gatewayIdFlag, "Required: Gateway ID")
	cmd.Flags().String(displayNameFlag, "", "Required: A user friendly name for the connection.")
	cmd.Flags().Bool(enabledFlag, true, "Enable the connection")
	cmd.Flags().StringToString(labelsFlag, nil, "Map of custom labels. Key and values must be a string with max 63 chars, start/end with alphanumeric. The key of a label follows the same rules as the LabelValue except that it cannot be empty. (example: foo=bar)")
	cmd.Flags().StringSlice(localSubnetsFlag, nil, "Defaults to 0.0.0.0/0 for Route-based VPN configurations. Mandatory for Policy-based.")
	cmd.Flags().StringSlice(remoteSubnetsFlag, nil, "Defaults to 0.0.0.0/0 for Route-based VPN configurations. Mandatory for Policy-based.")
	cmd.Flags().StringSlice(staticRoutesFlag, nil, "Use this for route-based VPN.")

	cmd.Flags().Int64(tunnel1BgpRemoteAsnFlag, 0, "Required: Tunnel 1 BGP Remote ASN.\nASN for private use (reserved by IANA), both 16Bit and 32Bit ranges are valid (RFC 6996).")
	cmd.Flags().String(tunnel1PeeringLocalAddressFlag, "", "Tunnel 1 Peering Local Address.\nThe peering object defines the point-to-point IP configuration for the Tunnel Interface. These addresses serve as next-hop identifiers and are used for BGP peering sessions and can be used in Static Route-Based connectivity.")
	cmd.Flags().String(tunnel1PeeringRemoteAddressFlag, "", "Tunnel 1 Peering Remote Address")
	tunnel1Phase1DhGroupsFlag.Register(cmd)
	tunnel1Phase1EncryptionAlgorithmsFlag.Register(cmd)
	tunnel1Phase1IntegrityAlgorithmsFlag.Register(cmd)
	cmd.Flags().Int64(tunnel1Phase1RekeyTimeFlag, 0, "Tunnel 1 Phase 1 Rekey Time.\nTime to schedule a IKE re-keying (in seconds).")
	tunnel1Phase2DhGroupsFlag.Register(cmd)
	tunnel1Phase2EncryptionAlgorithmsFlag.Register(cmd)
	tunnel1Phase2IntegrityAlgorithmsFlag.Register(cmd)
	cmd.Flags().Int64(tunnel1Phase2RekeyTimeFlag, 0, "Tunnel 1 Phase 2 Rekey Time.\nTime to schedule a Child SA re-keying (in seconds).")
	tunnel1Phase2DpdActionFlag.Register(cmd)
	tunnel1Phase2StartActionFlag.Register(cmd)
	cmd.Flags().String(tunnel1PreSharedKeyFlag, "", "Required: Tunnel 1 Pre Shared Key.\nA Pre-Shared Key for authentication. Required in create-requests, optional in update-requests and omitted in every response.")
	cmd.Flags().String(tunnel1RemoteAddressFlag, "", "Tunnel 1 Remote Address")

	cmd.Flags().Int64(tunnel2BgpRemoteAsnFlag, 0, "Tunnel 2 BGP Remote ASN")
	cmd.Flags().String(tunnel2PeeringLocalAddressFlag, "", "Tunnel 2 Peering Local Address.\nThe peering object defines the point-to-point IP configuration for the Tunnel Interface. These addresses serve as next-hop identifiers and are used for BGP peering sessions and can be used in Static Route-Based connectivity.")
	cmd.Flags().String(tunnel2PeeringRemoteAddressFlag, "", "Tunnel 2 Peering Remote Address")
	tunnel2Phase1DhGroupsFlag.Register(cmd)
	tunnel2Phase1EncryptionAlgorithmsFlag.Register(cmd)
	tunnel2Phase1IntegrityAlgorithmsFlag.Register(cmd)
	cmd.Flags().Int64(tunnel2Phase1RekeyTimeFlag, 0, "Tunnel 2 Phase 1 Rekey Time.\nTime to schedule a IKE re-keying (in seconds).")
	tunnel2Phase2DhGroupsFlag.Register(cmd)
	tunnel2Phase2EncryptionAlgorithmsFlag.Register(cmd)
	tunnel2Phase2IntegrityAlgorithmsFlag.Register(cmd)
	cmd.Flags().Int64(tunnel2Phase2RekeyTimeFlag, 0, "Tunnel 2 Phase 2 Rekey Time.\nTime to schedule a Child SA re-keying (in seconds).")
	tunnel2Phase2DpdActionFlag.Register(cmd)
	tunnel2Phase2StartActionFlag.Register(cmd)
	cmd.Flags().String(tunnel2PreSharedKeyFlag, "", "Required: Tunnel 2 Pre Shared Key.\nA Pre-Shared Key for authentication. Required in create-requests, optional in update-requests and omitted in every response.")
	cmd.Flags().String(tunnel2RemoteAddressFlag, "", "Tunnel 2 Remote Address")

	err := flags.MarkFlagsRequired(
		cmd,
		gatewayIdFlag, displayNameFlag,
		tunnel1RemoteAddressFlag,
		tunnel1PreSharedKeyFlag,
		tunnel1Phase1EncryptionAlgorithmsFlag.Name(), tunnel1Phase1IntegrityAlgorithmsFlag.Name(),
		tunnel1Phase2EncryptionAlgorithmsFlag.Name(), tunnel1Phase2IntegrityAlgorithmsFlag.Name(),
		tunnel2RemoteAddressFlag,
		tunnel2PreSharedKeyFlag,
		tunnel2Phase1EncryptionAlgorithmsFlag.Name(), tunnel2Phase1IntegrityAlgorithmsFlag.Name(),
		tunnel2Phase2EncryptionAlgorithmsFlag.Name(), tunnel2Phase2IntegrityAlgorithmsFlag.Name(),
	)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &cliErr.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		GatewayId:       flags.FlagToStringValue(p, cmd, gatewayIdFlag),

		DisplayName:   flags.FlagToStringValue(p, cmd, displayNameFlag),
		Enabled:       flags.FlagToBoolPointer(p, cmd, enabledFlag),
		Labels:        flags.FlagToStringToStringPointer(p, cmd, labelsFlag),
		LocalSubnets:  flags.FlagToStringSliceValue(p, cmd, localSubnetsFlag),
		RemoteSubnets: flags.FlagToStringSliceValue(p, cmd, remoteSubnetsFlag),
		StaticRoutes:  flags.FlagToStringSliceValue(p, cmd, staticRoutesFlag),

		Tunnel1BgpRemoteAsn:               flags.FlagToInt64Pointer(p, cmd, tunnel1BgpRemoteAsnFlag),
		Tunnel1PeeringLocalAddress:        flags.FlagToStringPointer(p, cmd, tunnel1PeeringLocalAddressFlag),
		Tunnel1PeeringRemoteAddress:       flags.FlagToStringPointer(p, cmd, tunnel1PeeringRemoteAddressFlag),
		Tunnel1Phase1DhGroups:             tunnel1Phase1DhGroupsFlag.Get(),
		Tunnel1Phase1EncryptionAlgorithms: tunnel1Phase1EncryptionAlgorithmsFlag.Get(),
		Tunnel1Phase1IntegrityAlgorithms:  tunnel1Phase1IntegrityAlgorithmsFlag.Get(),
		Tunnel1Phase1RekeyTime:            flags.FlagToInt32Pointer(p, cmd, tunnel1Phase1RekeyTimeFlag),
		Tunnel1Phase2DhGroups:             tunnel1Phase2DhGroupsFlag.Get(),
		Tunnel1Phase2EncryptionAlgorithms: tunnel1Phase2EncryptionAlgorithmsFlag.Get(),
		Tunnel1Phase2IntegrityAlgorithms:  tunnel1Phase2IntegrityAlgorithmsFlag.Get(),
		Tunnel1Phase2RekeyTime:            flags.FlagToInt32Pointer(p, cmd, tunnel1Phase2RekeyTimeFlag),
		Tunnel1Phase2DpdAction:            tunnel1Phase2DpdActionFlag.Ptr(),
		Tunnel1Phase2StartAction:          tunnel1Phase2StartActionFlag.Ptr(),
		Tunnel1PreSharedKey:               flags.FlagToStringValue(p, cmd, tunnel1PreSharedKeyFlag),
		Tunnel1RemoteAddress:              flags.FlagToStringValue(p, cmd, tunnel1RemoteAddressFlag),

		Tunnel2BgpRemoteAsn:               flags.FlagToInt64Pointer(p, cmd, tunnel2BgpRemoteAsnFlag),
		Tunnel2PeeringLocalAddress:        flags.FlagToStringPointer(p, cmd, tunnel2PeeringLocalAddressFlag),
		Tunnel2PeeringRemoteAddress:       flags.FlagToStringPointer(p, cmd, tunnel2PeeringRemoteAddressFlag),
		Tunnel2Phase1DhGroups:             tunnel2Phase1DhGroupsFlag.Get(),
		Tunnel2Phase1EncryptionAlgorithms: tunnel2Phase1EncryptionAlgorithmsFlag.Get(),
		Tunnel2Phase1IntegrityAlgorithms:  tunnel2Phase1IntegrityAlgorithmsFlag.Get(),
		Tunnel2Phase1RekeyTime:            flags.FlagToInt32Pointer(p, cmd, tunnel2Phase1RekeyTimeFlag),
		Tunnel2Phase2DhGroups:             tunnel2Phase2DhGroupsFlag.Get(),
		Tunnel2Phase2EncryptionAlgorithms: tunnel2Phase2EncryptionAlgorithmsFlag.Get(),
		Tunnel2Phase2IntegrityAlgorithms:  tunnel2Phase2IntegrityAlgorithmsFlag.Get(),
		Tunnel2Phase2RekeyTime:            flags.FlagToInt32Pointer(p, cmd, tunnel2Phase2RekeyTimeFlag),
		Tunnel2Phase2DpdAction:            tunnel2Phase2DpdActionFlag.Ptr(),
		Tunnel2Phase2StartAction:          tunnel2Phase2StartActionFlag.Ptr(),
		Tunnel2PreSharedKey:               flags.FlagToStringValue(p, cmd, tunnel2PreSharedKeyFlag),
		Tunnel2RemoteAddress:              flags.FlagToStringValue(p, cmd, tunnel2RemoteAddressFlag),
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *vpn.APIClient) (vpn.ApiCreateGatewayConnectionRequest, error) {
	req := apiClient.DefaultAPI.CreateGatewayConnection(ctx, model.ProjectId, model.Region, model.GatewayId)

	payload := vpn.CreateGatewayConnectionPayload{
		DisplayName:   model.DisplayName,
		Enabled:       model.Enabled,
		Labels:        model.Labels,
		LocalSubnets:  model.LocalSubnets,
		RemoteSubnets: model.RemoteSubnets,
		StaticRoutes:  model.StaticRoutes,
	}

	tunnel1 := vpn.TunnelConfiguration{
		RemoteAddress: model.Tunnel1RemoteAddress,
	}
	if model.Tunnel1BgpRemoteAsn != nil {
		tunnel1.Bgp = &vpn.BGPTunnelConfig{
			RemoteAsn: *model.Tunnel1BgpRemoteAsn,
		}
	}
	if model.Tunnel1PeeringLocalAddress != nil || model.Tunnel1PeeringRemoteAddress != nil {
		tunnel1.Peering = &vpn.PeeringConfig{
			LocalAddress:  model.Tunnel1PeeringLocalAddress,
			RemoteAddress: model.Tunnel1PeeringRemoteAddress,
		}
	}
	tunnel1.Phase1 = vpn.TunnelConfigurationPhase1{
		DhGroups:             model.Tunnel1Phase1DhGroups,
		EncryptionAlgorithms: model.Tunnel1Phase1EncryptionAlgorithms,
		IntegrityAlgorithms:  model.Tunnel1Phase1IntegrityAlgorithms,
		RekeyTime:            model.Tunnel1Phase1RekeyTime,
	}
	tunnel1.Phase2 = vpn.TunnelConfigurationPhase2{
		DhGroups:             model.Tunnel1Phase2DhGroups,
		EncryptionAlgorithms: model.Tunnel1Phase2EncryptionAlgorithms,
		IntegrityAlgorithms:  model.Tunnel1Phase2IntegrityAlgorithms,
		RekeyTime:            model.Tunnel1Phase2RekeyTime,
		DpdAction:            model.Tunnel1Phase2DpdAction,
		StartAction:          model.Tunnel1Phase2StartAction,
	}
	tunnel1.PreSharedKey = &model.Tunnel1PreSharedKey
	payload.Tunnel1 = tunnel1

	tunnel2 := vpn.TunnelConfiguration{
		RemoteAddress: model.Tunnel2RemoteAddress,
	}
	if model.Tunnel2BgpRemoteAsn != nil {
		tunnel2.Bgp = &vpn.BGPTunnelConfig{
			RemoteAsn: *model.Tunnel2BgpRemoteAsn,
		}
	}
	if model.Tunnel2PeeringLocalAddress != nil || model.Tunnel2PeeringRemoteAddress != nil {
		tunnel2.Peering = &vpn.PeeringConfig{
			LocalAddress:  model.Tunnel2PeeringLocalAddress,
			RemoteAddress: model.Tunnel2PeeringRemoteAddress,
		}
	}
	tunnel2.Phase1 = vpn.TunnelConfigurationPhase1{
		DhGroups:             model.Tunnel2Phase1DhGroups,
		EncryptionAlgorithms: model.Tunnel2Phase1EncryptionAlgorithms,
		IntegrityAlgorithms:  model.Tunnel2Phase1IntegrityAlgorithms,
		RekeyTime:            model.Tunnel2Phase1RekeyTime,
	}

	tunnel2.Phase2 = vpn.TunnelConfigurationPhase2{
		DhGroups:             model.Tunnel2Phase2DhGroups,
		EncryptionAlgorithms: model.Tunnel2Phase2EncryptionAlgorithms,
		IntegrityAlgorithms:  model.Tunnel2Phase2IntegrityAlgorithms,
		RekeyTime:            model.Tunnel2Phase2RekeyTime,
		DpdAction:            model.Tunnel2Phase2DpdAction,
		StartAction:          model.Tunnel2Phase2StartAction,
	}
	tunnel2.PreSharedKey = &model.Tunnel2PreSharedKey
	payload.Tunnel2 = tunnel2

	return req.CreateGatewayConnectionPayload(payload), nil
}

func outputResult(p *print.Printer, model *inputModel, projectLabel string, resp *vpn.ConnectionResponse) error {
	if resp == nil {
		return fmt.Errorf("create response is empty")
	}
	return p.OutputResult(model.OutputFormat, resp, func() error {
		p.Outputf("Created VPN connection %q for gateway %q in project %q.\n", utils.PtrString(resp.Id), model.GatewayId, projectLabel)
		return nil
	})
}
