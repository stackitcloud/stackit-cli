package create

import (
	"context"
	"fmt"

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
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/sca/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/spinner"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
)

const (
	environmentIDFlag = "environment-id"
	nameFlag          = "name"
	// Container config
	containerNameFlag = "container-name"
	imageFlag         = "image"
	cpuFlag           = "cpu"
	memoryFlag        = "memory"
	envVarsFlag       = "environment-vars"
	commandsFlag      = "commands"
	argsFlag          = "args"
	// Scaling
	instancesFlag    = "instances"
	minInstancesFlag = "min-instances"
	maxInstancesFlag = "max-instances"
	scaleToZeroFlag  = "scale-to-zero"
	rpsFlag          = "http-rule-rps"
	concurrencyFlag  = "http-rule-concurrency"
	// Networking
	publicFlag       = "public"
	externalPortFlag = "external-port"

	scalingTypeManual = "manual"
	scalingTypeAuto   = "auto"

	defaultPublic        = true
	defaultPort          = 8080
	defaultCPU           = 1000
	defaultMemory        = 1024
	defaultInstances     = 1
	defaultContainerName = "container-1"
)

var scalingTypeFlag = flags.StringEnumFlag(
	"scaling-type",
	[]string{scalingTypeManual, scalingTypeAuto},
	"Scaling type,",
	flags.StringEnumDefaultValue(scalingTypeManual),
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	EnvironmentID         string
	Name                  string
	ContainerName         string
	Image                 string
	CPU                   int32
	Memory                int32
	EnvironmentVars       map[string]string
	Commands              []string
	Args                  []string
	ScalingType           string
	Instances             int32
	MinInstances          int32
	MaxInstances          int32
	ScaleToZero           bool
	Concurrency           int32
	RPS                   int32
	Public                bool
	ContainerExternalPort int32
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a SCA application",
		Long:  "Create a STACKIT Container Applications (SCA) application.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Create a SCA application with name "application-name" and image "my-image" for an environment with ID "yyy"`,
				"$ stackit alpha sca application create --name application-name --image my-image --environment-id yyy"),
			examples.NewExample(
				`Create a SCA application with name "application-name" and image "my-image" with 2 instances`,
				"$ stackit alpha sca application create --name application-name --image my-image --instances 2"),
			examples.NewExample(
				`Create a SCA application with name "application-name" and image "my-image" exposing port 8888 of the container`,
				"$ stackit alpha sca application create --name application-name --image my-image --external-port 8888"),
			examples.NewExample(
				`Create a SCA application with name "application-name" and image "my-image" disabling public networking`,
				"$ stackit alpha sca application create --name application-name --image my-image --public=false"),
			examples.NewExample(
				`Create a SCA application with name "application-name" and image "my-image" and environment variables ENV1=value1 and ENV2=value2`,
				"$ stackit alpha sca application create --name application-name --image my-image --environment-vars ENV1=value1,ENV2=value2"),
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

			globalFlags := globalflags.Parse(params.Printer, cmd)
			if globalFlags.ProjectId == "" {
				return &errors.ProjectIdError{}
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("create application: %w", err)
			}

			if !model.Async {
				err := spinner.Run(params.Printer, fmt.Sprintf("Creating application with id %q", resp.GetId()), func() error {
					_, err := wait.CreateApplicationWaitHandler(ctx, apiClient.DefaultAPI, model.ProjectId, model.EnvironmentID, resp.GetId()).WaitWithContext(ctx)
					return err
				})
				if err != nil {
					return fmt.Errorf("wait for application creation: %w", err)
				}
			}

			return outputResult(params.Printer, model, resp)
		},
	}

	configureFlags(cmd)
	return cmd
}

func buildScalingConfig(model *inputModel) sca.Scaling {
	if model.ScalingType == scalingTypeAuto {
		return sca.Scaling{
			Type: sca.SCALINGTYPE_SCALING_TYPE_AUTO,
			AutoScaling: &sca.AutoScaling{
				MinInstances:     model.MinInstances,
				MaxInstances:     model.MaxInstances,
				AllowScaleToZero: &model.ScaleToZero,
				Rules:            buildScalingRules(model),
			},
		}
	}

	return sca.Scaling{
		Type: sca.SCALINGTYPE_SCALING_TYPE_MANUAL,
		ManualScaling: &sca.ManualScaling{
			Instances: model.Instances,
		},
	}
}

func buildScalingRules(model *inputModel) []sca.ScaleRule {
	var scaleRules []sca.ScaleRule
	if model.Concurrency != 0 || model.RPS != 0 {
		scaleRules = append(scaleRules, sca.ScaleRule{
			Type: sca.RULETYPE_RULE_TYPE_HTTP,
			Name: "http-scaling-rule",
			HttpRule: &sca.HttpScaleRule{
				Concurrency: &model.Concurrency,
				Rps:         &model.RPS,
			},
		})
	}

	return scaleRules
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *sca.APIClient) sca.ApiCreateApplicationRequest {
	network := sca.Network{
		PublicIngress: model.Public,
		Port:          &model.ContainerExternalPort,
	}

	envVars := utils.EnvironmentVariablesFromMap(model.EnvironmentVars)

	container := sca.Container{
		Name:                 defaultContainerName,
		Image:                model.Image,
		Cpu:                  &model.CPU,
		Memory:               &model.Memory,
		EnvironmentVariables: envVars,
		Command:              model.Commands,
		Args:                 model.Args,
	}

	payload := sca.CreateApplicationPayload{
		DisplayName: model.Name,
		Scaling:     buildScalingConfig(model),
		Network:     network,
		Containers:  []sca.Container{container},
	}

	return apiClient.DefaultAPI.CreateApplication(ctx, model.ProjectId, model.EnvironmentID).
		CreateApplicationPayload(payload)
}

func outputResult(p *print.Printer, model *inputModel, application *sca.Application) error {
	if application == nil {
		return fmt.Errorf("create application response is empty")
	}

	return p.OutputResult(model.OutputFormat, application, func() error {
		operationState := "Created"
		if model.Async {
			operationState = "Triggered creation of"
		}
		p.Outputf("%s application for environment %s. Application ID: %s\n", operationState, application.GetEnvironmentId(), application.GetId())
		return nil
	})
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), environmentIDFlag, "Environment ID (uses default environment if not set)")
	cmd.Flags().String(nameFlag, "", "Application display name")
	// container
	cmd.Flags().String(containerNameFlag, defaultContainerName, "Container name")
	cmd.Flags().String(imageFlag, "", "Container image")
	cmd.Flags().Int32(cpuFlag, defaultCPU, "The dedicated virtual CPU processing power allocated per container instance")
	cmd.Flags().Int32(memoryFlag, defaultMemory, "The total amount of memory (RAM) allocated per container instance")
	cmd.Flags().StringToString(envVarsFlag, nil, "Environment variables to inject into the application")
	cmd.Flags().StringSlice(commandsFlag, nil, "Commands to execute in the application container")
	cmd.Flags().StringSlice(argsFlag, nil, "Arguments to pass to the application container command")
	// scaling
	scalingTypeFlag.Register(cmd.Flags())
	cmd.Flags().Int32(instancesFlag, defaultInstances, "The number of application instances (if manually scaled)")
	cmd.Flags().Int32(minInstancesFlag, 1, "The minimum number of application instances (if autoscaling is enabled)")
	cmd.Flags().Int32(maxInstancesFlag, 0, "The maximum number of application instances (if autoscaling is enabled)")
	cmd.Flags().Int32(rpsFlag, 0, "Target number of requests per second to trigger autoscaling (if autoscaling is enabled)")
	cmd.Flags().Int32(concurrencyFlag, 0, "Target number of in-flight requests to trigger autoscaling (if autoscaling is enabled)")
	cmd.Flags().Bool(scaleToZeroFlag, false, "Enable scale to zero (if autoscaling is enabled)")
	// networking
	cmd.Flags().Bool(publicFlag, defaultPublic, "Exposes your application securely to the public internet via HTTPS endpoint")
	cmd.Flags().Int32(externalPortFlag, defaultPort, "Container external exposed port")

	cobra.CheckErr(flags.MarkFlagsRequired(cmd, nameFlag))
	cobra.CheckErr(flags.MarkFlagsRequired(cmd, imageFlag))
}

func parseInput(p *print.Printer, cmd *cobra.Command, _ []string) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	var scalingType string
	if scalingTypeFlagValue := scalingTypeFlag.Ptr(); scalingTypeFlagValue != nil && *scalingTypeFlagValue != "" {
		scalingType = *scalingTypeFlagValue
	}

	if err := validateScalingInput(cmd, scalingType); err != nil {
		return nil, err
	}

	instances := flags.FlagWithDefaultToInt32Value(p, cmd, instancesFlag)
	if instances < 0 || instances > 10 {
		return nil, &errors.FlagValidationError{
			Flag:    instancesFlag,
			Details: "must be an integer between 0 and 10",
		}
	}

	minInstances := flags.FlagWithDefaultToInt32Value(p, cmd, minInstancesFlag)
	if minInstances < 0 || minInstances > 10 {
		return nil, &errors.FlagValidationError{
			Flag:    minInstancesFlag,
			Details: "must be an integer between 0 and 10",
		}
	}

	maxInstances := minInstances
	if i := flags.FlagToInt32Pointer(p, cmd, maxInstancesFlag); i != nil {
		maxInstances = *i
	}
	if maxInstances < minInstances || maxInstances > 10 {
		return nil, &errors.FlagValidationError{
			Flag:    minInstancesFlag,
			Details: fmt.Sprintf("must be an integer between minInstances (%d) and 10", minInstances),
		}
	}

	extertalPort := flags.FlagWithDefaultToInt32Value(p, cmd, externalPortFlag)
	if extertalPort <= 1024 || extertalPort > 65535 {
		return nil, &errors.FlagValidationError{
			Flag:    externalPortFlag,
			Details: "must be a valid non-privileged port (from 1025 to 65535)",
		}
	}

	cpu := flags.FlagWithDefaultToInt32Value(p, cmd, cpuFlag)
	if cpu <= 0 || cpu%1000 != 0 {
		return nil, &errors.FlagValidationError{
			Flag:    cpuFlag,
			Details: "must be a valid value divisible by 1000",
		}
	}

	var envVars map[string]string
	if env := flags.FlagToStringToStringPointer(p, cmd, envVarsFlag); env != nil {
		envVars = *env
	}

	environmentID := flags.FlagToStringValue(p, cmd, environmentIDFlag)
	if environmentID == "" {
		environmentID = globalFlags.ProjectId
	}

	model := inputModel{
		GlobalFlagModel:       globalFlags,
		EnvironmentID:         environmentID,
		Name:                  flags.FlagToStringValue(p, cmd, nameFlag),
		ContainerName:         flags.FlagWithDefaultToStringValue(p, cmd, containerNameFlag),
		Image:                 flags.FlagToStringValue(p, cmd, imageFlag),
		CPU:                   cpu,
		Memory:                flags.FlagWithDefaultToInt32Value(p, cmd, memoryFlag),
		EnvironmentVars:       envVars,
		Commands:              flags.FlagToStringSliceValue(p, cmd, commandsFlag),
		Args:                  flags.FlagToStringSliceValue(p, cmd, argsFlag),
		ScalingType:           scalingType,
		Instances:             instances,
		MinInstances:          minInstances,
		MaxInstances:          maxInstances,
		ScaleToZero:           flags.FlagToBoolValue(p, cmd, scaleToZeroFlag),
		RPS:                   flags.FlagWithDefaultToInt32Value(p, cmd, rpsFlag),
		Concurrency:           flags.FlagWithDefaultToInt32Value(p, cmd, concurrencyFlag),
		Public:                flags.FlagToBoolValue(p, cmd, publicFlag),
		ContainerExternalPort: extertalPort,
	}

	p.DebugInputModel(model)
	return &model, nil
}

func validateScalingInput(cmd *cobra.Command, scalingType string) error {
	switch scalingType {
	case scalingTypeManual:
		autoFlags := []string{minInstancesFlag, maxInstancesFlag, scaleToZeroFlag, concurrencyFlag, rpsFlag}
		for _, flagName := range autoFlags {
			if cmd.Flags().Changed(flagName) {
				return &errors.FlagValidationError{
					Flag:    flagName,
					Details: fmt.Sprintf("is only valid when --scaling-type is %q", scalingTypeAuto),
				}
			}
		}
	case scalingTypeAuto:
		if cmd.Flags().Changed(instancesFlag) {
			return &errors.FlagValidationError{
				Flag:    instancesFlag,
				Details: fmt.Sprintf("is only valid when --scaling-type is %q", scalingTypeManual),
			}
		}
		if !cmd.Flags().Changed(rpsFlag) && !cmd.Flags().Changed(concurrencyFlag) {
			return &errors.OneOfFlagsIsMissing{
				MissingFlags: []string{rpsFlag, concurrencyFlag},
				SetFlag:      fmt.Sprintf("--%s=%s", scalingTypeFlag.Name(), scalingTypeAuto),
			}
		}
	default:
		return &errors.FlagValidationError{
			Flag:    scalingType,
			Details: "invalid scaling type",
		}
	}

	return nil
}
