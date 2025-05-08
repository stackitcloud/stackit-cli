package describe

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/goccy/go-yaml"

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
	machineTypeArg = "MACHINE_TYPE"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	MachineType string
}

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("describe %s", machineTypeArg),
		Short: "Shows details of a server machine type",
		Long:  "Shows details of a server machine type.",
		Args:  args.SingleArg(machineTypeArg, nil),
		Example: examples.Build(
			examples.NewExample(
				`Show details of a server machine type with name "xxx"`,
				"$ stackit server machine-type describe xxx",
			),
			examples.NewExample(
				`Show details of a server machine type with name "xxx" in JSON format`,
				"$ stackit server machine-type describe xxx --output-format json",
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			model, err := parseInput(params.Printer, cmd, args)
			if err != nil {
				return err
			}

			// Configure API client
			apiClient, err := client.ConfigureClient(params.Printer)
			if err != nil {
				return err
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("read server machine type: %w", err)
			}

			return outputResult(params.Printer, model.OutputFormat, resp)
		},
	}
	return cmd
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	machineType := inputArgs[0]

	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		MachineType:     machineType,
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

func buildRequest(ctx context.Context, model *inputModel, apiClient *iaas.APIClient) iaas.ApiGetMachineTypeRequest {
	return apiClient.GetMachineType(ctx, model.ProjectId, model.MachineType)
}

func outputResult(p *print.Printer, outputFormat string, machineType *iaas.MachineType) error {
	if machineType == nil {
		return fmt.Errorf("api response for machine type is empty")
	}
	switch outputFormat {
	case print.JSONOutputFormat:
		details, err := json.MarshalIndent(machineType, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal server machine type: %w", err)
		}
		p.Outputln(string(details))

		return nil
	case print.YAMLOutputFormat:
		details, err := yaml.MarshalWithOptions(machineType, yaml.IndentSequence(true), yaml.UseJSONMarshaler())
		if err != nil {
			return fmt.Errorf("marshal server machine type: %w", err)
		}
		p.Outputln(string(details))

		return nil
	default:
		table := tables.NewTable()
		table.AddRow("NAME", utils.PtrString(machineType.Name))
		table.AddSeparator()
		table.AddRow("VCPUS", utils.PtrString(machineType.Vcpus))
		table.AddSeparator()
		table.AddRow("RAM (in MB)", utils.PtrString(machineType.Ram))
		table.AddSeparator()
		table.AddRow("DISK SIZE (in GB)", utils.PtrString(machineType.Disk))
		table.AddSeparator()
		table.AddRow("DESCRIPTION", utils.PtrString(machineType.Description))
		table.AddSeparator()

		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}
		return nil
	}
}
