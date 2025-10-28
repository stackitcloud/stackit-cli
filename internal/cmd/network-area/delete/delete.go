package delete

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/iaas/client"
	iaasUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/iaas/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/iaas"
	"github.com/stackitcloud/stackit-sdk-go/services/iaas/wait"
)

const (
	areaIdArg          = "AREA_ID"
	organizationIdFlag = "organization-id"

	deprecationMessage = "The regional network area configuration %q for the area %q still exists.\n" +
		"The regional configuration of the network area was moved to the new command group `$ stackit network-area region`.\n" +
		"The regional area will be automatically deleted. This behavior is deprecated and will be removed after April 2026.\n" +
		"Use in the future the command `$ stackit network-area region delete` to delete the regional network area and afterwards delete the network-area with the command `$ stackit network-area delete`.\n"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	OrganizationId *string
	AreaId         string
}

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("delete %s", areaIdArg),
		Short: "Deletes a STACKIT Network Area (SNA)",
		Long: fmt.Sprintf("%s\n%s\n",
			"Deletes a STACKIT Network Area (SNA) in an organization.",
			"If the SNA is attached to any projects, the deletion will fail",
		),
		Args: args.SingleArg(areaIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Delete network area with ID "xxx" in organization with ID "yyy"`,
				"$ stackit network-area delete xxx --organization-id yyy",
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

			networkAreaLabel, err := iaasUtils.GetNetworkAreaName(ctx, apiClient, *model.OrganizationId, model.AreaId)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get network area name: %v", err)
				networkAreaLabel = model.AreaId
			}

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to delete network area %q?", networkAreaLabel)
				err = params.Printer.PromptForConfirmation(prompt)
				if err != nil {
					return err
				}
			}

			// Check if the network area has a regional configuration
			regionalArea, err := apiClient.GetNetworkAreaRegion(ctx, *model.OrganizationId, model.AreaId, model.Region).Execute()
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get regional area: %v", err)
			}
			if regionalArea != nil {
				params.Printer.Warn(deprecationMessage, model.Region, networkAreaLabel)
				err = apiClient.DeleteNetworkAreaRegion(ctx, *model.OrganizationId, model.AreaId, model.Region).Execute()
				if err != nil {
					return fmt.Errorf("delete network area region: %w", err)
				}
				_, err := wait.DeleteNetworkAreaRegionWaitHandler(ctx, apiClient, *model.OrganizationId, model.AreaId, model.Region).WaitWithContext(ctx)
				if err != nil {
					return fmt.Errorf("wait delete network area region: %w", err)
				}
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			err = req.Execute()
			if err != nil {
				return fmt.Errorf("delete network area: %w", err)
			}

			params.Printer.Outputf("Deleted STACKIT Network Area %q\n", networkAreaLabel)
			return nil
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), organizationIdFlag, "Organization ID")

	err := flags.MarkFlagsRequired(cmd, organizationIdFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	areaId := inputArgs[0]

	globalFlags := globalflags.Parse(p, cmd)

	model := inputModel{
		GlobalFlagModel: globalFlags,
		OrganizationId:  flags.FlagToStringPointer(p, cmd, organizationIdFlag),
		AreaId:          areaId,
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *iaas.APIClient) iaas.ApiDeleteNetworkAreaRequest {
	return apiClient.DeleteNetworkArea(ctx, *model.OrganizationId, model.AreaId)
}
