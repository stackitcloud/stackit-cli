package list

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/iaas/client"
	iaasUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/iaas/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/iaas"
)

const (
	networkAreaIdFlag  = "network-area-id"
	organizationIdFlag = "organization-id"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	OrganizationId string
	NetworkAreaId  string
}

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Lists all configured regions for a STACKIT Network Area (SNA)",
		Long:  "Lists all configured regions for a STACKIT Network Area (SNA).",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`List all configured region for a STACKIT Network Area with ID "xxx" in organization with ID "yyy"`,
				`$ stackit network-area region list --network-area-id xxx --organization-id yyy`,
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

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("list network area region: %w", err)
			}

			if resp == nil {
				return fmt.Errorf("empty response from API")
			}

			return outputResult(params.Printer, model.OutputFormat, networkAreaLabel, *resp)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), networkAreaIdFlag, "STACKIT Network Area (SNA) ID")
	cmd.Flags().Var(flags.UUIDFlag(), organizationIdFlag, "Organization ID")

	err := flags.MarkFlagsRequired(cmd, networkAreaIdFlag, organizationIdFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command, _ []string) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)

	model := inputModel{
		GlobalFlagModel: globalFlags,
		NetworkAreaId:   flags.FlagToStringValue(p, cmd, networkAreaIdFlag),
		OrganizationId:  flags.FlagToStringValue(p, cmd, organizationIdFlag),
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *iaas.APIClient) iaas.ApiListNetworkAreaRegionsRequest {
	return apiClient.ListNetworkAreaRegions(ctx, model.OrganizationId, model.NetworkAreaId)
}

func outputResult(p *print.Printer, outputFormat, areaLabel string, regionalArea iaas.RegionalAreaListResponse) error {
	return p.OutputResult(outputFormat, regionalArea, func() error {
		if regionalArea.Regions == nil || len(*regionalArea.Regions) == 0 {
			p.Outputf("No regions found for network area %q\n", areaLabel)
			return nil
		}

		table := tables.NewTable()
		table.SetHeader("REGION", "STATUS", "DNS NAME SERVERS", "NETWORK RANGES", "TRANSFER NETWORK")
		for region, regionConfig := range *regionalArea.Regions {
			var dnsNames string
			var networkRanges []string
			var transferNetwork string

			if ipv4 := regionConfig.Ipv4; ipv4 != nil {
				// Set dnsNames
				dnsNames = utils.JoinStringPtr(ipv4.DefaultNameservers, ",")

				// Set networkRanges
				if ipv4.NetworkRanges != nil && len(*ipv4.NetworkRanges) > 0 {
					for _, networkRange := range *ipv4.NetworkRanges {
						if networkRange.Prefix != nil {
							networkRanges = append(networkRanges, *networkRange.Prefix)
						}
					}
				}

				// Set transferNetwork
				transferNetwork = utils.PtrString(ipv4.TransferNetwork)
			}

			table.AddRow(
				region,
				utils.PtrString(regionConfig.Status),
				dnsNames,
				strings.Join(networkRanges, ","),
				transferNetwork,
			)
		}

		if err := table.Display(p); err != nil {
			return fmt.Errorf("render table: %w", err)
		}
		return nil
	})
}
