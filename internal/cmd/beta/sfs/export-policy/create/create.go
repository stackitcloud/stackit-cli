package create

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/projectname"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/sfs/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/sfs"
)

const (
	nameFlag  = "name"
	rulesFlag = "rules"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	Name  string
	Rules *[]sfs.CreateShareExportPolicyRequestRule
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Creates a export policy",
		Long:  "Creates a export policy.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Create a new export policy with name "EXPORT_POLICY_NAME"`,
				"$ stackit beta sfs export-policy create --name EXPORT_POLICY_NAME",
			),
			examples.NewExample(
				`Create a new export policy with name "EXPORT_POLICY_NAME" and rules from file "./rules.json"`,
				"$ stackit beta sfs export-policy create --name EXPORT_POLICY_NAME --rules @./rules.json",
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			model, err := parseInput(params.Printer, cmd, args)
			if err != nil {
				return fmt.Errorf("unable to parse input: %w", err)
			}

			// Configure client
			apiClient, err := client.ConfigureClient(params.Printer, params.CliVersion)
			if err != nil {
				return err
			}

			projectLabel, err := projectname.GetProjectName(ctx, params.Printer, params.CliVersion, cmd)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get project name: %v", err)
				projectLabel = model.ProjectId
			} else if projectLabel == "" {
				projectLabel = model.ProjectId
			}

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to create a export policy for project %q?", projectLabel)
				err = params.Printer.PromptForConfirmation(prompt)
				if err != nil {
					return err
				}
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("create export policy: %w", err)
			}

			return outputResult(params.Printer, model.OutputFormat, projectLabel, resp)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().String(nameFlag, "", "Export policy name")
	cmd.Flags().Var(flags.ReadFromFileFlag(), rulesFlag, "Rules of the export policy (format: json)")

	err := flags.MarkFlagsRequired(cmd, nameFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command, _ []string) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	rulesString := flags.FlagToStringPointer(p, cmd, rulesFlag)
	var rules *[]sfs.CreateShareExportPolicyRequestRule
	if rulesString != nil && *rulesString != "" {
		var r []sfs.CreateShareExportPolicyRequestRule
		err := json.Unmarshal([]byte(*rulesString), &r)
		if err != nil {
			return nil, fmt.Errorf("could not parse rules: %w", err)
		}
		rules = &r
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		Name:            flags.FlagToStringValue(p, cmd, nameFlag),
		Rules:           rules,
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *sfs.APIClient) sfs.ApiCreateShareExportPolicyRequest {
	req := apiClient.CreateShareExportPolicy(ctx, model.ProjectId, model.Region)
	req = req.CreateShareExportPolicyPayload(
		sfs.CreateShareExportPolicyPayload{
			Name:  utils.Ptr(model.Name),
			Rules: model.Rules,
		},
	)
	return req
}

func outputResult(p *print.Printer, outputFormat, projectLabel string, item *sfs.CreateShareExportPolicyResponse) error {
	return p.OutputResult(outputFormat, item, func() error {
		if item == nil || item.ShareExportPolicy == nil {
			return fmt.Errorf("no export policy found")
		}
		p.Outputf(
			"Created export policy %q for project %q.\nExport policy ID: %s\n",
			utils.PtrString(item.ShareExportPolicy.Name),
			projectLabel,
			utils.PtrString(item.ShareExportPolicy.Id),
		)
		return nil
	})
}
