package create

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
	payloadFlag = "payload"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	Payload *albwaf.CreateCustomRuleGroupPayload
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Creates an ALB WAF custom rule group",
		Long: fmt.Sprintf("%s\n%s\n%s",
			"Creates a STACKIT Application Load Balancer (ALB) Web Application Firewall (WAF) custom rule group.",
			"The payload can be provided as a JSON string or a file path prefixed with \"@\".",
			"See https://docs.api.stackit.cloud/documentation/alb-waf/version/v1 for information regarding the payload structure.",
		),
		Args: args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Create an ALB WAF custom rule group using an API payload sourced from the file "./payload.json"`,
				"$ stackit beta alb-waf custom-rule-group create --payload @./payload.json"),
			examples.NewExample(
				`Create an ALB WAF custom rule group using an API payload provided as a JSON string`,
				`$ stackit beta alb-waf custom-rule-group create --payload "{...}"`),
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

			prompt := fmt.Sprintf("Are you sure you want to create an ALB WAF custom rule group for project %q?", projectLabel)
			err = params.Printer.PromptForConfirmation(prompt)
			if err != nil {
				return err
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("create ALB WAF custom rule group: %w", err)
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

func parseInput(p *print.Printer, cmd *cobra.Command, _ []string) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	payloadValue := flags.FlagToStringPointer(p, cmd, payloadFlag)
	var payload *albwaf.CreateCustomRuleGroupPayload
	if payloadValue != nil {
		payload = &albwaf.CreateCustomRuleGroupPayload{}
		err := json.Unmarshal([]byte(*payloadValue), payload)
		if err != nil {
			return nil, fmt.Errorf("decode payload: %w", err)
		}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		Payload:         payload,
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *albwaf.APIClient) albwaf.ApiCreateCustomRuleGroupRequest {
	req := apiClient.DefaultAPI.CreateCustomRuleGroup(ctx, model.ProjectId, model.Region)
	req = req.CreateCustomRuleGroupPayload(*model.Payload)
	return req
}

func outputResult(p *print.Printer, outputFormat string, resp *albwaf.GetCustomRuleGroupResponse) error {
	return p.OutputResult(outputFormat, resp, func() error {
		if resp == nil {
			p.Outputf("Received empty custom rule group response")
			return nil
		}
		p.Outputf("Created custom rule group with name %q \n", resp.Name)
		return nil
	})
}
