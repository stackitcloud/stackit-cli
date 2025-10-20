package create

import (
	"context"
	"fmt"

	iaasClient "github.com/stackitcloud/stackit-cli/internal/pkg/services/iaas/client"

	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	cliErr "github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	iaasUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/iaas/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/serverbackup/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/serverbackup"
)

const (
	backupNameFlag            = "name"
	backupRetentionPeriodFlag = "retention-period"
	backupVolumeIdsFlag       = "volume-ids"
	serverIdFlag              = "server-id"

	defaultRetentionPeriod = 14
)

type inputModel struct {
	*globalflags.GlobalFlagModel

	ServerId              string
	BackupName            string
	BackupRetentionPeriod int64
	BackupVolumeIds       []string
}

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Creates a Server Backup.",
		Long:  "Creates a Server Backup. Operation always is async.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Create a Server Backup with name "mybackup"`,
				`$ stackit server backup create --server-id xxx --name=mybackup`),
			examples.NewExample(
				`Create a Server Backup with name "mybackup" and retention period of 5 days`,
				`$ stackit server backup create --server-id xxx --name=mybackup --retention-period=5`),
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

			serverLabel := model.ServerId
			// Get server name
			if iaasApiClient, err := iaasClient.ConfigureClient(params.Printer, params.CliVersion); err == nil {
				serverName, err := iaasUtils.GetServerName(ctx, iaasApiClient, model.ProjectId, model.ServerId)
				if err != nil {
					params.Printer.Debug(print.ErrorLevel, "get server name: %v", err)
				} else if serverName != "" {
					serverLabel = serverName
				}
			}

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to create a Backup for server %s?", model.ServerId)
				err = params.Printer.PromptForConfirmation(prompt)
				if err != nil {
					return err
				}
			}

			// Call API
			req, err := buildRequest(ctx, model, apiClient)
			if err != nil {
				return err
			}
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("create Server Backup: %w", err)
			}

			return outputResult(params.Printer, model.OutputFormat, serverLabel, *resp)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().VarP(flags.UUIDFlag(), serverIdFlag, "s", "Server ID")
	cmd.Flags().StringP(backupNameFlag, "b", "", "Backup name")
	cmd.Flags().Int64P(backupRetentionPeriodFlag, "d", defaultRetentionPeriod, "Backup retention period (in days)")
	cmd.Flags().VarP(flags.UUIDSliceFlag(), backupVolumeIdsFlag, "i", "Backup volume IDs, as comma separated UUID values.")

	err := flags.MarkFlagsRequired(cmd, serverIdFlag, backupNameFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &cliErr.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel:       globalFlags,
		ServerId:              flags.FlagToStringValue(p, cmd, serverIdFlag),
		BackupRetentionPeriod: flags.FlagWithDefaultToInt64Value(p, cmd, backupRetentionPeriodFlag),
		BackupName:            flags.FlagToStringValue(p, cmd, backupNameFlag),
		BackupVolumeIds:       flags.FlagToStringSliceValue(p, cmd, backupVolumeIdsFlag),
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *serverbackup.APIClient) (serverbackup.ApiCreateBackupRequest, error) {
	req := apiClient.CreateBackup(ctx, model.ProjectId, model.ServerId, model.Region)
	payload := serverbackup.CreateBackupPayload{
		Name:            &model.BackupName,
		RetentionPeriod: &model.BackupRetentionPeriod,
		VolumeIds:       &model.BackupVolumeIds,
	}
	if model.BackupVolumeIds == nil {
		payload.VolumeIds = nil
	}
	req = req.CreateBackupPayload(payload)
	return req, nil
}

func outputResult(p *print.Printer, outputFormat, serverLabel string, resp serverbackup.BackupJob) error {
	return p.OutputResult(outputFormat, resp, func() error {
		p.Outputf("Triggered creation of server backup for server %s. Backup ID: %s\n", serverLabel, utils.PtrString(resp.Id))
		return nil
	})
}
