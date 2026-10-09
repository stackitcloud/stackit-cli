package create

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	iaas "github.com/stackitcloud/stackit-sdk-go/services/iaas/v2alpha1api"
	"github.com/stackitcloud/stackit-sdk-go/services/iaas/v2alpha1api/wait"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/iaasalpha/client"
	iaasUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/iaasalpha/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/spinner"
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
		Use:   "create",
		Short: "Creates a regional configuration for a VPC",
		Long:  "Creates a regional configuration for a VPC. VPC can only be used in projects with enabled VPC.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Create a regional configuration "eu02" for a VPC with ID "xxx"`,
				`$ stackit alpha vpc region create --vpc-id xxx --region eu02`,
			),
			examples.NewExample(
				`Create a regional configuration with default DNS name servers, using the configured project and region`,
				`$ stackit alpha vpc region create --vpc-id xxx --ipv4-default-nameservers 8.8.8.8,8.8.4.4`,
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

			prompt := fmt.Sprintf("Are you sure you want to create the regional configuration %q for VPC %q?", model.Region, vpcLabel)
			if err := params.Printer.PromptForConfirmation(prompt); err != nil {
				return err
			}

			req := buildRequest(ctx, model, apiClient)

			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("create vpc region: %w", err)
			}
			if !model.Async {
				err := spinner.Run(params.Printer, "Create vpc region", func() error {
					resp, err = wait.CreateVPCRegionWaitHandler(ctx, apiClient.DefaultAPI, model.ProjectId, model.VpcId, model.Region).WaitWithContext(ctx)
					return err
				})
				if err != nil {
					return fmt.Errorf("wait for vpc region creation: %w", err)
				}
			}

			return outputResult(params.Printer, model.OutputFormat, model.Async, model.Region, vpcLabel, resp)
		},
	}
	configureFlags(cmd)

	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), vpcIdFlag, "VPC ID")
	cmd.Flags().StringSlice(ipv4DefaultNameserversFlag, nil, "List of default DNS name server IPs")

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
		GlobalFlagModel:        globalFlags,
		VpcId:                  flags.FlagToStringValue(p, cmd, vpcIdFlag),
		IPv4DefaultNameservers: flags.FlagToStringSliceValue(p, cmd, ipv4DefaultNameserversFlag),
	}

	p.DebugInputModel(model)
	return model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *iaas.APIClient) iaas.ApiCreateVPCRegionRequest {
	payload := iaas.CreateVPCRegionPayload{}
	if model.IPv4DefaultNameservers != nil {
		payload.Ipv4 = &iaas.RegionalVPCIPv4{
			DefaultNameservers: model.IPv4DefaultNameservers,
		}
	}

	return apiClient.DefaultAPI.CreateVPCRegion(ctx, model.ProjectId, model.VpcId, model.Region).CreateVPCRegionPayload(payload)
}

func outputResult(p *print.Printer, outputFormat string, async bool, region, vpcLabel string, resp *iaas.RegionalVPC) error {
	if resp == nil {
		return fmt.Errorf("create vpc region response is nil")
	}

	return p.OutputResult(outputFormat, resp, func() error {
		operationState := "Created"
		if async {
			operationState = "Triggered creation of"
		}

		p.Outputf("%s region configuration for VPC %q.\nRegion: %s\n", operationState, vpcLabel, region)

		return nil
	})
}
