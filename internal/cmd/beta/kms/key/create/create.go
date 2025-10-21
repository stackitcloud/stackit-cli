package create

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
	"github.com/stackitcloud/stackit-cli/internal/pkg/spinner"

	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/kms"
	"github.com/stackitcloud/stackit-sdk-go/services/kms/wait"
)

const (
	keyRingIdFlag = "keyring-id"

	algorithmFlag   = "algorithm"
	descriptionFlag = "description"
	displayNameFlag = "name"
	importOnlyFlag  = "import-only"
	purposeFlag     = "purpose"
	protectionFlag  = "protection"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	KeyRingId string

	Algorithm   *string
	Description *string
	Name        *string
	ImportOnly  bool // Default false
	Purpose     *string
	Protection  *string
}

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Creates a KMS key",
		Long:  "Creates a KMS key.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Create a Symmetric KMS key`,
				`$ stackit beta kms key create --keyring-id "MY_KEYRING_ID" --algorithm "rsa_2048_oaep_sha256" --name "my-key-name" --purpose "asymmetric_encrypt_decrypt" --protection "software"`),
			examples.NewExample(
				`Create a Message Authentication KMS key`,
				`$ stackit beta kms key create --keyring-id "MY_KEYRING_ID" --algorithm "hmac_sha512" --name "my-key-name" --purpose "message_authentication_code" --protection "software"`),
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
				err = params.Printer.PromptForConfirmation("Are you sure you want to create a KMS Key?")
				if err != nil {
					return err
				}
			}

			// Call API
			req, _ := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("create KMS key: %w", err)
			}

			// Wait for async operation, if async mode not enabled
			if !model.Async {
				s := spinner.New(params.Printer)
				s.Start("Creating key")
				_, err = wait.CreateOrUpdateKeyWaitHandler(ctx, apiClient, model.ProjectId, model.Region, model.KeyRingId, *resp.Id).WaitWithContext(ctx)
				if err != nil {
					return fmt.Errorf("wait for KMS key creation: %w", err)
				}
				s.Stop()
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
		Algorithm:       flags.FlagToStringPointer(p, cmd, algorithmFlag),
		Name:            flags.FlagToStringPointer(p, cmd, displayNameFlag),
		Description:     flags.FlagToStringPointer(p, cmd, descriptionFlag),
		ImportOnly:      flags.FlagToBoolValue(p, cmd, importOnlyFlag),
		Purpose:         flags.FlagToStringPointer(p, cmd, purposeFlag),
		Protection:      flags.FlagToStringPointer(p, cmd, protectionFlag),
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
		DisplayName: model.Name,
		Description: model.Description,
		Algorithm:   kms.CreateKeyPayloadGetAlgorithmAttributeType(model.Algorithm),
		Purpose:     kms.CreateKeyPayloadGetPurposeAttributeType(model.Purpose),
		ImportOnly:  &model.ImportOnly,
		Protection:  kms.CreateKeyPayloadGetProtectionAttributeType(model.Protection),
	})
	return req, nil
}

func outputResult(p *print.Printer, model *inputModel, resp *kms.Key) error {
	if resp == nil {
		return fmt.Errorf("response is nil")
	}

	switch model.OutputFormat {
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
		operationState := "Created"
		if model.Async {
			operationState = "Triggered creation of"
		}
		p.Outputf("%s the KMS key '%s'. Key ID: %s\n", operationState, utils.PtrString(resp.DisplayName), utils.PtrString(resp.Id))
	}
	return nil
}

func configureFlags(cmd *cobra.Command) {
	// Algorithm
	var algorithmFlagOptions []string
	for _, val := range kms.AllowedAlgorithmEnumValues {
		algorithmFlagOptions = append(algorithmFlagOptions, string(val))
	}
	cmd.Flags().Var(flags.EnumFlag(false, "", algorithmFlagOptions...), algorithmFlag, fmt.Sprintf("En-/Decryption / signing algorithm. Possible values: %q", algorithmFlagOptions))

	// Purpose
	var purposeFlagOptions []string
	for _, val := range kms.AllowedPurposeEnumValues {
		purposeFlagOptions = append(purposeFlagOptions, string(val))
	}
	cmd.Flags().Var(flags.EnumFlag(false, "", purposeFlagOptions...), purposeFlag, fmt.Sprintf("Purpose of the key. Possible values: %q", purposeFlagOptions))

	// Protection
	var protectionFlagOptions []string
	for _, val := range kms.AllowedProtectionEnumValues {
		protectionFlagOptions = append(protectionFlagOptions, string(val))
	}
	cmd.Flags().Var(flags.EnumFlag(false, "", protectionFlagOptions...), protectionFlag, fmt.Sprintf("The underlying system that is responsible for protecting the key material. Possible values: %q", purposeFlagOptions))

	// All further non Enum Flags
	cmd.Flags().Var(flags.UUIDFlag(), keyRingIdFlag, "ID of the KMS key ring")
	cmd.Flags().String(displayNameFlag, "", "The display name to distinguish multiple keys")
	cmd.Flags().String(descriptionFlag, "", "Optional description of the key")
	cmd.Flags().Bool(importOnlyFlag, false, "States whether versions can be created or only imported")

	err := flags.MarkFlagsRequired(cmd, keyRingIdFlag, algorithmFlag, purposeFlag, displayNameFlag, protectionFlag)
	cobra.CheckErr(err)
}
