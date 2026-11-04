package create

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
	vpcIdFlag         = "vpc-id"
	nameFlag          = "name"
	descriptionFlag   = "description"
	labelFlag         = "labels"
	dynamicRoutesFlag = "dynamic-routes"
	systemRoutesFlag  = "system-routes"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	VpcId         string
	Name          string
	Description   *string
	Labels        map[string]any
	DynamicRoutes bool
	SystemRoutes  bool
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Creates a VPC routing-table",
		Long:  "Creates a routing-table for a VPC in a region. VPC can only be used in projects with enabled VPC.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Create a routing-table with name "rt" for VPC "xxx" in region "eu02"`,
				`$ stackit alpha vpc routing-table create --vpc-id xxx --region eu02 --name rt`,
			),
			examples.NewExample(
				`Create a routing-table with labels and dynamic routes disabled, using the configured project and region`,
				`$ stackit alpha vpc routing-table create --vpc-id xxx --name rt --labels env=dev --dynamic-routes=false`,
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
				params.Printer.Debug(print.ErrorLevel, "get VPC name: %v", err)
				vpcLabel = model.VpcId
			} else if vpcLabel == "" {
				vpcLabel = model.VpcId
			}

			prompt := fmt.Sprintf("Are you sure you want to create the routing-table %q for VPC %q in region %q?", model.Name, vpcLabel, model.Region)
			if err := params.Printer.PromptForConfirmation(prompt); err != nil {
				return err
			}

			req := buildRequest(ctx, model, apiClient.DefaultAPI)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("create vpc routing-table: %w", err)
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
	cmd.Flags().Bool(dynamicRoutesFlag, true, "If set to false, prevents dynamic routes from propagating to the routing-table")
	cmd.Flags().Bool(systemRoutesFlag, true, "If set to false, disables automatic system routes")

	cobra.CheckErr(flags.MarkFlagsRequired(cmd, vpcIdFlag, nameFlag))
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
		Name:            flags.FlagToStringValue(p, cmd, nameFlag),
		Description:     flags.FlagToStringPointer(p, cmd, descriptionFlag),
		Labels:          flags.FlagToStringToAny(p, cmd, labelFlag),
		DynamicRoutes:   flags.FlagToBoolValue(p, cmd, dynamicRoutesFlag),
		SystemRoutes:    flags.FlagToBoolValue(p, cmd, systemRoutesFlag),
	}

	p.DebugInputModel(model)
	return model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient iaas.DefaultAPI) iaas.ApiAddVPCRoutingTableRequest {
	payload := iaas.AddVPCRoutingTablePayload{
		Name:          model.Name,
		Description:   model.Description,
		Labels:        model.Labels,
		DynamicRoutes: new(model.DynamicRoutes),
		SystemRoutes:  new(model.SystemRoutes),
	}

	return apiClient.AddVPCRoutingTable(ctx, model.ProjectId, model.VpcId, model.Region).AddVPCRoutingTablePayload(payload)
}

func outputResult(p *print.Printer, outputFormat string, resp *iaas.VPCRoutingTable) error {
	if resp == nil {
		return fmt.Errorf("create vpc routing-table response is nil")
	}
	if resp.Id == nil || *resp.Id == "" {
		return fmt.Errorf("create vpc routing-table id is empty")
	}

	return p.OutputResult(outputFormat, resp, func() error {
		p.Outputf("Created routing-table with ID %q.\n", *resp.Id)
		return nil
	})
}
