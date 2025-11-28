package list

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/alb/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/alb"
)

const (
	limitFlag = "limit"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	Limit *int64
}

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Lists all credentials",
		Long:  "Lists all credentials.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Lists all credentials`,
				"$ stackit beta alb observability-credentials list",
			),
			examples.NewExample(
				`Lists all credentials in JSON format`,
				"$ stackit beta alb observability-credentials list --output-format json",
			),
			examples.NewExample(
				`Lists up to 10 credentials`,
				"$ stackit beta alb observability-credentials list --limit 10",
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
				return fmt.Errorf("list credentials: %w", err)
			}
			items := resp.GetCredentials()

			// Truncate output
			if model.Limit != nil && len(items) > int(*model.Limit) {
				items = items[:*model.Limit]
			}

			return outputResult(params.Printer, model.OutputFormat, items)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Int64(limitFlag, 0, "Number of credentials to list")
}

func parseInput(p *print.Printer, cmd *cobra.Command, _ []string) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)

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
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *alb.APIClient) alb.ApiListCredentialsRequest {
	req := apiClient.ListCredentials(ctx, model.ProjectId, model.Region)
	return req
}

func outputResult(p *print.Printer, outputFormat string, items []alb.CredentialsResponse) error {
	return p.OutputResult(outputFormat, items, func() error {
		if len(items) == 0 {
			p.Outputf("No credentials found\n")
			return nil
		}

		table := tables.NewTable()
		table.SetHeader("CREDENTIAL REF", "DISPLAYNAME", "USERNAME", "REGION")

		for _, item := range items {
			table.AddRow(
				utils.PtrString(item.CredentialsRef),
				utils.PtrString(item.DisplayName),
				utils.PtrString(item.Username),
				utils.PtrString(item.Region),
			)
		}

		p.Outputln(table.Render())
		return nil
	})
}
