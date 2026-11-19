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
	"github.com/stackitcloud/stackit-cli/internal/pkg/projectname"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/modelexperiments/client"
	meUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/modelexperiments/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	modelexperiments "github.com/stackitcloud/stackit-sdk-go/services/modelexperiments/v1api"
)

const (
	instanceIDFlag = "instance-id"
	limitFlag      = "limit"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	InstanceID string
	Limit      *int64
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Lists all AI Model Experiments instance tokens",
		Long:  "Lists all auth tokens for an AI Model Experiments instance.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample("List all tokens for an instance", "$ stackit model-experiments token list --instance-id xxx"),
			examples.NewExample("List tokens in region eu02", "$ stackit model-experiments token list --instance-id xxx --region eu02"),
			examples.NewExample("List tokens in JSON format", "$ stackit model-experiments token list --instance-id xxx --output-format json"),
			examples.NewExample("List up to 10 tokens for an instance", "$ stackit model-experiments token list --instance-id xxx --limit 10"),
		),
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := context.Background()
			model, err := parseInput(params.Printer, cmd, nil)
			if err != nil {
				return err
			}
			apiClient, err := client.ConfigureClient(params.Printer, params.CliVersion)
			if err != nil {
				return err
			}
			resp, err := buildRequest(ctx, model, apiClient).Execute()
			if err != nil {
				return fmt.Errorf("list instance tokens: %w", err)
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
			applyLimit(resp, model.Limit)
			return outputResult(params.Printer, model, projectLabel, instanceLabel, resp)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), instanceIDFlag, "ID of the instance")
	cmd.Flags().Int64(limitFlag, 0, "Maximum number of entries to list")
	err := flags.MarkFlagsRequired(cmd, instanceIDFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command, _ []string) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &cliErr.ProjectIdError{}
	}
	limit := flags.FlagToInt64Pointer(p, cmd, limitFlag)
	if limit != nil && *limit < 1 {
		return nil, &cliErr.FlagValidationError{
			Flag:    limitFlag,
			Details: "must be greater than 0",
		}
	}
	model := &inputModel{
		GlobalFlagModel: globalFlags,
		InstanceID:      flags.FlagToStringValue(p, cmd, instanceIDFlag),
		Limit:           limit,
	}
	p.DebugInputModel(model)
	return model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *modelexperiments.APIClient) modelexperiments.ApiListInstanceTokensRequest {
	return apiClient.DefaultAPI.ListInstanceTokens(ctx, model.ProjectId, model.Region, model.InstanceID)
}

func applyLimit(resp *modelexperiments.ListInstanceTokensResponse, limit *int64) {
	if resp == nil || limit == nil || int64(len(resp.Tokens)) <= *limit {
		return
	}
	resp.Tokens = resp.Tokens[:int(*limit)]
}

func outputResult(p *print.Printer, model *inputModel, projectLabel, instanceLabel string, resp *modelexperiments.ListInstanceTokensResponse) error {
	if resp == nil {
		return fmt.Errorf("response is nil")
	}
	return p.OutputResult(model.OutputFormat, resp, func() error {
		if len(resp.Tokens) == 0 {
			p.Outputf("No tokens found for instance %q in project %q\n", instanceLabel, projectLabel)
			return nil
		}

		table := tables.NewTable()
		table.SetHeader("ID", "NAME", "DESCRIPTION", "REGION", "STATE", "VALID UNTIL")
		for _, token := range resp.Tokens {
			table.AddRow(token.Id, token.Name, utils.PtrString(token.Description), token.Region, token.State, token.ValidUntil.String())
		}
		return table.Display(p)
	})
}
