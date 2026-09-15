package update

import (
	"context"
	"fmt"

	telemetryrouter "github.com/stackitcloud/stackit-sdk-go/services/telemetryrouter/v1api"
	"github.com/stackitcloud/stackit-sdk-go/services/telemetryrouter/v1api/wait"

	"github.com/stackitcloud/stackit-cli/internal/pkg/types"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	cliErr "github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/telemetryrouter/client"
	telemetryrouterUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/telemetryrouter/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/spinner"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

const (
	instanceIdFlag   = "instance-id"
	displayNameFlag  = "display-name"
	descriptionFlag  = "description"
	accessTokenIdArg = "ACCESS_TOKEN_ID"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	InstanceId    string
	AccessTokenId string
	DisplayName   *string
	Description   *string
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("update %s", accessTokenIdArg),
		Short: "Updates a TelemetryRouter access token",
		Long:  "Updates a TelemetryRouter access token.",
		Args:  args.SingleArg(accessTokenIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Update the display name of the access token with ID "xxx"`,
				`$ stackit beta telemetryrouter access-token update xxx --instance-id yyy --display-name access-token-1`,
			),
			examples.NewExample(
				`Update the description of the access token with ID "xxx"`,
				`$ stackit beta telemetryrouter access-token update xxx --instance-id yyy --description "Access token for service"`,
			),
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

			instanceLabel, err := telemetryrouterUtils.GetInstanceName(ctx, apiClient.DefaultAPI, model.ProjectId, model.Region, model.InstanceId)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get instance name: %v", err)
			}
			if instanceLabel == "" {
				instanceLabel = model.InstanceId
			}

			accessTokenLabel, err := telemetryrouterUtils.GetAccessTokenName(ctx, apiClient.DefaultAPI, model.ProjectId, model.Region, model.InstanceId, model.AccessTokenId)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get access token: %v", err)
			}
			if accessTokenLabel == "" {
				accessTokenLabel = model.AccessTokenId
			}

			prompt := fmt.Sprintf("Are you sure you want to update access token %q for TelemetryRouter instance %q?", accessTokenLabel, instanceLabel)
			err = params.Printer.PromptForConfirmation(prompt)
			if err != nil {
				return err
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("update TelemetryRouter access token: %w", err)
			}
			if resp == nil {
				return fmt.Errorf("update TelemetryRouter access token: empty response from API")
			}

			// Wait for async operation, if async mode not enabled
			if !model.Async {
				err := spinner.Run(params.Printer, "Updating access token", func() error {
					_, err = wait.UpdateAccessTokenWaitHandler(ctx, apiClient.DefaultAPI, model.ProjectId, model.Region, model.InstanceId, model.AccessTokenId).WaitWithContext(ctx)
					return err
				})
				if err != nil {
					return fmt.Errorf("wait for TelemetryRouter access token update: %w", err)
				}
			}

			operationState := "Updated"
			if model.Async {
				operationState = "Triggered update of"
			}
			params.Printer.Outputf("%s access token %q\n", operationState, accessTokenLabel)
			return nil
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), instanceIdFlag, "ID of the TelemetryRouter instance")
	cmd.Flags().String(displayNameFlag, "", "Display name for the access token")
	cmd.Flags().String(descriptionFlag, "", "Description of the access token")

	err := flags.MarkFlagsRequired(cmd, instanceIdFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	accessTokenId := inputArgs[0]

	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &cliErr.ProjectIdError{}
	}

	displayName := flags.FlagToStringPointer(p, cmd, displayNameFlag)
	description := flags.FlagToStringPointer(p, cmd, descriptionFlag)

	if displayName == nil && description == nil {
		return nil, &cliErr.EmptyUpdateError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		InstanceId:      flags.FlagToStringValue(p, cmd, instanceIdFlag),
		AccessTokenId:   accessTokenId,
		DisplayName:     displayName,
		Description:     description,
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *telemetryrouter.APIClient) telemetryrouter.ApiUpdateAccessTokenRequest {
	req := apiClient.DefaultAPI.UpdateAccessToken(ctx, model.ProjectId, model.Region, model.InstanceId, model.AccessTokenId)

	payload := telemetryrouter.UpdateAccessTokenPayload{}
	// Only set fields that were actually provided via flags to prevent that it is cleared on API side.
	if model.DisplayName != nil {
		payload.DisplayName = *telemetryrouter.NewNullableString(model.DisplayName)
	}
	if model.Description != nil {
		payload.Description = *telemetryrouter.NewNullableString(model.Description)
	}

	return req.UpdateAccessTokenPayload(payload)
}
