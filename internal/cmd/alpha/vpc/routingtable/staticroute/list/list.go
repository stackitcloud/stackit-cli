package list

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
	iaasUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/iaasalpha/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

const (
	labelSelectorFlag  = "label-selector"
	limitFlag          = "limit"
	routingTableIdFlag = "routing-table-id"
	vpcIdFlag          = "vpc-id"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	LabelSelector  *string
	Limit          *int64
	RoutingTableId string
	VpcId          string
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Lists all static-routes within a VPC routing-table",
		Long:  "Lists all static-routes within a VPC routing-table.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`List all static-routes within a VPC routing-table`,
				`$ stackit alpha vpc routing-table static-route list --routing-table-id xxx --vpc-id yyy`,
			),
			examples.NewExample(
				`List all static-routes within a VPC routing-table with labels foo=bar,bar=baz`,
				`$ stackit alpha vpc routing-table static-route list --routing-table-id xxx --vpc-id yyy --label-selector foo=bar,bar=baz`,
			),
			examples.NewExample(
				`List up to 10 static-routes within a VPC routing-table`,
				`$ stackit alpha vpc routing-table static-route list --routing-table-id xxx --vpc-id yyy --limit 10`,
			),
		),
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := context.Background()

			model, err := parseInput(params.Printer, cmd, nil)
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
				return fmt.Errorf("list static-routes: %w", err)
			}

			routes := response.Items

			// Truncate output
			if model.Limit != nil && len(routes) > int(*model.Limit) {
				routes = routes[:*model.Limit]
			}

			return outputResult(params.Printer, model.OutputFormat, routes, model.VpcId, model.RoutingTableId)
		},
	}

	configureFlags(cmd)

	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Int64(limitFlag, 0, "Maximum number of entries to list")
	cmd.Flags().String(labelSelectorFlag, "", "Filter by label")
	cmd.Flags().Var(flags.UUIDFlag(), vpcIdFlag, "VPC ID")
	cmd.Flags().Var(flags.UUIDFlag(), routingTableIdFlag, "Routing-Table ID")

	err := flags.MarkFlagsRequired(cmd, vpcIdFlag, routingTableIdFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command, _ []string) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)

	limit := flags.FlagToInt64Pointer(p, cmd, limitFlag)
	if limit != nil && *limit < 1 {
		return nil, &errors.FlagValidationError{
			Flag:    limitFlag,
			Details: "must be greater than 0",
		}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		LabelSelector:   flags.FlagToStringPointer(p, cmd, labelSelectorFlag),
		Limit:           limit,
		RoutingTableId:  flags.FlagToStringValue(p, cmd, routingTableIdFlag),
		VpcId:           flags.FlagToStringValue(p, cmd, vpcIdFlag),
	}

	p.DebugInputModel(model)

	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient iaas.DefaultAPI) iaas.ApiListVPCStaticRoutesRequest {
	req := apiClient.ListVPCStaticRoutes(ctx, model.ProjectId, model.VpcId, model.Region, model.RoutingTableId)

	if model.LabelSelector != nil {
		req = req.LabelSelector(*model.LabelSelector)
	}

	return req
}

func outputResult(p *print.Printer, outputFormat string, routes []iaas.Route, vpcId, routeTableId string) error {
	if routes == nil {
		return fmt.Errorf("list static-routes response is nil")
	}

	return p.OutputResult(outputFormat, routes, func() error {
		if len(routes) == 0 {
			p.Outputf("No static-routes found for routing-table %q in VPC %q\n", routeTableId, vpcId)

			return nil
		}

		table := tables.NewTable()
		table.SetHeader("ID", "DESTINATION TYPE", "DESTINATION VALUE", "NEXTHOP TYPE", "NEXTHOP VALUE", "LABELS")

		for _, route := range routes {
			routeDetails := iaasUtils.ExtractRouteDetails(&route)
			table.AddRow(
				utils.PtrString(route.Id),
				routeDetails.DestType,
				routeDetails.DestValue,
				routeDetails.NextHopType,
				routeDetails.NextHopValue,
				routeDetails.Labels,
			)
			table.AddSeparator()
		}

		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}

		return nil
	})
}
