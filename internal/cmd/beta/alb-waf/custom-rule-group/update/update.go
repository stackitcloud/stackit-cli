package update

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/types"

	"github.com/spf13/cobra"
	albwaf "github.com/stackitcloud/stackit-sdk-go/services/albwaf/v1api"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/projectname"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/albwaf/client"
)

const (
	customRuleGroupNameArg = "CUSTOM_RULE_GROUP_NAME"

	payloadFlag = "payload"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	Name    string
	Payload *albwaf.UpdateCustomRuleGroupPayload
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("update %s", customRuleGroupNameArg),
		Short: "Updates an ALB WAF custom rule group",
		Long: fmt.Sprintf("%s\n%s\n%s\n%s",
			"Updates a STACKIT Application Load Balancer (ALB) Web Application Firewall (WAF) custom rule group.",
			"The rules of the custom rule group are replaced atomically: send the complete desired rule set, per-rule partial merge is not supported.",
			"The payload can be provided as a JSON string or a file path prefixed with \"@\". The rule IDs and the rule behavior severity are managed by the server and cannot be set.",
			"See https://docs.api.stackit.cloud/documentation/alb-waf/version/v1 for information regarding the payload structure.",
		),
		Args: args.SingleArg(customRuleGroupNameArg, nil),
		Example: examples.Build(
			examples.NewExample(
				`Update an ALB WAF custom rule group using an API payload sourced from the file "./payload.json"`,
				"$ stackit beta alb-waf custom-rule-group update my-custom-rule-group --payload @./payload.json"),
			examples.NewExample(
				`Update an ALB WAF custom rule group using an API payload provided as a JSON string`,
				`$ stackit beta alb-waf custom-rule-group update my-custom-rule-group --payload "{...}"`),
			examples.NewExample(
				`Generate a payload with the current values of an existing custom rule group, adapt it and update the custom rule group with it`,
				`$ stackit beta alb-waf custom-rule-group generate-payload --name my-custom-rule-group > ./payload.json`,
				`<Modify payload in file>`,
				`$ stackit beta alb-waf custom-rule-group update my-custom-rule-group --payload @./payload.json`),
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

			prompt := fmt.Sprintf("Are you sure you want to update the ALB WAF custom rule group %q for project %q?", model.Name, projectLabel)
			err = params.Printer.PromptForConfirmation(prompt)
			if err != nil {
				return err
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("update ALB WAF custom rule group: %w", err)
			}

			return outputResult(params.Printer, model.OutputFormat, resp)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.ReadFromFileFlag(), payloadFlag, `Request payload (JSON). Can be a string or a file path, if prefixed with "@" (example: @./payload.json).`)

	err := flags.MarkFlagsRequired(cmd, payloadFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	name := inputArgs[0]

	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	payloadValue := flags.FlagToStringPointer(p, cmd, payloadFlag)
	var payload *albwaf.UpdateCustomRuleGroupPayload
	if payloadValue != nil {
		payload = &albwaf.UpdateCustomRuleGroupPayload{}
		err := json.Unmarshal([]byte(*payloadValue), payload)
		if err != nil {
			return nil, fmt.Errorf("encode payload: %w", err)
		}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		Name:            name,
		Payload:         payload,
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *albwaf.APIClient) albwaf.ApiUpdateCustomRuleGroupRequest {
	req := apiClient.DefaultAPI.UpdateCustomRuleGroup(ctx, model.ProjectId, model.Region, model.Name)
	req = req.UpdateCustomRuleGroupPayload(*model.Payload)
	return req
}

func outputResult(p *print.Printer, outputFormat string, resp *albwaf.GetCustomRuleGroupResponse) error {
	return p.OutputResult(outputFormat, resp, func() error {
		if resp == nil {
			p.Outputf("Received empty custom rule group response")
			return nil
		}
		p.Outputf("Updated custom rule group %q.\n", resp.Name)
		return nil
	})
}
