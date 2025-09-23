package list

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/goccy/go-yaml"
	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/iaas/client"
	routeUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/routing-table/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/iaasalpha"
)

const (
	organizationIdFlag = "organization-id"
	networkAreaIdFlag  = "network-area-id"
	routingTableIdFlag = "routing-table-id"
	labelSelectorFlag  = "label-selector"
	limitFlag          = "limit"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	OrganizationId *string
	NetworkAreaId  *string
	RoutingTableId *string
	LabelSelector  *string
	Limit          *int64
}

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "list all routes within a routing-table",
		Long:  "list all routes within a routing-table",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`List all routes within a routing-table`,
				`$ stackit beta routing-table route list --routing-table-id xxx --organization-id yyy --network-area-id zzz`,
			),
			examples.NewExample(
				`List all routes within a routing-table with labels`,
				`$ stackit beta routing-table list --routing-table-id xxx --organization-id yyy --network-area-id zzz --label-selector env=dev,env=rc`,
			),
			examples.NewExample(
				`List all routes within a routing-tables with labels and limit to 10`,
				`$ stackit beta routing-table list --routing-table-id xxx --organization-id yyy --network-area-id zzz --label-selector env=dev,env=rc --limit 10`,
			),
		),
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := context.Background()
			model, err := parseInput(params.Printer, cmd)
			if err != nil {
				return err
			}

			// Configure API client
			apiClient, err := client.ConfigureAlphaClient(params.Printer, params.CliVersion)
			if err != nil {
				return err
			}

			// Call API
			request := apiClient.ListRoutesOfRoutingTable(
				ctx,
				*model.OrganizationId,
				*model.NetworkAreaId,
				model.Region,
				*model.RoutingTableId,
			)

			if model.LabelSelector != nil {
				request.LabelSelector(*model.LabelSelector)
			}

			response, err := request.Execute()
			if err != nil {
				return fmt.Errorf("list routes: %w", err)
			}

			if items := response.Items; items == nil || len(*items) == 0 {
				params.Printer.Info("No routes  found for routing-table %q\n", *model.RoutingTableId)
				return nil
			}

			// Truncate output
			items := *response.Items
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
	cmd.Flags().Var(flags.UUIDFlag(), organizationIdFlag, "Organization ID")
	cmd.Flags().Var(flags.UUIDFlag(), networkAreaIdFlag, "Network-Area ID")
	cmd.Flags().Var(flags.UUIDFlag(), routingTableIdFlag, "Routing-Table ID")
	cmd.Flags().String(labelSelectorFlag, "", "Filter by label")
	cmd.Flags().Int64(limitFlag, 0, "Maximum number of entries to list")

	err := flags.MarkFlagsRequired(cmd, organizationIdFlag, networkAreaIdFlag, routingTableIdFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command) (*inputModel, error) {
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
		Limit:           limit,
		NetworkAreaId:   flags.FlagToStringPointer(p, cmd, networkAreaIdFlag),
		OrganizationId:  flags.FlagToStringPointer(p, cmd, organizationIdFlag),
		RoutingTableId:  flags.FlagToStringPointer(p, cmd, routingTableIdFlag),
		LabelSelector:   flags.FlagToStringPointer(p, cmd, labelSelectorFlag),
	}

	p.DebugInputModel(model)
	return &model, nil
}

func outputResult(p *print.Printer, outputFormat string, items []iaasalpha.Route) error {
	switch outputFormat {
	case print.JSONOutputFormat:
		details, err := json.MarshalIndent(items, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal routes list: %w", err)
		}
		p.Outputln(string(details))

		return nil
	case print.YAMLOutputFormat:
		details, err := yaml.MarshalWithOptions(items, yaml.IndentSequence(true), yaml.UseJSONMarshaler())
		if err != nil {
			return fmt.Errorf("marshal routes list: %w", err)
		}
		p.Outputln(string(details))

		return nil
	default:
		table := tables.NewTable()
		table.SetHeader("ID", "DEST. TYPE", "DEST. VALUE", "NEXTHOP TYPE", "NEXTHOP VALUE", "LABELS", "CREATED", "UPDATED")
		for _, item := range items {
			destType, destValue, hopType, hopValue, labels := routeUtils.ExtractRouteDetails(item)

			table.AddRow(
				utils.PtrString(item.Id),
				destType,
				destValue,
				hopType,
				hopValue,
				labels,
				item.CreatedAt.String(),
				item.UpdatedAt.String(),
			)
		}
		return table.Display(p)
	}
}
