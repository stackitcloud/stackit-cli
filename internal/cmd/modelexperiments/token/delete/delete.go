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
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/modelexperiments/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	modelexperiments "github.com/stackitcloud/stackit-sdk-go/services/modelexperiments/v1api"
)

const (
	tokenIDArg     = "TOKEN_ID"
	instanceIDFlag = "instance-id"
	regionFlag     = "region"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	TokenID    string
	InstanceID string
	Region     string
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:     fmt.Sprintf("delete %s", tokenIDArg),
		Short:   "Deletes an instance token",
		Long:    "Deletes an auth token from an AI Model Experiments instance.",
		Args:    args.SingleArg(tokenIDArg, utils.ValidateUUID),
		Example: examples.Build(examples.NewExample(`Delete an auth token with ID "xxx"`, `$ stackit ai-model-experiments token delete xxx --instance-id yyy`)),
		RunE: func(cmd *cobra.Command, inputArgs []string) error {
			model, err := parseInput(params.Printer, cmd, inputArgs)
			if err != nil {
				return err
			}
			apiClient, err := client.ConfigureClient(params.Printer, params.CliVersion)
			if err != nil {
				return err
			}
			if err := params.Printer.PromptForConfirmation(fmt.Sprintf("Are you sure you want to delete instance token %q?", model.TokenID)); err != nil {
				return err
			}
			if _, err := buildRequest(context.Background(), model, apiClient).Execute(); err != nil {
				return fmt.Errorf("delete instance token: %w", err)
			}
			params.Printer.Outputf("Deleted instance token %q\n", model.TokenID)
			return nil
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), instanceIDFlag, "ID of the instance")
	_ = flags.MarkFlagsRequired(cmd, instanceIDFlag, regionFlag)
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &cliErr.ProjectIdError{}
	}
	model := &inputModel{GlobalFlagModel: globalFlags, TokenID: inputArgs[0], InstanceID: flags.FlagToStringValue(p, cmd, instanceIDFlag), Region: flags.FlagToStringValue(p, cmd, regionFlag)}
	p.DebugInputModel(model)
	return model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *modelexperiments.APIClient) modelexperiments.ApiDeleteInstanceTokenRequest {
	return apiClient.DefaultAPI.DeleteInstanceToken(ctx, model.ProjectId, model.Region, model.TokenID, model.InstanceID)
}
