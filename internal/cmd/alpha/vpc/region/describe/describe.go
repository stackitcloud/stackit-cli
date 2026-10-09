package describe

import (
	"context"
	"fmt"
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
		Use:   "describe",
		Short: "Describes a regional configuration for a VPC",
		Long:  "Describes a regional configuration for a VPC.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Describe the regional configuration "eu02" for a VPC with ID "xxx"`,
				`$ stackit alpha vpc region describe --vpc-id xxx --region eu02`,
			),
			examples.NewExample(
				`Describe a regional configuration in JSON format, using the configured project and region`,
				`$ stackit alpha vpc region describe --vpc-id xxx --output-format json`,
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

			vpcName, err := iaasUtils.GetVPCName(ctx, apiClient.DefaultAPI, model.ProjectId, model.VpcId)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get vpc name: %v", err)
				vpcName = ""
			}

			req := buildRequest(ctx, model, apiClient)

			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("describe vpc region: %w", err)
			}

			return outputResult(params.Printer, model.OutputFormat, model.Region, model.VpcId, vpcName, resp)
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
	if globalFlags.Region == "" {
		return nil, &errors.RegionError{}
	}

	model := &inputModel{
		GlobalFlagModel: globalFlags,
		VpcId:           flags.FlagToStringValue(p, cmd, vpcIdFlag),
	}

	p.DebugInputModel(model)
	return model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *iaas.APIClient) iaas.ApiGetVPCRegionRequest {
	return apiClient.DefaultAPI.GetVPCRegion(ctx, model.ProjectId, model.VpcId, model.Region)
}

func outputResult(p *print.Printer, outputFormat, region, vpcId, vpcName string, resp *iaas.RegionalVPC) error {
	if resp == nil {
		return fmt.Errorf("describe vpc region response is nil")
	}

	return p.OutputResult(outputFormat, resp, func() error {
		table := tables.NewTable()
		table.AddRow("ID", vpcId)
		table.AddSeparator()
		if vpcName != "" {
			table.AddRow("NAME", vpcName)
			table.AddSeparator()
		}
		table.AddRow("REGION", region)
		table.AddSeparator()
		table.AddRow("STATUS", utils.PtrString(resp.Status))
		if ipv4 := resp.Ipv4; ipv4 != nil {
			table.AddSeparator()
			table.AddRow("DNS NAME SERVERS", strings.Join(ipv4.DefaultNameservers, ","))
		}

		if err := table.Display(p); err != nil {
			return fmt.Errorf("render table: %w", err)
		}

		return nil
	})
}
