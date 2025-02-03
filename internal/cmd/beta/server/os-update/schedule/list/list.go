package list

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/goccy/go-yaml"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/serverosupdate/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/serverupdate"
)

const (
	limitFlag    = "limit"
	serverIdFlag = "server-id"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	ServerId string
	Limit    *int64
}

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Lists all server os-update schedules",
		Long:  "Lists all server os-update schedules.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`List all os-update schedules for a server with ID "xxx"`,
				"$ stackit beta server os-update schedule list --server-id xxx"),
			examples.NewExample(
				`List all os-update schedules for a server with ID "xxx" in JSON format`,
				"$ stackit beta server os-update schedule list --server-id xxx --output-format json"),
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

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("list server os-update schedules: %w", err)
			}
			schedules := *resp.Items
			if len(schedules) == 0 {
				p.Info("No os-update schedules found for server %s\n", model.ServerId)
				return nil
			}

			// Truncate output
			if model.Limit != nil && len(schedules) > int(*model.Limit) {
				schedules = schedules[:*model.Limit]
			}
			return outputResult(p, model.OutputFormat, schedules)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Int64(limitFlag, 0, "Maximum number of entries to list")
	cmd.Flags().VarP(flags.UUIDFlag(), serverIdFlag, "s", "Server ID")

	err := flags.MarkFlagsRequired(cmd, serverIdFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	limit := flags.FlagToInt64Pointer(p, cmd, limitFlag)
	if limit != nil && *limit < 1 {
		return nil, &errors.FlagValidationError{
			Flag:    limitFlag,
			Details: "must be greater than 0",
		}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		ServerId:        flags.FlagToStringValue(p, cmd, serverIdFlag),
		Limit:           limit,
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

func buildRequest(ctx context.Context, model *inputModel, apiClient *serverupdate.APIClient) serverupdate.ApiListUpdateSchedulesRequest {
	req := apiClient.ListUpdateSchedules(ctx, model.ProjectId, model.ServerId)
	return req
}

func outputResult(p *print.Printer, outputFormat string, schedules []serverupdate.UpdateSchedule) error {
	switch outputFormat {
	case print.JSONOutputFormat:
		details, err := json.MarshalIndent(schedules, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal Server os-update Schedules list: %w", err)
		}
		p.Outputln(string(details))

		return nil
	case print.YAMLOutputFormat:
		details, err := yaml.MarshalWithOptions(schedules, yaml.IndentSequence(true))
		if err != nil {
			return fmt.Errorf("marshal Server os-update Schedules list: %w", err)
		}
		p.Outputln(string(details))

		return nil
	default:
		table := tables.NewTable()
		table.SetHeader("SCHEDULE ID", "SCHEDULE NAME", "ENABLED", "RRULE", "MAINTENANCE WINDOW")
		for i := range schedules {
			s := schedules[i]
			table.AddRow(
				utils.PtrString(s.Id),
				utils.PtrString(s.Name),
				utils.PtrString(s.Enabled),
				utils.PtrString(s.Rrule),
				utils.PtrString(s.MaintenanceWindow),
			)
		}
		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}
		return nil
	}
}
