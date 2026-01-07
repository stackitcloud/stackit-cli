package describe

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/projectname"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/sfs/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/sfs"
)

const exportPolicyIdArg = "EXPORT_POLICY_ID"

type inputModel struct {
	*globalflags.GlobalFlagModel
	ExportPolicyId string
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("describe %s", exportPolicyIdArg),
		Short: "Shows details of a export policy",
		Long:  "Shows details of a export policy.",
		Args:  args.SingleArg(exportPolicyIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Describe a export policy with ID "xxx"`,
				"$ stackit beta sfs export-policy describe xxx",
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			model, err := parseInput(params.Printer, cmd, args)
			if err != nil {
				return fmt.Errorf("unable to parse input: %w", err)
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
				return fmt.Errorf("read export policy: %w", err)
			}

			// Get projectLabel
			projectLabel, err := projectname.GetProjectName(ctx, params.Printer, params.CliVersion, cmd)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get project name: %v", err)
				projectLabel = model.ProjectId
			} else if projectLabel == "" {
				projectLabel = model.ProjectId
			}

			return outputResult(params.Printer, model.OutputFormat, model.ExportPolicyId, projectLabel, resp)
		},
	}
	return cmd
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	exportPolicyId := inputArgs[0]

	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		ExportPolicyId:  exportPolicyId,
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *sfs.APIClient) sfs.ApiGetShareExportPolicyRequest {
	return apiClient.GetShareExportPolicy(ctx, model.ProjectId, model.Region, model.ExportPolicyId)
}

func outputResult(p *print.Printer, outputFormat, exportPolicyId, projectLabel string, exportPolicy *sfs.GetShareExportPolicyResponse) error {
	return p.OutputResult(outputFormat, exportPolicy, func() error {
		if exportPolicy == nil || exportPolicy.ShareExportPolicy == nil {
			p.Outputf("Export policy %q not found in project %q", exportPolicyId, projectLabel)
			return nil
		}

		var content []tables.Table

		table := tables.NewTable()
		table.SetTitle("Export Policy")
		policy := exportPolicy.ShareExportPolicy

		table.AddRow("ID", utils.PtrString(policy.Id))
		table.AddSeparator()
		table.AddRow("NAME", utils.PtrString(policy.Name))
		table.AddSeparator()
		table.AddRow("SHARES USING EXPORT POLICY", utils.PtrString(policy.SharesUsingExportPolicy))
		table.AddSeparator()
		table.AddRow("CREATED AT", utils.ConvertTimePToDateTimeString(policy.CreatedAt))

		content = append(content, table)

		if policy.Rules != nil && len(*policy.Rules) > 0 {
			rulesTable := tables.NewTable()
			rulesTable.SetTitle("Rules")

			rulesTable.SetHeader("ID", "ORDER", "DESCRIPTION", "IP ACL", "READ ONLY", "SET UUID", "SUPER USER", "CREATED AT")

			for _, rule := range *policy.Rules {
				var description string
				if rule.Description != nil {
					description = utils.PtrString(rule.Description.Get())
				}
				rulesTable.AddRow(
					utils.PtrString(rule.Id),
					utils.PtrString(rule.Order),
					description,
					utils.JoinStringPtr(rule.IpAcl, ", "),
					utils.PtrString(rule.ReadOnly),
					utils.PtrString(rule.SetUuid),
					utils.PtrString(rule.SuperUser),
					utils.ConvertTimePToDateTimeString(rule.CreatedAt),
				)
				rulesTable.AddSeparator()
			}

			content = append(content, rulesTable)
		}

		if err := tables.DisplayTables(p, content); err != nil {
			return fmt.Errorf("render tables: %w", err)
		}
		return nil
	})
}
