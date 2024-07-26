package create

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/goccy/go-yaml"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/iaas/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/iaas/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/iaas"

	"github.com/spf13/cobra"
)

const (
	organizationIdFlag = "organization-id"
	networkAreaIdFlag  = "network-area-id"
	prefixFlag         = "prefix"
	nexthopFlag        = "next-hop"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	OrganizationId *string
	NetworkAreaId  *string
	Prefix         *string
	Nexthop        *string
}

func NewCmd(p *print.Printer) *cobra.Command {
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
				"$ stackit beta network-area routes create --organization-id yyy --network-area-id xxx --prefix 1.1.1.0/24 --next-hop 1.1.1.1",
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			model, err := parseInput(p, cmd)
			if err != nil {
				return err
			}

			// Configure API client
			apiClient, err := client.ConfigureClient(p)
			if err != nil {
				return err
			}

			// Get network area label
			networkAreaLabel, err := utils.GetNetworkAreaName(ctx, apiClient, *model.OrganizationId, *model.NetworkAreaId)
			if err != nil {
				p.Debug(print.ErrorLevel, "get network area name: %v", err)
				networkAreaLabel = *model.NetworkAreaId
			}

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to create a static route for STACKIT Network Area (SNA) %q?", networkAreaLabel)
				err = p.PromptForConfirmation(prompt)
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

			route, err := utils.GetRouteFromAPIResponse(*model.Prefix, *model.Nexthop, resp.Items)
			if err != nil {
				return err
			}

			return outputResult(p, model, networkAreaLabel, route)
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

func buildRequest(ctx context.Context, model *inputModel, apiClient *iaas.APIClient) iaas.ApiCreateNetworkAreaRouteRequest {
	req := apiClient.CreateNetworkAreaRoute(ctx, *model.OrganizationId, *model.NetworkAreaId)
	payload := iaas.CreateNetworkAreaRoutePayload{
		Ipv4: &[]iaas.Route{
			{
				Prefix:  model.Prefix,
				Nexthop: model.Nexthop,
			},
		},
	}
	return req.CreateNetworkAreaRoutePayload(payload)
}

func outputResult(p *print.Printer, model *inputModel, networkAreaLabel string, route iaas.Route) error {
	switch model.OutputFormat {
	case print.JSONOutputFormat:
		details, err := json.MarshalIndent(route, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal static route: %w", err)
		}
		p.Outputln(string(details))

		return nil
	case print.YAMLOutputFormat:
		details, err := yaml.MarshalWithOptions(route, yaml.IndentSequence(true))
		if err != nil {
			return fmt.Errorf("marshal static route: %w", err)
		}
		p.Outputln(string(details))

		return nil
	default:
		p.Outputf("Created static route for SNA %q.\nStatic route ID: %s\n", networkAreaLabel, *route.RouteId)
		return nil
	}
}
