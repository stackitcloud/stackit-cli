package list

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/service-account/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/serviceaccount"
)

const (
	limitFlag               = "limit"
	serviceAccountEmailFlag = "email"
)

type inputModel struct {
	*globalflags.GlobalFlagModel

	ServiceAccountEmail string
	Limit               *int64
}

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Lists all service account keys",
		Long:  "Lists all service account keys.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`List all keys belonging to the service account with email "my-service-account-1234567@sa.stackit.cloud"`,
				"$ stackit service-account key list --email my-service-account-1234567@sa.stackit.cloud"),
			examples.NewExample(
				`List all keys belonging to the service account with email "my-service-account-1234567@sa.stackit.cloud" in JSON format`,
				"$ stackit service-account key list --email my-service-account-1234567@sa.stackit.cloud --output-format json"),
			examples.NewExample(
				`List up to 10 keys belonging to the service account with email "my-service-account-1234567@sa.stackit.cloud"`,
				"$ stackit service-account key list --email my-service-account-1234567@sa.stackit.cloud --limit 10"),
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
				return fmt.Errorf("list keys metadata: %w", err)
			}
			keys := *resp.Items
			if len(keys) == 0 {
				p.Info("No keys found for service account %s\n", model.ServiceAccountEmail)
				return nil
			}

			// Truncate output
			if model.Limit != nil && len(keys) > int(*model.Limit) {
				keys = keys[:*model.Limit]
			}

			return outputResult(cmd, model.OutputFormat, keys, p)
		},
	}

	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().StringP(serviceAccountEmailFlag, "e", "", "Service account email")
	cmd.Flags().Int64(limitFlag, 0, "Maximum number of entries to list")

	err := flags.MarkFlagsRequired(cmd, serviceAccountEmailFlag)
	cobra.CheckErr(err)
}

func parseInput(cmd *cobra.Command) (*inputModel, error) {
	globalFlags := globalflags.Parse(cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	email := flags.FlagToStringValue(cmd, serviceAccountEmailFlag)
	if email == "" {
		return nil, &errors.FlagValidationError{
			Flag:    serviceAccountEmailFlag,
			Details: "can't be empty",
		}
	}

	limit := flags.FlagToInt64Pointer(cmd, limitFlag)
	if limit != nil && *limit < 1 {
		return nil, &errors.FlagValidationError{
			Flag:    limitFlag,
			Details: "must be greater than 0",
		}
	}

	return &inputModel{
		GlobalFlagModel:     globalFlags,
		ServiceAccountEmail: email,
		Limit:               limit,
	}, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *serviceaccount.APIClient) serviceaccount.ApiListServiceAccountKeysRequest {
	req := apiClient.ListServiceAccountKeys(ctx, model.ProjectId, model.ServiceAccountEmail)
	return req
}

func outputResult(cmd *cobra.Command, outputFormat string, keys []serviceaccount.ServiceAccountKeyListResponse, p *print.Printer) error {
	switch outputFormat {
	case globalflags.JSONOutputFormat:
		details, err := json.MarshalIndent(keys, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal keys metadata: %w", err)
		}
		p.Outputln(string(details))

		return nil
	default:
		table := tables.NewTable()
		table.SetHeader("ID", "ACTIVE", "CREATED_AT", "VALID_UNTIL")
		for i := range keys {
			k := keys[i]
			validUntil := "does not expire"
			if k.ValidUntil != nil {
				validUntil = k.ValidUntil.String()
			}
			table.AddRow(*k.Id, *k.Active, *k.CreatedAt, validUntil)
		}
		err := table.Display(cmd)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}

		return nil
	}
}
