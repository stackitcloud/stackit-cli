package update

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/goccy/go-yaml"
	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	cliErr "github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/iaas/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/iaasalpha"
)

const (
	organizationIdFlag = "organization-id"
	networkAreaIdFlag  = "network-area-id"
	routingTableIdFlag = "routing-table-id"
	labelFlag          = "labels"
	routeIdArg         = "ROUTE_ID_ARG"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	OrganizationId *string
	NetworkAreaId  *string
	RoutingTableId *string
	RouteId        string
	Labels         *map[string]string
}

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("update %s", routeIdArg),
		Short: "Updates a route in a routing-table",
		Long:  "Updates a route in a routing-table.",
		Args:  args.SingleArg(routeIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Updates the label(s) of a route with ID "xxx" in a routing-table ID "xxx" in organization with ID "yyy" and network-area with ID "zzz"`,
				"$ stackit beta routing-table route update xxx --labels key=value,foo=bar --routing-table-id xxx --organization-id yyy --network-area-id zzz",
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			model, err := parseInput(params.Printer, cmd, args)
			if err != nil {
				return err
			}

			// Configure API client
			apiClient, err := client.ConfigureAlphaClient(params.Printer, params.CliVersion)
			if err != nil {
				return err
			}

			// Call API
			req := apiClient.UpdateRouteOfRoutingTable(
				ctx,
				*model.OrganizationId,
				*model.NetworkAreaId,
				model.Region,
				*model.RoutingTableId,
				model.RouteId,
			)

			payload := iaasalpha.UpdateRouteOfRoutingTablePayload{
				Labels: utils.ConvertStringMapToInterfaceMap(model.Labels),
			}
			req = req.UpdateRouteOfRoutingTablePayload(payload)

			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("update route %q of routing-table %q : %w", model.RouteId, *model.RoutingTableId, err)
			}

			return outputResult(params.Printer, model.OutputFormat, *model.RoutingTableId, *model.NetworkAreaId, *resp)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), organizationIdFlag, "Organization ID")
	cmd.Flags().Var(flags.UUIDFlag(), networkAreaIdFlag, "Network-Area ID")
	cmd.Flags().Var(flags.UUIDFlag(), routingTableIdFlag, "Routing-Table ID")
	cmd.Flags().StringToString(labelFlag, nil, "Labels are key-value string pairs which can be attached to a route. A label can be provided with the format key=value and the flag can be used multiple times to provide a list of labels")

	err := flags.MarkFlagsRequired(cmd, labelFlag, organizationIdFlag, networkAreaIdFlag, routingTableIdFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)

	if len(inputArgs) == 0 {
		return nil, fmt.Errorf("at least one argument is required")
	}
	routeId := inputArgs[0]

	labels := flags.FlagToStringToStringPointer(p, cmd, labelFlag)

	if labels == nil {
		return nil, &cliErr.EmptyUpdateError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		OrganizationId:  flags.FlagToStringPointer(p, cmd, organizationIdFlag),
		NetworkAreaId:   flags.FlagToStringPointer(p, cmd, networkAreaIdFlag),
		RoutingTableId:  flags.FlagToStringPointer(p, cmd, routingTableIdFlag),
		RouteId:         routeId,
		Labels:          labels,
	}

	p.DebugInputModel(model)
	return &model, nil
}

func outputResult(p *print.Printer, outputFormat, routingTableId, networkAreaId string, route iaasalpha.Route) error {
	switch outputFormat {
	case print.JSONOutputFormat:
		details, err := json.MarshalIndent(route, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal route: %w", err)
		}
		p.Outputln(string(details))

		return nil
	case print.YAMLOutputFormat:
		details, err := yaml.MarshalWithOptions(route, yaml.IndentSequence(true), yaml.UseJSONMarshaler())
		if err != nil {
			return fmt.Errorf("marshal route: %w", err)
		}
		p.Outputln(string(details))

		return nil
	default:
		p.Outputf("Updated route %q for routing-table %q in network-area %q.", *route.Id, routingTableId, networkAreaId)
		return nil
	}
}
