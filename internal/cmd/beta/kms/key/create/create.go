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

	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/projectname"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/kms"
)

const (
	keyRingIdFlag = "key-ring"

	algorithmFlag   = "algorithm"
	backendFlag     = "backend"
	descriptionFlag = "description"
	displayNameFlag = "name"
	importOnlyFlag  = "import-only"
	purposeFlag     = "purpose"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	KeyRingId string

	Algorithm   *string
	Backend     string // Keep "backend" as a variable, but set the default to "software" (see UI)
	Description *string
	Name        *string
	ImportOnly  *bool // Default false
	Purpose     *string
}

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Creates a KMS Key",
		Long:  "Creates a KMS Key.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Create a Symmetric KMS Key`,
				`$ stakit beta kms key create --key-ring "my-keyring-id" --algorithm "rsa_2048_oaep_sha256" --name "my-key-name" --purpose "symmetric_encrypt_decrypt"`),
			examples.NewExample(
				`Create a Message Authentication KMS Key`,
				`$ stakit beta kms key create --key-ring "my-keyring-id" --algorithm "hmac_sha512" --name "my-key-name" --purpose "message_authentication_code"`),
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

			projectLabel, err := projectname.GetProjectName(ctx, params.Printer, params.CliVersion, cmd)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get project name: %v", err)
				projectLabel = model.ProjectId
			}

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to create a KMS Key for project %q?", projectLabel)
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

			key, err := req.Execute()
			if err != nil {
				return fmt.Errorf("create KMS Key: %w", err)
			}

			// No wait exists for the key creation
			return outputResult(params.Printer, model.OutputFormat, projectLabel, key)
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
		Algorithm:       flags.FlagToStringPointer(p, cmd, algorithmFlag),
		Backend:         flags.FlagWithDefaultToStringValue(p, cmd, backendFlag),
		Name:            flags.FlagToStringPointer(p, cmd, displayNameFlag),
		Description:     flags.FlagToStringPointer(p, cmd, descriptionFlag),
		ImportOnly:      flags.FlagToBoolPointer(p, cmd, importOnlyFlag),
		Purpose:         flags.FlagToStringPointer(p, cmd, purposeFlag),
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
	CreateKey(ctx context.Context, projectId string, regionId string, keyRingId string) kms.ApiCreateKeyRequest
}

func buildRequest(ctx context.Context, model *inputModel, apiClient kmsKeyClient) (kms.ApiCreateKeyRequest, error) {
	req := apiClient.CreateKey(ctx, model.ProjectId, model.Region, model.KeyRingId)

	// Question: Should there be additional checks here?
	req = req.CreateKeyPayload(kms.CreateKeyPayload{
		DisplayName: model.Name,
		Description: model.Description,
		Algorithm:   kms.CreateKeyPayloadGetAlgorithmAttributeType(model.Algorithm),
		Backend:     kms.CreateKeyPayloadGetBackendAttributeType(&model.Backend),
		Purpose:     kms.CreateKeyPayloadGetPurposeAttributeType(model.Purpose),
		ImportOnly:  model.ImportOnly,
	})
	return req, nil
}

func outputResult(p *print.Printer, outputFormat, projectLabel string, resp *kms.Key) error {
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
		p.Outputf("Created Key for project %q. Key ID: %s\n", projectLabel, utils.PtrString(resp.Id))
		return nil
	}
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), keyRingIdFlag, "ID of the KMS Key Ring")
	cmd.Flags().String(algorithmFlag, "", "En-/Decryption / signing algorithm")
	cmd.Flags().String(backendFlag, "software", "The backend that is responsible for maintaining this key")
	cmd.Flags().String(displayNameFlag, "", "The display name to distinguish multiple keys")
	cmd.Flags().String(descriptionFlag, "", "Optinal description of the Key")
	cmd.Flags().Bool(importOnlyFlag, false, "States whether versions can be created or only imported")
	cmd.Flags().String(purposeFlag, "", "Purpose of the Key. Enum: 'symmetric_encrypt_decrypt', 'asymmetric_encrypt_decrypt', 'message_authentication_code', 'asymmetric_sign_verify' ")

	err := flags.MarkFlagsRequired(cmd, keyRingIdFlag, algorithmFlag, purposeFlag, displayNameFlag)
	cobra.CheckErr(err)
}
