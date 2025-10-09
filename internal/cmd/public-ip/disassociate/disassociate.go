package disassociate

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	cliErr "github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/iaas/client"
	iaasUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/iaas/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/iaas"

	"github.com/spf13/cobra"
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
		Use:   fmt.Sprintf("disassociate %s", publicIpIdArg),
		Short: "Disassociates a Public IP from a network interface or a virtual IP",
		Long:  "Disassociates a Public IP from a network interface or a virtual IP.",
		Args:  args.SingleArg(publicIpIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Disassociate public IP with ID "xxx" from a resource (network interface or virtual IP)`,
				`$ stackit public-ip disassociate xxx`,
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

			publicIpLabel, associatedResourceId, err := iaasUtils.GetPublicIP(ctx, apiClient, model.ProjectId, model.PublicIpId)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get public IP: %v", err)
				publicIpLabel = model.PublicIpId
			} else if publicIpLabel == "" {
				publicIpLabel = model.PublicIpId
			}

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to disassociate public IP %q from the associated resource %q?", publicIpLabel, associatedResourceId)
				err = params.Printer.PromptForConfirmation(prompt)
				if err != nil {
					return err
				}
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			_, err = req.Execute()
			if err != nil {
				return fmt.Errorf("disassociate public IP: %w", err)
			}

			params.Printer.Outputf("Disassociated public IP %q from the associated resource %q.\n", publicIpLabel, associatedResourceId)
			return nil
		},
	}
	return cmd
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	publicIpId := inputArgs[0]

	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &cliErr.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		PublicIpId:      publicIpId,
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *iaas.APIClient) iaas.ApiUpdatePublicIPRequest {
	req := apiClient.UpdatePublicIP(ctx, model.ProjectId, model.PublicIpId)

	payload := iaas.UpdatePublicIPPayload{
		NetworkInterface: iaas.NewNullableString(nil),
	}

	return req.UpdatePublicIPPayload(payload)
}
