package delete

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/types"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/iaas/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/iaas"
)

const (
	nicIdArg = "NIC_ID"

	networkIdFlag = "network-id"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	NetworkId string
	NicId     string
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("delete %s", nicIdArg),
		Short: "Deletes a network interface",
		Long:  "Deletes a network interface.",
		Args:  args.SingleArg(nicIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Delete network interface with nic id "xxx" and network ID "yyy"`,
				`$ stackit network-interface delete xxx --network-id yyy`,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			model, err := parseInput(params.Printer, cmd, args)
			if err != nil {
				return err
			}

			// Configure API client
			apiClient, err := client.ConfigureClient(params.Printer, params.CliVersion)
			if err != nil {
				return err
			}

			prompt := fmt.Sprintf("Are you sure you want to delete the network interface %q?  (This cannot be undone)", model.NicId)
			err = params.Printer.PromptForConfirmation(prompt)
			if err != nil {
				return err
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			err = req.Execute()
			if err != nil {
				return fmt.Errorf("delete network interface: %w", err)
			}

			params.Printer.Info("Deleted network interface %q\n", model.NicId)

			return nil
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), networkIdFlag, "Network ID")

	err := flags.MarkFlagsRequired(cmd, networkIdFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	nicId := inputArgs[0]
	globalFlags := globalflags.Parse(p, cmd)

	model := inputModel{
		GlobalFlagModel: globalFlags,
		NetworkId:       flags.FlagToStringValue(p, cmd, networkIdFlag),
		NicId:           nicId,
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *iaas.APIClient) iaas.ApiDeleteNicRequest {
	req := apiClient.DeleteNic(ctx, model.ProjectId, model.Region, model.NetworkId, model.NicId)
	return req
}
