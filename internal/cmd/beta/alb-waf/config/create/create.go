package create

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
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/albwaf/client"

	"github.com/spf13/cobra"
	albwaf "github.com/stackitcloud/stackit-sdk-go/services/albwaf/v1api"
)

const (
	nameFlag                = "name"
	managedRuleSetNameFlag  = "managed-rule-set-name"
	customRuleGroupNameFlag = "custom-rule-group-name"
	labelsFlag              = "labels"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	Name                string
	ManagedRuleSetName  *string
	CustomRuleGroupName *string
	Labels              *map[string]string
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Creates an ALB WAF configuration",
		Long:  "Creates a STACKIT Application Load Balancer (ALB) Web Application Firewall (WAF) configuration.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Create an ALB WAF configuration with name "my-waf-config"`,
				"$ stackit beta alb-waf config create --name my-waf-config"),
			examples.NewExample(
				`Create an ALB WAF configuration with a managed rule set, a custom rule group and labels`,
				"$ stackit beta alb-waf config create --name my-waf-config --managed-rule-set-name my-managed-rule-set --custom-rule-group-name my-custom-rule-group --labels key1=value1,key2=value2"),
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

			projectLabel, err := projectname.GetProjectName(ctx, params.Printer, params.CliVersion, cmd)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get project name: %v", err)
				projectLabel = model.ProjectId
			}

			prompt := fmt.Sprintf("Are you sure you want to create an ALB WAF configuration for project %q?", projectLabel)
			err = params.Printer.PromptForConfirmation(prompt)
			if err != nil {
				return err
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("create ALB WAF configuration: %w", err)
			}

			return outputResult(params.Printer, model.OutputFormat, projectLabel, resp)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().String(nameFlag, "", "Name of the WAF configuration")
	cmd.Flags().String(managedRuleSetNameFlag, "", "Name of the managed rule set configuration to attach to the WAF")
	cmd.Flags().String(customRuleGroupNameFlag, "", "Name of the custom rule group configuration to attach to the WAF")
	cmd.Flags().StringToString(labelsFlag, nil, "Labels are key-value string pairs which can be attached to the WAF configuration. E.g. '--labels key1=value1,key2=value2,...'")

	err := flags.MarkFlagsRequired(cmd, nameFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command, _ []string) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel:     globalFlags,
		Name:                flags.FlagToStringValue(p, cmd, nameFlag),
		ManagedRuleSetName:  flags.FlagToStringPointer(p, cmd, managedRuleSetNameFlag),
		CustomRuleGroupName: flags.FlagToStringPointer(p, cmd, customRuleGroupNameFlag),
		Labels:              flags.FlagToStringToStringPointer(p, cmd, labelsFlag),
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *albwaf.APIClient) albwaf.ApiCreateWAFRequest {
	req := apiClient.DefaultAPI.CreateWAF(ctx, model.ProjectId, model.Region)
	payload := albwaf.CreateWAFPayload{
		Name:                model.Name,
		ManagedRuleSetName:  model.ManagedRuleSetName,
		CustomRuleGroupName: model.CustomRuleGroupName,
		Labels:              model.Labels,
	}
	return req.CreateWAFPayload(payload)
}

func outputResult(p *print.Printer, outputFormat, projectLabel string, resp *albwaf.GetWAFResponse) error {
	return p.OutputResult(outputFormat, resp, func() error {
		if resp == nil {
			return fmt.Errorf("create WAF configuration response is empty")
		}
		p.Outputf("Created WAF configuration %q for project %q.\n", resp.Name, projectLabel)
		return nil
	})
}
