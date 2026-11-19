package describe

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
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
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
		Use:   fmt.Sprintf("describe %s", tokenIDArg),
		Short: "Shows details of an AI Model Experiments instance token",
		Long:  "Shows details of an auth token for an AI Model Experiments instance.",
		Args:  args.SingleArg(tokenIDArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Show details of an auth token with ID "xxx"`,
				`$ stackit model-experiments token describe xxx --instance-id yyy`,
			),
			examples.NewExample(
				`Show details of an auth token with ID "xxx" in JSON format`,
				`$ stackit model-experiments token describe xxx --instance-id yyy --output-format json`,
			),
		),
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
				return fmt.Errorf("describe instance token: %w", err)
			}
			return outputResult(params.Printer, model, resp)
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
	model := &inputModel{
		GlobalFlagModel: globalFlags,
		TokenID:         inputArgs[0],
		InstanceID:      flags.FlagToStringValue(p, cmd, instanceIDFlag),
	}
	p.DebugInputModel(model)
	return model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *modelexperiments.APIClient) modelexperiments.ApiGetInstanceTokenRequest {
	return apiClient.DefaultAPI.GetInstanceToken(ctx, model.ProjectId, model.Region, model.TokenID, model.InstanceID)
}

func outputResult(p *print.Printer, model *inputModel, resp *modelexperiments.GetInstanceTokenResponse) error {
	if resp == nil {
		return fmt.Errorf("response token is nil")
	}
	return p.OutputResult(model.OutputFormat, resp, func() error {
		table := tables.NewTable()
		table.AddRow("ID", resp.Token.Id)
		table.AddSeparator()
		table.AddRow("NAME", resp.Token.Name)
		table.AddSeparator()
		table.AddRow("DESCRIPTION", utils.PtrString(resp.Token.Description))
		table.AddSeparator()
		table.AddRow("REGION", resp.Token.Region)
		table.AddSeparator()
		table.AddRow("STATE", string(resp.Token.State))
		table.AddSeparator()
		table.AddRow("VALID UNTIL", resp.Token.ValidUntil.String())
		table.AddSeparator()
		labels := ""
		if resp.Token.Labels != nil {
			labels = fmt.Sprint(*resp.Token.Labels)
		}
		table.AddRow("LABELS", labels)
		table.AddSeparator()
		if err := table.Display(p); err != nil {
			return fmt.Errorf("render table: %w", err)
		}
		return nil
	})
}
