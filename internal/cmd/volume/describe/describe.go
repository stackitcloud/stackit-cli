package describe

import (
	"context"
	"fmt"
	"strings"

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

	"github.com/spf13/cobra"
)

const (
	volumeIdArg = "VOLUME_ID"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	VolumeId string
}

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("describe %s", volumeIdArg),
		Short: "Shows details of a volume",
		Long:  "Shows details of a volume.",
		Args:  args.SingleArg(volumeIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Show details of a volume with ID "xxx"`,
				"$ stackit volume describe xxx",
			),
			examples.NewExample(
				`Show details of a volume with ID "xxx" in JSON format`,
				"$ stackit volume describe xxx --output-format json",
			),
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
				return fmt.Errorf("read volume: %w", err)
			}

			return outputResult(params.Printer, model.OutputFormat, resp)
		},
	}
	return cmd
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	volumeId := inputArgs[0]

	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		VolumeId:        volumeId,
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *iaas.APIClient) iaas.ApiGetVolumeRequest {
	return apiClient.GetVolume(ctx, model.ProjectId, model.Region, model.VolumeId)
}

func outputResult(p *print.Printer, outputFormat string, volume *iaas.Volume) error {
	if volume == nil {
		return fmt.Errorf("volume response is empty")
	}
	return p.OutputResult(outputFormat, volume, func() error {
		table := tables.NewTable()
		table.AddRow("ID", utils.PtrString(volume.Id))
		table.AddSeparator()
		table.AddRow("NAME", utils.PtrString(volume.Name))
		table.AddSeparator()
		table.AddRow("STATE", utils.PtrString(volume.Status))
		table.AddSeparator()
		table.AddRow("VOLUME SIZE (GB)", utils.PtrString(volume.Size))
		table.AddSeparator()
		table.AddRow("PERFORMANCE CLASS", utils.PtrString(volume.PerformanceClass))
		table.AddSeparator()
		table.AddRow("AVAILABILITY ZONE", utils.PtrString(volume.AvailabilityZone))
		table.AddSeparator()

		if volume.Source != nil {
			sourceId := *volume.Source.Id
			table.AddRow("SOURCE", sourceId)
			table.AddSeparator()
		}

		if volume.ServerId != nil {
			serverId := *volume.ServerId
			table.AddRow("SERVER", serverId)
			table.AddSeparator()
		}

		if volume.Labels != nil && len(*volume.Labels) > 0 {
			labels := []string{}
			for key, value := range *volume.Labels {
				labels = append(labels, fmt.Sprintf("%s: %s", key, value))
			}
			table.AddRow("LABELS", strings.Join(labels, "\n"))
		}

		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}
		return nil
	})
}
