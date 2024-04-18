package describe

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/argus/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/argus"
)

const (
	instanceIdArg = "INSTANCE_ID"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	InstanceId string
}

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("describe %s", instanceIdArg),
		Short: "Shows details of an Argus instance",
		Long:  "Shows details of an Argus instance.",
		Args:  args.SingleArg(instanceIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Get details of an Argus instance with ID "xxx"`,
				"$ stackit argus instance describe xxx"),
			examples.NewExample(
				`Get details of an Argus instance with ID "xxx" in a table format`,
				"$ stackit argus instance describe xxx --output-format pretty"),
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
				return fmt.Errorf("read Argus instance: %w", err)
			}

			return outputResult(p, model.OutputFormat, resp)
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

	return &inputModel{
		GlobalFlagModel: globalFlags,
		InstanceId:      instanceId,
	}, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *argus.APIClient) argus.ApiGetInstanceRequest {
	req := apiClient.GetInstance(ctx, model.InstanceId, model.ProjectId)
	return req
}

func outputResult(p *print.Printer, outputFormat string, instance *argus.GetInstanceResponse) error {
	switch outputFormat {
	case globalflags.PrettyOutputFormat:

		table := tables.NewTable()
		table.AddRow("ID", *instance.Id)
		table.AddSeparator()
		table.AddRow("NAME", *instance.Name)
		table.AddSeparator()
		table.AddRow("STATUS", *instance.Status)
		table.AddSeparator()
		table.AddRow("PLAN NAME", *instance.PlanName)
		table.AddSeparator()
		table.AddRow("METRIC SAMPLES (PER MIN)", *instance.Instance.Plan.TotalMetricSamples)
		table.AddSeparator()
		table.AddRow("LOGS (GB)", *instance.Instance.Plan.LogsStorage)
		table.AddSeparator()
		table.AddRow("TRACES (GB)", *instance.Instance.Plan.TracesStorage)
		table.AddSeparator()
		table.AddRow("NOTIFICATION RULES", *instance.Instance.Plan.AlertRules)
		table.AddSeparator()
		table.AddRow("GRAFANA USERS", *instance.Instance.Plan.GrafanaGlobalUsers)
		table.AddSeparator()
		table.AddRow("GRAFANA URL", *instance.Instance.GrafanaUrl)
		table.AddSeparator()
		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}

		return nil
	default:
		details, err := json.MarshalIndent(instance, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal Argus instance: %w", err)
		}
		p.Outputln(string(details))

		return nil
	}
}
