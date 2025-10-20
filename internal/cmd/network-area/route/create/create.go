package create

import (
	"context"
	"fmt"

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

	"github.com/spf13/cobra"
)

const (
	organizationIdFlag = "organization-id"
	networkAreaIdFlag  = "network-area-id"
	prefixFlag         = "prefix"
	nexthopFlag        = "next-hop"
	labelFlag          = "labels"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	OrganizationId *string
	NetworkAreaId  *string
	Prefix         *string
	Nexthop        *string
	Labels         *map[string]string
}

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Creates a static route in a STACKIT Network Area (SNA)",
		Long: fmt.Sprintf("%s\n%s\n",
			"Creates a static route in a STACKIT Network Area (SNA).",
			"This command is currently asynchonous only due to limitations in the waiting functionality of the SDK. This will be updated in a future release.",
		),
		Args: args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Create a static route with prefix "1.1.1.0/24" and next hop "1.1.1.1" in a STACKIT Network Area with ID "xxx" in organization with ID "yyy"`,
				"$ stackit network-area route create --organization-id yyy --network-area-id xxx --prefix 1.1.1.0/24 --next-hop 1.1.1.1",
			),
			examples.NewExample(
				`Create a static route with labels "key:value" and "foo:bar" with prefix "1.1.1.0/24" and next hop "1.1.1.1" in a STACKIT Network Area with ID "xxx" in organization with ID "yyy"`,
				"$ stackit network-area route create --labels key=value,foo=bar --organization-id yyy --network-area-id xxx --prefix 1.1.1.0/24 --next-hop 1.1.1.1",
			),
		),
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := context.Background()
			model, err := parseInput(params.Printer, cmd)
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

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to create a static route for STACKIT Network Area (SNA) %q?", networkAreaLabel)
				err = params.Printer.PromptForConfirmation(prompt)
				if err != nil {
					return err
				}
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("create static route: %w", err)
			}

			if resp.Items == nil || len(*resp.Items) == 0 {
				return fmt.Errorf("empty response from API")
			}

			route, err := iaasUtils.GetRouteFromAPIResponse(*model.Prefix, *model.Nexthop, resp.Items)
			if err != nil {
				return err
			}

			return outputResult(params.Printer, model.OutputFormat, networkAreaLabel, route)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), organizationIdFlag, "Organization ID")
	cmd.Flags().Var(flags.UUIDFlag(), networkAreaIdFlag, "STACKIT Network Area ID")
	cmd.Flags().Var(flags.CIDRFlag(), prefixFlag, "Static route prefix")
	cmd.Flags().String(nexthopFlag, "", "Next hop IP address. Must be a valid IPv4")
	cmd.Flags().StringToString(labelFlag, nil, "Labels are key-value string pairs which can be attached to a route. A label can be provided with the format key=value and the flag can be used multiple times to provide a list of labels")

	err := flags.MarkFlagsRequired(cmd, organizationIdFlag, networkAreaIdFlag, prefixFlag, nexthopFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)

	model := inputModel{
		GlobalFlagModel: globalFlags,
		OrganizationId:  flags.FlagToStringPointer(p, cmd, organizationIdFlag),
		NetworkAreaId:   flags.FlagToStringPointer(p, cmd, networkAreaIdFlag),
		Prefix:          flags.FlagToStringPointer(p, cmd, prefixFlag),
		Nexthop:         flags.FlagToStringPointer(p, cmd, nexthopFlag),
		Labels:          flags.FlagToStringToStringPointer(p, cmd, labelFlag),
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *iaas.APIClient) iaas.ApiCreateNetworkAreaRouteRequest {
	req := apiClient.CreateNetworkAreaRoute(ctx, *model.OrganizationId, *model.NetworkAreaId)

	payload := iaas.CreateNetworkAreaRoutePayload{
		Ipv4: &[]iaas.Route{
			{
				Prefix:  model.Prefix,
				Nexthop: model.Nexthop,
				Labels:  utils.ConvertStringMapToInterfaceMap(model.Labels),
			},
		},
	}
	return req.CreateNetworkAreaRoutePayload(payload)
}

func outputResult(p *print.Printer, outputFormat, networkAreaLabel string, route iaas.Route) error {
	return p.OutputResult(outputFormat, route, func() error {
		p.Outputf("Created static route for SNA %q.\nStatic route ID: %s\n", networkAreaLabel, utils.PtrString(route.RouteId))
		return nil
	})
}
