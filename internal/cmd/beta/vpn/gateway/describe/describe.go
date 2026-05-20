package describe

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
		Use:   fmt.Sprintf("describe %s", gatewayIdArg),
		Short: "Shows details of a gateway",
		Long:  "Shows details of a gateway.",
		Args:  args.SingleArg(gatewayIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Describe a gateway with ID "xxx"`,
				"$ stackit beta vpn gateway describe xxx",
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

func buildRequest(ctx context.Context, model *inputModel, apiClient *vpn.APIClient) vpn.ApiGetGatewayRequest {
	return apiClient.DefaultAPI.GetGateway(ctx, model.ProjectId, vpn.Region(model.Region), model.GatewayId)
}

func outputResult(p *print.Printer, outputFormat, gatewayId, projectLabel string, gateway *vpn.GatewayResponse) error {
	return p.OutputResult(outputFormat, gateway, func() error {
		if gateway == nil {
			p.Outputf("gateway %q not found in project %q\n", gatewayId, projectLabel)
			return nil
		}

		table := tables.NewTable()
		table.SetTitle("Gateway")

		table.AddRow("ID", utils.PtrString(gateway.Id))
		table.AddSeparator()
		table.AddRow("NAME", gateway.DisplayName)
		table.AddSeparator()
		table.AddRow("STATE", utils.PtrString(gateway.State))

		if err := table.Display(p); err != nil {
			return fmt.Errorf("render tables: %w", err)
		}
		return nil
	})
}
