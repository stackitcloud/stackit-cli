package describe

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"
	automation "github.com/stackitcloud/stackit-sdk-go/services/automation/v1api"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	cliErr "github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/projectname"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/automation/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

const (
	automationIdArg = "AUTOMATION_ID"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	AutomationId string
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("describe %s", automationIdArg),
		Short: "Shows details of a Volume Automation",
		Long:  "Shows details of a Volume Automation.",
		Args:  args.SingleArg(automationIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Describe the Volume Automation with ID "xxx"`,
				"$ stackit beta volume automation describe xxx"),
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
			resp, err := buildRequest(ctx, model, apiClient).Execute()
			if err != nil {
				return fmt.Errorf("describe volume automation: %w", err)
			}

			// Get projectLabel
			projectLabel, err := projectname.GetProjectName(ctx, params.Printer, params.CliVersion, cmd)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get project name: %v", err)
				projectLabel = model.ProjectId
			} else if projectLabel == "" {
				projectLabel = model.ProjectId
			}

			return outputResult(params.Printer, model.OutputFormat, model.AutomationId, projectLabel, resp)
		},
	}
	return cmd
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *automation.APIClient) automation.ApiGetVolumeAutomationRequest {
	return apiClient.DefaultAPI.GetVolumeAutomation(ctx, model.ProjectId, model.Region, model.AutomationId)
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	automationId := inputArgs[0]

	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &cliErr.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		AutomationId:    automationId,
	}

	p.DebugInputModel(model)
	return &model, nil
}

func outputResult(p *print.Printer, outputFormat, automationId, projectLabel string, item *automation.VolumeAutomation) error {
	return p.OutputResult(outputFormat, item, func() error {
		if item == nil {
			p.Outputf("Volume automation %q not found in project %q\n", automationId, projectLabel)
			return nil
		}

		var content []tables.Table

		table := tables.NewTable()
		table.SetTitle("Volume Automation")

		table.AddRow("ID", item.Id)
		table.AddSeparator()
		table.AddRow("NAME", utils.PtrString(item.Name))
		table.AddSeparator()
		table.AddRow("DESCRIPTION", utils.PtrString(item.Description))
		table.AddSeparator()
		table.AddRow("TEMPLATE ID", utils.PtrString(item.TemplateId))
		table.AddSeparator()
		table.AddRow("CREATE TIME", item.CreateTime.Format(time.DateTime))
		table.AddSeparator()
		table.AddRow("UPDATE TIME", item.UpdateTime.Format(time.DateTime))
		table.AddSeparator()

		content = append(content, table)

		if input := item.Input.Get(); input != nil {
			inputTable := tables.NewTable()
			inputTable.SetTitle("Input")

			inputTable.AddRow("KIND", input.Kind)
			inputTable.AddSeparator()
			for key, value := range utils.FlattenMap(input.AdditionalProperties) {
				inputTable.AddRow(key, value)
				inputTable.AddSeparator()
			}

			content = append(content, inputTable)
		}

		if triggers := item.Triggers.Get(); triggers != nil {
			// As long only schedule with rrule exists as trigger, the whole table should be printed only when rrule is set
			if schedule := triggers.Schedule.Get(); schedule != nil {
				triggersTable := tables.NewTable()
				triggersTable.SetTitle("Triggers")
				triggersTable.AddRow("SCHEDULE - RRULE", schedule.Rrule)
				triggersTable.AddSeparator()

				content = append(content, triggersTable)
			}
		}

		if item.Output != nil && len(item.Output.Steps) > 0 {
			outputStepsTable := tables.NewTable()
			outputStepsTable.SetTitle("Output Steps")
			outputStepsTable.SetHeader("NAME", "KIND", "DETAILS")
			for _, step := range item.Output.Steps {
				var kind, details string
				if step.Result != nil {
					kind = step.Result.Kind
					// step.Result is an open object, so we need to read here the AdditionalProperties
					for key, value := range utils.FlattenMap(step.Result.AdditionalProperties) {
						details += fmt.Sprintf("%s: %v\n", key, value)
					}
				}

				outputStepsTable.AddRow(step.Name, kind, details)
				outputStepsTable.AddSeparator()
			}

			content = append(content, outputStepsTable)
		}

		err := tables.DisplayTables(p, content)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}
		return nil
	})
}
