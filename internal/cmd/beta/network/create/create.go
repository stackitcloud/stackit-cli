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
	nameFlag         = "name"
	dnsServersFlag   = "dns-servers"
	prefixLengthFlag = "prefix-length"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	Name         *string
	DnsServers   *[]string
	PrefixLength *int64
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
				`Create a network with name "network-1" with dns servers and a prefix length`,
				`$ stackit beta network create --name network-1  --dns-servers "1.1.1.1,8.8.8.8,9.9.9.9" --prefix-length 25`,
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
	cmd.Flags().StringSlice(dnsServersFlag, []string{}, "List of DNS servers/nameservers IPs")
	cmd.Flags().Int64(prefixLengthFlag, 0, "The default prefix length for networks")

	err := flags.MarkFlagsRequired(cmd, nameFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &cliErr.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		Name:            flags.FlagToStringPointer(p, cmd, nameFlag),
		DnsServers:      flags.FlagToStringSlicePointer(p, cmd, dnsServersFlag),
		PrefixLength:    flags.FlagToInt64Pointer(p, cmd, prefixLengthFlag),
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

	payload := iaas.CreateNetworkPayload{
		Name: model.Name,
		AddressFamily: &iaas.CreateNetworkAddressFamily{
			Ipv4: &iaas.CreateNetworkIPv4{
				Nameservers:  model.DnsServers,
				PrefixLength: model.PrefixLength,
			},
		},
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
