package delete

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"

	iaas "github.com/stackitcloud/stackit-sdk-go/services/iaas/v2alpha1api"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/iaasalpha/client"
	iaasUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/iaasalpha/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

const (
	networkRangeIdArg = "NETWORK_RANGE_ID"

	vpcIdFlag = "vpc-id"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	VpcId          string
	NetworkRangeId string
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("delete %s", networkRangeIdArg),
		Short: "Deletes a regional network range in a VPC",
		Long:  "Deletes a regional network range in a VPC.",
		Args:  args.SingleArg(networkRangeIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Delete network range with id "xxx" in a VPC with ID "yyy"`,
				`$ stackit alpha vpc network-range delete xxx --vpc-id yyy`,
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

			vpcLabel, err := iaasUtils.GetVPCName(ctx, apiClient.DefaultAPI, model.ProjectId, model.VpcId)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get vpc name: %v", err)
				vpcLabel = model.VpcId
			}

			networkRangeLabel, err := iaasUtils.GetVPCNetworkRangePrefix(ctx, apiClient.DefaultAPI, model.ProjectId, model.VpcId, model.Region, model.NetworkRangeId)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get vpc network range prefix: %v", err)
				networkRangeLabel = model.NetworkRangeId
			} else if networkRangeLabel == "" {
				networkRangeLabel = model.NetworkRangeId
			}

			prompt := fmt.Sprintf("Are you sure you want to delete network range %q on VPC %q?", networkRangeLabel, vpcLabel)
			err = params.Printer.PromptForConfirmation(prompt)
			if err != nil {
				return err
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			err = req.Execute()
			if err != nil {
				return fmt.Errorf("delete vpc network range: %w", err)
			}

			params.Printer.Outputf("Deleted network range %q on VPC %q\n", networkRangeLabel, vpcLabel)

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

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	networkRangeId := inputArgs[0]
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		VpcId:           flags.FlagToStringValue(p, cmd, vpcIdFlag),
		NetworkRangeId:  networkRangeId,
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *iaas.APIClient) iaas.ApiDeleteVPCNetworkRangeRequest {
	return apiClient.DefaultAPI.DeleteVPCNetworkRange(ctx, model.ProjectId, model.VpcId, model.Region, model.NetworkRangeId)
}
