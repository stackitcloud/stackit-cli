package list

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/goccy/go-yaml"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/iaas/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-sdk-go/services/iaas"

	"github.com/spf13/cobra"
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

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Lists all network interfaces of a network",
		Long:  "Lists all network interfaces of a network.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Lists all network interfaces with network ID "xxx"`,
				`$ stackit beta network-interface list --network-id xxx`,
			),
			examples.NewExample(
				`Lists all network interfaces with network ID "xxx" which contains the label xxx`,
				`$ stackit beta network-interface list --network-id xxx --label-selector xxx`,
			),
			examples.NewExample(
				`Lists all network interfaces with network ID "xxx" in JSON format`,
				`$ stackit beta network-interface list --network-id xxx --output-format json`,
			),
			examples.NewExample(
				`Lists up to 10 network interfaces with network ID "xxx"`,
				`$ stackit beta network-interface list --network-id xxx --limit 10`,
			),
		),
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := context.Background()
			model, err := parseInput(p, cmd)
			if err != nil {
				return err
			}

			// Configure API client
			apiClient, err := client.ConfigureClient(p)
			if err != nil {
				return err
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("list network interfaces: %w", err)
			}

			if resp.Items == nil || len(*resp.Items) == 0 {
				p.Info("No network interfaces found for network %d\n", model.NetworkId)
				return nil
			}

			// Truncate output
			items := *resp.Items
			if model.Limit != nil && len(items) > int(*model.Limit) {
				items = items[:*model.Limit]
			}

			return outputResult(p, model.OutputFormat, items)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), networkIdFlag, "Network ID")
	cmd.Flags().Int64(limitFlag, 0, "Maximum number of entries to list")
	cmd.Flags().String(labelSelectorFlag, "", "Filter by label")

	err := flags.MarkFlagsRequired(cmd, networkIdFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command) (*inputModel, error) {
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

	if p.IsVerbosityDebug() {
		modelStr, err := print.BuildDebugStrFromInputModel(model)
		if err != nil {
			p.Debug(print.ErrorLevel, "convert model to string for debugging: %v", err)
		} else {
			p.Debug(print.DebugLevel, "parsed input values: %s", modelStr)
		}
	}

	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *iaas.APIClient) iaas.ApiListNicsRequest {
	req := apiClient.ListNics(ctx, model.ProjectId, *model.NetworkId)
	if model.LabelSelector != nil {
		req = req.LabelSelector(*model.LabelSelector)
	}

	return req
}

func outputResult(p *print.Printer, outputFormat string, nics []iaas.NIC) error {
	switch outputFormat {
	case print.JSONOutputFormat:
		details, err := json.MarshalIndent(nics, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal nics: %w", err)
		}
		p.Outputln(string(details))

		return nil
	case print.YAMLOutputFormat:
		details, err := yaml.MarshalWithOptions(nics, yaml.IndentSequence(true))
		if err != nil {
			return fmt.Errorf("marshal nics: %w", err)
		}
		p.Outputln(string(details))

		return nil
	default:
		table := tables.NewTable()
		table.SetHeader("ID", "NAME", "NIC SECURITY", "STATUS", "TYPE")

		for _, nic := range nics {
			name := ""
			if nic.Name != nil {
				name = *nic.Name
			}
			table.AddRow(*nic.Id, name, *nic.NicSecurity, *nic.Status, *nic.Type)
			table.AddSeparator()
		}

		p.Outputln(table.Render())
		return nil
	}
}
