package delete

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/projectname"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/iaas/client"
	iaasUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/iaas/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/iaas"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	AffinityGroupId string
}

const (
	affinityGroupIdArg = "AFFINITY_GROUP"
)

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("delete %s", affinityGroupIdArg),
		Short: "Deletes an affinity group",
		Long:  `Deletes an affinity group.`,
		Args:  args.SingleArg(affinityGroupIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Delete an affinity group with ID "xxx"`,
				"$ stackit affinity-group delete xxx",
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

			projectLabel, err := projectname.GetProjectName(ctx, params.Printer, cmd)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get project name: %v", err)
				projectLabel = model.ProjectId
			}

			affinityGroupLabel, err := iaasUtils.GetAffinityGroupName(ctx, apiClient, model.ProjectId, model.AffinityGroupId)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get affinity group name: %v", err)
				affinityGroupLabel = model.AffinityGroupId
			}

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to delete affinity group %q?", affinityGroupLabel)
				err = params.Printer.PromptForConfirmation(prompt)
				if err != nil {
					return err
				}
			}

			// Call API
			request := buildRequest(ctx, *model, apiClient)
			err = request.Execute()
			if err != nil {
				return fmt.Errorf("delete affinity group: %w", err)
			}
			params.Printer.Info("Deleted affinity group %q for %q\n", affinityGroupLabel, projectLabel)

			return nil
		},
	}
	return cmd
}

func buildRequest(ctx context.Context, model inputModel, apiClient *iaas.APIClient) iaas.ApiDeleteAffinityGroupRequest {
	return apiClient.DeleteAffinityGroup(ctx, model.ProjectId, model.AffinityGroupId)
}

func parseInput(p *print.Printer, cmd *cobra.Command, cliArgs []string) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		AffinityGroupId: cliArgs[0],
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
