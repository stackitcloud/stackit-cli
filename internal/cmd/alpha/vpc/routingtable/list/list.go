package list

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
	labelSelectorFlag = "label-selector"
	filterFlag        = "filter"
	limitFlag         = "limit"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	VpcId         string
	LabelSelector *string
	Filter        *string
	Limit         *int64
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Lists VPC routing-tables",
		Long:  "Lists routing-tables for a VPC in a region.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`List routing-tables for VPC "xxx" in region "eu02"`,
				`$ stackit alpha vpc routing-table list --vpc-id xxx --region eu02`,
			),
			examples.NewExample(
				`List up to ten routing-tables with labels, using the configured project and region`,
				`$ stackit alpha vpc routing-table list --vpc-id xxx --label-selector env=dev --limit 10`,
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
				return fmt.Errorf("list vpc routing-tables: %w", err)
			}
			if resp == nil {
				return fmt.Errorf("list vpc routing-tables response is nil")
			}

			if model.Limit != nil && int64(len(resp.Items)) > *model.Limit {
				resp.Items = resp.Items[:*model.Limit]
			}

			return outputResult(params.Printer, model.OutputFormat, model.VpcId, resp)
		},
	}
	configureFlags(cmd)

	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), vpcIdFlag, "VPC ID")
	cmd.Flags().String(labelSelectorFlag, "", "Filter by label")
	cmd.Flags().String(filterFlag, "", "Filter resources by fields. A subset of expr-lang is supported. See https://expr-lang.org/docs/language-definition for usage details")
	cmd.Flags().Int64(limitFlag, 0, "Maximum number of entries to list")

	cobra.CheckErr(flags.MarkFlagsRequired(cmd, vpcIdFlag))
}

func parseInput(p *print.Printer, cmd *cobra.Command, _ []string) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}
	if globalFlags.Region == "" {
		return nil, &errors.RegionError{}
	}

	limit := flags.FlagToInt64Pointer(p, cmd, limitFlag)
	if limit != nil && *limit < 1 {
		return nil, &errors.FlagValidationError{
			Flag:    limitFlag,
			Details: "must be greater than 0",
		}
	}

	model := &inputModel{
		GlobalFlagModel: globalFlags,
		VpcId:           flags.FlagToStringValue(p, cmd, vpcIdFlag),
		LabelSelector:   flags.FlagToStringPointer(p, cmd, labelSelectorFlag),
		Filter:          flags.FlagToStringPointer(p, cmd, filterFlag),
		Limit:           limit,
	}

	p.DebugInputModel(model)
	return model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient iaas.DefaultAPI) iaas.ApiListVPCRoutingTablesRequest {
	req := apiClient.ListVPCRoutingTables(ctx, model.ProjectId, model.VpcId, model.Region)
	if model.LabelSelector != nil {
		req = req.LabelSelector(*model.LabelSelector)
	}
	if model.Filter != nil {
		req = req.Filter(*model.Filter)
	}

	return req
}

func outputResult(p *print.Printer, outputFormat, vpcId string, resp *iaas.VPCRoutingTableList) error {
	if resp == nil {
		return fmt.Errorf("list vpc routing-tables response is nil")
	}

	return p.OutputResult(outputFormat, resp, func() error {
		if len(resp.Items) == 0 {
			p.Outputf("No routing-tables found for VPC %q\n", vpcId)
			return nil
		}

		table := tables.NewTable()
		table.SetHeader("ID", "NAME", "DESCRIPTION", "LABELS", "SYSTEM ROUTES", "DYNAMIC ROUTES", "CREATED AT", "UPDATED AT")
		for _, routingTable := range resp.Items {
			var labels []string
			for _, key := range slices.Sorted(maps.Keys(routingTable.Labels)) {
				labels = append(labels, fmt.Sprintf("%s: %v", key, routingTable.Labels[key]))
			}

			table.AddRow(
				utils.PtrString(routingTable.Id),
				routingTable.Name,
				utils.PtrString(routingTable.Description),
				strings.Join(labels, "\n"),
				utils.PtrString(routingTable.SystemRoutes),
				utils.PtrString(routingTable.DynamicRoutes),
				utils.ConvertTimePToDateTimeString(routingTable.CreatedAt),
				utils.ConvertTimePToDateTimeString(routingTable.UpdatedAt),
			)
		}

		if err := table.Display(p); err != nil {
			return fmt.Errorf("render table: %w", err)
		}
		return nil
	})
}
