package createfrompayload

import (
	"context"
	"encoding/json"
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
	payloadFlag       = "payload"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	EnvironmentID string
	Payload       *sca.CreateApplicationPayload
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-from-payload",
		Short: "Create a SCA application from payload",
		Long:  "Create a STACKIT Kubernetes Engine (SCA) application from payload.",
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

func buildRequest(ctx context.Context, model *inputModel, apiClient *sca.APIClient) sca.ApiCreateApplicationRequest {
	return apiClient.DefaultAPI.CreateApplication(ctx, model.ProjectId, model.EnvironmentID).
		CreateApplicationPayload(*model.Payload)
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
	cmd.Flags().Var(flags.ReadFromFileFlag(), payloadFlag, `Request payload (JSON). Can be a string or a file path, if prefixed with "@" (example: @./payload.json). If unset, will use a default payload (you can check it by running "stackit sca application generate-payload")`)

	cobra.CheckErr(flags.MarkFlagsRequired(cmd, nameFlag))
}

func parseInput(p *print.Printer, cmd *cobra.Command) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	environmentID := flags.FlagToStringValue(p, cmd, environmentIDFlag)
	if environmentID == "" {
		environmentID = globalFlags.ProjectId
	}

	payloadValue := flags.FlagToStringPointer(p, cmd, payloadFlag)
	var payload *sca.CreateApplicationPayload
	if payloadValue != nil {
		payload = &sca.CreateApplicationPayload{}
		err := json.Unmarshal([]byte(*payloadValue), payload)
		if err != nil {
			return nil, fmt.Errorf("enconde payload: %w", err)
		}
	}

	payload.DisplayName = flags.FlagToStringValue(p, cmd, nameFlag)

	model := inputModel{
		GlobalFlagModel: globalFlags,
		EnvironmentID:   environmentID,
		Payload:         payload,
	}

	p.DebugInputModel(model)
	return &model, nil
}
