package delete

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/iaas/client"
	iaasUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/iaas/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/spinner"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/iaas"
	"github.com/stackitcloud/stackit-sdk-go/services/iaas/wait"

	"github.com/spf13/cobra"
)

const (
	areaIdArg          = "AREA_ID"
	organizationIdFlag = "organization-id"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	OrganizationId *string
	AreaId         string
}

func NewCmd(p *print.Printer) *cobra.Command {
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
				"$ stackit beta network-area delete xxx --organization-id yyy",
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			model, err := parseInput(p, cmd, args)
			if err != nil {
				return err
			}

			// Configure API client
			apiClient, err := client.ConfigureClient(p)
			if err != nil {
				return err
			}

			networkAreaLabel, err := iaasUtils.GetNetworkAreaName(ctx, apiClient, *model.OrganizationId, model.AreaId)
			if err != nil {
				p.Debug(print.ErrorLevel, "get network area name: %v", err)
				networkAreaLabel = model.AreaId
			}

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to delete network area %q?", networkAreaLabel)
				err = p.PromptForConfirmation(prompt)
				if err != nil {
					return err
				}
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			err = req.Execute()
			if err != nil {
				return fmt.Errorf("delete network area: %w", err)
			}

			// Wait for async operation, if async mode not enabled
			if !model.Async {
				s := spinner.New(p)
				s.Start("Deleting network area")
				_, err = wait.DeleteNetworkAreaWaitHandler(ctx, apiClient, *model.OrganizationId, model.AreaId).WaitWithContext(ctx)
				if err != nil {
					return fmt.Errorf("wait for network area deletion: %w", err)
				}
				s.Stop()
			}

			operationState := "Deleted"
			if model.Async {
				operationState = "Triggered deletion of"
			}
			p.Info("%s STACKIT Network Area %q\n", operationState, networkAreaLabel)
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

	if p.IsVerbosityDebug() {
		modelStr, err := print.BuildDebugStrFromInputModel(model)
		if err != nil {
			p.Debug(print.ErrorLevel, "convert model to string for debugging: %v", err)
		} else {
			p.Debug(print.DebugLevel, "parsed input values: %s", modelStr)
		}
	}

	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *iaas.APIClient) iaas.ApiDeleteNetworkAreaRequest {
	return apiClient.DeleteNetworkArea(ctx, *model.OrganizationId, model.AreaId)
}
