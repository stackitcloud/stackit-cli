package update

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/goccy/go-yaml"
	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	cliErr "github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/serverosupdate/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/serverupdate"
)

const (
	scheduleIdArg = "SCHEDULE_ID"

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
	ScheduleId        string
	ScheduleName      *string
	Enabled           *bool
	Rrule             *string
	MaintenanceWindow *int64
}

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("update %s", scheduleIdArg),
		Short: "Updates a Server os-update Schedule",
		Long:  "Updates a Server os-update Schedule.",
		Example: examples.Build(
			examples.NewExample(
				`Update the name of the os-update schedule "zzz" of server "xxx"`,
				"$ stackit beta server os-update schedule update zzz --server-id=xxx --name=newname"),
		),
		Args: args.SingleArg(scheduleIdArg, nil),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()

			model, err := parseInput(p, cmd, args)
			if err != nil {
				return err
			}

			// Configure API client
			apiClient, err := client.ConfigureClient(p)
			if err != nil {
				return err
			}

			currentSchedule, err := apiClient.GetUpdateScheduleExecute(ctx, model.ProjectId, model.ServerId, model.ScheduleId)
			if err != nil {
				p.Debug(print.ErrorLevel, "get current server os-update schedule: %v", err)
				return err
			}

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to update Server os-update Schedule %q?", model.ScheduleId)
				err = p.PromptForConfirmation(prompt)
				if err != nil {
					return err
				}
			}

			// Call API
			req, err := buildRequest(ctx, model, apiClient, *currentSchedule)
			if err != nil {
				return err
			}
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("update Server os-update Schedule: %w", err)
			}

			return outputResult(p, model.OutputFormat, *resp)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().VarP(flags.UUIDFlag(), serverIdFlag, "s", "Server ID")

	cmd.Flags().StringP(nameFlag, "n", "", "os-update schedule name")
	cmd.Flags().Int64P(maintenanceWindowFlag, "d", defaultMaintenanceWindow, "Maintenance window (in hours, 1-24)")
	cmd.Flags().BoolP(enabledFlag, "e", defaultEnabled, "Is the server os-update schedule enabled")
	cmd.Flags().StringP(rruleFlag, "r", defaultRrule, "os-update RRULE (recurrence rule)")

	err := flags.MarkFlagsRequired(cmd, serverIdFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	scheduleId := inputArgs[0]

	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &cliErr.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel:   globalFlags,
		ScheduleId:        scheduleId,
		ScheduleName:      flags.FlagToStringPointer(p, cmd, nameFlag),
		ServerId:          flags.FlagToStringValue(p, cmd, serverIdFlag),
		MaintenanceWindow: flags.FlagToInt64Pointer(p, cmd, maintenanceWindowFlag),
		Rrule:             flags.FlagToStringPointer(p, cmd, rruleFlag),
		Enabled:           flags.FlagToBoolPointer(p, cmd, enabledFlag),
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

func buildRequest(ctx context.Context, model *inputModel, apiClient *serverupdate.APIClient, old serverupdate.UpdateSchedule) (serverupdate.ApiUpdateUpdateScheduleRequest, error) {
	req := apiClient.UpdateUpdateSchedule(ctx, model.ProjectId, model.ServerId, model.ScheduleId)

	if model.MaintenanceWindow != nil {
		old.MaintenanceWindow = model.MaintenanceWindow
	}
	if model.Enabled != nil {
		old.Enabled = model.Enabled
	}
	if model.ScheduleName != nil {
		old.Name = model.ScheduleName
	}
	if model.Rrule != nil {
		old.Rrule = model.Rrule
	}

	req = req.UpdateUpdateSchedulePayload(serverupdate.UpdateUpdateSchedulePayload{
		Enabled:           old.Enabled,
		Name:              old.Name,
		Rrule:             old.Rrule,
		MaintenanceWindow: old.MaintenanceWindow,
	})
	return req, nil
}

func outputResult(p *print.Printer, outputFormat string, resp serverupdate.UpdateSchedule) error {
	switch outputFormat {
	case print.JSONOutputFormat:
		details, err := json.MarshalIndent(resp, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal update server os-update schedule: %w", err)
		}
		p.Outputln(string(details))

		return nil
	case print.YAMLOutputFormat:
		details, err := yaml.MarshalWithOptions(resp, yaml.IndentSequence(true), yaml.UseJSONMarshaler())
		if err != nil {
			return fmt.Errorf("marshal update server os-update schedule: %w", err)
		}
		p.Outputln(string(details))

		return nil
	default:
		p.Info("Updated server os-update schedule %s\n", utils.PtrString(resp.Id))
		return nil
	}
}
