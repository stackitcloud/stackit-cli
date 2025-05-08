package delete

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	iaasClient "github.com/stackitcloud/stackit-cli/internal/pkg/services/iaas/client"
	iaasUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/iaas/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/serverbackup/client"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/serverbackup"
)

const (
	scheduleIdArg = "SCHEDULE_ID"
	serverIdFlag  = "server-id"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	ScheduleId string
	ServerId   string
}

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("delete %s", scheduleIdArg),
		Short: "Deletes a Server Backup Schedule",
		Long:  "Deletes a Server Backup Schedule.",
		Args:  args.SingleArg(scheduleIdArg, nil),
		Example: examples.Build(
			examples.NewExample(
				`Delete a Server Backup Schedule with ID "xxx" for server "zzz"`,
				"$ stackit server backup schedule delete xxx --server-id=zzz"),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			model, err := parseInput(params.Printer, cmd, args)
			if err != nil {
				return err
			}

			// Configure API client
			apiClient, err := client.ConfigureClient(params.Printer)
			if err != nil {
				return err
			}

			serverLabel := model.ServerId
			// Get server name
			if iaasApiClient, err := iaasClient.ConfigureClient(params.Printer); err == nil {
				serverName, err := iaasUtils.GetServerName(ctx, iaasApiClient, model.ProjectId, model.ServerId)
				if err != nil {
					params.Printer.Debug(print.ErrorLevel, "get server name: %v", err)
				} else if serverName != "" {
					serverLabel = serverName
				}
			}

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to delete server backup schedule %q? (This cannot be undone)", serverLabel)
				err = params.Printer.PromptForConfirmation(prompt)
				if err != nil {
					return err
				}
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			err = req.Execute()
			if err != nil {
				return fmt.Errorf("delete Server Backup Schedule: %w", err)
			}

			params.Printer.Info("Deleted server backup schedule %q\n", model.ScheduleId)
			return nil
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().VarP(flags.UUIDFlag(), serverIdFlag, "s", "Server ID")

	err := flags.MarkFlagsRequired(cmd, serverIdFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	scheduleId := inputArgs[0]

	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		ScheduleId:      scheduleId,
		ServerId:        flags.FlagToStringValue(p, cmd, serverIdFlag),
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

func buildRequest(ctx context.Context, model *inputModel, apiClient *serverbackup.APIClient) serverbackup.ApiDeleteBackupScheduleRequest {
	req := apiClient.DeleteBackupSchedule(ctx, model.ProjectId, model.ServerId, model.Region, model.ScheduleId)
	return req
}
