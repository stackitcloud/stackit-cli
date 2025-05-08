package describe

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/goccy/go-yaml"
	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/iaas/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/iaas"

	"github.com/spf13/cobra"
)

const (
	securityGroupRuleIdArg = "SECURITY_GROUP_RULE_ID"

	securityGroupIdFlag = "security-group-id"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	SecurityGroupRuleId string
	SecurityGroupId     *string
}

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("describe %s", securityGroupRuleIdArg),
		Short: "Shows details of a security group rule",
		Long:  "Shows details of a security group rule.",
		Args:  args.SingleArg(securityGroupRuleIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Show details of a security group rule with ID "xxx" in security group with ID "yyy"`,
				"$ stackit security-group rule describe xxx --security-group-id yyy",
			),
			examples.NewExample(
				`Show details of a security group rule with ID "xxx" in security group with ID "yyy" in JSON format`,
				"$ stackit security-group rule describe xxx --security-group-id yyy --output-format json",
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			model, err := parseInput(params.Printer, cmd, args)
			if err != nil {
				return err
			}

			// Configure API client
			apiClient, err := client.ConfigureClient(params.Printer)
			if err != nil {
				return err
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("read security group rule: %w", err)
			}

			return outputResult(params.Printer, model.OutputFormat, resp)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), securityGroupIdFlag, `The security group ID`)

	err := flags.MarkFlagsRequired(cmd, securityGroupIdFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	securityGroupRuleId := inputArgs[0]

	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel:     globalFlags,
		SecurityGroupRuleId: securityGroupRuleId,
		SecurityGroupId:     flags.FlagToStringPointer(p, cmd, securityGroupIdFlag),
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

func buildRequest(ctx context.Context, model *inputModel, apiClient *iaas.APIClient) iaas.ApiGetSecurityGroupRuleRequest {
	return apiClient.GetSecurityGroupRule(ctx, model.ProjectId, *model.SecurityGroupId, model.SecurityGroupRuleId)
}

func outputResult(p *print.Printer, outputFormat string, securityGroupRule *iaas.SecurityGroupRule) error {
	if securityGroupRule == nil {
		return fmt.Errorf("security group rule is empty")
	}
	switch outputFormat {
	case print.JSONOutputFormat:
		details, err := json.MarshalIndent(securityGroupRule, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal security group rule: %w", err)
		}
		p.Outputln(string(details))

		return nil
	case print.YAMLOutputFormat:
		details, err := yaml.MarshalWithOptions(securityGroupRule, yaml.IndentSequence(true), yaml.UseJSONMarshaler())
		if err != nil {
			return fmt.Errorf("marshal security group rule: %w", err)
		}
		p.Outputln(string(details))

		return nil
	default:
		table := tables.NewTable()
		table.AddRow("ID", utils.PtrString(securityGroupRule.Id))
		table.AddSeparator()

		if securityGroupRule.Protocol != nil {
			if securityGroupRule.Protocol.Name != nil {
				table.AddRow("PROTOCOL NAME", *securityGroupRule.Protocol.Name)
				table.AddSeparator()
			}

			if securityGroupRule.Protocol.Number != nil {
				table.AddRow("PROTOCOL NUMBER", *securityGroupRule.Protocol.Number)
				table.AddSeparator()
			}
		}

		table.AddRow("DIRECTION", utils.PtrString(securityGroupRule.Direction))
		table.AddSeparator()

		if securityGroupRule.PortRange != nil {
			if securityGroupRule.PortRange.Min != nil {
				table.AddRow("START PORT", *securityGroupRule.PortRange.Min)
				table.AddSeparator()
			}

			if securityGroupRule.PortRange.Max != nil {
				table.AddRow("END PORT", *securityGroupRule.PortRange.Max)
				table.AddSeparator()
			}
		}

		if securityGroupRule.Ethertype != nil {
			table.AddRow("ETHER TYPE", *securityGroupRule.Ethertype)
			table.AddSeparator()
		}

		if securityGroupRule.IpRange != nil {
			table.AddRow("IP RANGE", *securityGroupRule.IpRange)
			table.AddSeparator()
		}

		if securityGroupRule.RemoteSecurityGroupId != nil {
			table.AddRow("REMOTE SECURITY GROUP", *securityGroupRule.RemoteSecurityGroupId)
			table.AddSeparator()
		}

		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}
		return nil
	}
}
