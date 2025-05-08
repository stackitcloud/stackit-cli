package delete

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/iaas/client"
	iaasUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/iaas/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/iaas"
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
		Use:   fmt.Sprintf("delete %s", securityGroupRuleIdArg),
		Short: "Deletes a security group rule",
		Long: fmt.Sprintf("%s\n%s\n",
			"Deletes a security group rule.",
			"If the security group rule is still in use, the deletion will fail",
		),
		Args: args.SingleArg(securityGroupRuleIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Delete security group rule with ID "xxx" in security group with ID "yyy"`,
				"$ stackit security-group rule delete xxx --security-group-id yyy",
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

			securityGroupLabel, err := iaasUtils.GetSecurityGroupName(ctx, apiClient, model.ProjectId, *model.SecurityGroupId)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get security group name: %v", err)
				securityGroupLabel = *model.SecurityGroupId
			}

			securityGroupRuleLabel, err := iaasUtils.GetSecurityGroupRuleName(ctx, apiClient, model.ProjectId, model.SecurityGroupRuleId, *model.SecurityGroupId)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get security group rule name: %v", err)
				securityGroupRuleLabel = model.SecurityGroupRuleId
			}

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to delete security group rule %q from security group %q?", securityGroupRuleLabel, securityGroupLabel)
				err = params.Printer.PromptForConfirmation(prompt)
				if err != nil {
					return err
				}
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			err = req.Execute()
			if err != nil {
				return fmt.Errorf("delete security group rule: %w", err)
			}

			params.Printer.Info("Deleted security group rule %q from security group %q\n", securityGroupRuleLabel, securityGroupLabel)
			return nil
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

func buildRequest(ctx context.Context, model *inputModel, apiClient *iaas.APIClient) iaas.ApiDeleteSecurityGroupRuleRequest {
	return apiClient.DeleteSecurityGroupRule(ctx, model.ProjectId, *model.SecurityGroupId, model.SecurityGroupRuleId)
}
