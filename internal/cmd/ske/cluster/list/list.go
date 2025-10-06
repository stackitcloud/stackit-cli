package list

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/projectname"
	serviceEnablementClient "github.com/stackitcloud/stackit-cli/internal/pkg/services/service-enablement/client"
	serviceEnablementUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/service-enablement/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/ske/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

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

func NewCmd(params *params.CmdParams) *cobra.Command {
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
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := context.Background()
			model, err := parseInput(params.Printer, cmd)
			if err != nil {
				return err
			}

			// Configure API client
			apiClient, err := client.ConfigureClient(params.Printer, params.CliVersion)
			if err != nil {
				return err
			}

			// Configure ServiceEnable API client
			serviceEnablementApiClient, err := serviceEnablementClient.ConfigureClient(params.Printer, params.CliVersion)
			if err != nil {
				return err
			}

			// Check if SKE is enabled for this project
			enabled, err := serviceEnablementUtils.ProjectEnabled(ctx, serviceEnablementApiClient, model.ProjectId, model.Region)
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

			// Truncate output
			if model.Limit != nil && len(clusters) > int(*model.Limit) {
				clusters = clusters[:*model.Limit]
			}

			projectLabel := model.ProjectId
			if len(clusters) == 0 {
				projectLabel, err = projectname.GetProjectName(ctx, params.Printer, params.CliVersion, cmd)
				if err != nil {
					params.Printer.Debug(print.ErrorLevel, "get project name: %v", err)
				}
			}

			return outputResult(params.Printer, model.OutputFormat, projectLabel, clusters)
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

	model := inputModel{
		GlobalFlagModel: globalFlags,
		Limit:           flags.FlagToInt64Pointer(p, cmd, limitFlag),
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *ske.APIClient) ske.ApiListClustersRequest {
	req := apiClient.ListClusters(ctx, model.ProjectId, model.Region)
	return req
}

func outputResult(p *print.Printer, outputFormat, projectLabel string, clusters []ske.Cluster) error {
	return p.OutputResult(outputFormat, clusters, func() error {
		if len(clusters) == 0 {
			p.Outputf("No clusters found for project %q\n", projectLabel)
			return nil
		}

		table := tables.NewTable()
		table.SetHeader("NAME", "STATE", "VERSION", "POOLS", "MONITORING")
		for i := range clusters {
			c := clusters[i]
			monitoring := "Disabled"
			if c.Extensions != nil && c.Extensions.Observability != nil && *c.Extensions.Observability.Enabled {
				monitoring = "Enabled"
			}
			statusAggregated, kubernetesVersion := "", ""
			if c.HasStatus() {
				statusAggregated = utils.PtrString(c.Status.Aggregated)
			}
			if c.Kubernetes != nil {
				kubernetesVersion = utils.PtrString(c.Kubernetes.Version)
			}
			countNodepools := 0
			if c.Nodepools != nil {
				countNodepools = len(*c.Nodepools)
			}
			table.AddRow(
				utils.PtrString(c.Name),
				statusAggregated,
				kubernetesVersion,
				countNodepools,
				monitoring,
			)
		}
		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}

		return nil
	})
}
