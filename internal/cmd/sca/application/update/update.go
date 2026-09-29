package update

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
	"github.com/stackitcloud/stackit-cli/internal/pkg/spinner"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
)

const (
	applicationIDArg  = "APPLICATION_ID"
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
	stoppedFlag       = "stopped"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	EnvironmentID         string
	ApplicationID         string
	Image                 string
	ScalingType           string
	Public                *bool
	ContainerExternalPort int32
	CPU                   int32
	Memory                int32
	Instances             int32
	MinInstances          int32
	MaxInstances          int32
	ScaleToZero           bool
	Concurrency           int32
	RPS                   int32
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
				`Update a SCA application with ID "xxx" from an environment with ID "yyy"`,
				"$ stackit sca application describe xxx --environment-id yyy"),
			examples.NewExample(
				`Get details of all SCA application with ID "xxx" from an environment with ID "yyy" in JSON format`,
				"$ stackit sca application describe xxx --environment-id yyy --output-format json"),
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

func containersFromInput(model *inputModel) []sca.Container {
	container := sca.Container{}
	updateContainer := false
	if model.Image != "" {
		updateContainer = true
		container.Image = model.Image
	}
	if model.CPU != 0 {
		updateContainer = true
		container.Cpu = &model.CPU
	}
	if model.Memory != 0 {
		updateContainer = true
		container.Memory = &model.Memory
	}

	if updateContainer {
		return []sca.Container{container}
	}

	return nil
}

func networkFromInput(model *inputModel) *sca.Network {
	if model.Public == nil || model.ContainerExternalPort == 0 {
		return nil
	}

	return &sca.Network{
		PublicIngress: *model.Public,
		Port:          &model.ContainerExternalPort,
	}
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *sca.APIClient) sca.ApiUpdateApplicationRequest {
	payload := sca.UpdateApplicationPayload{
		Stopped: model.Stopped,
	}
	payload.Containers = containersFromInput(model)
	payload.Network = networkFromInput(model)

	return apiClient.DefaultAPI.UpdateApplication(ctx, model.ProjectId, model.EnvironmentID, model.ApplicationID).
		UpdateApplicationPayload(payload)
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

	model := inputModel{
		GlobalFlagModel: globalFlags,
		EnvironmentID:   environmentID,
		ApplicationID:   applicationID,
		Stopped:         flags.FlagToBoolPointer(p, cmd, stoppedFlag),
	}

	p.DebugInputModel(model)
	return &model, nil
}
