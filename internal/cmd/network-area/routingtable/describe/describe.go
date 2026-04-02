package describe

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/iaas"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/iaas/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

const (
	networkAreaIdFlag  = "network-area-id"
	organizationIdFlag = "organization-id"
	routingTableIdArg  = "ROUTING_TABLE_ID"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	NetworkAreaId  string
	OrganizationId string
	RoutingTableId string
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("describe %s", routingTableIdArg),
		Short: "Describes a routing-table",
		Long:  "Describes a routing-table",
		Args:  args.SingleArg(routingTableIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Describe a routing-table`,
				`$ stackit network-area routing-table describe xxx --organization-id xxx --network-area-id yyy`,
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
			request := apiClient.GetRoutingTableOfArea(
				ctx,
				model.OrganizationId,
				model.NetworkAreaId,
				model.Region,
				model.RoutingTableId,
			)

			response, err := request.Execute()
			if err != nil {
				return fmt.Errorf("describe routing-tables: %w", err)
			}

			return outputResult(params.Printer, model.OutputFormat, response)
		},
	}

	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), organizationIdFlag, "Organization ID")
	cmd.Flags().Var(flags.UUIDFlag(), networkAreaIdFlag, "Network-Area ID")

	err := flags.MarkFlagsRequired(cmd, organizationIdFlag, networkAreaIdFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	routingTableId := inputArgs[0]

	model := inputModel{
		GlobalFlagModel: globalFlags,
		NetworkAreaId:   flags.FlagToStringValue(p, cmd, networkAreaIdFlag),
		OrganizationId:  flags.FlagToStringValue(p, cmd, organizationIdFlag),
		RoutingTableId:  routingTableId,
	}

	p.DebugInputModel(model)
	return &model, nil
}

func outputResult(p *print.Printer, outputFormat string, routingTable *iaas.RoutingTable) error {
	if routingTable == nil {
		return fmt.Errorf("describe routingtable response is empty")
	}

	return p.OutputResult(outputFormat, routingTable, func() error {
		table := tables.NewTable()

		table.AddRow("ID", utils.PtrString(routingTable.Id))
		table.AddSeparator()

		table.AddRow("NAME", utils.PtrString(routingTable.Name))
		table.AddSeparator()

		table.AddRow("DESCRIPTION", utils.PtrString(routingTable.Description))
		table.AddSeparator()

		table.AddRow("DEFAULT", utils.PtrString(routingTable.Default))
		table.AddSeparator()

		if routingTable.Labels != nil && len(*routingTable.Labels) > 0 {
			var labels []string
			for key, value := range *routingTable.Labels {
				labels = append(labels, fmt.Sprintf("%s: %s", key, value))
			}
			table.AddRow("LABELS", strings.Join(labels, "\n"))
			table.AddSeparator()
		}

		table.AddRow("SYSTEM ROUTES", utils.PtrString(routingTable.SystemRoutes))
		table.AddSeparator()

		table.AddRow("DYNAMIC ROUTES", utils.PtrString(routingTable.DynamicRoutes))
		table.AddSeparator()

		table.AddRow("CREATED AT", utils.ConvertTimePToDateTimeString(routingTable.CreatedAt))
		table.AddSeparator()

		table.AddRow("UPDATED AT", utils.ConvertTimePToDateTimeString(routingTable.UpdatedAt))
		table.AddSeparator()

		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}
		return nil
	})
}
