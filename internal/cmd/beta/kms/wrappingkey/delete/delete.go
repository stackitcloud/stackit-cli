package delete

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
	keyRingIdFlag     = "key-ring"
	wrappingKeyIdFlag = "wrapping-key"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	KeyRingId   string
	WrappingKey string
}

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Deletes a KMS wrapping key",
		Long:  "Deletes a KMS wrapping key inside a specific key ring.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Delete a KMS wrapping key "my-wrapping-key-id" inside the key ring "my-key-ring-id"`,
				`$ stackit beta kms keyring delete --key-ring "my-key-ring-id" --wrapping-key "my-wrapping-key-id"`),
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

			wrappingKeyName, err := kmsUtils.GetWrappingKeyName(ctx, apiClient, model.ProjectId, model.Region, model.KeyRingId, model.WrappingKey)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get wrapping key name: %v", err)
				wrappingKeyName = model.WrappingKey
			}

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to delete key ring %q? (This cannot be undone)", wrappingKeyName)
				err = params.Printer.PromptForConfirmation(prompt)
				if err != nil {
					return err
				}
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			err = req.Execute()
			if err != nil {
				return fmt.Errorf("delete KMS wrapping key: %w", err)
			}

			// Wait for async operation not relevant. Wrapping key deletion is synchronous
			return outputResult(params.Printer, model.OutputFormat, wrappingKeyName)
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
		WrappingKey:     flags.FlagToStringValue(p, cmd, wrappingKeyIdFlag),
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

func buildRequest(ctx context.Context, model *inputModel, apiClient *kms.APIClient) kms.ApiDeleteWrappingKeyRequest {
	req := apiClient.DeleteWrappingKey(ctx, model.ProjectId, model.Region, model.KeyRingId, model.WrappingKey)
	return req
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), keyRingIdFlag, "ID of the KMS key ring where the wrapping key is stored")
	cmd.Flags().Var(flags.UUIDFlag(), wrappingKeyIdFlag, "ID of the actual wrapping key")
	err := flags.MarkFlagsRequired(cmd, keyRingIdFlag, wrappingKeyIdFlag)
	cobra.CheckErr(err)
}

func outputResult(p *print.Printer, outputFormat, wrappingKeyName string) error {
	switch outputFormat {
	case print.JSONOutputFormat:
		details := struct {
			WrappingKeyName string `json:"wrappingKeyName"`
			Status          string `json:"status"`
		}{
			WrappingKeyName: wrappingKeyName,
			Status:          "Deleted wrapping key.",
		}
		b, err := json.MarshalIndent(details, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal output to JSON: %w", err)
		}
		p.Outputln(string(b))
		return nil

	case print.YAMLOutputFormat:
		details := struct {
			WrappingKeyName string `yaml:"wrappingKeyName"`
			Status          string `yaml:"status"`
		}{
			WrappingKeyName: wrappingKeyName,
			Status:          "Deleted wrapping key.",
		}
		b, err := yaml.Marshal(details)
		if err != nil {
			return fmt.Errorf("marshal output to YAML: %w", err)
		}
		p.Outputln(string(b))
		return nil

	default:
		p.Outputf("Deleted wrapping key: %s\n", wrappingKeyName)
		return nil
	}
}
