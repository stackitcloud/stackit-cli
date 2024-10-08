package update

import (
	"context"
	"fmt"

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
	ipv6DnsNameServersFlag = "ipv6-dns-name-servers"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	NetworkId          string
	Name               *string
	IPv4DnsNameServers *[]string
	IPv6DnsNameServers *[]string
}

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Updates a network",
		Long:  "Updates a network.",
		Args:  args.SingleArg(networkIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Update network with ID "xxx" with new name "network-1-new"`,
				`$ stackit beta network update xxx --name network-1-new`,
			),
			examples.NewExample(
				`Update IPv4 network with ID "xxx" with new name "network-1-new" and new DNS name servers`,
				`$ stackit beta network update xxx --name network-1-new --ipv4-dns-name-servers "2.2.2.2"`,
			),
			examples.NewExample(
				`Update IPv6 network with ID "xxx" with new name "network-1-new" and new DNS name servers`,
				`$ stackit beta network update xxx --name network-1-new --ipv6-dns-name-servers "2001:4860:4860::8888"`,
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

			networkLabel, err := iaasUtils.GetNetworkName(ctx, apiClient, model.ProjectId, model.NetworkId)
			if err != nil {
				p.Debug(print.ErrorLevel, "get network name: %v", err)
				networkLabel = model.NetworkId
			}

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to update network %q?", networkLabel)
				err = p.PromptForConfirmation(prompt)
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
				s := spinner.New(p)
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
			p.Info("%s network %q\n", operationState, networkLabel)
			return nil
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().StringP(nameFlag, "n", "", "Network name")
	cmd.Flags().StringSlice(ipv4DnsNameServersFlag, nil, "List of DNS name servers IPv4")
	cmd.Flags().StringSlice(ipv6DnsNameServersFlag, nil, "List of DNS name servers for IPv6")
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
		IPv6DnsNameServers: flags.FlagToStringSlicePointer(p, cmd, ipv6DnsNameServersFlag),
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

	if model.IPv6DnsNameServers != nil {
		addressFamily.Ipv6 = &iaas.UpdateNetworkIPv6Body{
			Nameservers: model.IPv6DnsNameServers,
		}
	}

	if model.IPv4DnsNameServers != nil {
		addressFamily.Ipv4 = &iaas.UpdateNetworkIPv4Body{
			Nameservers: model.IPv4DnsNameServers,
		}
	}

	payload := iaas.PartialUpdateNetworkPayload{
		Name: model.Name,
	}

	if addressFamily.Ipv4 != nil || addressFamily.Ipv6 != nil {
		payload.AddressFamily = addressFamily
	}

	return req.PartialUpdateNetworkPayload(payload)
}
