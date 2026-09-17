package describe

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	automation "github.com/stackitcloud/stackit-sdk-go/services/automation/v1betaapi"

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
	templateIdArg = "TEMPLATE_ID"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	TemplateId string
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "describe",
		Short: "Shows details of a Volume Automation Template",
		Long:  "Shows details of a Volume Automation Template.",
		Args:  args.SingleArg(templateIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Describe the Volume Automation Template with ID "xxx"`,
				"$ stackit beta volume automation template describe xxx"),
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
				return fmt.Errorf("describe volume automation template: %w", err)
			}

			// Get projectLabel
			projectLabel, err := projectname.GetProjectName(ctx, params.Printer, params.CliVersion, cmd)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get project name: %v", err)
				projectLabel = model.ProjectId
			} else if projectLabel == "" {
				projectLabel = model.ProjectId
			}

			return outputResult(params.Printer, model.OutputFormat, model.TemplateId, projectLabel, resp)
		},
	}
	return cmd
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *automation.APIClient) automation.ApiGetVolumeTemplateRequest {
	req := apiClient.DefaultAPI.GetVolumeTemplate(ctx, model.ProjectId, model.Region, model.TemplateId)
	return req
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	templateId := inputArgs[0]

	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &cliErr.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		TemplateId:      templateId,
	}

	p.DebugInputModel(model)
	return &model, nil
}

func outputResult(p *print.Printer, outputFormat, templateId, projectLabel string, template *automation.GetVolumeTemplateResponse) error {
	return p.OutputResult(outputFormat, template, func() error {
		if template == nil {
			p.Outputf("Volume automation template %q not found in project %q\n", templateId, projectLabel)
			return nil
		}

		var content []tables.Table

		table := tables.NewTable()

		table.SetTitle("Volume Automation Template")
		table.AddRow("ID", template.Id)
		table.AddSeparator()
		table.AddRow("NAME", template.Name)
		table.AddSeparator()
		table.AddRow("DESCRIPTION", template.Description)
		table.AddSeparator()
		if template.Input != nil {
			table.AddRow("INPUT KIND", template.Input.Kind)
			table.AddSeparator()
		}

		content = append(content, table)

		if template.Output != nil && len(template.Output.Steps) > 0 {
			outputStepsTable := tables.NewTable()
			outputStepsTable.SetTitle("OUTPUT STEPS")
			for _, step := range template.Output.Steps {
				outputStepsTable.AddRow("NAME", step.Name)
				outputStepsTable.AddSeparator()
				if step.Result != nil {
					if step.Result.GetVolumeIDsResult != nil {
						volumeIds := ""
						if len(step.Result.GetVolumeIDsResult.VolumeIDs) > 0 {
							volumeIds = strings.Join(step.Result.GetVolumeIDsResult.VolumeIDs, "\n")
						}
						outputStepsTable.AddRow("VOLUME IDS", volumeIds)
						outputStepsTable.AddSeparator()
					}
					if step.Result.CreateSnapshotsResult != nil {
						createdSnapshotIds := ""
						if len(step.Result.CreateSnapshotsResult.CreatedSnapshotIDs) > 0 {
							createdSnapshotIds = strings.Join(step.Result.CreateSnapshotsResult.CreatedSnapshotIDs, "\n")
						}
						outputStepsTable.AddRow("CREATED SNAPSHOT IDS", createdSnapshotIds)
						outputStepsTable.AddSeparator()
					}
				}
			}
			content = append(content, outputStepsTable)
		}

		err := tables.DisplayTables(p, content)
		if err != nil {
			return fmt.Errorf("render tables: %w", err)
		}
		return nil
	})
}
