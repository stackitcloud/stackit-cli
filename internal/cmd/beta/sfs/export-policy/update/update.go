package update

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
	sfsUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/sfs/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/sfs"
)

const (
	exportPolicyArg = "EXPORT_POLICY_ID"

	rulesFlag       = "rules"
	removeRulesFlag = "remove-rules"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	ExportPolicyId string
	Rules          *[]sfs.UpdateShareExportPolicyBodyRule
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("update %s", exportPolicyArg),
		Short: "Updates a export policy",
		Long:  "Updates a export policy.",
		Args:  args.SingleArg(exportPolicyArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Update a export policy with ID "xxx" and with rules from file "./rules.json"`,
				"$ stackit beta sfs export-policy update xxx --rules @./rules.json",
			),
			examples.NewExample(
				`Update a export policy with ID "xxx" and remove the rules`,
				"$ stackit beta sfs export-policy update XXX --remove-rules",
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			model, err := parseInput(params.Printer, cmd, args)
			if err != nil {
				return fmt.Errorf("unable to parse input: %w", err)
			}

			// Configure client
			apiClient, err := client.ConfigureClient(params.Printer, params.CliVersion)
			if err != nil {
				return err
			}

			exportPolicyLabel, err := sfsUtils.GetExportPolicyName(ctx, apiClient, model.ProjectId, model.Region, model.ExportPolicyId)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get export policy name: %v", err)
				exportPolicyLabel = model.ExportPolicyId
			} else if exportPolicyLabel == "" {
				exportPolicyLabel = model.ExportPolicyId
			}

			projectLabel, err := projectname.GetProjectName(ctx, params.Printer, params.CliVersion, cmd)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get project name: %v", err)
				projectLabel = model.ProjectId
			} else if projectLabel == "" {
				projectLabel = model.ProjectId
			}

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to update export policy %q for project %q?", exportPolicyLabel, projectLabel)
				err = params.Printer.PromptForConfirmation(prompt)
				if err != nil {
					return err
				}
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("update export policy: %w", err)
			}

			return outputResult(params.Printer, model.OutputFormat, projectLabel, exportPolicyLabel, resp)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.ReadFromFileFlag(), rulesFlag, "Rules of the export policy")
	cmd.Flags().Bool(removeRulesFlag, false, "Remove the export policy rules")

	rulesFlags := []string{rulesFlag, removeRulesFlag}
	cmd.MarkFlagsMutuallyExclusive(rulesFlags...)
	cmd.MarkFlagsOneRequired(rulesFlags...) // Because the update endpoint supports only rules at the moment, one of the flags must be required
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *sfs.APIClient) sfs.ApiUpdateShareExportPolicyRequest {
	req := apiClient.UpdateShareExportPolicy(ctx, model.ProjectId, model.Region, model.ExportPolicyId)

	payload := sfs.UpdateShareExportPolicyPayload{
		Rules: model.Rules,
	}
	return req.UpdateShareExportPolicyPayload(payload)
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	exportPolicyId := inputArgs[0]

	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	var rules *[]sfs.UpdateShareExportPolicyBodyRule
	noRulesErr := fmt.Errorf("no rules specified")
	if rulesString := flags.FlagToStringPointer(p, cmd, rulesFlag); rulesString != nil {
		var r []sfs.UpdateShareExportPolicyBodyRule
		err := json.Unmarshal([]byte(*rulesString), &r)
		if err != nil {
			return nil, fmt.Errorf("could not parse rules: %w", err)
		}
		if r == nil {
			return nil, noRulesErr
		}
		rules = &r
	}

	if removeRules := flags.FlagToBoolPointer(p, cmd, removeRulesFlag); removeRules != nil {
		// Create an empty slice for the patch request
		rules = &[]sfs.UpdateShareExportPolicyBodyRule{}
	}

	// Because the update endpoint supports only rules at the moment, this should not be empty
	if rules == nil {
		return nil, noRulesErr
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		ExportPolicyId:  exportPolicyId,
		Rules:           rules,
	}

	p.DebugInputModel(model)
	return &model, nil
}

func outputResult(p *print.Printer, outputFormat, projectLabel, exportPolicyLabel string, resp *sfs.UpdateShareExportPolicyResponse) error {
	return p.OutputResult(outputFormat, resp, func() error {
		if resp == nil {
			p.Outputln("Empty export policy response")
			return nil
		}
		p.Outputf("Updated export policy %q for project %q\n", exportPolicyLabel, projectLabel)
		return nil
	})
}
