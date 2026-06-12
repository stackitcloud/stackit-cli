package status

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	vpn "github.com/stackitcloud/stackit-sdk-go/services/vpn/v1api"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/vpn/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
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
		Use:   fmt.Sprintf("status %s", connectionIdArg),
		Short: "Shows the status of a VPN connection",
		Long:  "Shows the status of a VPN connection.",
		Args:  args.SingleArg(connectionIdArg, nil),
		Example: examples.Build(
			examples.NewExample(
				`Show status of a VPN connection`,
				"$ stackit beta vpn connection status xxx --gateway-id yyy"),
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
				return fmt.Errorf("get VPN connection status: %w", err)
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
		return nil, &errors.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		GatewayId:       flags.FlagToStringPointer(p, cmd, gatewayIdFlag),
		ConnectionId:    connectionId,
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *vpn.APIClient) (vpn.ApiGetGatewayConnectionStatusRequest, error) {
	req := apiClient.DefaultAPI.GetGatewayConnectionStatus(ctx, model.ProjectId, model.Region, *model.GatewayId, model.ConnectionId)
	return req, nil
}

func outputResult(p *print.Printer, model *inputModel, resp *vpn.ConnectionStatusResponse) error {
	if resp == nil {
		return fmt.Errorf("status response is empty")
	}

	return p.OutputResult(model.OutputFormat, resp, func() error {
		mainTable := tables.NewTable()
		mainTable.AddRow("ID", utils.PtrString(resp.Id))
		mainTable.AddRow("Name", utils.PtrString(resp.DisplayName))
		mainTable.AddRow("Enabled", utils.PtrString(resp.Enabled))

		ts := []tables.Table{
			mainTable,
		}
		for _, tunnel := range resp.Tunnels {
			ts = append(ts, tunnelTables(&tunnel)...)
		}

		return tables.DisplayTables(p, ts)
	})
}

func tunnelTables(tunnel *vpn.TunnelStatus) []tables.Table {
	title := "Tunnel"
	if tunnel.Name != nil {
		title = string(*tunnel.Name)
	}

	table := tables.NewTable()
	table.SetTitle(title)
	table.AddRow("Established", utils.PtrString(tunnel.Established))

	res := []tables.Table{table}

	if tunnel.Phase1 != nil {
		phase1Table := tables.NewTable()
		phase1Table.SetTitle(fmt.Sprintf("%s Phase 1", title))
		phase1Table.AddRow("State", utils.PtrString(tunnel.Phase1.State))
		phase1Table.AddRow("DH Group", utils.PtrString(tunnel.Phase1.DhGroup))
		phase1Table.AddRow("Encryption Algo", utils.PtrString(tunnel.Phase1.EncryptionAlgorithm))
		phase1Table.AddRow("Integrity Algo", utils.PtrString(tunnel.Phase1.IntegrityAlgorithm))
		res = append(res, phase1Table)
	}

	if tunnel.Phase2 != nil {
		phase2Table := tables.NewTable()
		phase2Table.SetTitle(fmt.Sprintf("%s Phase 2", title))
		phase2Table.AddRow("State", utils.PtrString(tunnel.Phase2.State))
		phase2Table.AddRow("Protocol", utils.PtrString(tunnel.Phase2.Protocol))
		phase2Table.AddRow("DH Group", utils.PtrString(tunnel.Phase2.DhGroup))
		phase2Table.AddRow("Encryption Algo", utils.PtrString(tunnel.Phase2.EncryptionAlgorithm))
		phase2Table.AddRow("Integrity Algo", utils.PtrString(tunnel.Phase2.IntegrityAlgorithm))
		phase2Table.AddRow("Encap", utils.PtrString(tunnel.Phase2.Encap))
		phase2Table.AddRow("Bytes In/Out", fmt.Sprintf("%s / %s", utils.PtrString(tunnel.Phase2.BytesIn), utils.PtrString(tunnel.Phase2.BytesOut)))
		phase2Table.AddRow("Packets In/Out", fmt.Sprintf("%s / %s", utils.PtrString(tunnel.Phase2.PacketsIn), utils.PtrString(tunnel.Phase2.PacketsOut)))
		res = append(res, phase2Table)
	}

	return res
}
