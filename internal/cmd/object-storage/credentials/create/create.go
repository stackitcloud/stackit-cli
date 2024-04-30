package create

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	objectStorageUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/object-storage/utils"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/object-storage/client"
	"github.com/stackitcloud/stackit-sdk-go/services/objectstorage"
)

const (
	expireDateFlag         = "expire-date"
	credentialsGroupIdFlag = "credentials-group-id"
	expirationTimeFormat   = time.RFC3339
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	ExpireDate         *time.Time
	CredentialsGroupId string
	HidePassword       bool
}

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Creates credentials for an Object Storage credentials group",
		Long:  "Creates credentials for an Object Storage credentials group. The credentials are only displayed upon creation, and will not be retrievable later.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Create credentials for a credentials group with ID "xxx"`,
				"$ stackit object-storage credentials create --credentials-group-id xxx"),
			examples.NewExample(
				`Create credentials for a credentials group with ID "xxx", including a specific expiration date`,
				"$ stackit object-storage credentials create --credentials-group-id xxx --expire-date 2024-03-06T00:00:00.000Z"),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
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

			credentialsGroupLabel, err := objectStorageUtils.GetCredentialsGroupName(ctx, apiClient, model.ProjectId, model.CredentialsGroupId)
			if err != nil {
				p.Debug(print.ErrorLevel, "get credentials group name: %v", err)
				credentialsGroupLabel = model.CredentialsGroupId
			}

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to create credentials in group %q?", credentialsGroupLabel)
				err = p.PromptForConfirmation(prompt)
				if err != nil {
					return err
				}
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("create Object Storage credentials: %w", err)
			}

			return outputResult(p, model, credentialsGroupLabel, resp)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().String(expireDateFlag, "", "Expiration date for the credentials, in a date-time with the RFC3339 layout format, e.g. 2024-01-01T00:00:00Z")
	cmd.Flags().Var(flags.UUIDFlag(), credentialsGroupIdFlag, "Credentials Group ID")

	err := flags.MarkFlagsRequired(cmd, credentialsGroupIdFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	expireDate, err := flags.FlagToDateTimePointer(p, cmd, expireDateFlag, expirationTimeFormat)
	if err != nil {
		return nil, &errors.FlagValidationError{
			Flag:    expireDateFlag,
			Details: err.Error(),
		}
	}

	model := inputModel{
		GlobalFlagModel:    globalFlags,
		ExpireDate:         expireDate,
		CredentialsGroupId: flags.FlagToStringValue(p, cmd, credentialsGroupIdFlag),
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

func buildRequest(ctx context.Context, model *inputModel, apiClient *objectstorage.APIClient) objectstorage.ApiCreateAccessKeyRequest {
	req := apiClient.CreateAccessKey(ctx, model.ProjectId)
	req = req.CredentialsGroup(model.CredentialsGroupId)
	req = req.CreateAccessKeyPayload(objectstorage.CreateAccessKeyPayload{
		Expires: model.ExpireDate,
	})
	return req
}

func outputResult(p *print.Printer, model *inputModel, credentialsGroupLabel string, resp *objectstorage.CreateAccessKeyResponse) error {
	switch model.OutputFormat {
	case print.JSONOutputFormat:
		details, err := json.MarshalIndent(resp, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal Object Storage credentials: %w", err)
		}
		p.Outputln(string(details))

		return nil
	default:
		expireDate := "Never"
		if resp.Expires != nil && *resp.Expires != "" {
			expireDate = *resp.Expires
		}

		p.Outputf("Created credentials in group %q. Credentials ID: %s\n\n", credentialsGroupLabel, *resp.KeyId)
		p.Outputf("Access Key ID: %s\n", *resp.AccessKey)
		p.Outputf("Secret Access Key: %s\n", *resp.SecretAccessKey)
		p.Outputf("Expire Date: %s\n", expireDate)

		return nil
	}
}
