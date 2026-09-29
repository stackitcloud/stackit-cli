package delete

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	automation "github.com/stackitcloud/stackit-sdk-go/services/automation/v1api"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	cliErr "github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/projectname"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/automation/client"
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
		Use:   fmt.Sprintf("delete %s", automationIdArg),
		Short: "Deletes a Volume Automation",
		Long:  "Deletes a Volume Automation.",
		Args:  args.SingleArg(automationIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Delete a Volume Automation with ID "xxx"`,
				"$ stackit beta volume automation delete xxx"),
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

			projectLabel, err := projectname.GetProjectName(ctx, params.Printer, params.CliVersion, cmd)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get project name: %v", err)
				projectLabel = model.ProjectId
			} else if projectLabel == "" {
				projectLabel = model.ProjectId
			}

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to delete volume automation %q in project %q?", model.AutomationId, projectLabel)
				err = params.Printer.PromptForConfirmation(prompt)
				if err != nil {
					return err
				}
			}

			// Call API
			err = buildRequest(ctx, model, apiClient).Execute()
			if err != nil {
				return fmt.Errorf("delete volume automation: %w", err)
			}

			params.Printer.Outputf("Deleted volume automation %q\n", model.AutomationId)
			return nil
		},
	}
	return cmd
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *automation.APIClient) automation.ApiDeleteVolumeAutomationRequest {
	return apiClient.DefaultAPI.DeleteVolumeAutomation(ctx, model.ProjectId, model.Region, model.AutomationId)
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
