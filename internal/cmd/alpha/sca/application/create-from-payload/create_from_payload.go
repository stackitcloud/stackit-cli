package createfrompayload

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
		Long: fmt.Sprintf("%s\n%s\n%s",
			"Create a STACKIT Kubernetes Engine (SCA) application from payload.",
			`The payload can be provided as a JSON string or a file path prefixed with "@".`,
			"See https://docs.api.stackit.cloud/documentation/sca/version/v1alpha#tag/Applications/operation/Applications_CreateApplication for information regarding the payload structure.",
		),
		Args: args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Create a SCA application using an API payload sourced from the file "./payload.json"`,
				"$ stackit alpha sca cluster create-from-payload --name application-name --payload @./payload.json"),
			examples.NewExample(
				`Create a SCA application using an API payload provided as a JSON string`,
				`$ stackit alpha sca cluster create-from-payload --name application-name --payload "{...}"`),
			examples.NewExample(
				`Generate a payload with default values, and adapt it with custom values for the different configuration options`,
				`$ stackit alpha sca application generate-payload --file-path ./payload.json`,
				`<Modify payload in file, if needed>`,
				`$ stackit alpha sca application create-from-payload --name application-name --payload @./payload.json`),
		),
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := context.Background()

			model, err := parseInput(params.Printer, cmd, nil)
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

func parseInput(p *print.Printer, cmd *cobra.Command, _ []string) (*inputModel, error) {
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
