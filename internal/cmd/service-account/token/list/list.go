package list

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/goccy/go-yaml"
	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/service-account/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

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

func NewCmd(params *params.CmdParams) *cobra.Command {
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
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := context.Background()
			model, err := parseInput(params.Printer, cmd)
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
				return fmt.Errorf("list tokens metadata: %w", err)
			}
			tokensMetadata := *resp.Items
			if len(tokensMetadata) == 0 {
				params.Printer.Info("No tokens found for service account with email %q\n", model.ServiceAccountEmail)
				return nil
			}

			// Truncate output
			if model.Limit != nil && len(tokensMetadata) > int(*model.Limit) {
				tokensMetadata = tokensMetadata[:*model.Limit]
			}

			return outputResult(params.Printer, model.OutputFormat, tokensMetadata)
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

func parseInput(p *print.Printer, cmd *cobra.Command) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	email := flags.FlagToStringValue(p, cmd, serviceAccountEmailFlag)
	if email == "" {
		return nil, &errors.FlagValidationError{
			Flag:    serviceAccountEmailFlag,
			Details: "can't be empty.",
		}
	}

	limit := flags.FlagToInt64Pointer(p, cmd, limitFlag)
	if limit != nil && *limit < 1 {
		return nil, &errors.FlagValidationError{
			Flag:    limitFlag,
			Details: "must be greater than 0",
		}
	}

	model := inputModel{
		GlobalFlagModel:     globalFlags,
		ServiceAccountEmail: email,
		Limit:               limit,
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

func buildRequest(ctx context.Context, model *inputModel, apiClient *serviceaccount.APIClient) serviceaccount.ApiListAccessTokensRequest {
	req := apiClient.ListAccessTokens(ctx, model.ProjectId, model.ServiceAccountEmail)
	return req
}

func outputResult(p *print.Printer, outputFormat string, tokensMetadata []serviceaccount.AccessTokenMetadata) error {
	switch outputFormat {
	case print.JSONOutputFormat:
		details, err := json.MarshalIndent(tokensMetadata, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal tokens metadata: %w", err)
		}
		p.Outputln(string(details))

		return nil
	case print.YAMLOutputFormat:
		details, err := yaml.MarshalWithOptions(tokensMetadata, yaml.IndentSequence(true), yaml.UseJSONMarshaler())
		if err != nil {
			return fmt.Errorf("marshal tokens metadata: %w", err)
		}
		p.Outputln(string(details))

		return nil
	default:
		table := tables.NewTable()
		table.SetHeader("ID", "ACTIVE", "CREATED_AT", "VALID_UNTIL")
		for i := range tokensMetadata {
			t := tokensMetadata[i]
			table.AddRow(
				utils.PtrString(t.Id),
				utils.PtrString(t.Active),
				utils.PtrString(t.CreatedAt),
				utils.PtrString(t.ValidUntil),
			)
		}
		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}

		return nil
	}
}
