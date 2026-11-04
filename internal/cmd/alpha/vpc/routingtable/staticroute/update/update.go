package update

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	iaas "github.com/stackitcloud/stackit-sdk-go/services/iaas/v2alpha1api"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	cliErr "github.com/stackitcloud/stackit-cli/internal/pkg/errors"
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
	labelFlag          = "labels"
	routeIdArg         = "ROUTE_ID"
	routingTableIdFlag = "routing-table-id"
	vpcIdFlag          = "vpc-id"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	Labels         map[string]any
	RouteId        string
	RoutingTableId string
	VpcId          string
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("update %s", routeIdArg),
		Short: "Updates a static-route in a VPC routing-table",
		Long:  "Updates a static-route in a VPC routing-table.",
		Args:  args.SingleArg(routeIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Updates the label(s) of a static-route with ID "xxx" in a routing-table ID "yyy" in VPC "zzz"`,
				"$ stackit alpha vpc routing-table static-route update xxx --routing-table-id yyy --vpc-id zzz --labels key=value,foo=bar",
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

			prompt := fmt.Sprintf("Are you sure you want to update static-route %q for routing-table %q?", model.RouteId, routingTableLabel)
			err = params.Printer.PromptForConfirmation(prompt)
			if err != nil {
				return err
			}

			req := buildRequest(ctx, model, apiClient.DefaultAPI)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("update static-route %q of routing-table %q: %w", model.RouteId, model.RoutingTableId, err)
			}

			return outputResult(params.Printer, model.OutputFormat, model.RoutingTableId, model.VpcId, resp)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), vpcIdFlag, "VPC ID")
	cmd.Flags().Var(flags.UUIDFlag(), routingTableIdFlag, "Routing-Table ID")
	cmd.Flags().StringToString(labelFlag, nil, "Labels are key-value string pairs which can be attached to a static route. E.g. '--labels key1=value1,key2=value2,...'")

	err := flags.MarkFlagsRequired(cmd, labelFlag, vpcIdFlag, routingTableIdFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)

	routeId := inputArgs[0]

	labels := flags.FlagToStringToAny(p, cmd, labelFlag)

	if labels == nil {
		return nil, &cliErr.EmptyUpdateError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		Labels:          labels,
		RouteId:         routeId,
		RoutingTableId:  flags.FlagToStringValue(p, cmd, routingTableIdFlag),
		VpcId:           flags.FlagToStringValue(p, cmd, vpcIdFlag),
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient iaas.DefaultAPI) iaas.ApiUpdateVPCStaticRouteRequest {
	req := apiClient.UpdateVPCStaticRoute(ctx, model.ProjectId, model.VpcId, model.Region, model.RoutingTableId, model.RouteId)

	payload := iaas.UpdateVPCStaticRoutePayload{
		Labels: model.Labels,
	}

	return req.UpdateVPCStaticRoutePayload(payload)
}

func outputResult(p *print.Printer, outputFormat, routingTableId, vpcId string, route *iaas.Route) error {
	return p.OutputResult(outputFormat, route, func() error {
		if route == nil {
			return fmt.Errorf("update static-route response is empty")
		}

		if route.Id == nil || *route.Id == "" {
			return fmt.Errorf("update static-route response has empty id")
		}

		p.Outputf("Updated static-route %q for routing-table %q in VPC %q.\n", *route.Id, routingTableId, vpcId)
		return nil
	})
}
