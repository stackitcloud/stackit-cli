package describe

import (
	"context"
	"fmt"
	"strings"

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

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "describe",
		Short: "Describes a regional configuration for a STACKIT Network Area (SNA)",
		Long:  "Describes a regional configuration for a STACKIT Network Area (SNA).",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Describe a regional configuration "eu02" for a STACKIT Network Area with ID "xxx" in organization with ID "yyy"`,
				`$ stackit network-area region describe --network-area-id xxx --region eu02 --organization-id yyy`,
			),
			examples.NewExample(
				`Describe a regional configuration "eu02" for a STACKIT Network Area with ID "xxx" in organization with ID "yyy", using the set region config`,
				`$ stackit config set --region eu02`,
				`$ stackit network-area region describe --network-area-id xxx --organization-id yyy`,
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
			networkAreaName, err := iaasUtils.GetNetworkAreaName(ctx, apiClient, model.OrganizationId, model.NetworkAreaId)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get network area name: %v", err)
				// Set explicit the networkAreaName to empty string and not to the ID, because this is used for the table output
				networkAreaName = ""
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("describe network area region: %w", err)
			}

			if resp == nil || resp.Ipv4 == nil {
				return fmt.Errorf("empty response from API")
			}

			return outputResult(params.Printer, model.OutputFormat, model.Region, model.NetworkAreaId, networkAreaName, *resp)
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
	if globalFlags.Region == "" {
		return nil, &errors.RegionError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		NetworkAreaId:   flags.FlagToStringValue(p, cmd, networkAreaIdFlag),
		OrganizationId:  flags.FlagToStringValue(p, cmd, organizationIdFlag),
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *iaas.APIClient) iaas.ApiGetNetworkAreaRegionRequest {
	return apiClient.GetNetworkAreaRegion(ctx, model.OrganizationId, model.NetworkAreaId, model.Region)
}

func outputResult(p *print.Printer, outputFormat, region, areaId, areaName string, regionalArea iaas.RegionalArea) error {
	return p.OutputResult(outputFormat, regionalArea, func() error {
		table := tables.NewTable()
		table.AddRow("ID", areaId)
		table.AddSeparator()
		if areaName != "" {
			table.AddRow("NAME", areaName)
			table.AddSeparator()
		}
		table.AddRow("REGION", region)
		table.AddSeparator()
		table.AddRow("STATUS", utils.PtrString(regionalArea.Status))
		table.AddSeparator()
		if ipv4 := regionalArea.Ipv4; ipv4 != nil {
			if ipv4.NetworkRanges != nil {
				var networkRanges []string
				for _, networkRange := range *ipv4.NetworkRanges {
					if networkRange.Prefix != nil {
						networkRanges = append(networkRanges, *networkRange.Prefix)
					}
				}
				table.AddRow("NETWORK RANGES", strings.Join(networkRanges, ","))
				table.AddSeparator()
			}
			if transferNetwork := ipv4.TransferNetwork; transferNetwork != nil {
				table.AddRow("TRANSFER RANGE", utils.PtrString(transferNetwork))
				table.AddSeparator()
			}
			if defaultNameserver := ipv4.DefaultNameservers; defaultNameserver != nil && len(*defaultNameserver) > 0 {
				table.AddRow("DNS NAME SERVERS", strings.Join(*defaultNameserver, ","))
				table.AddSeparator()
			}
			if defaultPrefixLength := ipv4.DefaultPrefixLen; defaultPrefixLength != nil {
				table.AddRow("DEFAULT PREFIX LENGTH", utils.PtrString(defaultPrefixLength))
				table.AddSeparator()
			}
			if maxPrefixLength := ipv4.MaxPrefixLen; maxPrefixLength != nil {
				table.AddRow("MAX PREFIX LENGTH", utils.PtrString(maxPrefixLength))
				table.AddSeparator()
			}
			if minPrefixLen := ipv4.MinPrefixLen; minPrefixLen != nil {
				table.AddRow("MIN PREFIX LENGTH", utils.PtrString(minPrefixLen))
				table.AddSeparator()
			}
		}

		if err := table.Display(p); err != nil {
			return fmt.Errorf("render table: %w", err)
		}
		return nil
	})
}
