package list

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	cliErr "github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/iaas/client"
	iaasUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/iaas/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/iaas"
)

const (
	limitFlag          = "limit"
	organizationIdFlag = "organization-id"
	networkAreaIdFlag  = "network-area-id"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	Limit          *int64
	OrganizationId *string
	NetworkAreaId  *string
}

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Lists all static routes in a STACKIT Network Area (SNA)",
		Long:  "Lists all static routes in a STACKIT Network Area (SNA).",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Lists all static routes in a STACKIT Network Area with ID "xxx" in organization with ID "yyy"`,
				"$ stackit network-area route list --network-area-id xxx --organization-id yyy",
			),
			examples.NewExample(
				`Lists all static routes in a STACKIT Network Area with ID "xxx" in organization with ID "yyy" in JSON format`,
				"$ stackit network-area route list --network-area-id xxx --organization-id yyy --output-format json",
			),
			examples.NewExample(
				`Lists up to 10 static routes in a STACKIT Network Area with ID "xxx" in organization with ID "yyy"`,
				"$ stackit network-area route list --network-area-id xxx --organization-id yyy --limit 10",
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
				return fmt.Errorf("list static routes: %w", err)
			}

			if resp.Items == nil || len(*resp.Items) == 0 {
				var networkAreaLabel string
				networkAreaLabel, err = iaasUtils.GetNetworkAreaName(ctx, apiClient, *model.OrganizationId, *model.NetworkAreaId)
				if err != nil {
					params.Printer.Debug(print.ErrorLevel, "get network area name: %v", err)
					networkAreaLabel = *model.NetworkAreaId
				}
				params.Printer.Info("No static routes found for STACKIT Network Area %q\n", networkAreaLabel)
				return nil
			}

			// Truncate output
			items := *resp.Items
			if model.Limit != nil && len(items) > int(*model.Limit) {
				items = items[:*model.Limit]
			}

			return outputResult(params.Printer, model.OutputFormat, items)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Int64(limitFlag, 0, "Maximum number of entries to list")
	cmd.Flags().Var(flags.UUIDFlag(), organizationIdFlag, "Organization ID")
	cmd.Flags().Var(flags.UUIDFlag(), networkAreaIdFlag, "STACKIT Network Area ID")

	err := flags.MarkFlagsRequired(cmd, organizationIdFlag, networkAreaIdFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command, _ []string) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	limit := flags.FlagToInt64Pointer(p, cmd, limitFlag)
	if limit != nil && *limit < 1 {
		return nil, &cliErr.FlagValidationError{
			Flag:    limitFlag,
			Details: "must be greater than 0",
		}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		Limit:           limit,
		OrganizationId:  flags.FlagToStringPointer(p, cmd, organizationIdFlag),
		NetworkAreaId:   flags.FlagToStringPointer(p, cmd, networkAreaIdFlag),
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *iaas.APIClient) iaas.ApiListNetworkAreaRoutesRequest {
	return apiClient.ListNetworkAreaRoutes(ctx, *model.OrganizationId, *model.NetworkAreaId, model.Region)
}

func outputResult(p *print.Printer, outputFormat string, routes []iaas.Route) error {
	return p.OutputResult(outputFormat, routes, func() error {
		table := tables.NewTable()
		table.SetHeader("Static Route ID", "Next Hop", "Next Hop Type", "Destination")

		for _, route := range routes {
			var nextHop string
			var nextHopType string
			var destination string
			if routeDest := route.Destination; routeDest != nil {
				if routeDest.DestinationCIDRv4 != nil {
					destination = *routeDest.DestinationCIDRv4.Value
				}
				if routeDest.DestinationCIDRv6 != nil {
					destination = *routeDest.DestinationCIDRv6.Value
				}
			}
			if routeNexthop := route.Nexthop; routeNexthop != nil {
				if routeNexthop.NexthopIPv4 != nil {
					nextHop = *routeNexthop.NexthopIPv4.Value
					nextHopType = *routeNexthop.NexthopIPv4.Type
				} else if routeNexthop.NexthopIPv6 != nil {
					nextHop = *routeNexthop.NexthopIPv6.Value
					nextHopType = *routeNexthop.NexthopIPv6.Type
				} else if routeNexthop.NexthopBlackhole != nil {
					nextHopType = *routeNexthop.NexthopBlackhole.Type
				} else if routeNexthop.NexthopInternet != nil {
					nextHopType = *routeNexthop.NexthopInternet.Type
				}
			}

			table.AddRow(
				utils.PtrString(route.Id),
				nextHop,
				nextHopType,
				destination,
			)
		}

		p.Outputln(table.Render())
		return nil
	})
}
