package create

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	vpn "github.com/stackitcloud/stackit-sdk-go/services/vpn/v1api"

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

type tunnelInputModel struct {
	BgpRemoteAsn               *int64
	PeeringLocalAddress        *string
	PeeringRemoteAddress       *string
	Phase1DhGroups             []vpn.PhaseDhGroupsInner
	Phase1EncryptionAlgorithms []vpn.PhaseEncryptionAlgorithmsInner
	Phase1IntegrityAlgorithms  []vpn.PhaseIntegrityAlgorithmsInner
	Phase1RekeyTime            *int32
	Phase2DhGroups             []vpn.PhaseDhGroupsInner
	Phase2EncryptionAlgorithms []vpn.PhaseEncryptionAlgorithmsInner
	Phase2IntegrityAlgorithms  []vpn.PhaseIntegrityAlgorithmsInner
	Phase2RekeyTime            *int32
	Phase2DpdAction            *vpn.TunnelConfigurationPhase2AllOfDpdAction
	Phase2StartAction          *vpn.TunnelConfigurationPhase2AllOfStartAction
	PreSharedKey               string
	RemoteAddress              string
}

type inputModel struct {
	*globalflags.GlobalFlagModel
	GatewayId string

	DisplayName   string
	Enabled       *bool
	Labels        *map[string]string
	LocalSubnets  []string
	RemoteSubnets []string
	StaticRoutes  []string

	Tunnel1 tunnelInputModel
	Tunnel2 tunnelInputModel
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

		Tunnel1: tunnelInputModel{
			BgpRemoteAsn:               flags.FlagToInt64Pointer(p, cmd, tunnel1BgpRemoteAsnFlag),
			PeeringLocalAddress:        flags.FlagToStringPointer(p, cmd, tunnel1PeeringLocalAddressFlag),
			PeeringRemoteAddress:       flags.FlagToStringPointer(p, cmd, tunnel1PeeringRemoteAddressFlag),
			Phase1DhGroups:             tunnel1Phase1DhGroupsFlag.Get(),
			Phase1EncryptionAlgorithms: tunnel1Phase1EncryptionAlgorithmsFlag.Get(),
			Phase1IntegrityAlgorithms:  tunnel1Phase1IntegrityAlgorithmsFlag.Get(),
			Phase1RekeyTime:            flags.FlagToInt32Pointer(p, cmd, tunnel1Phase1RekeyTimeFlag),
			Phase2DhGroups:             tunnel1Phase2DhGroupsFlag.Get(),
			Phase2EncryptionAlgorithms: tunnel1Phase2EncryptionAlgorithmsFlag.Get(),
			Phase2IntegrityAlgorithms:  tunnel1Phase2IntegrityAlgorithmsFlag.Get(),
			Phase2RekeyTime:            flags.FlagToInt32Pointer(p, cmd, tunnel1Phase2RekeyTimeFlag),
			Phase2DpdAction:            tunnel1Phase2DpdActionFlag.Ptr(),
			Phase2StartAction:          tunnel1Phase2StartActionFlag.Ptr(),
			PreSharedKey:               flags.FlagToStringValue(p, cmd, tunnel1PreSharedKeyFlag),
			RemoteAddress:              flags.FlagToStringValue(p, cmd, tunnel1RemoteAddressFlag),
		},

		Tunnel2: tunnelInputModel{
			BgpRemoteAsn:               flags.FlagToInt64Pointer(p, cmd, tunnel2BgpRemoteAsnFlag),
			PeeringLocalAddress:        flags.FlagToStringPointer(p, cmd, tunnel2PeeringLocalAddressFlag),
			PeeringRemoteAddress:       flags.FlagToStringPointer(p, cmd, tunnel2PeeringRemoteAddressFlag),
			Phase1DhGroups:             tunnel2Phase1DhGroupsFlag.Get(),
			Phase1EncryptionAlgorithms: tunnel2Phase1EncryptionAlgorithmsFlag.Get(),
			Phase1IntegrityAlgorithms:  tunnel2Phase1IntegrityAlgorithmsFlag.Get(),
			Phase1RekeyTime:            flags.FlagToInt32Pointer(p, cmd, tunnel2Phase1RekeyTimeFlag),
			Phase2DhGroups:             tunnel2Phase2DhGroupsFlag.Get(),
			Phase2EncryptionAlgorithms: tunnel2Phase2EncryptionAlgorithmsFlag.Get(),
			Phase2IntegrityAlgorithms:  tunnel2Phase2IntegrityAlgorithmsFlag.Get(),
			Phase2RekeyTime:            flags.FlagToInt32Pointer(p, cmd, tunnel2Phase2RekeyTimeFlag),
			Phase2DpdAction:            tunnel2Phase2DpdActionFlag.Ptr(),
			Phase2StartAction:          tunnel2Phase2StartActionFlag.Ptr(),
			PreSharedKey:               flags.FlagToStringValue(p, cmd, tunnel2PreSharedKeyFlag),
			RemoteAddress:              flags.FlagToStringValue(p, cmd, tunnel2RemoteAddressFlag),
		},
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildTunnelConfiguration(model tunnelInputModel) vpn.TunnelConfiguration {
	tunnel := vpn.TunnelConfiguration{
		RemoteAddress: model.RemoteAddress,
	}
	if model.BgpRemoteAsn != nil {
		tunnel.Bgp = &vpn.BGPTunnelConfig{
			RemoteAsn: *model.BgpRemoteAsn,
		}
	}
	if model.PeeringLocalAddress != nil || model.PeeringRemoteAddress != nil {
		tunnel.Peering = &vpn.PeeringConfig{
			LocalAddress:  model.PeeringLocalAddress,
			RemoteAddress: model.PeeringRemoteAddress,
		}
	}
	tunnel.Phase1 = vpn.TunnelConfigurationPhase1{
		DhGroups:             model.Phase1DhGroups,
		EncryptionAlgorithms: model.Phase1EncryptionAlgorithms,
		IntegrityAlgorithms:  model.Phase1IntegrityAlgorithms,
		RekeyTime:            model.Phase1RekeyTime,
	}
	tunnel.Phase2 = vpn.TunnelConfigurationPhase2{
		DhGroups:             model.Phase2DhGroups,
		EncryptionAlgorithms: model.Phase2EncryptionAlgorithms,
		IntegrityAlgorithms:  model.Phase2IntegrityAlgorithms,
		RekeyTime:            model.Phase2RekeyTime,
		DpdAction:            model.Phase2DpdAction,
		StartAction:          model.Phase2StartAction,
	}
	tunnel.PreSharedKey = &model.PreSharedKey
	return tunnel
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

	payload.Tunnel1 = buildTunnelConfiguration(model.Tunnel1)
	payload.Tunnel2 = buildTunnelConfiguration(model.Tunnel2)

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
