package update

import (
	"context"
	"fmt"
	"slices"

	"github.com/spf13/cobra"
	sca "github.com/stackitcloud/stackit-sdk-go/services/sca/v1alphaapi"
	"github.com/stackitcloud/stackit-sdk-go/services/sca/v1alphaapi/wait"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/sca/client"
	scautils "github.com/stackitcloud/stackit-cli/internal/pkg/services/sca/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/spinner"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
)

const (
	applicationIDArg  = "APPLICATION_ID"
	environmentIDFlag = "environment-id"
	imageFlag         = "image"
	publicFlag        = "public"
	externalPortFlag  = "external-port"
	cpuFlag           = "cpu"
	memoryFlag        = "memory"
	instancesFlag     = "instances"
	minInstancesFlag  = "min-instances"
	maxInstancesFlag  = "max-instances"
	rpsFlag           = "http-rule-rps"
	concurrencyFlag   = "http-rule-concurrency"
	scaleToZeroFlag   = "scale-to-zero"
	envVarsFlag       = "environment-vars"
	commandsFlag      = "commands"
	argsFlag          = "args"
	stoppedFlag       = "stopped"
)

const (
	scalingTypeManual = "manual"
	scalingTypeAuto   = "auto"
)

var scalingTypeFlag = flags.StringEnumFlag(
	"scaling-type",
	[]string{scalingTypeManual, scalingTypeAuto},
	"Scaling type",
	flags.StringEnumDefaultValue(scalingTypeManual),
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	EnvironmentID         string
	ApplicationID         string
	Image                 string
	EnvironmentVariables  map[string]string
	Commands              []string
	Args                  []string
	ScalingType           *string
	Public                *bool
	ContainerExternalPort *int32
	CPU                   *int32
	Memory                *int32
	Instances             *int32
	MinInstances          *int32
	MaxInstances          *int32
	ScaleToZero           *bool
	Concurrency           *int32
	RPS                   *int32
	Stopped               *bool
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update a SCA application",
		Long:  "Update a STACKIT Kubernetes Engine (SCA) application.",
		Args:  args.SingleArg(applicationIDArg, nil),
		Example: examples.Build(
			// TODO: fix examples
			examples.NewExample(
				`Update the number of instances of a SCA application with ID "xxx" from an environment with ID "yyy"`,
				"$ stackit sca application update xxx --instances 2 --environment-id yyy"),
		),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			updateFlags := []string{
				imageFlag,
				publicFlag,
				externalPortFlag,
				cpuFlag,
				memoryFlag,
				instancesFlag,
				minInstancesFlag,
				maxInstancesFlag,
				scaleToZeroFlag,
				concurrencyFlag,
				rpsFlag,
				envVarsFlag,
				commandsFlag,
				argsFlag,
				stoppedFlag,
			}

			if !slices.ContainsFunc(updateFlags, cmd.Flags().Changed) {
				return fmt.Errorf("at least one field should be updated")
			}

			return nil
		},
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

			globalFlags := globalflags.Parse(params.Printer, cmd)
			if globalFlags.ProjectId == "" {
				return &errors.ProjectIdError{}
			}

			// Call API
			req, err := buildRequest(ctx, model, apiClient)
			if err != nil {
				return err
			}
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("update application: %w", err)
			}

			if !model.Async {
				err := spinner.Run(params.Printer, fmt.Sprintf("Updating application with id %q", resp.GetId()), func() error {
					_, err := wait.UpdateApplicationWaitHandler(ctx, apiClient.DefaultAPI, model.ProjectId, model.EnvironmentID, resp.GetId()).WaitWithContext(ctx)
					return err
				})
				if err != nil {
					return fmt.Errorf("wait for application update: %w", err)
				}
			}

			outputResult(params.Printer, model, resp)

			return nil
		},
	}

	configureFlags(cmd)
	return cmd
}

func containersFromInput(containers []sca.Container, model *inputModel) []sca.Container {
	if len(containers) == 0 {
		return []sca.Container{{
			Name:                 "container-1",
			Image:                model.Image,
			Cpu:                  model.CPU,
			Memory:               model.Memory,
			EnvironmentVariables: scautils.EnvironmentVariablesFromMap(model.EnvironmentVariables),
		}}
	}

	updateContainer := false
	if model.Image != "" {
		updateContainer = true
		containers[0].Image = model.Image
	}
	if model.CPU != nil {
		updateContainer = true
		containers[0].Cpu = model.CPU
	}
	if model.Memory != nil {
		updateContainer = true
		containers[0].Memory = model.Memory
	}

	if model.Commands != nil {
		containers[0].Command = model.Commands
	}

	if model.Args != nil {
		containers[0].Args = model.Args
	}

	for key, val := range model.EnvironmentVariables {
		found := false
		updateContainer = true
		for i, envVar := range containers[0].EnvironmentVariables {
			if envVar.Key == key {
				containers[0].EnvironmentVariables[i].Origin = sca.ENVVARTYPE_ENV_FROM_SOURCE_TYPE_MANUAL.Ptr()
				containers[0].EnvironmentVariables[i].Value = val
				found = true
			}
		}
		if !found {
			containers[0].EnvironmentVariables = append(containers[0].EnvironmentVariables, sca.EnvVar{
				Key:    key,
				Origin: sca.ENVVARTYPE_ENV_FROM_SOURCE_TYPE_MANUAL.Ptr(),
				Value:  val,
			})
		}

	}

	if updateContainer {
		return containers
	}

	return nil
}

func networkFromInput(model *inputModel) *sca.Network {
	if model.Public == nil && model.ContainerExternalPort == nil {
		return nil
	}

	network := &sca.Network{}

	if model.Public != nil {
		network.PublicIngress = *model.Public
	}

	if model.ContainerExternalPort != nil {
		network.Port = model.ContainerExternalPort
	}

	return network
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *sca.APIClient) (*sca.ApiUpdateApplicationRequest, error) {
	current, err := apiClient.DefaultAPI.GetApplication(ctx, model.ProjectId, model.EnvironmentID, model.ApplicationID).Execute()
	if err != nil {
		return nil, fmt.Errorf("failed to get current application instance: %w", err)
	}

	// scaling
	updatedScaling, err := buildScalingConfig(model, current)
	if err != nil {
		return nil, err
	}

	payload := sca.UpdateApplicationPayload{
		Scaling: updatedScaling,
		Stopped: model.Stopped,
	}

	payload.Containers = containersFromInput(current.Containers, model)
	payload.Network = networkFromInput(model)

	return new(apiClient.DefaultAPI.UpdateApplication(ctx, model.ProjectId, model.EnvironmentID, model.ApplicationID).
		UpdateApplicationPayload(payload)), nil
}

func buildScalingConfig(model *inputModel, current *sca.Application) (*sca.Scaling, error) {
	// input's scaling type is different from current
	if model.ScalingType != nil && *model.ScalingType != scautils.HumanReadableScalingType(current.GetScaling().Type) {
		switch *model.ScalingType {
		case scalingTypeManual:
			if model.Instances == nil {
				return nil, &errors.FlagValidationError{
					Flag:    instancesFlag,
					Details: "required flag if replacing scaling type to manual",
				}
			}
			return &sca.Scaling{
				Type: sca.SCALINGTYPE_SCALING_TYPE_MANUAL,
				ManualScaling: &sca.ManualScaling{
					Instances: *model.Instances,
				},
			}, nil
		case scalingTypeAuto:
			autoScaling, err := buildAutoscalingConfig(model, nil)
			if err != nil {
				return nil, err
			}
			return &sca.Scaling{
				Type:        sca.SCALINGTYPE_SCALING_TYPE_AUTO,
				AutoScaling: autoScaling,
			}, nil
		}
	}

	// input's scaling type is ni or equal to current

	switch current.Scaling.Type {
	case sca.SCALINGTYPE_SCALING_TYPE_MANUAL:
		instances := current.Scaling.ManualScaling.Instances
		if model.Instances != nil {
			instances = *model.Instances
		}
		return &sca.Scaling{
			Type: sca.SCALINGTYPE_SCALING_TYPE_MANUAL,
			ManualScaling: &sca.ManualScaling{
				Instances: instances,
			},
		}, nil
	case sca.SCALINGTYPE_SCALING_TYPE_AUTO:
		autoScaling, err := buildAutoscalingConfig(model, current.Scaling.AutoScaling)
		if err != nil {
			return nil, err
		}
		return &sca.Scaling{
			Type:        sca.SCALINGTYPE_SCALING_TYPE_AUTO,
			AutoScaling: autoScaling,
		}, nil
	default:
		return &current.Scaling, nil
	}
}

func buildAutoscalingConfig(model *inputModel, current *sca.AutoScaling) (*sca.AutoScaling, error) {
	if model.Instances != nil {
		return nil, &errors.FlagValidationError{
			Flag:    instancesFlag,
			Details: "Instances of an application with autoscale config can't be updated",
		}
	}

	autoScaling := &sca.AutoScaling{}
	if current != nil {
		autoScaling = current
	}

	if model.ScaleToZero != nil {
		autoScaling.AllowScaleToZero = model.ScaleToZero
	}

	if model.MinInstances != nil {
		autoScaling.MinInstances = *model.MinInstances
	}

	if model.MaxInstances != nil {
		autoScaling.MaxInstances = *model.MaxInstances
	}

	if autoScaling.MaxInstances == 0 {
		autoScaling.MaxInstances = autoScaling.MinInstances
	}

	if model.RPS != nil || model.Concurrency != nil {
		autoScaling.Rules = buildScalingRules(model, autoScaling.Rules)
	}

	if len(autoScaling.Rules) == 0 {
		return nil, &errors.OneOfFlagsIsMissing{
			MissingFlags: []string{rpsFlag, concurrencyFlag},
			SetFlag:      fmt.Sprintf("--%s=%s", scalingTypeFlag.Name(), scalingTypeAuto),
		}
	}

	return autoScaling, nil
}

func buildScalingRules(model *inputModel, current []sca.ScaleRule) []sca.ScaleRule {
	if model.Concurrency == nil && model.RPS == nil {
		return current
	}
	found := false
	for _, rule := range current {
		if rule.Type == sca.RULETYPE_RULE_TYPE_HTTP {
			found = true
			if model.Concurrency != nil {
				rule.HttpRule.Concurrency = model.Concurrency
			}
			if model.RPS != nil {
				rule.HttpRule.Rps = model.RPS
			}
		}
	}

	if found {
		return current
	}

	return append(current, sca.ScaleRule{
		Type: sca.RULETYPE_RULE_TYPE_HTTP,
		Name: "http-scaling-rule",
		HttpRule: &sca.HttpScaleRule{
			Concurrency: model.Concurrency,
			Rps:         model.RPS,
		},
	})
}

func outputResult(p *print.Printer, model *inputModel, application *sca.Application) {
	operationState := "Updated"
	if model.Async {
		operationState = "Triggered update of"
	}

	p.Outputf("%s application for environment %s. Application ID: %s\n", operationState, application.GetEnvironmentId(), application.GetId())
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), environmentIDFlag, "Environment ID (uses default environment if not set)")
	cmd.Flags().Bool(stoppedFlag, false, "Stopped")
	// container
	cmd.Flags().String(imageFlag, "", "Container image")
	cmd.Flags().Int32(cpuFlag, 0, "The dedicated virtual CPU processing power allocated per container instance")
	cmd.Flags().Int32(memoryFlag, 0, "The total amount of memory (RAM) allocated per container instance")
	cmd.Flags().StringToString(envVarsFlag, nil, "Environment variables to inject into the application")
	cmd.Flags().StringSlice(commandsFlag, nil, "Commands to execute in the application container")
	cmd.Flags().StringSlice(argsFlag, nil, "Arguments to pass to the application container command")
	// scaling
	scalingTypeFlag.Register(cmd.Flags())
	cmd.Flags().Int32(instancesFlag, 0, "The number of application instances (if manually scaled)")
	cmd.Flags().Int32(minInstancesFlag, 0, "The minimum number of application instances (if autoscaling is enabled)")
	cmd.Flags().Int32(maxInstancesFlag, 0, "The maximum number of application instances (if autoscaling is enabled)")
	cmd.Flags().Bool(scaleToZeroFlag, false, "Enable scale to zero (if autoscaling is enabled)")
	cmd.Flags().Int32(concurrencyFlag, 0, "Target number of in-flight requests to trigger autoscaling (if autoscaling is enabled)")
	cmd.Flags().Int32(rpsFlag, 0, "Target number of requests per second to trigger autoscaling (if autoscaling is enabled)")
	// networking
	cmd.Flags().Bool(publicFlag, false, "Exposes your application securely to the public internet via HTTPS endpoint")
	cmd.Flags().Int32(externalPortFlag, 0, "Container external exposed port")

	cmd.MarkFlagsMutuallyExclusive(instancesFlag, minInstancesFlag)
	cmd.MarkFlagsMutuallyExclusive(instancesFlag, maxInstancesFlag)
	cmd.MarkFlagsMutuallyExclusive(instancesFlag, scaleToZeroFlag)
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	applicationID := inputArgs[0]

	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	environmentID := flags.FlagToStringValue(p, cmd, environmentIDFlag)
	if environmentID == "" {
		environmentID = globalFlags.ProjectId
	}

	var envVars map[string]string
	if env := flags.FlagToStringToStringPointer(p, cmd, envVarsFlag); env != nil {
		envVars = *env
	}

	model := inputModel{
		GlobalFlagModel:       globalFlags,
		EnvironmentID:         environmentID,
		ApplicationID:         applicationID,
		Stopped:               flags.FlagToBoolPointer(p, cmd, stoppedFlag),
		Memory:                flags.FlagToInt32Pointer(p, cmd, memoryFlag),
		CPU:                   flags.FlagToInt32Pointer(p, cmd, cpuFlag),
		EnvironmentVariables:  envVars,
		Commands:              flags.FlagToStringSliceValue(p, cmd, commandsFlag),
		Args:                  flags.FlagToStringSliceValue(p, cmd, argsFlag),
		ScalingType:           scalingTypeFlag.Ptr(),
		Instances:             flags.FlagToInt32Pointer(p, cmd, instancesFlag),
		MinInstances:          flags.FlagToInt32Pointer(p, cmd, minInstancesFlag),
		MaxInstances:          flags.FlagToInt32Pointer(p, cmd, maxInstancesFlag),
		ScaleToZero:           flags.FlagToBoolPointer(p, cmd, scaleToZeroFlag),
		Concurrency:           flags.FlagToInt32Pointer(p, cmd, concurrencyFlag),
		RPS:                   flags.FlagToInt32Pointer(p, cmd, rpsFlag),
		Public:                flags.FlagToBoolPointer(p, cmd, publicFlag),
		ContainerExternalPort: flags.FlagToInt32Pointer(p, cmd, externalPortFlag),
	}

	p.DebugInputModel(model)
	return &model, nil
}
