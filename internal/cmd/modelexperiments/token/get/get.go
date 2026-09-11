package get

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
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	TokenID    string
	InstanceID string
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:     fmt.Sprintf("get %s", tokenIDArg),
		Short:   "Gets an instance token",
		Long:    "Gets an auth token for an AI Model Experiments instance.",
		Args:    args.SingleArg(tokenIDArg, utils.ValidateUUID),
		Example: examples.Build(examples.NewExample(`Get an auth token with ID "xxx"`, `$ stackit ai-model-experiments token get xxx --instance-id yyy`)),
		RunE: func(cmd *cobra.Command, inputArgs []string) error {
			model, err := parseInput(params.Printer, cmd, inputArgs)
			if err != nil {
				return err
			}
			apiClient, err := client.ConfigureClient(params.Printer, params.CliVersion)
			if err != nil {
				return err
			}
			resp, err := buildRequest(context.Background(), model, apiClient).Execute()
			if err != nil {
				return fmt.Errorf("get instance token: %w", err)
			}
			return outputResult(params.Printer, model.OutputFormat, resp)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), instanceIDFlag, "ID of the instance")
	_ = flags.MarkFlagsRequired(cmd, instanceIDFlag)
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &cliErr.ProjectIdError{}
	}
	model := &inputModel{GlobalFlagModel: globalFlags, TokenID: inputArgs[0], InstanceID: flags.FlagToStringValue(p, cmd, instanceIDFlag)}
	p.DebugInputModel(model)
	return model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *modelexperiments.APIClient) modelexperiments.ApiGetInstanceTokenRequest {
	return apiClient.DefaultAPI.GetInstanceToken(ctx, model.ProjectId, model.Region, model.TokenID, model.InstanceID)
}

func outputResult(p *print.Printer, outputFormat string, resp *modelexperiments.GetInstanceTokenResponse) error {
	if resp == nil || resp.Token.Name == "" {
		return fmt.Errorf("response token is nil")
	}
	return p.OutputResult(outputFormat, resp.Token, func() error {
		p.Outputf("Instance token %q (ID: %s)\n", resp.Token.Name, resp.Token.Id)
		return nil
	})
}
