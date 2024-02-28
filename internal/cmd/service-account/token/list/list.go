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

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Lists access tokens of a service account",
		Long: fmt.Sprintf("%s\n%s\n%s",
			"Lists access tokens of a service account.",
			"Only the metadata about the access tokens is shown, and not the tokens themselves.",
			"Access tokens (including revoked tokens) are returned until they are expired.",
		),
		Args: args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`List all access tokens of the service account with email "my-service-account-1234567@sa.stackit.cloud"`,
				"$ stackit service-account token list --email my-service-account-1234567@sa.stackit.cloud"),
			examples.NewExample(
				`List all access tokens of the service account with email "my-service-account-1234567@sa.stackit.cloud" in JSON format`,
				"$ stackit service-account token list --email my-service-account-1234567@sa.stackit.cloud --output-format json"),
			examples.NewExample(
				`List up to 10 access tokens of the service account with email "my-service-account-1234567@sa.stackit.cloud"`,
				"$ stackit service-account token list --email my-service-account-1234567@sa.stackit.cloud --limit 10"),
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
				return fmt.Errorf("list tokens metadata: %w", err)
			}
			tokensMetadata := *resp.Items
			if len(tokensMetadata) == 0 {
				cmd.Printf("No tokens found for service account with email %q\n", model.ServiceAccountEmail)
				return nil
			}

			// Truncate output
			if model.Limit != nil && len(tokensMetadata) > int(*model.Limit) {
				tokensMetadata = tokensMetadata[:*model.Limit]
			}

			return outputResult(cmd, model.OutputFormat, tokensMetadata)
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
			Details: "can't be empty.",
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

func buildRequest(ctx context.Context, model *inputModel, apiClient *serviceaccount.APIClient) serviceaccount.ApiListAccessTokensRequest {
	req := apiClient.ListAccessTokens(ctx, model.ProjectId, model.ServiceAccountEmail)
	return req
}

func outputResult(cmd *cobra.Command, outputFormat string, tokensMetadata []serviceaccount.AccessTokenMetadata) error {
	switch outputFormat {
	case globalflags.JSONOutputFormat:
		details, err := json.MarshalIndent(tokensMetadata, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal tokens metadata: %w", err)
		}
		cmd.Println(string(details))

		return nil
	default:
		table := tables.NewTable()
		table.SetHeader("ID", "ACTIVE", "CREATED_AT", "VALID_UNTIL")
		for i := range tokensMetadata {
			t := tokensMetadata[i]
			table.AddRow(*t.Id, *t.Active, *t.CreatedAt, *t.ValidUntil)
		}
		err := table.Display(cmd)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}

		return nil
	}
}
