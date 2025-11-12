package update

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	cliErr "github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/iaas/client"
	iaasUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/iaas/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/iaas"
)

const (
	routeIdArg = "ROUTE_ID"

	organizationIdFlag = "organization-id"
	networkAreaIdFlag  = "network-area-id"

	labelFlag = "labels"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	OrganizationId *string
	NetworkAreaId  *string
	RouteId        string
	Labels         *map[string]string
}

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("update %s", routeIdArg),
		Short: "Updates a static route in a STACKIT Network Area (SNA)",
		Long: fmt.Sprintf("%s\n%s\n",
			"Updates a static route in a STACKIT Network Area (SNA).",
			"This command is currently asynchonous only due to limitations in the waiting functionality of the SDK. This will be updated in a future release.",
		),
		Args: args.SingleArg(routeIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Updates the label(s) of a static route with ID "xxx" in a STACKIT Network Area with ID "yyy" in organization with ID "zzz"`,
				"$ stackit network-area route update xxx --labels key=value,foo=bar --organization-id yyy --network-area-id zzz",
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
			networkAreaLabel, err := iaasUtils.GetNetworkAreaName(ctx, apiClient, *model.OrganizationId, *model.NetworkAreaId)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get network area name: %v", err)
				networkAreaLabel = *model.NetworkAreaId
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("create static route: %w", err)
			}

			return outputResult(params.Printer, model.OutputFormat, networkAreaLabel, *resp)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), organizationIdFlag, "Organization ID")
	cmd.Flags().Var(flags.UUIDFlag(), networkAreaIdFlag, "STACKIT Network Area ID")
	cmd.Flags().StringToString(labelFlag, nil, "Labels are key-value string pairs which can be attached to a route. A label can be provided with the format key=value and the flag can be used multiple times to provide a list of labels")

	err := flags.MarkFlagsRequired(cmd, organizationIdFlag, networkAreaIdFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	routeId := inputArgs[0]
	globalFlags := globalflags.Parse(p, cmd)

	labels := flags.FlagToStringToStringPointer(p, cmd, labelFlag)

	if labels == nil {
		return nil, &cliErr.EmptyUpdateError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		OrganizationId:  flags.FlagToStringPointer(p, cmd, organizationIdFlag),
		NetworkAreaId:   flags.FlagToStringPointer(p, cmd, networkAreaIdFlag),
		RouteId:         routeId,
		Labels:          labels,
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *iaas.APIClient) iaas.ApiUpdateNetworkAreaRouteRequest {
	req := apiClient.UpdateNetworkAreaRoute(ctx, *model.OrganizationId, *model.NetworkAreaId, model.Region, model.RouteId)

	payload := iaas.UpdateNetworkAreaRoutePayload{
		Labels: utils.ConvertStringMapToInterfaceMap(model.Labels),
	}
	req = req.UpdateNetworkAreaRoutePayload(payload)

	return req
}

func outputResult(p *print.Printer, outputFormat, networkAreaLabel string, route iaas.Route) error {
	return p.OutputResult(outputFormat, route, func() error {
		p.Outputf("Updated static route for SNA %q.\nStatic route ID: %s\n", networkAreaLabel, utils.PtrString(route.Id))
		return nil
	})
}
