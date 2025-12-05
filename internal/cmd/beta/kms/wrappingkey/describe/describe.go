package describe

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/types"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/kms/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/kms"
)

const (
	argWrappingKeyID = "WRAPPING_KEY_ID"
	flagKeyRingID    = "keyring-id"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	WrappingKeyID string
	KeyRingID     string
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("describe %s", argWrappingKeyID),
		Short: "Describe a KMS wrapping key",
		Long:  "Describe a KMS wrapping key",
		Args:  args.SingleArg(argWrappingKeyID, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Describe a KMS wrapping key with ID xxx of keyring yyy`,
				`$ stackit beta kms wrappingkey describe xxx --keyring-id yyy`,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()

			model, err := parseInput(params.Printer, cmd, args)
			if err != nil {
				return err
			}

			apiClient, err := client.ConfigureClient(params.Printer, params.CliVersion)
			if err != nil {
				return err
			}

			req := buildRequest(ctx, model, apiClient)

			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("get wrapping key: %w", err)
			}

			return outputResult(params.Printer, model.OutputFormat, resp)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), flagKeyRingID, "Key Ring ID")
	err := flags.MarkFlagsRequired(cmd, flagKeyRingID)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}
	model := &inputModel{
		GlobalFlagModel: globalFlags,
		WrappingKeyID:   inputArgs[0],
		KeyRingID:       flags.FlagToStringValue(p, cmd, flagKeyRingID),
	}
	p.DebugInputModel(model)
	return model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *kms.APIClient) kms.ApiGetWrappingKeyRequest {
	return apiClient.GetWrappingKey(ctx, model.ProjectId, model.Region, model.KeyRingID, model.WrappingKeyID)
}

func outputResult(p *print.Printer, outputFormat string, wrappingKey *kms.WrappingKey) error {
	if wrappingKey == nil {
		return fmt.Errorf("wrapping key response is empty")
	}
	return p.OutputResult(outputFormat, wrappingKey, func() error {
		table := tables.NewTable()
		table.AddRow("ID", utils.PtrString(wrappingKey.Id))
		table.AddSeparator()
		table.AddRow("DISPLAY NAME", utils.PtrString(wrappingKey.DisplayName))
		table.AddSeparator()
		table.AddRow("CREATED AT", utils.PtrString(wrappingKey.CreatedAt))
		table.AddSeparator()
		table.AddRow("STATE", utils.PtrString(wrappingKey.State))
		table.AddSeparator()
		table.AddRow("DESCRIPTION", utils.PtrString(wrappingKey.Description))
		table.AddSeparator()
		table.AddRow("ACCESS SCOPE", utils.PtrString(wrappingKey.AccessScope))
		table.AddSeparator()
		table.AddRow("ALGORITHM", utils.PtrString(wrappingKey.Algorithm))
		table.AddSeparator()
		table.AddRow("EXPIRES AT", utils.PtrString(wrappingKey.ExpiresAt))
		table.AddSeparator()
		table.AddRow("KEYRING ID", utils.PtrString(wrappingKey.KeyRingId))
		table.AddSeparator()
		table.AddRow("PROTECTION", utils.PtrString(wrappingKey.Protection))
		table.AddSeparator()
		table.AddRow("PUBLIC KEY", utils.PtrString(wrappingKey.PublicKey))
		table.AddSeparator()
		table.AddRow("PURPOSE", utils.PtrString(wrappingKey.Purpose))

		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("display table: %w", err)
		}
		return nil
	})
}
