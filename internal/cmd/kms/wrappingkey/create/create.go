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
)

var (
	algorithmFlag = flags.StringEnumFlag(
		"algorithm",
		kms.AllowedWrappingAlgorithmEnumValues,
		"En-/Decryption / signing algorithm.",
	)
	purposeFlag = flags.StringEnumFlag(
		"purpose",
		kms.AllowedWrappingPurposeEnumValues,
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

	Algorithm   kms.WrappingAlgorithm
	Description *string
	Name        *string
	Purpose     kms.WrappingPurpose
	Protection  kms.Protection
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Creates a KMS wrapping key",
		Long:  "Creates a KMS wrapping key.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Create a symmetric (RSA + AES) KMS wrapping key with name "my-wrapping-key-name" in key ring with ID "my-keyring-id"`,
				`$ stackit kms wrapping-key create --keyring-id "my-keyring-id" --algorithm "rsa_2048_oaep_sha256_aes_256_key_wrap" --name "my-wrapping-key-name" --purpose "wrap_symmetric_key" --protection "software"`),
			examples.NewExample(
				`Create an asymmetric (RSA) KMS wrapping key with name "my-wrapping-key-name" in key ring with ID "my-keyring-id"`,
				`$ stackit kms wrapping-key create --keyring-id "my-keyring-id" --algorithm "rsa_3072_oaep_sha256" --name "my-wrapping-key-name" --purpose "wrap_asymmetric_key" --protection "software"`),
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

			err = params.Printer.PromptForConfirmation("Are you sure you want to create a KMS wrapping key?")
			if err != nil {
				return err
			}

			// Call API
			req, _ := buildRequest(ctx, model, apiClient.DefaultAPI)
			wrappingKey, err := req.Execute()
			if err != nil {
				return fmt.Errorf("create KMS wrapping key: %w", err)
			}

			// Wait for async operation, if async mode not enabled
			if !model.Async {
				err := spinner.Run(params.Printer, "Creating wrapping key", func() error {
					_, err = wait.CreateWrappingKeyWaitHandler(ctx, apiClient.DefaultAPI, model.ProjectId, model.Region, wrappingKey.KeyRingId, wrappingKey.Id).WaitWithContext(ctx)
					return err
				})
				if err != nil {
					return fmt.Errorf("wait for KMS wrapping key creation: %w", err)
				}
			}

			return outputResult(params.Printer, model, wrappingKey)
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

	// All values are mandatory strings. No additional type check required.
	model := inputModel{
		GlobalFlagModel: globalFlags,
		KeyRingId:       flags.FlagToStringValue(p, cmd, keyRingIdFlag),
		Algorithm:       algorithmFlag.Get(),
		Name:            flags.FlagToStringPointer(p, cmd, displayNameFlag),
		Description:     flags.FlagToStringPointer(p, cmd, descriptionFlag),
		Purpose:         purposeFlag.Get(),
		Protection:      protectionFlag.Get(),
	}

	p.DebugInputModel(model)
	return &model, nil
}

type kmsWrappingKeyClient interface {
	CreateWrappingKey(ctx context.Context, projectId string, regionId string, keyRingId string) kms.ApiCreateWrappingKeyRequest
}

func buildRequest(ctx context.Context, model *inputModel, apiClient kmsWrappingKeyClient) (kms.ApiCreateWrappingKeyRequest, error) {
	req := apiClient.CreateWrappingKey(ctx, model.ProjectId, model.Region, model.KeyRingId)

	req = req.CreateWrappingKeyPayload(kms.CreateWrappingKeyPayload{
		DisplayName: utils.PtrString(model.Name),
		Description: model.Description,
		Algorithm:   model.Algorithm,
		Purpose:     model.Purpose,
		Protection:  model.Protection,
	})
	return req, nil
}

func outputResult(p *print.Printer, model *inputModel, resp *kms.WrappingKey) error {
	if resp == nil {
		return fmt.Errorf("response is nil")
	}

	return p.OutputResult(model.OutputFormat, resp, func() error {
		operationState := "Created"
		if model.Async {
			operationState = "Triggered creation of"
		}
		p.Outputf("%s wrapping key. Wrapping key ID: %s\n", operationState, resp.Id)
		return nil
	})
}

func configureFlags(cmd *cobra.Command) {
	algorithmFlag.Register(cmd)
	purposeFlag.Register(cmd)
	protectionFlag.Register(cmd)

	// All further non Enum Flags
	cmd.Flags().Var(flags.UUIDFlag(), keyRingIdFlag, "ID of the KMS key ring")
	cmd.Flags().String(displayNameFlag, "", "The display name to distinguish multiple wrapping keys")
	cmd.Flags().String(descriptionFlag, "", "Optional description of the wrapping key")

	err := flags.MarkFlagsRequired(cmd, keyRingIdFlag, algorithmFlag.Name(), purposeFlag.Name(), displayNameFlag, protectionFlag.Name())
	cobra.CheckErr(err)
}
