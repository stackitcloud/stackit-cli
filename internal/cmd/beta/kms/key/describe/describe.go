package describe

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
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/kms"
)

const (
	argKeyID      = "KEY_ID"
	flagKeyRingID = "keyring-id"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	KeyID     string
	KeyRingID string
}

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("describe %s", argKeyID),
		Short: "Describe a KMS key",
		Long:  "Describe a KMS key",
		Args:  args.SingleArg(argKeyID, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Describe a KMS key with ID xxx of keyring yyy`,
				`$ stackit beta kms key describe xxx --keyring-id yyy`,
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
				return fmt.Errorf("get key: %w", err)
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

func parseInput(p *print.Printer, cmd *cobra.Command, args []string) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}
	model := &inputModel{
		GlobalFlagModel: globalFlags,
		KeyID:           args[0],
		KeyRingID:       flags.FlagToStringValue(p, cmd, flagKeyRingID),
	}
	p.DebugInputModel(model)
	return model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *kms.APIClient) kms.ApiGetKeyRequest {
	return apiClient.GetKey(ctx, model.ProjectId, model.Region, model.KeyRingID, model.KeyID)
}

func outputResult(p *print.Printer, outputFormat string, key *kms.Key) error {
	if key == nil {
		return fmt.Errorf("key response is empty")
	}
	return p.OutputResult(outputFormat, key, func() error {
		table := tables.NewTable()
		table.AddRow("ID", utils.PtrString(key.Id))
		table.AddSeparator()
		table.AddRow("DISPLAY NAME", utils.PtrString(key.DisplayName))
		table.AddSeparator()
		table.AddRow("CREATED AT", utils.PtrString(key.CreatedAt))
		table.AddSeparator()
		table.AddRow("STATE", utils.PtrString(key.State))
		table.AddSeparator()
		table.AddRow("DESCRIPTION", utils.PtrString(key.Description))
		table.AddSeparator()
		table.AddRow("ACCESS SCOPE", utils.PtrString(key.AccessScope))
		table.AddSeparator()
		table.AddRow("ALGORITHM", utils.PtrString(key.Algorithm))
		table.AddSeparator()
		table.AddRow("DELETION DATE", utils.PtrString(key.DeletionDate))
		table.AddSeparator()
		table.AddRow("IMPORT ONLY", utils.PtrString(key.ImportOnly))
		table.AddSeparator()
		table.AddRow("KEYRING ID", utils.PtrString(key.KeyRingId))
		table.AddSeparator()
		table.AddRow("PROTECTION", utils.PtrString(key.Protection))
		table.AddSeparator()
		table.AddRow("PURPOSE", utils.PtrString(key.Purpose))

		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("display table: %w", err)
		}
		return nil
	})
}
