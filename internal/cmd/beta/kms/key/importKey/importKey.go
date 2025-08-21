package importKey

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/goccy/go-yaml"
	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	cliErr "github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/kms/client"
	kmsUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/kms/utils"

	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-sdk-go/services/kms"
)

const (
	keyRingIdFlag = "key-ring"
	keyIdFlag     = "key"

	wrappedKeyFlag    = "wrapped-key"
	wrappingKeyIdFlag = "wrapping-key-id"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	KeyRingId string
	KeyId     string

	WrappedKey    *string
	WrappingKeyId *string
}

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "import",
		Short: "Import a KMS Key Version",
		Long:  "Improt a new version to the given KMS key.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Import a new version for the given KMS Key "my-key"`,
				`$ stakit beta kms key improt --key-ring "my-keyring-id" --key "my-key-id" --wrapped-key "base64-encoded-wrapped-key-material" --wrapping-key-id "my-wrapping-key-id"`),
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

			keyName, err := kmsUtils.GetKeyName(ctx, apiClient, model.ProjectId, model.Region, model.KeyRingId, model.KeyId)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get Key name: %v", err)
				keyName = model.KeyId
			}
			keyRingName, err := kmsUtils.GetKeyRingName(ctx, apiClient, model.ProjectId, model.KeyRingId, model.Region)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get Key Ring name: %v", err)
				keyRingName = model.KeyRingId
			}

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to import a new version for the KMS Key %q inside the Key Ring %q?", keyName, keyRingName)
				err = params.Printer.PromptForConfirmation(prompt)
				if err != nil {
					return err
				}
			}

			// Call API
			req, _ := buildRequest(ctx, model, apiClient)
			if err != nil {
				return err
			}

			keyVersion, err := req.Execute()
			if err != nil {
				return fmt.Errorf("import KMS Key: %w", err)
			}

			// No wait exists for the wrapped key import
			return outputResult(params.Printer, model.OutputFormat, keyRingName, keyName, keyVersion)
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

	// What checks could this need?
	// I would rather let the creation fail instead of checking all possible algorithms
	model := inputModel{
		GlobalFlagModel: globalFlags,
		KeyRingId:       flags.FlagToStringValue(p, cmd, keyRingIdFlag),
		KeyId:           flags.FlagToStringValue(p, cmd, keyIdFlag),
		WrappedKey:      flags.FlagToStringPointer(p, cmd, wrappedKeyFlag),
		WrappingKeyId:   flags.FlagToStringPointer(p, cmd, wrappingKeyIdFlag),
	}

	if p.IsVerbosityDebug() {
		modelStr, err := print.BuildDebugStrFromInputModel(model)
		if err != nil {
			p.Debug(print.ErrorLevel, "convert model to string for debugging: %v", err)
		} else {
			p.Debug(print.ErrorLevel, "parsed input values: %s", modelStr)
		}
	}

	return &model, nil
}

type kmsKeyClient interface {
	ImportKey(ctx context.Context, projectId string, regionId string, keyRingId string, keyId string) kms.ApiImportKeyRequest
}

func buildRequest(ctx context.Context, model *inputModel, apiClient kmsKeyClient) (kms.ApiImportKeyRequest, error) {
	req := apiClient.ImportKey(ctx, model.ProjectId, model.Region, model.KeyRingId, model.KeyId)

	// Question: Should there be additional checks here?
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
			return fmt.Errorf("marshal KMS Key: %w", err)
		}
		p.Outputln(string(details))
		return nil

	case print.YAMLOutputFormat:
		details, err := yaml.MarshalWithOptions(resp, yaml.IndentSequence(true), yaml.UseJSONMarshaler())
		if err != nil {
			return fmt.Errorf("marshal KMS Key: %w", err)
		}
		p.Outputln(string(details))
		return nil

	default:
		p.Outputf("Imported a new version for the Key %q inside the Key Ring %q\n", keyName, keyRingName)
		return nil
	}
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), keyRingIdFlag, "ID of the KMS Key Ring")
	cmd.Flags().Var(flags.UUIDFlag(), keyIdFlag, "ID of the KMS Key")
	cmd.Flags().String(wrappedKeyFlag, "", "The wrapped key material that has to be imported. Encoded in base64")
	cmd.Flags().Var(flags.UUIDFlag(), wrappingKeyIdFlag, "he unique id of the wrapping key the key material has been wrapped with")

	err := flags.MarkFlagsRequired(cmd, keyRingIdFlag, keyIdFlag, wrappedKeyFlag, wrappingKeyIdFlag)
	cobra.CheckErr(err)
}
