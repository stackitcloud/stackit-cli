package delete

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	cliErr "github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/projectname"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/modelexperiments/client"
	meUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/modelexperiments/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/spinner"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	modelexperiments "github.com/stackitcloud/stackit-sdk-go/services/modelexperiments/v1api"
	"github.com/stackitcloud/stackit-sdk-go/services/modelexperiments/v1api/wait"
)

const (
	instanceIdArg = "INSTANCE_ID"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	InstanceId string
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("delete %s", instanceIdArg),
		Short: "Deletes an AI Model Experiments instance",
		Long:  "Deletes an AI Model Experiments instance from a STACKIT project.",
		Args:  args.SingleArg(instanceIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Delete an AI Model Experiments instance with ID "xxx"`,
				`$ stackit model-experiments instance delete xxx`,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()

			model, err := parseInput(params.Printer, cmd, args)
			if err != nil {
				return err
			}

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

			instanceLabel, err := meUtils.GetInstanceName(ctx, apiClient.DefaultAPI, model.ProjectId, model.Region, model.InstanceId)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get instance name: %v", err)
				instanceLabel = model.InstanceId
			}

			prompt := fmt.Sprintf("Are you sure you want to delete AI Model Experiments instance %q from project %q?", instanceLabel, projectLabel)

			if err := params.Printer.PromptForConfirmation(prompt); err != nil {
				return err
			}

			req := buildRequest(ctx, model, apiClient)

			_, err = req.Execute()
			if err != nil {
				return fmt.Errorf("delete AI Model Experiments instance: %w", err)
			}

			if !model.Async {
				err := spinner.Run(params.Printer, "Deleting AI Model Experiments instance", func() error {
					_, err = wait.DeleteInstanceWaitHandler(ctx, apiClient.DefaultAPI, model.Region, model.ProjectId, model.InstanceId).WaitWithContext(ctx)
					return err
				})
				if err != nil {
					return fmt.Errorf("wait for AI Model Experiments instance deletion: %w", err)
				}
			}

			operationState := "Deleted"
			if model.Async {
				operationState = "Triggered deletion of"
			}
			params.Printer.Outputf("%s AI Model Experiments instance %q\n", operationState, instanceLabel)
			return nil
		},
	}

	return cmd
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)

	if globalFlags.ProjectId == "" {
		return nil, &cliErr.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		InstanceId:      inputArgs[0],
	}

	p.DebugInputModel(model)

	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *modelexperiments.APIClient) modelexperiments.ApiDeleteInstanceRequest {
	return apiClient.DefaultAPI.DeleteInstance(ctx, model.ProjectId, model.Region, model.InstanceId)
}
