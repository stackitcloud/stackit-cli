package describe

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
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/serviceaccount"
)

const (
	keyIdArg = "KEY_ID"

	serviceAccountEmailFlag = "email"
)

type inputModel struct {
	*globalflags.GlobalFlagModel

	ServiceAccountEmail string
	KeyId               string
}

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("describe %s", keyIdArg),
		Short: "Shows details of a service account key",
		Long:  "Shows details of a service account key. Only JSON output is supported.",
		Args:  args.SingleArg(keyIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Get details of a service account key with ID "xxx" belonging to the service account with email "my-service-account-1234567@sa.stackit.cloud"`,
				"$ stackit service-account key describe xxx --email my-service-account-1234567@sa.stackit.cloud"),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			model, err := parseInput(p, cmd, args)
			if err != nil {
				return err
			}
			// Configure API client
			apiClient, err := client.ConfigureClient(p)
			if err != nil {
				return err
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("read service account key: %w", err)
			}

			return outputResult(p, resp)
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

	return &inputModel{
		GlobalFlagModel:     globalFlags,
		ServiceAccountEmail: email,
		KeyId:               keyId,
	}, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *serviceaccount.APIClient) serviceaccount.ApiGetServiceAccountKeyRequest {
	req := apiClient.GetServiceAccountKey(ctx, model.ProjectId, model.ServiceAccountEmail, model.KeyId)
	return req
}

func outputResult(p *print.Printer, key *serviceaccount.GetServiceAccountKeyResponse) error {
	marshaledKey, err := json.MarshalIndent(key, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal service account key: %w", err)
	}
	p.Outputln(string(marshaledKey))
	return nil
}
