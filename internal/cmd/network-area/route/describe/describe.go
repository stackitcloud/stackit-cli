package describe

import (
	"context"
	"fmt"
	"strings"

	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/iaas/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/iaas"

	"github.com/spf13/cobra"
)

const (
	routeIdArg = "ROUTE_ID"

	organizationIdFlag = "organization-id"
	networkAreaIdFlag  = "network-area-id"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	OrganizationId *string
	NetworkAreaId  *string
	RouteId        string
}

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("describe %s", routeIdArg),
		Short: "Shows details of a static route in a STACKIT Network Area (SNA)",
		Long:  "Shows details of a static route in a STACKIT Network Area (SNA).",
		Args:  args.SingleArg(routeIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Show details of a static route with ID "xxx" in a STACKIT Network Area with ID "yyy" in organization with ID "zzz"`,
				`$ stackit network-area route describe xxx --network-area-id yyy --organization-id zzz`,
			),
			examples.NewExample(
				`Show details of a static route with ID "xxx" in a STACKIT Network Area with ID "yyy" in organization with ID "zzz" in JSON format`,
				`$ stackit network-area route describe xxx --network-area-id yyy --organization-id zzz --output-format json`,
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

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("describe static route: %w", err)
			}

			return outputResult(params.Printer, model.OutputFormat, *resp)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), organizationIdFlag, "Organization ID")
	cmd.Flags().Var(flags.UUIDFlag(), networkAreaIdFlag, "STACKIT Network Area ID")

	err := flags.MarkFlagsRequired(cmd, organizationIdFlag, networkAreaIdFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	routeId := inputArgs[0]
	globalFlags := globalflags.Parse(p, cmd)

	model := inputModel{
		GlobalFlagModel: globalFlags,
		OrganizationId:  flags.FlagToStringPointer(p, cmd, organizationIdFlag),
		NetworkAreaId:   flags.FlagToStringPointer(p, cmd, networkAreaIdFlag),
		RouteId:         routeId,
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *iaas.APIClient) iaas.ApiGetNetworkAreaRouteRequest {
	req := apiClient.GetNetworkAreaRoute(ctx, *model.OrganizationId, *model.NetworkAreaId, model.Region, model.RouteId)
	return req
}

func outputResult(p *print.Printer, outputFormat string, route iaas.Route) error {
	return p.OutputResult(outputFormat, route, func() error {
		table := tables.NewTable()
		table.AddRow("ID", utils.PtrString(route.Id))
		table.AddSeparator()
		if destination := route.Destination; destination != nil {
			if destination.DestinationCIDRv4 != nil {
				table.AddRow("DESTINATION TYPE", utils.PtrString(destination.DestinationCIDRv4.Type))
				table.AddSeparator()
				table.AddRow("DESTINATION", utils.PtrString(destination.DestinationCIDRv4.Value))
				table.AddSeparator()
			} else if destination.DestinationCIDRv6 != nil {
				table.AddRow("DESTINATION TYPE", utils.PtrString(destination.DestinationCIDRv6.Type))
				table.AddSeparator()
				table.AddRow("DESTINATION", utils.PtrString(destination.DestinationCIDRv6.Value))
				table.AddSeparator()
			}
		}
		if nexthop := route.Nexthop; nexthop != nil {
			if nexthop.NexthopIPv4 != nil {
				table.AddRow("NEXTHOP", utils.PtrString(nexthop.NexthopIPv4.Value))
				table.AddSeparator()
				table.AddRow("NEXTHOP TYPE", utils.PtrString(nexthop.NexthopIPv4.Type))
				table.AddSeparator()
			} else if nexthop.NexthopIPv6 != nil {
				table.AddRow("NEXTHOP", utils.PtrString(nexthop.NexthopIPv6.Value))
				table.AddSeparator()
				table.AddRow("NEXTHOP TYPE", utils.PtrString(nexthop.NexthopIPv6.Type))
				table.AddSeparator()
			} else if nexthop.NexthopBlackhole != nil {
				table.AddRow("NEXTHOP TYPE", utils.PtrString(nexthop.NexthopBlackhole.Type))
				table.AddSeparator()
			} else if nexthop.NexthopInternet != nil {
				table.AddRow("NEXTHOP TYPE", utils.PtrString(nexthop.NexthopInternet.Type))
				table.AddSeparator()
			}
		}
		if route.Labels != nil && len(*route.Labels) > 0 {
			labels := []string{}
			for key, value := range *route.Labels {
				labels = append(labels, fmt.Sprintf("%s: %s", key, value))
			}
			table.AddSeparator()
			table.AddRow("LABELS", strings.Join(labels, "\n"))
		}

		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}
		return nil
	})
}
