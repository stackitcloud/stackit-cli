package list

import (
	"context"
	"encoding/json"
	"fmt"

	"stackit/internal/pkg/args"
	"stackit/internal/pkg/errors"
	"stackit/internal/pkg/examples"
	"stackit/internal/pkg/flags"
	"stackit/internal/pkg/globalflags"
	"stackit/internal/pkg/projectname"
	"stackit/internal/pkg/services/service-account/client"
	"stackit/internal/pkg/tables"

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

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all service accounts",
		Long:  "List all service accounts",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`List all service accounts`,
				"$ stackit service-account list"),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			model, err := parseInput(cmd)
			if err != nil {
				return err
			}

			// Configure API client
			apiClient, err := client.ConfigureClient(cmd)
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
				projectLabel, err := projectname.GetProjectName(ctx, cmd)
				if err != nil {
					projectLabel = model.ProjectId
				}
				cmd.Printf("No service accounts found for project %s\n", projectLabel)
				return nil
			}

			// Truncate output
			if model.Limit != nil && len(serviceAccounts) > int(*model.Limit) {
				serviceAccounts = serviceAccounts[:*model.Limit]
			}

			return outputResult(cmd, model.OutputFormat, serviceAccounts)
		},
	}

	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Int64(limitFlag, 0, "Maximum number of entries to list")
}

func parseInput(cmd *cobra.Command) (*inputModel, error) {
	globalFlags := globalflags.Parse(cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	limit := flags.FlagToInt64Pointer(cmd, limitFlag)
	if limit != nil && *limit < 1 {
		return nil, &errors.FlagValidationError{
			Flag:    limitFlag,
			Details: "must be greater than 0",
		}
	}

	return &inputModel{
		GlobalFlagModel: globalFlags,
		Limit:           flags.FlagToInt64Pointer(cmd, limitFlag),
	}, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *serviceaccount.APIClient) serviceaccount.ApiListServiceAccountsRequest {
	req := apiClient.ListServiceAccounts(ctx, model.ProjectId)
	return req
}

func outputResult(cmd *cobra.Command, outputFormat string, serviceAccounts []serviceaccount.ServiceAccount) error {
	switch outputFormat {
	case globalflags.JSONOutputFormat:
		details, err := json.MarshalIndent(serviceAccounts, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal service accounts list: %w", err)
		}
		cmd.Println(string(details))
	default:
		table := tables.NewTable()
		table.SetHeader("ID", "EMAIL")
		for i := range serviceAccounts {
			account := serviceAccounts[i]
			table.AddRow(*account.Id, *account.Email)
		}
		err := table.Display(cmd)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}
	}

	return nil
}
