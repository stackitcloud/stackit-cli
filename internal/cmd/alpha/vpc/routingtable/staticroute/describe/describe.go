package describe

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
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
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
		Use:   fmt.Sprintf("describe %s", routeIdArg),
		Short: "Describes a static-route within a VPC routing-table",
		Long:  "Describes a static-route within a VPC routing-table.",
		Args:  args.SingleArg(routeIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Describe a static-route within a VPC routing-table`,
				`$ stackit alpha vpc routing-table static-route describe xxx --routing-table-id yyy --vpc-id zzz`,
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

			// Call API
			request := buildRequest(ctx, model, apiClient.DefaultAPI)
			response, err := request.Execute()
			if err != nil {
				return fmt.Errorf("describe static-route: %w", err)
			}

			return outputResult(params.Printer, model.OutputFormat, response)
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

func buildRequest(ctx context.Context, model *inputModel, apiClient iaas.DefaultAPI) iaas.ApiGetVPCStaticRouteRequest {
	return apiClient.GetVPCStaticRoute(ctx, model.ProjectId, model.VpcId, model.Region, model.RoutingTableId, model.RouteId)
}

func outputResult(p *print.Printer, outputFormat string, route *iaas.Route) error {
	if route == nil {
		return fmt.Errorf("describe static-route response is empty")
	}

	return p.OutputResult(outputFormat, route, func() error {
		routeDetails := iaasUtils.ExtractRouteDetails(route)

		table := tables.NewTable()

		table.AddRow("ID", utils.PtrString(route.Id))
		table.AddSeparator()

		table.AddRow("DESTINATION TYPE", routeDetails.DestType)
		table.AddSeparator()

		table.AddRow("DESTINATION VALUE", routeDetails.DestValue)
		table.AddSeparator()

		table.AddRow("NEXTHOP TYPE", routeDetails.NextHopType)
		table.AddSeparator()

		table.AddRow("NEXTHOP VALUE", routeDetails.NextHopValue)
		table.AddSeparator()

		table.AddRow("LABELS", routeDetails.Labels)
		table.AddSeparator()

		table.AddRow("CREATED AT", routeDetails.CreatedAt)
		table.AddSeparator()

		table.AddRow("UPDATED AT", routeDetails.UpdatedAt)
		table.AddSeparator()

		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}

		return nil
	})
}
