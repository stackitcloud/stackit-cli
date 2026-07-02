package delete

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	vpn "github.com/stackitcloud/stackit-sdk-go/services/vpn/v1api"
	"github.com/stackitcloud/stackit-sdk-go/services/vpn/v1api/wait"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/vpn/client"
	vpnUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/vpn/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/spinner"
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
		Use:   fmt.Sprintf("delete %s", gatewayIdArg),
		Short: "Deletes a vpn gateway",
		Long:  "Deletes a vpn gateway.",
		Args:  args.SingleArg(gatewayIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Delete a vpn gateway with the ID "xxx"`,
				"$ stackit beta vpn gateway delete xxx",
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

			prompt := fmt.Sprintf("Are you sure you want to delete the vpn gateway %q? (This cannot be undone)", gatewayLabel)
			err = params.Printer.PromptForConfirmation(prompt)
			if err != nil {
				return err
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			err = req.Execute()
			if err != nil {
				return fmt.Errorf("delete vpn gateway: %w", err)
			}

			// Wait for async operation, if async mode not enabled
			if !model.Async {
				err := spinner.Run(params.Printer, "Deleting gateway", func() error {
					_, err = wait.DeleteGatewayWaitHandler(ctx, apiClient.DefaultAPI, model.ProjectId, model.Region, model.GatewayId).WaitWithContext(ctx)
					return err
				})
				if err != nil {
					return fmt.Errorf("waiting for gateway deletion: %w", err)
				}
			}

			operation := "Deleted"
			if model.Async {
				operation = "Triggered deletion of"
			}

			params.Printer.Outputf("%s gateway %q\n", operation, gatewayLabel)
			return nil
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

func buildRequest(ctx context.Context, model *inputModel, apiClient *vpn.APIClient) vpn.ApiDeleteGatewayRequest {
	return apiClient.DefaultAPI.DeleteGateway(ctx, model.ProjectId, model.Region, model.GatewayId)
}
