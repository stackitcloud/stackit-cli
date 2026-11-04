package delete

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	iaas "github.com/stackitcloud/stackit-sdk-go/services/iaas/v2alpha1api"
	"github.com/stackitcloud/stackit-sdk-go/services/iaas/v2alpha1api/wait"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/iaasalpha/client"
	iaasUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/iaasalpha/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/spinner"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
)

const vpcIdFlag = "vpc-id"

type inputModel struct {
	*globalflags.GlobalFlagModel
	VpcId string
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Deletes a regional configuration for a VPC",
		Long:  "Deletes a regional configuration for a VPC.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Delete the regional configuration "eu02" for a VPC with ID "xxx"`,
				`$ stackit alpha vpc region delete --vpc-id xxx --region eu02`,
			),
			examples.NewExample(
				`Delete a regional configuration using the configured project and region`,
				`$ stackit alpha vpc region delete --vpc-id xxx`,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			model, err := parseInput(params.Printer, cmd, args)
			if err != nil {
				return err
			}

			apiClient, err := client.ConfigureClient(params.Printer, params.CliVersion)
			if err != nil {
				return err
			}

			vpcLabel, err := iaasUtils.GetVPCName(ctx, apiClient.DefaultAPI, model.ProjectId, model.VpcId)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get vpc name: %v", err)
			}
			if vpcLabel == "" {
				vpcLabel = model.VpcId
			}

			prompt := fmt.Sprintf("Are you sure you want to delete the regional configuration %q for VPC %q?", model.Region, vpcLabel)
			if err := params.Printer.PromptForConfirmation(prompt); err != nil {
				return err
			}

			req := buildRequest(ctx, model, apiClient.DefaultAPI)

			if err := req.Execute(); err != nil {
				return fmt.Errorf("delete vpc region: %w", err)
			}

			if !model.Async {
				err := spinner.Run(params.Printer, "Delete vpc region", func() error {
					_, err := wait.DeleteVPCRegionWaitHandler(ctx, apiClient.DefaultAPI, model.ProjectId, model.VpcId, model.Region).WaitWithContext(ctx)
					return err
				})
				if err != nil {
					return fmt.Errorf("wait for vpc region deletion: %w", err)
				}
			}

			operationState := "Deleted"
			if model.Async {
				operationState = "Triggered deletion of"
			}

			params.Printer.Outputf("%s region configuration for VPC %q.\nRegion: %s\n", operationState, vpcLabel, model.Region)

			return nil
		},
	}
	configureFlags(cmd)

	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), vpcIdFlag, "VPC ID")

	cobra.CheckErr(flags.MarkFlagsRequired(cmd, vpcIdFlag))
}

func parseInput(p *print.Printer, cmd *cobra.Command, _ []string) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}
	if globalFlags.Region == "" {
		return nil, &errors.RegionError{}
	}

	model := &inputModel{
		GlobalFlagModel: globalFlags,
		VpcId:           flags.FlagToStringValue(p, cmd, vpcIdFlag),
	}

	p.DebugInputModel(model)
	return model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient iaas.DefaultAPI) iaas.ApiDeleteVPCRegionRequest {
	return apiClient.DeleteVPCRegion(ctx, model.ProjectId, model.VpcId, model.Region)
}
