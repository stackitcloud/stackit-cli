package update

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/types"

	"github.com/spf13/cobra"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	cliErr "github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/projectname"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/intake/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/spinner"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/intake"
	"github.com/stackitcloud/stackit-sdk-go/services/intake/wait"
)

const (
	runnerIdArg = "RUNNER_ID"
)

const (
	displayNameFlag        = "display-name"
	maxMessageSizeKiBFlag  = "max-message-size-kib"
	maxMessagesPerHourFlag = "max-messages-per-hour"
	descriptionFlag        = "description"
	labelFlag              = "labels"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	RunnerId           string
	DisplayName        *string
	MaxMessageSizeKiB  *int64
	MaxMessagesPerHour *int64
	Description        *string
	Labels             *map[string]string
}

func NewCmd(p *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("update %s", runnerIdArg),
		Short: "Updates an Intake Runner",
		Long:  "Updates an Intake Runner. Only the specified fields are updated.",
		Args:  args.SingleArg(runnerIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Update the display name of an Intake Runner with ID "xxx"`,
				`$ stackit beta intake runner update xxx --display-name "new-runner-name"`),
			examples.NewExample(
				`Update the message capacity limits for an Intake Runner with ID "xxx"`,
				`$ stackit beta intake runner update xxx --max-message-size-kib 1000 --max-messages-per-hour 10000`),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			model, err := parseInput(p.Printer, cmd, args)
			if err != nil {
				return err
			}

			// Configure API client
			apiClient, err := client.ConfigureClient(p.Printer, p.CliVersion)
			if err != nil {
				return err
			}

			projectLabel, err := projectname.GetProjectName(ctx, p.Printer, p.CliVersion, cmd)
			if err != nil {
				p.Printer.Debug(print.ErrorLevel, "get project name: %v", err)
				projectLabel = model.ProjectId
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("update Intake Runner: %w", err)
			}

			// Wait for async operation, if async mode not enabled
			if !model.Async {
				s := spinner.New(p.Printer)
				s.Start("Updating STACKIT Intake Runner")
				_, err = wait.CreateOrUpdateIntakeRunnerWaitHandler(ctx, apiClient, model.ProjectId, model.Region, model.RunnerId).WaitWithContext(ctx)
				if err != nil {
					return fmt.Errorf("wait for STACKIT Intake Runner update: %w", err)
				}
				s.Stop()
			}

			return outputResult(p.Printer, model, projectLabel, resp)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().String(displayNameFlag, "", "Display name")
	cmd.Flags().Int64(maxMessageSizeKiBFlag, 0, "Maximum message size in KiB. Note: Overall message capacity cannot be decreased.")
	cmd.Flags().Int64(maxMessagesPerHourFlag, 0, "Maximum number of messages per hour. Note: Overall message capacity cannot be decreased.")
	cmd.Flags().String(descriptionFlag, "", "Description")
	cmd.Flags().StringToString(labelFlag, nil, `Labels in key=value format, separated by commas. Example: --labels "key1=value1,key2=value2".`)
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	runnerId := inputArgs[0]

	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &cliErr.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel:    globalFlags,
		RunnerId:           runnerId,
		DisplayName:        flags.FlagToStringPointer(p, cmd, displayNameFlag),
		MaxMessageSizeKiB:  flags.FlagToInt64Pointer(p, cmd, maxMessageSizeKiBFlag),
		MaxMessagesPerHour: flags.FlagToInt64Pointer(p, cmd, maxMessagesPerHourFlag),
		Description:        flags.FlagToStringPointer(p, cmd, descriptionFlag),
		Labels:             flags.FlagToStringToStringPointer(p, cmd, labelFlag),
	}

	if model.DisplayName == nil && model.MaxMessageSizeKiB == nil && model.MaxMessagesPerHour == nil && model.Description == nil && model.Labels == nil {
		return nil, &cliErr.EmptyUpdateError{}
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *intake.APIClient) intake.ApiUpdateIntakeRunnerRequest {
	req := apiClient.UpdateIntakeRunner(ctx, model.ProjectId, model.Region, model.RunnerId)

	payload := intake.UpdateIntakeRunnerPayload{}
	if model.DisplayName != nil {
		payload.DisplayName = model.DisplayName
	}
	if model.MaxMessageSizeKiB != nil {
		payload.MaxMessageSizeKiB = model.MaxMessageSizeKiB
	}
	if model.MaxMessagesPerHour != nil {
		payload.MaxMessagesPerHour = model.MaxMessagesPerHour
	}
	if model.Description != nil {
		payload.Description = model.Description
	}
	if model.Labels != nil {
		payload.Labels = model.Labels
	}

	req = req.UpdateIntakeRunnerPayload(payload)
	return req
}

func outputResult(p *print.Printer, model *inputModel, projectLabel string, resp *intake.IntakeRunnerResponse) error {
	return p.OutputResult(model.OutputFormat, resp, func() error {
		if resp == nil {
			p.Outputf("Triggered update of Intake Runner for project %q, but no runner ID was returned.\n", projectLabel)
			return nil
		}

		operationState := "Updated"
		if model.Async {
			operationState = "Triggered update of"
		}
		p.Outputf("%s Intake Runner for project %q. Runner ID: %s\n", operationState, projectLabel, utils.PtrString(resp.Id))
		return nil
	})
}
