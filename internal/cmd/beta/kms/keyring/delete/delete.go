package delete

import (
	"context"
	"fmt"

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
				`Delete a KMS key ring with ID "my-key-ring-id"`,
				`$ stackit beta kms keyring delete "my-key-ring-id"`),
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

			// Don't output anything. It's a deletion.
			params.Printer.Info("Deleted the key ring %q\n", keyRingLabel)
			return nil
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
