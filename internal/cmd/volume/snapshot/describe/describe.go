package describe

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/goccy/go-yaml"
	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/iaas/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/stackitcloud/stackit-sdk-go/services/iaas"
)

const (
	snapshotIdArg = "SNAPSHOT_ID"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	SnapshotId string
}

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("describe %s", snapshotIdArg),
		Short: "Describes a snapshot",
		Long:  "Describes a snapshot by its ID.",
		Args:  args.SingleArg(snapshotIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Get details of a snapshot`,
				"$ stackit volume snapshot describe xxx-xxx-xxx"),
			examples.NewExample(
				`Get details of a snapshot in JSON format`,
				"$ stackit volume snapshot describe xxx-xxx-xxx --output-format json"),
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

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("get snapshot details: %w", err)
			}

			return outputResult(params.Printer, model.OutputFormat, resp)
		},
	}
	return cmd
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	snapshotId := inputArgs[0]

	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		SnapshotId:      snapshotId,
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

func buildRequest(ctx context.Context, model *inputModel, apiClient *iaas.APIClient) iaas.ApiGetSnapshotRequest {
	return apiClient.GetSnapshot(ctx, model.ProjectId, model.SnapshotId)
}

func outputResult(p *print.Printer, outputFormat string, snapshot *iaas.Snapshot) error {
	if snapshot == nil {
		return fmt.Errorf("get snapshot response is empty")
	}

	switch outputFormat {
	case print.JSONOutputFormat:
		details, err := json.MarshalIndent(snapshot, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal snapshot: %w", err)
		}
		p.Outputln(string(details))
		return nil

	case print.YAMLOutputFormat:
		details, err := yaml.MarshalWithOptions(snapshot, yaml.IndentSequence(true), yaml.UseJSONMarshaler())
		if err != nil {
			return fmt.Errorf("marshal snapshot: %w", err)
		}
		p.Outputln(string(details))
		return nil

	default:
		table := tables.NewTable()
		table.SetHeader("ID", "NAME", "SIZE", "STATUS", "VOLUME ID", "LABELS", "CREATED AT", "UPDATED AT")

		table.AddRow(
			utils.PtrString(snapshot.Id),
			utils.PtrString(snapshot.Name),
			utils.PtrByteSizeDefault((*int64)(snapshot.Size), ""),
			utils.PtrString(snapshot.Status),
			utils.PtrString(snapshot.VolumeId),
			utils.PtrStringDefault(snapshot.Labels, ""),
			utils.ConvertTimePToDateTimeString(snapshot.CreatedAt),
			utils.ConvertTimePToDateTimeString(snapshot.UpdatedAt),
		)

		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}

		return nil
	}
}
