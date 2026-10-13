package describe

import (
	"context"
	"fmt"
	"strings"

	"github.com/stackitcloud/stackit-cli/internal/pkg/types"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/albwaf/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
	albwaf "github.com/stackitcloud/stackit-sdk-go/services/albwaf/v1api"
)

const (
	nameArg = "NAME"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	Name string
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("describe %s", nameArg),
		Short: "Shows details of an ALB WAF configuration",
		Long:  "Shows details of a STACKIT Application Load Balancer (ALB) Web Application Firewall (WAF) configuration.",
		Args:  args.SingleArg(nameArg, nil),
		Example: examples.Build(
			examples.NewExample(
				`Show details of an ALB WAF configuration with name "my-waf-config"`,
				`$ stackit beta alb-waf config describe my-waf-config`,
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
				return fmt.Errorf("read ALB WAF configuration: %w", err)
			}

			return outputResult(params.Printer, model.OutputFormat, resp)
		},
	}
	return cmd
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	name := inputArgs[0]
	model := inputModel{
		GlobalFlagModel: globalFlags,
		Name:            name,
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *albwaf.APIClient) albwaf.ApiGetWAFRequest {
	return apiClient.DefaultAPI.GetWAF(ctx, model.ProjectId, model.Region, model.Name)
}

func outputResult(p *print.Printer, outputFormat string, resp *albwaf.GetWAFResponse) error {
	return p.OutputResult(outputFormat, resp, func() error {
		if resp == nil {
			return fmt.Errorf("no WAF configuration found")
		}
		table := tables.NewTable()
		table.SetTitle("WAF Configuration")
		table.AddRow("NAME", resp.Name)
		table.AddSeparator()

		managedRuleSet := "-"
		if resp.ManagedRuleSetName != nil && *resp.ManagedRuleSetName != "" {
			managedRuleSet = *resp.ManagedRuleSetName
		}
		table.AddRow("MANAGED RULE SET", managedRuleSet)
		table.AddSeparator()

		customRuleGroup := "-"
		if resp.CustomRuleGroupName != nil && *resp.CustomRuleGroupName != "" {
			customRuleGroup = *resp.CustomRuleGroupName
		}
		table.AddRow("CUSTOM RULE GROUP", customRuleGroup)
		table.AddSeparator()

		if resp.Labels != nil && len(*resp.Labels) > 0 {
			table.AddRow("LABELS", formatLabels(*resp.Labels))
			table.AddSeparator()
		}

		usageCount := ""
		if resp.Usage != nil {
			usageCount = utils.PtrString(resp.Usage.Count)
		}
		table.AddRow("USAGE COUNT", usageCount)
		table.AddSeparator()

		if resp.Usage != nil && len(resp.Usage.Items) > 0 {
			table.AddRow("USED BY", formatUsageItems(resp.Usage.Items))
			table.AddSeparator()
		}

		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}
		return nil
	})
}

func formatLabels(labels map[string]string) string {
	pairs := make([]string, 0, len(labels))
	for k, v := range labels {
		pairs = append(pairs, fmt.Sprintf("%s=%s", k, v))
	}
	return strings.Join(pairs, "\n")
}

// formatUsageItems renders each WAF usage item as "<load-balancer-name> (listeners: l1, l2, ...)",
// or just "<load-balancer-name>" if no listener names are set, one per line.
func formatUsageItems(items []albwaf.WAFUsageItem) string {
	lines := make([]string, 0, len(items))
	for i := range items {
		item := &items[i]
		if len(item.ListenerNames) > 0 {
			lines = append(lines, fmt.Sprintf("%s (listeners: %s)", item.LoadBalancerName, strings.Join(item.ListenerNames, ", ")))
		} else {
			lines = append(lines, item.LoadBalancerName)
		}
	}
	return strings.Join(lines, "\n")
}
