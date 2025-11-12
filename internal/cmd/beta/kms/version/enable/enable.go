package enable

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/goccy/go-yaml"
	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/kms/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/spinner"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/kms"
	"github.com/stackitcloud/stackit-sdk-go/services/kms/wait"
)

const (
	versionNumberArg = "VERSION_NUMBER"

	keyRingIdFlag = "keyring-id"
	keyIdFlag     = "key-id"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	KeyRingId     string
	KeyId         string
	VersionNumber int64
}

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("enable %s", versionNumberArg),
		Short: "Enable a key version",
		Long:  "Enable the given key version.",
		Args:  args.SingleArg(versionNumberArg, nil),
		Example: examples.Build(
			examples.NewExample(
				`Enable key version "42" for the key "my-key-id" inside the key ring "my-keyring-id"`,
				`$ stackit beta kms version enable 42 --key-id "my-key-id" --keyring-id "my-keyring-id"`),
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

			// This operation can be undone. Don't ask for confirmation!

			// Call API
			req := buildRequest(ctx, model, apiClient)
			err = req.Execute()
			if err != nil {
				return fmt.Errorf("enable key version: %w", err)
			}

			// Wait for async operation, if async mode not enabled
			if !model.Async {
				s := spinner.New(params.Printer)
				s.Start("Enabling key version")
				_, err = wait.EnableKeyVersionWaitHandler(ctx, apiClient, model.ProjectId, model.Region, model.KeyRingId, model.KeyId, model.VersionNumber).WaitWithContext(ctx)
				if err != nil {
					return fmt.Errorf("wait for key version to be enabled: %w", err)
				}
				s.Stop()
			}

			// Get the key version in its state afterwards
			resp, err := apiClient.GetVersionExecute(ctx, model.ProjectId, model.Region, model.KeyRingId, model.KeyId, model.VersionNumber)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get key version: %v", err)
			}

			return outputResult(params.Printer, model.OutputFormat, resp)
		},
	}

	configureFlags(cmd)
	return cmd
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	versionStr := inputArgs[0]
	versionNumber, err := strconv.ParseInt(versionStr, 10, 64)
	if err != nil || versionNumber < 0 {
		return nil, &errors.ArgValidationError{
			Arg:     versionNumberArg,
			Details: fmt.Sprintf("invalid value %q: must be a positive integer", versionStr),
		}
	}

	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		KeyRingId:       flags.FlagToStringValue(p, cmd, keyRingIdFlag),
		KeyId:           flags.FlagToStringValue(p, cmd, keyIdFlag),
		VersionNumber:   versionNumber,
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *kms.APIClient) kms.ApiEnableVersionRequest {
	return apiClient.EnableVersion(ctx, model.ProjectId, model.Region, model.KeyRingId, model.KeyId, model.VersionNumber)
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), keyRingIdFlag, "ID of the KMS key ring")
	cmd.Flags().Var(flags.UUIDFlag(), keyIdFlag, "ID of the key")

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
			return fmt.Errorf("marshal KMS key: %w", err)
		}
		p.Outputln(string(details))

	case print.YAMLOutputFormat:
		details, err := yaml.MarshalWithOptions(resp, yaml.IndentSequence(true), yaml.UseJSONMarshaler())
		if err != nil {
			return fmt.Errorf("marshal KMS key: %w", err)
		}
		p.Outputln(string(details))

	default:
		p.Outputf("Enabled version %d of the key %q\n", utils.PtrValue(resp.Number), utils.PtrValue(resp.KeyId))
	}

	return nil
}
