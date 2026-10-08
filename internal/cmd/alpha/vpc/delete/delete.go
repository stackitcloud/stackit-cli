package delete

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	iaas "github.com/stackitcloud/stackit-sdk-go/services/iaas/v2alpha1api"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/iaasalpha/client"
	iaasUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/iaasalpha/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

const (
	vpcIdArg = "VPC_ID"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	VpcId string
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("delete %s", vpcIdArg),
		Short: "Deletes a VPC",
		Long:  "Deletes a VPC.",
		Args:  args.SingleArg(vpcIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Delete a vpc with ID "xxx"`,
				"$ stackit alpha vpc delete xxx"),
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

			vpcLabel, err := iaasUtils.GetVPCName(ctx, apiClient.DefaultAPI, model.ProjectId, model.VpcId)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get VPC name: %v", err)
				vpcLabel = model.VpcId
			} else if vpcLabel == "" {
				vpcLabel = model.VpcId
			}

			prompt := fmt.Sprintf("Are you sure you want to delete vpc %q?", vpcLabel)
			err = params.Printer.PromptForConfirmation(prompt)
			if err != nil {
				return err
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			err = req.Execute()
			if err != nil {
				return fmt.Errorf("delete vpc: %w", err)
			}

			params.Printer.Outputf("Deleted vpc %q\n", vpcLabel)
			return nil
		},
	}
	return cmd
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	vpcId := inputArgs[0]

	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		VpcId:           vpcId,
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *iaas.APIClient) iaas.ApiDeleteVPCRequest {
	return apiClient.DefaultAPI.DeleteVPC(ctx, model.ProjectId, model.VpcId)
}
