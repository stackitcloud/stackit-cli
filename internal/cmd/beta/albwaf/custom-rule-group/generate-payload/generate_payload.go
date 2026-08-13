package generatepayload

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/types"

	albwaf "github.com/stackitcloud/stackit-sdk-go/services/albwaf/v1api"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/fileutils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/albwaf/client"
	albwafUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/albwaf/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

const (
	nameFlag     = "name"
	filePathFlag = "file-path"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	Name     *string
	FilePath *string
}

var defaultCreateCustomRuleGroupPayload = albwaf.CreateCustomRuleGroupPayload{
	Name: "example-custom-rule-group",
	Rules: []albwaf.CreateCustomRule{
		{
			Behavior: albwaf.Behavior{
				Action: albwaf.ACTION_ACTION_DENY,
				Log:    new(true),
				LogMsg: new(""),
			},
			Conditions: []albwaf.Condition{
				{
					Operator: albwaf.ConditionOperator{
						Type:  albwaf.OPERATOR_OPERATOR_VALIDATE_UTF8_ENCODING,
						Value: new(""),
					},
					Variable: albwaf.ConditionVariable{
						Type:  albwaf.VARIABLE_VARIABLE_RESPONSE_STATUS,
						Value: new(""),
					},
				},
			},
			Description: new(""),
		},
	},
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "generate-payload",
		Short: "Generates a payload to create/update an ALB WAF custom rule group",
		Long: fmt.Sprintf("%s\n%s\n%s",
			"Generates a JSON payload with values to be used as --payload input for ALB WAF custom rule group creation or update.",
			"If --name is set, the payload is generated with the current values of the given custom rule group, to be used with the update command. If unset, a payload with default values for the create command is generated.",
			"See https://docs.api.stackit.cloud/documentation/alb-waf/version/v1 for information regarding the payload structure.",
		),
		Args: args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Generate a payload with default values, and adapt it with custom values for the different configuration options`,
				`$ stackit beta alb-waf custom-rule-group generate-payload --file-path ./payload.json`,
				`<Modify payload in file, if needed>`,
				`$ stackit beta alb-waf custom-rule-group create --payload @./payload.json`),
			examples.NewExample(
				`Generate a payload with the current values of an existing custom rule group, and adapt it with custom values for the different configuration options`,
				`$ stackit beta alb-waf custom-rule-group generate-payload --name my-custom-rule-group --file-path ./payload.json`,
				`<Modify payload in file>`,
				`$ stackit beta alb-waf custom-rule-group update my-custom-rule-group --payload @./payload.json`),
			examples.NewExample(
				`Generate a payload with the current values of an existing custom rule group, and preview it in the terminal`,
				`$ stackit beta alb-waf custom-rule-group generate-payload --name my-custom-rule-group`),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			model, err := parseInput(params.Printer, cmd, args)
			if err != nil {
				return err
			}

			if model.Name == nil {
				payload := defaultCreateCustomRuleGroupPayload
				return outputResult(params.Printer, model.FilePath, &payload)
			}

			// Configure API client
			apiClient, err := client.ConfigureClient(params.Printer, params.CliVersion)
			if err != nil {
				return err
			}

			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("read ALB WAF custom rule group: %w", err)
			}

			payload := &albwaf.UpdateCustomRuleGroupPayload{
				Name:  resp.Name,
				Rules: albwafUtils.ToCreateCustomRules(resp.Rules),
			}
			return outputResult(params.Printer, model.FilePath, payload)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().StringP(nameFlag, "n", "", "If set, generates the payload with the current values of the given custom rule group (to be used with the update command). If unset, generates the payload with default values (to be used with the create command)")
	cmd.Flags().StringP(filePathFlag, "f", "", "If set, writes the payload to the given file. If unset, writes the payload to the standard output")
}

func parseInput(p *print.Printer, cmd *cobra.Command, _ []string) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)

	name := flags.FlagToStringPointer(p, cmd, nameFlag)
	// If name is provided, the custom rule group is fetched from the API, so a project ID is needed as well
	if name != nil && globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		Name:            name,
		FilePath:        flags.FlagToStringPointer(p, cmd, filePathFlag),
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *albwaf.APIClient) albwaf.ApiGetCustomRuleGroupRequest {
	return apiClient.DefaultAPI.GetCustomRuleGroup(ctx, model.ProjectId, model.Region, *model.Name)
}

func outputResult(p *print.Printer, filePath *string, payload any) error {
	if payload == nil {
		return fmt.Errorf("payload is empty")
	}
	payloadBytes, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal payload: %w", err)
	}

	if filePath != nil {
		err = fileutils.WriteToFile(utils.PtrString(filePath), string(payloadBytes))
		if err != nil {
			return fmt.Errorf("write payload to the file: %w", err)
		}
	} else {
		p.Outputln(string(payloadBytes))
	}

	return nil
}
