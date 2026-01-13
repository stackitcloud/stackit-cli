package create

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/stackitcloud/stackit-cli/internal/pkg/types"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	cliErr "github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/iaas/client"
	rmClient "github.com/stackitcloud/stackit-cli/internal/pkg/services/resourcemanager/client"
	rmUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/resourcemanager/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/spinner"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/iaas"
	"github.com/stackitcloud/stackit-sdk-go/services/iaas/wait"

	"github.com/spf13/cobra"
)

const (
	nameFlag           = "name"
	organizationIdFlag = "organization-id"
	// Deprecated: dnsNameServersFlag is deprecated, because with iaas v2 the create endpoint for network area was separated, remove this after April 2026.
	dnsNameServersFlag = "dns-name-servers"
	// Deprecated: networkRangesFlag is deprecated, because with iaas v2 the create endpoint for network area was separated, remove this after April 2026.
	networkRangesFlag = "network-ranges"
	// Deprecated: transferNetworkFlag is deprecated, because with iaas v2 the create endpoint for network area was separated, remove this after April 2026.
	transferNetworkFlag = "transfer-network"
	// Deprecated: defaultPrefixLengthFlag is deprecated, because with iaas v2 the create endpoint for network area was separated, remove this after April 2026.
	defaultPrefixLengthFlag = "default-prefix-length"
	// Deprecated: maxPrefixLengthFlag is deprecated, because with iaas v2 the create endpoint for network area was separated, remove this after April 2026.
	maxPrefixLengthFlag = "max-prefix-length"
	// Deprecated: minPrefixLengthFlag is deprecated, because with iaas v2 the create endpoint for network area was separated, remove this after April 2026.
	minPrefixLengthFlag = "min-prefix-length"
	labelFlag           = "labels"

	deprecationMessage = "Deprecated and will be removed after April 2026. Use instead the new command `$ stackit network-area region` to configure these options for a network area."
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	Name           *string
	OrganizationId string
	// Deprecated: DnsNameServers is deprecated, because with iaas v2 the create endpoint for network area was separated, remove this after April 2026.
	DnsNameServers *[]string
	// Deprecated: NetworkRanges is deprecated, because with iaas v2 the create endpoint for network area was separated, remove this after April 2026.
	NetworkRanges *[]string
	// Deprecated: TransferNetwork is deprecated, because with iaas v2 the create endpoint for network area was separated, remove this after April 2026.
	TransferNetwork *string
	// Deprecated: DefaultPrefixLength is deprecated, because with iaas v2 the create endpoint for network area was separated, remove this after April 2026.
	DefaultPrefixLength *int64
	// Deprecated: MaxPrefixLength is deprecated, because with iaas v2 the create endpoint for network area was separated, remove this after April 2026.
	MaxPrefixLength *int64
	// Deprecated: MinPrefixLength is deprecated, because with iaas v2 the create endpoint for network area was separated, remove this after April 2026.
	MinPrefixLength *int64
	Labels          *map[string]string
}

// NetworkAreaResponses is a workaround, to keep the two responses of the iaas v2 api together for the json and yaml output
// Should be removed when the deprecated flags are removed
type NetworkAreaResponses struct {
	NetworkArea  iaas.NetworkArea   `json:"network_area"`
	RegionalArea *iaas.RegionalArea `json:"regional_area"`
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Creates a STACKIT Network Area (SNA)",
		Long:  "Creates a STACKIT Network Area (SNA) in an organization.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Create a network area with name "network-area-1" in organization with ID "xxx"`,
				`$ stackit network-area create --name network-area-1 --organization-id xxx"`,
			),
			examples.NewExample(
				`Create a network area with name "network-area-1" in organization with ID "xxx" with labels "key=value,key1=value1"`,
				`$ stackit network-area create --name network-area-1 --organization-id xxx --labels key=value,key1=value1`,
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

			var orgLabel string
			rmApiClient, err := rmClient.ConfigureClient(params.Printer, params.CliVersion)
			if err == nil {
				orgLabel, err = rmUtils.GetOrganizationName(ctx, rmApiClient, model.OrganizationId)
				if err != nil {
					params.Printer.Debug(print.ErrorLevel, "get organization name: %v", err)
					orgLabel = model.OrganizationId
				} else if orgLabel == "" {
					orgLabel = model.OrganizationId
				}
			} else {
				params.Printer.Debug(print.ErrorLevel, "configure resource manager client: %v", err)
			}

			prompt := fmt.Sprintf("Are you sure you want to create a network area for organization %q?", orgLabel)
			err = params.Printer.PromptForConfirmation(prompt)
			if err != nil {
				return err
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("create network area: %w", err)
			}
			if resp == nil || resp.Id == nil {
				return fmt.Errorf("create network area: empty response")
			}

			responses := &NetworkAreaResponses{
				NetworkArea: *resp,
			}

			if hasDeprecatedFlagsSet(model) {
				deprecatedFlags := getConfiguredDeprecatedFlags(model)
				params.Printer.Warn("the flags %q are deprecated and will be removed after April 2026. Use `$ stackit network-area region` to configure these options for a network area.\n", strings.Join(deprecatedFlags, ","))
				if resp == nil || resp.Id == nil {
					return fmt.Errorf("create network area: empty response")
				}
				reqNetworkArea := buildRequestNetworkAreaRegion(ctx, model, *resp.Id, apiClient)
				respNetworkArea, err := reqNetworkArea.Execute()
				if err != nil {
					return fmt.Errorf("create network area region: %w", err)
				}
				if !model.Async {
					s := spinner.New(params.Printer)
					s.Start("Create network area region")
					_, err = wait.CreateNetworkAreaRegionWaitHandler(ctx, apiClient, model.OrganizationId, *resp.Id, model.Region).WaitWithContext(ctx)
					if err != nil {
						return fmt.Errorf("wait for creating network area region %w", err)
					}
					s.Stop()
				}
				responses.RegionalArea = respNetworkArea
			}

			return outputResult(params.Printer, model.OutputFormat, orgLabel, responses)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().StringP(nameFlag, "n", "", "Network area name")
	cmd.Flags().Var(flags.UUIDFlag(), organizationIdFlag, "Organization ID")
	cmd.Flags().StringToString(labelFlag, nil, "Labels are key-value string pairs which can be attached to a network-area. E.g. '--labels key1=value1,key2=value2,...'")
	cmd.Flags().StringSlice(dnsNameServersFlag, nil, "List of DNS name server IPs")
	cmd.Flags().Var(flags.CIDRSliceFlag(), networkRangesFlag, "List of network ranges")
	cmd.Flags().Var(flags.CIDRFlag(), transferNetworkFlag, "Transfer network in CIDR notation")
	cmd.Flags().Int64(defaultPrefixLengthFlag, 0, "The default prefix length for networks in the network area")
	cmd.Flags().Int64(maxPrefixLengthFlag, 0, "The maximum prefix length for networks in the network area")
	cmd.Flags().Int64(minPrefixLengthFlag, 0, "The minimum prefix length for networks in the network area")

	cobra.CheckErr(cmd.Flags().MarkDeprecated(dnsNameServersFlag, deprecationMessage))
	cobra.CheckErr(cmd.Flags().MarkDeprecated(networkRangesFlag, deprecationMessage))
	cobra.CheckErr(cmd.Flags().MarkDeprecated(transferNetworkFlag, deprecationMessage))
	cobra.CheckErr(cmd.Flags().MarkDeprecated(defaultPrefixLengthFlag, deprecationMessage))
	cobra.CheckErr(cmd.Flags().MarkDeprecated(maxPrefixLengthFlag, deprecationMessage))
	cobra.CheckErr(cmd.Flags().MarkDeprecated(minPrefixLengthFlag, deprecationMessage))
	// Set the output for deprecation warnings to stderr
	cmd.Flags().SetOutput(os.Stderr)

	cmd.MarkFlagsRequiredTogether(networkRangesFlag, transferNetworkFlag)

	err := flags.MarkFlagsRequired(cmd, nameFlag, organizationIdFlag)
	cobra.CheckErr(err)
}

func hasDeprecatedFlagsSet(model *inputModel) bool {
	deprecatedFlags := getConfiguredDeprecatedFlags(model)
	return len(deprecatedFlags) > 0
}

func getConfiguredDeprecatedFlags(model *inputModel) []string {
	var result []string
	if model.DnsNameServers != nil {
		result = append(result, dnsNameServersFlag)
	}
	if model.NetworkRanges != nil {
		result = append(result, networkRangesFlag)
	}
	if model.TransferNetwork != nil {
		result = append(result, transferNetworkFlag)
	}
	if model.DefaultPrefixLength != nil {
		result = append(result, defaultPrefixLengthFlag)
	}
	if model.MaxPrefixLength != nil {
		result = append(result, maxPrefixLengthFlag)
	}
	if model.MinPrefixLength != nil {
		result = append(result, minPrefixLengthFlag)
	}
	return result
}

func parseInput(p *print.Printer, cmd *cobra.Command, _ []string) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)

	model := inputModel{
		GlobalFlagModel:     globalFlags,
		Name:                flags.FlagToStringPointer(p, cmd, nameFlag),
		OrganizationId:      flags.FlagToStringValue(p, cmd, organizationIdFlag),
		DnsNameServers:      flags.FlagToStringSlicePointer(p, cmd, dnsNameServersFlag),
		NetworkRanges:       flags.FlagToStringSlicePointer(p, cmd, networkRangesFlag),
		TransferNetwork:     flags.FlagToStringPointer(p, cmd, transferNetworkFlag),
		DefaultPrefixLength: flags.FlagToInt64Pointer(p, cmd, defaultPrefixLengthFlag),
		MaxPrefixLength:     flags.FlagToInt64Pointer(p, cmd, maxPrefixLengthFlag),
		MinPrefixLength:     flags.FlagToInt64Pointer(p, cmd, minPrefixLengthFlag),
		Labels:              flags.FlagToStringToStringPointer(p, cmd, labelFlag),
	}

	// Check if any of the deprecated **optional** fields are set and if no of the associated deprecated **required** fields is set.
	hasAllRequiredRegionalAreaFieldsSet := model.NetworkRanges != nil && model.TransferNetwork != nil
	hasOptionalRegionalAreaFieldsSet := model.DnsNameServers != nil || model.DefaultPrefixLength != nil || model.MaxPrefixLength != nil || model.MinPrefixLength != nil
	if hasOptionalRegionalAreaFieldsSet && !hasAllRequiredRegionalAreaFieldsSet {
		return nil, &cliErr.MultipleFlagsAreMissing{
			MissingFlags: []string{networkRangesFlag, transferNetworkFlag},
			SetFlags:     []string{dnsNameServersFlag, defaultPrefixLengthFlag, minPrefixLengthFlag, maxPrefixLengthFlag},
		}
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *iaas.APIClient) iaas.ApiCreateNetworkAreaRequest {
	req := apiClient.CreateNetworkArea(ctx, model.OrganizationId)

	payload := iaas.CreateNetworkAreaPayload{
		Name:   model.Name,
		Labels: utils.ConvertStringMapToInterfaceMap(model.Labels),
	}

	return req.CreateNetworkAreaPayload(payload)
}

func buildRequestNetworkAreaRegion(ctx context.Context, model *inputModel, networkAreaId string, apiClient *iaas.APIClient) iaas.ApiCreateNetworkAreaRegionRequest {
	req := apiClient.CreateNetworkAreaRegion(ctx, model.OrganizationId, networkAreaId, model.Region)

	var networkRanges []iaas.NetworkRange
	if model.NetworkRanges != nil {
		networkRanges = make([]iaas.NetworkRange, len(*model.NetworkRanges))
		for i, networkRange := range *model.NetworkRanges {
			networkRanges[i] = iaas.NetworkRange{
				Prefix: utils.Ptr(networkRange),
			}
		}
	}

	payload := iaas.CreateNetworkAreaRegionPayload{
		Ipv4: &iaas.RegionalAreaIPv4{
			DefaultNameservers: model.DnsNameServers,
			NetworkRanges:      utils.Ptr(networkRanges),
			TransferNetwork:    model.TransferNetwork,
			DefaultPrefixLen:   model.DefaultPrefixLength,
			MaxPrefixLen:       model.MaxPrefixLength,
			MinPrefixLen:       model.MinPrefixLength,
		},
	}

	return req.CreateNetworkAreaRegionPayload(payload)
}

func outputResult(p *print.Printer, outputFormat, orgLabel string, responses *NetworkAreaResponses) error {
	if responses == nil {
		return fmt.Errorf("network area is nil")
	}

	prettyOutputFunc := func() error {
		p.Outputf("Created STACKIT Network Area for organization %q.\nNetwork area ID: %s\n", orgLabel, utils.PtrString(responses.NetworkArea.Id))
		return nil
	}
	// If RegionalArea is NOT set in the response, then no deprecated Flags were set.
	// In this case, only the response of NetworkArea should be printed in JSON and yaml output, to avoid breaking changes after the deprecated fields are removed
	if responses.RegionalArea == nil {
		return p.OutputResult(outputFormat, responses.NetworkArea, prettyOutputFunc)
	}
	return p.OutputResult(outputFormat, responses, prettyOutputFunc)
}
