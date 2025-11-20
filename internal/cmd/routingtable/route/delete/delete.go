package delete

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/iaas/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

const (
	networkAreaIdFlag  = "network-area-id"
	organizationIdFlag = "organization-id"
	routeIdArg         = "ROUTE_ID_ARG"
	routingTableIdFlag = "routing-table-id"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	NetworkAreaId  string
	OrganizationId string
	RouteID        string
	RoutingTableId string
}

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("delete %s", routingTableIdFlag),
		Short: "Deletes a route within a routing-table",
		Long:  "Deletes a route within a routing-table",
		Args:  args.SingleArg(routeIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Deletes a route within a routing-table`,
				`$ stackit routing-table route delete xxx --routing-table-id xxx --organization-id yyy --network-area-id zzz`,
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
				prompt := fmt.Sprintf("Are you sure you want to delete the route %q in routing-table %q for network-area-id %q?", model.RouteID, model.RoutingTableId, model.OrganizationId)
				err = params.Printer.PromptForConfirmation(prompt)
				if err != nil {
					return err
				}
			}

			// Call API
			req := apiClient.DeleteRouteFromRoutingTable(
				ctx,
				model.OrganizationId,
				model.NetworkAreaId,
				model.Region,
				model.RoutingTableId,
				model.RouteID,
			)
			err = req.Execute()
			if err != nil {
				return fmt.Errorf("delete route from routing-table: %w", err)
			}

			params.Printer.Outputf("Route %q from routing-table %q deleted.", model.RouteID, model.RoutingTableId)
			return nil
		},
	}

	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), networkAreaIdFlag, "Network-Area ID")
	cmd.Flags().Var(flags.UUIDFlag(), organizationIdFlag, "Organization ID")
	cmd.Flags().Var(flags.UUIDFlag(), routingTableIdFlag, "Routing-Table ID")

	err := flags.MarkFlagsRequired(cmd, organizationIdFlag, networkAreaIdFlag, routingTableIdFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)

	if len(inputArgs) == 0 {
		return nil, fmt.Errorf("at least one argument is required")
	}
	routeId := inputArgs[0]

	model := inputModel{
		GlobalFlagModel: globalFlags,
		NetworkAreaId:   flags.FlagToStringValue(p, cmd, networkAreaIdFlag),
		OrganizationId:  flags.FlagToStringValue(p, cmd, organizationIdFlag),
		RouteID:         routeId,
		RoutingTableId:  flags.FlagToStringValue(p, cmd, routingTableIdFlag),
	}

	p.DebugInputModel(model)
	return &model, nil
}
