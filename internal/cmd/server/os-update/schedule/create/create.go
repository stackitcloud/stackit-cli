package create

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/types"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	cliErr "github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	iaasClient "github.com/stackitcloud/stackit-cli/internal/pkg/services/iaas/client"
	iaasUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/iaas/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/serverosupdate/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/serverupdate"
)

const (
	nameFlag              = "name"
	enabledFlag           = "enabled"
	rruleFlag             = "rrule"
	maintenanceWindowFlag = "maintenance-window"
	serverIdFlag          = "server-id"

	defaultRrule             = "DTSTART;TZID=Europe/Sofia:20200803T023000 RRULE:FREQ=DAILY;INTERVAL=1"
	defaultMaintenanceWindow = 23
	defaultEnabled           = true
)

type inputModel struct {
	*globalflags.GlobalFlagModel

	ServerId          string
	ScheduleName      string
	Enabled           bool
	Rrule             string
	MaintenanceWindow int64
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Creates a Server os-update Schedule",
		Long:  "Creates a Server os-update Schedule.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Create a Server os-update Schedule with name "myschedule"`,
				`$ stackit server os-update schedule create --server-id xxx --name=myschedule`),
			examples.NewExample(
				`Create a Server os-update Schedule with name "myschedule" and maintenance window for 14 o'clock`,
				`$ stackit server os-update schedule create --server-id xxx --name=myschedule --maintenance-window=14`),
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

			serverLabel := model.ServerId
			// Get server name
			if iaasApiClient, err := iaasClient.ConfigureClient(params.Printer, params.CliVersion); err == nil {
				serverName, err := iaasUtils.GetServerName(ctx, iaasApiClient, model.ProjectId, model.Region, model.ServerId)
				if err != nil {
					params.Printer.Debug(print.ErrorLevel, "get server name: %v", err)
				} else if serverName != "" {
					serverLabel = serverName
				}
			}

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to create a os-update Schedule for server %s?", serverLabel)
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
				return fmt.Errorf("create Server os-update Schedule: %w", err)
			}

			return outputResult(params.Printer, model.OutputFormat, serverLabel, *resp)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().VarP(flags.UUIDFlag(), serverIdFlag, "s", "Server ID")
	cmd.Flags().StringP(nameFlag, "n", "", "os-update schedule name")
	cmd.Flags().Int64P(maintenanceWindowFlag, "d", defaultMaintenanceWindow, "os-update maintenance window (in hours, 1-24)")
	cmd.Flags().BoolP(enabledFlag, "e", defaultEnabled, "Is the server os-update schedule enabled")
	cmd.Flags().StringP(rruleFlag, "r", defaultRrule, "os-update RRULE (recurrence rule)")

	err := flags.MarkFlagsRequired(cmd, serverIdFlag, nameFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command, _ []string) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &cliErr.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel:   globalFlags,
		ServerId:          flags.FlagToStringValue(p, cmd, serverIdFlag),
		MaintenanceWindow: flags.FlagWithDefaultToInt64Value(p, cmd, maintenanceWindowFlag),
		ScheduleName:      flags.FlagToStringValue(p, cmd, nameFlag),
		Rrule:             flags.FlagWithDefaultToStringValue(p, cmd, rruleFlag),
		Enabled:           flags.FlagToBoolValue(p, cmd, enabledFlag),
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *serverupdate.APIClient) (serverupdate.ApiCreateUpdateScheduleRequest, error) {
	req := apiClient.CreateUpdateSchedule(ctx, model.ProjectId, model.ServerId, model.Region)
	req = req.CreateUpdateSchedulePayload(serverupdate.CreateUpdateSchedulePayload{
		Enabled:           &model.Enabled,
		Name:              &model.ScheduleName,
		Rrule:             &model.Rrule,
		MaintenanceWindow: &model.MaintenanceWindow,
	})
	return req, nil
}

func outputResult(p *print.Printer, outputFormat, serverLabel string, resp serverupdate.UpdateSchedule) error {
	return p.OutputResult(outputFormat, resp, func() error {
		p.Outputf("Created server os-update schedule for server %s. os-update Schedule ID: %s\n", serverLabel, utils.PtrString(resp.Id))
		return nil
	})
}
