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
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/iaas"

	"github.com/spf13/cobra"
)

const (
	routeIdArg = "ROUTE_ID"

	organizationIdFlag = "organization-id"
	networkAreaIdFlag  = "network-area-id"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	OrganizationId *string
	NetworkAreaId  *string
	RouteId        string
}

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Deletes a static route in a STACKIT Network Area (SNA)",
		Long:  "Deletes a static route in a STACKIT Network Area (SNA).",
		Args:  args.SingleArg(routeIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Delete a static route with ID "xxx" in a STACKIT Network Area with ID "yyy" in organization with ID "zzz"`,
				"$ stackit beta network-area routes delete xxx --organization-id zzz --network-area-id yyy",
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

			networkAreaLabel, err := iaasUtils.GetNetworkAreaName(ctx, apiClient, *model.OrganizationId, *model.NetworkAreaId)
			if err != nil {
				p.Debug(print.ErrorLevel, "get network area name: %v", err)
				networkAreaLabel = *model.NetworkAreaId
			}

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to delete static route %q on STACKIT Network Area (SNA) %q?", model.RouteId, networkAreaLabel)
				err = p.PromptForConfirmation(prompt)
				if err != nil {
					return err
				}
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			err = req.Execute()
			if err != nil {
				return fmt.Errorf("delete static route: %w", err)
			}

			p.Info("Deleted static route %q on SNA %q\n", model.RouteId, networkAreaLabel)
			return nil
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), organizationIdFlag, "Organization ID")
	cmd.Flags().Var(flags.UUIDFlag(), networkAreaIdFlag, "STACKIT Network Area ID")

	err := flags.MarkFlagsRequired(cmd, organizationIdFlag, networkAreaIdFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	routeId := inputArgs[0]
	globalFlags := globalflags.Parse(p, cmd)

	model := inputModel{
		GlobalFlagModel: globalFlags,
		OrganizationId:  flags.FlagToStringPointer(p, cmd, organizationIdFlag),
		NetworkAreaId:   flags.FlagToStringPointer(p, cmd, networkAreaIdFlag),
		RouteId:         routeId,
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

func buildRequest(ctx context.Context, model *inputModel, apiClient *iaas.APIClient) iaas.ApiDeleteNetworkAreaRouteRequest {
	req := apiClient.DeleteNetworkAreaRoute(ctx, *model.OrganizationId, *model.NetworkAreaId, model.RouteId)
	return req
}
