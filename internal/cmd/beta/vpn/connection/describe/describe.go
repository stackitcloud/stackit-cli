package describe

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	cliErr "github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/vpn/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	vpn "github.com/stackitcloud/stackit-sdk-go/services/vpn/v1api"
)

const (
	connectionIdArg = "CONNECTION_ID"

	gatewayIdFlag = "gateway-id"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	GatewayId    *string
	ConnectionId string
}

func NewCmd(p *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("describe %s", connectionIdArg),
		Short: "Shows details of a VPN connection",
		Long:  "Shows details of a VPN connection.",
		Args:  args.SingleArg(connectionIdArg, nil),
		Example: examples.Build(
			examples.NewExample(
				`Show details of a VPN connection`,
				"$ stackit beta vpn connection describe xxx --gateway-id yyy"),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			model, err := parseInput(p.Printer, cmd, args)
			if err != nil {
				return err
			}

			// Configure API client
			apiClient, err := client.ConfigureClient(p.Printer, p.CliVersion)
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
				return fmt.Errorf("describe VPN connection: %w", err)
			}

			return outputResult(p.Printer, model, resp)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), gatewayIdFlag, "Gateway ID")

	err := flags.MarkFlagsRequired(cmd, gatewayIdFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	connectionId := inputArgs[0]

	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &cliErr.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		GatewayId:       flags.FlagToStringPointer(p, cmd, gatewayIdFlag),
		ConnectionId:    connectionId,
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *vpn.APIClient) (vpn.ApiGetGatewayConnectionRequest, error) {
	req := apiClient.DefaultAPI.GetGatewayConnection(ctx, model.ProjectId, model.Region, *model.GatewayId, model.ConnectionId)
	return req, nil
}

func outputResult(p *print.Printer, model *inputModel, resp *vpn.ConnectionResponse) error {
	if resp == nil {
		return fmt.Errorf("describe response is empty")
	}

	return p.OutputResult(model.OutputFormat, resp, func() error {
		mainTable := tables.NewTable()
		mainTable.AddRow("ID", utils.PtrString(resp.Id))
		mainTable.AddRow("Name", resp.DisplayName)
		mainTable.AddRow("Enabled", utils.PtrString(resp.Enabled))
		var labels string
		if resp.Labels != nil {
			labels = utils.JoinStringMap(*resp.Labels, "=", ", ")
		}
		mainTable.AddRow("Labels", labels)
		mainTable.AddRow("Local Subnets", strings.Join(resp.LocalSubnets, ", "))
		mainTable.AddRow("Remote Subnets", strings.Join(resp.RemoteSubnets, ", "))
		mainTable.AddRow("Static Routes", strings.Join(resp.StaticRoutes, ", "))

		ts := []tables.Table{
			mainTable,
		}
		ts = append(ts, tunnelTables(resp.Tunnel1, "Tunnel 1")...)
		ts = append(ts, tunnelTables(resp.Tunnel2, "Tunnel 2")...)
		return tables.DisplayTables(p, ts)
	})
}

func tunnelTables(tunnel vpn.TunnelConfiguration, title string) []tables.Table {
	table := tables.NewTable()
	table.SetTitle(title)
	table.AddRow("IP Address", tunnel.RemoteAddress)
	var bgp string
	if tunnel.Bgp != nil {
		bgp = fmt.Sprintf("%d", tunnel.Bgp.RemoteAsn)
	}
	table.AddRow("BGP ASN", bgp)
	var peering string
	if tunnel.Peering != nil {
		peering = fmt.Sprintf("%s/%s", utils.PtrString(tunnel.Peering.LocalAddress), utils.PtrString(tunnel.Peering.RemoteAddress))
	}
	table.AddRow("Peering (local/remote)", peering)

	phase1Table := tables.NewTable()
	phase1Table.SetTitle(fmt.Sprintf("%s Phase 1", title))
	phase1Table.AddRow("DH Groups", utils.JoinStringPtr(&tunnel.Phase1.DhGroups, ", "))
	phase1Table.AddRow("Encryption Algos", utils.JoinStringPtr(&tunnel.Phase1.EncryptionAlgorithms, ", "))
	phase1Table.AddRow("Integrity Algos", utils.JoinStringPtr(&tunnel.Phase1.IntegrityAlgorithms, ", "))
	phase1Table.AddRow("Rekey Time", utils.PtrString(tunnel.Phase1.RekeyTime))

	phase2Table := tables.NewTable()
	phase2Table.SetTitle(fmt.Sprintf("%s Phase 2", title))
	phase2Table.AddRow("DH Groups", utils.JoinStringPtr(&tunnel.Phase2.DhGroups, ", "))
	phase2Table.AddRow("Encryption Algos", utils.JoinStringPtr(&tunnel.Phase2.EncryptionAlgorithms, ", "))
	phase2Table.AddRow("Integrity Algos", utils.JoinStringPtr(&tunnel.Phase2.IntegrityAlgorithms, ", "))
	phase2Table.AddRow("Rekey Time", utils.PtrString(tunnel.Phase1.RekeyTime))
	phase2Table.AddRow("Dpd Action", utils.PtrString(tunnel.Phase2.DpdAction))
	phase2Table.AddRow("Start Action", utils.PtrString(tunnel.Phase2.StartAction))

	return []tables.Table{
		table,
		phase1Table,
		phase2Table,
	}
}
