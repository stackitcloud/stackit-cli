package delete

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/iaas/client"
	iaasUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/iaas/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/iaas"
)

const (
	publicIpIdArg = "PUBLIC_IP_ID"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	PublicIpId string
}

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("delete %s", publicIpIdArg),
		Short: "Deletes a Public IP",
		Long: fmt.Sprintf("%s\n%s\n",
			"Deletes a Public IP.",
			"If the public IP is still in use, the deletion will fail",
		),
		Args: args.SingleArg(publicIpIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Delete public IP with ID "xxx"`,
				"$ stackit public-ip delete xxx",
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

			publicIpLabel, _, err := iaasUtils.GetPublicIP(ctx, apiClient, model.ProjectId, model.PublicIpId)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get public IP: %v", err)
				publicIpLabel = model.PublicIpId
			} else if publicIpLabel == "" {
				publicIpLabel = model.PublicIpId
			}

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to delete public IP %q? (This cannot be undone)", publicIpLabel)
				err = params.Printer.PromptForConfirmation(prompt)
				if err != nil {
					return err
				}
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			err = req.Execute()
			if err != nil {
				return fmt.Errorf("delete public IP: %w", err)
			}

			params.Printer.Info("Deleted public IP %q\n", publicIpLabel)
			return nil
		},
	}
	return cmd
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	publicIpId := inputArgs[0]

	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		PublicIpId:      publicIpId,
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

func buildRequest(ctx context.Context, model *inputModel, apiClient *iaas.APIClient) iaas.ApiDeletePublicIPRequest {
	return apiClient.DeletePublicIP(ctx, model.ProjectId, model.PublicIpId)
}
