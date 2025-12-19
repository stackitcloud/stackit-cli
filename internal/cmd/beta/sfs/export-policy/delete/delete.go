package delete

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/sfs/client"
	sfsUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/sfs/utils"
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
		Use:   fmt.Sprintf("delete %s", exportPolicyIdArg),
		Short: "Deletes a export policy",
		Long:  "Deletes a export policy.",
		Args:  args.SingleArg(exportPolicyIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Delete a export policy with ID "xxx"`,
				"$ stackit beta sfs export-policy delete xxx",
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

			exportPolicyLabel, err := sfsUtils.GetExportPolicyName(ctx, apiClient, model.ProjectId, model.Region, model.ExportPolicyId)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get export policy name: %v", err)
				exportPolicyLabel = model.ExportPolicyId
			} else if exportPolicyLabel == "" {
				exportPolicyLabel = model.ExportPolicyId
			}

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to delete export policy %q? (This cannot be undone)", exportPolicyLabel)
				err = params.Printer.PromptForConfirmation(prompt)
				if err != nil {
					return err
				}
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			_, err = req.Execute()
			if err != nil {
				return fmt.Errorf("delete export policy: %w", err)
			}

			params.Printer.Outputf("Deleted export policy %q\n", exportPolicyLabel)
			return nil
		},
	}
	return cmd
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *sfs.APIClient) sfs.ApiDeleteShareExportPolicyRequest {
	return apiClient.DeleteShareExportPolicy(ctx, model.ProjectId, model.Region, model.ExportPolicyId)
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
