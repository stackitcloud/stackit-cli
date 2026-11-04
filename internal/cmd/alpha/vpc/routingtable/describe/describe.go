package describe

import (
	"context"
	"fmt"
	"maps"
	"slices"
	"strings"

	"github.com/spf13/cobra"
	iaas "github.com/stackitcloud/stackit-sdk-go/services/iaas/v2alpha1api"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/iaasalpha/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

const (
	vpcIdFlag         = "vpc-id"
	routingTableIdArg = "ROUTING_TABLE_ID"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	VpcId          string
	RoutingTableId string
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("describe %s", routingTableIdArg),
		Short: "Describes a VPC routing-table",
		Long:  "Describes a routing-table for a VPC in a region.",
		Args:  args.SingleArg(routingTableIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Describe routing-table "xxx" for VPC "yyy" in region "eu02"`,
				`$ stackit alpha vpc routing-table describe xxx --vpc-id yyy --region eu02`,
			),
			examples.NewExample(
				`Describe a routing-table in JSON format, using the configured project and region`,
				`$ stackit alpha vpc routing-table describe xxx --vpc-id yyy --output-format json`,
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

			req := buildRequest(ctx, model, apiClient.DefaultAPI)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("describe vpc routing-table: %w", err)
			}

			return outputResult(params.Printer, model.OutputFormat, resp)
		},
	}
	configureFlags(cmd)

	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), vpcIdFlag, "VPC ID")

	cobra.CheckErr(flags.MarkFlagsRequired(cmd, vpcIdFlag))
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
	}

	p.DebugInputModel(model)
	return model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient iaas.DefaultAPI) iaas.ApiGetVPCRoutingTableRequest {
	return apiClient.GetVPCRoutingTable(ctx, model.ProjectId, model.VpcId, model.Region, model.RoutingTableId)
}

func outputResult(p *print.Printer, outputFormat string, resp *iaas.VPCRoutingTable) error {
	if resp == nil {
		return fmt.Errorf("describe vpc routing-table response is nil")
	}

	return p.OutputResult(outputFormat, resp, func() error {
		table := tables.NewTable()
		table.AddRow("ID", utils.PtrString(resp.Id))
		table.AddSeparator()
		table.AddRow("NAME", resp.Name)
		table.AddSeparator()
		table.AddRow("DESCRIPTION", utils.PtrString(resp.Description))
		table.AddSeparator()

		if len(resp.Labels) > 0 {
			var labels []string
			for _, key := range slices.Sorted(maps.Keys(resp.Labels)) {
				labels = append(labels, fmt.Sprintf("%s: %v", key, resp.Labels[key]))
			}
			table.AddRow("LABELS", strings.Join(labels, "\n"))
			table.AddSeparator()
		}

		table.AddRow("SYSTEM ROUTES", utils.PtrString(resp.SystemRoutes))
		table.AddSeparator()
		table.AddRow("DYNAMIC ROUTES", utils.PtrString(resp.DynamicRoutes))
		table.AddSeparator()
		table.AddRow("CREATED AT", utils.ConvertTimePToDateTimeString(resp.CreatedAt))
		table.AddSeparator()
		table.AddRow("UPDATED AT", utils.ConvertTimePToDateTimeString(resp.UpdatedAt))

		if err := table.Display(p); err != nil {
			return fmt.Errorf("render table: %w", err)
		}
		return nil
	})
}
