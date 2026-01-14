// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2025 STACKIT GmbH & Co. KG

package list

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
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/edge"
)

const (
	limitFlag = "limit"
)

// Struct to model user input (arguments and/or flags)
type inputModel struct {
	*globalflags.GlobalFlagModel
	Limit *int64
}

// listRequestSpec captures the details of the request for testing.
type listRequestSpec struct {
	// Exported fields allow tests to inspect the request inputs
	ProjectID string
	Region    string
	Limit     *int64

	// Execute is a closure that wraps the actual SDK call
	Execute func() (*edge.InstanceList, error)
}

// Command constructor
func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Lists edge instances",
		Long:  "Lists STACKIT Edge Cloud (STEC) instances of a project.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Lists all edge instances of a given project`,
				`$ stackit beta edge-cloud instance list`),
			examples.NewExample(
				`Lists all edge instances of a given project and limits the output to two instances`,
				fmt.Sprintf(`$ stackit beta edge-cloud instance list --%s 2`, limitFlag)),
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

			// Call API
			resp, err := run(ctx, model, apiClient)
			if err != nil {
				return err
			}

			return outputResult(params.Printer, model.OutputFormat, projectLabel, resp)
		},
	}

	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Int64(limitFlag, 0, "Maximum number of entries to list")
}

// Parse user input (arguments and/or flags)
func parseInput(p *print.Printer, cmd *cobra.Command) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &cliErr.ProjectIdError{}
	}

	// Parse and validate user input then add it to the model
	limit := flags.FlagToInt64Pointer(p, cmd, limitFlag)
	if limit != nil && *limit < 1 {
		return nil, &cliErr.FlagValidationError{
			Flag:    limitFlag,
			Details: "must be greater than 0",
		}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		Limit:           limit,
	}

	// Log the parsed model if --verbosity is set to debug
	p.DebugInputModel(model)
	return &model, nil
}

// Run is the main execution function used by the command runner.
// It is decoupled from TTY output to have the ability to mock the API client during testing.
func run(ctx context.Context, model *inputModel, apiClient client.APIClient) ([]edge.Instance, error) {
	spec, err := buildRequest(ctx, model, apiClient)
	if err != nil {
		return nil, err
	}

	resp, err := spec.Execute()
	if err != nil {
		return nil, cliErr.NewRequestFailedError(err)
	}
	if resp == nil {
		return nil, fmt.Errorf("list instances: empty response from API")
	}
	if resp.Instances == nil {
		return nil, fmt.Errorf("list instances: instances missing in response")
	}
	instances := *resp.Instances

	// Truncate output if limit is set
	if spec.Limit != nil && len(instances) > int(*spec.Limit) {
		instances = instances[:*spec.Limit]
	}

	return instances, nil
}

// buildRequest constructs the spec that can be tested.
func buildRequest(ctx context.Context, model *inputModel, apiClient client.APIClient) (*listRequestSpec, error) {
	req := apiClient.GetInstances(ctx, model.ProjectId, model.Region)

	return &listRequestSpec{
		ProjectID: model.ProjectId,
		Region:    model.Region,
		Limit:     model.Limit,
		Execute:   req.Execute,
	}, nil
}

// Output result based on the configured output format
func outputResult(p *print.Printer, outputFormat, projectLabel string, instances []edge.Instance) error {
	return p.OutputResult(outputFormat, instances, func() error {
		// No instances found for project
		if len(instances) == 0 {
			p.Outputf("No instances found for project %q\n", projectLabel)
			return nil
		}

		// Display instances found for project in a table
		table := tables.NewTable()
		// List: only output the most important fields. Be sure to filter for any non-required fields.
		table.SetHeader("ID", "NAME", "UI", "STATE")
		for i := range instances {
			instance := instances[i]
			table.AddRow(
				utils.PtrString(instance.Id),
				utils.PtrString(instance.DisplayName),
				utils.PtrString(instance.FrontendUrl),
				utils.PtrString(instance.Status))
		}
		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}
		return nil
	})
}
