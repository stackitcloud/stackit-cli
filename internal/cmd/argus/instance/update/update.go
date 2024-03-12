package update

import (
	"context"
	"errors"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/confirm"
	cliErr "github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/argus/client"
	argusUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/argus/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/spinner"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/argus"
	"github.com/stackitcloud/stackit-sdk-go/services/argus/wait"
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

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("update %s", instanceIdArg),
		Short: "Updates an Argus instance",
		Long:  "Updates an Argus instance.",
		Args:  args.SingleArg(instanceIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Update the plan of an Argus instance with ID "xxx" by specifying the plan ID`,
				"$ stackit argus instance update xxx --plan-id yyy"),
			examples.NewExample(
				`Update the plan of an Argus instance with ID "xxx" by specifying the plan name`,
				"$ stackit argus instance update xxx --plan-name yyy"),
			examples.NewExample(
				`Update the name of an Argus instance with ID "xxx"`,
				"$ stackit argus instance update xxx --name new-instance-name"),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			model, err := parseInput(cmd, args)
			if err != nil {
				return err
			}

			// Configure API client
			apiClient, err := client.ConfigureClient(cmd)
			if err != nil {
				return err
			}

			instanceLabel, err := argusUtils.GetInstanceName(ctx, apiClient, model.InstanceId, model.ProjectId)
			if err != nil {
				instanceLabel = model.InstanceId
			}

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to update instance %q?", instanceLabel)
				err = confirm.PromptForConfirmation(cmd, prompt)
				if err != nil {
					return err
				}
			}

			// Call API
			req, err := buildRequest(ctx, model, apiClient)
			if err != nil {
				var argusInvalidPlanError *cliErr.ArgusInvalidPlanError
				if !errors.As(err, &argusInvalidPlanError) {
					return fmt.Errorf("build Argus instance update request: %w", err)
				}
				return err
			}

			_, err = req.Execute()
			if err != nil {
				return fmt.Errorf("update Argus instance: %w", err)
			}
			instanceId := model.InstanceId

			// Wait for async operation, if async mode not enabled
			if !model.Async {
				s := spinner.New(cmd)
				s.Start("Updating instance")
				_, err = wait.UpdateInstanceWaitHandler(ctx, apiClient, instanceId, model.ProjectId).WaitWithContext(ctx)
				if err != nil {
					return fmt.Errorf("wait for Argus instance update: %w", err)
				}
				s.Stop()
			}

			operationState := "Updated"
			if model.Async {
				operationState = "Triggered update of"
			}
			cmd.Printf("%s instance %q\n", operationState, instanceLabel)
			return nil
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), planIdFlag, "Plan ID")
	cmd.Flags().String(planNameFlag, "", "Plan name")
	cmd.Flags().String(instanceNameFlag, "", "Instance name")
}

func parseInput(cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	instanceId := inputArgs[0]

	globalFlags := globalflags.Parse(cmd)
	if globalFlags.ProjectId == "" {
		return nil, &cliErr.ProjectIdError{}
	}

	planId := flags.FlagToStringPointer(cmd, planIdFlag)
	planName := flags.FlagToStringValue(cmd, planNameFlag)
	instanceName := flags.FlagToStringPointer(cmd, instanceNameFlag)

	if planId != nil && (planName != "") {
		return nil, &cliErr.ArgusInputPlanError{
			Cmd: cmd,
		}
	}

	if planId == nil && planName == "" {
		if instanceName == nil {
			return nil, &cliErr.EmptyUpdateError{}
		} else if *instanceName == "" {
			return nil, fmt.Errorf("new instance name cannot be empty.")
		}
	}

	return &inputModel{
		GlobalFlagModel: globalFlags,
		InstanceId:      instanceId,
		PlanId:          planId,
		PlanName:        planName,
		InstanceName:    instanceName,
	}, nil
}

type argusClient interface {
	UpdateInstance(ctx context.Context, instanceId, projectId string) argus.ApiUpdateInstanceRequest
	ListPlansExecute(ctx context.Context, projectId string) (*argus.PlansResponse, error)
	GetInstanceExecute(ctx context.Context, instanceId, projectId string) (*argus.GetInstanceResponse, error)
}

func buildRequest(ctx context.Context, model *inputModel, apiClient argusClient) (argus.ApiUpdateInstanceRequest, error) {
	req := apiClient.UpdateInstance(ctx, model.InstanceId, model.ProjectId)

	var planId *string
	var err error

	plans, err := apiClient.ListPlansExecute(ctx, model.ProjectId)
	if err != nil {
		return req, fmt.Errorf("get Argus plans: %w", err)
	}

	if model.PlanId == nil && model.PlanName != "" {
		planId, err = argusUtils.LoadPlanId(model.PlanName, plans)
		if err != nil {
			var argusInvalidPlanError *cliErr.ArgusInvalidPlanError
			if !errors.As(err, &argusInvalidPlanError) {
				return req, fmt.Errorf("load plan ID: %w", err)
			}
			return req, err
		}
	} else if model.PlanId == nil && model.PlanName == "" {
		planId, err = argusUtils.GetInstancePlanId(ctx, apiClient, model.InstanceId, model.ProjectId)
		if err != nil {
			return req, fmt.Errorf("get Argus instance plan ID: %w", err)
		}
	} else {
		if model.PlanId != nil {
			err := argusUtils.ValidatePlanId(*model.PlanId, plans)
			if err != nil {
				return req, err
			}
		}
		planId = model.PlanId
	}

	req = req.UpdateInstancePayload(argus.UpdateInstancePayload{
		PlanId: planId,
		Name:   model.InstanceName,
	})
	return req, nil
}
