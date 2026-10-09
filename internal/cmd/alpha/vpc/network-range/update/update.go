package update

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
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

const (
	networkRangeIdArg = "NETWORK_RANGE_ID"

	vpcIdFlag               = "vpc-id"
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
		iaas.AllowedV1UpdateVPCNetworkRangeIPv4IpVersionEnumValues,
		"IP version of the network-range",
	)
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	NetworkRangeId      string
	VpcId               string
	IpVersion           iaas.V1UpdateVPCNetworkRangeIPv4IpVersion
	Description         *string
	DefaultPrefixLength *int64
	MaxPrefixLength     *int64
	MinPrefixLength     *int64
	Nameservers         []string
	Labels              map[string]any
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("update %s", networkRangeIdArg),
		Short: "Updates a regional network range in a VPC",
		Long:  "Updates a regional network range in a VPC.",
		Args:  args.SingleArg(networkRangeIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Update a network range with ID "xxx" in a VPC with ID "yyy" and with ip version "ipv4" to new nameservers "1.2.3.4,5.6.7.8"`,
				`$ stackit alpha vpc network-range update xxx --vpc-id yyy --ip-version ipv4 --nameservers "1.2.3.4,5.6.7.8"`,
			),
			examples.NewExample(
				`Update a network range with ID "xxx" in a VPC with ID "yyy" and with ip version "ipv4" to new description "updated network range" and default prefix length 24`,
				`$ stackit alpha vpc network-range update xxx --vpc-id yyy --ip-version ipv4 --description "updated network range" --default-prefix-length 24`,
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

			networkRangeLabel, err := iaasUtils.GetVPCNetworkRangePrefix(ctx, apiClient.DefaultAPI, model.ProjectId, model.VpcId, model.Region, model.NetworkRangeId)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get vpc network range prefix: %v", err)
				networkRangeLabel = model.NetworkRangeId
			} else if networkRangeLabel == "" {
				networkRangeLabel = model.NetworkRangeId
			}

			prompt := fmt.Sprintf("Are you sure you want to update the network range %q for VPC %q?", networkRangeLabel, vpcLabel)
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
				return fmt.Errorf("update network range: %w", err)
			}

			return outputResult(params.Printer, model.OutputFormat, vpcLabel, resp)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), vpcIdFlag, "VPC ID")
	cmd.Flags().String(descriptionFlag, "", "Description of the network range")
	ipVersionFlag.Register(cmd.Flags())
	cmd.Flags().Int64(defaultPrefixLengthFlag, 0, "The default prefix length for network ranges in the VPC")
	cmd.Flags().Int64(maxPrefixLengthFlag, 0, "The maximal prefix length for network ranges in the VPC")
	cmd.Flags().Int64(minPrefixLengthFlag, 0, "The minimal prefix length for network ranges in the VPC")
	cmd.Flags().StringSlice(nameserversFlag, nil, "A list containing DNS Servers")
	cmd.Flags().StringToString(labelsFlag, nil, "Labels are key-value string pairs which can be attached to a network range. E.g. '--labels key1=value1,key2=value2,...'")

	cmd.MarkFlagsOneRequired(
		descriptionFlag,
		defaultPrefixLengthFlag,
		maxPrefixLengthFlag,
		minPrefixLengthFlag,
		nameserversFlag,
		labelsFlag,
	)
	cobra.CheckErr(flags.MarkFlagsRequired(cmd, vpcIdFlag, ipVersionFlag.Name()))
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	networkRangeId := inputArgs[0]

	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel:     globalFlags,
		NetworkRangeId:      networkRangeId,
		VpcId:               flags.FlagToStringValue(p, cmd, vpcIdFlag),
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

func buildRequest(ctx context.Context, model *inputModel, apiClient *iaas.APIClient) (iaas.ApiUpdateVPCNetworkRangeRequest, error) {
	if model.IpVersion != "ipv4" {
		return iaas.ApiUpdateVPCNetworkRangeRequest{}, fmt.Errorf("unsupported IP version: %q", model.IpVersion)
	}

	payload := iaas.UpdateVPCNetworkRangePayload{
		V1UpdateVPCNetworkRangeIPv4: &iaas.V1UpdateVPCNetworkRangeIPv4{
			DefaultPrefixLen: model.DefaultPrefixLength,
			Description:      model.Description,
			IpVersion:        model.IpVersion,
			MaxPrefixLen:     model.MaxPrefixLength,
			MinPrefixLen:     model.MinPrefixLength,
			Nameservers:      model.Nameservers,
			Labels:           model.Labels,
		},
	}
	return apiClient.DefaultAPI.UpdateVPCNetworkRange(ctx, model.ProjectId, model.VpcId, model.Region, model.NetworkRangeId).UpdateVPCNetworkRangePayload(payload), nil
}

func outputResult(p *print.Printer, outputFormat, vpcLabel string, networkRange *iaas.VPCNetworkRange) error {
	if networkRange == nil {
		return fmt.Errorf("update vpc network range response is nil")
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

		p.Outputf("Updated network range for VPC %q.\nNetwork range ID: %s\n", vpcLabel, networkRangeId)
		return nil
	})
}
