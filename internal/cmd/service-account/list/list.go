package list

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/goccy/go-yaml"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/projectname"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/service-account/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/serviceaccount"
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
		Short: "Lists all service accounts",
		Long:  "Lists all service accounts.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`List all service accounts`,
				"$ stackit service-account list"),
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
				return fmt.Errorf("list service accounts: %w", err)
			}
			serviceAccounts := *resp.Items
			if len(serviceAccounts) == 0 {
				projectLabel, err := projectname.GetProjectName(ctx, p, cmd)
				if err != nil {
					p.Debug(print.ErrorLevel, "get project name: %v", err)
					projectLabel = model.ProjectId
				}
				p.Info("No service accounts found for project %q\n", projectLabel)
				return nil
			}

			// Truncate output
			if model.Limit != nil && len(serviceAccounts) > int(*model.Limit) {
				serviceAccounts = serviceAccounts[:*model.Limit]
			}

			return outputResult(p, model.OutputFormat, serviceAccounts)
		},
	}

	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Int64(limitFlag, 0, "Maximum number of entries to list")
}

func parseInput(p *print.Printer, cmd *cobra.Command) (*inputModel, error) {
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
		Limit:           flags.FlagToInt64Pointer(p, cmd, limitFlag),
	}

	if p.IsVerbosityDebug() {
		modelStr, err := print.BuildDebugStrFromInputModel(model)
		if err != nil {
			p.Debug(print.ErrorLevel, "convert model to string for debugging: %v", err)
		} else {
			p.Debug(print.DebugLevel, "parsed input values: %s", modelStr)
		}
	}

	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *serviceaccount.APIClient) serviceaccount.ApiListServiceAccountsRequest {
	req := apiClient.ListServiceAccounts(ctx, model.ProjectId)
	return req
}

func outputResult(p *print.Printer, outputFormat string, serviceAccounts []serviceaccount.ServiceAccount) error {
	switch outputFormat {
	case print.JSONOutputFormat:
		details, err := json.MarshalIndent(serviceAccounts, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal service accounts list: %w", err)
		}
		p.Outputln(string(details))
	case print.YAMLOutputFormat:
		details, err := yaml.MarshalWithOptions(serviceAccounts, yaml.IndentSequence(true))
		if err != nil {
			return fmt.Errorf("marshal service accounts list: %w", err)
		}
		p.Outputln(string(details))
	default:
		table := tables.NewTable()
		table.SetHeader("ID", "EMAIL")
		for i := range serviceAccounts {
			account := serviceAccounts[i]
			table.AddRow(
				utils.PtrString(account.Id),
				utils.PtrString(account.Email),
			)
		}
		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}
	}

	return nil
}
