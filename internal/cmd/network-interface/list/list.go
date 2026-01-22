package list

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
	NetworkId     string
}

type ExecutableRequest interface {
	Execute() (*iaas.NICListResponse, error)
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Lists all network interfaces of a network",
		Long:  "Lists all network interfaces of a network.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Lists all network interfaces in your current project`,
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

			// Call API
			var req ExecutableRequest
			var networkLabel = ""
			if model.NetworkId == "" {
				// Return all NICs in the Project
				req = buildProjectRequest(ctx, model, apiClient)
			} else {
				// Return the NICs for one Network
				req = buildRequest(ctx, model, apiClient)

				networkLabel, err = iaasUtils.GetNetworkName(ctx, apiClient, model.ProjectId, model.Region, model.NetworkId)
				if err != nil {
					params.Printer.Debug(print.ErrorLevel, "get network name: %v", err)
					networkLabel = model.NetworkId
				} else if networkLabel == "" {
					networkLabel = model.NetworkId
				}
			}
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("list network interfaces: %w", err)
			}

			// Truncate output
			items := *resp.Items
			if model.Limit != nil && len(items) > int(*model.Limit) {
				items = items[:*model.Limit]
			}

			return outputResult(params.Printer, model.OutputFormat, items, networkLabel)
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
		NetworkId:       flags.FlagToStringValue(p, cmd, networkIdFlag),
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildProjectRequest(ctx context.Context, model *inputModel, apiClient *iaas.APIClient) ExecutableRequest {
	req := apiClient.ListProjectNICs(ctx, model.ProjectId, model.Region)
	if model.LabelSelector != nil {
		req = req.LabelSelector(*model.LabelSelector)
	}

	return req
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *iaas.APIClient) ExecutableRequest {
	req := apiClient.ListNics(ctx, model.ProjectId, model.Region, model.NetworkId)
	if model.LabelSelector != nil {
		req = req.LabelSelector(*model.LabelSelector)
	}

	return req
}

func outputResult(p *print.Printer, outputFormat string, nics []iaas.NIC, networkLabel string) error {
	return p.OutputResult(outputFormat, nics, func() error {
		if len(nics) == 0 {
			if networkLabel == "" {
				p.Outputf("No network interfaces found for your current project\n")
			} else {
				p.Outputf("No network interfaces found for network %q\n", networkLabel)
			}
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
