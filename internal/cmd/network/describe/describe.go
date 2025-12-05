package describe

import (
	"context"
	"fmt"
	"strings"

	"github.com/stackitcloud/stackit-cli/internal/pkg/types"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/iaas/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/iaas"

	"github.com/spf13/cobra"
)

const (
	networkIdArg = "NETWORK_ID"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	NetworkId string
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("describe %s", networkIdArg),
		Short: "Shows details of a network",
		Long:  "Shows details of a network.",
		Args:  args.SingleArg(networkIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Show details of a network with ID "xxx"`,
				"$ stackit network describe xxx",
			),
			examples.NewExample(
				`Show details of a network with ID "xxx" in JSON format`,
				"$ stackit network describe xxx --output-format json",
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

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("read network: %w", err)
			}

			return outputResult(params.Printer, model.OutputFormat, resp)
		},
	}
	return cmd
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	networkId := inputArgs[0]

	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		NetworkId:       networkId,
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *iaas.APIClient) iaas.ApiGetNetworkRequest {
	return apiClient.GetNetwork(ctx, model.ProjectId, model.Region, model.NetworkId)
}

func outputResult(p *print.Printer, outputFormat string, network *iaas.Network) error {
	if network == nil {
		return fmt.Errorf("network cannot be nil")
	}
	return p.OutputResult(outputFormat, network, func() error {
		// IPv4
		var ipv4Nameservers, ipv4Prefixes []string
		var publicIp, ipv4Gateway *string
		if ipv4 := network.Ipv4; ipv4 != nil {
			if ipv4.Nameservers != nil {
				ipv4Nameservers = append(ipv4Nameservers, *ipv4.Nameservers...)
			}
			if ipv4.Prefixes != nil {
				ipv4Prefixes = append(ipv4Prefixes, *ipv4.Prefixes...)
			}
			if ipv4.PublicIp != nil {
				publicIp = ipv4.PublicIp
			}
			if ipv4.Gateway != nil && ipv4.Gateway.IsSet() {
				ipv4Gateway = ipv4.Gateway.Get()
			}
		}

		// IPv6
		var ipv6Nameservers, ipv6Prefixes []string
		var ipv6Gateway *string
		if ipv6 := network.Ipv6; ipv6 != nil {
			if ipv6.Nameservers != nil {
				ipv6Nameservers = append(ipv6Nameservers, *ipv6.Nameservers...)
			}
			if ipv6.Prefixes != nil {
				ipv6Prefixes = append(ipv6Prefixes, *ipv6.Prefixes...)
			}
			if ipv6.Gateway != nil && ipv6.Gateway.IsSet() {
				ipv6Gateway = ipv6.Gateway.Get()
			}
		}

		table := tables.NewTable()
		table.AddRow("ID", utils.PtrString(network.Id))
		table.AddSeparator()
		table.AddRow("NAME", utils.PtrString(network.Name))
		table.AddSeparator()
		table.AddRow("STATE", utils.PtrString(network.Status))
		table.AddSeparator()

		if publicIp != nil {
			table.AddRow("PUBLIC IP", *publicIp)
			table.AddSeparator()
		}

		routed := false
		if network.Routed != nil {
			routed = *network.Routed
		}

		table.AddRow("ROUTED", routed)
		table.AddSeparator()

		if network.RoutingTableId != nil {
			table.AddRow("ROUTING TABLE ID", utils.PtrString(network.RoutingTableId))
			table.AddSeparator()
		}

		if ipv4Gateway != nil {
			table.AddRow("IPv4 GATEWAY", *ipv4Gateway)
			table.AddSeparator()
		}

		if len(ipv4Nameservers) > 0 {
			table.AddRow("IPv4 NAME SERVERS", strings.Join(ipv4Nameservers, ", "))
		}
		table.AddSeparator()
		if len(ipv4Prefixes) > 0 {
			table.AddRow("IPv4 PREFIXES", strings.Join(ipv4Prefixes, ", "))
		}
		table.AddSeparator()

		if ipv6Gateway != nil {
			table.AddRow("IPv6 GATEWAY", *ipv6Gateway)
			table.AddSeparator()
		}

		if len(ipv6Nameservers) > 0 {
			table.AddRow("IPv6 NAME SERVERS", strings.Join(ipv6Nameservers, ", "))
			table.AddSeparator()
		}
		if len(ipv6Prefixes) > 0 {
			table.AddRow("IPv6 PREFIXES", strings.Join(ipv6Prefixes, ", "))
			table.AddSeparator()
		}
		if network.Labels != nil && len(*network.Labels) > 0 {
			var labels []string
			for key, value := range *network.Labels {
				labels = append(labels, fmt.Sprintf("%s: %s", key, value))
			}
			table.AddRow("LABELS", strings.Join(labels, "\n"))
			table.AddSeparator()
		}

		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}
		return nil
	})
}
