package list

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/projectname"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/iaas/client"
	iaasUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/iaas/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/goccy/go-yaml"
	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/iaas"
)

const (
	limitFlag = "limit"

	securityGroupIdFlag = "security-group-id"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	Limit           *int64
	SecurityGroupId *string
}

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Lists all security group rules in a security group of a project",
		Long:  "Lists all security group rules in a security group of a project.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Lists all security group rules in security group with ID "xxx"`,
				"$ stackit beta security-group rule list --security-group-id xxx",
			),
			examples.NewExample(
				`Lists all security group rules in security group with ID "xxx" in JSON format`,
				"$ stackit beta security-group rule list --security-group-id xxx --output-format json",
			),
			examples.NewExample(
				`Lists up to 10 security group rules in security group with ID "xxx"`,
				"$ stackit beta security-group rule list --security-group-id xxx --limit 10",
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
				return fmt.Errorf("list security group rules: %w", err)
			}

			if resp.Items == nil || len(*resp.Items) == 0 {
				securityGroupLabel, err := iaasUtils.GetSecurityGroupName(ctx, apiClient, model.ProjectId, *model.SecurityGroupId)
				if err != nil {
					p.Debug(print.ErrorLevel, "get security group name: %v", err)
					securityGroupLabel = *model.SecurityGroupId
				}

				projectLabel, err := projectname.GetProjectName(ctx, p, cmd)
				if err != nil {
					p.Debug(print.ErrorLevel, "get project name: %v", err)
					projectLabel = model.ProjectId
				}
				p.Info("No rules found in security group %q for project %q\n", securityGroupLabel, projectLabel)
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
	cmd.Flags().Int64(limitFlag, 0, `Maximum number of entries to list`)
	cmd.Flags().Var(flags.UUIDFlag(), securityGroupIdFlag, `The security group ID`)

	err := flags.MarkFlagsRequired(cmd, securityGroupIdFlag)
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
		SecurityGroupId: flags.FlagToStringPointer(p, cmd, securityGroupIdFlag),
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

func buildRequest(ctx context.Context, model *inputModel, apiClient *iaas.APIClient) iaas.ApiListSecurityGroupRulesRequest {
	return apiClient.ListSecurityGroupRules(ctx, model.ProjectId, *model.SecurityGroupId)
}

func outputResult(p *print.Printer, outputFormat string, securityGroupRules []iaas.SecurityGroupRule) error {
	switch outputFormat {
	case print.JSONOutputFormat:
		details, err := json.MarshalIndent(securityGroupRules, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal security group rules: %w", err)
		}
		p.Outputln(string(details))

		return nil
	case print.YAMLOutputFormat:
		details, err := yaml.MarshalWithOptions(securityGroupRules, yaml.IndentSequence(true), yaml.UseJSONMarshaler())
		if err != nil {
			return fmt.Errorf("marshal security group rules: %w", err)
		}
		p.Outputln(string(details))

		return nil
	default:
		table := tables.NewTable()
		table.SetHeader("ID", "ETHER TYPE", "DIRECTION", "PROTOCOL", "REMOTE SECURITY GROUP ID")

		for _, securityGroupRule := range securityGroupRules {
			etherType := utils.PtrStringDefault(securityGroupRule.Ethertype, "")

			protocolName := ""
			if securityGroupRule.Protocol != nil {
				if securityGroupRule.Protocol.Name != nil {
					protocolName = *securityGroupRule.Protocol.Name
				}
			}

			table.AddRow(
				utils.PtrString(securityGroupRule.Id),
				etherType,
				utils.PtrString(securityGroupRule.Direction),
				protocolName,
				utils.PtrString(securityGroupRule.RemoteSecurityGroupId),
			)
			table.AddSeparator()
		}

		p.Outputln(table.Render())
		return nil
	}
}
