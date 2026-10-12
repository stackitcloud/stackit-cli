package update

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
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

const (
	vpcIdFlag         = "vpc-id"
	routingTableIdArg = "ROUTING_TABLE_ID"
	nameFlag          = "name"
	descriptionFlag   = "description"
	labelFlag         = "labels"
	dynamicRoutesFlag = "dynamic-routes"
	systemRoutesFlag  = "system-routes"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	VpcId          string
	RoutingTableId string
	Name           *string
	Description    *string
	Labels         map[string]any
	DynamicRoutes  *bool
	SystemRoutes   *bool
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("update %s", routingTableIdArg),
		Short: "Updates a VPC routing-table",
		Long:  "Updates a routing-table for a VPC in a region.",
		Args:  args.SingleArg(routingTableIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Update the name of routing-table "xxx" for VPC "yyy" in region "eu02"`,
				`$ stackit alpha vpc routing-table update xxx --vpc-id yyy --region eu02 --name rt`,
			),
			examples.NewExample(
				`Disable dynamic routes, using the configured project and region`,
				`$ stackit alpha vpc routing-table update xxx --vpc-id yyy --dynamic-routes=false`,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			model, err := parseInput(params.Printer, cmd, args)
			if err != nil {
				return err
			}

			apiClient, err := client.ConfigureClient(params.Printer, params.CliVersion)
			if err != nil {
				return err
			}

			prompt := fmt.Sprintf("Are you sure you want to update the routing-table %q for VPC %q in region %q?", model.RoutingTableId, model.VpcId, model.Region)
			if err := params.Printer.PromptForConfirmation(prompt); err != nil {
				return err
			}

			req := buildRequest(ctx, model, apiClient.DefaultAPI)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("update vpc routing-table: %w", err)
			}

			return outputResult(params.Printer, model.OutputFormat, resp)
		},
	}
	configureFlags(cmd)

	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), vpcIdFlag, "VPC ID")
	cmd.Flags().String(nameFlag, "", "Name of the routing-table")
	cmd.Flags().String(descriptionFlag, "", "Description of the routing-table")
	cmd.Flags().StringToString(labelFlag, nil, "Key=value labels")
	cmd.Flags().Bool(dynamicRoutesFlag, false, "If set to false, prevents dynamic routes from propagating to the routing-table")
	cmd.Flags().Bool(systemRoutesFlag, false, "If set to false, disables automatic system routes")

	cobra.CheckErr(flags.MarkFlagsRequired(cmd, vpcIdFlag))
	cmd.MarkFlagsOneRequired(nameFlag, descriptionFlag, labelFlag, dynamicRoutesFlag, systemRoutesFlag)
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}
	if globalFlags.Region == "" {
		return nil, &errors.RegionError{}
	}

	model := &inputModel{
		GlobalFlagModel: globalFlags,
		VpcId:           flags.FlagToStringValue(p, cmd, vpcIdFlag),
		RoutingTableId:  inputArgs[0],
		Name:            flags.FlagToStringPointer(p, cmd, nameFlag),
		Description:     flags.FlagToStringPointer(p, cmd, descriptionFlag),
		Labels:          flags.FlagToStringToAny(p, cmd, labelFlag),
		DynamicRoutes:   flags.FlagToBoolPointer(p, cmd, dynamicRoutesFlag),
		SystemRoutes:    flags.FlagToBoolPointer(p, cmd, systemRoutesFlag),
	}

	p.DebugInputModel(model)
	return model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient iaas.DefaultAPI) iaas.ApiUpdateVPCRoutingTableRequest {
	payload := iaas.UpdateVPCRoutingTablePayload{
		Name:          model.Name,
		Description:   model.Description,
		Labels:        model.Labels,
		DynamicRoutes: model.DynamicRoutes,
		SystemRoutes:  model.SystemRoutes,
	}

	return apiClient.UpdateVPCRoutingTable(ctx, model.ProjectId, model.VpcId, model.Region, model.RoutingTableId).UpdateVPCRoutingTablePayload(payload)
}

func outputResult(p *print.Printer, outputFormat string, resp *iaas.VPCRoutingTable) error {
	if resp == nil {
		return fmt.Errorf("update vpc routing-table response is nil")
	}
	if resp.Id == nil || *resp.Id == "" {
		return fmt.Errorf("update vpc routing-table id is empty")
	}

	return p.OutputResult(outputFormat, resp, func() error {
		p.Outputf("Updated routing-table with ID %q.\n", *resp.Id)
		return nil
	})
}
