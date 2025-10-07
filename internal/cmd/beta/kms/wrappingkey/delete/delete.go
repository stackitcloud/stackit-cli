package delete

import (
	"context"
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
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/stackitcloud/stackit-cli/internal/pkg/services/kms/client"
	"github.com/stackitcloud/stackit-sdk-go/services/kms"
)

const (
	wrappingKeyIdArg = "WRAPPING_KEY_ID"

	keyRingIdFlag = "key-ring-id"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	WrappingKeyId string
	KeyRingId     string
}

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("delete %s", wrappingKeyIdArg),
		Short: "Deletes a KMS wrapping key",
		Long:  "Deletes a KMS wrapping key inside a specific key ring.",
		Args:  args.SingleArg(wrappingKeyIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Delete a KMS wrapping key "my-wrapping-key-id" inside the key ring "my-key-ring-id"`,
				`$ stackit beta kms keyring delete "my-wrapping-key-id" --key-ring "my-key-ring-id"`),
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

			wrappingKeyName, err := kmsUtils.GetWrappingKeyName(ctx, apiClient, model.ProjectId, model.Region, model.KeyRingId, model.WrappingKeyId)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get wrapping key name: %v", err)
				wrappingKeyName = model.WrappingKeyId
			}

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to delete the wrapping key %q? (This cannot be undone)", wrappingKeyName)
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

			// Don't output anything. It's a deletion.
			params.Printer.Info("Deleted wrapping key %q\n", wrappingKeyName)
			return nil
		},
	}

	configureFlags(cmd)
	return cmd
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	wrappingKeyId := inputArgs[0]

	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		KeyRingId:       flags.FlagToStringValue(p, cmd, keyRingIdFlag),
		WrappingKeyId:   wrappingKeyId,
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
	req := apiClient.DeleteWrappingKey(ctx, model.ProjectId, model.Region, model.KeyRingId, model.WrappingKeyId)
	return req
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), keyRingIdFlag, "ID of the KMS key ring where the wrapping key is stored")
	err := flags.MarkFlagsRequired(cmd, keyRingIdFlag)
	cobra.CheckErr(err)
}
