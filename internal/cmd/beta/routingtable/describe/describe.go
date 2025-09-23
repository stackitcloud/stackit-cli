package describe

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/goccy/go-yaml"
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
	"github.com/stackitcloud/stackit-sdk-go/services/iaasalpha"
)

const (
	organizationIdFlag = "organization-id"
	networkAreaIdFlag  = "network-area-id"
	routingTableArg    = "ROUTING_TABLE_ID_ARG"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	OrganizationId *string
	NetworkAreaId  *string
	RoutingTableId *string
}

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("describe %s", routingTableArg),
		Short: "Describe a routing-table",
		Long:  "Describe a routing-table",
		Args:  args.SingleArg(routingTableArg, nil),
		Example: examples.Build(
			examples.NewExample(
				`Describe a routing-table`,
				`$ stackit beta routing-table describe xxxx-xxxx-xxxx-xxxx --organization-id xxx --network-area-id yyy`,
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
			request := apiClient.GetRoutingTableOfArea(
				ctx,
				*model.OrganizationId,
				*model.NetworkAreaId,
				model.Region,
				*model.RoutingTableId,
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

func parseInput(p *print.Printer, cmd *cobra.Command, args []string) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)

	if len(args) == 0 {
		return nil, fmt.Errorf("at least one argument is required")
	}
	routingTableId := args[0]

	model := inputModel{
		GlobalFlagModel: globalFlags,
		NetworkAreaId:   flags.FlagToStringPointer(p, cmd, networkAreaIdFlag),
		OrganizationId:  flags.FlagToStringPointer(p, cmd, organizationIdFlag),
		RoutingTableId:  &routingTableId,
	}

	p.DebugInputModel(model)
	return &model, nil
}

func outputResult(p *print.Printer, outputFormat string, routingTable *iaasalpha.RoutingTable) error {
	switch outputFormat {
	case print.JSONOutputFormat:
		details, err := json.MarshalIndent(routingTable, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal routing-table describe: %w", err)
		}
		p.Outputln(string(details))

		return nil
	case print.YAMLOutputFormat:
		details, err := yaml.MarshalWithOptions(routingTable, yaml.IndentSequence(true), yaml.UseJSONMarshaler())
		if err != nil {
			return fmt.Errorf("marshal routing-table describe: %w", err)
		}
		p.Outputln(string(details))

		return nil
	default:
		var labels []string
		for key, value := range *routingTable.Labels {
			labels = append(labels, fmt.Sprintf("%s: %s", key, value))
		}

		table := tables.NewTable()
		table.SetHeader("ID", "NAME", "DESCRIPTION", "CREATED_AT", "UPDATED_AT", "DEFAULT", "LABELS", "SYSTEM_ROUTES")
		table.AddRow(
			utils.PtrString(routingTable.Id),
			utils.PtrString(routingTable.Name),
			utils.PtrString(routingTable.Description),
			routingTable.CreatedAt.String(),
			routingTable.UpdatedAt.String(),
			utils.PtrString(routingTable.Default),
			strings.Join(labels, "\n"),
			utils.PtrString(routingTable.SystemRoutes),
		)

		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}

		return nil
	}
}
