package create

import (
	"context"
	"errors"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	cliErr "github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/projectname"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/argus/client"
	argusUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/argus/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/spinner"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/argus"
	"github.com/stackitcloud/stackit-sdk-go/services/argus/wait"
)

const (
	instanceNameFlag = "name"
	planIdFlag       = "plan-id"
	planNameFlag     = "plan-name"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	PlanName string

	InstanceName *string
	PlanId       *string
}

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Creates an Argus instance",
		Long:  "Creates an Argus instance.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Create an Argus instance with name "my-instance" and specify plan by name`,
				"$ stackit argus instance create --name my-instance --plan-name Monitoring-Starter-EU01"),
			examples.NewExample(
				`Create an Argus instance with name "my-instance" and specify plan by ID`,
				"$ stackit argus instance create --name my-instance --plan-id xxx"),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			model, err := parseInput(cmd)
			if err != nil {
				return err
			}

			// Configure API client
			apiClient, err := client.ConfigureClient(p)
			if err != nil {
				return err
			}

			projectLabel, err := projectname.GetProjectName(ctx, cmd, p)
			if err != nil {
				projectLabel = model.ProjectId
			}

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to create an Argus instance for project %q?", projectLabel)
				err = p.PromptForConfirmation(prompt)
				if err != nil {
					return err
				}
			}

			// Call API
			req, err := buildRequest(ctx, model, apiClient)
			if err != nil {
				var argusInvalidPlanError *cliErr.ArgusInvalidPlanError
				if !errors.As(err, &argusInvalidPlanError) {
					return fmt.Errorf("build Argus instance creation request: %w", err)
				}
				return err
			}
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("create Argus instance: %w", err)
			}
			instanceId := *resp.InstanceId

			// Wait for async operation, if async mode not enabled
			if !model.Async {
				s := spinner.New(p)
				s.Start("Creating instance")
				_, err = wait.CreateInstanceWaitHandler(ctx, apiClient, instanceId, model.ProjectId).WaitWithContext(ctx)
				if err != nil {
					return fmt.Errorf("wait for Argus instance creation: %w", err)
				}
				s.Stop()
			}

			operationState := "Created"
			if model.Async {
				operationState = "Triggered creation of"
			}
			p.Outputf("%s instance for project %q. Instance ID: %s\n", operationState, projectLabel, instanceId)
			return nil
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().StringP(instanceNameFlag, "n", "", "Instance name")
	cmd.Flags().Var(flags.UUIDFlag(), planIdFlag, "Plan ID")
	cmd.Flags().String(planNameFlag, "", "Plan name")

	err := flags.MarkFlagsRequired(cmd, instanceNameFlag)
	cobra.CheckErr(err)
}

func parseInput(cmd *cobra.Command) (*inputModel, error) {
	globalFlags := globalflags.Parse(cmd)
	if globalFlags.ProjectId == "" {
		return nil, &cliErr.ProjectIdError{}
	}

	planId := flags.FlagToStringPointer(cmd, planIdFlag)
	planName := flags.FlagToStringValue(cmd, planNameFlag)

	if planId == nil && (planName == "") {
		return nil, &cliErr.ArgusInputPlanError{
			Cmd: cmd,
		}
	}
	if planId != nil && (planName != "") {
		return nil, &cliErr.ArgusInputPlanError{
			Cmd: cmd,
		}
	}

	return &inputModel{
		GlobalFlagModel: globalFlags,
		InstanceName:    flags.FlagToStringPointer(cmd, instanceNameFlag),
		PlanId:          planId,
		PlanName:        planName,
	}, nil
}

type argusClient interface {
	CreateInstance(ctx context.Context, projectId string) argus.ApiCreateInstanceRequest
	ListPlansExecute(ctx context.Context, projectId string) (*argus.PlansResponse, error)
}

func buildRequest(ctx context.Context, model *inputModel, apiClient argusClient) (argus.ApiCreateInstanceRequest, error) {
	req := apiClient.CreateInstance(ctx, model.ProjectId)

	var planId *string
	var err error

	plans, err := apiClient.ListPlansExecute(ctx, model.ProjectId)
	if err != nil {
		return req, fmt.Errorf("get Argus plans: %w", err)
	}

	if model.PlanId == nil {
		planId, err = argusUtils.LoadPlanId(model.PlanName, plans)
		if err != nil {
			var argusInvalidPlanError *cliErr.ArgusInvalidPlanError
			if !errors.As(err, &argusInvalidPlanError) {
				return req, fmt.Errorf("load plan ID: %w", err)
			}
			return req, err
		}
	} else {
		err := argusUtils.ValidatePlanId(*model.PlanId, plans)
		if err != nil {
			return req, err
		}
		planId = model.PlanId
	}

	req = req.CreateInstancePayload(argus.CreateInstancePayload{
		Name:   model.InstanceName,
		PlanId: planId,
	})
	return req, nil
}
