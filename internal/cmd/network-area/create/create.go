package create

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/goccy/go-yaml"
	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/iaas/client"
	rmClient "github.com/stackitcloud/stackit-cli/internal/pkg/services/resourcemanager/client"
	rmUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/resourcemanager/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/iaas"

	"github.com/spf13/cobra"
)

const (
	nameFlag                = "name"
	organizationIdFlag      = "organization-id"
	dnsNameServersFlag      = "dns-name-servers"
	networkRangesFlag       = "network-ranges"
	transferNetworkFlag     = "transfer-network"
	defaultPrefixLengthFlag = "default-prefix-length"
	maxPrefixLengthFlag     = "max-prefix-length"
	minPrefixLengthFlag     = "min-prefix-length"
	labelFlag               = "labels"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	Name                *string
	OrganizationId      *string
	DnsNameServers      *[]string
	NetworkRanges       *[]string
	TransferNetwork     *string
	DefaultPrefixLength *int64
	MaxPrefixLength     *int64
	MinPrefixLength     *int64
	Labels              *map[string]string
}

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Creates a STACKIT Network Area (SNA)",
		Long:  "Creates a STACKIT Network Area (SNA) in an organization.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Create a network area with name "network-area-1" in organization with ID "xxx" with network ranges and a transfer network`,
				`$ stackit network-area create --name network-area-1 --organization-id xxx --network-ranges "1.1.1.0/24,192.123.1.0/24" --transfer-network "192.160.0.0/24"`,
			),
			examples.NewExample(
				`Create a network area with name "network-area-2" in organization with ID "xxx" with network ranges, transfer network and DNS name server`,
				`$ stackit network-area create --name network-area-2 --organization-id xxx --network-ranges "1.1.1.0/24,192.123.1.0/24" --transfer-network "192.160.0.0/24" --dns-name-servers "1.1.1.1"`,
			),
			examples.NewExample(
				`Create a network area with name "network-area-3" in organization with ID "xxx" with network ranges, transfer network and additional options`,
				`$ stackit network-area create --name network-area-3 --organization-id xxx --network-ranges "1.1.1.0/24,192.123.1.0/24" --transfer-network "192.160.0.0/24" --default-prefix-length 25 --max-prefix-length 29 --min-prefix-length 24`,
			),
			examples.NewExample(
				`Create a network area with name "network-area-1" in organization with ID "xxx" with network ranges and a transfer network and labels "key=value,key1=value1"`,
				`$ stackit network-area create --name network-area-1 --organization-id xxx --network-ranges "1.1.1.0/24,192.123.1.0/24" --transfer-network "192.160.0.0/24" --labels key=value,key1=value1`,
			),
		),
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := context.Background()
			model, err := parseInput(params.Printer, cmd)
			if err != nil {
				return err
			}

			// Configure API client
			apiClient, err := client.ConfigureClient(params.Printer, params.CliVersion)
			if err != nil {
				return err
			}

			var orgLabel string
			rmApiClient, err := rmClient.ConfigureClient(params.Printer, params.CliVersion)
			if err == nil {
				orgLabel, err = rmUtils.GetOrganizationName(ctx, rmApiClient, *model.OrganizationId)
				if err != nil {
					params.Printer.Debug(print.ErrorLevel, "get organization name: %v", err)
					orgLabel = *model.OrganizationId
				} else if orgLabel == "" {
					orgLabel = *model.OrganizationId
				}
			} else {
				params.Printer.Debug(print.ErrorLevel, "configure resource manager client: %v", err)
			}

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to create a network area for organization %q?", orgLabel)
				err = params.Printer.PromptForConfirmation(prompt)
				if err != nil {
					return err
				}
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("create network area: %w", err)
			}

			return outputResult(params.Printer, model.OutputFormat, orgLabel, resp)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().StringP(nameFlag, "n", "", "Network area name")
	cmd.Flags().Var(flags.UUIDFlag(), organizationIdFlag, "Organization ID")
	cmd.Flags().StringSlice(dnsNameServersFlag, nil, "List of DNS name server IPs")
	cmd.Flags().Var(flags.CIDRSliceFlag(), networkRangesFlag, "List of network ranges")
	cmd.Flags().Var(flags.CIDRFlag(), transferNetworkFlag, "Transfer network in CIDR notation")
	cmd.Flags().Int64(defaultPrefixLengthFlag, 0, "The default prefix length for networks in the network area")
	cmd.Flags().Int64(maxPrefixLengthFlag, 0, "The maximum prefix length for networks in the network area")
	cmd.Flags().Int64(minPrefixLengthFlag, 0, "The minimum prefix length for networks in the network area")
	cmd.Flags().StringToString(labelFlag, nil, "Labels are key-value string pairs which can be attached to a network-area. E.g. '--labels key1=value1,key2=value2,...'")

	err := flags.MarkFlagsRequired(cmd, nameFlag, organizationIdFlag, networkRangesFlag, transferNetworkFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)

	model := inputModel{
		GlobalFlagModel:     globalFlags,
		Name:                flags.FlagToStringPointer(p, cmd, nameFlag),
		OrganizationId:      flags.FlagToStringPointer(p, cmd, organizationIdFlag),
		DnsNameServers:      flags.FlagToStringSlicePointer(p, cmd, dnsNameServersFlag),
		NetworkRanges:       flags.FlagToStringSlicePointer(p, cmd, networkRangesFlag),
		TransferNetwork:     flags.FlagToStringPointer(p, cmd, transferNetworkFlag),
		DefaultPrefixLength: flags.FlagToInt64Pointer(p, cmd, defaultPrefixLengthFlag),
		MaxPrefixLength:     flags.FlagToInt64Pointer(p, cmd, maxPrefixLengthFlag),
		MinPrefixLength:     flags.FlagToInt64Pointer(p, cmd, minPrefixLengthFlag),
		Labels:              flags.FlagToStringToStringPointer(p, cmd, labelFlag),
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

func buildRequest(ctx context.Context, model *inputModel, apiClient *iaas.APIClient) iaas.ApiCreateNetworkAreaRequest {
	req := apiClient.CreateNetworkArea(ctx, *model.OrganizationId)

	networkRanges := make([]iaas.NetworkRange, len(*model.NetworkRanges))
	for i, networkRange := range *model.NetworkRanges {
		networkRanges[i] = iaas.NetworkRange{
			Prefix: utils.Ptr(networkRange),
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

	payload := iaas.CreateNetworkAreaPayload{
		Name:   model.Name,
		Labels: labelsMap,
		AddressFamily: &iaas.CreateAreaAddressFamily{
			Ipv4: &iaas.CreateAreaIPv4{
				DefaultNameservers: model.DnsNameServers,
				NetworkRanges:      utils.Ptr(networkRanges),
				TransferNetwork:    model.TransferNetwork,
				DefaultPrefixLen:   model.DefaultPrefixLength,
				MaxPrefixLen:       model.MaxPrefixLength,
				MinPrefixLen:       model.MinPrefixLength,
			},
		},
	}

	return req.CreateNetworkAreaPayload(payload)
}

func outputResult(p *print.Printer, outputFormat, orgLabel string, networkArea *iaas.NetworkArea) error {
	if networkArea == nil {
		return fmt.Errorf("network area is nil")
	}
	switch outputFormat {
	case print.JSONOutputFormat:
		details, err := json.MarshalIndent(networkArea, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal network area: %w", err)
		}
		p.Outputln(string(details))

		return nil
	case print.YAMLOutputFormat:
		details, err := yaml.MarshalWithOptions(networkArea, yaml.IndentSequence(true), yaml.UseJSONMarshaler())
		if err != nil {
			return fmt.Errorf("marshal network area: %w", err)
		}
		p.Outputln(string(details))

		return nil
	default:
		p.Outputf("Created STACKIT Network Area for organization %q.\nNetwork area ID: %s\n", orgLabel, utils.PtrString(networkArea.AreaId))
		return nil
	}
}
