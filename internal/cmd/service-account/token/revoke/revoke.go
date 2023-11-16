package revoke

import (
	"context"
	"fmt"

	"stackit/internal/pkg/args"
	"stackit/internal/pkg/confirm"
	"stackit/internal/pkg/errors"
	"stackit/internal/pkg/examples"
	"stackit/internal/pkg/flags"
	"stackit/internal/pkg/globalflags"
	"stackit/internal/pkg/services/service-account/client"
	"stackit/internal/pkg/utils"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/serviceaccount"
)

const (
	serviceAccountEmailFlag = "email"
	tokenIdArg              = "TOKEN_ID"
)

type inputModel struct {
	*globalflags.GlobalFlagModel

	ServiceAccountEmail string
	TokenId             string
}

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("revoke %s", tokenIdArg),
		Short: "Revoke an access token of a service account",
		Long: fmt.Sprintf("%s\n%s\n%s",
			"Revoke an access token of a service account.",
			"The access token is instantly revoked, any following calls with the token will be unauthorized.",
			"The token metadata is still stored until the expiration time.",
		),
		Args: args.SingleArg(tokenIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Revoke an access token with ID "xxx" of the service account with email "my-service-account-1234567@sa.stackit.cloud"`,
				"$ stackit service-account token revoke xxx --email my-service-account-1234567@sa.stackit.cloud"),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			model, err := parseInput(cmd, args)
			if err != nil {
				return err
			}

			// Configure API client
			apiClient, err := client.ConfigureClient(cmd)
			if err != nil {
				return err
			}

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to revoke the access token with ID %s?", model.TokenId)
				err = confirm.PromptForConfirmation(cmd, prompt)
				if err != nil {
					return err
				}
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			err = req.Execute()
			if err != nil {
				return fmt.Errorf("revoke access token: %w", err)
			}

			cmd.Printf("Revoked access token with ID %s\n", model.TokenId)
			return nil
		},
	}

	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().StringP(serviceAccountEmailFlag, "e", "", "Service account email")

	err := flags.MarkFlagsRequired(cmd, serviceAccountEmailFlag)
	cobra.CheckErr(err)
}

func parseInput(cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	tokenId := inputArgs[0]

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

	return &inputModel{
		GlobalFlagModel:     globalFlags,
		ServiceAccountEmail: email,
		TokenId:             tokenId,
	}, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *serviceaccount.APIClient) serviceaccount.ApiDeleteAccessTokenRequest {
	req := apiClient.DeleteAccessToken(ctx, model.ProjectId, model.ServiceAccountEmail, model.TokenId)
	return req
}
