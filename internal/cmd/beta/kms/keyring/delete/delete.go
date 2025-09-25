package delete

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
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	kmsUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/kms/utils"

	"github.com/stackitcloud/stackit-cli/internal/pkg/services/kms/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/kms"
)

const (
	keyRingIdArg = "KEYRING_ID"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	KeyRingId string
}

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("delete %s", keyRingIdArg),
		Short: "Deletes a KMS key ring",
		Long:  "Deletes a KMS key ring.",
		Args:  args.SingleArg(keyRingIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Delete a KMS key ring with ID "xxx"`,
				"$ stackit beta kms keyring delete xxx"),
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

			keyRingLabel, err := kmsUtils.GetKeyRingName(ctx, apiClient, model.ProjectId, model.KeyRingId, model.Region)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get key ring name: %v", err)
				keyRingLabel = model.KeyRingId
			}

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to delete key ring %q? (This cannot be undone)", keyRingLabel)
				err = params.Printer.PromptForConfirmation(prompt)
				if err != nil {
					return err
				}
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			err = req.Execute()
			if err != nil {
				return fmt.Errorf("delete KMS key ring: %w", err)
			}

			// Wait for async operation not relevant. Key ring deletion is synchronous.

			return outputResult(params.Printer, model.OutputFormat, keyRingLabel)
		},
	}
	return cmd
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	keyRingId := inputArgs[0]

	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		KeyRingId:       keyRingId,
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

func buildRequest(ctx context.Context, model *inputModel, apiClient *kms.APIClient) kms.ApiDeleteKeyRingRequest {
	req := apiClient.DeleteKeyRing(ctx, model.ProjectId, model.Region, model.KeyRingId)
	return req
}

func outputResult(p *print.Printer, outputFormat, keyRingLabel string) error {
	switch outputFormat {
	case print.JSONOutputFormat:
		details := struct {
			KeyRingLabel string `json:"keyRingLabel"`
			Status       string `json:"status"`
		}{
			KeyRingLabel: keyRingLabel,
			Status:       "Key ring deleted.",
		}
		b, err := json.MarshalIndent(details, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal output to JSON: %w", err)
		}
		p.Outputln(string(b))
		return nil

	case print.YAMLOutputFormat:
		details := struct {
			KeyRingLabel string `yaml:"keyRingLabel"`
			Status       string `yaml:"status"`
		}{
			KeyRingLabel: keyRingLabel,
			Status:       "Key ring deleted.",
		}
		b, err := yaml.Marshal(details)
		if err != nil {
			return fmt.Errorf("marshal output to YAML: %w", err)
		}
		p.Outputln(string(b))
		return nil

	default:
		p.Outputf("Deleted key ring: %s\n", keyRingLabel)
		return nil
	}
}
