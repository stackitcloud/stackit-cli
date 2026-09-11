package delete

import (
	"context"
	"fmt"

	telemetryrouter "github.com/stackitcloud/stackit-sdk-go/services/telemetryrouter/v1api"
	"github.com/stackitcloud/stackit-sdk-go/services/telemetryrouter/v1api/wait"

	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	cliErr "github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
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
		Short: "Deletes a TelemetryRouter access token",
		Long:  "Deletes a TelemetryRouter access token.",
		Args:  args.SingleArg(accessTokenIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Delete access token with ID "xxx" for the TelemetryRouter instance "yyy"`,
				"$ stackit beta telemetryrouter access-token delete xxx --instance-id yyy",
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

			accessTokenLabel, err := telemetryrouterUtils.GetAccessTokenName(ctx, apiClient.DefaultAPI, model.ProjectId, model.Region, model.InstanceId, model.AccessTokenId)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get access token: %v", err)
			}
			if accessTokenLabel == "" {
				accessTokenLabel = model.AccessTokenId
			}

			prompt := fmt.Sprintf("Are you sure you want to delete access token %q? (This cannot be undone)", accessTokenLabel)
			err = params.Printer.PromptForConfirmation(prompt)
			if err != nil {
				return err
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			err = req.Execute()
			if err != nil {
				return fmt.Errorf("delete TelemetryRouter access token: %w", err)
			}

			// Wait for async operation, if async mode not enabled
			if !model.Async {
				err := spinner.Run(params.Printer, "Deleting access token", func() error {
					_, err = wait.DeleteAccessTokenWaitHandler(ctx, apiClient.DefaultAPI, model.ProjectId, model.Region, model.InstanceId, model.AccessTokenId).WaitWithContext(ctx)
					return err
				})
				if err != nil {
					return fmt.Errorf("wait for TelemetryRouter access token deletion: %w", err)
				}
			}

			operationState := "Deleted"
			if model.Async {
				operationState = "Triggered deletion of"
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
		InstanceId:      flags.FlagToStringValue(p, cmd, instanceIdFlag),
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *telemetryrouter.APIClient) telemetryrouter.ApiDeleteAccessTokenRequest {
	return apiClient.DefaultAPI.DeleteAccessToken(ctx, model.ProjectId, model.Region, model.InstanceId, model.AccessTokenId)
}
