package delete

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-sdk-go/services/logs"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/logs/client"
	logUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/logs/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

const (
	instanceIdFlag   = "instance-id"
	accessTokenIdArg = "ACCESS_TOKEN_ID"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	InstanceId    string
	AccessTokenId string
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("delete %s", accessTokenIdArg),
		Short: "Deletes a Logs access token",
		Long:  "Deletes a Logs access token.",
		Args:  args.SingleArg(accessTokenIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Delete access token with ID "xxx" in instance "yyy"`,
				"$ stackit logs access-token delete xxx --instance-id yyy",
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
			accessTokenLabel, err := logUtils.GetAccessTokenName(ctx, apiClient, model.ProjectId, model.Region, model.InstanceId, model.AccessTokenId)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get access token: %v", err)
			}
			if accessTokenLabel == "" {
				accessTokenLabel = model.AccessTokenId
			}

			prompt := fmt.Sprintf("Are you sure you want to delete access token %q?", accessTokenLabel)
			err = params.Printer.PromptForConfirmation(prompt)
			if err != nil {
				return err
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			err = req.Execute()
			if err != nil {
				return fmt.Errorf("delete access token: %w", err)
			}

			params.Printer.Outputf("Deleted access token %q\n", accessTokenLabel)
			return nil
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), instanceIdFlag, "ID of the Logs instance")

	err := flags.MarkFlagsRequired(cmd, instanceIdFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	accessTokenId := inputArgs[0]

	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		AccessTokenId:   accessTokenId,
		InstanceId:      flags.FlagToStringValue(p, cmd, instanceIdFlag),
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *logs.APIClient) logs.ApiDeleteAccessTokenRequest {
	return apiClient.DeleteAccessToken(ctx, model.ProjectId, model.Region, model.InstanceId, model.AccessTokenId)
}
