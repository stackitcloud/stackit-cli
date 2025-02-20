package describe

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/goccy/go-yaml"
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

func NewCmd(p *print.Printer) *cobra.Command {
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
			model, err := parseInput(p, cmd, args)
			if err != nil {
				return err
			}

			// Configure API client
			apiClient, err := client.ConfigureClient(p)
			if err != nil {
				return err
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("describe static route: %w", err)
			}

			return outputResult(p, model.OutputFormat, *resp)
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

	if p.IsVerbosityDebug() {
		modelStr, err := print.BuildDebugStrFromInputModel(model)
		if err != nil {
			p.Debug(print.ErrorLevel, "convert model to string for debugging: %v", err)
		} else {
			p.Debug(print.DebugLevel, "parsed input values: %s", modelStr)
		}
	}

	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *iaas.APIClient) iaas.ApiGetNetworkAreaRouteRequest {
	req := apiClient.GetNetworkAreaRoute(ctx, *model.OrganizationId, *model.NetworkAreaId, model.RouteId)
	return req
}

func outputResult(p *print.Printer, outputFormat string, route iaas.Route) error {
	switch outputFormat {
	case print.JSONOutputFormat:
		details, err := json.MarshalIndent(route, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal static route: %w", err)
		}
		p.Outputln(string(details))

		return nil
	case print.YAMLOutputFormat:
		details, err := yaml.MarshalWithOptions(route, yaml.IndentSequence(true), yaml.UseJSONMarshaler())
		if err != nil {
			return fmt.Errorf("marshal static route: %w", err)
		}
		p.Outputln(string(details))

		return nil
	default:
		table := tables.NewTable()
		table.AddRow("ID", utils.PtrString(route.RouteId))
		table.AddSeparator()
		table.AddRow("PREFIX", utils.PtrString(route.Prefix))
		table.AddSeparator()
		table.AddRow("NEXTHOP", utils.PtrString(route.Nexthop))
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
	}
}
