package create

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/goccy/go-yaml"
	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/service-account/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/serviceaccount"
)

const (
	serviceAccountEmailFlag = "email"
	ttlDaysFlag             = "ttl-days"

	defaultTTLDays = 90
)

type inputModel struct {
	*globalflags.GlobalFlagModel

	ServiceAccountEmail string
	TTLDays             *int64
}

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Creates an access token for a service account",
		Long: fmt.Sprintf("%s\n%s\n%s",
			"Creates an access token for a service account.",
			"The access token can be then used for API calls (where enabled).",
			"The token is only displayed upon creation, and it will not be recoverable later.",
		),
		Args: args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Create an access token for the service account with email "my-service-account-1234567@sa.stackit.cloud" with a default time to live`,
				"$ stackit service-account token create --email my-service-account-1234567@sa.stackit.cloud"),
			examples.NewExample(
				`Create an access token for the service account with email "my-service-account-1234567@sa.stackit.cloud" with a time to live of 10 days`,
				"$ stackit service-account token create --email my-service-account-1234567@sa.stackit.cloud --ttl-days 10"),
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

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to create an access token for service account %s?", model.ServiceAccountEmail)
				err = params.Printer.PromptForConfirmation(prompt)
				if err != nil {
					return err
				}
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			token, err := req.Execute()
			if err != nil {
				return fmt.Errorf("create access token: %w", err)
			}

			return outputResult(params.Printer, model.OutputFormat, model.ServiceAccountEmail, token)
		},
	}

	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().StringP(serviceAccountEmailFlag, "e", "", "Service account email")
	cmd.Flags().Int64(ttlDaysFlag, defaultTTLDays, "How long (in days) the new access token is valid")

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
			Details: "can't be empty",
		}
	}

	ttlDays := flags.FlagWithDefaultToInt64Value(p, cmd, ttlDaysFlag)
	if ttlDays < 1 {
		return nil, &errors.FlagValidationError{
			Flag:    serviceAccountEmailFlag,
			Details: "time to live should be at least 1 day",
		}
	}

	model := inputModel{
		GlobalFlagModel:     globalFlags,
		ServiceAccountEmail: email,
		TTLDays:             &ttlDays,
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *serviceaccount.APIClient) serviceaccount.ApiCreateAccessTokenRequest {
	req := apiClient.CreateAccessToken(ctx, model.ProjectId, model.ServiceAccountEmail)
	req = req.CreateAccessTokenPayload(serviceaccount.CreateAccessTokenPayload{
		TtlDays: model.TTLDays,
	})
	return req
}

func outputResult(p *print.Printer, outputFormat, serviceAccountEmail string, token *serviceaccount.AccessToken) error {
	if token == nil {
		return fmt.Errorf("token is nil")
	}

	switch outputFormat {
	case print.JSONOutputFormat:
		details, err := json.MarshalIndent(token, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal service account access token: %w", err)
		}
		p.Outputln(string(details))

		return nil
	case print.YAMLOutputFormat:
		details, err := yaml.MarshalWithOptions(token, yaml.IndentSequence(true), yaml.UseJSONMarshaler())
		if err != nil {
			return fmt.Errorf("marshal service account access token: %w", err)
		}
		p.Outputln(string(details))

		return nil
	default:
		p.Outputf("Created access token for service account %s. Token ID: %s\n\n", serviceAccountEmail, utils.PtrString(token.Id))
		p.Outputf("Valid until: %s\n", utils.PtrString(token.ValidUntil))
		p.Outputf("Token: %s\n", utils.PtrString(token.Token))
		return nil
	}
}
