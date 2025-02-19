package delete

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/iaas/client"
	iaasUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/iaas/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/spinner"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/iaas"
	"github.com/stackitcloud/stackit-sdk-go/services/iaas/wait"

	"github.com/spf13/cobra"
)

const (
	networkIdArg = "NETWORK_ID"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	NetworkId string
}

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("delete %s", networkIdArg),
		Short: "Deletes a network",
		Long: fmt.Sprintf("%s\n%s\n",
			"Deletes a network.",
			"If the network is still in use, the deletion will fail",
		),
		Args: args.SingleArg(networkIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Delete network with ID "xxx"`,
				"$ stackit beta network delete xxx",
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

			networkLabel, err := iaasUtils.GetNetworkName(ctx, apiClient, model.ProjectId, model.NetworkId)
			if err != nil {
				p.Debug(print.ErrorLevel, "get network name: %v", err)
				networkLabel = model.NetworkId
			}
			if networkLabel == "" {
				networkLabel = model.NetworkId
			}

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to delete network %q?", networkLabel)
				err = p.PromptForConfirmation(prompt)
				if err != nil {
					return err
				}
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			err = req.Execute()
			if err != nil {
				return fmt.Errorf("delete network: %w", err)
			}

			// Wait for async operation, if async mode not enabled
			if !model.Async {
				s := spinner.New(p)
				s.Start("Deleting network")
				_, err = wait.DeleteNetworkWaitHandler(ctx, apiClient, model.ProjectId, model.NetworkId).WaitWithContext(ctx)
				if err != nil {
					return fmt.Errorf("wait for network deletion: %w", err)
				}
				s.Stop()
			}

			operationState := "Deleted"
			if model.Async {
				operationState = "Triggered deletion of"
			}
			p.Info("%s network %q\n", operationState, networkLabel)
			return nil
		},
	}
	return cmd
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	networkId := inputArgs[0]

	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		NetworkId:       networkId,
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

func buildRequest(ctx context.Context, model *inputModel, apiClient *iaas.APIClient) iaas.ApiDeleteNetworkRequest {
	return apiClient.DeleteNetwork(ctx, model.ProjectId, model.NetworkId)
}
