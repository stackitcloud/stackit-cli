package create

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
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

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Creates a network",
		Long:  "Creates a network.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Create a network with name "network-1"`,
				`$ stackit network create --name network-1`,
			),
			examples.NewExample(
				`Create a non-routed network with name "network-1"`,
				`$ stackit network create --name network-1 --non-routed`,
			),
			examples.NewExample(
				`Create a network with name "network-1" and no gateway`,
				`$ stackit network create --name network-1 --no-ipv4-gateway`,
			),
			examples.NewExample(
				`Create a network with name "network-1" and labels "key=value,key1=value1"`,
				`$ stackit network create --name network-1 --labels key=value,key1=value1`,
			),
			examples.NewExample(
				`Create an IPv4 network with name "network-1" with DNS name servers, a prefix and a gateway`,
				`$ stackit network create --name network-1 --non-routed --ipv4-dns-name-servers "1.1.1.1,8.8.8.8,9.9.9.9" --ipv4-prefix "10.1.2.0/24" --ipv4-gateway "10.1.2.3"`,
			),
			examples.NewExample(
				`Create an IPv6 network with name "network-1" with DNS name servers, a prefix and a gateway`,
				`$ stackit network create --name network-1  --ipv6-dns-name-servers "2001:4860:4860::8888,2001:4860:4860::8844" --ipv6-prefix "2001:4860:4860::8888" --ipv6-gateway "2001:4860:4860::8888"`,
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

			projectLabel, err := projectname.GetProjectName(ctx, params.Printer, params.CliVersion, cmd)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get project name: %v", err)
				projectLabel = model.ProjectId
			} else if projectLabel == "" {
				projectLabel = model.ProjectId
			}

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to create a network for project %q?", projectLabel)
				err = params.Printer.PromptForConfirmation(prompt)
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

			if resp == nil || resp.Id == nil {
				return fmt.Errorf("create network : empty response")
			}
			networkId := *resp.Id

			// Wait for async operation, if async mode not enabled
			if !model.Async {
				s := spinner.New(params.Printer)
				s.Start("Creating network")
				_, err = wait.CreateNetworkWaitHandler(ctx, apiClient, model.ProjectId, model.Region, networkId).WaitWithContext(ctx)
				if err != nil {
					return fmt.Errorf("wait for network creation: %w", err)
				}
				s.Stop()
			}

			return outputResult(params.Printer, model.OutputFormat, model.Async, projectLabel, resp)
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

	// IPv4 checks
	cmd.MarkFlagsMutuallyExclusive(ipv4PrefixFlag, ipv4PrefixLengthFlag)
	cmd.MarkFlagsMutuallyExclusive(ipv4GatewayFlag, ipv4PrefixLengthFlag)
	cmd.MarkFlagsMutuallyExclusive(ipv4GatewayFlag, noIpv4GatewayFlag)
	cmd.MarkFlagsMutuallyExclusive(noIpv4GatewayFlag, ipv4PrefixLengthFlag)

	// IPv6 checks
	cmd.MarkFlagsMutuallyExclusive(ipv6PrefixFlag, ipv6PrefixLengthFlag)
	cmd.MarkFlagsMutuallyExclusive(ipv6GatewayFlag, ipv6PrefixLengthFlag)
	cmd.MarkFlagsMutuallyExclusive(ipv6GatewayFlag, noIpv6GatewayFlag)

	err := flags.MarkFlagsRequired(cmd, nameFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command, _ []string) (*inputModel, error) {
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

	// IPv4 nameserver can not be set alone. IPv4 Prefix || IPv4 Prefix length must be set as well
	isIPv4NameserverSet := model.IPv4DnsNameServers != nil && len(*model.IPv4DnsNameServers) > 0
	isIPv4PrefixOrPrefixLengthSet := model.IPv4Prefix != nil || model.IPv4PrefixLength != nil
	if isIPv4NameserverSet && !isIPv4PrefixOrPrefixLengthSet {
		return nil, &cliErr.OneOfFlagsIsMissing{
			MissingFlags: []string{ipv4PrefixLengthFlag, ipv4PrefixFlag},
			SetFlag:      ipv4DnsNameServersFlag,
		}
	}
	isIPv4GatewaySet := model.IPv4Gateway != nil
	isIPv4PrefixSet := model.IPv4Prefix != nil
	if isIPv4GatewaySet && !isIPv4PrefixSet {
		return nil, &cliErr.DependingFlagIsMissing{
			MissingFlag: ipv4PrefixFlag,
			SetFlag:     ipv4GatewayFlag,
		}
	}

	// IPv6 nameserver can not be set alone. IPv6 Prefix || IPv6 Prefix length must be set as well
	isIPv6NameserverSet := model.IPv6DnsNameServers != nil && len(*model.IPv6DnsNameServers) > 0
	isIPv6PrefixOrPrefixLengthSet := model.IPv6Prefix != nil || model.IPv6PrefixLength != nil
	if isIPv6NameserverSet && !isIPv6PrefixOrPrefixLengthSet {
		return nil, &cliErr.OneOfFlagsIsMissing{
			MissingFlags: []string{ipv6PrefixLengthFlag, ipv6PrefixFlag},
			SetFlag:      ipv6DnsNameServersFlag,
		}
	}
	isIPv6GatewaySet := model.IPv6Gateway != nil
	isIPv6PrefixSet := model.IPv6Prefix != nil
	if isIPv6GatewaySet && !isIPv6PrefixSet {
		return nil, &cliErr.DependingFlagIsMissing{
			MissingFlag: ipv6PrefixFlag,
			SetFlag:     ipv6GatewayFlag,
		}
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *iaas.APIClient) iaas.ApiCreateNetworkRequest {
	req := apiClient.CreateNetwork(ctx, model.ProjectId, model.Region)
	var ipv4Network *iaas.CreateNetworkIPv4
	var ipv6Network *iaas.CreateNetworkIPv6

	if model.IPv6Prefix != nil {
		ipv6Network = &iaas.CreateNetworkIPv6{
			CreateNetworkIPv6WithPrefix: &iaas.CreateNetworkIPv6WithPrefix{
				Prefix:      model.IPv6Prefix,
				Nameservers: model.IPv6DnsNameServers,
			},
		}

		if model.NoIPv6Gateway {
			ipv6Network.CreateNetworkIPv6WithPrefix.Gateway = iaas.NewNullableString(nil)
		} else if model.IPv6Gateway != nil {
			ipv6Network.CreateNetworkIPv6WithPrefix.Gateway = iaas.NewNullableString(model.IPv6Gateway)
		}
	} else if model.IPv6PrefixLength != nil {
		ipv6Network = &iaas.CreateNetworkIPv6{
			CreateNetworkIPv6WithPrefixLength: &iaas.CreateNetworkIPv6WithPrefixLength{
				PrefixLength: model.IPv6PrefixLength,
				Nameservers:  model.IPv6DnsNameServers,
			},
		}
	}

	if model.IPv4Prefix != nil {
		ipv4Network = &iaas.CreateNetworkIPv4{
			CreateNetworkIPv4WithPrefix: &iaas.CreateNetworkIPv4WithPrefix{
				Prefix:      model.IPv4Prefix,
				Nameservers: model.IPv4DnsNameServers,
			},
		}

		if model.NoIPv4Gateway {
			ipv4Network.CreateNetworkIPv4WithPrefix.Gateway = iaas.NewNullableString(nil)
		} else if model.IPv4Gateway != nil {
			ipv4Network.CreateNetworkIPv4WithPrefix.Gateway = iaas.NewNullableString(model.IPv4Gateway)
		}
	} else if model.IPv4PrefixLength != nil {
		ipv4Network = &iaas.CreateNetworkIPv4{
			CreateNetworkIPv4WithPrefixLength: &iaas.CreateNetworkIPv4WithPrefixLength{
				PrefixLength: model.IPv4PrefixLength,
				Nameservers:  model.IPv4DnsNameServers,
			},
		}
	}

	routed := true
	if model.NonRouted {
		routed = false
	}

	payload := iaas.CreateNetworkPayload{
		Name:   model.Name,
		Labels: utils.ConvertStringMapToInterfaceMap(model.Labels),
		Routed: &routed,
		Ipv4:   ipv4Network,
		Ipv6:   ipv6Network,
	}

	return req.CreateNetworkPayload(payload)
}

func outputResult(p *print.Printer, outputFormat string, async bool, projectLabel string, network *iaas.Network) error {
	if network == nil {
		return fmt.Errorf("network cannot be nil")
	}
	return p.OutputResult(outputFormat, network, func() error {
		operationState := "Created"
		if async {
			operationState = "Triggered creation of"
		}
		p.Outputf("%s network for project %q.\nNetwork ID: %s\n", operationState, projectLabel, utils.PtrString(network.Id))
		return nil
	})
}
