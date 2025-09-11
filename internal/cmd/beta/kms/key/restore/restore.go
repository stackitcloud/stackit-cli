package restore

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	kmsUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/kms/utils"
	"gopkg.in/yaml.v2"

	"github.com/stackitcloud/stackit-cli/internal/pkg/services/kms/client"
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
		Use:   "restore",
		Short: "Resotre a key",
		Long:  "Restores the given key from being deleted.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Restore a KMS key "my-key-id" inside the key ring "my-key-ring-id" that was scheduled for deletion.`,
				`$ stackit beta kms keyring restore --key-ring "my-key-ring-id" --key "my-key-id"`),
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
				prompt := fmt.Sprintf("Are you sure you want to restore key %q? (This cannot be undone)", keyName)
				err = params.Printer.PromptForConfirmation(prompt)
				if err != nil {
					return err
				}
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			err = req.Execute()
			if err != nil {
				return fmt.Errorf("restore KMS key: %w", err)
			}

			return outputResult(params.Printer, model.OutputFormat, model.KeyId, keyName)
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

	model := inputModel{
		GlobalFlagModel: globalFlags,
		KeyRingId:       flags.FlagToStringValue(p, cmd, keyRingIdFlag),
		KeyId:           flags.FlagToStringValue(p, cmd, keyIdFlag),
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

func buildRequest(ctx context.Context, model *inputModel, apiClient *kms.APIClient) kms.ApiRestoreKeyRequest {
	req := apiClient.RestoreKey(ctx, model.ProjectId, model.Region, model.KeyRingId, model.KeyId)
	return req
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), keyRingIdFlag, "ID of the KMS Key Ring where the Key is stored")
	cmd.Flags().Var(flags.UUIDFlag(), keyIdFlag, "ID of the actual Key")
	err := flags.MarkFlagsRequired(cmd, keyRingIdFlag, keyIdFlag)
	cobra.CheckErr(err)
}

func outputResult(p *print.Printer, outputFormat, keyId, keyName string) error {
	switch outputFormat {
	case print.JSONOutputFormat:
		details := struct {
			KeyId   string `json:"keyId"`
			KeyName string `json:"keyName"`
			Status  string `json:"status"`
		}{
			KeyId:   keyId,
			KeyName: keyName,
			Status:  "key restored",
		}
		b, err := json.MarshalIndent(details, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal output to JSON: %w", err)
		}
		p.Outputln(string(b))
		return nil

	case print.YAMLOutputFormat:
		details := struct {
			KeyId   string `yaml:"keyId"`
			KeyName string `yaml:"keyName"`
			Status  string `yaml:"status"`
		}{
			KeyId:   keyId,
			KeyName: keyName,
			Status:  "key restored",
		}
		b, err := yaml.Marshal(details)
		if err != nil {
			return fmt.Errorf("marshal output to YAML: %w", err)
		}
		p.Outputln(string(b))
		return nil

	default:
		p.Outputf("Successfully restored KMS key %q\n", keyName)
		return nil
	}
}
