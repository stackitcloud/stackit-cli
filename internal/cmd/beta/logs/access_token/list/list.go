package list

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/types"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/projectname"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/logs/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/logs"

	"github.com/spf13/cobra"
)

const (
	limitFlag      = "limit"
	instanceIdFlag = "instance-id"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	InstanceId string
	Limit      *int64
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Lists all access tokens of a project",
		Long:  "Lists all access tokens of a project.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Lists all access tokens of the instance "xxx"`,
				"$ stackit logs access-token list --instance-id xxx",
			),
			examples.NewExample(
				`Lists all access tokens in JSON format`,
				"$ stackit logs access-token list --instance-id xxx --output-format json",
			),
			examples.NewExample(
				`Lists up to 10 access-token`,
				"$ stackit logs access-token list --instance-id xxx --limit 10",
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

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("list access tokens: %w", err)
			}

			projectLabel, err := projectname.GetProjectName(ctx, params.Printer, params.CliVersion, cmd)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get project name: %v", err)
				projectLabel = model.ProjectId
			} else if projectLabel == "" {
				projectLabel = model.ProjectId
			}

			// Truncate output
			items := *resp.Tokens
			if model.Limit != nil && len(items) > int(*model.Limit) {
				items = items[:*model.Limit]
			}

			return outputResult(params.Printer, model.OutputFormat, items, projectLabel)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Int64(limitFlag, 0, "Maximum number of entries to list")
	cmd.Flags().Var(flags.UUIDFlag(), instanceIdFlag, "ID of the logs instance")

	err := flags.MarkFlagsRequired(cmd, instanceIdFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command, _ []string) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	limit := flags.FlagToInt64Pointer(p, cmd, limitFlag)
	if limit != nil && *limit < 1 {
		return nil, &errors.FlagValidationError{
			Flag:    limitFlag,
			Details: "must be greater than 0",
		}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		Limit:           limit,
		InstanceId:      flags.FlagToStringValue(p, cmd, instanceIdFlag),
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *logs.APIClient) logs.ApiListAccessTokensRequest {
	return apiClient.ListAccessTokens(ctx, model.ProjectId, model.Region, model.InstanceId)
}

func outputResult(p *print.Printer, outputFormat string, tokens []logs.AccessToken, projectLabel string) error {
	return p.OutputResult(outputFormat, tokens, func() error {
		if len(tokens) == 0 {
			p.Outputf("No access token found for project %q\n", projectLabel)
			return nil
		}

		table := tables.NewTable()
		table.SetHeader("ID", "NAME", "DESCRIPTION", "PERMISSIONS", "VALID UNTIL", "STATUS")

		for _, token := range tokens {
			table.AddRow(
				utils.PtrString(token.Id),
				utils.PtrString(token.DisplayName),
				utils.PtrString(token.Description),
				utils.PtrString(token.Permissions),
				utils.PtrString(token.ValidUntil),
				utils.PtrString(token.Status),
			)
			table.AddSeparator()
		}

		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}
		return nil
	})
}
