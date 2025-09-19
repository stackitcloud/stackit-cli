package update

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	cliErr "github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/intake/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/intake"
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

func NewUpdateCmd(p *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("update %s", runnerIdArg),
		Short: "Updates an Intake Runner",
		Long:  "Updates an Intake Runner. Only the specified fields are updated.",
		Args:  args.SingleArg(runnerIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Update the display name of an Intake Runner with ID "xxx"`,
				`$ stackit intake runner update xxx --display-name "new-runner-name"`),
			examples.NewExample(
				`Update the message capacity limits for an Intake Runner with ID "xxx"`,
				`$ stackit intake runner update xxx --max-message-size-kib 2000 --max-messages-per-hour 10000`),
			examples.NewExample(
				`Clear the labels of an Intake Runner with ID "xxx" by providing an empty value`,
				`$ stackit intake runner update xxx --labels ""`),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			model, err := parseInput(p.Printer, cmd, args)
			if err != nil {
				return err
			}

			// Configure API client
			apiClient, err := client.ConfigureClient(p.Printer)
			if err != nil {
				return err
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			if err := req.Execute(); err != nil {
				return fmt.Errorf("update Intake Runner: %w", err)
			}

			p.Printer.Info("Update request for Intake Runner %q sent successfully.\n", model.RunnerId)
			return nil
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
	cmd.Flags().StringToString(labelFlag, nil, "Labels in key=value format. To clear all labels, provide an empty string, e.g. --labels \"\"")
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

	if p.IsVerbosityDebug() {
		modelStr, err := print.BuildDebugStrFromInputModel(model)
		if err != nil {
			p.Debug(print.ErrorLevel, "convert model to string for debugging: %v", err)
		} else {
			p.Debug(print.DebugLevel, "parsed input values: %s", modelStr)
		}
	}

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
