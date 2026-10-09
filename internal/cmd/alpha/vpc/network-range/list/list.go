package list

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/types"

	iaas "github.com/stackitcloud/stackit-sdk-go/services/iaas/v2alpha1api"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	cliErr "github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/iaasalpha/client"
	iaasUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/iaasalpha/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

const (
	limitFlag = "limit"
	vpcIdFlag = "vpc-id"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	Limit *int64
	VpcId string
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Lists all network ranges in a VPC",
		Long:  "Lists all network ranges in a VPC.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Lists all network ranges in a VPC with ID "xxx"`,
				"$ stackit alpha vpc network-range list --vpc-id xxx",
			),
			examples.NewExample(
				`Lists all network ranges in a VPC with ID "xxx" in JSON format`,
				"$ stackit alpha vpc network-range list --vpc-id xxx --output-format json",
			),
			examples.NewExample(
				`Lists up to 10 network ranges in a VPC with ID "xxx"`,
				"$ stackit alpha vpc network-range list --vpc-id xxx --limit 10",
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
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("list network ranges: %w", err)
			}

			items := resp.GetItems()

			vpcLabel, err := iaasUtils.GetVPCName(ctx, apiClient.DefaultAPI, model.ProjectId, model.VpcId)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get vpc name: %v", err)
				vpcLabel = model.VpcId
			}

			// Truncate output
			if model.Limit != nil && len(items) > int(*model.Limit) {
				items = items[:*model.Limit]
			}

			return outputResult(params.Printer, model.OutputFormat, vpcLabel, items)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Int64(limitFlag, 0, "Maximum number of entries to list")
	cmd.Flags().Var(flags.UUIDFlag(), vpcIdFlag, "VPC ID")

	err := flags.MarkFlagsRequired(cmd, vpcIdFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command, _ []string) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &cliErr.ProjectIdError{}
	}

	limit := flags.FlagToInt64Pointer(p, cmd, limitFlag)
	if limit != nil && *limit < 1 {
		return nil, &cliErr.FlagValidationError{
			Flag:    limitFlag,
			Details: "must be greater than 0",
		}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		Limit:           limit,
		VpcId:           flags.FlagToStringValue(p, cmd, vpcIdFlag),
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *iaas.APIClient) iaas.ApiListVPCNetworkRangesRequest {
	return apiClient.DefaultAPI.ListVPCNetworkRanges(ctx, model.ProjectId, model.VpcId, model.Region)
}

func outputResult(p *print.Printer, outputFormat, vpcLabel string, networkRanges []iaas.VPCNetworkRange) error {
	return p.OutputResult(outputFormat, networkRanges, func() error {
		if len(networkRanges) == 0 {
			p.Outputf("No network ranges found for VPC %q\n", vpcLabel)
			return nil
		}
		table := tables.NewTable()
		table.SetHeader("ID", "PREFIX", "DESCRIPTION", "IP VERSION", "STATUS")

		for _, networkRange := range networkRanges {
			if ipv4 := networkRange.VPCNetworkRangeIPv4; ipv4 != nil {
				table.AddRow(
					utils.PtrString(ipv4.Id),
					ipv4.Prefix,
					utils.PtrString(ipv4.Description),
					ipv4.IpVersion,
					utils.PtrString(ipv4.Status),
				)
			} else if ipv6 := networkRange.VPCNetworkRangeIPv6; ipv6 != nil {
				table.AddRow(
					utils.PtrString(ipv6.Id),
					ipv6.Prefix,
					utils.PtrString(ipv6.Description),
					ipv6.IpVersion,
					utils.PtrString(ipv6.Status),
				)
			}
		}

		p.Outputln(table.Render())
		return nil
	})
}
