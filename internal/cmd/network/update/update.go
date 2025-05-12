package update

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
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/iaas/client"
	iaasUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/iaas/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/spinner"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/iaas"
	"github.com/stackitcloud/stackit-sdk-go/services/iaas/wait"

	"github.com/spf13/cobra"
)

const (
	networkIdArg = "NETWORK_ID"

	nameFlag               = "name"
	ipv4DnsNameServersFlag = "ipv4-dns-name-servers"
	ipv4GatewayFlag        = "ipv4-gateway"
	ipv6DnsNameServersFlag = "ipv6-dns-name-servers"
	ipv6GatewayFlag        = "ipv6-gateway"
	noIpv4GatewayFlag      = "no-ipv4-gateway"
	noIpv6GatewayFlag      = "no-ipv6-gateway"
	labelFlag              = "labels"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	NetworkId          string
	Name               *string
	IPv4DnsNameServers *[]string
	IPv4Gateway        *string
	IPv6DnsNameServers *[]string
	IPv6Gateway        *string
	NoIPv4Gateway      bool
	NoIPv6Gateway      bool
	Labels             *map[string]string
}

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("update %s", networkIdArg),
		Short: "Updates a network",
		Long:  "Updates a network.",
		Args:  args.SingleArg(networkIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Update network with ID "xxx" with new name "network-1-new"`,
				`$ stackit network update xxx --name network-1-new`,
			),
			examples.NewExample(
				`Update network with ID "xxx" with no gateway`,
				`$ stackit network update --no-ipv4-gateway`,
			),
			examples.NewExample(
				`Update IPv4 network with ID "xxx" with new name "network-1-new", new gateway and new DNS name servers`,
				`$ stackit network update xxx --name network-1-new --ipv4-dns-name-servers "2.2.2.2" --ipv4-gateway "10.1.2.3"`,
			),
			examples.NewExample(
				`Update IPv6 network with ID "xxx" with new name "network-1-new", new gateway and new DNS name servers`,
				`$ stackit network update xxx --name network-1-new --ipv6-dns-name-servers "2001:4860:4860::8888" --ipv6-gateway "2001:4860:4860::8888"`,
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

			networkLabel, err := iaasUtils.GetNetworkName(ctx, apiClient, model.ProjectId, model.NetworkId)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get network name: %v", err)
				networkLabel = model.NetworkId
			} else if networkLabel == "" {
				networkLabel = model.NetworkId
			}

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to update network %q?", networkLabel)
				err = params.Printer.PromptForConfirmation(prompt)
				if err != nil {
					return err
				}
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			err = req.Execute()
			if err != nil {
				return fmt.Errorf("update network area: %w", err)
			}
			networkId := model.NetworkId

			// Wait for async operation, if async mode not enabled
			if !model.Async {
				s := spinner.New(params.Printer)
				s.Start("Updating network")
				_, err = wait.UpdateNetworkWaitHandler(ctx, apiClient, model.ProjectId, networkId).WaitWithContext(ctx)
				if err != nil {
					return fmt.Errorf("wait for network update: %w", err)
				}
				s.Stop()
			}

			operationState := "Updated"
			if model.Async {
				operationState = "Triggered update of"
			}
			params.Printer.Info("%s network %q\n", operationState, networkLabel)
			return nil
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().StringP(nameFlag, "n", "", "Network name")
	cmd.Flags().StringSlice(ipv4DnsNameServersFlag, nil, "List of DNS name servers IPv4. Nameservers cannot be defined for routed networks")
	cmd.Flags().String(ipv4GatewayFlag, "", "The IPv4 gateway of a network. If not specified, the first IP of the network will be assigned as the gateway")
	cmd.Flags().StringSlice(ipv6DnsNameServersFlag, nil, "List of DNS name servers for IPv6. Nameservers cannot be defined for routed networks")
	cmd.Flags().String(ipv6GatewayFlag, "", "The IPv6 gateway of a network. If not specified, the first IP of the network will be assigned as the gateway")
	cmd.Flags().Bool(noIpv4GatewayFlag, false, "If set to true, the network doesn't have an IPv4 gateway")
	cmd.Flags().Bool(noIpv6GatewayFlag, false, "If set to true, the network doesn't have an IPv6 gateway")
	cmd.Flags().StringToString(labelFlag, nil, "Labels are key-value string pairs which can be attached to a network. E.g. '--labels key1=value1,key2=value2,...'")
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	networkId := inputArgs[0]

	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &cliErr.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel:    globalFlags,
		Name:               flags.FlagToStringPointer(p, cmd, nameFlag),
		NetworkId:          networkId,
		IPv4DnsNameServers: flags.FlagToStringSlicePointer(p, cmd, ipv4DnsNameServersFlag),
		IPv4Gateway:        flags.FlagToStringPointer(p, cmd, ipv4GatewayFlag),
		IPv6DnsNameServers: flags.FlagToStringSlicePointer(p, cmd, ipv6DnsNameServersFlag),
		IPv6Gateway:        flags.FlagToStringPointer(p, cmd, ipv6GatewayFlag),
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

func buildRequest(ctx context.Context, model *inputModel, apiClient *iaas.APIClient) iaas.ApiPartialUpdateNetworkRequest {
	req := apiClient.PartialUpdateNetwork(ctx, model.ProjectId, model.NetworkId)
	addressFamily := &iaas.UpdateNetworkAddressFamily{}

	var labelsMap *map[string]interface{}
	if model.Labels != nil && len(*model.Labels) > 0 {
		// convert map[string]string to map[string]interface{}
		labelsMap = utils.Ptr(map[string]interface{}{})
		for k, v := range *model.Labels {
			(*labelsMap)[k] = v
		}
	}

	if model.IPv6DnsNameServers != nil || model.NoIPv6Gateway || model.IPv6Gateway != nil {
		addressFamily.Ipv6 = &iaas.UpdateNetworkIPv6Body{
			Nameservers: model.IPv6DnsNameServers,
		}

		if model.NoIPv6Gateway {
			addressFamily.Ipv6.Gateway = iaas.NewNullableString(nil)
		} else if model.IPv6Gateway != nil {
			addressFamily.Ipv6.Gateway = iaas.NewNullableString(model.IPv6Gateway)
		}
	}

	if model.IPv4DnsNameServers != nil || model.NoIPv4Gateway || model.IPv4Gateway != nil {
		addressFamily.Ipv4 = &iaas.UpdateNetworkIPv4Body{
			Nameservers: model.IPv4DnsNameServers,
		}

		if model.NoIPv4Gateway {
			addressFamily.Ipv4.Gateway = iaas.NewNullableString(nil)
		} else if model.IPv4Gateway != nil {
			addressFamily.Ipv4.Gateway = iaas.NewNullableString(model.IPv4Gateway)
		}
	}

	payload := iaas.PartialUpdateNetworkPayload{
		Name:   model.Name,
		Labels: labelsMap,
	}

	if addressFamily.Ipv4 != nil || addressFamily.Ipv6 != nil {
		payload.AddressFamily = addressFamily
	}

	return req.PartialUpdateNetworkPayload(payload)
}
