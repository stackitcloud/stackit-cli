package list

import (
	"cmp"
	"context"
	"fmt"
	"slices"

	"github.com/stackitcloud/stackit-cli/internal/pkg/projectname"
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
	limitFlag         = "limit"
	labelSelectorFlag = "label-selector"
	networkIdFlag     = "network-id"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	Limit         *int64
	LabelSelector *string
	NetworkId     *string
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Lists all network interfaces of a network",
		Long:  "Lists all network interfaces of a network.",
		Args:  args.NoArgs,
		Example: examples.Build(
			// Note: this subcommand uses two different API enpoints, which makes the implementation somewhat messy
			examples.NewExample(
				`Lists all network interfaces`,
				`$ stackit network-interface list`,
			),
			examples.NewExample(
				`Lists all network interfaces with network ID "xxx"`,
				`$ stackit network-interface list --network-id xxx`,
			),
			examples.NewExample(
				`Lists all network interfaces with network ID "xxx" which contains the label xxx`,
				`$ stackit network-interface list --network-id xxx --label-selector xxx`,
			),
			examples.NewExample(
				`Lists all network interfaces with network ID "xxx" in JSON format`,
				`$ stackit network-interface list --network-id xxx --output-format json`,
			),
			examples.NewExample(
				`Lists up to 10 network interfaces with network ID "xxx"`,
				`$ stackit network-interface list --network-id xxx --limit 10`,
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

			if model.NetworkId == nil {
				// Call API to get all NICs in the Project
				req := buildProjectRequest(ctx, model, apiClient)

				resp, err := req.Execute()
				if err != nil {
					return fmt.Errorf("list network interfaces: %w", err)
				}

				projectLabel, err := projectname.GetProjectName(ctx, params.Printer, params.CliVersion, cmd)
				if err != nil {
					projectLabel = model.ProjectId
				}

				// Truncate output
				items := utils.GetSliceFromPointer(resp.Items)
				if model.Limit != nil && len(items) > int(*model.Limit) {
					items = items[:*model.Limit]
				}

				return outputProjectResult(params.Printer, model.OutputFormat, items, projectLabel)
			}

			// Call API to get NICs for one Network
			req := buildNetworkRequest(ctx, model, apiClient)

			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("list network interfaces: %w", err)
			}

			networkLabel, err := iaasUtils.GetNetworkName(ctx, apiClient, model.ProjectId, model.Region, *model.NetworkId)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get network name: %v", err)
				networkLabel = *model.NetworkId
			} else if networkLabel == "" {
				networkLabel = *model.NetworkId
			}

			// Truncate output
			items := utils.GetSliceFromPointer(resp.Items)
			if model.Limit != nil && len(items) > int(*model.Limit) {
				items = items[:*model.Limit]
			}

			return outputNetworkResult(params.Printer, model.OutputFormat, items, networkLabel)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), networkIdFlag, "Network ID")
	cmd.Flags().Int64(limitFlag, 0, "Maximum number of entries to list")
	cmd.Flags().String(labelSelectorFlag, "", "Filter by label")
}

func parseInput(p *print.Printer, cmd *cobra.Command, _ []string) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	limit := flags.FlagToInt64Pointer(p, cmd, limitFlag)
	if limit != nil && *limit < 1 {
		return nil, &errors.FlagValidationError{
			Flag:    limitFlag,
			Details: "must be greater than 0",
		}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		Limit:           limit,
		LabelSelector:   flags.FlagToStringPointer(p, cmd, labelSelectorFlag),
		NetworkId:       flags.FlagToStringPointer(p, cmd, networkIdFlag),
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildProjectRequest(ctx context.Context, model *inputModel, apiClient *iaas.APIClient) iaas.ApiListProjectNICsRequest {
	req := apiClient.ListProjectNICs(ctx, model.ProjectId, model.Region)
	if model.LabelSelector != nil {
		req = req.LabelSelector(*model.LabelSelector)
	}

	return req
}

func buildNetworkRequest(ctx context.Context, model *inputModel, apiClient *iaas.APIClient) iaas.ApiListNicsRequest {
	req := apiClient.ListNics(ctx, model.ProjectId, model.Region, *model.NetworkId)
	if model.LabelSelector != nil {
		req = req.LabelSelector(*model.LabelSelector)
	}

	return req
}

func outputProjectResult(p *print.Printer, outputFormat string, nics []iaas.NIC, projectLabel string) error {
	return p.OutputResult(outputFormat, nics, func() error {
		if len(nics) == 0 {
			p.Outputf("No network interfaces found for project %q\n", projectLabel)
			return nil
		}

		slices.SortFunc(nics, func(a, b iaas.NIC) int {
			return cmp.Compare(utils.PtrValue(a.NetworkId), utils.PtrValue(b.NetworkId))
		})

		table := tables.NewTable()
		table.SetHeader("ID", "NAME", "NETWORK ID", "NIC SECURITY", "DEVICE ID", "IPv4 ADDRESS", "STATUS", "TYPE")

		for _, nic := range nics {
			table.AddRow(
				utils.PtrString(nic.Id),
				utils.PtrString(nic.Name),
				utils.PtrString(nic.NetworkId),
				utils.PtrString(nic.NicSecurity),
				utils.PtrString(nic.Device),
				utils.PtrString(nic.Ipv4),
				utils.PtrString(nic.Status),
				utils.PtrString(nic.Type),
			)
			table.AddSeparator()
		}

		p.Outputln(table.Render())
		return nil
	})
}

func outputNetworkResult(p *print.Printer, outputFormat string, nics []iaas.NIC, networkLabel string) error {
	return p.OutputResult(outputFormat, nics, func() error {
		if len(nics) == 0 {
			p.Outputf("No network interfaces found for network %q\n", networkLabel)
			return nil
		}

		table := tables.NewTable()
		table.SetHeader("ID", "NAME", "NIC SECURITY", "DEVICE ID", "IPv4 ADDRESS", "STATUS", "TYPE")

		for _, nic := range nics {
			table.AddRow(
				utils.PtrString(nic.Id),
				utils.PtrString(nic.Name),
				utils.PtrString(nic.NicSecurity),
				utils.PtrString(nic.Device),
				utils.PtrString(nic.Ipv4),
				utils.PtrString(nic.Status),
				utils.PtrString(nic.Type),
			)
			table.AddSeparator()
		}

		p.Outputln(table.Render())
		return nil
	})
}
