package update

import (
	"context"
	"fmt"
	"os"
	"strings"

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
	areaIdArg = "AREA_ID"

	nameFlag           = "name"
	organizationIdFlag = "organization-id"
	areaIdFlag         = "area-id"
	// Deprecated: dnsNameServersFlag is deprecated, because with iaas v2 the create endpoint for network area was separated, remove this after April 2026.
	dnsNameServersFlag = "dns-name-servers"
	// Deprecated: defaultPrefixLengthFlag is deprecated, because with iaas v2 the create endpoint for network area was separated, remove this after April 2026.
	defaultPrefixLengthFlag = "default-prefix-length"
	// Deprecated: maxPrefixLengthFlag is deprecated, because with iaas v2 the create endpoint for network area was separated, remove this after April 2026.
	maxPrefixLengthFlag = "max-prefix-length"
	// Deprecated: minPrefixLengthFlag is deprecated, because with iaas v2 the create endpoint for network area was separated, remove this after April 2026.
	minPrefixLengthFlag = "min-prefix-length"
	labelFlag           = "labels"

	deprecationMessage = "Deprecated and will be removed after April 2026. Use instead the new command `$ stackit network-area region` to configure these options for a network area."
)

// NetworkAreaResponses is a workaround, to keep the two responses of the iaas v2 api together for the json and yaml output
// Should be removed when the deprecated flags are removed
type NetworkAreaResponses struct {
	NetworkArea  iaas.NetworkArea   `json:"network_area"`
	RegionalArea *iaas.RegionalArea `json:"regional_area"`
}

type inputModel struct {
	*globalflags.GlobalFlagModel
	AreaId         string
	Name           *string
	OrganizationId *string
	// Deprecated: DnsNameServers is deprecated, because with iaas v2 the create endpoint for network area was separated, remove this after April 2026.
	DnsNameServers *[]string
	// Deprecated: DefaultPrefixLength is deprecated, because with iaas v2 the create endpoint for network area was separated, remove this after April 2026.
	DefaultPrefixLength *int64
	// Deprecated: MaxPrefixLength is deprecated, because with iaas v2 the create endpoint for network area was separated, remove this after April 2026.
	MaxPrefixLength *int64
	// Deprecated: MinPrefixLength is deprecated, because with iaas v2 the create endpoint for network area was separated, remove this after April 2026.
	MinPrefixLength *int64
	Labels          *map[string]string
}

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("update %s", areaIdArg),
		Short: "Updates a STACKIT Network Area (SNA)",
		Long:  "Updates a STACKIT Network Area (SNA) in an organization.",
		Args:  args.SingleArg(areaIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Update network area with ID "xxx" in organization with ID "yyy" with new name "network-area-1-new"`,
				"$ stackit network-area update xxx --organization-id yyy --name network-area-1-new",
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
				prompt := fmt.Sprintf("Are you sure you want to update a network area for organization %q?", orgLabel)
				err = params.Printer.PromptForConfirmation(prompt)
				if err != nil {
					return err
				}
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("update network area: %w", err)
			}

			if resp == nil || resp.Id == nil {
				return fmt.Errorf("create network area: empty response")
			}

			responses := NetworkAreaResponses{
				NetworkArea: *resp,
			}

			if hasDeprecatedFlagsSet(model) {
				deprecatedFlags := getConfiguredDeprecatedFlags(model)
				params.Printer.Warn("the flags %q are deprecated and will be removed after April 2026. Use `$ stackit network-area region` to configure these options for a network area.\n", strings.Join(deprecatedFlags, ","))
				if resp == nil || resp.Id == nil {
					return fmt.Errorf("create network area: empty response")
				}
				reqNetworkArea := buildRequestNetworkAreaRegion(ctx, model, apiClient)
				respNetworkArea, err := reqNetworkArea.Execute()
				if err != nil {
					return fmt.Errorf("create network area region: %w", err)
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
	cmd.Flags().StringSlice(dnsNameServersFlag, nil, "List of DNS name server IPs")
	cmd.Flags().Int64(defaultPrefixLengthFlag, 0, "The default prefix length for networks in the network area")
	cmd.Flags().Int64(maxPrefixLengthFlag, 0, "The maximum prefix length for networks in the network area")
	cmd.Flags().Int64(minPrefixLengthFlag, 0, "The minimum prefix length for networks in the network area")
	cmd.Flags().StringToString(labelFlag, nil, "Labels are key-value string pairs which can be attached to a network-area. E.g. '--labels key1=value1,key2=value2,...'")

	cobra.CheckErr(cmd.Flags().MarkDeprecated(dnsNameServersFlag, deprecationMessage))
	cobra.CheckErr(cmd.Flags().MarkDeprecated(defaultPrefixLengthFlag, deprecationMessage))
	cobra.CheckErr(cmd.Flags().MarkDeprecated(maxPrefixLengthFlag, deprecationMessage))
	cobra.CheckErr(cmd.Flags().MarkDeprecated(minPrefixLengthFlag, deprecationMessage))
	// Set the output for deprecation warnings to stderr
	cmd.Flags().SetOutput(os.Stderr)

	err := flags.MarkFlagsRequired(cmd, organizationIdFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	areaId := inputArgs[0]

	globalFlags := globalflags.Parse(p, cmd)

	model := inputModel{
		GlobalFlagModel:     globalFlags,
		Name:                flags.FlagToStringPointer(p, cmd, nameFlag),
		OrganizationId:      flags.FlagToStringPointer(p, cmd, organizationIdFlag),
		AreaId:              areaId,
		DnsNameServers:      flags.FlagToStringSlicePointer(p, cmd, dnsNameServersFlag),
		DefaultPrefixLength: flags.FlagToInt64Pointer(p, cmd, defaultPrefixLengthFlag),
		MaxPrefixLength:     flags.FlagToInt64Pointer(p, cmd, maxPrefixLengthFlag),
		MinPrefixLength:     flags.FlagToInt64Pointer(p, cmd, minPrefixLengthFlag),
		Labels:              flags.FlagToStringToStringPointer(p, cmd, labelFlag),
	}

	p.DebugInputModel(model)
	return &model, nil
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

func buildRequest(ctx context.Context, model *inputModel, apiClient *iaas.APIClient) iaas.ApiPartialUpdateNetworkAreaRequest {
	req := apiClient.PartialUpdateNetworkArea(ctx, *model.OrganizationId, model.AreaId)

	payload := iaas.PartialUpdateNetworkAreaPayload{
		Name:   model.Name,
		Labels: utils.ConvertStringMapToInterfaceMap(model.Labels),
	}

	return req.PartialUpdateNetworkAreaPayload(payload)
}

func buildRequestNetworkAreaRegion(ctx context.Context, model *inputModel, apiClient *iaas.APIClient) iaas.ApiUpdateNetworkAreaRegionRequest {
	req := apiClient.UpdateNetworkAreaRegion(ctx, *model.OrganizationId, model.AreaId, model.Region)

	payload := iaas.UpdateNetworkAreaRegionPayload{
		Ipv4: &iaas.UpdateRegionalAreaIPv4{
			DefaultNameservers: model.DnsNameServers,
			DefaultPrefixLen:   model.DefaultPrefixLength,
			MaxPrefixLen:       model.MaxPrefixLength,
			MinPrefixLen:       model.MinPrefixLength,
		},
	}

	return req.UpdateNetworkAreaRegionPayload(payload)
}

func outputResult(p *print.Printer, outputFormat, projectLabel string, responses NetworkAreaResponses) error {
	prettyOutputFunc := func() error {
		p.Outputf("Updated STACKIT Network Area for project %q.\n", projectLabel)
		return nil
	}

	// If RegionalArea is NOT set in the reponses, then no deprecated Flags were set.
	// In this case, only the response of NetworkArea should be printed in JSON and yaml output, to avoid breaking changes after the deprecated fields are removed
	if responses.RegionalArea == nil {
		return p.OutputResult(outputFormat, responses.NetworkArea, prettyOutputFunc)
	}

	return p.OutputResult(outputFormat, responses, func() error {
		p.Outputf("Updated STACKIT Network Area for project %q.\n", projectLabel)
		return nil
	})
}
