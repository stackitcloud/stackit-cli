package list

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/goccy/go-yaml"
	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/projectname"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/intake/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-sdk-go/services/intake"
)

const (
	limitFlag = "limit"
)

// inputModel struct holds all the input parameters for the command
type inputModel struct {
	*globalflags.GlobalFlagModel
	Limit *int64
}

// NewListCmd creates a new cobra command for listing Intake Runners
func NewListCmd(p *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Lists all Intake Runners",
		Long:  "Lists all Intake Runners for the current project.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`List all Intake Runners`,
				`$ stackit beta intake runner list`),
			examples.NewExample(
				`List all Intake Runners in JSON format`,
				`$ stackit beta intake runner list --output-format json`),
			examples.NewExample(
				`List up to 5 Intake Runners`,
				`$ stackit beta intake runner list --limit 5`),
		),
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := context.Background()
			model, err := parseInput(p.Printer, cmd)
			if err != nil {
				return err
			}

			// Configure API client
			apiClient, err := client.ConfigureClient(p.Printer, p.CliVersion)
			if err != nil {
				return err
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("list Intake Runners: %w", err)
			}
			runners := resp.GetIntakeRunners()

			// Truncate output
			if model.Limit != nil && len(runners) > int(*model.Limit) {
				runners = runners[:*model.Limit]
			}

			projectLabel := model.ProjectId
			if len(runners) == 0 {
				projectLabel, err = projectname.GetProjectName(ctx, p.Printer, p.CliVersion, cmd)
				if err != nil {
					p.Printer.Debug(print.ErrorLevel, "get project name: %v", err)
				}
			}

			return outputResult(p.Printer, model.OutputFormat, projectLabel, runners)
		},
	}
	configureFlags(cmd)
	return cmd
}

// configureFlags adds the --limit flag to the command
func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Int64(limitFlag, 0, "Maximum number of entries to list")
}

// parseInput parses the command flags into a standardized model
func parseInput(p *print.Printer, cmd *cobra.Command) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	limit := flags.FlagToInt64Pointer(p, cmd, limitFlag)
	if limit != nil && *limit < 1 {
		return nil, &errors.FlagValidationError{
			Flag:    limitFlag,
			Details: "must be greater than 0",
		}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		Limit:           limit,
	}

	p.DebugInputModel(model)

	return &model, nil
}

// buildRequest creates the API request to list Intake Runners
func buildRequest(ctx context.Context, model *inputModel, apiClient *intake.APIClient) intake.ApiListIntakeRunnersRequest {
	req := apiClient.ListIntakeRunners(ctx, model.ProjectId, model.Region)
	// Note: we do support API pagination, but for consistency with other services, we fetch all items and apply
	// client-side limit.
	// A more advanced implementation could use the --limit flag to set the API's PageSize.
	return req
}

// outputResult formats the API response and prints it to the console
func outputResult(p *print.Printer, outputFormat, projectLabel string, runners []intake.IntakeRunnerResponse) error {
	switch outputFormat {
	case print.JSONOutputFormat:
		details, err := json.MarshalIndent(runners, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal Intake Runner list: %w", err)
		}
		p.Outputln(string(details))
		return nil

	case print.YAMLOutputFormat:
		details, err := yaml.MarshalWithOptions(runners, yaml.IndentSequence(true), yaml.UseJSONMarshaler())
		if err != nil {
			return fmt.Errorf("marshal Intake Runner list: %w", err)
		}
		p.Outputln(string(details))
		return nil

	default:
		if len(runners) == 0 {
			p.Outputf("No intake runners found for project %q\n", projectLabel)
			return nil
		}

		table := tables.NewTable()

		table.SetHeader("ID", "NAME")
		for _, runner := range runners {
			table.AddRow(
				runner.GetId(),
				runner.GetDisplayName(),
			)
		}
		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}
		return nil
	}
}
