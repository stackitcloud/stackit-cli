package delete

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/iaas/client"
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
		Use:   fmt.Sprintf("delete %s", routingTableIdArg),
		Short: "Deletes a routing-table",
		Long:  "Deletes a routing-table",
		Args:  args.SingleArg(routingTableIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Deletes a a routing-table`,
				`$ stackit routing-table delete xxx --organization-id yyy --network-area-id zzz`,
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

			prompt := fmt.Sprintf("Are you sure you want to delete the routing-table %q?", model.RoutingTableId)
			err = params.Printer.PromptForConfirmation(prompt)
			if err != nil {
				return err
			}

			// Call API
			req := apiClient.DeleteRoutingTableFromArea(
				ctx,
				model.OrganizationId,
				model.NetworkAreaId,
				model.Region,
				model.RoutingTableId,
			)
			err = req.Execute()
			if err != nil {
				return fmt.Errorf("delete routing-table: %w", err)
			}

			params.Printer.Outputf("Routing-table %q deleted.", model.RoutingTableId)
			return nil
		},
	}

	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), networkAreaIdFlag, "Network-Area ID")
	cmd.Flags().Var(flags.UUIDFlag(), organizationIdFlag, "Organization ID")

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
