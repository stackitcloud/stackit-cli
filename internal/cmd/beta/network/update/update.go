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

	nameFlag       = "name"
	dnsServersFlag = "dns-servers"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	NetworkId  string
	Name       *string
	DnsServers *[]string
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
				`Update network with ID "xxx" with new name "network-1-new" and new dns servers`,
				`$ stackit beta network update xxx --name network-1-new --dns-servers "2.2.2.2"`,
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
	cmd.Flags().StringSlice(dnsServersFlag, nil, "List of DNS servers/nameservers IPs")
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	networkId := inputArgs[0]

	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &cliErr.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		Name:            flags.FlagToStringPointer(p, cmd, nameFlag),
		NetworkId:       networkId,
		DnsServers:      flags.FlagToStringSlicePointer(p, cmd, dnsServersFlag),
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

	payload := iaas.PartialUpdateNetworkPayload{
		Name: model.Name,
		AddressFamily: &iaas.UpdateNetworkAddressFamily{
			Ipv4: &iaas.UpdateNetworkIPv4{
				Nameservers: model.DnsServers,
			},
		},
	}

	return req.PartialUpdateNetworkPayload(payload)
}
