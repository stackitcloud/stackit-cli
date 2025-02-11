package create

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/goccy/go-yaml"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	cliErr "github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
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

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Creates a Server os-update Schedule",
		Long:  "Creates a Server os-update Schedule.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Create a Server os-update Schedule with name "myschedule"`,
				`$ stackit beta server os-update schedule create --server-id xxx --name=myschedule`),
			examples.NewExample(
				`Create a Server os-update Schedule with name "myschedule" and maintenance window for 14 o'clock`,
				`$ stackit beta server os-update schedule create --server-id xxx --name=myschedule --maintenance-window=14`),
		),
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := context.Background()

			model, err := parseInput(p, cmd)
			if err != nil {
				return err
			}

			// Configure API client
			apiClient, err := client.ConfigureClient(p)
			if err != nil {
				return err
			}

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to create a os-update Schedule for server %s?", model.ServerId)
				err = p.PromptForConfirmation(prompt)
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

			return outputResult(p, model, resp)
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

func parseInput(p *print.Printer, cmd *cobra.Command) (*inputModel, error) {
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

func buildRequest(ctx context.Context, model *inputModel, apiClient *serverupdate.APIClient) (serverupdate.ApiCreateUpdateScheduleRequest, error) {
	req := apiClient.CreateUpdateSchedule(ctx, model.ProjectId, model.ServerId)
	req = req.CreateUpdateSchedulePayload(serverupdate.CreateUpdateSchedulePayload{
		Enabled:           &model.Enabled,
		Name:              &model.ScheduleName,
		Rrule:             &model.Rrule,
		MaintenanceWindow: &model.MaintenanceWindow,
	})
	return req, nil
}

func outputResult(p *print.Printer, model *inputModel, resp *serverupdate.UpdateSchedule) error {
	switch model.OutputFormat {
	case print.JSONOutputFormat:
		details, err := json.MarshalIndent(resp, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal server os-update schedule: %w", err)
		}
		p.Outputln(string(details))

		return nil
	case print.YAMLOutputFormat:
		details, err := yaml.MarshalWithOptions(resp, yaml.IndentSequence(true), yaml.UseJSONMarshaler())
		if err != nil {
			return fmt.Errorf("marshal server os-update schedule: %w", err)
		}
		p.Outputln(string(details))

		return nil
	default:
		p.Outputf("Created server os-update schedule for server %s. os-update Schedule ID: %s\n", model.ServerId, utils.PtrString(resp.Id))
		return nil
	}
}
