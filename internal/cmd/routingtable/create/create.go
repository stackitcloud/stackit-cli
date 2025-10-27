package create

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"
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
)

const (
	descriptionFlag      = "description"
	labelFlag            = "labels"
	nameFlag             = "name"
	networkAreaIdFlag    = "network-area-id"
	nonDynamicRoutesFlag = "non-dynamic-routes"
	nonSystemRoutesFlag  = "non-system-routes"
	organizationIdFlag   = "organization-id"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	Description      *string
	Labels           *map[string]string
	Name             *string
	NetworkAreaId    *string
	NonSystemRoutes  bool
	NonDynamicRoutes bool
	OrganizationId   *string
}

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Creates a routing-table",
		Long:  "Creates a routing-table.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				"Create a routing-table with name `rt`",
				`stackit routing-table create --organization-id xxx --network-area-id yyy --name "rt"`,
			),
			examples.NewExample(
				"Create a routing-table with name `rt` and description `some description`",
				`stackit routing-table create --organization-id xxx --network-area-id yyy --name "rt" --description "some description"`,
			),
			examples.NewExample(
				"Create a routing-table with name `rt` with system_routes disabled",
				`stackit routing-table create --organization-id xxx --network-area-id yyy --name "rt" --non-system-routes`,
			),
			examples.NewExample(
				"Create a routing-table with name `rt` with dynamic_routes disabled",
				`stackit routing-table create --organization-id xxx --network-area-id yyy --name "rt" --non-dynamic-routes`,
			),
		),
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := context.Background()
			model, err := parseInput(params.Printer, cmd, nil)
			if err != nil {
				return err
			}

			apiClient, err := client.ConfigureClient(params.Printer, params.CliVersion)
			if err != nil {
				return err
			}

			if !model.AssumeYes {
				prompt := "Are you sure you want to create a routing-table?"
				if err := params.Printer.PromptForConfirmation(prompt); err != nil {
					return err
				}
			}

			req, err := buildRequest(ctx, model, apiClient)
			if err != nil {
				return err
			}

			routingTableResp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("create routing-table request failed: %w", err)
			}

			return outputResult(params.Printer, model.OutputFormat, routingTableResp)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().String(descriptionFlag, "", "Description of the routing-table")
	cmd.Flags().StringToString(labelFlag, nil, "Key=value labels")
	cmd.Flags().String(nameFlag, "", "Name of the routing-table")
	cmd.Flags().Var(flags.UUIDFlag(), networkAreaIdFlag, "Network-Area ID")
	cmd.Flags().Bool(nonDynamicRoutesFlag, false, "If true, preventing dynamic routes from propagating to the routing-table.")
	cmd.Flags().Bool(nonSystemRoutesFlag, false, "If true, automatically disables routes for project-to-project communication.")
	cmd.Flags().Var(flags.UUIDFlag(), organizationIdFlag, "Organization ID")

	err := flags.MarkFlagsRequired(cmd, organizationIdFlag, networkAreaIdFlag, nameFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command, _ []string) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)

	model := &inputModel{
		GlobalFlagModel:  globalFlags,
		Description:      flags.FlagToStringPointer(p, cmd, descriptionFlag),
		NonDynamicRoutes: flags.FlagToBoolValue(p, cmd, nonDynamicRoutesFlag),
		Labels:           flags.FlagToStringToStringPointer(p, cmd, labelFlag),
		Name:             flags.FlagToStringPointer(p, cmd, nameFlag),
		NetworkAreaId:    flags.FlagToStringPointer(p, cmd, networkAreaIdFlag),
		OrganizationId:   flags.FlagToStringPointer(p, cmd, organizationIdFlag),
		NonSystemRoutes:  flags.FlagToBoolValue(p, cmd, nonSystemRoutesFlag),
	}

	p.DebugInputModel(model)
	return model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *iaas.APIClient) (iaas.ApiAddRoutingTableToAreaRequest, error) {
	systemRoutes := true
	if model.NonSystemRoutes {
		systemRoutes = false
	}

	dynamicRoutes := true
	if model.NonDynamicRoutes {
		dynamicRoutes = false
	}

	payload := iaas.AddRoutingTableToAreaPayload{
		Description:   model.Description,
		Name:          model.Name,
		Labels:        utils.ConvertStringMapToInterfaceMap(model.Labels),
		SystemRoutes:  utils.Ptr(systemRoutes),
		DynamicRoutes: utils.Ptr(dynamicRoutes),
	}

	return apiClient.AddRoutingTableToArea(
		ctx,
		*model.OrganizationId,
		*model.NetworkAreaId,
		model.Region,
	).AddRoutingTableToAreaPayload(payload), nil
}

func outputResult(p *print.Printer, outputFormat string, routingTable *iaas.RoutingTable) error {
	if routingTable == nil {
		return fmt.Errorf("create routing-table response is empty")
	}

	if routingTable.Id == nil {
		return fmt.Errorf("routing-table Id is empty")
	}

	return p.OutputResult(outputFormat, routingTable, func() error {
		var labels []string
		if routingTable.Labels != nil && len(*routingTable.Labels) > 0 {
			for key, value := range *routingTable.Labels {
				labels = append(labels, fmt.Sprintf("%s: %s", key, value))
			}
		}

		createdAt := ""
		if routingTable.CreatedAt != nil {
			createdAt = routingTable.CreatedAt.Format(time.RFC3339)
		}

		updatedAt := ""
		if routingTable.UpdatedAt != nil {
			updatedAt = routingTable.UpdatedAt.Format(time.RFC3339)
		}

		table := tables.NewTable()
		table.SetHeader("ID", "NAME", "DESCRIPTION", "CREATED_AT", "UPDATED_AT", "DEFAULT", "LABELS", "SYSTEM_ROUTES", "DYNAMIC_ROUTES")
		table.AddRow(
			utils.PtrString(routingTable.Id),
			utils.PtrString(routingTable.Name),
			utils.PtrString(routingTable.Description),
			createdAt,
			updatedAt,
			utils.PtrString(routingTable.Default),
			strings.Join(labels, "\n"),
			utils.PtrString(routingTable.SystemRoutes),
			utils.PtrString(routingTable.DynamicRoutes),
		)

		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}
		return nil
	})
}
