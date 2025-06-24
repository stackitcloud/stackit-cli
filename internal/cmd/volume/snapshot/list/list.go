package list

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/goccy/go-yaml"
	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/iaas/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/stackitcloud/stackit-cli/internal/pkg/projectname"
	"github.com/stackitcloud/stackit-sdk-go/services/iaas"
)

const (
	limitFlag         = "limit"
	labelSelectorFlag = "label-selector"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	Limit         *int64
	LabelSelector *string
}

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Lists all snapshots",
		Long:  "Lists all snapshots in a project.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`List all snapshots`,
				"$ stackit volume snapshot list"),
			examples.NewExample(
				`List snapshots with a limit of 10`,
				"$ stackit volume snapshot list --limit 10"),
			examples.NewExample(
				`List snapshots filtered by label`,
				"$ stackit volume snapshot list --label-selector key1=value1"),
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

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("list snapshots: %w", err)
			}

			// Check if response is empty
			if resp.Items == nil || len(*resp.Items) == 0 {
				projectLabel, err := projectname.GetProjectName(ctx, params.Printer, params.CliVersion, cmd)
				if err != nil {
					params.Printer.Debug(print.ErrorLevel, "get project name: %v", err)
					projectLabel = model.ProjectId
				}
				params.Printer.Info("No snapshots found for project %q\n", projectLabel)
				return nil
			}

			snapshots := *resp.Items

			// Apply limit if specified
			if model.Limit != nil && int(*model.Limit) < len(snapshots) {
				snapshots = snapshots[:*model.Limit]
			}

			return outputResult(params.Printer, model.OutputFormat, snapshots)
		},
	}

	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Int64(limitFlag, 0, "Maximum number of entries to list")
	cmd.Flags().String(labelSelectorFlag, "", "Filter snapshots by labels")
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

	labelSelector := flags.FlagToStringPointer(p, cmd, labelSelectorFlag)

	model := inputModel{
		GlobalFlagModel: globalFlags,
		Limit:           limit,
		LabelSelector:   labelSelector,
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

func buildRequest(ctx context.Context, model *inputModel, apiClient *iaas.APIClient) iaas.ApiListSnapshotsRequest {
	req := apiClient.ListSnapshots(ctx, model.ProjectId)
	if model.LabelSelector != nil {
		req = req.LabelSelector(*model.LabelSelector)
	}
	return req
}

func outputResult(p *print.Printer, outputFormat string, snapshots []iaas.Snapshot) error {
	if snapshots == nil {
		return fmt.Errorf("list snapshots response is empty")
	}

	switch outputFormat {
	case print.JSONOutputFormat:
		details, err := json.MarshalIndent(snapshots, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal snapshots: %w", err)
		}
		p.Outputln(string(details))
		return nil

	case print.YAMLOutputFormat:
		details, err := yaml.MarshalWithOptions(snapshots, yaml.IndentSequence(true), yaml.UseJSONMarshaler())
		if err != nil {
			return fmt.Errorf("marshal snapshots: %w", err)
		}
		p.Outputln(string(details))
		return nil

	default:
		table := tables.NewTable()
		table.SetHeader("ID", "NAME", "SIZE", "STATUS", "VOLUME ID", "LABELS", "CREATED AT", "UPDATED AT")

		for _, snapshot := range snapshots {
			var labelsString string
			if snapshot.Labels != nil {
				var labels []string
				for key, value := range *snapshot.Labels {
					labels = append(labels, fmt.Sprintf("%s: %s", key, value))
				}
				labelsString = strings.Join(labels, "\n")
			}
			table.AddRow(
				utils.PtrString(snapshot.Id),
				utils.PtrString(snapshot.Name),
				utils.PtrByteSizeDefault(snapshot.Size, ""),
				utils.PtrString(snapshot.Status),
				utils.PtrString(snapshot.VolumeId),
				labelsString,
				utils.ConvertTimePToDateTimeString(snapshot.CreatedAt),
				utils.ConvertTimePToDateTimeString(snapshot.UpdatedAt),
			)
			table.AddSeparator()
		}

		p.Outputln(table.Render())
		return nil
	}
}
