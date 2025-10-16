package importKey

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/goccy/go-yaml"
	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	cliErr "github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/kms/client"
	kmsUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/kms/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-sdk-go/services/kms"
)

const (
	keyIdArg = "KEY_ID"

	keyRingIdFlag     = "key-ring-id"
	wrappedKeyFlag    = "wrapped-key"
	wrappingKeyIdFlag = "wrapping-key-id"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	KeyRingId     string
	KeyId         string
	WrappedKey    *string
	WrappingKeyId *string
}

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("import %s", keyIdArg),
		Short: "Import a KMS key",
		Long:  "Import a new version to the given KMS key.",
		Args:  args.SingleArg(keyIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Import a new version for the given KMS key "my-key-id"`,
				`$ stackit beta kms key import "my-key-id" --key-ring "my-keyring-id" --wrapped-key "base64-encoded-wrapped-key-material" --wrapping-key-id "my-wrapping-key-id"`),
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

			keyName, err := kmsUtils.GetKeyName(ctx, apiClient, model.ProjectId, model.Region, model.KeyRingId, model.KeyId)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get key name: %v", err)
				keyName = model.KeyId
			}
			keyRingName, err := kmsUtils.GetKeyRingName(ctx, apiClient, model.ProjectId, model.KeyRingId, model.Region)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get key ring name: %v", err)
				keyRingName = model.KeyRingId
			}

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to import a new version for the KMS Key %q inside the key ring %q?", keyName, keyRingName)
				err = params.Printer.PromptForConfirmation(prompt)
				if err != nil {
					return err
				}
			}

			// Call API
			req, _ := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("import KMS key: %w", err)
			}

			return outputResult(params.Printer, model.OutputFormat, keyRingName, keyName, resp)
		},
	}
	configureFlags(cmd)
	return cmd
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	keyId := inputArgs[0]

	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &cliErr.ProjectIdError{}
	}

	// WrappedKey needs to be base64 encoded
	var wrappedKey *string = flags.FlagToStringPointer(p, cmd, wrappedKeyFlag)
	_, err := base64.StdEncoding.DecodeString(*wrappedKey)
	if err != nil || *wrappedKey == "" {
		return nil, &cliErr.FlagValidationError{
			Flag:    wrappedKeyFlag,
			Details: "The 'wrappedKey' argument is required and needs to be base64 encoded.",
		}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		KeyId:           keyId,
		KeyRingId:       flags.FlagToStringValue(p, cmd, keyRingIdFlag),
		WrappedKey:      wrappedKey,
		WrappingKeyId:   flags.FlagToStringPointer(p, cmd, wrappingKeyIdFlag),
	}

	p.DebugInputModel(model)
	return &model, nil
}

type kmsKeyClient interface {
	ImportKey(ctx context.Context, projectId string, regionId string, keyRingId string, keyId string) kms.ApiImportKeyRequest
}

func buildRequest(ctx context.Context, model *inputModel, apiClient kmsKeyClient) (kms.ApiImportKeyRequest, error) {
	req := apiClient.ImportKey(ctx, model.ProjectId, model.Region, model.KeyRingId, model.KeyId)

	req = req.ImportKeyPayload(kms.ImportKeyPayload{
		WrappedKey:    model.WrappedKey,
		WrappingKeyId: model.WrappingKeyId,
	})
	return req, nil
}

func outputResult(p *print.Printer, outputFormat, keyRingName, keyName string, resp *kms.Version) error {
	if resp == nil {
		return fmt.Errorf("response is nil")
	}

	switch outputFormat {
	case print.JSONOutputFormat:
		details, err := json.MarshalIndent(resp, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal KMS key: %w", err)
		}
		p.Outputln(string(details))

	case print.YAMLOutputFormat:
		details, err := yaml.MarshalWithOptions(resp, yaml.IndentSequence(true), yaml.UseJSONMarshaler())
		if err != nil {
			return fmt.Errorf("marshal KMS key: %w", err)
		}
		p.Outputln(string(details))

	default:
		p.Outputf("Imported a new version for the key %q inside the key ring %q\n", keyName, keyRingName)
	}

	return nil
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), keyRingIdFlag, "ID of the KMS key ring")
	cmd.Flags().String(wrappedKeyFlag, "", "The wrapped key material that has to be imported. Encoded in base64")
	cmd.Flags().Var(flags.UUIDFlag(), wrappingKeyIdFlag, "The unique id of the wrapping key the key material has been wrapped with")

	err := flags.MarkFlagsRequired(cmd, keyRingIdFlag, wrappedKeyFlag, wrappingKeyIdFlag)
	cobra.CheckErr(err)
}
