package create

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/types"

	"github.com/spf13/cobra"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	cliErr "github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/kms/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/spinner"

	kms "github.com/stackitcloud/stackit-sdk-go/services/kms/v1api"
	"github.com/stackitcloud/stackit-sdk-go/services/kms/v1api/wait"

	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

const (
	keyRingIdFlag = "keyring-id"

	descriptionFlag = "description"
	displayNameFlag = "name"
	importOnlyFlag  = "import-only"
)

var (
	algorithmFlag = flags.StringEnumFlag(
		"algorithm",
		kms.AllowedAlgorithmEnumValues,
		"En-/Decryption / signing algorithm.",
	)
	purposeFlag = flags.StringEnumFlag(
		"purpose",
		kms.AllowedPurposeEnumValues,
		"Purpose of the key.",
	)
	protectionFlag = flags.StringEnumFlag(
		"protection",
		kms.AllowedProtectionEnumValues,
		"The underlying system that is responsible for protecting the key material.")
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	KeyRingId string

	Algorithm   kms.Algorithm
	Description *string
	Name        *string
	ImportOnly  bool // Default false
	Purpose     kms.Purpose
	Protection  kms.Protection
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Creates a KMS key",
		Long:  "Creates a KMS key.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Create a symmetric AES key (AES-256) with the name "symm-aes-gcm" under the key ring "my-keyring-id"`,
				`$ stackit kms key create --keyring-id "my-keyring-id" --algorithm "aes_256_gcm" --name "symm-aes-gcm" --purpose "symmetric_encrypt_decrypt" --protection "software"`),
			examples.NewExample(
				`Create an asymmetric RSA encryption key (RSA-2048)`,
				`$ stackit kms key create --keyring-id "my-keyring-id" --algorithm "rsa_2048_oaep_sha256" --name "prod-orders-rsa" --purpose "asymmetric_encrypt_decrypt" --protection "software"`),
			examples.NewExample(
				`Create a message authentication key (HMAC-SHA512)`,
				`$ stackit kms key create --keyring-id "my-keyring-id" --algorithm "hmac_sha512" --name "api-mac-key" --purpose "message_authentication_code" --protection "software"`),
			examples.NewExample(
				`Create an ECDSA P-256 key for signing & verification`,
				`$ stackit kms key create --keyring-id "my-keyring-id" --algorithm "ecdsa_p256_sha256" --name "signing-ecdsa-p256" --purpose "asymmetric_sign_verify" --protection "software"`),
			examples.NewExample(
				`Create an import-only key (versions must be imported)`,
				`$ stackit kms key create --keyring-id "my-keyring-id" --algorithm "rsa_2048_oaep_sha256" --name "ext-managed-rsa" --purpose "asymmetric_encrypt_decrypt" --protection "software" --import-only`),
			examples.NewExample(
				`Create a key and print the result as YAML`,
				`$ stackit kms key create --keyring-id "my-keyring-id" --algorithm "rsa_2048_oaep_sha256" --name "yaml-output-rsa" --purpose "asymmetric_encrypt_decrypt" --protection "software" --output yaml`),
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

			err = params.Printer.PromptForConfirmation("Are you sure you want to create a KMS Key?")
			if err != nil {
				return err
			}

			// Call API
			req, _ := buildRequest(ctx, model, apiClient.DefaultAPI)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("create KMS key: %w", err)
			}

			// Wait for async operation, if async mode not enabled
			if !model.Async {
				err := spinner.Run(params.Printer, "Creating key", func() error {
					_, err = wait.CreateOrUpdateKeyWaitHandler(ctx, apiClient.DefaultAPI, model.ProjectId, model.Region, model.KeyRingId, resp.Id).WaitWithContext(ctx)
					return err
				})
				if err != nil {
					return fmt.Errorf("wait for KMS key creation: %w", err)
				}
			}

			return outputResult(params.Printer, model, resp)
		},
	}
	configureFlags(cmd)
	return cmd
}

func parseInput(p *print.Printer, cmd *cobra.Command) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &cliErr.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		KeyRingId:       flags.FlagToStringValue(p, cmd, keyRingIdFlag),
		Algorithm:       algorithmFlag.Get(),
		Name:            flags.FlagToStringPointer(p, cmd, displayNameFlag),
		Description:     flags.FlagToStringPointer(p, cmd, descriptionFlag),
		ImportOnly:      flags.FlagToBoolValue(p, cmd, importOnlyFlag),
		Purpose:         purposeFlag.Get(),
		Protection:      protectionFlag.Get(),
	}

	p.DebugInputModel(model)
	return &model, nil
}

type kmsKeyClient interface {
	CreateKey(ctx context.Context, projectId string, regionId string, keyRingId string) kms.ApiCreateKeyRequest
}

func buildRequest(ctx context.Context, model *inputModel, apiClient kmsKeyClient) (kms.ApiCreateKeyRequest, error) {
	req := apiClient.CreateKey(ctx, model.ProjectId, model.Region, model.KeyRingId)

	req = req.CreateKeyPayload(kms.CreateKeyPayload{
		DisplayName: utils.PtrString(model.Name),
		Description: model.Description,
		Algorithm:   model.Algorithm,
		Purpose:     model.Purpose,
		ImportOnly:  &model.ImportOnly,
		Protection:  model.Protection,
	})
	return req, nil
}

func outputResult(p *print.Printer, model *inputModel, resp *kms.Key) error {
	if resp == nil {
		return fmt.Errorf("response is nil")
	}

	return p.OutputResult(model.OutputFormat, resp, func() error {
		operationState := "Created"
		if model.Async {
			operationState = "Triggered creation of"
		}
		p.Outputf("%s the KMS key %q. Key ID: %s\n", operationState, resp.DisplayName, resp.Id)
		return nil
	})
}

func configureFlags(cmd *cobra.Command) {
	algorithmFlag.Register(cmd)
	purposeFlag.Register(cmd)
	protectionFlag.Register(cmd)

	// All further non Enum Flags
	cmd.Flags().Var(flags.UUIDFlag(), keyRingIdFlag, "ID of the KMS key ring")
	cmd.Flags().String(displayNameFlag, "", "The display name to distinguish multiple keys")
	cmd.Flags().String(descriptionFlag, "", "Optional description of the key")
	cmd.Flags().Bool(importOnlyFlag, false, "States whether versions can be created or only imported")

	err := flags.MarkFlagsRequired(cmd, keyRingIdFlag, algorithmFlag.Name(), purposeFlag.Name(), displayNameFlag, protectionFlag.Name())
	cobra.CheckErr(err)
}
