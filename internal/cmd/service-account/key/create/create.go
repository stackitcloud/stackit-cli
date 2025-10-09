package create

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/service-account/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/serviceaccount"
)

const (
	serviceAccountEmailFlag = "email"
	expiredInDaysFlag       = "expires-in-days"
	publicKeyFlag           = "public-key"
)

type inputModel struct {
	*globalflags.GlobalFlagModel

	ServiceAccountEmail string
	ExpiresInDays       *int64
	PublicKey           *string
}

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Creates a service account key",
		Long: fmt.Sprintf("%s\n%s\n%s",
			"Creates a service account key.",
			"You can generate an RSA keypair and provide the public key.",
			"If you do not provide a public key, the service will generate a new key-pair and the private key is included in the response. You won't be able to retrieve it later.",
		),
		Args: args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Create a key for the service account with email "my-service-account-1234567@sa.stackit.cloud with no expiration date"`,
				"$ stackit service-account key create --email my-service-account-1234567@sa.stackit.cloud"),
			examples.NewExample(
				`Create a key for the service account with email "my-service-account-1234567@sa.stackit.cloud" expiring in 10 days`,
				"$ stackit service-account key create --email my-service-account-1234567@sa.stackit.cloud --expires-in-days 10"),
			examples.NewExample(
				`Create a key for the service account with email "my-service-account-1234567@sa.stackit.cloud" and provide the public key in a .pem file"`,
				`$ stackit service-account key create --email my-service-account-1234567@sa.stackit.cloud --public-key @./public.pem`),
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
				validUntilInfo := "The key will be valid until deleted"
				if model.ExpiresInDays != nil {
					validUntilInfo = fmt.Sprintf("The key will be valid for %d days", *model.ExpiresInDays)
				}
				prompt := fmt.Sprintf("Are you sure you want to create a key for service account %s? %s", model.ServiceAccountEmail, validUntilInfo)
				err = params.Printer.PromptForConfirmation(prompt)
				if err != nil {
					return err
				}
			}

			// Call API
			req := buildRequest(ctx, model, apiClient, time.Now())
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("create service account key: %w", err)
			}

			params.Printer.Info("Created key for service account %s with ID %q\n", model.ServiceAccountEmail, *resp.Id)

			key, err := json.MarshalIndent(resp, "", "  ")
			if err != nil {
				return fmt.Errorf("marshal key: %w", err)
			}
			params.Printer.Outputln(string(key))
			return nil
		},
	}

	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().StringP(serviceAccountEmailFlag, "e", "", "Service account email")
	cmd.Flags().Int64P(expiredInDaysFlag, "", 0, "Number of days until expiration. When omitted, the key is valid until deleted")
	cmd.Flags().Var(flags.ReadFromFileFlag(), publicKeyFlag, `Public key of the user generated RSA 2048 key-pair. Must be in x509 format. Can be a string or path to the .pem file, if prefixed with "@"`)

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

	expriresInDays := flags.FlagToInt64Pointer(p, cmd, expiredInDaysFlag)
	if expriresInDays != nil && *expriresInDays < 1 {
		return nil, &errors.FlagValidationError{
			Flag:    expiredInDaysFlag,
			Details: "must be greater than 0",
		}
	}

	model := inputModel{
		GlobalFlagModel:     globalFlags,
		ServiceAccountEmail: email,
		ExpiresInDays:       expriresInDays,
		PublicKey:           flags.FlagToStringPointer(p, cmd, publicKeyFlag),
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *serviceaccount.APIClient, now time.Time) serviceaccount.ApiCreateServiceAccountKeyRequest {
	req := apiClient.CreateServiceAccountKey(ctx, model.ProjectId, model.ServiceAccountEmail)

	var validUntil *time.Time
	validUntil = nil
	if model.ExpiresInDays != nil {
		validUntil = utils.Ptr(daysFromNow(now, *model.ExpiresInDays))
	}

	req = req.CreateServiceAccountKeyPayload(serviceaccount.CreateServiceAccountKeyPayload{
		ValidUntil: validUntil,
		PublicKey:  model.PublicKey,
	})
	return req
}

func daysFromNow(now time.Time, days int64) time.Time {
	validUntil := now.AddDate(0, 0, int(days))
	return validUntil
}
