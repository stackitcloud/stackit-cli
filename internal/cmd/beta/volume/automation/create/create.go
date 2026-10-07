package create

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	automation "github.com/stackitcloud/stackit-sdk-go/services/automation/v1api"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	cliErr "github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/projectname"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/automation/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
)

const (
	templateIdFlag           = "template-id"
	nameFlag                 = "name"
	descriptionFlag          = "description"
	inputFlag                = "input"
	triggerScheduleRruleFlag = "trigger-schedule-rrule"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	TemplateId           string
	Name                 *string
	Description          *string
	Input                *automation.VolumeAutomationInput
	TriggerScheduleRrule *string
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Creates a Volume Automation",
		Long: fmt.Sprintf("%s\n%s\n%s\n",
			"Creates a Volume Automation.",
			"The input for the automation can be provided as a JSON string or a file path prefixed with \"@\".",
			"See https://docs.api.stackit.cloud/documentation/automation-service/version/v1#tag/Volume-Automations/operation/CreateVolumeAutomation for information regarding the payload structure.",
		),
		Args: args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Creates a Volume Automation with name "my-automation", template ID "xxx" and input from the file "./input-payload.json"`,
				"$ stackit beta volume automation create --name my-automation --template-id xxx --input @./input-payload.json"),
			examples.NewExample(
				`Creates a Volume Automation with name "my-automation", description "CLI Example", template ID "xxx" and input from the file "./input-payload.json"`,
				"$ stackit beta volume automation create --name my-automation --description \"CLI Example\" --template-id xxx --input @./input-payload.json"),
			examples.NewExample(
				`Creates a Volume Automation with name "my-automation", template ID "xxx", rrule schedule trigger "FREQ=WEEKLY;BYDAY=MO,TU,WE,TH,FR" and input from the file "./input-payload.json"`,
				"$ stackit beta volume automation create --name my-automation --template-id xxx --trigger-schedule-rrule \"FREQ=WEEKLY;BYDAY=MO,TU,WE,TH,FR\" --input @./input-payload.json"),
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

			// Call API
			createdAutomation, err := buildRequest(ctx, model, apiClient).Execute()
			if err != nil {
				return fmt.Errorf("create volume automation: %w", err)
			}

			return outputResult(params.Printer, model.OutputFormat, projectLabel, createdAutomation)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), templateIdFlag, "Template ID which should be used for the automation")
	cmd.Flags().String(nameFlag, "", "Name of the automation")
	cmd.Flags().String(descriptionFlag, "", "Description of the automation")
	cmd.Flags().Var(flags.ReadFromFileFlag(), inputFlag, `Input for the automation (JSON). Can be a string or a file path, if prefixed with "@" (example: @./payload.json).`)
	cmd.Flags().String(triggerScheduleRruleFlag, "", "Trigger schedule (RRULE) for the automation")

	err := flags.MarkFlagsRequired(cmd, templateIdFlag)
	cobra.CheckErr(err)
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *automation.APIClient) automation.ApiCreateVolumeAutomationRequest {
	payload := automation.CreateVolumeAutomationPayload{
		Name:        *automation.NewNullableString(model.Name),
		Description: *automation.NewNullableString(model.Description),
		Input:       *automation.NewNullableVolumeAutomationInput(model.Input),
		TemplateId:  model.TemplateId,
	}
	if model.TriggerScheduleRrule != nil {
		payload.Triggers = *automation.NewNullableAutomationTriggers(&automation.AutomationTriggers{
			Schedule: *automation.NewNullableAutomationScheduleTrigger(&automation.AutomationScheduleTrigger{
				Rrule: *model.TriggerScheduleRrule,
			}),
		})
	}
	return apiClient.DefaultAPI.CreateVolumeAutomation(ctx, model.ProjectId, model.Region).CreateVolumeAutomationPayload(payload)
}

func parseInput(p *print.Printer, cmd *cobra.Command, _ []string) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &cliErr.ProjectIdError{}
	}

	var inputConfig *automation.VolumeAutomationInput
	if inputJson := flags.FlagToStringPointer(p, cmd, inputFlag); inputJson != nil {
		inputConfig = &automation.VolumeAutomationInput{}
		err := inputConfig.UnmarshalJSON([]byte(*inputJson))
		if err != nil {
			return nil, fmt.Errorf("the input from --%s can not be parsed: %w", inputFlag, err)
		}
	}

	model := inputModel{
		GlobalFlagModel:      globalFlags,
		TemplateId:           flags.FlagToStringValue(p, cmd, templateIdFlag),
		Name:                 flags.FlagToStringPointer(p, cmd, nameFlag),
		Description:          flags.FlagToStringPointer(p, cmd, descriptionFlag),
		Input:                inputConfig,
		TriggerScheduleRrule: flags.FlagToStringPointer(p, cmd, triggerScheduleRruleFlag),
	}

	p.DebugInputModel(model)
	return &model, nil
}

func outputResult(p *print.Printer, outputFormat, projectLabel string, item *automation.VolumeAutomation) error {
	if item == nil {
		return fmt.Errorf("volume automation response is empty")
	}
	return p.OutputResult(outputFormat, item, func() error {
		p.Outputf("Created volume automation for project %q. Automation ID: %q\n", projectLabel, item.Id)
		return nil
	})
}
