package list

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

	modelexperiments "github.com/stackitcloud/stackit-sdk-go/services/modelexperiments/v1api"
)

const (
	instanceIDFlag = "instance-id"
	regionFlag     = "region"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	InstanceID string
	Region     string
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{Use: "list", Short: "Lists instance tokens", Long: "Lists all auth tokens for an AI Model Experiments instance.", Args: args.NoArgs, Example: examples.Build(examples.NewExample("List all tokens for an instance", "$ stackit ai-model-experiments token list --instance-id xxx")), RunE: func(cmd *cobra.Command, _ []string) error {
		model, err := parseInput(params.Printer, cmd, nil)
		if err != nil {
			return err
		}
		apiClient, err := client.ConfigureClient(params.Printer, params.CliVersion)
		if err != nil {
			return err
		}
		resp, err := buildRequest(context.Background(), model, apiClient).Execute()
		if err != nil {
			return fmt.Errorf("list instance tokens: %w", err)
		}
		return outputResult(params.Printer, model.OutputFormat, resp)
	}}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), instanceIDFlag, "ID of the instance")
	_ = flags.MarkFlagsRequired(cmd, instanceIDFlag, regionFlag)
}

func parseInput(p *print.Printer, cmd *cobra.Command, _ []string) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &cliErr.ProjectIdError{}
	}
	model := &inputModel{GlobalFlagModel: globalFlags, InstanceID: flags.FlagToStringValue(p, cmd, instanceIDFlag), Region: flags.FlagToStringValue(p, cmd, regionFlag)}
	p.DebugInputModel(model)
	return model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *modelexperiments.APIClient) modelexperiments.ApiListInstanceTokensRequest {
	return apiClient.DefaultAPI.ListInstanceTokens(ctx, model.ProjectId, model.Region, model.InstanceID)
}

func outputResult(p *print.Printer, outputFormat string, resp *modelexperiments.ListInstanceTokensResponse) error {
	if resp == nil {
		return fmt.Errorf("response is nil")
	}
	return p.OutputResult(outputFormat, resp.Tokens, func() error { return nil })
}
