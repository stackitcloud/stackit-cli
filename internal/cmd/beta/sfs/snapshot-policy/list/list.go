package list

import (
	"context"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	sfs "github.com/stackitcloud/stackit-sdk-go/services/sfs/v1api"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/projectname"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/sfs/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

const (
	limitFlag     = "limit"
	immutableFlag = "immutable"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	Limit     *int64
	Immutable bool
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Lists all snapshot policies of a project",
		Long:  "Lists all snapshot policies of a project.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`List all snapshot policies`,
				"$ stackit beta sfs snapshot-policy list",
			),
			examples.NewExample(
				`List all immutable snapshot policies`,
				"$ stackit beta sfs snapshot-policy list --immutable",
			),
			examples.NewExample(
				`List up to 10 snapshot policies`,
				"$ stackit beta sfs snapshot-policy list --limit 10",
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
				return fmt.Errorf("list snapshot policies: %w", err)
			}

			// Get projectLabel
			projectLabel, err := projectname.GetProjectName(ctx, params.Printer, params.CliVersion, cmd)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get project name: %v", err)
				projectLabel = model.ProjectId
			} else if projectLabel == "" {
				projectLabel = model.ProjectId
			}

			// Truncate output
			items := utils.GetSliceFromPointer(&resp.SnapshotPolicies)
			if model.Limit != nil && len(items) > int(*model.Limit) {
				items = items[:*model.Limit]
			}

			return outputResult(params.Printer, model.OutputFormat, projectLabel, items)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Int64(limitFlag, 0, "Maximum number of entries to list")
	cmd.Flags().Bool(immutableFlag, false, "Immutable snapshot policy")
}

func parseInput(p *print.Printer, cmd *cobra.Command, _ []string) (*inputModel, error) {
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
		Limit:           limit,
		Immutable:       flags.FlagToBoolValue(p, cmd, immutableFlag),
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *sfs.APIClient) sfs.ApiListSnapshotPoliciesRequest {
	req := apiClient.DefaultAPI.ListSnapshotPolicies(ctx, model.ProjectId)
	if model.Immutable {
		req = req.Immutable(true)
	}
	return req
}

func outputResult(p *print.Printer, outputFormat, projectLabel string, snapshotPolicies []sfs.SnapshotPolicy) error {
	return p.OutputResult(outputFormat, snapshotPolicies, func() error {
		if len(snapshotPolicies) == 0 {
			p.Outputf("No snapshot policies found for project %q\n", projectLabel)
			return nil
		}

		table := tables.NewTable()
		table.SetHeader("ID", "NAME", "COMMENT", "ENABLED", "AMOUNT OF SNAPSHOT SCHEDULES", "CREATED AT")

		for _, snapshotPolicy := range snapshotPolicies {
			amountSnapshotSchedules := "-"
			if snapshotPolicy.SnapshotSchedules != nil {
				amountSnapshotSchedules = strconv.Itoa(len(snapshotPolicy.SnapshotSchedules))
			}
			table.AddRow(
				utils.PtrString(snapshotPolicy.Id),
				utils.PtrString(snapshotPolicy.Name),
				utils.PtrString(snapshotPolicy.Comment),
				utils.PtrString(snapshotPolicy.Enabled),
				amountSnapshotSchedules,
				utils.ConvertTimePToDateTimeString(snapshotPolicy.CreatedAt),
			)
		}
		p.Outputln(table.Render())
		return nil
	})
}
