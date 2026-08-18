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
	nameFlag = "name"
)

var typeFlag = flags.StringEnumFlag(
	"type",
	albwaf.AllowedTypeEnumValues,
	"Managed rule set type",
	flags.StringEnumDefaultValue(albwaf.TYPE_TYPE_OWASP_CRS),
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	Name string
	Type albwaf.Type
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Creates a managed rule set for the ALB WAF",
		Long:  "Creates a managed rule set (MRS) for the Web Application Firewall (WAF) of application loadbalancers.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Create a managed rule set with name "my-managed-rule-set"`,
				"$ stackit beta alb-waf managed-rule-set create --name my-managed-rule-set"),
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

			prompt := fmt.Sprintf("Are you sure you want to create a managed rule set for project %q?", projectLabel)
			err = params.Printer.PromptForConfirmation(prompt)
			if err != nil {
				return err
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("create managed rule set: %w", err)
			}

			return outputResult(params.Printer, model.OutputFormat, projectLabel, resp)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().String(nameFlag, "", "Name of the managed rule set")
	typeFlag.Register(cmd.Flags())

	err := flags.MarkFlagsRequired(cmd, nameFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command, _ []string) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		Name:            flags.FlagToStringValue(p, cmd, nameFlag),
		Type:            typeFlag.Get(),
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *albwaf.APIClient) albwaf.ApiCreateManagedRuleSetRequest {
	req := apiClient.DefaultAPI.CreateManagedRuleSet(ctx, model.ProjectId, model.Region)
	payload := albwaf.CreateManagedRuleSetPayload{
		Name: model.Name,
		Type: model.Type,
	}
	return req.CreateManagedRuleSetPayload(payload)
}

func outputResult(p *print.Printer, outputFormat, projectLabel string, resp *albwaf.GetManagedRuleSetResponse) error {
	if resp == nil {
		return fmt.Errorf("create managed rule set response is empty")
	}
	return p.OutputResult(outputFormat, resp, func() error {
		p.Outputf("Created managed rule set %q for project %q.\n", resp.Name, projectLabel)
		return nil
	})
}
