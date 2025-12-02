package rotate

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/types"

	"github.com/goccy/go-yaml"
	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	kmsUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/kms/utils"

	"github.com/stackitcloud/stackit-cli/internal/pkg/services/kms/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/kms"
)

const (
	keyIdArg = "KEY_ID"

	keyRingIdFlag = "keyring-id"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	KeyId     string
	KeyRingId string
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("rotate %s", keyIdArg),
		Short: "Rotate a key",
		Long:  "Rotates the given key.",
		Args:  args.SingleArg(keyIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Rotate a KMS key "MY_KEY_ID" and increase its version inside the key ring "my-keyring-id".`,
				`$ stackit beta kms key rotate "MY_KEY_ID" --keyring-id "my-keyring-id"`),
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

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to rotate the key %q? (this cannot be undone)", keyName)
				err = params.Printer.PromptForConfirmation(prompt)
				if err != nil {
					return err
				}
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("rotate KMS key: %w", err)
			}

			return outputResult(params.Printer, model.OutputFormat, resp)
		},
	}

	configureFlags(cmd)
	return cmd
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	keyId := inputArgs[0]

	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		KeyRingId:       flags.FlagToStringValue(p, cmd, keyRingIdFlag),
		KeyId:           keyId,
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *kms.APIClient) kms.ApiRotateKeyRequest {
	req := apiClient.RotateKey(ctx, model.ProjectId, model.Region, model.KeyRingId, model.KeyId)
	return req
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), keyRingIdFlag, "ID of the KMS key ring where the key is stored")

	err := flags.MarkFlagsRequired(cmd, keyRingIdFlag)
	cobra.CheckErr(err)
}

func outputResult(p *print.Printer, outputFormat string, resp *kms.Version) error {
	if resp == nil {
		return fmt.Errorf("response is nil")
	}

	switch outputFormat {
	case print.JSONOutputFormat:
		details, err := json.MarshalIndent(resp, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal KMS key version: %w", err)
		}
		p.Outputln(string(details))

	case print.YAMLOutputFormat:
		details, err := yaml.MarshalWithOptions(resp, yaml.IndentSequence(true), yaml.UseJSONMarshaler())
		if err != nil {
			return fmt.Errorf("marshal KMS key version: %w", err)
		}
		p.Outputln(string(details))

	default:
		p.Outputf("Rotated key %s\n", utils.PtrString(resp.KeyId))
	}

	return nil
}
