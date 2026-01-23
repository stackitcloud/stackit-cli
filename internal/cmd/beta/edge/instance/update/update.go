// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2025 STACKIT GmbH & Co. KG

package update

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
	commonValidation "github.com/stackitcloud/stackit-cli/internal/pkg/services/edge/common/validation"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"

	"github.com/stackitcloud/stackit-cli/internal/pkg/spinner"
	"github.com/stackitcloud/stackit-sdk-go/services/edge"
	"github.com/stackitcloud/stackit-sdk-go/services/edge/wait"
)

// Struct to model user input (arguments and/or flags)
type inputModel struct {
	*globalflags.GlobalFlagModel
	identifier  *commonValidation.Identifier
	Description *string
	PlanId      *string
}

// updateRequestSpec captures the details of the request for testing.
type updateRequestSpec struct {
	// Exported fields allow tests to inspect the request inputs
	ProjectID    string
	Region       string
	InstanceId   string // Set if updating by ID
	InstanceName string // Set if updating by Name
	Payload      edge.UpdateInstancePayload

	// Execute is a closure that wraps the actual SDK call
	Execute func() error
}

// OpenApi generated code will have different types for by-instance-id and by-display-name API calls and therefore different wait handlers.
// InstanceWaiter is an interface to abstract the different wait handlers so they can be used interchangeably.
type instanceWaiter interface {
	WaitWithContext(context.Context) (*edge.Instance, error)
}

// A function that creates an instance waiter
type instanceWaiterFactory = func(client *edge.APIClient) instanceWaiter

// Command constructor
// Instance id and displayname are likely to be refactored in future. For the time being we decided to use flags
// instead of args to provide the instance-id xor displayname to uniquely identify an instance. The displayname
// is guaranteed to be unique within a given project as of today. The chosen flag over args approach ensures we
// won't need a breaking change of the CLI when we refactor the commands to take the identifier as arg at some point.
func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Updates an edge instance",
		Long:  "Updates a STACKIT Edge Cloud (STEC) instance.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				fmt.Sprintf(`Updates the description of an edge instance with %s "xxx"`, commonInstance.InstanceIdFlag),
				fmt.Sprintf(`$ stackit beta edge-cloud instance update --%s "xxx" --%s "yyy"`, commonInstance.InstanceIdFlag, commonInstance.DescriptionFlag)),
			examples.NewExample(
				fmt.Sprintf(`Updates the plan of an edge instance with %s "xxx"`, commonInstance.DisplayNameFlag),
				fmt.Sprintf(`$ stackit beta edge-cloud instance update --%s "xxx" --%s "yyy"`, commonInstance.DisplayNameFlag, commonInstance.PlanIdFlag)),
			examples.NewExample(
				fmt.Sprintf(`Updates the description and plan of an edge instance with %s "xxx"`, commonInstance.InstanceIdFlag),
				fmt.Sprintf(`$ stackit beta edge-cloud instance update --%s "xxx" --%s "yyy" --%s "zzz"`, commonInstance.InstanceIdFlag, commonInstance.DescriptionFlag, commonInstance.PlanIdFlag)),
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
				prompt := fmt.Sprintf("Are you sure you want to update the edge instance %q of project %q?", model.identifier.Value, projectLabel)
				err = params.Printer.PromptForConfirmation(prompt)
				if err != nil {
					return err
				}
			}

			// Call API
			err = run(ctx, model, apiClient)
			if err != nil {
				return err
			}

			// Wait for async operation, if async mode not enabled
			operationState := "Triggered update of"
			if !model.Async {
				// Wait for async operation, if async mode not enabled
				// Show spinner while waiting
				s := spinner.New(params.Printer)
				s.Start("Updating instance")
				// Determine identifier and waiter to use
				waiterFactory, err := getWaiterFactory(ctx, model)
				if err != nil {
					return err
				}
				// The waiter handler needs a concrete client type. We can safely cast here as the real implementation will always match.
				client, ok := apiClient.(*edge.APIClient)
				if !ok {
					return fmt.Errorf("failed to configure API client")
				}
				waiter := waiterFactory(client)

				if _, err = waiter.WaitWithContext(ctx); err != nil {
					return fmt.Errorf("wait for edge instance update: %w", err)
				}
				operationState = "Updated"
				s.Stop()
			}

			params.Printer.Info("%s instance with %q %q of project %q.\n", operationState, model.identifier.Flag, model.identifier.Value, projectLabel)

			return nil
		},
	}

	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().StringP(commonInstance.InstanceIdFlag, commonInstance.InstanceIdShorthand, "", commonInstance.InstanceIdUsage)
	cmd.Flags().StringP(commonInstance.DisplayNameFlag, commonInstance.DisplayNameShorthand, "", commonInstance.DisplayNameUsage)
	cmd.Flags().StringP(commonInstance.DescriptionFlag, commonInstance.DescriptionShorthand, "", commonInstance.DescriptionUsage)
	cmd.Flags().StringP(commonInstance.PlanIdFlag, "", "", commonInstance.PlanIdUsage)

	identifierFlags := []string{commonInstance.InstanceIdFlag, commonInstance.DisplayNameFlag}
	cmd.MarkFlagsMutuallyExclusive(identifierFlags...) // InstanceId xor DisplayName
	cmd.MarkFlagsOneRequired(identifierFlags...)

	// Make sure at least one updatable field is provided, otherwise it would be a no-op
	updatedFields := []string{commonInstance.DescriptionFlag, commonInstance.PlanIdFlag}
	cmd.MarkFlagsOneRequired(updatedFields...)
}

// Parse user input (arguments and/or flags)
func parseInput(p *print.Printer, cmd *cobra.Command) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &cliErr.ProjectIdError{}
	}

	// Generate input model based on chosen flags
	model := inputModel{
		GlobalFlagModel: globalFlags,
	}

	// Parse and validate user input then add it to the model
	id, err := commonValidation.GetValidatedInstanceIdentifier(p, cmd)
	if err != nil {
		return nil, err
	}
	model.identifier = id

	if planIdValue := flags.FlagToStringPointer(p, cmd, commonInstance.PlanIdFlag); planIdValue != nil {
		if err := commonInstance.ValidatePlanId(planIdValue); err != nil {
			return nil, err
		}
		model.PlanId = planIdValue
	}

	if descriptionValue := flags.FlagToStringPointer(p, cmd, commonInstance.DescriptionFlag); descriptionValue != nil {
		if err := commonInstance.ValidateDescription(*descriptionValue); err != nil {
			return nil, err
		}
		model.Description = descriptionValue
	}

	// Log the parsed model if --verbosity is set to debug
	p.DebugInputModel(model)
	return &model, nil
}

// Run is the main execution function used by the command runner.
// It is decoupled from TTY output to have the ability to mock the API client during testing.
func run(ctx context.Context, model *inputModel, apiClient client.APIClient) error {
	spec, err := buildRequest(ctx, model, apiClient)
	if err != nil {
		return err
	}

	err = spec.Execute()
	if err != nil {
		return cliErr.NewRequestFailedError(err)
	}

	return nil
}

// buildRequest constructs the spec that can be tested.
// It handles the logic of choosing between UpdateInstance and UpdateInstanceByName.
func buildRequest(ctx context.Context, model *inputModel, apiClient client.APIClient) (*updateRequestSpec, error) {
	if model == nil || model.identifier == nil {
		return nil, commonErr.NewNoIdentifierError("")
	}

	spec := &updateRequestSpec{
		ProjectID: model.ProjectId,
		Region:    model.Region,
		Payload: edge.UpdateInstancePayload{
			Description: model.Description,
			PlanId:      model.PlanId,
		},
	}

	// Switch the concrete client based on the identifier flag used
	switch model.identifier.Flag {
	case commonInstance.InstanceIdFlag:
		spec.InstanceId = model.identifier.Value
		req := apiClient.UpdateInstance(ctx, model.ProjectId, model.Region, model.identifier.Value)
		req = req.UpdateInstancePayload(spec.Payload)
		spec.Execute = req.Execute
	case commonInstance.DisplayNameFlag:
		spec.InstanceName = model.identifier.Value
		req := apiClient.UpdateInstanceByName(ctx, model.ProjectId, model.Region, model.identifier.Value)
		req = req.UpdateInstanceByNamePayload(edge.UpdateInstanceByNamePayload{
			Description: spec.Payload.Description,
			PlanId:      spec.Payload.PlanId,
		})
		spec.Execute = req.Execute
	default:
		return nil, fmt.Errorf("%w: %w", cliErr.NewBuildRequestError("invalid identifier flag", nil), commonErr.NewInvalidIdentifierError(model.identifier.Flag))
	}

	return spec, nil
}

// Returns a factory function to create the appropriate waiter based on the input model.
func getWaiterFactory(ctx context.Context, model *inputModel) (instanceWaiterFactory, error) {
	if model == nil || model.identifier == nil {
		return nil, commonErr.NewNoIdentifierError("")
	}

	switch model.identifier.Flag {
	case commonInstance.InstanceIdFlag:
		factory := func(c *edge.APIClient) instanceWaiter {
			return wait.CreateOrUpdateInstanceWaitHandler(ctx, c, model.ProjectId, model.Region, model.identifier.Value)
		}
		return factory, nil
	case commonInstance.DisplayNameFlag:
		factory := func(c *edge.APIClient) instanceWaiter {
			return wait.CreateOrUpdateInstanceByNameWaitHandler(ctx, c, model.ProjectId, model.Region, model.identifier.Value)
		}
		return factory, nil
	default:
		return nil, commonErr.NewInvalidIdentifierError(model.identifier.Flag)
	}
}
