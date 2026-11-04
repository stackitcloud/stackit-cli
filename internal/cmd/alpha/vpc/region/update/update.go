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
	iaasUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/iaasalpha/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
)

const (
	vpcIdFlag                  = "vpc-id"
	ipv4DefaultNameserversFlag = "ipv4-default-nameservers"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	VpcId                  string
	IPv4DefaultNameservers []string
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Updates a regional configuration for a VPC",
		Long:  "Updates a regional configuration for a VPC.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Update the default DNS name servers for a VPC with ID "xxx" in region "eu02"`,
				`$ stackit alpha vpc region update --vpc-id xxx --region eu02 --ipv4-default-nameservers 8.8.8.8,8.8.4.4`,
			),
			examples.NewExample(
				`Clear the default DNS name servers, using the configured project and region`,
				`$ stackit alpha vpc region update --vpc-id xxx --ipv4-default-nameservers ""`,
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

			prompt := fmt.Sprintf("Are you sure you want to update the regional configuration %q for VPC %q?", model.Region, vpcLabel)
			if err := params.Printer.PromptForConfirmation(prompt); err != nil {
				return err
			}

			req := buildRequest(ctx, model, apiClient.DefaultAPI)

			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("update vpc region: %w", err)
			}

			return outputResult(params.Printer, model.OutputFormat, model.Region, vpcLabel, resp)
		},
	}
	configureFlags(cmd)

	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), vpcIdFlag, "VPC ID")
	cmd.Flags().StringSlice(ipv4DefaultNameserversFlag, nil, "List of default DNS name server IPs")

	cobra.CheckErr(flags.MarkFlagsRequired(cmd, vpcIdFlag, ipv4DefaultNameserversFlag))
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
		GlobalFlagModel:        globalFlags,
		VpcId:                  flags.FlagToStringValue(p, cmd, vpcIdFlag),
		IPv4DefaultNameservers: flags.FlagToStringSliceValue(p, cmd, ipv4DefaultNameserversFlag),
	}

	p.DebugInputModel(model)
	return model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient iaas.DefaultAPI) iaas.ApiUpdateVPCRegionRequest {
	payload := iaas.UpdateVPCRegionPayload{
		Ipv4: &iaas.RegionalVPCIPv4{
			DefaultNameservers: model.IPv4DefaultNameservers,
		},
	}

	return apiClient.UpdateVPCRegion(ctx, model.ProjectId, model.VpcId, model.Region).UpdateVPCRegionPayload(payload)
}

func outputResult(p *print.Printer, outputFormat, region, vpcLabel string, resp *iaas.RegionalVPC) error {
	if resp == nil {
		return fmt.Errorf("update vpc region response is nil")
	}

	return p.OutputResult(outputFormat, resp, func() error {
		p.Outputf("Updated region configuration for VPC %q.\nRegion: %s\n", vpcLabel, region)

		return nil
	})
}
