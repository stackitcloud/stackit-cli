package create

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/types"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
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
)

const (
	networkAreaIdFlag           = "network-area-id"
	organizationIdFlag          = "organization-id"
	ipv4DefaultNameservers      = "ipv4-default-nameservers"
	ipv4DefaultPrefixLengthFlag = "ipv4-default-prefix-length"
	ipv4MaxPrefixLengthFlag     = "ipv4-max-prefix-length"
	ipv4MinPrefixLengthFlag     = "ipv4-min-prefix-length"
	ipv4NetworkRangesFlag       = "ipv4-network-ranges"
	ipv4TransferNetworkFlag     = "ipv4-transfer-network"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	OrganizationId string
	NetworkAreaId  string

	IPv4DefaultNameservers  *[]string
	IPv4DefaultPrefixLength *int64
	IPv4MaxPrefixLength     *int64
	IPv4MinPrefixLength     *int64
	IPv4NetworkRanges       []string
	IPv4TransferNetwork     string
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Creates a new regional configuration for a STACKIT Network Area (SNA)",
		Long:  "Creates a new regional configuration for a STACKIT Network Area (SNA).",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Create a new regional configuration "eu02" for a STACKIT Network Area with ID "xxx" in organization with ID "yyy", ipv4 network range "192.168.0.0/24" and ipv4 transfer network "192.168.1.0/24"`,
				`$ stackit network-area region create --network-area-id xxx --region eu02 --organization-id yyy --ipv4-network-ranges 192.168.0.0/24 --ipv4-transfer-network 192.168.1.0/24`,
			),
			examples.NewExample(
				`Create a new regional configuration "eu02" for a STACKIT Network Area with ID "xxx" in organization with ID "yyy", using the set region config`,
				`$ stackit config set --region eu02`,
				`$ stackit network-area region create --network-area-id xxx --organization-id yyy --ipv4-network-ranges 192.168.0.0/24 --ipv4-transfer-network 192.168.1.0/24`,
			),
			examples.NewExample(
				`Create a new regional configuration for a STACKIT Network Area with ID "xxx" in organization with ID "yyy", ipv4 network range "192.168.0.0/24", ipv4 transfer network "192.168.1.0/24", default prefix length "24", max prefix length "25" and min prefix length "20"`,
				`$ stackit network-area region create --network-area-id xxx --organization-id yyy --ipv4-network-ranges 192.168.0.0/24 --ipv4-transfer-network 192.168.1.0/24 --region "eu02" --ipv4-default-prefix-length 24 --ipv4-max-prefix-length 25 --ipv4-min-prefix-length 20`,
			),
			examples.NewExample(
				`Create a new regional configuration for a STACKIT Network Area with ID "xxx" in organization with ID "yyy", ipv4 network range "192.168.0.0/24", ipv4 transfer network "192.168.1.0/24", default prefix length "24", max prefix length "25" and min prefix length "20"`,
				`$ stackit network-area region create --network-area-id xxx --organization-id yyy --ipv4-network-ranges 192.168.0.0/24 --ipv4-transfer-network 192.168.1.0/24 --region "eu02" --ipv4-default-prefix-length 24 --ipv4-max-prefix-length 25 --ipv4-min-prefix-length 20`,
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

			// Get network area label
			networkAreaLabel, err := iaasUtils.GetNetworkAreaName(ctx, apiClient, model.OrganizationId, model.NetworkAreaId)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get network area name: %v", err)
				networkAreaLabel = model.NetworkAreaId
			}

			prompt := fmt.Sprintf("Are you sure you want to create the regional configuration %q for STACKIT Network Area (SNA) %q?", model.Region, networkAreaLabel)
			err = params.Printer.PromptForConfirmation(prompt)
			if err != nil {
				return err
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("create network area region: %w", err)
			}

			if resp == nil || resp.Ipv4 == nil {
				return fmt.Errorf("empty response from API")
			}

			if !model.Async {
				s := spinner.New(params.Printer)
				s.Start("Create network area region")
				_, err = wait.CreateNetworkAreaRegionWaitHandler(ctx, apiClient, model.OrganizationId, model.NetworkAreaId, model.Region).WaitWithContext(ctx)
				if err != nil {
					return fmt.Errorf("wait for network area region creation: %w", err)
				}
				s.Stop()
			}

			return outputResult(params.Printer, model.OutputFormat, model.Region, networkAreaLabel, *resp)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), networkAreaIdFlag, "STACKIT Network Area (SNA) ID")
	cmd.Flags().Var(flags.UUIDFlag(), organizationIdFlag, "Organization ID")
	cmd.Flags().Var(flags.CIDRSliceFlag(), ipv4NetworkRangesFlag, "Network range to create in CIDR notation")
	cmd.Flags().Var(flags.CIDRFlag(), ipv4TransferNetworkFlag, "Transfer network in CIDR notation")
	cmd.Flags().StringSlice(ipv4DefaultNameservers, nil, "List of default DNS name server IPs")
	cmd.Flags().Int64(ipv4DefaultPrefixLengthFlag, 0, "The default prefix length for networks in the network area")
	cmd.Flags().Int64(ipv4MaxPrefixLengthFlag, 0, "The maximum prefix length for networks in the network area")
	cmd.Flags().Int64(ipv4MinPrefixLengthFlag, 0, "The minimum prefix length for networks in the network area")

	err := flags.MarkFlagsRequired(cmd, networkAreaIdFlag, organizationIdFlag, ipv4NetworkRangesFlag, ipv4TransferNetworkFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command, _ []string) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.Region == "" {
		return nil, &errors.RegionError{}
	}

	model := inputModel{
		GlobalFlagModel:         globalFlags,
		NetworkAreaId:           flags.FlagToStringValue(p, cmd, networkAreaIdFlag),
		OrganizationId:          flags.FlagToStringValue(p, cmd, organizationIdFlag),
		IPv4DefaultNameservers:  flags.FlagToStringSlicePointer(p, cmd, ipv4DefaultNameservers),
		IPv4DefaultPrefixLength: flags.FlagToInt64Pointer(p, cmd, ipv4DefaultPrefixLengthFlag),
		IPv4MaxPrefixLength:     flags.FlagToInt64Pointer(p, cmd, ipv4MaxPrefixLengthFlag),
		IPv4MinPrefixLength:     flags.FlagToInt64Pointer(p, cmd, ipv4MinPrefixLengthFlag),
		IPv4NetworkRanges:       flags.FlagToStringSliceValue(p, cmd, ipv4NetworkRangesFlag),
		IPv4TransferNetwork:     flags.FlagToStringValue(p, cmd, ipv4TransferNetworkFlag),
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *iaas.APIClient) iaas.ApiCreateNetworkAreaRegionRequest {
	req := apiClient.CreateNetworkAreaRegion(ctx, model.OrganizationId, model.NetworkAreaId, model.Region)

	var networkRange []iaas.NetworkRange
	if len(model.IPv4NetworkRanges) > 0 {
		networkRange = make([]iaas.NetworkRange, len(model.IPv4NetworkRanges))
		for i := range model.IPv4NetworkRanges {
			networkRange[i] = iaas.NetworkRange{
				Prefix: utils.Ptr(model.IPv4NetworkRanges[i]),
			}
		}
	}

	payload := iaas.CreateNetworkAreaRegionPayload{
		Ipv4: &iaas.RegionalAreaIPv4{
			DefaultNameservers: model.IPv4DefaultNameservers,
			DefaultPrefixLen:   model.IPv4DefaultPrefixLength,
			MaxPrefixLen:       model.IPv4MaxPrefixLength,
			MinPrefixLen:       model.IPv4MinPrefixLength,
			NetworkRanges:      utils.Ptr(networkRange),
			TransferNetwork:    utils.Ptr(model.IPv4TransferNetwork),
		},
	}
	return req.CreateNetworkAreaRegionPayload(payload)
}

func outputResult(p *print.Printer, outputFormat, region, networkAreaLabel string, regionalArea iaas.RegionalArea) error {
	return p.OutputResult(outputFormat, regionalArea, func() error {
		p.Outputf("Create region configuration for SNA %q.\nRegion: %s\n", networkAreaLabel, region)
		return nil
	})
}
