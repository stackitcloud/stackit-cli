package describe

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/goccy/go-yaml"
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

func NewCmd(p *print.Printer) *cobra.Command {
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
			model, err := parseInput(p, cmd, args)
			if err != nil {
				return err
			}

			// Configure API client
			apiClient, err := client.ConfigureClient(p)
			if err != nil {
				return err
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("read network: %w", err)
			}

			return outputResult(p, model.OutputFormat, resp)
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

func buildRequest(ctx context.Context, model *inputModel, apiClient *iaas.APIClient) iaas.ApiGetNetworkRequest {
	return apiClient.GetNetwork(ctx, model.ProjectId, model.NetworkId)
}

func outputResult(p *print.Printer, outputFormat string, network *iaas.Network) error {
	if network == nil {
		return fmt.Errorf("network cannot be nil")
	}
	switch outputFormat {
	case print.JSONOutputFormat:
		details, err := json.MarshalIndent(network, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal network: %w", err)
		}
		p.Outputln(string(details))

		return nil
	case print.YAMLOutputFormat:
		details, err := yaml.MarshalWithOptions(network, yaml.IndentSequence(true), yaml.UseJSONMarshaler())
		if err != nil {
			return fmt.Errorf("marshal network: %w", err)
		}
		p.Outputln(string(details))

		return nil
	default:
		var ipv4nameservers []string
		if network.Nameservers != nil {
			ipv4nameservers = append(ipv4nameservers, *network.Nameservers...)
		}

		var ipv4prefixes []string
		if network.Prefixes != nil {
			ipv4prefixes = append(ipv4prefixes, *network.Prefixes...)
		}

		var ipv6nameservers []string
		if network.NameserversV6 != nil {
			ipv6nameservers = append(ipv6nameservers, *network.NameserversV6...)
		}

		var ipv6prefixes []string
		if network.PrefixesV6 != nil {
			ipv6prefixes = append(ipv6prefixes, *network.PrefixesV6...)
		}

		table := tables.NewTable()
		table.AddRow("ID", utils.PtrString(network.NetworkId))
		table.AddSeparator()
		table.AddRow("NAME", utils.PtrString(network.Name))
		table.AddSeparator()
		table.AddRow("STATE", utils.PtrString(network.State))
		table.AddSeparator()

		if network.PublicIp != nil {
			table.AddRow("PUBLIC IP", *network.PublicIp)
			table.AddSeparator()
		}

		routed := false
		if network.Routed != nil {
			routed = *network.Routed
		}

		table.AddRow("ROUTED", routed)
		table.AddSeparator()

		if network.Gateway != nil {
			table.AddRow("IPv4 GATEWAY", *network.Gateway.Get())
			table.AddSeparator()
		}

		if len(ipv4nameservers) > 0 {
			table.AddRow("IPv4 NAME SERVERS", strings.Join(ipv4nameservers, ", "))
		}
		table.AddSeparator()
		if len(ipv4prefixes) > 0 {
			table.AddRow("IPv4 PREFIXES", strings.Join(ipv4prefixes, ", "))
		}
		table.AddSeparator()

		if network.Gatewayv6 != nil {
			table.AddRow("IPv6 GATEWAY", *network.Gatewayv6.Get())
			table.AddSeparator()
		}

		if len(ipv6nameservers) > 0 {
			table.AddRow("IPv6 NAME SERVERS", strings.Join(ipv6nameservers, ", "))
		}
		table.AddSeparator()
		if len(ipv6prefixes) > 0 {
			table.AddRow("IPv6 PREFIXES", strings.Join(ipv6prefixes, ", "))
		}
		table.AddSeparator()
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
	}
}
