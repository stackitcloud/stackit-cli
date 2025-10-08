package associate

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	cliErr "github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
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

	associatedResourceIdFlag = "associated-resource-id"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	PublicIpId           string
	AssociatedResourceId *string
}

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("associate %s", publicIpIdArg),
		Short: "Associates a Public IP with a network interface or a virtual IP",
		Long:  "Associates a Public IP with a network interface or a virtual IP.",
		Args:  args.SingleArg(publicIpIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Associate public IP with ID "xxx" to a resource (network interface or virtual IP) with ID "yyy"`,
				`$ stackit public-ip associate xxx --associated-resource-id yyy`,
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
				prompt := fmt.Sprintf("Are you sure you want to associate public IP %q with resource %v?", publicIpLabel, *model.AssociatedResourceId)
				err = params.Printer.PromptForConfirmation(prompt)
				if err != nil {
					return err
				}
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("associate public IP: %w", err)
			}

			params.Printer.Outputf("Associated public IP %q with resource %v.\n", publicIpLabel, utils.PtrString(resp.GetNetworkInterface()))
			return nil
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), associatedResourceIdFlag, "Associates the public IP with a network interface or virtual IP (ID)")

	err := flags.MarkFlagsRequired(cmd, associatedResourceIdFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	publicIpId := inputArgs[0]

	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &cliErr.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel:      globalFlags,
		AssociatedResourceId: flags.FlagToStringPointer(p, cmd, associatedResourceIdFlag),
		PublicIpId:           publicIpId,
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *iaas.APIClient) iaas.ApiUpdatePublicIPRequest {
	req := apiClient.UpdatePublicIP(ctx, model.ProjectId, model.PublicIpId)

	payload := iaas.UpdatePublicIPPayload{
		NetworkInterface: iaas.NewNullableString(model.AssociatedResourceId),
	}

	return req.UpdatePublicIPPayload(payload)
}
