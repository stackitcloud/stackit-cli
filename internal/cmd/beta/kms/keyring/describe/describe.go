package describe

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/types"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/kms/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/kms"
)

const (
	argKeyRingID = "KEYRING_ID"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	KeyRingID string
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("describe %s", argKeyRingID),
		Short: "Describe a KMS key ring",
		Long:  "Describe a KMS key ring",
		Args:  args.SingleArg(argKeyRingID, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Describe a KMS key ring with ID xxx`,
				`$ stackit beta kms keyring describe xxx`,
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
				return fmt.Errorf("get key ring: %w", err)
			}

			return outputResult(params.Printer, model.OutputFormat, resp)
		},
	}
	return cmd
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}
	model := &inputModel{
		GlobalFlagModel: globalFlags,
		KeyRingID:       inputArgs[0],
	}
	p.DebugInputModel(model)
	return model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *kms.APIClient) kms.ApiGetKeyRingRequest {
	return apiClient.GetKeyRing(ctx, model.ProjectId, model.Region, model.KeyRingID)
}

func outputResult(p *print.Printer, outputFormat string, keyRing *kms.KeyRing) error {
	if keyRing == nil {
		return fmt.Errorf("key ring response is empty")
	}
	return p.OutputResult(outputFormat, keyRing, func() error {
		table := tables.NewTable()
		table.AddRow("ID", utils.PtrString(keyRing.Id))
		table.AddSeparator()
		table.AddRow("DISPLAY NAME", utils.PtrString(keyRing.DisplayName))
		table.AddSeparator()
		table.AddRow("CREATED AT", utils.PtrString(keyRing.CreatedAt))
		table.AddSeparator()
		table.AddRow("STATE", utils.PtrString(keyRing.State))
		table.AddSeparator()
		table.AddRow("DESCRIPTION", utils.PtrString(keyRing.Description))

		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("display table: %w", err)
		}
		return nil
	})
}
