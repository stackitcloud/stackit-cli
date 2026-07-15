package status

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	vpn "github.com/stackitcloud/stackit-sdk-go/services/vpn/v1api"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/projectname"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/vpn/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

const (
	gatewayIdArg = "GATEWAY_ID"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	GatewayId string
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("status %s", gatewayIdArg),
		Short: "Shows the status of a gateway",
		Long:  "Shows the status of a gateway.",
		Args:  args.SingleArg(gatewayIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Show the status of the gateway with the ID "xxx"`,
				"$ stackit beta vpn gateway status xxx",
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

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("describe vpn gateway: %w", err)
			}

			projectLabel, err := projectname.GetProjectName(ctx, params.Printer, params.CliVersion, cmd)
			if err != nil || projectLabel == "" {
				projectLabel = model.ProjectId
			}

			return outputResult(params.Printer, model.OutputFormat, model.GatewayId, projectLabel, resp)
		},
	}
	return cmd
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	gatewayId := inputArgs[0]
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		GatewayId:       gatewayId,
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *vpn.APIClient) vpn.ApiGetGatewayStatusRequest {
	return apiClient.DefaultAPI.GetGatewayStatus(ctx, model.ProjectId, model.Region, model.GatewayId)
}

func outputResult(p *print.Printer, outputFormat, gatewayId, projectLabel string, gateway *vpn.GatewayStatusResponse) error {
	return p.OutputResult(outputFormat, gateway, func() error {
		if gateway == nil {
			p.Outputf("gateway %q not found in project %q\n", gatewayId, projectLabel)
			return nil
		}

		mainTable := tables.NewTable()
		mainTable.SetTitle("Gateway Status")

		mainTable.AddRow("ID", gateway.GetId())
		mainTable.AddSeparator()
		mainTable.AddRow("NAME", gateway.GetDisplayName())
		mainTable.AddSeparator()
		mainTable.AddRow("STATUS", gateway.GetGatewayStatus())
		if gateway.ErrorMessage != nil {
			mainTable.AddSeparator()
			mainTable.AddRow("ERROR MESSAGE", *gateway.ErrorMessage)
		}

		ts := []tables.Table{
			mainTable,
		}
		for _, tunnel := range gateway.Tunnels {
			ts = append(ts, tunnelTable(tunnel))
		}

		return tables.DisplayTables(p, ts)
	})
}

func tunnelTable(tunnel vpn.VPNTunnels) tables.Table {
	title := "Tunnel"
	if tunnel.Name != nil {
		title = string(*tunnel.Name)
	}

	table := tables.NewTable()
	table.SetTitle(title)

	table.AddSeparator()
	table.AddRow("PUBLIC IP", tunnel.GetPublicIP())
	table.AddSeparator()
	table.AddRow("INTERNAL NEXT HOP IP", tunnel.GetInternalNextHopIP())
	table.AddSeparator()
	table.AddRow("STATE", tunnel.GetInstanceState())

	if tunnel.BgpStatus.IsSet() {
		table.AddSeparator()
		routeString := ""
		for _, route := range tunnel.BgpStatus.Get().Routes {
			if route.Network != "" {
				routeString += fmt.Sprintf("Network: %s; ", route.Network)
			}
			if route.Origin != "" {
				routeString += fmt.Sprintf("Origin: %s; ", route.Origin)
			}
			if route.Path != "" {
				routeString += fmt.Sprintf("Path: %s; ", route.Path)
			}
			if route.PeerId != "" {
				routeString += fmt.Sprintf("PeerId: %s; ", route.PeerId)
			}
			routeString += fmt.Sprintf("Weight: %d\n", route.Weight)
		}
		table.AddRow("BGP Routes", routeString)
		table.AddSeparator()
		bgpPeers := ""
		for _, peer := range tunnel.BgpStatus.Get().Peers {
			if peer.PeerUptime != "" {
				bgpPeers += fmt.Sprintf("PeerUptime: %s; ", peer.PeerUptime)
			}
			if peer.RemoteIP != "" {
				bgpPeers += fmt.Sprintf("RemoteIP: %s; ", peer.RemoteIP)
			}
			if peer.State != "" {
				bgpPeers += fmt.Sprintf("State: %s; ", peer.State)
			}
			bgpPeers += fmt.Sprintf("LocalAsn: %d; ", peer.LocalAs)
			bgpPeers += fmt.Sprintf("PfxRcd: %d; ", peer.PfxRcd)
			bgpPeers += fmt.Sprintf("PfxSnt: %d; ", peer.PfxSnt)
			bgpPeers += fmt.Sprintf("RemoteAs: %d\n", peer.RemoteAs)
		}
		table.AddRow("BGP Peers", bgpPeers)
	}

	return table
}
