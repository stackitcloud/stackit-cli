package list

import (
	"context"
	"fmt"
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
	iaasUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/iaasalpha/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

const vpcIdFlag = "vpc-id"

type inputModel struct {
	*globalflags.GlobalFlagModel
	VpcId string
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Lists all configured regions for a VPC",
		Long:  "Lists all configured regions for a VPC.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`List all configured regions for a VPC with ID "xxx"`,
				`$ stackit alpha vpc region list --vpc-id xxx`,
			),
			examples.NewExample(
				`List all configured regions for a VPC in JSON format`,
				`$ stackit alpha vpc region list --vpc-id xxx --output-format json`,
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

			vpcLabel, err := iaasUtils.GetVPCName(ctx, apiClient.DefaultAPI, model.ProjectId, model.VpcId)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get vpc name: %v", err)
			}
			if vpcLabel == "" {
				vpcLabel = model.VpcId
			}

			req := buildRequest(ctx, model, apiClient)

			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("list vpc regions: %w", err)
			}

			return outputResult(params.Printer, model.OutputFormat, vpcLabel, resp)
		},
	}
	configureFlags(cmd)

	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), vpcIdFlag, "VPC ID")

	cobra.CheckErr(flags.MarkFlagsRequired(cmd, vpcIdFlag))
}

func parseInput(p *print.Printer, cmd *cobra.Command, _ []string) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	model := &inputModel{
		GlobalFlagModel: globalFlags,
		VpcId:           flags.FlagToStringValue(p, cmd, vpcIdFlag),
	}

	p.DebugInputModel(model)
	return model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient iaas.DefaultAPI) iaas.ApiListVPCRegionsRequest {
	return apiClient.DefaultAPI.ListVPCRegions(ctx, model.ProjectId, model.VpcId)
}

func outputResult(p *print.Printer, outputFormat, vpcLabel string, resp *iaas.RegionalVPCList) error {
	if resp == nil {
		return fmt.Errorf("list vpc regions response is nil")
	}

	return p.OutputResult(outputFormat, resp, func() error {
		if len(resp.Regions) == 0 {
			p.Outputf("No regions found for VPC %q\n", vpcLabel)
			return nil
		}

		regionIds := slices.Sorted(maps.Keys(resp.Regions))
		for region := range resp.Regions {
			regions = append(regions, region)
		}

		slices.Sort(regions)

		table := tables.NewTable()
		table.SetHeader("REGION", "STATUS", "DNS NAME SERVERS")
		for _, region := range regions {
			regionConfig := resp.Regions[region]
			var dnsNames string
			if ipv4 := regionConfig.Ipv4; ipv4 != nil {
				dnsNames = strings.Join(ipv4.DefaultNameservers, ",")
			}

			table.AddRow(region, utils.PtrString(regionConfig.Status), dnsNames)
		}

		if err := table.Display(p); err != nil {
			return fmt.Errorf("render table: %w", err)
		}

		return nil
	})
}
