package describe

import (
	"context"
	"fmt"
	"strings"

	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"

	iaas "github.com/stackitcloud/stackit-sdk-go/services/iaas/v2alpha1api"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/iaasalpha/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

const (
	networkRangeIdArg = "NETWORK_RANGE_ID"

	vpcIdFlag = "vpc-id"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	VpcId          string
	NetworkRangeId string
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("describe %s", networkRangeIdArg),
		Short: "Shows details of a network range in a VPC",
		Long:  "Shows details of a network range in a VPC.",
		Args:  args.SingleArg(networkRangeIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Show details of a network range with ID "xxx" in a VPC with ID "yyy"`,
				`$ stackit alpha vpc network-range describe xxx --vpc-id yyy`,
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

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("describe network range: %w", err)
			}

			return outputResult(params.Printer, model.OutputFormat, resp)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), vpcIdFlag, "VPC ID")

	err := flags.MarkFlagsRequired(cmd, vpcIdFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	networkRangeId := inputArgs[0]
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		VpcId:           flags.FlagToStringValue(p, cmd, vpcIdFlag),
		NetworkRangeId:  networkRangeId,
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *iaas.APIClient) iaas.ApiGetVPCNetworkRangeRequest {
	return apiClient.DefaultAPI.GetVPCNetworkRange(ctx, model.ProjectId, model.VpcId, model.Region, model.NetworkRangeId)
}

func outputResult(p *print.Printer, outputFormat string, networkRange *iaas.VPCNetworkRange) error {
	if networkRange == nil {
		return fmt.Errorf("network range is nil")
	}

	return p.OutputResult(outputFormat, networkRange, func() error {
		table := tables.NewTable()

		// iaas.VPCNetworkRange is a oneOf, either IPv4 or IPv6 is set
		if ipv4 := networkRange.VPCNetworkRangeIPv4; ipv4 != nil {
			table.AddRow("ID", utils.PtrString(ipv4.Id))
			table.AddSeparator()
			table.AddRow("PREFIX", ipv4.Prefix)
			table.AddSeparator()
			table.AddRow("DESCRIPTION", utils.PtrString(ipv4.Description))
			table.AddSeparator()
			table.AddRow("NAMESERVERS", strings.Join(ipv4.Nameservers, ","))
			table.AddSeparator()
			table.AddRow("DEFAULT PREFIX LENGTH", utils.PtrString(ipv4.DefaultPrefixLen))
			table.AddSeparator()
			table.AddRow("MAX PREFIX LENGTH", utils.PtrString(ipv4.MaxPrefixLen))
			table.AddSeparator()
			table.AddRow("MIN PREFIX LENGTH", utils.PtrString(ipv4.MinPrefixLen))
			table.AddSeparator()
			table.AddRow("IP VERSION", ipv4.IpVersion)
			table.AddSeparator()
			table.AddRow("STATUS", utils.PtrString(ipv4.Status))
			table.AddSeparator()

			var labels []string
			for key, value := range ipv4.Labels {
				labels = append(labels, fmt.Sprintf("%s: %s", key, value))
			}
			table.AddRow("LABELS", strings.Join(labels, "\n"))
			table.AddSeparator()
			table.AddRow("CREATED AT", utils.ConvertTimePToDateTimeString(ipv4.CreatedAt))
			table.AddSeparator()
			table.AddRow("UPDATED AT", utils.ConvertTimePToDateTimeString(ipv4.UpdatedAt))
			table.AddSeparator()
		} else if ipv6 := networkRange.VPCNetworkRangeIPv6; ipv6 != nil {
			table.AddRow("ID", utils.PtrString(ipv6.Id))
			table.AddSeparator()
			table.AddRow("PREFIX", ipv6.Prefix)
			table.AddSeparator()
			table.AddRow("DESCRIPTION", utils.PtrString(ipv6.Description))
			table.AddSeparator()
			table.AddRow("NAMESERVERS", strings.Join(ipv6.Nameservers, ","))
			table.AddSeparator()
			table.AddRow("DEFAULT PREFIX LENGTH", utils.PtrString(ipv6.DefaultPrefixLen))
			table.AddSeparator()
			table.AddRow("MAX PREFIX LENGTH", utils.PtrString(ipv6.MaxPrefixLen))
			table.AddSeparator()
			table.AddRow("MIN PREFIX LENGTH", utils.PtrString(ipv6.MinPrefixLen))
			table.AddSeparator()
			table.AddRow("IP VERSION", ipv6.IpVersion)
			table.AddSeparator()
			table.AddRow("STATUS", utils.PtrString(ipv6.Status))
			table.AddSeparator()

			var labels []string
			for key, value := range ipv6.Labels {
				labels = append(labels, fmt.Sprintf("%s: %s", key, value))
			}
			table.AddRow("LABELS", strings.Join(labels, "\n"))
			table.AddSeparator()
			table.AddRow("CREATED AT", utils.ConvertTimePToDateTimeString(ipv6.CreatedAt))
			table.AddSeparator()
			table.AddRow("UPDATED AT", utils.ConvertTimePToDateTimeString(ipv6.UpdatedAt))
			table.AddSeparator()
		} else {
			return fmt.Errorf("unsupported or empty network range payload")
		}

		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}
		return nil
	})
}
