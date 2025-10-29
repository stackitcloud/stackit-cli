package delete

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/iaas/client"
	iaasUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/iaas/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/spinner"
	"github.com/stackitcloud/stackit-sdk-go/services/iaas"
	"github.com/stackitcloud/stackit-sdk-go/services/iaas/wait"
)

const (
	networkAreaIdFlag  = "network-area-id"
	organizationIdFlag = "organization-id"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	OrganizationId string
	NetworkAreaId  string
}

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Deletes a regional configuration for a STACKIT Network Area (SNA)",
		Long:  "Deletes a regional configuration for a STACKIT Network Area (SNA).",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Delete a regional configuration "eu02" for a STACKIT Network Area with ID "xxx" in organization with ID "yyy"`,
				`$ stackit network-area region delete --network-area-id xxx --region eu02 --organization-id yyy`,
			),
			examples.NewExample(
				`Delete a regional configuration "eu02" for a STACKIT Network Area with ID "xxx" in organization with ID "yyy", using the set region config`,
				`$ stackit config set --region eu02`,
				`$ stackit network-area region delete --network-area-id xxx --organization-id yyy`,
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

			// Get network area label
			networkAreaName, err := iaasUtils.GetNetworkAreaName(ctx, apiClient, model.OrganizationId, model.NetworkAreaId)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get network area name: %v", err)
				networkAreaName = model.NetworkAreaId
			}

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to delete the regional configuration %q for STACKIT Network Area (SNA) %q?", model.Region, networkAreaName)
				err = params.Printer.PromptForConfirmation(prompt)
				if err != nil {
					return err
				}
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			err = req.Execute()
			if err != nil {
				return fmt.Errorf("delete network area region: %w", err)
			}

			if !model.Async {
				s := spinner.New(params.Printer)
				s.Start("Delete network area region")
				_, err = wait.DeleteNetworkAreaRegionWaitHandler(ctx, apiClient, model.OrganizationId, model.NetworkAreaId, model.Region).WaitWithContext(ctx)
				if err != nil {
					return fmt.Errorf("wait for network area region deletion: %w", err)
				}
				s.Stop()
			}

			params.Printer.Outputf("Delete regional network area %q for %q\n", model.Region, networkAreaName)
			return nil
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), networkAreaIdFlag, "STACKIT Network Area (SNA) ID")
	cmd.Flags().Var(flags.UUIDFlag(), organizationIdFlag, "Organization ID")

	err := flags.MarkFlagsRequired(cmd, networkAreaIdFlag, organizationIdFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command, _ []string) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.Region == "" {
		return nil, &errors.RegionError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		NetworkAreaId:   flags.FlagToStringValue(p, cmd, networkAreaIdFlag),
		OrganizationId:  flags.FlagToStringValue(p, cmd, organizationIdFlag),
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *iaas.APIClient) iaas.ApiDeleteNetworkAreaRegionRequest {
	return apiClient.DeleteNetworkAreaRegion(ctx, model.OrganizationId, model.NetworkAreaId, model.Region)
}
