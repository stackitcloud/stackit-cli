package disable

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
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/kms/client"
	kmsUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/kms/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/spinner"
	"github.com/stackitcloud/stackit-sdk-go/services/kms"
	"github.com/stackitcloud/stackit-sdk-go/services/kms/wait"
)

const (
	keyRingIdFlag     = "key-ring"
	keyIdFlag         = "key"
	versionNumberFlag = "version"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	KeyRingId     string
	KeyId         string
	VersionNumber *int64
}

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "disable",
		Short: "Disable a Key Versions",
		Long:  "Disable the given key version.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Disable key version "0" for the key "my-key-id" inside the key ring "my-key-ring-id"`,
				`$ stackit beta kms version disable --key "my-key-id" --key-ring "my-key-ring-id" --version 0`),
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

			// This operatio can be undone. Don't ask for confirmation!

			// Call API
			req := buildRequest(ctx, model, apiClient)
			err = req.Execute()
			if err != nil {
				return fmt.Errorf("disable Key Version: %w", err)
			}

			// Wait for async operation, if async mode not enabled
			if !model.Async {
				s := spinner.New(params.Printer)
				s.Start("Disableing key version")
				_, err = wait.DisableKeyVersionWaitHandler(ctx, apiClient, model.ProjectId, model.Region, model.KeyRingId, model.KeyId, *model.VersionNumber).WaitWithContext(ctx)
				if err != nil {
					return fmt.Errorf("wait for Key Version to be disabled: %w", err)
				}
				s.Stop()
			}

			params.Printer.Info("Disabled version %d of Key %q\n", *model.VersionNumber, keyName)
			return nil
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
		VersionNumber:   flags.FlagToInt64Pointer(p, cmd, versionNumberFlag),
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

func buildRequest(ctx context.Context, model *inputModel, apiClient *kms.APIClient) kms.ApiDisableVersionRequest {
	return apiClient.DisableVersion(ctx, model.ProjectId, model.Region, model.KeyRingId, model.KeyId, *model.VersionNumber)
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), keyRingIdFlag, "ID of the KMS Key Ring")
	cmd.Flags().Var(flags.UUIDFlag(), keyIdFlag, "ID of the Key")
	cmd.Flags().Int64(versionNumberFlag, 0, "Version number of the key")

	err := flags.MarkFlagsRequired(cmd, keyRingIdFlag, keyIdFlag, versionNumberFlag)
	cobra.CheckErr(err)
}
