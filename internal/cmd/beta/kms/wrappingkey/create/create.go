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
	"github.com/stackitcloud/stackit-cli/internal/pkg/projectname"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/kms"
	"github.com/stackitcloud/stackit-sdk-go/services/kms/wait"
)

const (
	keyRingIdFlag = "key-ring-id"

	algorithmFlag   = "algorithm"
	descriptionFlag = "description"
	displayNameFlag = "name"
	purposeFlag     = "purpose"
	protectionFlag  = "protection"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	KeyRingId string

	Algorithm   *string
	Description *string
	Name        *string
	Purpose     *string
	Protection  *string
}

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Creates a KMS wrapping key",
		Long:  "Creates a KMS wrapping key.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Create a Symmetric KMS wrapping key`,
				`$ stackit beta kms wrapping-key create --key-ring-id "my-keyring-id" --algorithm "rsa_2048_oaep_sha256" --name "my-wrapping-key-name" --purpose "wrap_symmetric_key" --protection "software"`),
			examples.NewExample(
				`Create an Asymmetric KMS wrapping key with a description`,
				`$ stackit beta kms wrapping-key create --key-ring-id "my-keyring-id" --algorithm "hmac_sha256" --name "my-wrapping-key-name" --description "my-description" --purpose "wrap_asymmetric_key" --protection "software"`),
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
				prompt := fmt.Sprintf("Are you sure you want to create a KMS wrapping key for project %q?", projectLabel)
				err = params.Printer.PromptForConfirmation(prompt)
				if err != nil {
					return err
				}
			}

			// Call API
			req, _ := buildRequest(ctx, model, apiClient)
			wrappingKey, err := req.Execute()
			if err != nil {
				return fmt.Errorf("create KMS wrapping key: %w", err)
			}

			// Wait for async operation, if async mode not enabled
			if !model.Async {
				s := spinner.New(params.Printer)
				s.Start("Creating instance")
				_, err = wait.CreateWrappingKeyWaitHandler(ctx, apiClient, model.ProjectId, model.Region, *wrappingKey.KeyRingId, *wrappingKey.Id).WaitWithContext(ctx)
				if err != nil {
					return fmt.Errorf("wait for KMS wrapping key creation: %w", err)
				}
				s.Stop()
			}

			return outputResult(params.Printer, model.OutputFormat, projectLabel, wrappingKey)
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
		Algorithm:       flags.FlagToStringPointer(p, cmd, algorithmFlag),
		Name:            flags.FlagToStringPointer(p, cmd, displayNameFlag),
		Description:     flags.FlagToStringPointer(p, cmd, descriptionFlag),
		Purpose:         flags.FlagToStringPointer(p, cmd, purposeFlag),
		Protection:      flags.FlagToStringPointer(p, cmd, protectionFlag),
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
		DisplayName: model.Name,
		Description: model.Description,
		Algorithm:   kms.CreateWrappingKeyPayloadGetAlgorithmAttributeType(model.Algorithm),
		Purpose:     kms.CreateWrappingKeyPayloadGetPurposeAttributeType(model.Purpose),
		Protection:  kms.CreateWrappingKeyPayloadGetProtectionAttributeType(model.Protection),
	})
	return req, nil
}

func outputResult(p *print.Printer, outputFormat, projectLabel string, resp *kms.WrappingKey) error {
	if resp == nil {
		return fmt.Errorf("response is nil")
	}

	switch outputFormat {
	case print.JSONOutputFormat:
		details, err := json.MarshalIndent(resp, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal KMS wrapping key: %w", err)
		}
		p.Outputln(string(details))

	case print.YAMLOutputFormat:
		details, err := yaml.MarshalWithOptions(resp, yaml.IndentSequence(true), yaml.UseJSONMarshaler())
		if err != nil {
			return fmt.Errorf("marshal KMS wrapping key: %w", err)
		}
		p.Outputln(string(details))

	default:
		p.Outputf("Created wrapping key for project %q. wrapping key ID: %s\n", projectLabel, utils.PtrString(resp.Id))
	}

	return nil
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), keyRingIdFlag, "ID of the KMS key ring")
	cmd.Flags().String(algorithmFlag, "", "En-/Decryption algorithm")
	cmd.Flags().String(displayNameFlag, "", "The display name to distinguish multiple wrapping keys")
	cmd.Flags().String(descriptionFlag, "", "Optional description of the wrapping key")
	cmd.Flags().String(purposeFlag, "", "Purpose of the wrapping key. Enum: 'wrap_symmetric_key', 'wrap_asymmetric_key' ")

	// backend was deprectaed in /v1beta, but protection is a required attribute with value "software"
	cmd.Flags().String(protectionFlag, "", "Protection of the wrapping key. Value: 'software' ")

	err := flags.MarkFlagsRequired(cmd, keyRingIdFlag, algorithmFlag, purposeFlag, displayNameFlag, protectionFlag)
	cobra.CheckErr(err)
}
