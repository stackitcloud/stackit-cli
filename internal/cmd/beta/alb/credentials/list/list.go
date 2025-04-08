package list

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/alb/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"

	"github.com/goccy/go-yaml"
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

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Lists all credentials",
		Long:  "Lists all credentials.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Lists all credentials`,
				"$ stackit beta alb credential list",
			),
			examples.NewExample(
				`Lists all credentials in JSON format`,
				"$ stackit beta alb credential list --output-format json",
			),
			examples.NewExample(
				`Lists up to 10 credentials`,
				"$ stackit beta alb credential list --limit 10",
			),
		),
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := context.Background()
			model, err := parseInput(p, cmd)
			if err != nil {
				return err
			}

			// Configure API client
			apiClient, err := client.ConfigureClient(p)
			if err != nil {
				return err
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("list credentials: %w", err)
			}

			if resp.Credentials == nil || len(*resp.Credentials) == 0 {
				p.Info("No credentials found\n")
				return nil
			}

			items := *resp.Credentials
			if model.Limit != nil && len(items) > int(*model.Limit) {
				items = items[:*model.Limit]
			}

			return outputResult(p, model.OutputFormat, items)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Int64(limitFlag, 0, "Number of credentials to list")
}

func parseInput(p *print.Printer, cmd *cobra.Command) (*inputModel, error) {
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

	if p.IsVerbosityDebug() {
		modelStr, err := print.BuildDebugStrFromInputModel(model)
		if err != nil {
			p.Debug(print.ErrorLevel, "convert model to string for debugging: %v", err)
		} else {
			p.Debug(print.InfoLevel, modelStr)
		}
	}

	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *alb.APIClient) alb.ApiListCredentialsRequest {
	req := apiClient.ListCredentials(ctx, model.ProjectId, model.Region)
	return req
}

func outputResult(p *print.Printer, outputFormat string, items []alb.CredentialsResponse) error {
	if items == nil {
		p.Outputln("no credentials found")
		return nil
	}
	switch outputFormat {
	case print.JSONOutputFormat:
		details, err := json.MarshalIndent(items, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal credentials: %w", err)
		}
		p.Outputln(string(details))

	case print.YAMLOutputFormat:
		details, err := yaml.MarshalWithOptions(items, yaml.IndentSequence(true), yaml.UseJSONMarshaler())
		if err != nil {
			return fmt.Errorf("marshal credentials: %w", err)
		}
		p.Outputln(string(details))

	default:
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
	}
	return nil
}
