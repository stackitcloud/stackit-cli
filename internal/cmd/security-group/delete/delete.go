package delete

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/types"

	"github.com/spf13/cobra"
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
	SecurityGroupId string
}

const groupIdArg = "GROUP_ID"

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("delete %s", groupIdArg),
		Short: "Deletes a security group",
		Long:  "Deletes a security group by its internal ID.",
		Args:  args.SingleArg(groupIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(`Delete a named group with ID "xxx"`, `$ stackit security-group delete xxx`),
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

			groupLabel, err := iaasUtils.GetSecurityGroupName(ctx, apiClient, model.ProjectId, model.Region, model.SecurityGroupId)
			if err != nil {
				params.Printer.Warn("get security group name: %v", err)
				groupLabel = model.SecurityGroupId
			}

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to delete the security group %q for %q?", groupLabel, projectLabel)
				err = params.Printer.PromptForConfirmation(prompt)
				if err != nil {
					return err
				}
			}

			// Call API
			request := buildRequest(ctx, model, apiClient)

			if err := request.Execute(); err != nil {
				return fmt.Errorf("delete security group: %w", err)
			}
			params.Printer.Info("Deleted security group %q for %q\n", groupLabel, projectLabel)

			return nil
		},
	}

	return cmd
}

func parseInput(p *print.Printer, cmd *cobra.Command, cliArgs []string) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		SecurityGroupId: cliArgs[0],
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *iaas.APIClient) iaas.ApiDeleteSecurityGroupRequest {
	request := apiClient.DeleteSecurityGroup(ctx, model.ProjectId, model.Region, model.SecurityGroupId)
	return request
}
