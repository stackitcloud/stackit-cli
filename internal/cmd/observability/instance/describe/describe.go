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
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/observability/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/observability"
)

const (
	instanceIdArg = "INSTANCE_ID"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	InstanceId string
}

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("describe %s", instanceIdArg),
		Short: "Shows details of an Observability instance",
		Long:  "Shows details of an Observability instance.",
		Args:  args.SingleArg(instanceIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Get details of an Observability instance with ID "xxx"`,
				"$ stackit observability instance describe xxx"),
			examples.NewExample(
				`Get details of an Observability instance with ID "xxx" in JSON format`,
				"$ stackit observability instance describe xxx --output-format json"),
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
				return fmt.Errorf("read Observability instance: %w", err)
			}

			return outputResult(params.Printer, model.OutputFormat, resp)
		},
	}
	return cmd
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	instanceId := inputArgs[0]

	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		InstanceId:      instanceId,
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

func buildRequest(ctx context.Context, model *inputModel, apiClient *observability.APIClient) observability.ApiGetInstanceRequest {
	req := apiClient.GetInstance(ctx, model.InstanceId, model.ProjectId)
	return req
}

func outputResult(p *print.Printer, outputFormat string, instance *observability.GetInstanceResponse) error {
	if instance == nil {
		return fmt.Errorf("instance is nil")
	}

	switch outputFormat {
	case print.JSONOutputFormat:
		details, err := json.MarshalIndent(instance, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal Observability instance: %w", err)
		}
		p.Outputln(string(details))

		return nil
	case print.YAMLOutputFormat:
		details, err := yaml.MarshalWithOptions(instance, yaml.IndentSequence(true), yaml.UseJSONMarshaler())
		if err != nil {
			return fmt.Errorf("marshal Observability instance: %w", err)
		}
		p.Outputln(string(details))

		return nil
	default:
		table := tables.NewTable()
		table.AddRow("ID", utils.PtrString(instance.Id))
		table.AddSeparator()
		table.AddRow("NAME", utils.PtrString(instance.Name))
		table.AddSeparator()
		table.AddRow("STATUS", utils.PtrString(instance.Status))
		table.AddSeparator()
		table.AddRow("PLAN NAME", utils.PtrString(instance.PlanName))
		table.AddSeparator()
		if inst := instance.Instance; inst != nil {
			if plan := inst.Plan; plan != nil {
				table.AddRow("METRIC SAMPLES (PER MIN)", utils.PtrString(plan.TotalMetricSamples))
				table.AddSeparator()
				table.AddRow("LOGS (GB)", utils.PtrString(plan.LogsStorage))
				table.AddSeparator()
				table.AddRow("TRACES (GB)", utils.PtrString(plan.TracesStorage))
				table.AddSeparator()
				table.AddRow("NOTIFICATION RULES", utils.PtrString(plan.AlertRules))
				table.AddSeparator()
				table.AddRow("GRAFANA USERS", utils.PtrString(plan.GrafanaGlobalUsers))
				table.AddSeparator()
			}
			table.AddRow("GRAFANA URL", utils.PtrString(inst.GrafanaUrl))
			table.AddSeparator()
		}
		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}

		return nil
	}
}
