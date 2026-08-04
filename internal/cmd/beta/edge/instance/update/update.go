package update

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	cliErr "github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/projectname"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/edge/client"
	edgeUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/edge/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"

	edge "github.com/stackitcloud/stackit-sdk-go/services/edge/v1beta1api"
	"github.com/stackitcloud/stackit-sdk-go/services/edge/v1beta1api/wait"

	"github.com/stackitcloud/stackit-cli/internal/pkg/spinner"
)

const (
	instanceIdArg = "INSTANCE_ID"

	descriptionFlag = "description"
	planIdFlag      = "plan-id"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	InstanceId  string
	Description *string
	PlanId      *string
}

// NewCmd https://aip.stackit.cloud/aip/general/0121/
// We have decided to eliminate the usage of display name flag
// To be the AIP compliant, and align with the standard CLI implementation, we will use the InstanceID arg
func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("update %s", instanceIdArg),
		Short: "Updates an Edge Cloud instance",
		Long:  "Updates a STACKIT Edge Cloud (STEC) instance.",
		Args:  args.SingleArg(instanceIdArg, nil),
		Example: examples.Build(
			examples.NewExample(
				`Updates the description of an Edge Cloud instance with ID "xxx"`,
				`$ stackit beta edge-cloud instance update xxx --description yyy`),
			examples.NewExample(
				`Updates the plan of an Edge Cloud instance with ID "xxx"`,
				`$ stackit beta edge-cloud instance update xxx --plan-id yyy`),
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

			instanceLabel, err := edgeUtils.GetInstanceName(ctx, apiClient.DefaultAPI, model.ProjectId, model.Region, model.InstanceId)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get instance name: %v", err)
				instanceLabel = model.InstanceId
			}

			prompt := fmt.Sprintf("Are you sure you want to update the Edge Cloud instance %q of project %q?", instanceLabel, projectLabel)
			err = params.Printer.PromptForConfirmation(prompt)
			if err != nil {
				return err
			}

			// Call API
			req, err := buildRequest(ctx, model, apiClient)
			if err != nil {
				return err
			}
			err = req.Execute()
			if err != nil {
				return fmt.Errorf("update Edge Cloud instance: %w", err)
			}

			if !model.Async {
				err := spinner.Run(params.Printer, "Updating instance", func() error {
					_, err = wait.CreateOrUpdateInstanceWaitHandler(ctx, apiClient.DefaultAPI, model.ProjectId, model.Region, model.InstanceId).WaitWithContext(ctx)
					return err
				})
				if err != nil {
					return fmt.Errorf("wait for Edge Cloud instance update: %w", err)
				}
			}

			operationState := "Updated"
			if model.Async {
				operationState = "Triggered update of"
			}
			params.Printer.Outputf("%s instance %q\n", operationState, instanceLabel)
			return nil
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().StringP(descriptionFlag, "d", "", "A user chosen description to distinguish multiple instances.")
	cmd.Flags().String(planIdFlag, "", "Service plan configures the size of the Instance.")

	// Make sure at least one updatable field is provided, otherwise it would be a no-op
	cmd.MarkFlagsOneRequired(descriptionFlag, planIdFlag)
}

// Parse user input (arguments and/or flags)
func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	instanceId := inputArgs[0]

	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &cliErr.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		PlanId:          flags.FlagToStringPointer(p, cmd, planIdFlag),
		InstanceId:      instanceId,
		Description:     flags.FlagToStringPointer(p, cmd, descriptionFlag),
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *edge.APIClient) (req edge.ApiUpdateInstanceRequest, err error) {
	req = apiClient.DefaultAPI.UpdateInstance(ctx, model.ProjectId, model.Region, model.InstanceId)

	payload := edge.UpdateInstancePayload{
		Description: model.Description,
		PlanId:      model.PlanId,
	}

	return req.UpdateInstancePayload(payload), nil
}
