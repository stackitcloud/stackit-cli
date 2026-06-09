package list

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
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	vpn "github.com/stackitcloud/stackit-sdk-go/services/vpn/v1api"
)

const (
	gatewayIdFlag = "gateway-id"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	GatewayId *string
}

func NewCmd(p *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Lists all VPN connections of a gateway",
		Long:  "Lists all VPN connections of a gateway.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`List all VPN connections of a gateway`,
				"$ stackit beta vpn connection list --gateway-id xxx"),
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

			// Call API
			req, err := buildRequest(ctx, model, apiClient)
			if err != nil {
				return err
			}
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("list VPN connections: %w", err)
			}

			projectLabel, err := projectname.GetProjectName(ctx, p.Printer, p.CliVersion, cmd)
			if err != nil {
				p.Printer.Debug(print.ErrorLevel, "get project name: %v", err)
				projectLabel = model.ProjectId
			}

			return outputResult(p.Printer, model, projectLabel, resp)
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

func parseInput(p *print.Printer, cmd *cobra.Command) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &cliErr.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		GatewayId:       flags.FlagToStringPointer(p, cmd, gatewayIdFlag),
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *vpn.APIClient) (vpn.ApiListGatewayConnectionsRequest, error) {
	req := apiClient.DefaultAPI.ListGatewayConnections(ctx, model.ProjectId, model.Region, *model.GatewayId)
	return req, nil
}

func outputResult(p *print.Printer, model *inputModel, projectLabel string, resp *vpn.ConnectionList) error {
	if resp == nil || resp.Connections == nil {
		return fmt.Errorf("list connections response is empty")
	}

	return p.OutputResult(model.OutputFormat, resp.Connections, func() error {
		table := tables.NewTable()
		table.SetHeader("ID", "NAME", "ENABLED", "LABELS")
		for _, c := range resp.Connections {
			id := utils.PtrString(c.Id)
			name := c.DisplayName
			enabled := utils.PtrString(c.Enabled)
			var labels string
			if c.Labels != nil {
				labels = utils.JoinStringMap(*c.Labels, "=", ", ")
			}
			table.AddRow(id, name, enabled, labels)
		}
		p.Outputln(table.Render())
		return nil
	})
}
