package create

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/sca/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/spinner"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	sca "github.com/stackitcloud/stackit-sdk-go/services/sca/v1alphaapi"
	"github.com/stackitcloud/stackit-sdk-go/services/sca/v1alphaapi/wait"
)

const (
	environmentIDFlag = "environment-id"
	nameFlag          = "name"
	imageFlag         = "image"
	publicFlag        = "public"
	externalPortFlag  = "external-port"
	cpuFlag           = "cpu"
	memoryFlag        = "memory"
	instancesFlag     = "instances"
	minInstancesFlag  = "min-instances"
	maxInstancesFlag  = "max-instances"
	scaleToZeroFlag   = "scale-to-zero"
)

var scalingTypeFlag = flags.StringEnumFlag(
	"scaling-type",
	[]string{"manual", "auto"},
	"Scaling type,",
	flags.StringEnumDefaultValue("manual"),
)

const (
	defaultPublic      = true
	defaultPort        = 8080
	defaultCPU         = 1000
	defaultMemory      = 1024
	defaultScalingType = "manual"
	defaultInstances   = 1
	// defaultMinInstances = 1
	// defaultMaxInstances = 1
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	EnvironmentID         string
	Name                  string
	Image                 string
	ScalingType           string
	Public                bool
	ContainerExternalPort int32
	CPU                   int32
	Memory                int32
	Instances             int32
	MinInstances          int32
	MaxInstances          int32
	ScaleToZero           bool
	Concurrency           int32
	RPS                   int32
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a SCA application",
		Long:  "Create a STACKIT Kubernetes Engine (SCA) application.",
		Args:  args.NoArgs,
		Example: examples.Build(
			// TODO: fix examples
			examples.NewExample(
				`Create a SCA application with ID "xxx" from an environment with ID "yyy"`,
				"$ stackit sca application describe xxx --environment-id yyy"),
			examples.NewExample(
				`Get details of all SCA application with ID "xxx" from an environment with ID "yyy" in JSON format`,
				"$ stackit sca application describe xxx --environment-id yyy --output-format json"),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
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

			outputResult(params.Printer, model, resp)

			return nil
		},
	}

	configureFlags(cmd)
	return cmd
}

func buildScalingRules(model *inputModel) []sca.ScaleRule {
	scaleRules := []sca.ScaleRule{}
	if model.Concurrency != 0 || model.RPS != 0 {
		scaleRules = append(scaleRules, sca.ScaleRule{
			Type: sca.RULETYPE_RULE_TYPE_HTTP,
			HttpRule: &sca.HttpScaleRule{
				Concurrency: &model.Concurrency,
				Rps:         &model.RPS,
			},
		})
	}

	return scaleRules
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *sca.APIClient) sca.ApiCreateApplicationRequest {
	scaling := sca.Scaling{
		Type: sca.SCALINGTYPE_SCALING_TYPE_MANUAL,
		ManualScaling: &sca.ManualScaling{
			Instances: model.Instances,
		},
	}
	if model.ScalingType == "auto" {
		rules := []sca.ScaleRule{
			{
				HttpRule: &sca.HttpScaleRule{
					Concurrency: model.Concurrency,
					Rps:         model.RPS,
				},
			},
		}

		scaling = sca.Scaling{
			Type: sca.SCALINGTYPE_SCALING_TYPE_AUTO,
			AutoScaling: &sca.AutoScaling{
				MinInstances:     model.MinInstances,
				MaxInstances:     model.MaxInstances,
				AllowScaleToZero: &model.ScaleToZero,
				Rules:            rules,
			},
		}

	}

	network := sca.Network{
		PublicIngress: model.Public,
		Port:          &model.ContainerExternalPort,
	}

	container := sca.Container{
		Name:   "container-1",
		Image:  model.Image,
		Cpu:    &model.CPU,
		Memory: &model.Memory,
	}

	payload := sca.CreateApplicationPayload{
		DisplayName: model.Name,
		Scaling:     scaling,
		Network:     network,
		Containers:  []sca.Container{container},
	}

	return apiClient.DefaultAPI.CreateApplication(ctx, model.ProjectId, model.EnvironmentID).
		CreateApplicationPayload(payload)
}

func outputResult(p *print.Printer, model *inputModel, application *sca.Application) {
	operationState := "Created"
	if model.Async {
		operationState = "Triggered creation of"
	}

	p.Outputf("%s application for environment %s. Application ID: %s\n", operationState, application.GetEnvironmentId(), application.GetId())
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), environmentIDFlag, "Environment ID (uses default environment if not set)")
	cmd.Flags().String(nameFlag, "", "Application display name")
	cmd.Flags().String(imageFlag, "", "Container image")
	cmd.Flags().Bool(publicFlag, defaultPublic, "Exposes your application securely to the public internet via HTTPS endpoint")
	cmd.Flags().Int32(externalPortFlag, defaultPort, "Container external exposed port")
	cmd.Flags().Int32(cpuFlag, defaultCPU, "The dedicated virtual CPU processing power allocated per container instance")
	cmd.Flags().Int32(memoryFlag, defaultMemory, "The total amount of memory (RAM) allocated per container instance")
	scalingTypeFlag.Register(cmd.Flags())
	// cmd.Flags().String(scalingTypeFlag, defaultScalingType, "")
	cmd.Flags().Int32(instancesFlag, defaultInstances, "")
	cmd.Flags().Int32(minInstancesFlag, 0, "")
	cmd.Flags().Int32(maxInstancesFlag, 0, "")
	cmd.Flags().Bool(maxInstancesFlag, false, "")

	cobra.CheckErr(flags.MarkFlagsRequired(cmd, nameFlag))
	cobra.CheckErr(flags.MarkFlagsRequired(cmd, imageFlag))
}

func parseInput(p *print.Printer, cmd *cobra.Command) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	extertalPort := flags.FlagWithDefaultToInt32Value(p, cmd, externalPortFlag)
	if extertalPort <= 1024 || extertalPort > 65535 {
		return nil, &errors.FlagValidationError{
			Flag:    externalPortFlag,
			Details: "must be a valid non-privileged port (from 1025 to 65535)",
		}
	}

	environmentID := flags.FlagToStringValue(p, cmd, environmentIDFlag)
	if environmentID == "" {
		environmentID = globalFlags.ProjectId
	}

	model := inputModel{
		GlobalFlagModel:       globalFlags,
		EnvironmentID:         environmentID,
		Name:                  flags.FlagToStringValue(p, cmd, nameFlag),
		Image:                 flags.FlagToStringValue(p, cmd, imageFlag),
		Public:                flags.FlagToBoolValue(p, cmd, publicFlag),
		ContainerExternalPort: extertalPort,
		CPU:                   flags.FlagWithDefaultToInt32Value(p, cmd, cpuFlag),
		Memory:                flags.FlagWithDefaultToInt32Value(p, cmd, memoryFlag),
		Instances:             flags.FlagWithDefaultToInt32Value(p, cmd, instancesFlag),
		ScalingType:           flags.FlagToStringValue(p, cmd, scalingTypeFlag),
		MinInstances:          flags.FlagWithDefaultToInt32Value(p, cmd, minInstancesFlag),
		MaxInstances:          flags.FlagWithDefaultToInt32Value(p, cmd, maxInstancesFlag),
		ScaleToZero:           flags.FlagToBoolValue(p, cmd, scaleToZeroFlag),
	}

	p.DebugInputModel(model)
	return &model, nil
}
