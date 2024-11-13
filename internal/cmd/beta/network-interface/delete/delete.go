package delete

import (
	"context"
	"fmt"

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
	NetworkId *string
	NicId     string
}

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Deletes a network interface",
		Long:  "Deletes a network interface.",
		Args:  args.SingleArg(nicIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Delete network interface with nic id "xxx" and network ID "yyy"`,
				`$ stackit beta network-interface delete xxx --network-id yyy`,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			model, err := parseInput(p, cmd, args)
			if err != nil {
				return err
			}

			// Configure API client
			apiClient, err := client.ConfigureClient(p)
			if err != nil {
				return err
			}

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to delete the network interface %q?  (This cannot be undone)", model.NicId)
				err = p.PromptForConfirmation(prompt)
				if err != nil {
					return err
				}
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			err = req.Execute()
			if err != nil {
				return fmt.Errorf("delete network interface: %w", err)
			}

			p.Info("Deleted network interface %q\n", model.NicId)

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
		NetworkId:       flags.FlagToStringPointer(p, cmd, networkIdFlag),
		NicId:           nicId,
	}

	if p.IsVerbosityDebug() {
		modelStr, err := print.BuildDebugStrFromInputModel(model)
		if err != nil {
			p.Debug(print.ErrorLevel, "convert model to string for debugging: %v", err)
		} else {
			p.Debug(print.DebugLevel, "parsed input values: %s", modelStr)
		}
	}

	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *iaas.APIClient) iaas.ApiDeleteNicRequest {
	req := apiClient.DeleteNic(ctx, model.ProjectId, *model.NetworkId, model.NicId)
	return req
}
