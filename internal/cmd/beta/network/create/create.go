package create

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/goccy/go-yaml"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	cliErr "github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/projectname"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/iaas/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/spinner"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/iaas"
	"github.com/stackitcloud/stackit-sdk-go/services/iaas/wait"

	"github.com/spf13/cobra"
)

const (
	nameFlag               = "name"
	ipv4DnsNameServersFlag = "ipv4-dns-name-servers"
	ipv4PrefixLengthFlag   = "ipv4-prefix-length"
	ipv4PrefixFlag         = "ipv4-prefix"
	ipv4GatewayFlag        = "ipv4-gateway"
	ipv6DnsNameServersFlag = "ipv6-dns-name-servers"
	ipv6PrefixLengthFlag   = "ipv6-prefix-length"
	ipv6PrefixFlag         = "ipv6-prefix"
	ipv6GatewayFlag        = "ipv6-gateway"
	nonRoutedFlag          = "non-routed"
	noIpv4GatewayFlag      = "no-ipv4-gateway"
	noIpv6GatewayFlag      = "no-ipv6-gateway"
	labelFlag              = "labels"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	Name               *string
	IPv4DnsNameServers *[]string
	IPv4PrefixLength   *int64
	IPv4Prefix         *string
	IPv4Gateway        *string
	IPv6DnsNameServers *[]string
	IPv6PrefixLength   *int64
	IPv6Prefix         *string
	IPv6Gateway        *string
	NonRouted          bool
	NoIPv4Gateway      bool
	NoIPv6Gateway      bool
	Labels             *map[string]string
}

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Creates a network",
		Long:  "Creates a network.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Create a network with name "network-1"`,
				`$ stackit beta network create --name network-1`,
			),
			examples.NewExample(
				`Create a non-routed network with name "network-1"`,
				`$ stackit beta network create --name network-1 --non-routed`,
			),
			examples.NewExample(
				`Create a network with name "network-1" and no gateway`,
				`$ stackit beta network create --name network-1 --no-ipv4-gateway`,
			),
			examples.NewExample(
				`Create a network with name "network-1" and labels "key=value,key1=value1"`,
				`$ stackit beta network create --name network-1 --labels key=value,key1=value1`,
			),
			examples.NewExample(
				`Create an IPv4 network with name "network-1" with DNS name servers, a prefix and a gateway`,
				`$ stackit beta network create --name network-1  --ipv4-dns-name-servers "1.1.1.1,8.8.8.8,9.9.9.9" --ipv4-prefix "10.1.2.0/24" --ipv4-gateway "10.1.2.3"`,
			),
			examples.NewExample(
				`Create an IPv6 network with name "network-1" with DNS name servers, a prefix and a gateway`,
				`$ stackit beta network create --name network-1  --ipv6-dns-name-servers "2001:4860:4860::8888,2001:4860:4860::8844" --ipv6-prefix "2001:4860:4860::8888" --ipv6-gateway "2001:4860:4860::8888"`,
			),
		),
		RunE: func(cmd *cobra.Command, _ []string) error {
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

			projectLabel, err := projectname.GetProjectName(ctx, p, cmd)
			if err != nil {
				p.Debug(print.ErrorLevel, "get project name: %v", err)
				projectLabel = model.ProjectId
			} else if projectLabel == "" {
				projectLabel = model.ProjectId
			}

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to create a network for project %q?", projectLabel)
				err = p.PromptForConfirmation(prompt)
				if err != nil {
					return err
				}
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("create network : %w", err)
			}
			networkId := *resp.NetworkId

			// Wait for async operation, if async mode not enabled
			if !model.Async {
				s := spinner.New(p)
				s.Start("Creating network")
				_, err = wait.CreateNetworkWaitHandler(ctx, apiClient, model.ProjectId, networkId).WaitWithContext(ctx)
				if err != nil {
					return fmt.Errorf("wait for network creation: %w", err)
				}
				s.Stop()
			}

			return outputResult(p, model.OutputFormat, model.Async, projectLabel, resp)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().StringP(nameFlag, "n", "", "Network name")
	cmd.Flags().StringSlice(ipv4DnsNameServersFlag, []string{}, "List of DNS name servers for IPv4. Nameservers cannot be defined for routed networks")
	cmd.Flags().Int64(ipv4PrefixLengthFlag, 0, "The prefix length of the IPv4 network")
	cmd.Flags().String(ipv4PrefixFlag, "", "The IPv4 prefix of the network (CIDR)")
	cmd.Flags().String(ipv4GatewayFlag, "", "The IPv4 gateway of a network. If not specified, the first IP of the network will be assigned as the gateway")
	cmd.Flags().StringSlice(ipv6DnsNameServersFlag, []string{}, "List of DNS name servers for IPv6. Nameservers cannot be defined for routed networks")
	cmd.Flags().Int64(ipv6PrefixLengthFlag, 0, "The prefix length of the IPv6 network")
	cmd.Flags().String(ipv6PrefixFlag, "", "The IPv6 prefix of the network (CIDR)")
	cmd.Flags().String(ipv6GatewayFlag, "", "The IPv6 gateway of a network. If not specified, the first IP of the network will be assigned as the gateway")
	cmd.Flags().Bool(nonRoutedFlag, false, "If set to true, the network is not routed and therefore not accessible from other networks")
	cmd.Flags().Bool(noIpv4GatewayFlag, false, "If set to true, the network doesn't have an IPv4 gateway")
	cmd.Flags().Bool(noIpv6GatewayFlag, false, "If set to true, the network doesn't have an IPv6 gateway")
	cmd.Flags().StringToString(labelFlag, nil, "Labels are key-value string pairs which can be attached to a network. E.g. '--labels key1=value1,key2=value2,...'")

	err := flags.MarkFlagsRequired(cmd, nameFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &cliErr.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel:    globalFlags,
		Name:               flags.FlagToStringPointer(p, cmd, nameFlag),
		IPv4DnsNameServers: flags.FlagToStringSlicePointer(p, cmd, ipv4DnsNameServersFlag),
		IPv4PrefixLength:   flags.FlagToInt64Pointer(p, cmd, ipv4PrefixLengthFlag),
		IPv4Prefix:         flags.FlagToStringPointer(p, cmd, ipv4PrefixFlag),
		IPv4Gateway:        flags.FlagToStringPointer(p, cmd, ipv4GatewayFlag),
		IPv6DnsNameServers: flags.FlagToStringSlicePointer(p, cmd, ipv6DnsNameServersFlag),
		IPv6PrefixLength:   flags.FlagToInt64Pointer(p, cmd, ipv6PrefixLengthFlag),
		IPv6Prefix:         flags.FlagToStringPointer(p, cmd, ipv6PrefixFlag),
		IPv6Gateway:        flags.FlagToStringPointer(p, cmd, ipv6GatewayFlag),
		NonRouted:          flags.FlagToBoolValue(p, cmd, nonRoutedFlag),
		NoIPv4Gateway:      flags.FlagToBoolValue(p, cmd, noIpv4GatewayFlag),
		NoIPv6Gateway:      flags.FlagToBoolValue(p, cmd, noIpv6GatewayFlag),
		Labels:             flags.FlagToStringToStringPointer(p, cmd, labelFlag),
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

func buildRequest(ctx context.Context, model *inputModel, apiClient *iaas.APIClient) iaas.ApiCreateNetworkRequest {
	req := apiClient.CreateNetwork(ctx, model.ProjectId)
	addressFamily := &iaas.CreateNetworkAddressFamily{}

	if model.IPv6DnsNameServers != nil || model.IPv6PrefixLength != nil || model.IPv6Prefix != nil || model.NoIPv6Gateway || model.IPv6Gateway != nil {
		addressFamily.Ipv6 = &iaas.CreateNetworkIPv6Body{
			Nameservers:  model.IPv6DnsNameServers,
			PrefixLength: model.IPv6PrefixLength,
			Prefix:       model.IPv6Prefix,
		}

		if model.NoIPv6Gateway {
			addressFamily.Ipv6.Gateway = iaas.NewNullableString(nil)
		} else if model.IPv6Gateway != nil {
			addressFamily.Ipv6.Gateway = iaas.NewNullableString(model.IPv6Gateway)
		}
	}

	if model.IPv4DnsNameServers != nil || model.IPv4PrefixLength != nil || model.IPv4Prefix != nil || model.NoIPv4Gateway || model.IPv4Gateway != nil {
		addressFamily.Ipv4 = &iaas.CreateNetworkIPv4Body{
			Nameservers:  model.IPv4DnsNameServers,
			PrefixLength: model.IPv4PrefixLength,
			Prefix:       model.IPv4Prefix,
		}

		if model.NoIPv4Gateway {
			addressFamily.Ipv4.Gateway = iaas.NewNullableString(nil)
		} else if model.IPv4Gateway != nil {
			addressFamily.Ipv4.Gateway = iaas.NewNullableString(model.IPv4Gateway)
		}
	}

	var labelsMap *map[string]interface{}
	if model.Labels != nil && len(*model.Labels) > 0 {
		// convert map[string]string to map[string]interface{}
		labelsMap = utils.Ptr(map[string]interface{}{})
		for k, v := range *model.Labels {
			(*labelsMap)[k] = v
		}
	}

	routed := true
	if model.NonRouted {
		routed = false
	}

	payload := iaas.CreateNetworkPayload{
		Name:   model.Name,
		Labels: labelsMap,
		Routed: &routed,
	}

	if addressFamily.Ipv4 != nil || addressFamily.Ipv6 != nil {
		payload.AddressFamily = addressFamily
	}

	return req.CreateNetworkPayload(payload)
}

func outputResult(p *print.Printer, outputFormat string, async bool, projectLabel string, network *iaas.Network) error {
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
		operationState := "Created"
		if async {
			operationState = "Triggered creation of"
		}
		p.Outputf("%s network for project %q.\nNetwork ID: %s\n", operationState, projectLabel, utils.PtrString(network.NetworkId))
		return nil
	}
}
