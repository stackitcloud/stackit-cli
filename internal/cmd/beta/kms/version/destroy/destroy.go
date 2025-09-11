package destroy

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
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/kms/client"
	kmsUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/kms/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/kms"
	"gopkg.in/yaml.v2"
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
		Use:   "destroy",
		Short: "Destroy a key version",
		Long:  "Removes the key material of a version.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Destroy key version "0" for the key "my-key-id" inside the key ring "my-key-ring-id"`,
				`$ stackit beta kms version destroy --key "my-key-id" --key-ring "my-key-ring-id" --version 0`),
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
			// This operation can be undone. Don't ask for confirmation!

			// Call API
			req := buildRequest(ctx, model, apiClient)
			err = req.Execute()
			if err != nil {
				return fmt.Errorf("destroy key Version: %w", err)
			}

			return outputResult(params.Printer, model.OutputFormat, *model.VersionNumber, model.KeyId, keyName)
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

func buildRequest(ctx context.Context, model *inputModel, apiClient *kms.APIClient) kms.ApiDestroyVersionRequest {
	return apiClient.DestroyVersion(ctx, model.ProjectId, model.Region, model.KeyRingId, model.KeyId, *model.VersionNumber)
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), keyRingIdFlag, "ID of the KMS key ring")
	cmd.Flags().Var(flags.UUIDFlag(), keyIdFlag, "ID of the key")
	cmd.Flags().Int64(versionNumberFlag, 0, "Version number of the key")

	err := flags.MarkFlagsRequired(cmd, keyRingIdFlag, keyIdFlag, versionNumberFlag)
	cobra.CheckErr(err)
}

func outputResult(p *print.Printer, outputFormat string, versionNumber int64, keyId, keyName string) error {
	switch outputFormat {
	case print.JSONOutputFormat:
		details := struct {
			KeyId         string `json:"keyId"`
			KeyName       string `json:"keyName"`
			VersionNumber int64  `json:"versionNumber"`
			Status        string `json:"status"`
		}{
			KeyId:         keyId,
			KeyName:       keyName,
			VersionNumber: versionNumber,
			Status:        fmt.Sprintf("Destroyed version %d of key '%s'.", versionNumber, keyName),
		}
		b, err := json.MarshalIndent(details, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal output to JSON: %w", err)
		}
		p.Outputln(string(b))
		return nil

	case print.YAMLOutputFormat:
		details := struct {
			KeyId         string `yaml:"keyId"`
			KeyName       string `yaml:"keyName"`
			VersionNumber int64  `yaml:"versionNumber"`
			Status        string `yaml:"status"`
		}{
			KeyId:         keyId,
			KeyName:       keyName,
			VersionNumber: versionNumber,
			Status:        fmt.Sprintf("Destroyed version %d of key '%s'.", versionNumber, keyName),
		}
		b, err := yaml.Marshal(details)
		if err != nil {
			return fmt.Errorf("marshal output to YAML: %w", err)
		}
		p.Outputln(string(b))
		return nil

	default:
		p.Outputf("Destroyed version %d of key '%q'\n", versionNumber, keyName)
		return nil
	}
}
