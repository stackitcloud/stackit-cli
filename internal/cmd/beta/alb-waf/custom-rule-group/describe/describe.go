package describe

import (
	"context"
	"fmt"
	"strings"

	"github.com/stackitcloud/stackit-cli/internal/pkg/types"

	"github.com/spf13/cobra"
	sdkUtils "github.com/stackitcloud/stackit-sdk-go/core/utils"
	albwaf "github.com/stackitcloud/stackit-sdk-go/services/albwaf/v1api"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/albwaf/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

const (
	customRuleGroupNameArg = "CUSTOM_RULE_GROUP_NAME"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	Name string
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("describe %s", customRuleGroupNameArg),
		Short: "Shows details of an ALB WAF custom rule group",
		Long:  "Shows details of a STACKIT Application Load Balancer (ALB) Web Application Firewall (WAF) custom rule group.",
		Args:  args.SingleArg(customRuleGroupNameArg, nil),
		Example: examples.Build(
			examples.NewExample(
				`Show details of an ALB WAF custom rule group with name "my-custom-rule-group"`,
				`$ stackit beta alb-waf custom-rule-group describe my-custom-rule-group`,
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
				return fmt.Errorf("read ALB WAF custom rule group: %w", err)
			}

			return outputResult(params.Printer, model.OutputFormat, resp)
		},
	}
	return cmd
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	name := inputArgs[0]

	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		Name:            name,
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *albwaf.APIClient) albwaf.ApiGetCustomRuleGroupRequest {
	return apiClient.DefaultAPI.GetCustomRuleGroup(ctx, model.ProjectId, model.Region, model.Name)
}

func outputResult(p *print.Printer, outputFormat string, crg *albwaf.GetCustomRuleGroupResponse) error {
	return p.OutputResult(outputFormat, crg, func() error {
		if crg == nil {
			return fmt.Errorf("custom rule group response is empty")
		}
		content := []tables.Table{buildOverviewTable(crg)}

		for i := range crg.Rules {
			content = append(content, buildRuleTable(&crg.Rules[i]))
		}

		err := tables.DisplayTables(p, content)
		if err != nil {
			return fmt.Errorf("display table: %w", err)
		}
		return nil
	})
}

func buildOverviewTable(crg *albwaf.GetCustomRuleGroupResponse) tables.Table {
	table := tables.NewTable()
	table.SetTitle("Custom Rule Group")
	table.AddRow("NAME", crg.Name)
	table.AddSeparator()
	table.AddRow("RULES", len(crg.Rules))
	table.AddSeparator()
	if crg.Usage != nil {
		table.AddRow("USED BY WAF CONFIGS", strings.Join(crg.Usage.Items, "\n"))
		table.AddSeparator()
	}
	return table
}

func buildRuleTable(rule *albwaf.GetCustomRule) tables.Table {
	table := tables.NewTable()
	table.SetTitle(fmt.Sprintf("Rule %d", rule.Id))
	table.AddRow("ID", rule.Id)
	table.AddSeparator()
	table.AddRow("DESCRIPTION", utils.PtrString(rule.Description))
	table.AddSeparator()
	table.AddRow("ACTION", string(rule.Behavior.Action))
	table.AddSeparator()
	table.AddRow("SEVERITY", string(rule.Behavior.Severity))
	table.AddSeparator()
	table.AddRow("LOG", rule.Behavior.Log)
	table.AddSeparator()
	table.AddRow("LOG MESSAGE", utils.PtrString(rule.Behavior.LogMsg))
	table.AddSeparator()
	for i := range rule.Conditions {
		condition := &rule.Conditions[i]
		table.AddRow("CONDITION VARIABLE", conditionVariableString(condition.Variable))
		table.AddSeparator()
		table.AddRow("CONDITION OPERATOR", conditionOperatorString(condition.Operator))
		table.AddSeparator()
		if len(condition.Transformations) > 0 {
			table.AddRow("CONDITION TRANSFORMATIONS", strings.Join(sdkUtils.EnumSliceToStringSlice(condition.Transformations), ", "))
			table.AddSeparator()
		}
	}
	return table
}

// conditionVariableString renders a condition variable as "TYPE" or "TYPE:value".
func conditionVariableString(variable albwaf.ConditionVariable) string {
	if variable.Value != nil && *variable.Value != "" {
		return fmt.Sprintf("%s:%s", string(variable.Type), *variable.Value)
	}
	return string(variable.Type)
}

// conditionOperatorString renders a condition operator as "TYPE value".
func conditionOperatorString(operator albwaf.ConditionOperator) string {
	if operator.Value != nil && *operator.Value != "" {
		return fmt.Sprintf("%s %q", string(operator.Type), *operator.Value)
	}
	return string(operator.Type)
}
