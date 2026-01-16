package update

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/stackitcloud/stackit-cli/internal/pkg/types"

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
	keyIdArg = "KEY_ID"

	serviceAccountEmailFlag = "email"
	expiredInDaysFlag       = "expires-in-days"
	activateFlag            = "activate"
	deactivateFlag          = "deactivate"
)

type inputModel struct {
	*globalflags.GlobalFlagModel

	ServiceAccountEmail string
	KeyId               string
	ExpiresInDays       *int64
	Activate            bool
	Deactivate          bool
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("update %s", keyIdArg),
		Short: "Updates a service account key",
		Long: fmt.Sprintf("%s\n%s",
			"Updates a service account key.",
			"You can temporarily activate or deactivate the key and/or update its date of expiration.",
		),
		Args: args.SingleArg(keyIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Temporarily deactivate a key with ID "xxx" of the service account with email "my-service-account-1234567@sa.stackit.cloud"`,
				"$ stackit service-account key update xxx --email my-service-account-1234567@sa.stackit.cloud --deactivate"),
			examples.NewExample(
				`Activate a key of the service account with email "my-service-account-1234567@sa.stackit.cloud"`,
				"$ stackit service-account key update xxx --email my-service-account-1234567@sa.stackit.cloud --activate"),
			examples.NewExample(
				`Update the expiration date of a key of the service account with email "my-service-account-1234567@sa.stackit.cloud"`,
				"$ stackit service-account key update xxx --email my-service-account-1234567@sa.stackit.cloud --expires-in-days 30"),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			model, err := parseInput(params.Printer, cmd, args)
			if err != nil {
				return err
			}

			// Configure API client
			apiClient, err := client.ConfigureClient(params.Printer, params.CliVersion)
			if err != nil {
				return err
			}

			prompt := fmt.Sprintf("Are you sure you want to update the key with ID %q?", model.KeyId)
			err = params.Printer.PromptForConfirmation(prompt)
			if err != nil {
				return err
			}

			// Call API
			req := buildRequest(ctx, model, apiClient, time.Now())
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("create service account key: %w", err)
			}

			key, err := json.MarshalIndent(resp, "", "  ")
			if err != nil {
				return fmt.Errorf("marshal key: %w", err)
			}
			params.Printer.Info("%s", string(key))
			return nil
		},
	}

	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().StringP(serviceAccountEmailFlag, "e", "", "Service account email")
	cmd.Flags().Int64P(expiredInDaysFlag, "", 0, "Number of days until expiration")
	cmd.Flags().Bool(activateFlag, false, "If set, activates the service account key")
	cmd.Flags().Bool(deactivateFlag, false, "If set, temporarily deactivates the service account key")

	err := flags.MarkFlagsRequired(cmd, serviceAccountEmailFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	keyId := inputArgs[0]

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

	activate := flags.FlagToBoolValue(p, cmd, activateFlag)
	deactivate := flags.FlagToBoolValue(p, cmd, deactivateFlag)
	if activate && deactivate {
		return nil, fmt.Errorf("only one of %q and %q can be set", activateFlag, deactivateFlag)
	}

	model := inputModel{
		GlobalFlagModel:     globalFlags,
		ServiceAccountEmail: email,
		KeyId:               keyId,
		ExpiresInDays:       expriresInDays,
		Activate:            activate,
		Deactivate:          deactivate,
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *serviceaccount.APIClient, now time.Time) serviceaccount.ApiPartialUpdateServiceAccountKeyRequest {
	req := apiClient.PartialUpdateServiceAccountKey(ctx, model.ProjectId, model.ServiceAccountEmail, model.KeyId)

	var validUntil *time.Time
	validUntil = nil
	if model.ExpiresInDays != nil {
		validUntil = utils.Ptr(daysFromNow(now, *model.ExpiresInDays))
	}

	var active *bool
	active = nil
	if model.Deactivate {
		active = utils.Ptr(false)
	}
	if model.Activate {
		active = utils.Ptr(true)
	}

	req = req.PartialUpdateServiceAccountKeyPayload(serviceaccount.PartialUpdateServiceAccountKeyPayload{
		ValidUntil: validUntil,
		Active:     active,
	})
	return req
}

func daysFromNow(now time.Time, days int64) time.Time {
	validUntil := now.AddDate(0, 0, int(days))
	return validUntil
}
