package list

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/core/experimental/paginate"
	automation "github.com/stackitcloud/stackit-sdk-go/services/automation/v1betaapi"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	cliErr "github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/projectname"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/automation/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
)

const (
	limitFlag = "limit"

	// maxPageSize is the maximum number of items the API returns per page.
	maxPageSize = 100
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	Limit *int64
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all Volume Automation Templates",
		Long:  "List all Volume Automation Templates.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`List all Volume Automation Templates`,
				"$ stackit beta volume automation template list"),
			examples.NewExample(
				`List up to 10 Volume Automation Templates`,
				"$ stackit beta volume automation template list --limit 10"),
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

			// Call API
			resp, err := fetchTemplates(ctx, model, apiClient)
			if err != nil {
				return fmt.Errorf("list volume automation templates: %w", err)
			}

			// Get projectLabel
			projectLabel, err := projectname.GetProjectName(ctx, params.Printer, params.CliVersion, cmd)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get project name: %v", err)
				projectLabel = model.ProjectId
			} else if projectLabel == "" {
				projectLabel = model.ProjectId
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

func parseInput(p *print.Printer, cmd *cobra.Command, _ []string) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &cliErr.ProjectIdError{}
	}

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

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *automation.APIClient) automation.ApiListVolumeTemplatesRequest {
	return apiClient.DefaultAPI.ListVolumeTemplates(ctx, model.ProjectId, model.Region)
}

func fetchTemplates(ctx context.Context, model *inputModel, apiClient *automation.APIClient) ([]automation.Template, error) {
	req := buildRequest(ctx, model, apiClient)
	opts := []paginate.Option{paginate.WithPageSize(maxPageSize)}
	if model.Limit != nil {
		opts = append(opts, paginate.WithLimit(int(*model.Limit)))
	}

	items, err := paginate.All[automation.Template](req, opts...)
	if err != nil {
		return nil, fmt.Errorf("list volume automation templates: %w", err)
	}
	if items == nil {
		items = []automation.Template{}
	}

	return items, nil
}

func outputResult(p *print.Printer, outputFormat, projectLabel string, templates []automation.Template) error {
	return p.OutputResult(outputFormat, templates, func() error {
		if len(templates) == 0 {
			p.Outputf("No volume automation templates found in project %q\n", projectLabel)
			return nil
		}

		table := tables.NewTable()

		table.SetTitle("Volume Automation Templates")
		table.SetHeader("ID", "NAME", "DESCRIPTION")
		table.AddSeparator()
		for _, template := range templates {
			table.AddRow(
				template.Id,
				template.Name,
				template.Description,
			)
			table.AddSeparator()
		}

		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("render tables: %w", err)
		}
		return nil
	})
}
