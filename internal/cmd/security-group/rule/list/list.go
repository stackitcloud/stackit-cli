package list

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/types"

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
	SecurityGroupId string
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Lists all security group rules in a security group of a project",
		Long:  "Lists all security group rules in a security group of a project.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Lists all security group rules in security group with ID "xxx"`,
				"$ stackit security-group rule list --security-group-id xxx",
			),
			examples.NewExample(
				`Lists all security group rules in security group with ID "xxx" in JSON format`,
				"$ stackit security-group rule list --security-group-id xxx --output-format json",
			),
			examples.NewExample(
				`Lists up to 10 security group rules in security group with ID "xxx"`,
				"$ stackit security-group rule list --security-group-id xxx --limit 10",
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
				return fmt.Errorf("list security group rules: %w", err)
			}

			if resp.Items == nil || len(*resp.Items) == 0 {
				securityGroupLabel, err := iaasUtils.GetSecurityGroupName(ctx, apiClient, model.ProjectId, model.Region, model.SecurityGroupId)
				if err != nil {
					params.Printer.Debug(print.ErrorLevel, "get security group name: %v", err)
					securityGroupLabel = model.SecurityGroupId
				}

				projectLabel, err := projectname.GetProjectName(ctx, params.Printer, params.CliVersion, cmd)
				if err != nil {
					params.Printer.Debug(print.ErrorLevel, "get project name: %v", err)
					projectLabel = model.ProjectId
				}
				params.Printer.Info("No rules found in security group %q for project %q\n", securityGroupLabel, projectLabel)
				return nil
			}

			// Truncate output
			items := *resp.Items
			if model.Limit != nil && len(items) > int(*model.Limit) {
				items = items[:*model.Limit]
			}

			return outputResult(params.Printer, model.OutputFormat, items)
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
		SecurityGroupId: flags.FlagToStringValue(p, cmd, securityGroupIdFlag),
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *iaas.APIClient) iaas.ApiListSecurityGroupRulesRequest {
	return apiClient.ListSecurityGroupRules(ctx, model.ProjectId, model.Region, model.SecurityGroupId)
}

func outputResult(p *print.Printer, outputFormat string, securityGroupRules []iaas.SecurityGroupRule) error {
	return p.OutputResult(outputFormat, securityGroupRules, func() error {
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
	})
}
