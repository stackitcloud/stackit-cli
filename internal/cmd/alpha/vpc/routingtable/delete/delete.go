package delete

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	iaas "github.com/stackitcloud/stackit-sdk-go/services/iaas/v2alpha1api"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/iaasalpha/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

const (
	vpcIdFlag         = "vpc-id"
	routingTableIdArg = "ROUTING_TABLE_ID"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	VpcId          string
	RoutingTableId string
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("delete %s", routingTableIdArg),
		Short: "Deletes a VPC routing-table",
		Long:  "Deletes a routing-table for a VPC in a region.",
		Args:  args.SingleArg(routingTableIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Delete routing-table "xxx" for VPC "yyy" in region "eu02"`,
				`$ stackit alpha vpc routing-table delete xxx --vpc-id yyy --region eu02`,
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

			prompt := fmt.Sprintf("Are you sure you want to delete the routing-table %q for VPC %q in region %q?", model.RoutingTableId, model.VpcId, model.Region)
			if err := params.Printer.PromptForConfirmation(prompt); err != nil {
				return err
			}

			req := buildRequest(ctx, model, apiClient.DefaultAPI)
			if err := req.Execute(); err != nil {
				return fmt.Errorf("delete vpc routing-table: %w", err)
			}

			params.Printer.Outputf("Deleted routing-table %q.\n", model.RoutingTableId)
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
		RoutingTableId:  inputArgs[0],
	}

	p.DebugInputModel(model)
	return model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient iaas.DefaultAPI) iaas.ApiDeleteVPCRoutingTableRequest {
	return apiClient.DeleteVPCRoutingTable(ctx, model.ProjectId, model.VpcId, model.Region, model.RoutingTableId)
}
