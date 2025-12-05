package update

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	cliErr "github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/iaas/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/iaas"
)

const (
	labelFlag          = "labels"
	networkAreaIdFlag  = "network-area-id"
	organizationIdFlag = "organization-id"
	routeIdArg         = "ROUTE_ID"
	routingTableIdFlag = "routing-table-id"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	Labels         *map[string]string
	NetworkAreaId  string
	OrganizationId string
	RouteId        string
	RoutingTableId string
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("update %s", routeIdArg),
		Short: "Updates a route in a routing-table",
		Long:  "Updates a route in a routing-table.",
		Args:  args.SingleArg(routeIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Updates the label(s) of a route with ID "xxx" in a routing-table ID "xxx" in organization with ID "yyy" and network-area with ID "zzz"`,
				"$ stackit routing-table route update xxx --labels key=value,foo=bar --routing-table-id xxx --organization-id yyy --network-area-id zzz",
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

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to update route %q for routing-table with id %q?", model.RouteId, model.RoutingTableId)
				if err := params.Printer.PromptForConfirmation(prompt); err != nil {
					return err
				}
			}

			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("update route %q of routing-table %q : %w", model.RouteId, model.RoutingTableId, err)
			}

			return outputResult(params.Printer, model.OutputFormat, model.RoutingTableId, model.NetworkAreaId, resp)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().StringToString(labelFlag, nil, "Labels are key-value string pairs which can be attached to a route. A label can be provided with the format key=value and the flag can be used multiple times to provide a list of labels")
	cmd.Flags().Var(flags.UUIDFlag(), networkAreaIdFlag, "Network-Area ID")
	cmd.Flags().Var(flags.UUIDFlag(), organizationIdFlag, "Organization ID")
	cmd.Flags().Var(flags.UUIDFlag(), routingTableIdFlag, "Routing-Table ID")

	err := flags.MarkFlagsRequired(cmd, labelFlag, organizationIdFlag, networkAreaIdFlag, routingTableIdFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)

	routeId := inputArgs[0]

	labels := flags.FlagToStringToStringPointer(p, cmd, labelFlag)

	if labels == nil {
		return nil, &cliErr.EmptyUpdateError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		Labels:          labels,
		NetworkAreaId:   flags.FlagToStringValue(p, cmd, networkAreaIdFlag),
		OrganizationId:  flags.FlagToStringValue(p, cmd, organizationIdFlag),
		RouteId:         routeId,
		RoutingTableId:  flags.FlagToStringValue(p, cmd, routingTableIdFlag),
	}

	p.DebugInputModel(model)
	return &model, nil
}

func outputResult(p *print.Printer, outputFormat, routingTableId, networkAreaId string, route *iaas.Route) error {
	if route == nil {
		return fmt.Errorf("update route response is empty")
	}

	if route.Id == nil || *route.Id == "" {
		return fmt.Errorf("update route response has empty id")
	}

	return p.OutputResult(outputFormat, route, func() error {
		if route == nil {
			return fmt.Errorf("update route response is empty")
		}

		if route.Id == nil || *route.Id == "" {
			return fmt.Errorf("update route response has empty id")
		}

		p.Outputf("Updated route %q for routing-table %q in network-area %q.", *route.Id, routingTableId, networkAreaId)
		return nil
	})
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *iaas.APIClient) iaas.ApiUpdateRouteOfRoutingTableRequest {
	req := apiClient.UpdateRouteOfRoutingTable(
		ctx,
		model.OrganizationId,
		model.NetworkAreaId,
		model.Region,
		model.RoutingTableId,
		model.RouteId,
	)

	payload := iaas.UpdateRouteOfRoutingTablePayload{
		Labels: utils.ConvertStringMapToInterfaceMap(model.Labels),
	}

	return req.UpdateRouteOfRoutingTablePayload(payload)
}
