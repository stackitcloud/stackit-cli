package create

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/service-account/client"

	"github.com/spf13/cobra"
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

func NewCmd(p *print.Printer) *cobra.Command {
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

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to create an access token for service account %s?", model.ServiceAccountEmail)
				err = p.PromptForConfirmation(cmd, prompt)
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

			p.Outputf("Created access token for service account %s. Token ID: %s\n\n", model.ServiceAccountEmail, *token.Id)
			p.Outputf("Valid until: %s\n", *token.ValidUntil)
			p.Outputf("Token: %s\n", *token.Token)
			return nil
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

	ttlDays := flags.FlagWithDefaultToInt64Value(cmd, ttlDaysFlag)
	if ttlDays < 1 {
		return nil, &errors.FlagValidationError{
			Flag:    serviceAccountEmailFlag,
			Details: "time to live should be at least 1 day",
		}
	}

	return &inputModel{
		GlobalFlagModel:     globalFlags,
		ServiceAccountEmail: email,
		TTLDays:             &ttlDays,
	}, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *serviceaccount.APIClient) serviceaccount.ApiCreateAccessTokenRequest {
	req := apiClient.CreateAccessToken(ctx, model.ProjectId, model.ServiceAccountEmail)
	req = req.CreateAccessTokenPayload(serviceaccount.CreateAccessTokenPayload{
		TtlDays: model.TTLDays,
	})
	return req
}
