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
	"github.com/stackitcloud/stackit-sdk-go/services/iaas"
	"github.com/stackitcloud/stackit-sdk-go/services/iaas/wait"

	"github.com/spf13/cobra"
)

const (
	nameFlag               = "name"
	ipv4DnsNameServersFlag = "ipv4-dns-name-servers"
	ipv4PrefixLengthFlag   = "ipv4-prefix-length"
	ipv6DnsNameServersFlag = "ipv6-dns-name-servers"
	ipv6PrefixLengthFlag   = "ipv6-prefix-length"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	Name               *string
	IPv4DnsNameServers *[]string
	IPv4PrefixLength   *int64
	IPv6DnsNameServers *[]string
	IPv6PrefixLength   *int64
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
				`Create an IPv4 network with name "network-1" with DNS name servers and a prefix length`,
				`$ stackit beta network create --name network-1  --ipv4-dns-name-servers "1.1.1.1,8.8.8.8,9.9.9.9" --ipv4-prefix-length 25`,
			),
			examples.NewExample(
				`Create an IPv6 network with name "network-1" with DNS name servers and a prefix length`,
				`$ stackit beta network create --name network-1  --ipv6-dns-name-servers "2001:4860:4860::8888,2001:4860:4860::8844" --ipv6-prefix-length 56`,
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

			projectLabel, err := projectname.GetProjectName(ctx, p, cmd)
			if err != nil {
				p.Debug(print.ErrorLevel, "get project name: %v", err)
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

			return outputResult(p, model, projectLabel, resp)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().StringP(nameFlag, "n", "", "Network name")
	cmd.Flags().StringSlice(ipv4DnsNameServersFlag, []string{}, "List of DNS name servers for IPv4")
	cmd.Flags().Int64(ipv4PrefixLengthFlag, 0, "The prefix length of the IPv4 network")
	cmd.Flags().StringSlice(ipv6DnsNameServersFlag, []string{}, "List of DNS name servers for IPv6")
	cmd.Flags().Int64(ipv6PrefixLengthFlag, 0, "The prefix length of the IPv6 network")

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
		IPv6DnsNameServers: flags.FlagToStringSlicePointer(p, cmd, ipv6DnsNameServersFlag),
		IPv6PrefixLength:   flags.FlagToInt64Pointer(p, cmd, ipv6PrefixLengthFlag),
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

	if model.IPv6DnsNameServers != nil || model.IPv6PrefixLength != nil {
		addressFamily.Ipv6 = &iaas.CreateNetworkIPv6Body{
			Nameservers:  model.IPv6DnsNameServers,
			PrefixLength: model.IPv6PrefixLength,
		}
	}

	if model.IPv4DnsNameServers != nil || model.IPv4PrefixLength != nil {
		addressFamily.Ipv4 = &iaas.CreateNetworkIPv4Body{
			Nameservers:  model.IPv4DnsNameServers,
			PrefixLength: model.IPv4PrefixLength,
		}
	}

	payload := iaas.CreateNetworkPayload{
		Name: model.Name,
	}

	if addressFamily.Ipv4 != nil || addressFamily.Ipv6 != nil {
		payload.AddressFamily = addressFamily
	}

	return req.CreateNetworkPayload(payload)
}

func outputResult(p *print.Printer, model *inputModel, projectLabel string, network *iaas.Network) error {
	switch model.OutputFormat {
	case print.JSONOutputFormat:
		details, err := json.MarshalIndent(network, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal network: %w", err)
		}
		p.Outputln(string(details))

		return nil
	case print.YAMLOutputFormat:
		details, err := yaml.MarshalWithOptions(network, yaml.IndentSequence(true))
		if err != nil {
			return fmt.Errorf("marshal network: %w", err)
		}
		p.Outputln(string(details))

		return nil
	default:
		p.Outputf("Created network for project %q.\nNetwork ID: %s\n", projectLabel, *network.NetworkId)
		return nil
	}
}
