package describe

import (
	"context"
	"fmt"
	"strings"

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
	argInstanceID = "INSTANCE_ID"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	InstanceID string
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("describe %s", argInstanceID),
		Short: "Shows details of a TelemetryRouter instance",
		Long:  "Shows details of a TelemetryRouter instance",
		Args:  args.SingleArg(argInstanceID, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Get details of a TelemetryRouter instance with ID "xxx"`,
				`$ stackit beta telemetryrouter instance describe xxx`,
			),
			examples.NewExample(
				`Get details of a TelemetryRouter instance with ID "xxx" in JSON format`,
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

			req := buildRequest(ctx, model, apiClient)

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
		InstanceID:      inputArgs[0],
	}
	p.DebugInputModel(model)
	return model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *telemetryrouter.APIClient) telemetryrouter.ApiGetTelemetryRouterRequest {
	return apiClient.DefaultAPI.GetTelemetryRouter(ctx, model.ProjectId, model.Region, model.InstanceID)
}

func outputResult(p *print.Printer, outputFormat string, instance *telemetryrouter.TelemetryRouterResponse) error {
	if instance == nil {
		return fmt.Errorf("instance response is empty")
	}
	return p.OutputResult(outputFormat, instance, func() error {
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
		table.AddRow("CREATION TIME", instance.CreationTime)
		table.AddSeparator()
		table.AddRow("FILTER ATTRIBUTES", formatFilterAttributes(instance.Filter))

		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("display table: %w", err)
		}
		return nil
	})
}

// formatFilterAttributes renders a ConfigFilter's attributes for table output, using the same
// "key=...;level=...;matcher=...;values=..." syntax accepted by the --filter-attribute flag.
func formatFilterAttributes(filter *telemetryrouter.ConfigFilter) string {
	if filter == nil || len(filter.Attributes) == 0 {
		return "-"
	}

	lines := make([]string, 0, len(filter.Attributes))
	for _, attr := range filter.Attributes {
		lines = append(lines, fmt.Sprintf("key=%s;level=%s;matcher=%s;values=%s",
			attr.Key, attr.Level, attr.Matcher, strings.Join(attr.Values, ",")))
	}
	return strings.Join(lines, "\n")
}
