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
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/postgresflex/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/postgresflex"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
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
		Short: "Lists all PostgreSQL Flex instances",
		Long:  "Lists all PostgreSQL Flex instances.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`List all PostgreSQL Flex instances`,
				"$ stackit postgresflex instance list"),
			examples.NewExample(
				`List all PostgreSQL Flex instances in JSON format`,
				"$ stackit postgresflex instance list --output-format json"),
			examples.NewExample(
				`List up to 10 PostgreSQL Flex instances`,
				"$ stackit postgresflex instance list --limit 10"),
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
				return fmt.Errorf("get PostgreSQL Flex instances: %w", err)
			}
			if resp.Items == nil || len(*resp.Items) == 0 {
				projectLabel, err := projectname.GetProjectName(ctx, cmd, p)
				if err != nil {
					p.Debug(print.ErrorLevel, "get project name: %v", err)
					projectLabel = model.ProjectId
				}
				p.Info("No instances found for project %q\n", projectLabel)
				return nil
			}
			instances := *resp.Items

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

func buildRequest(ctx context.Context, model *inputModel, apiClient *postgresflex.APIClient) postgresflex.ApiListInstancesRequest {
	req := apiClient.ListInstances(ctx, model.ProjectId)
	return req
}

func outputResult(p *print.Printer, outputFormat string, instances []postgresflex.InstanceListInstance) error {
	switch outputFormat {
	case globalflags.JSONOutputFormat:
		details, err := json.MarshalIndent(instances, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal PostgreSQL Flex instance list: %w", err)
		}
		p.Outputln(string(details))

		return nil
	default:
		caser := cases.Title(language.English)
		table := tables.NewTable()
		table.SetHeader("ID", "NAME", "STATUS")
		for i := range instances {
			instance := instances[i]
			table.AddRow(*instance.Id, *instance.Name, caser.String(*instance.Status))
		}
		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}

		return nil
	}
}
