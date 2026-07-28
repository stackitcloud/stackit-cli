package updateschedule

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/types"

	"github.com/spf13/cobra"
	postgresflexLegacy "github.com/stackitcloud/stackit-sdk-go/services/postgresflex/v2api"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	cliErr "github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/postgresflex/client"
	postgresflexUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/postgresflex/utils"
)

const (
	// Deprecated: Will be removed after 2027-01-31.
	instanceIdFlag = "instance-id"
	// Deprecated: Will be removed after 2027-01-31.
	scheduleFlag = "schedule"
)

// Deprecated: Will be removed after 2027-01-31.
type inputModel struct {
	*globalflags.GlobalFlagModel

	InstanceId     string
	BackupSchedule string
}

// Deprecated: Will be removed after 2027-01-31.
func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:        "update-schedule",
		Short:      "Updates backup schedule for a PostgreSQL Flex instance",
		Long:       `Updates backup schedule for a PostgreSQL Flex instance. The current backup schedule can be seen in the output of the "stackit postgresflex instance describe" command.`,
		Args:       args.NoArgs,
		Deprecated: `Command "stackit postgresflex backup update-schedule" is deprecated and will be removed after 2027-01-31. Please use "stackit postgresflex instance update --backup-schedule" instead.`,
		Example: examples.Build(
			examples.NewExample(
				`Update the backup schedule of a PostgreSQL Flex instance with ID "xxx"`,
				"$ stackit postgresflex backup update-schedule --instance-id xxx --schedule '6 6 * * *'"),
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

			apiClientLegacy, err := client.ConfigureClientLegacy(params.Printer, params.CliVersion)
			if err != nil {
				return err
			}

			instanceLabel, err := postgresflexUtils.GetInstanceName(ctx, apiClient.DefaultAPI, model.ProjectId, model.Region, model.InstanceId)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get instance name: %v", err)
				instanceLabel = model.InstanceId
			}

			prompt := fmt.Sprintf("Are you sure you want to update backup schedule of instance %q?", instanceLabel)
			err = params.Printer.PromptForConfirmation(prompt)
			if err != nil {
				return err
			}

			// Call API
			req := buildRequest(ctx, model, apiClientLegacy)
			err = req.Execute()
			if err != nil {
				return fmt.Errorf("update backup schedule of PostgreSQL Flex instance: %w", err)
			}

			cmd.Printf("Updated backup schedule of instance %q\n", instanceLabel)
			return nil
		},
	}
	configureFlags(cmd)
	return cmd
}

// Deprecated: Will be removed after 2027-01-31.
func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), instanceIdFlag, "Instance ID")
	cmd.Flags().String(scheduleFlag, "", "Backup schedule, in the cron scheduling system format e.g. '0 0 * * *'")

	err := flags.MarkFlagsRequired(cmd, instanceIdFlag, scheduleFlag)
	cobra.CheckErr(err)
}

// Deprecated: Will be removed after 2027-01-31.
func parseInput(p *print.Printer, cmd *cobra.Command, _ []string) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &cliErr.ProjectIdError{}
	}

	return &inputModel{
		GlobalFlagModel: globalFlags,
		InstanceId:      flags.FlagToStringValue(p, cmd, instanceIdFlag),
		BackupSchedule:  flags.FlagToStringValue(p, cmd, scheduleFlag),
	}, nil
}

// Deprecated: Will be removed after 2027-01-31.
func buildRequest(ctx context.Context, model *inputModel, apiClient *postgresflexLegacy.APIClient) postgresflexLegacy.ApiUpdateBackupScheduleRequest {
	req := apiClient.DefaultAPI.UpdateBackupSchedule(ctx, model.ProjectId, model.Region, model.InstanceId)
	req = req.UpdateBackupSchedulePayload(postgresflexLegacy.UpdateBackupSchedulePayload{
		BackupSchedule: model.BackupSchedule,
	})
	return req
}
