package create

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/goccy/go-yaml"
	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	cliErr "github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/projectname"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/observability/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/spinner"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
	observabilityUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/observability/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/observability"
	"github.com/stackitcloud/stackit-sdk-go/services/observability/wait"
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

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Creates an Observability instance",
		Long:  "Creates an Observability instance.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Create an Observability instance with name "my-instance" and specify plan by name`,
				"$ stackit observability instance create --name my-instance --plan-name Monitoring-Starter-EU01"),
			examples.NewExample(
				`Create an Observability instance with name "my-instance" and specify plan by ID`,
				"$ stackit observability instance create --name my-instance --plan-id xxx"),
		),
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := context.Background()
			model, err := parseInput(params.Printer, cmd)
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

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to create an Observability instance for project %q?", projectLabel)
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
					return fmt.Errorf("build Observability instance creation request: %w", err)
				}
				return err
			}
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("create Observability instance: %w", err)
			}
			instanceId := *resp.InstanceId

			// Wait for async operation, if async mode not enabled
			if !model.Async {
				s := spinner.New(params.Printer)
				s.Start("Creating instance")
				_, err = wait.CreateInstanceWaitHandler(ctx, apiClient, instanceId, model.ProjectId).WaitWithContext(ctx)
				if err != nil {
					return fmt.Errorf("wait for Observability instance creation: %w", err)
				}
				s.Stop()
			}

			return outputResult(params.Printer, model.OutputFormat, model.Async, projectLabel, resp)
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

func parseInput(p *print.Printer, cmd *cobra.Command) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &cliErr.ProjectIdError{}
	}

	planId := flags.FlagToStringPointer(p, cmd, planIdFlag)
	planName := flags.FlagToStringValue(p, cmd, planNameFlag)

	if planId == nil && (planName == "") {
		return nil, &cliErr.ObservabilityInputPlanError{
			Cmd: cmd,
		}
	}
	if planId != nil && (planName != "") {
		return nil, &cliErr.ObservabilityInputPlanError{
			Cmd: cmd,
		}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		InstanceName:    flags.FlagToStringPointer(p, cmd, instanceNameFlag),
		PlanId:          planId,
		PlanName:        planName,
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

type observabilityClient interface {
	CreateInstance(ctx context.Context, projectId string) observability.ApiCreateInstanceRequest
	ListPlansExecute(ctx context.Context, projectId string) (*observability.PlansResponse, error)
}

func buildRequest(ctx context.Context, model *inputModel, apiClient observabilityClient) (observability.ApiCreateInstanceRequest, error) {
	req := apiClient.CreateInstance(ctx, model.ProjectId)

	var planId *string
	var err error

	plans, err := apiClient.ListPlansExecute(ctx, model.ProjectId)
	if err != nil {
		return req, fmt.Errorf("get Observability plans: %w", err)
	}

	if model.PlanId == nil {
		planId, err = observabilityUtils.LoadPlanId(model.PlanName, plans)
		if err != nil {
			var observabilityInvalidPlanError *cliErr.ObservabilityInvalidPlanError
			if !errors.As(err, &observabilityInvalidPlanError) {
				return req, fmt.Errorf("load plan ID: %w", err)
			}
			return req, err
		}
	} else {
		err := observabilityUtils.ValidatePlanId(*model.PlanId, plans)
		if err != nil {
			return req, err
		}
		planId = model.PlanId
	}

	req = req.CreateInstancePayload(observability.CreateInstancePayload{
		Name:   model.InstanceName,
		PlanId: planId,
	})
	return req, nil
}

func outputResult(p *print.Printer, outputFormat string, async bool, projectLabel string, resp *observability.CreateInstanceResponse) error {
	if resp == nil {
		return fmt.Errorf("resp is empty")
	}

	switch outputFormat {
	case print.JSONOutputFormat:
		details, err := json.MarshalIndent(resp, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal Observability instance: %w", err)
		}
		p.Outputln(string(details))

		return nil
	case print.YAMLOutputFormat:
		details, err := yaml.MarshalWithOptions(resp, yaml.IndentSequence(true), yaml.UseJSONMarshaler())
		if err != nil {
			return fmt.Errorf("marshal Observability instance: %w", err)
		}
		p.Outputln(string(details))

		return nil
	default:
		operationState := "Created"
		if async {
			operationState = "Triggered creation of"
		}
		p.Outputf("%s instance for project %q. Instance ID: %s\n", operationState, projectLabel, utils.PtrString(resp.InstanceId))
		return nil
	}
}
