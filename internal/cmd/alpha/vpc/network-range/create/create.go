package create

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	iaas "github.com/stackitcloud/stackit-sdk-go/services/iaas/v2alpha1api"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/iaasalpha/client"
	iaasUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/iaasalpha/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
)

const (
	vpcIdFlag               = "vpc-id"
	prefixFlag              = "prefix"
	descriptionFlag         = "description"
	defaultPrefixLengthFlag = "default-prefix-length"
	maxPrefixLengthFlag     = "max-prefix-length"
	minPrefixLengthFlag     = "min-prefix-length"
	nameserversFlag         = "nameservers"
	labelsFlag              = "labels"
)

var (
	ipVersionFlag = flags.StringEnumFlag(
		"ip-version",
		iaas.AllowedNetworkRangeIPv4RequestIpVersionEnumValues,
		"IP version of the network-range",
	)
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	VpcId               string
	Prefix              string
	IpVersion           iaas.NetworkRangeIPv4RequestIpVersion
	Description         *string
	DefaultPrefixLength *int64
	MaxPrefixLength     *int64
	MinPrefixLength     *int64
	Nameservers         []string
	Labels              map[string]any
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Creates a regional network range in a VPC",
		Long: fmt.Sprintf("%s\n%s",
			"Creates a regional network range in a VPC.",
			"The vpc needs to be enabled in the specified region before creating a network range.",
		),
		Args: args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Create a network range in a VPC with ID "xxx" with prefix "1.1.1.0/24" and with ip version 4`,
				`$ stackit alpha vpc network-range create --vpc-id xxx --prefix "1.1.1.0/24" --ip-version ipv4`,
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

			// Get vpc label
			vpcLabel, err := iaasUtils.GetVPCName(ctx, apiClient.DefaultAPI, model.ProjectId, model.VpcId)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get vpc name: %v", err)
				vpcLabel = model.VpcId
			}

			prompt := fmt.Sprintf("Are you sure you want to create a network range for VPC %q?", vpcLabel)
			err = params.Printer.PromptForConfirmation(prompt)
			if err != nil {
				return err
			}

			// Call API
			req, err := buildRequest(ctx, model, apiClient)
			if err != nil {
				return err
			}
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("create network range: %w", err)
			}

			return outputResult(params.Printer, model.OutputFormat, vpcLabel, resp)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), vpcIdFlag, "VPC ID")
	cmd.Flags().Var(flags.CIDRFlag(), prefixFlag, "Network range to create in CIDR notation")
	cmd.Flags().String(descriptionFlag, "", "Description of the network range")
	ipVersionFlag.Register(cmd.Flags())
	cmd.Flags().Int64(defaultPrefixLengthFlag, 0, "The default prefix length for network ranges in the VPC")
	cmd.Flags().Int64(maxPrefixLengthFlag, 0, "The maximal prefix length for network ranges in the VPC")
	cmd.Flags().Int64(minPrefixLengthFlag, 0, "The minimal prefix length for network ranges in the VPC")
	cmd.Flags().StringSlice(nameserversFlag, nil, "A list containing DNS Servers")
	cmd.Flags().StringToString(labelsFlag, nil, "Labels are key-value string pairs which can be attached to a network range. E.g. '--labels key1=value1,key2=value2,...'")

	err := flags.MarkFlagsRequired(cmd, vpcIdFlag, prefixFlag, ipVersionFlag.Name())
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command, _ []string) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel:     globalFlags,
		VpcId:               flags.FlagToStringValue(p, cmd, vpcIdFlag),
		Prefix:              flags.FlagToStringValue(p, cmd, prefixFlag),
		IpVersion:           ipVersionFlag.Get(),
		Description:         flags.FlagToStringPointer(p, cmd, descriptionFlag),
		DefaultPrefixLength: flags.FlagToInt64Pointer(p, cmd, defaultPrefixLengthFlag),
		MaxPrefixLength:     flags.FlagToInt64Pointer(p, cmd, maxPrefixLengthFlag),
		MinPrefixLength:     flags.FlagToInt64Pointer(p, cmd, minPrefixLengthFlag),
		Nameservers:         flags.FlagToStringSliceValue(p, cmd, nameserversFlag),
		Labels:              flags.FlagToStringToAny(p, cmd, labelsFlag),
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *iaas.APIClient) (iaas.ApiCreateVPCNetworkRangeRequest, error) {
	if model.IpVersion != "ipv4" {
		return iaas.ApiCreateVPCNetworkRangeRequest{}, fmt.Errorf("unsupported IP version: %q", model.IpVersion)
	}

	payload := iaas.CreateVPCNetworkRangePayload{
		NetworkRangeIPv4Request: &iaas.NetworkRangeIPv4Request{
			DefaultPrefixLen: model.DefaultPrefixLength,
			Description:      model.Description,
			Prefix:           model.Prefix,
			IpVersion:        model.IpVersion,
			MaxPrefixLen:     model.MaxPrefixLength,
			MinPrefixLen:     model.MinPrefixLength,
			Nameservers:      model.Nameservers,
			Labels:           model.Labels,
		},
	}
	return apiClient.DefaultAPI.CreateVPCNetworkRange(ctx, model.ProjectId, model.VpcId, model.Region).CreateVPCNetworkRangePayload(payload), nil
}

func outputResult(p *print.Printer, outputFormat, vpcLabel string, networkRange *iaas.VPCNetworkRange) error {
	if networkRange == nil {
		return fmt.Errorf("create vpc network range response is nil")
	}
	return p.OutputResult(outputFormat, networkRange, func() error {
		var networkRangeId string
		if ipv4 := networkRange.VPCNetworkRangeIPv4; ipv4 != nil && ipv4.Id != nil {
			networkRangeId = *ipv4.Id
		} else if ipv6 := networkRange.VPCNetworkRangeIPv6; ipv6 != nil && ipv6.Id != nil {
			networkRangeId = *ipv6.Id
		} else {
			return fmt.Errorf("network range ID is missing in API response")
		}

		p.Outputf("Created network range for VPC %q.\nNetwork range ID: %s\n", vpcLabel, networkRangeId)
		return nil
	})
}
