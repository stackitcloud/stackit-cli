package list

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/projectname"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/argus/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/argus"
)

const (
	limitFlag = "limit"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	Limit *int64
}

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Lists all Argus instances",
		Long:  "Lists all Argus instances.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`List all Argus instances`,
				"$ stackit argus instance list"),
			examples.NewExample(
				`List all Argus instances in JSON format`,
				"$ stackit argus instance list --output-format json"),
			examples.NewExample(
				`List up to 10 Argus instances`,
				"$ stackit argus instance list --limit 10"),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			model, err := parseInput(cmd, p)
			if err != nil {
				return err
			}

			// Configure API client
			apiClient, err := client.ConfigureClient(p)
			if err != nil {
				return err
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("get Argus instances: %w", err)
			}
			instances := *resp.Instances
			if len(instances) == 0 {
				projectLabel, err := projectname.GetProjectName(ctx, cmd, p)
				if err != nil {
					p.Debug(print.ErrorLevel, "get project name: %v", err)
					projectLabel = model.ProjectId
				}
				p.Info("No instances found for project %q\n", projectLabel)
				return nil
			}

			// Truncate output
			if model.Limit != nil && len(instances) > int(*model.Limit) {
				instances = instances[:*model.Limit]
			}

			return outputResult(p, model.OutputFormat, instances)
		},
	}

	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Int64(limitFlag, 0, "Maximum number of entries to list")
}

func parseInput(cmd *cobra.Command, p *print.Printer) (*inputModel, error) {
	globalFlags := globalflags.Parse(cmd, p)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	limit := flags.FlagToInt64Pointer(cmd, limitFlag, p)
	if limit != nil && *limit < 1 {
		return nil, &errors.FlagValidationError{
			Flag:    limitFlag,
			Details: "must be greater than 0",
		}
	}

	return &inputModel{
		GlobalFlagModel: globalFlags,
		Limit:           flags.FlagToInt64Pointer(cmd, limitFlag, p),
	}, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *argus.APIClient) argus.ApiListInstancesRequest {
	req := apiClient.ListInstances(ctx, model.ProjectId)
	return req
}

func outputResult(p *print.Printer, outputFormat string, instances []argus.ProjectInstanceFull) error {
	switch outputFormat {
	case globalflags.JSONOutputFormat:
		details, err := json.MarshalIndent(instances, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal Argus instance list: %w", err)
		}
		p.Outputln(string(details))

		return nil
	default:
		table := tables.NewTable()
		table.SetHeader("ID", "NAME", "PLAN", "STATUS")
		for i := range instances {
			instance := instances[i]
			table.AddRow(*instance.Id, *instance.Name, *instance.PlanName, *instance.Status)
		}
		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}

		return nil
	}
}
