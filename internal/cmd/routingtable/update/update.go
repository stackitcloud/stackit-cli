package update

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
	"github.com/stackitcloud/stackit-sdk-go/services/iaas"
)

const (
	descriptionFlag      = "description"
	labelFlag            = "labels"
	nameFlag             = "name"
	networkAreaIdFlag    = "network-area-id"
	nonDynamicRoutesFlag = "non-dynamic-routes"
	organizationIdFlag   = "organization-id"
	routingTableIdArg    = "ROUTE_TABLE_ID_ARG"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	OrganizationId   string
	NetworkAreaId    string
	NonDynamicRoutes bool
	RoutingTableId   string
	Description      *string
	Labels           *map[string]string
	Name             *string
}

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("update %s", routingTableIdArg),
		Short: "Updates a routing-table",
		Long:  "Updates a routing-table.",
		Args:  args.SingleArg(routingTableIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Updates the label(s) of a routing-table with ID "xxx" in organization with ID "yyy" and network-area with ID "zzz"`,
				"$ stackit routing-table update xxx --labels key=value,foo=bar --organization-id yyy --network-area-id zzz",
			),
			examples.NewExample(
				`Updates the name of a routing-table with ID "xxx" in organization with ID "yyy" and network-area with ID "zzz"`,
				"$ stackit routing-table update xxx --name foo --organization-id yyy --network-area-id zzz",
			),
			examples.NewExample(
				`Updates the description of a routing-table with ID "xxx" in organization with ID "yyy" and network-area with ID "zzz"`,
				"$ stackit routing-table update xxx --description foo --organization-id yyy --network-area-id zzz",
			),
			examples.NewExample(
				`Disables the dynamic_routes of a routing-table with ID "xxx" in organization with ID "yyy" and network-area with ID "zzz"`,
				"$ stackit routing-table update xxx --organization-id yyy --network-area-id zzz --non-dynamic-routes",
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
				prompt := fmt.Sprintf("Are you sure you want to update routing-table %q?", model.RoutingTableId)
				err = params.Printer.PromptForConfirmation(prompt)
				if err != nil {
					return err
				}
			}

			// Call API
			req := apiClient.UpdateRoutingTableOfArea(
				ctx,
				model.OrganizationId,
				model.NetworkAreaId,
				model.Region,
				model.RoutingTableId,
			)

			dynamicRoutes := true
			if model.NonDynamicRoutes {
				dynamicRoutes = false
			}

			payload := iaas.UpdateRoutingTableOfAreaPayload{
				Labels:        utils.ConvertStringMapToInterfaceMap(model.Labels),
				Name:          model.Name,
				Description:   model.Description,
				DynamicRoutes: &dynamicRoutes,
			}
			req = req.UpdateRoutingTableOfAreaPayload(payload)

			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("update routing-table %q : %w", model.RoutingTableId, err)
			}

			return outputResult(params.Printer, model.OutputFormat, model.NetworkAreaId, resp)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().String(descriptionFlag, "", "Description of the routing-table")
	cmd.Flags().String(nameFlag, "", "Name of the routing-table")
	cmd.Flags().StringToString(labelFlag, nil, "Key=value labels")
	cmd.Flags().Var(flags.UUIDFlag(), networkAreaIdFlag, "Network-Area ID")
	cmd.Flags().Bool(nonDynamicRoutesFlag, false, "If true, preventing dynamic routes from propagating to the routing-table.")
	cmd.Flags().Var(flags.UUIDFlag(), organizationIdFlag, "Organization ID")

	err := flags.MarkFlagsRequired(cmd, organizationIdFlag, networkAreaIdFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)

	if len(inputArgs) == 0 {
		return nil, fmt.Errorf("at least one argument is required")
	}
	routeTableId := inputArgs[0]

	model := inputModel{
		GlobalFlagModel:  globalFlags,
		Description:      flags.FlagToStringPointer(p, cmd, descriptionFlag),
		Labels:           flags.FlagToStringToStringPointer(p, cmd, labelFlag),
		Name:             flags.FlagToStringPointer(p, cmd, nameFlag),
		NetworkAreaId:    flags.FlagToStringValue(p, cmd, networkAreaIdFlag),
		NonDynamicRoutes: flags.FlagToBoolValue(p, cmd, nonDynamicRoutesFlag),
		OrganizationId:   flags.FlagToStringValue(p, cmd, organizationIdFlag),
		RoutingTableId:   routeTableId,
	}

	p.DebugInputModel(model)
	return &model, nil
}

func outputResult(p *print.Printer, outputFormat, networkAreaId string, routingTable *iaas.RoutingTable) error {
	if routingTable == nil {
		return fmt.Errorf("update routing-table response is empty")
	}

	if routingTable.Id == nil {
		return fmt.Errorf("update routing-table response is empty")
	}

	return p.OutputResult(outputFormat, routingTable, func() error {
		p.Outputf("Updated routing-table %q in network-area %q.", *routingTable.Id, networkAreaId)
		return nil
	})
}
