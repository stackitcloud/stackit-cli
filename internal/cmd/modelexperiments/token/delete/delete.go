package delete

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	cliErr "github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/projectname"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/modelexperiments/client"
	meUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/modelexperiments/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/spinner"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	modelexperiments "github.com/stackitcloud/stackit-sdk-go/services/modelexperiments/v1api"
	"github.com/stackitcloud/stackit-sdk-go/services/modelexperiments/v1api/wait"
)

const (
	tokenIDArg     = "TOKEN_ID"
	instanceIDFlag = "instance-id"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	TokenID    string
	InstanceID string
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("delete %s", tokenIDArg),
		Short: "Deletes an AI Model Experiments instance token",
		Long:  "Deletes an auth token from an AI Model Experiments instance.",
		Args:  args.SingleArg(tokenIDArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Delete an auth token with ID "xxx"`,
				`$ stackit model-experiments token delete xxx --instance-id yyy`,
			),
		),
		RunE: func(cmd *cobra.Command, inputArgs []string) error {
			ctx := context.Background()

			model, err := parseInput(params.Printer, cmd, inputArgs)
			if err != nil {
				return err
			}

			apiClient, err := client.ConfigureClient(params.Printer, params.CliVersion)
			if err != nil {
				return err
			}

			projectLabel, err := projectname.GetProjectName(ctx, params.Printer, params.CliVersion, cmd)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get project name: %v", err)
				projectLabel = model.ProjectId
			} else if projectLabel == "" {
				projectLabel = model.ProjectId
			}

			instanceLabel, err := meUtils.GetInstanceName(ctx, apiClient.DefaultAPI, model.ProjectId, model.Region, model.InstanceID)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get instance name: %v", err)
				instanceLabel = model.InstanceID
			}

			tokenLabel, err := meUtils.GetTokenName(ctx, apiClient.DefaultAPI, model.ProjectId, model.Region, model.InstanceID, model.TokenID)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get token name: %v", err)
				tokenLabel = model.TokenID
			}

			prompt := fmt.Sprintf("Are you sure you want to delete token %q for instance %q in project %q?", tokenLabel, instanceLabel, projectLabel)
			if err := params.Printer.PromptForConfirmation(prompt); err != nil {
				return err
			}

			if _, err := buildRequest(ctx, model, apiClient).Execute(); err != nil {
				return fmt.Errorf("delete instance token: %w", err)
			}

			if !model.Async {
				err := spinner.Run(params.Printer, "Deleting instance token", func() error {
					_, err = wait.DeleteInstanceTokenWaitHandler(ctx, apiClient.DefaultAPI, model.Region, model.ProjectId, model.InstanceID, model.TokenID).WaitWithContext(ctx)
					return err
				})
				if err != nil {
					return fmt.Errorf("wait for instance token deletion: %w", err)
				}
			}

			operationState := "Deleted"
			if model.Async {
				operationState = "Triggered deletion of"
			}
			params.Printer.Outputf("%s token %q for instance %q in project %q\n", operationState, tokenLabel, instanceLabel, projectLabel)
			return nil
		},
	}

	configureFlags(cmd)

	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), instanceIDFlag, "ID of the instance")
	err := flags.MarkFlagsRequired(cmd, instanceIDFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &cliErr.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		TokenID:         inputArgs[0],
		InstanceID:      flags.FlagToStringValue(p, cmd, instanceIDFlag),
	}

	p.DebugInputModel(model)

	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *modelexperiments.APIClient) modelexperiments.ApiDeleteInstanceTokenRequest {
	return apiClient.DefaultAPI.DeleteInstanceToken(ctx, model.ProjectId, model.Region, model.TokenID, model.InstanceID)
}
