package describe

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"
	automation "github.com/stackitcloud/stackit-sdk-go/services/automation/v1betaapi"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	cliErr "github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/projectname"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/automation/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

const (
	executionIdArg   = "EXECUTION_ID"
	automationIdFlag = "automation-id"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	AutomationId string
	ExecutionId  string
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("describe %s", executionIdArg),
		Short: "Shows details of a Volume Automation Execution",
		Long:  "Shows details of a Volume Automation Execution.",
		Args:  args.SingleArg(executionIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Describe the Volume Automation Execution with ID "xxx" for Automation with ID "yyy"`,
				"$ stackit beta volume automation execution describe xxx --automation-id yyy"),
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
				return fmt.Errorf("describe volume automation execution: %w", err)
			}

			// Get projectLabel
			projectLabel, err := projectname.GetProjectName(ctx, params.Printer, params.CliVersion, cmd)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get project name: %v", err)
				projectLabel = model.ProjectId
			} else if projectLabel == "" {
				projectLabel = model.ProjectId
			}

			return outputResult(params.Printer, model.OutputFormat, model.ExecutionId, projectLabel, resp)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), automationIdFlag, "Automation ID")

	err := flags.MarkFlagsRequired(cmd, automationIdFlag)
	cobra.CheckErr(err)
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *automation.APIClient) automation.ApiGetVolumeExecutionRequest {
	return apiClient.DefaultAPI.GetVolumeExecution(ctx, model.ProjectId, model.Region, model.AutomationId, model.ExecutionId)
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	executionId := inputArgs[0]

	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &cliErr.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		AutomationId:    flags.FlagToStringValue(p, cmd, automationIdFlag),
		ExecutionId:     executionId,
	}

	p.DebugInputModel(model)
	return &model, nil
}

func outputResult(p *print.Printer, outputFormat, executionId, projectLabel string, execution *automation.VolumeExecutionResponse) error {
	return p.OutputResult(outputFormat, execution, func() error {
		if execution == nil {
			p.Outputf("Volume automation execution %q not found in project %q\n", executionId, projectLabel)
			return nil
		}

		var content []tables.Table

		baseTable := tables.NewTable()

		baseTable.SetTitle("Volume Automation Execution")
		baseTable.AddRow("ID", execution.Id)
		baseTable.AddSeparator()
		baseTable.AddRow("Status", execution.Status)
		baseTable.AddSeparator()
		baseTable.AddRow("CREATE TIME", execution.CreateTime.Format(time.DateTime))
		baseTable.AddSeparator()
		baseTable.AddRow("END TIME", utils.ConvertTimePToDateTimeString(execution.EndTime))
		baseTable.AddSeparator()

		content = append(content, baseTable)

		if execution.Automation != nil {
			automationConfigTable := tables.NewTable()
			automationConfigTable.SetTitle("Automation config")
			automationConfigTable.AddRow("NAME", utils.PtrString(execution.Automation.Config.Name))
			automationConfigTable.AddSeparator()
			automationConfigTable.AddRow("DESCRIPTION", utils.PtrString(execution.Automation.Config.Description))
			automationConfigTable.AddSeparator()
			automationConfigTable.AddRow("TEMPLATE ID", utils.PtrString(execution.Automation.Config.TemplateId))
			automationConfigTable.AddSeparator()
			automationConfigTable.AddRow("CREATE TIME", execution.Automation.Config.CreateTime.Format(time.DateTime))
			automationConfigTable.AddSeparator()
			automationConfigTable.AddRow("UPDATE TIME", execution.Automation.Config.UpdateTime.Format(time.DateTime))
			automationConfigTable.AddSeparator()

			content = append(content, automationConfigTable)

			if execution.Automation.Output != nil && len(execution.Automation.Output.Steps) > 0 {
				outputStepsTable := tables.NewTable()
				outputStepsTable.SetTitle("OUTPUT STEPS")
				outputStepsTable.SetHeader("NAME", "MESSAGE", "STATUS")
				for _, step := range execution.Automation.Output.Steps {
					outputStepsTable.AddRow(step.Name, utils.PtrString(step.Message), step.Status)
					outputStepsTable.AddSeparator()
				}
				content = append(content, outputStepsTable)
			}
		}

		err := tables.DisplayTables(p, content)
		if err != nil {
			return fmt.Errorf("render tables: %w", err)
		}
		return nil
	})
}
