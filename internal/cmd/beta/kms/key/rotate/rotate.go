package rotate

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/goccy/go-yaml"
	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
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
	keyRingIdFlag = "key-ring"
	keyIdFlag     = "key"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	KeyRingId string
	KeyId     string
}

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rotate",
		Short: "Rotate a key",
		Long:  "Rotates the given key.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Rotate a KMS Key "my-key-id" and increase it's version inside the Key Ring "my-key-ring-id".`,
				`$ stackit beta kms keyring rotate --key-ring "my-key-ring-id" --key "my-key-id"`),
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
				params.Printer.Debug(print.ErrorLevel, "get key name: %v", err)
				keyName = model.KeyId
			}

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to rotate the key %q? (This cannot be undone)", keyName)
				err = params.Printer.PromptForConfirmation(prompt)
				if err != nil {
					return err
				}
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("rotate KMS Key: %w", err)
			}

			return outputResult(params.Printer, model.OutputFormat, resp)
		},
	}

	configureFlags(cmd)
	return cmd
}

func parseInput(p *print.Printer, cmd *cobra.Command) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	keyRingId := flags.FlagToStringValue(p, cmd, keyRingIdFlag)
	keyId := flags.FlagToStringValue(p, cmd, keyIdFlag)

	// Validate the uuid format of the IDs
	errKeyRing := utils.ValidateUUID(keyRingId)
	errKey := utils.ValidateUUID(keyId)
	if errKeyRing != nil || errKey != nil {
		return nil, &errors.DSAInputPlanError{
			Cmd: cmd,
		}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		KeyRingId:       keyRingId,
		KeyId:           keyId,
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

func buildRequest(ctx context.Context, model *inputModel, apiClient *kms.APIClient) kms.ApiRotateKeyRequest {
	req := apiClient.RotateKey(ctx, model.ProjectId, model.Region, model.KeyRingId, model.KeyId)
	return req
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), keyRingIdFlag, "ID of the KMS Key Ring where the Key is stored")
	cmd.Flags().Var(flags.UUIDFlag(), keyIdFlag, "ID of the actual Key")
	err := flags.MarkFlagsRequired(cmd, keyRingIdFlag, keyIdFlag)
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
			return fmt.Errorf("marshal KMS Key Version: %w", err)
		}
		p.Outputln(string(details))
		return nil

	case print.YAMLOutputFormat:
		details, err := yaml.MarshalWithOptions(resp, yaml.IndentSequence(true), yaml.UseJSONMarshaler())
		if err != nil {
			return fmt.Errorf("marshal KMS Key Version: %w", err)
		}
		p.Outputln(string(details))
		return nil

	default:
		p.Outputf("Rotated Key %s\n", utils.PtrString(resp.KeyId))
		return nil
	}
}
