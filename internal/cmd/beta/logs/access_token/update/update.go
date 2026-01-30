package update

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-sdk-go/services/logs"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	cliErr "github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/logs/client"
	logUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/logs/utils"
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
	Description   *string
	DisplayName   *string
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("update %s", accessTokenIdArg),
		Short: "Updates a access token",
		Long:  "Updates a access token.",
		Args:  args.SingleArg(accessTokenIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Update access token with ID "xxx" with new name "access-token-1"`,
				`$ stackit logs access-token update xxx --display-name access-token-1`,
			),
			examples.NewExample(
				`Update access token with ID "xxx" with new description "Access token for Service XY"`,
				`$ stackit logs access-token update xxx --description "Access token for Service XY"`,
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

			// Get the display name for confirmation
			instanceLabel, err := logUtils.GetInstanceName(ctx, apiClient, model.ProjectId, model.Region, model.InstanceId)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get log instance: %v", err)
			}
			if instanceLabel == "" {
				instanceLabel = model.InstanceId
			}

			// Get the display name for confirmation
			accessTokenLabel, err := logUtils.GetAccessTokenName(ctx, apiClient, model.ProjectId, model.Region, model.InstanceId, model.AccessTokenId)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get access token: %v", err)
			}
			if accessTokenLabel == "" {
				accessTokenLabel = model.AccessTokenId
			}

			prompt := fmt.Sprintf("Are you sure you want to update access token %q for instance %q?", accessTokenLabel, instanceLabel)
			err = params.Printer.PromptForConfirmation(prompt)
			if err != nil {
				return err
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			err = req.Execute()
			if err != nil {
				return fmt.Errorf("update access token: %w", err)
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
	cmd.Flags().Var(flags.UUIDFlag(), instanceIdFlag, "ID of the logs instance")
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

	model := inputModel{
		GlobalFlagModel: globalFlags,
		AccessTokenId:   accessTokenId,
		DisplayName:     flags.FlagToStringPointer(p, cmd, displayNameFlag),
		InstanceId:      flags.FlagToStringValue(p, cmd, instanceIdFlag),
		Description:     flags.FlagToStringPointer(p, cmd, descriptionFlag),
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *logs.APIClient) logs.ApiUpdateAccessTokenRequest {
	req := apiClient.UpdateAccessToken(ctx, model.ProjectId, model.Region, model.InstanceId, model.AccessTokenId)

	payload := logs.UpdateAccessTokenPayload{
		DisplayName: model.DisplayName,
		Description: model.Description,
	}

	return req.UpdateAccessTokenPayload(payload)
}
