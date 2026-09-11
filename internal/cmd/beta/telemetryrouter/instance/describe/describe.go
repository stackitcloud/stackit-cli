package describe

import (
	"context"
	"fmt"
	"strings"
	"time"

	telemetryrouter "github.com/stackitcloud/stackit-sdk-go/services/telemetryrouter/v1api"

	"github.com/stackitcloud/stackit-cli/internal/pkg/types"

	"github.com/spf13/cobra"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/telemetryrouter/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

const (
	instanceIdArg = "INSTANCE_ID"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	InstanceId string
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("describe %s", instanceIdArg),
		Short: "Shows details of a Telemetry Router instance",
		Long:  "Shows details of a Telemetry Router instance",
		Args:  args.SingleArg(instanceIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Get details of a Telemetry Router instance with ID "xxx"`,
				`$ stackit beta telemetryrouter instance describe xxx`,
			),
			examples.NewExample(
				`Get details of a Telemetry Router instance with ID "xxx" in JSON format`,
				"$ stackit beta telemetryrouter instance describe xxx --output-format json"),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()

			model, err := parseInput(params.Printer, cmd, args)
			if err != nil {
				return err
			}

			apiClient, err := client.ConfigureClient(params.Printer, params.CliVersion)
			if err != nil {
				return err
			}

			req := buildRequest(ctx, model, apiClient.DefaultAPI)

			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("get instance: %w", err)
			}

			return outputResult(params.Printer, model.OutputFormat, resp)
		},
	}
	return cmd
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}
	model := &inputModel{
		GlobalFlagModel: globalFlags,
		InstanceId:      inputArgs[0],
	}
	p.DebugInputModel(model)
	return model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient telemetryrouter.DefaultAPI) telemetryrouter.ApiGetTelemetryRouterRequest {
	return apiClient.GetTelemetryRouter(ctx, model.ProjectId, model.Region, model.InstanceId)
}

func outputResult(p *print.Printer, outputFormat string, instance *telemetryrouter.TelemetryRouterResponse) error {
	if instance == nil {
		return fmt.Errorf("instance response is empty")
	}
	return p.OutputResult(outputFormat, instance, func() error {
		var content []tables.Table

		table := tables.NewTable()
		table.AddRow("ID", instance.Id)
		table.AddSeparator()
		table.AddRow("DISPLAY NAME", instance.DisplayName)
		table.AddSeparator()
		table.AddRow("DESCRIPTION", utils.PtrString(instance.Description))
		table.AddSeparator()
		table.AddRow("STATUS", instance.Status)
		table.AddSeparator()
		table.AddRow("URI", instance.Uri)
		table.AddSeparator()
		table.AddRow("CREATION TIME", instance.CreationTime.Format(time.DateTime))
		content = append(content, table)

		if instance.Filter != nil && len(instance.Filter.Attributes) > 0 {
			filterTable := tables.NewTable()
			filterTable.SetTitle("FILTER ATTRIBUTES")
			filterTable.SetHeader("KEY", "LEVEL", "MATCHER", "VALUES")
			for _, attr := range instance.Filter.Attributes {
				filterTable.AddRow(attr.Key, attr.Level, attr.Matcher, strings.Join(attr.Values, ","))
			}
			content = append(content, filterTable)
		}

		if err := tables.DisplayTables(p, content); err != nil {
			return fmt.Errorf("display tables: %w", err)
		}
		return nil
	})
}
