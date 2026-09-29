package updatefrompayload

import (
	"context"
	"encoding/json"
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
	payloadFlag       = "payload"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	EnvironmentID string
	ApplicationID string
	Payload       *sca.UpdateApplicationPayload
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-from-payload",
		Short: "Update a SCA application from payload",
		Long:  "Update a STACKIT Kubernetes Engine (SCA) application from payload.",
		Args:  args.SingleArg(applicationIDArg, nil),
		Example: examples.Build(
			examples.NewExample(
				`Update a SCA application using an API payload sourced from the file "./payload.json"`,
				"$ stackit sca application update-from-payload my-application-id --payload @./payload.json"),
			examples.NewExample(
				`Update a SCA application using an API payload provided as a JSON string`,
				`$ stackit sca application update-from-payload my-application-id --payload "{...}"`),
			examples.NewExample(
				`Generate a payload with the current values of an application, and adapt it with custom values for the different configuration options`,
				`$ stackit sca application generate-payload --application-id application-id > ./payload.json`,
				`<Modify payload in file>`,
				`$ stackit sca application update-from-payload application-id --payload @./payload.json`),
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

func buildRequest(ctx context.Context, model *inputModel, apiClient *sca.APIClient) sca.ApiUpdateApplicationRequest {
	return apiClient.DefaultAPI.UpdateApplication(ctx, model.ProjectId, model.EnvironmentID, model.ApplicationID).
		UpdateApplicationPayload(*model.Payload)
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
	cmd.Flags().Var(flags.ReadFromFileFlag(), payloadFlag, `Request payload (JSON). Can be a string or a file path, if prefixed with "@" (example: @./payload.json). If unset, will use a default payload (you can check it by running "stackit sca application generate-payload")`)
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

	payloadValue := flags.FlagToStringPointer(p, cmd, payloadFlag)
	var payload *sca.UpdateApplicationPayload
	if payloadValue != nil {
		payload = &sca.UpdateApplicationPayload{}
		err := json.Unmarshal([]byte(*payloadValue), payload)
		if err != nil {
			return nil, fmt.Errorf("enconde payload: %w", err)
		}
	}

	payload.AdditionalProperties = nil

	fmt.Printf("%+v\n", payload)

	model := inputModel{
		GlobalFlagModel: globalFlags,
		EnvironmentID:   environmentID,
		ApplicationID:   applicationID,
		Payload:         payload,
	}

	p.DebugInputModel(model)
	return &model, nil
}
