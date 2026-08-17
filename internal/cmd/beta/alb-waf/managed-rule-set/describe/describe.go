package describe

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
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
		Short: "Describes a managed rule set of the ALB WAF",
		Long:  "Describes a managed rule set (MRS) of the Web Application Firewall (WAF) for application loadbalancers.",
		Args:  args.SingleArg(nameArg, nil),
		Example: examples.Build(
			examples.NewExample(
				`Get details about a managed rule set with name "my-managed-rule-set"`,
				"$ stackit beta alb-waf managed-rule-set describe my-managed-rule-set",
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
				return fmt.Errorf("read managed rule set: %w", err)
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

func buildRequest(ctx context.Context, model *inputModel, apiClient *albwaf.APIClient) albwaf.ApiGetManagedRuleSetRequest {
	return apiClient.DefaultAPI.GetManagedRuleSet(ctx, model.ProjectId, model.Region, model.Name)
}

func outputResult(p *print.Printer, outputFormat string, resp *albwaf.GetManagedRuleSetResponse) error {
	if resp == nil {
		return fmt.Errorf("no managed rule set found")
	}
	return p.OutputResult(outputFormat, resp, func() error {
		usageCount := ""
		if resp.Usage != nil {
			usageCount = utils.PtrString(resp.Usage.Count)
		}

		content := []tables.Table{}

		table := tables.NewTable()
		table.SetTitle("Managed Rule Set")
		table.AddRow("NAME", resp.Name)
		table.AddSeparator()
		table.AddRow("TYPE", string(resp.Type))
		table.AddSeparator()
		table.AddRow("VERSION", resp.Version)
		table.AddSeparator()
		table.AddRow("USAGE COUNT", usageCount)
		content = append(content, table)

		if resp.Groups != nil && len(*resp.Groups) > 0 {
			for _, groupId := range utils.SortedKeys(*resp.Groups) {
				group := (*resp.Groups)[groupId]
				if group.Rules == nil || len(*group.Rules) == 0 {
					continue
				}
				groupTable := tables.NewTable()
				groupTable.SetTitle(fmt.Sprintf("Rule Group: %s (%s)", groupId, group.GroupName))
				groupTable.SetHeader("DESCRIPTION", "RULE ID", "SEVERITY", "MODE")
				for _, ruleId := range utils.SortedKeys(*group.Rules) {
					rule := (*group.Rules)[ruleId]
					groupTable.AddRow(utils.Truncate(new(rule.Description), 70), ruleId, rule.Severity, string(rule.Mode))
					groupTable.AddSeparator()
				}
				content = append(content, groupTable)
			}
		}

		err := tables.DisplayTables(p, content)
		if err != nil {
			return fmt.Errorf("render tables: %w", err)
		}
		if resp.Groups != nil && len(*resp.Groups) > 0 {
			p.Outputln("\nUse --output-format json/yaml to see the untruncated rule descriptions and severity.")
		}
		return nil
	})
}
