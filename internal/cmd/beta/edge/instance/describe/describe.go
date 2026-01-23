// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2025 STACKIT GmbH & Co. KG

package describe

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	cliErr "github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/edge/client"
	commonErr "github.com/stackitcloud/stackit-cli/internal/pkg/services/edge/common/error"
	commonInstance "github.com/stackitcloud/stackit-cli/internal/pkg/services/edge/common/instance"
	commonValidation "github.com/stackitcloud/stackit-cli/internal/pkg/services/edge/common/validation"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/edge"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	identifier *commonValidation.Identifier
}

// describeRequestSpec captures the details of the request for testing.
type describeRequestSpec struct {
	// Exported fields allow tests to inspect the request inputs
	ProjectID    string
	Region       string
	InstanceId   string // Set if describing by ID
	InstanceName string // Set if describing by Name

	// Execute is a closure that wraps the actual SDK call
	Execute func() (*edge.Instance, error)
}

// Command constructor
// Instance id and displayname are likely to be refactored in future. For the time being we decided to use flags
// instead of args to provide the instance-id xor displayname to uniquely identify an instance. The displayname
// is guaranteed to be unique within a given project as of today. The chosen flag over args approach ensures we
// won't need a breaking change of the CLI when we refactor the commands to take the identifier as arg at some point.
func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "describe",
		Short: "Describes an edge instance",
		Long:  "Describes a STACKIT Edge Cloud (STEC) instance.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				fmt.Sprintf(`Describe an edge instance with %s "xxx"`, commonInstance.InstanceIdFlag),
				fmt.Sprintf(`$ stackit beta edge-cloud instance describe --%s <ID>`, commonInstance.InstanceIdFlag)),
			examples.NewExample(
				fmt.Sprintf(`Describe an edge instance with %s "xxx"`, commonInstance.DisplayNameFlag),
				fmt.Sprintf(`$ stackit beta edge-cloud instance describe --%s <NAME>`, commonInstance.DisplayNameFlag)),
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

			// Call API
			resp, err := run(ctx, model, apiClient)
			if err != nil {
				return err
			}

			// Handle output to printer
			return outputResult(params.Printer, model.OutputFormat, resp)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().StringP(commonInstance.InstanceIdFlag, commonInstance.InstanceIdShorthand, "", commonInstance.InstanceIdUsage)
	cmd.Flags().StringP(commonInstance.DisplayNameFlag, commonInstance.DisplayNameShorthand, "", commonInstance.DisplayNameUsage)

	identifierFlags := []string{commonInstance.InstanceIdFlag, commonInstance.DisplayNameFlag}
	cmd.MarkFlagsMutuallyExclusive(identifierFlags...) // InstanceId xor DisplayName
	cmd.MarkFlagsOneRequired(identifierFlags...)
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
// It handles the logic of choosing between GetInstance and GetInstanceByName.
func buildRequest(ctx context.Context, model *inputModel, apiClient client.APIClient) (*describeRequestSpec, error) {
	if model == nil || model.identifier == nil {
		return nil, commonErr.NewNoIdentifierError("")
	}

	spec := &describeRequestSpec{
		ProjectID: model.ProjectId,
		Region:    model.Region,
	}

	// Switch the concrete client based on the identifier flag used
	switch model.identifier.Flag {
	case commonInstance.InstanceIdFlag:
		spec.InstanceId = model.identifier.Value
		req := apiClient.GetInstance(ctx, model.ProjectId, model.Region, model.identifier.Value)
		spec.Execute = req.Execute
	case commonInstance.DisplayNameFlag:
		spec.InstanceName = model.identifier.Value
		req := apiClient.GetInstanceByName(ctx, model.ProjectId, model.Region, model.identifier.Value)
		spec.Execute = req.Execute
	default:
		return nil, fmt.Errorf("%w: %w", cliErr.NewBuildRequestError("invalid identifier flag", nil), commonErr.NewInvalidIdentifierError(model.identifier.Flag))
	}

	return spec, nil
}

// Output result based on the configured output format
func outputResult(p *print.Printer, outputFormat string, instance *edge.Instance) error {
	if instance == nil {
		// This is only to prevent nil pointer deref.
		// As long as the API behaves as defined by it's spec, instance can not be empty (HTTP 200 with an empty body)
		return commonErr.NewNoInstanceError("")
	}

	return p.OutputResult(outputFormat, instance, func() error {
		table := tables.NewTable()
		// Describe: output all fields. Be sure to filter for any non-required fields.
		table.AddRow("CREATED", utils.PtrString(instance.Created))
		table.AddSeparator()
		table.AddRow("ID", utils.PtrString(instance.Id))
		table.AddSeparator()
		table.AddRow("NAME", utils.PtrString(instance.DisplayName))
		table.AddSeparator()
		if instance.HasDescription() {
			table.AddRow("DESCRIPTION", utils.PtrString(instance.Description))
			table.AddSeparator()
		}
		table.AddRow("UI", utils.PtrString(instance.FrontendUrl))
		table.AddSeparator()
		table.AddRow("STATE", utils.PtrString(instance.Status))
		table.AddSeparator()
		table.AddRow("PLAN", utils.PtrString(instance.PlanId))
		table.AddSeparator()

		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}
		return nil
	})
}
