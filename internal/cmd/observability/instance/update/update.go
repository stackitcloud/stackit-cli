package update

import (
	"context"
	"errors"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/types"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	cliErr "github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/observability/client"
	observabilityUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/observability/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/spinner"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/observability"
	"github.com/stackitcloud/stackit-sdk-go/services/observability/wait"
)

const (
	instanceIdArg = "INSTANCE_ID"

	instanceNameFlag = "name"
	planIdFlag       = "plan-id"
	planNameFlag     = "plan-name"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	InstanceId string
	PlanName   string

	InstanceName *string
	PlanId       *string
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("update %s", instanceIdArg),
		Short: "Updates an Observability instance",
		Long:  "Updates an Observability instance.",
		Args:  args.SingleArg(instanceIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Update the plan of an Observability instance with ID "xxx" by specifying the plan ID`,
				"$ stackit observability instance update xxx --plan-id yyy"),
			examples.NewExample(
				`Update the plan of an Observability instance with ID "xxx" by specifying the plan name`,
				"$ stackit observability instance update xxx --plan-name Frontend-Starter-EU01"),
			examples.NewExample(
				`Update the name of an Observability instance with ID "xxx"`,
				"$ stackit observability instance update xxx --name new-instance-name"),
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

			instanceLabel, err := observabilityUtils.GetInstanceName(ctx, apiClient, model.InstanceId, model.ProjectId)
			if err != nil || instanceLabel == "" {
				params.Printer.Debug(print.ErrorLevel, "get instance name: %v", err)
				instanceLabel = model.InstanceId
			}

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to update instance %q?", instanceLabel)
				err = params.Printer.PromptForConfirmation(prompt)
				if err != nil {
					return err
				}
			}

			// Call API
			req, err := buildRequest(ctx, model, apiClient)
			if err != nil {
				var observabilityInvalidPlanError *cliErr.ObservabilityInvalidPlanError
				if !errors.As(err, &observabilityInvalidPlanError) {
					return fmt.Errorf("build Observability instance update request: %w", err)
				}
				return err
			}

			_, err = req.Execute()
			if err != nil {
				return fmt.Errorf("update Observability instance: %w", err)
			}
			instanceId := model.InstanceId

			// Wait for async operation, if async mode not enabled
			if !model.Async {
				s := spinner.New(params.Printer)
				s.Start("Updating instance")
				_, err = wait.UpdateInstanceWaitHandler(ctx, apiClient, instanceId, model.ProjectId).WaitWithContext(ctx)
				if err != nil {
					return fmt.Errorf("wait for Observability instance update: %w", err)
				}
				s.Stop()
			}

			operationState := "Updated"
			if model.Async {
				operationState = "Triggered update of"
			}
			params.Printer.Info("%s instance %q\n", operationState, instanceLabel)
			return nil
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), planIdFlag, "Plan ID")
	cmd.Flags().String(planNameFlag, "", "Plan name")
	cmd.Flags().StringP(instanceNameFlag, "n", "", "Instance name")
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	instanceId := inputArgs[0]

	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &cliErr.ProjectIdError{}
	}

	planId := flags.FlagToStringPointer(p, cmd, planIdFlag)
	planName := flags.FlagToStringValue(p, cmd, planNameFlag)
	instanceName := flags.FlagToStringPointer(p, cmd, instanceNameFlag)

	if planId != nil && (planName != "") {
		return nil, &cliErr.ObservabilityInputPlanError{
			Cmd: cmd,
		}
	}

	if planId == nil && planName == "" && instanceName == nil {
		return nil, &cliErr.EmptyUpdateError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		InstanceId:      instanceId,
		PlanId:          planId,
		PlanName:        planName,
		InstanceName:    instanceName,
	}

	p.DebugInputModel(model)
	return &model, nil
}

type observabilityClient interface {
	UpdateInstance(ctx context.Context, instanceId, projectId string) observability.ApiUpdateInstanceRequest
	ListPlansExecute(ctx context.Context, projectId string) (*observability.PlansResponse, error)
	GetInstanceExecute(ctx context.Context, instanceId, projectId string) (*observability.GetInstanceResponse, error)
}

func buildRequest(ctx context.Context, model *inputModel, apiClient observabilityClient) (observability.ApiUpdateInstanceRequest, error) {
	req := apiClient.UpdateInstance(ctx, model.InstanceId, model.ProjectId)

	var err error

	plans, err := apiClient.ListPlansExecute(ctx, model.ProjectId)
	if err != nil {
		return req, fmt.Errorf("get Observability plans: %w", err)
	}

	currentInstance, err := apiClient.GetInstanceExecute(ctx, model.InstanceId, model.ProjectId)
	if err != nil {
		return req, fmt.Errorf("get Observability instance: %w", err)
	}

	payload := observability.UpdateInstancePayload{
		PlanId: currentInstance.PlanId,
		Name:   currentInstance.Name,
	}

	if model.PlanId == nil && model.PlanName != "" {
		payload.PlanId, err = observabilityUtils.LoadPlanId(model.PlanName, plans)
		if err != nil {
			var observabilityInvalidPlanError *cliErr.ObservabilityInvalidPlanError
			if !errors.As(err, &observabilityInvalidPlanError) {
				return req, fmt.Errorf("load plan ID: %w", err)
			}
			return req, err
		}
	} else if model.PlanId != nil && model.PlanName == "" {
		err := observabilityUtils.ValidatePlanId(*model.PlanId, plans)
		if err != nil {
			var observabilityInvalidPlanError *cliErr.ObservabilityInvalidPlanError
			if !errors.As(err, &observabilityInvalidPlanError) {
				return req, fmt.Errorf("validate plan ID: %w", err)
			}
			return req, err
		}
		payload.PlanId = model.PlanId
	}

	if model.InstanceName != nil {
		payload.Name = model.InstanceName
	}

	req = req.UpdateInstancePayload(payload)
	return req, nil
}
