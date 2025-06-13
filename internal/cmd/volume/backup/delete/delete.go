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
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/iaas/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/spinner"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/stackitcloud/stackit-sdk-go/services/iaas"
	"github.com/stackitcloud/stackit-sdk-go/services/iaas/wait"
)

const (
	backupIdArg = "BACKUP_ID"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	BackupId string
}

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("delete %s", backupIdArg),
		Short: "Deletes a backup",
		Long:  "Deletes a backup by its ID.",
		Args:  args.SingleArg(backupIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Delete a backup with ID "xxx"`, "$ stackit volume backup delete xxx"),
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

			backup, err := apiClient.GetBackup(ctx, model.ProjectId, model.BackupId).Execute()
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get backup name: %v", err)
			}
			backupLabel := model.BackupId
			if backup != nil && backup.Name != nil {
				backupLabel = *backup.Name
			}

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to delete backup %q? (This cannot be undone)", backupLabel)
				err = params.Printer.PromptForConfirmation(prompt)
				if err != nil {
					return err
				}
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			err = req.Execute()
			if err != nil {
				return fmt.Errorf("delete backup: %w", err)
			}

			// Wait for async operation, if async mode not enabled
			if !model.Async {
				s := spinner.New(params.Printer)
				s.Start("Deleting backup")
				_, err = wait.DeleteBackupWaitHandler(ctx, apiClient, model.ProjectId, model.BackupId).WaitWithContext(ctx)
				if err != nil {
					return fmt.Errorf("wait for backup deletion: %w", err)
				}
				s.Stop()
			}

			if model.Async {
				params.Printer.Info("Triggered deletion of backup %q\n", backupLabel)
			} else {
				params.Printer.Info("Deleted backup %q\n", backupLabel)
			}
			return nil
		},
	}
	return cmd
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	backupId := inputArgs[0]

	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		BackupId:        backupId,
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

func buildRequest(ctx context.Context, model *inputModel, apiClient *iaas.APIClient) iaas.ApiDeleteBackupRequest {
	req := apiClient.DeleteBackup(ctx, model.ProjectId, model.BackupId)
	return req
}
