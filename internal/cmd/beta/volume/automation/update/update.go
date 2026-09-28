package update

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
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

const (
	automationIdArg = "AUTOMATION_ID"

	nameFlag                 = "name"
	descriptionFlag          = "description"
	inputFlag                = "input"
	triggerScheduleRruleFlag = "trigger-schedule-rrule"
	disableTriggerSchedule   = "disable-trigger-schedule"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	AutomationId           string
	Name                   *string
	Description            *string
	Input                  *automation.VolumeAutomationInput
	TriggerScheduleRrule   *string
	DisableTriggerSchedule *bool
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("update %s", automationIdArg),
		Short: "Updates a Volume Automation",
		Long: fmt.Sprintf("%s\n%s\n%s\n%s\n",
			"Updates a Volume Automation.",
			"The input for the automation can be provided as a JSON string or a file path prefixed with \"@\".",
			"Updates are always applied as a patch. Omitted fields in the JSON will be ignored and remain unchanged after the update. Fields can be removed by setting them explicitly to null.",
			"See https://docs.api.stackit.cloud/documentation/automation-service/version/v1#tag/Volume-Automations/operation/PartialUpdateVolumeAutomation for information regarding the payload structure.",
		),
		Args: args.SingleArg(automationIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Update volume automation with ID "xxx" with input from the file "./input-payload.json"`,
				"$ stackit beta volume automation update xxx --input @./input-payload.json"),
			examples.NewExample(
				`Update volume automation with ID "xxx" with name "my-updated-automation"`,
				"$ stackit beta volume automation update xxx --name \"my-updated-automation\""),
			examples.NewExample(
				`Update volume automation with ID "xxx" with description "Updated automation from STACKIT CLI"`,
				"$ stackit beta volume automation update xxx --description \"Updated automation from STACKIT CLI\""),
			examples.NewExample(
				`Update volume automation with ID "xxx" with disabling trigger schedule rrule`,
				"$ stackit beta volume automation update xxx --disable-trigger-schedule"),
			examples.NewExample(
				`Update volume automation with ID "xxx" with trigger schedule rrule "FREQ=WEEKLY;BYDAY=MO"`,
				"$ stackit beta volume automation update xxx --trigger-schedule-rrule \"FREQ=WEEKLY;BYDAY=MO\""),
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
			updatedAutomation, err := buildRequest(ctx, model, apiClient).Execute()
			if err != nil {
				return fmt.Errorf("update volume automation: %w", err)
			}

			return outputResult(params.Printer, model.OutputFormat, projectLabel, updatedAutomation)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().String(nameFlag, "", "Name of the automation")
	cmd.Flags().String(descriptionFlag, "", "Description of the automation")
	cmd.Flags().Var(flags.ReadFromFileFlag(), inputFlag, `Input for the automation (JSON). Can be a string or a file path, if prefixed with "@" (example: @./payload.json).`)
	cmd.Flags().String(triggerScheduleRruleFlag, "", "Trigger schedule (RRULE) for the automation")
	cmd.Flags().Bool(disableTriggerSchedule, false, "If set, trigger schedule will be disabled")

	cmd.MarkFlagsOneRequired(nameFlag, descriptionFlag, inputFlag, triggerScheduleRruleFlag, disableTriggerSchedule)
	cmd.MarkFlagsMutuallyExclusive(triggerScheduleRruleFlag, disableTriggerSchedule)
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *automation.APIClient) automation.ApiPartialUpdateVolumeAutomationRequest {
	// The update payload fields are nullable, therefore we need to check if they are set in the model.
	// Otherwise, the NullableString would be potentially set with nil as value, and they would be unwanted removed in the update
	payload := automation.PartialUpdateVolumeAutomationPayload{}
	if model.Name != nil {
		payload.Name = *automation.NewNullableString(model.Name)
	}
	if model.Description != nil {
		payload.Description = *automation.NewNullableString(model.Description)
	}
	if model.Input != nil {
		payload.Input = *automation.NewNullableVolumeAutomationInput(model.Input)
	}
	if model.DisableTriggerSchedule != nil && *model.DisableTriggerSchedule {
		payload.Triggers = *automation.NewNullableAutomationTriggers(&automation.AutomationTriggers{
			Schedule: *automation.NewNullableAutomationScheduleTrigger(nil),
		})
	} else if model.TriggerScheduleRrule != nil {
		payload.Triggers = *automation.NewNullableAutomationTriggers(&automation.AutomationTriggers{
			Schedule: *automation.NewNullableAutomationScheduleTrigger(&automation.AutomationScheduleTrigger{
				Rrule: *model.TriggerScheduleRrule,
			}),
		})
	}
	return apiClient.DefaultAPI.PartialUpdateVolumeAutomation(ctx, model.ProjectId, model.Region, model.AutomationId).PartialUpdateVolumeAutomationPayload(payload)
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	automationId := inputArgs[0]

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
		GlobalFlagModel:        globalFlags,
		AutomationId:           automationId,
		Name:                   flags.FlagToStringPointer(p, cmd, nameFlag),
		Description:            flags.FlagToStringPointer(p, cmd, descriptionFlag),
		Input:                  inputConfig,
		TriggerScheduleRrule:   flags.FlagToStringPointer(p, cmd, triggerScheduleRruleFlag),
		DisableTriggerSchedule: flags.FlagToBoolPointer(p, cmd, disableTriggerSchedule),
	}

	p.DebugInputModel(model)
	return &model, nil
}

func outputResult(p *print.Printer, outputFormat, projectLabel string, item *automation.VolumeAutomation) error {
	if item == nil {
		return fmt.Errorf("volume automation response is empty")
	}
	return p.OutputResult(outputFormat, item, func() error {
		p.Outputf("Updated volume automation for project %q. Automation ID: %q\n", projectLabel, item.Id)
		return nil
	})
}
