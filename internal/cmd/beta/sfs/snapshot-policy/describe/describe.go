package describe

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	sfs "github.com/stackitcloud/stackit-sdk-go/services/sfs/v1api"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/projectname"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/sfs/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

const snapshotPolicyIdArg = "SNAPSHOT_POLICY_ID"

type inputModel struct {
	*globalflags.GlobalFlagModel
	SnapshotPolicyId string
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("describe %s", snapshotPolicyIdArg),
		Short: "Shows details of a snapshot policy",
		Long:  "Shows details of a snapshot policy.",
		Args:  args.SingleArg(snapshotPolicyIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Describe a snapshot policy with ID "xxx"`,
				"$ stackit beta sfs snapshot-policy describe xxx",
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			model, err := parseInput(params.Printer, cmd, args)
			if err != nil {
				return fmt.Errorf("unable to parse input: %w", err)
			}

			// Configure API client
			apiClient, err := client.ConfigureClient(params.Printer, params.CliVersion)
			if err != nil {
				return err
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("read snapshot policy: %w", err)
			}

			// Get projectLabel
			projectLabel, err := projectname.GetProjectName(ctx, params.Printer, params.CliVersion, cmd)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get project name: %v", err)
				projectLabel = model.ProjectId
			} else if projectLabel == "" {
				projectLabel = model.ProjectId
			}

			return outputResult(params.Printer, model.OutputFormat, model.SnapshotPolicyId, projectLabel, resp)
		},
	}
	return cmd
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	snapshotPolicyId := inputArgs[0]

	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel:  globalFlags,
		SnapshotPolicyId: snapshotPolicyId,
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *sfs.APIClient) sfs.ApiGetSnapshotPolicyRequest {
	return apiClient.DefaultAPI.GetSnapshotPolicy(ctx, model.ProjectId, model.SnapshotPolicyId)
}

func outputResult(p *print.Printer, outputFormat, snapshotPolicyId, projectLabel string, snapshotPolicy *sfs.GetSnapshotPolicyResponse) error {
	return p.OutputResult(outputFormat, snapshotPolicy, func() error {
		if snapshotPolicy == nil || snapshotPolicy.SnapshotPolicy == nil {
			p.Outputf("Snapshot policy %q not found in project %q", snapshotPolicyId, projectLabel)
			return nil
		}

		var content []tables.Table

		table := tables.NewTable()
		table.SetTitle("Snapshot Policy")
		policy := snapshotPolicy.SnapshotPolicy

		table.AddRow("ID", utils.PtrString(policy.Id))
		table.AddSeparator()
		table.AddRow("NAME", utils.PtrString(policy.Name))
		table.AddSeparator()
		table.AddRow("ENABLED", utils.PtrString(policy.Enabled))
		table.AddSeparator()
		table.AddRow("COMMENT", utils.PtrString(policy.Comment))
		table.AddSeparator()
		table.AddRow("CREATED AT", utils.ConvertTimePToDateTimeString(policy.CreatedAt))

		content = append(content, table)

		if len(policy.SnapshotSchedules) > 0 {
			snapshotSchedulesTable := tables.NewTable()
			snapshotSchedulesTable.SetTitle("Snapshot Schedules")

			snapshotSchedulesTable.SetHeader("ID", "NAME", "INTERVAL", "PREFIX", "RETENTION COUNT", "RETENTION PERIOD", "CREATED AT")

			for _, snapshotSchedule := range policy.SnapshotSchedules {
				snapshotSchedulesTable.AddRow(
					utils.PtrString(snapshotSchedule.Id),
					utils.PtrString(snapshotSchedule.Name),
					utils.PtrString(snapshotSchedule.Interval),
					utils.PtrString(snapshotSchedule.Prefix),
					utils.PtrString(snapshotSchedule.RetentionCount),
					utils.PtrString(snapshotSchedule.RetentionPeriod),
					utils.ConvertTimePToDateTimeString(snapshotSchedule.CreatedAt),
				)
				snapshotSchedulesTable.AddSeparator()
			}

			content = append(content, snapshotSchedulesTable)
		}

		if err := tables.DisplayTables(p, content); err != nil {
			return fmt.Errorf("render tables: %w", err)
		}
		return nil
	})
}
