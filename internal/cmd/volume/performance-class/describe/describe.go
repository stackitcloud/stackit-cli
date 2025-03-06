package describe

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/goccy/go-yaml"

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
	volumePerformanceClassArg = "VOLUME_PERFORMANCE_CLASS"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	VolumePerformanceClass string
}

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("describe %s", volumePerformanceClassArg),
		Short: "Shows details of a volume performance class",
		Long:  "Shows details of a volume performance class.",
		Args:  args.SingleArg(volumePerformanceClassArg, nil),
		Example: examples.Build(
			examples.NewExample(
				`Show details of a volume performance class with name "xxx"`,
				"$ stackit volume performance-class describe xxx",
			),
			examples.NewExample(
				`Show details of a volume performance class with name "xxx" in JSON format`,
				"$ stackit volume performance-class describe xxx --output-format json",
			),
		),
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

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("read volume performance class: %w", err)
			}

			return outputResult(p, model.OutputFormat, resp)
		},
	}
	return cmd
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	volumePerformanceClass := inputArgs[0]

	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel:        globalFlags,
		VolumePerformanceClass: volumePerformanceClass,
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

func buildRequest(ctx context.Context, model *inputModel, apiClient *iaas.APIClient) iaas.ApiGetVolumePerformanceClassRequest {
	return apiClient.GetVolumePerformanceClass(ctx, model.ProjectId, model.VolumePerformanceClass)
}

func outputResult(p *print.Printer, outputFormat string, performanceClass *iaas.VolumePerformanceClass) error {
	if performanceClass == nil {
		return fmt.Errorf("performanceClass response is empty")
	}
	switch outputFormat {
	case print.JSONOutputFormat:
		details, err := json.MarshalIndent(performanceClass, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal volume performance class: %w", err)
		}
		p.Outputln(string(details))

		return nil
	case print.YAMLOutputFormat:
		details, err := yaml.MarshalWithOptions(performanceClass, yaml.IndentSequence(true), yaml.UseJSONMarshaler())
		if err != nil {
			return fmt.Errorf("marshal volume performance class: %w", err)
		}
		p.Outputln(string(details))

		return nil
	default:
		table := tables.NewTable()
		table.AddRow("NAME", utils.PtrString(performanceClass.Name))
		table.AddSeparator()
		table.AddRow("DESCRIPTION", utils.PtrString(performanceClass.Description))
		table.AddSeparator()
		table.AddRow("IOPS", utils.PtrString(performanceClass.Iops))
		table.AddSeparator()
		table.AddRow("THROUGHPUT", utils.PtrString(performanceClass.Throughput))
		table.AddSeparator()

		if performanceClass.Labels != nil && len(*performanceClass.Labels) > 0 {
			labels := []string{}
			for key, value := range *performanceClass.Labels {
				labels = append(labels, fmt.Sprintf("%s: %s", key, value))
			}
			table.AddRow("LABELS", strings.Join(labels, "\n"))
			table.AddSeparator()
		}

		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}
		return nil
	}
}
