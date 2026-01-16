// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2025 STACKIT GmbH & Co. KG

package create

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	cliErr "github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/projectname"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/edge/client"
	commonErr "github.com/stackitcloud/stackit-cli/internal/pkg/services/edge/common/error"
	commonInstance "github.com/stackitcloud/stackit-cli/internal/pkg/services/edge/common/instance"
	"github.com/stackitcloud/stackit-cli/internal/pkg/spinner"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/edge"
	"github.com/stackitcloud/stackit-sdk-go/services/edge/wait"
)

// Command constructor
// Instance id and displayname are likely to be refactored in future. For the time being we decided to use flags
// instead of args to provide the instance-id xor displayname to uniquely identify an instance. The displayname
// is guaranteed to be unique within a given project as of today. The chosen flag over args approach ensures we
// won't need a breaking change of the CLI when we refactor the commands to take the identifier as arg at some point.
func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Creates an edge instance",
		Long:  "Creates a STACKIT Edge Cloud (STEC) instance. The instance will take a moment to become fully functional.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				fmt.Sprintf(`Creates an edge instance with the %s "xxx" and %s "yyy"`, commonInstance.DisplayNameFlag, commonInstance.PlanIdFlag),
				fmt.Sprintf(`$ stackit beta edge-cloud instance create --%s "xxx" --%s "yyy"`, commonInstance.DisplayNameFlag, commonInstance.PlanIdFlag)),
		),
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := context.Background()

			// Parse user input (arguments and/or flags)
			model, err := parseInput(params.Printer, cmd)
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
				// If project label can't be determined, fall back to project ID
				projectLabel = model.ProjectId
			}

			// Prompt for confirmation
			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to create a new edge instance for project %q?", projectLabel)
				err = params.Printer.PromptForConfirmation(prompt)
				if err != nil {
					return err
				}
			}

			// Call API
			resp, err := run(ctx, model, apiClient)
			if err != nil {
				return err
			}
			if resp == nil {
				return fmt.Errorf("create instance: empty response from API")
			}
			if resp.Id == nil {
				return fmt.Errorf("create instance: instance id missing in response")
			}
			instanceId := *resp.Id

			// Wait for async operation, if async mode not enabled
			if !model.Async {
				s := spinner.New(params.Printer)
				s.Start("Creating instance")
				// The waiter handler needs a concrete client type. We can safely cast here as the real implementation will always match.
				client, ok := apiClient.(*edge.APIClient)
				if !ok {
					return fmt.Errorf("failed to configure API client")
				}
				_, err = wait.CreateOrUpdateInstanceWaitHandler(ctx, client, model.ProjectId, model.Region, instanceId).WaitWithContext(ctx)

				if err != nil {
					return fmt.Errorf("wait for edge instance creation: %w", err)
				}
				s.Stop()
			}

			// Handle output to printer
			return outputResult(params.Printer, model.OutputFormat, model.Async, projectLabel, resp)
		},
	}

	configureFlags(cmd)
	return cmd
}

// inputModel represents the user input for creating an edge instance.
type inputModel struct {
	*globalflags.GlobalFlagModel
	DisplayName string
	Description string
	PlanId      string
}

// createRequestSpec captures the details of the request for testing.
type createRequestSpec struct {
	// Exported fields allow tests to inspect the request inputs
	ProjectID string
	Region    string
	Payload   edge.CreateInstancePayload

	// Execute is a closure that wraps the actual SDK call
	Execute func() (*edge.Instance, error)
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().StringP(commonInstance.DisplayNameFlag, commonInstance.DisplayNameShorthand, "", commonInstance.DisplayNameUsage)
	cmd.Flags().StringP(commonInstance.DescriptionFlag, commonInstance.DescriptionShorthand, "", commonInstance.DescriptionUsage)
	cmd.Flags().String(commonInstance.PlanIdFlag, "", commonInstance.PlanIdUsage)

	cobra.CheckErr(flags.MarkFlagsRequired(cmd, commonInstance.DisplayNameFlag, commonInstance.PlanIdFlag))
}

// Parse user input (arguments and/or flags)
func parseInput(p *print.Printer, cmd *cobra.Command) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &cliErr.ProjectIdError{}
	}

	// Parse and validate user input then add it to the model
	displayNameValue := flags.FlagToStringPointer(p, cmd, commonInstance.DisplayNameFlag)
	if err := commonInstance.ValidateDisplayName(displayNameValue); err != nil {
		return nil, err
	}

	planIdValue := flags.FlagToStringPointer(p, cmd, commonInstance.PlanIdFlag)
	if err := commonInstance.ValidatePlanId(planIdValue); err != nil {
		return nil, err
	}

	descriptionValue := flags.FlagWithDefaultToStringValue(p, cmd, commonInstance.DescriptionFlag)
	if err := commonInstance.ValidateDescription(descriptionValue); err != nil {
		return nil, err
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		DisplayName:     *displayNameValue,
		Description:     descriptionValue,
		PlanId:          *planIdValue,
	}

	// Log the parsed model if --verbosity is set to debug
	p.DebugInputModel(model)
	return &model, nil
}

// Run is the main execution function used by the command runner.
// It is decoupled from TTY output to have the ability to mock the API client during testing.
func run(ctx context.Context, model *inputModel, apiClient client.APIClient) (*edge.Instance, error) {
	spec, err := buildRequest(ctx, model, apiClient)
	if err != nil {
		return nil, err
	}

	resp, err := spec.Execute()
	if err != nil {
		return nil, cliErr.NewRequestFailedError(err)
	}

	return resp, nil
}

// buildRequest constructs the spec that can be tested.
func buildRequest(ctx context.Context, model *inputModel, apiClient client.APIClient) (*createRequestSpec, error) {
	req := apiClient.CreateInstance(ctx, model.ProjectId, model.Region)

	// Build request payload
	payload := edge.CreateInstancePayload{
		DisplayName: &model.DisplayName,
		Description: &model.Description,
		PlanId:      &model.PlanId,
	}
	req = req.CreateInstancePayload(payload)

	return &createRequestSpec{
		ProjectID: model.ProjectId,
		Region:    model.Region,
		Payload:   payload,
		Execute:   req.Execute,
	}, nil
}

// Output result based on the configured output format
func outputResult(p *print.Printer, outputFormat string, async bool, projectLabel string, instance *edge.Instance) error {
	if instance == nil {
		// This is only to prevent nil pointer deref.
		// As long as the API behaves as defined by it's spec, instance can not be empty (HTTP 200 with an empty body)
		return commonErr.NewNoInstanceError("")
	}

	return p.OutputResult(outputFormat, instance, func() error {
		operationState := "Created"
		if async {
			operationState = "Triggered creation of"
		}
		p.Outputf("%s instance for project %q. Instance ID: %q.\n", operationState, projectLabel, utils.PtrString(instance.Id))
		return nil
	})
}
