package delete

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/cmd/params"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/projectname"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/git/client"
	gitUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/git/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/spinner"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/git"
	"github.com/stackitcloud/stackit-sdk-go/services/git/wait"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	InstanceId string
}

const instanceIdArg = "INSTANCE_ID"

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("delete %s", instanceIdArg),
		Short: "Deletes STACKIT Git instance",
		Long:  "Deletes a STACKIT Git instance by its internal ID.",
		Args:  args.SingleArg(instanceIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(`Delete a instance with ID "xxx"`,
				`$ stackit git instance delete xxx`),
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

			projectName, err := projectname.GetProjectName(ctx, params.Printer, params.CliVersion, cmd)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get project name: %v", err)
				projectName = model.ProjectId
			}

			instanceName, err := gitUtils.GetInstanceName(ctx, apiClient, model.ProjectId, model.InstanceId)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get stackit git intance name: %v", err)
				instanceName = model.InstanceId
			} else if instanceName == "" {
				instanceName = model.InstanceId
			}

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to delete the stackit git instance %q for %q?", instanceName, projectName)
				err = params.Printer.PromptForConfirmation(prompt)
				if err != nil {
					return err
				}
			}

			// Call API
			request := buildRequest(ctx, model, apiClient)

			err = request.Execute()
			if err != nil {
				return fmt.Errorf("delete instance: %w", err)
			}

			// Wait for async operation, if async mode not enabled
			if !model.Async {
				s := spinner.New(params.Printer)
				s.Start("Deleting stackit git instance")
				_, err = wait.DeleteGitInstanceWaitHandler(ctx, apiClient, model.ProjectId, model.InstanceId).WaitWithContext(ctx)
				if err != nil {
					return fmt.Errorf("wait for stackit git instance deletion: %w", err)
				}
				s.Stop()
			}

			operationState := "Deleted"
			if model.Async {
				operationState = "Triggered deletion of"
			}
			params.Printer.Info("%s stackit git instance %s \n", operationState, model.InstanceId)

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
		InstanceId:      cliArgs[0],
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *git.APIClient) git.ApiDeleteInstanceRequest {
	return apiClient.DeleteInstance(ctx, model.ProjectId, model.InstanceId)
}
