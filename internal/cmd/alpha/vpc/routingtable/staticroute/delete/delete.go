package delete

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	iaas "github.com/stackitcloud/stackit-sdk-go/services/iaas/v2alpha1api"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/iaasalpha/client"
	iaasUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/iaasalpha/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

const (
	routeIdArg         = "ROUTE_ID"
	routingTableIdFlag = "routing-table-id"
	vpcIdFlag          = "vpc-id"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	RouteId        string
	RoutingTableId string
	VpcId          string
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("delete %s", routeIdArg),
		Short: "Deletes a static-route within a VPC routing-table",
		Long:  "Deletes a static-route within a VPC routing-table.",
		Args:  args.SingleArg(routeIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Deletes a static-route within a VPC routing-table`,
				`$ stackit alpha vpc routing-table static-route delete xxx --routing-table-id yyy --vpc-id zzz`,
			),
		),
		RunE: func(cmd *cobra.Command, inputArgs []string) error {
			ctx := context.Background()

			model, err := parseInput(params.Printer, cmd, inputArgs)
			if err != nil {
				return err
			}

			// Configure API client
			apiClient, err := client.ConfigureClient(params.Printer, params.CliVersion)
			if err != nil {
				return err
			}

			routingTableLabel, err := iaasUtils.GetVPCRoutingTableName(ctx, apiClient.DefaultAPI, model.ProjectId, model.VpcId, model.Region, model.RoutingTableId)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get vpc routing-table name: %v", err)

				routingTableLabel = model.RoutingTableId
			} else if routingTableLabel == "" {
				routingTableLabel = model.RoutingTableId
			}

			prompt := fmt.Sprintf("Are you sure you want to delete static-route %q in routing-table %q?", model.RouteId, routingTableLabel)

			err = params.Printer.PromptForConfirmation(prompt)
			if err != nil {
				return err
			}

			// Call API
			req := buildRequest(ctx, model, apiClient.DefaultAPI)

			err = req.Execute()
			if err != nil {
				return fmt.Errorf("delete static-route from routing-table: %w", err)
			}

			params.Printer.Outputf("Deleted static-route %q from routing-table %q.\n", model.RouteId, model.RoutingTableId)

			return nil
		},
	}

	configureFlags(cmd)

	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), vpcIdFlag, "VPC ID")
	cmd.Flags().Var(flags.UUIDFlag(), routingTableIdFlag, "Routing-Table ID")

	err := flags.MarkFlagsRequired(cmd, vpcIdFlag, routingTableIdFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)

	routeId := inputArgs[0]

	model := inputModel{
		GlobalFlagModel: globalFlags,
		RouteId:         routeId,
		RoutingTableId:  flags.FlagToStringValue(p, cmd, routingTableIdFlag),
		VpcId:           flags.FlagToStringValue(p, cmd, vpcIdFlag),
	}

	p.DebugInputModel(model)

	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient iaas.DefaultAPI) iaas.ApiDeleteVPCStaticRouteRequest {
	return apiClient.DeleteVPCStaticRoute(ctx, model.ProjectId, model.VpcId, model.Region, model.RoutingTableId, model.RouteId)
}
