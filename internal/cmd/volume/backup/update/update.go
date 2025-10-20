package update

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
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/iaas/client"
	iaasutils "github.com/stackitcloud/stackit-cli/internal/pkg/services/iaas/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/stackitcloud/stackit-sdk-go/services/iaas"
)

const (
	backupIdArg = "BACKUP_ID"
	nameFlag    = "name"
	labelsFlag  = "labels"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	BackupId string
	Name     *string
	Labels   map[string]string
}

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("update %s", backupIdArg),
		Short: "Updates a backup",
		Long:  "Updates a backup by its ID.",
		Args:  args.SingleArg(backupIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Update the name of a backup with ID "xxx"`,
				"$ stackit volume backup update xxx --name new-name"),
			examples.NewExample(
				`Update the labels of a backup with ID "xxx"`,
				"$ stackit volume backup update xxx --labels key1=value1,key2=value2"),
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

			backupLabel, err := iaasutils.GetBackupName(ctx, apiClient, model.ProjectId, model.BackupId)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get backup name: %v", err)
			}

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to update backup %q? (This cannot be undone)", model.BackupId)
				err = params.Printer.PromptForConfirmation(prompt)
				if err != nil {
					return err
				}
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("update backup: %w", err)
			}

			return outputResult(params.Printer, model.OutputFormat, backupLabel, resp)
		},
	}

	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().String(nameFlag, "", "Name of the backup")
	cmd.Flags().StringToString(labelsFlag, nil, "Key-value string pairs as labels")
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	backupId := inputArgs[0]

	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	name := flags.FlagToStringPointer(p, cmd, nameFlag)
	labels := flags.FlagToStringToStringPointer(p, cmd, labelsFlag)
	if labels == nil {
		labels = &map[string]string{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		BackupId:        backupId,
		Name:            name,
		Labels:          *labels,
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *iaas.APIClient) iaas.ApiUpdateBackupRequest {
	req := apiClient.UpdateBackup(ctx, model.ProjectId, model.BackupId)

	payload := iaas.UpdateBackupPayload{
		Name:   model.Name,
		Labels: utils.ConvertStringMapToInterfaceMap(utils.Ptr(model.Labels)),
	}

	req = req.UpdateBackupPayload(payload)
	return req
}

func outputResult(p *print.Printer, outputFormat, backupLabel string, backup *iaas.Backup) error {
	if backup == nil {
		return fmt.Errorf("backup response is empty")
	}

	return p.OutputResult(outputFormat, backup, func() error {
		p.Outputf("Updated backup %q\n", backupLabel)
		return nil
	})
}
