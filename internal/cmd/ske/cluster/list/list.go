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
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/ske/client"
	skeUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/ske/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/ske"
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
		Short: "Lists all SKE clusters",
		Long:  "Lists all STACKIT Kubernetes Engine (SKE) clusters.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`List all SKE clusters`,
				"$ stackit ske cluster list"),
			examples.NewExample(
				`List all SKE clusters in JSON format`,
				"$ stackit ske cluster list --output-format json"),
			examples.NewExample(
				`List up to 10 SKE clusters`,
				"$ stackit ske cluster list --limit 10"),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			model, err := parseInput(p, cmd)
			if err != nil {
				return err
			}

			// Configure API client
			apiClient, err := client.ConfigureClient(p)
			if err != nil {
				return err
			}

			// Check if SKE is enabled for this project
			enabled, err := skeUtils.ProjectEnabled(ctx, apiClient, model.ProjectId)
			if err != nil {
				return err
			}
			if !enabled {
				return fmt.Errorf("SKE isn't enabled for this project, please run 'stackit ske enable'")
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("get SKE clusters: %w", err)
			}
			clusters := *resp.Items
			if len(clusters) == 0 {
				projectLabel, err := projectname.GetProjectName(ctx, p, cmd)
				if err != nil {
					p.Debug(print.ErrorLevel, "get project name: %v", err)
					projectLabel = model.ProjectId
				}
				p.Info("No clusters found for project %q\n", projectLabel)
				return nil
			}

			// Truncate output
			if model.Limit != nil && len(clusters) > int(*model.Limit) {
				clusters = clusters[:*model.Limit]
			}

			return outputResult(p, model.OutputFormat, clusters)
		},
	}

	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Int64(limitFlag, 0, "Maximum number of entries to list")
}

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

	return &inputModel{
		GlobalFlagModel: globalFlags,
		Limit:           flags.FlagToInt64Pointer(p, cmd, limitFlag),
	}, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *ske.APIClient) ske.ApiListClustersRequest {
	req := apiClient.ListClusters(ctx, model.ProjectId)
	return req
}

func outputResult(p *print.Printer, outputFormat string, clusters []ske.Cluster) error {
	switch outputFormat {
	case globalflags.JSONOutputFormat:
		details, err := json.MarshalIndent(clusters, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal SKE cluster list: %w", err)
		}
		p.Outputln(string(details))

		return nil
	default:
		table := tables.NewTable()
		table.SetHeader("NAME", "STATE", "VERSION", "POOLS", "MONITORING")
		for i := range clusters {
			c := clusters[i]
			monitoring := "Disabled"
			if c.Extensions != nil && c.Extensions.Argus != nil && *c.Extensions.Argus.Enabled {
				monitoring = "Enabled"
			}
			table.AddRow(*c.Name, *c.Status.Aggregated, *c.Kubernetes.Version, len(*c.Nodepools), monitoring)
		}
		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}

		return nil
	}
}
