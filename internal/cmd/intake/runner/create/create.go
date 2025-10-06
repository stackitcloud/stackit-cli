package create

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/goccy/go-yaml"
	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/intake"

	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	cliErr "github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/projectname"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/intake/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

const (
	displayNameFlag        = "display-name"
	maxMessageSizeKiBFlag  = "max-message-size-kib"
	maxMessagesPerHourFlag = "max-messages-per-hour"
	descriptionFlag        = "description"
	labelsFlag             = "labels"
)

// inputModel struct holds all the input parameters for the command
type inputModel struct {
	*globalflags.GlobalFlagModel
	DisplayName        *string
	MaxMessageSizeKiB  *int64
	MaxMessagesPerHour *int64
	Description        *string
	Labels             *map[string]string
}

func NewCreateCmd(p *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Creates a new Intake Runner",
		Long:  "Creates a new Intake Runner.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Create a new Intake Runner with a display name and message capacity limits`,
				`$ stackit intake runner create --display-name my-runner --max-message-size-kib 1000 --max-messages-per-hour 5000`),
			examples.NewExample(
				`Create a new Intake Runner with a description and labels`,
				`$ stackit intake runner create --display-name my-runner --max-message-size-kib 1000 --max-messages-per-hour 5000 --description "Main runner for production" --labels="env=prod,team=billing"`),
		),
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := context.Background()
			model, err := parseInput(p.Printer, cmd)
			if err != nil {
				return err
			}

			// Configure API client
			apiClient, err := client.ConfigureClient(p.Printer)
			if err != nil {
				return err
			}

			projectLabel, err := projectname.GetProjectName(ctx, p.Printer, p.CliVersion, cmd)
			if err != nil {
				p.Printer.Debug(print.ErrorLevel, "get project name: %v", err)
				projectLabel = model.ProjectId
			}

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to create an Intake Runner for project %q?", projectLabel)
				err = p.Printer.PromptForConfirmation(prompt)
				if err != nil {
					return err
				}
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("create Intake Runner: %w", err)
			}

			return outputResult(p.Printer, model.OutputFormat, projectLabel, resp)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().String(displayNameFlag, "", "Display name")
	cmd.Flags().Int64(maxMessageSizeKiBFlag, 0, "Maximum message size in KiB")
	cmd.Flags().Int64(maxMessagesPerHourFlag, 0, "Maximum number of messages per hour")
	cmd.Flags().String(descriptionFlag, "", "Description")
	cmd.Flags().StringToString(labelsFlag, nil, "Labels in key=value format, separated by commas. Example: --labels \"key1=value1,key2=value2\"")

	err := flags.MarkFlagsRequired(cmd, displayNameFlag, maxMessageSizeKiBFlag, maxMessagesPerHourFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &cliErr.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel:    globalFlags,
		DisplayName:        flags.FlagToStringPointer(p, cmd, displayNameFlag),
		MaxMessageSizeKiB:  flags.FlagToInt64Pointer(p, cmd, maxMessageSizeKiBFlag),
		MaxMessagesPerHour: flags.FlagToInt64Pointer(p, cmd, maxMessagesPerHourFlag),
		Description:        flags.FlagToStringPointer(p, cmd, descriptionFlag),
		Labels:             flags.FlagToStringToStringPointer(p, cmd, labelsFlag),
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

func buildRequest(ctx context.Context, model *inputModel, apiClient *intake.APIClient) intake.ApiCreateIntakeRunnerRequest {
	// Start building the request by calling the base method with path parameters
	req := apiClient.CreateIntakeRunner(ctx, model.ProjectId, model.Region)

	// Create the payload struct with data from the input model
	payload := intake.CreateIntakeRunnerPayload{
		DisplayName:        model.DisplayName,
		MaxMessageSizeKiB:  model.MaxMessageSizeKiB,
		MaxMessagesPerHour: model.MaxMessagesPerHour,
		Description:        model.Description,
		Labels:             model.Labels,
	}
	// Attach the payload to the request builder
	req = req.CreateIntakeRunnerPayload(payload)

	return req
}

func outputResult(p *print.Printer, outputFormat, projectLabel string, resp *intake.IntakeRunnerResponse) error {
	switch outputFormat {
	case print.JSONOutputFormat:
		details, err := json.MarshalIndent(resp, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal instance: %w", err)
		}
		p.Outputln(string(details))

		return nil
	case print.YAMLOutputFormat:
		details, err := yaml.MarshalWithOptions(resp, yaml.IndentSequence(true), yaml.UseJSONMarshaler())
		if err != nil {
			return fmt.Errorf("marshal instance: %w", err)
		}
		p.Outputln(string(details))

		return nil
	default:
		p.Outputf("Created Intake Runner for project %q. Runner ID: %s\n", projectLabel, utils.PtrString(resp.Id))
		return nil
	}
}
