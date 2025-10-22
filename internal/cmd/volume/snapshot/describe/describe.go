package describe

import (
	"context"
	"fmt"
	"strings"

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
				`Get details of a snapshot with ID "xxx"`,
				"$ stackit volume snapshot describe xxx"),
			examples.NewExample(
				`Get details of a snapshot with ID "xxx" in JSON format`,
				"$ stackit volume snapshot describe xxx --output-format json"),
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

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *iaas.APIClient) iaas.ApiGetSnapshotRequest {
	return apiClient.GetSnapshot(ctx, model.ProjectId, model.SnapshotId)
}

func outputResult(p *print.Printer, outputFormat string, snapshot *iaas.Snapshot) error {
	if snapshot == nil {
		return fmt.Errorf("get snapshot response is empty")
	}

	return p.OutputResult(outputFormat, snapshot, func() error {
		table := tables.NewTable()
		table.AddRow("ID", utils.PtrString(snapshot.Id))
		table.AddSeparator()
		table.AddRow("NAME", utils.PtrString(snapshot.Name))
		table.AddSeparator()
		table.AddRow("SIZE", utils.PtrGigaByteSizeDefault(snapshot.Size, "n/a"))
		table.AddSeparator()
		table.AddRow("STATUS", utils.PtrString(snapshot.Status))
		table.AddSeparator()
		table.AddRow("VOLUME ID", utils.PtrString(snapshot.VolumeId))
		table.AddSeparator()

		if snapshot.Labels != nil && len(*snapshot.Labels) > 0 {
			labels := []string{}
			for key, value := range *snapshot.Labels {
				labels = append(labels, fmt.Sprintf("%s: %s", key, value))
			}
			table.AddRow("LABELS", strings.Join(labels, "\n"))
			table.AddSeparator()
		}

		table.AddRow("CREATED AT", utils.ConvertTimePToDateTimeString(snapshot.CreatedAt))
		table.AddSeparator()
		table.AddRow("UPDATED AT", utils.ConvertTimePToDateTimeString(snapshot.UpdatedAt))

		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}

		return nil
	})
}
