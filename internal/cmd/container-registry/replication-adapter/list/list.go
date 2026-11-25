package list

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	registry "github.com/stackitcloud/stackit-sdk-go/services/registry/v1api"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/projectname"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/container-registry/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
)

const (
	limitFlag = "limit"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	Limit *int64
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Lists all Container Registry replication adapters",
		Long:  "Lists all available adapters for replication to/from other registries.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`List all replication adapters`,
				"$ stackit container-registry replication-adapter list",
			),
			examples.NewExample(
				`List all replication adapters in JSON format`,
				"$ stackit container-registry replication-adapter list --output-format json",
			),
			examples.NewExample(
				`List up to 10 replication adapters`,
				"$ stackit container-registry replication-adapter list --limit 10",
			),
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
			adapters, err := fetchReplicationAdapters(ctx, model, apiClient)
			if err != nil {
				return err
			}

			projectLabel, err := projectname.GetProjectName(ctx, params.Printer, params.CliVersion, cmd)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get project name: %v", err)
				projectLabel = model.ProjectId
			}

			return outputResult(params.Printer, model.OutputFormat, projectLabel, adapters)
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

func buildRequest(ctx context.Context, model *inputModel, apiClient *registry.APIClient) registry.ApiListReplicationAdaptersRequest {
	return apiClient.DefaultAPI.ListReplicationAdapters(ctx, model.ProjectId, model.Region)
}

func fetchReplicationAdapters(ctx context.Context, model *inputModel, apiClient *registry.APIClient) ([]registry.ReplicationAdapter, error) {
	// Note: the ListReplicationAdapters endpoint does not accept page_token/page_size
	// query parameters (the SDK only exposes Execute()), so the next page token of the
	// response cannot be forwarded. The list of adapter types is small and always
	// fits in a single page.
	resp, err := buildRequest(ctx, model, apiClient).Execute()
	if err != nil {
		return nil, fmt.Errorf("list replication adapters: %w", err)
	}
	if resp == nil {
		return nil, fmt.Errorf("list replication adapters: empty response")
	}

	adapters := resp.GetReplicationAdapters()
	if adapters == nil {
		adapters = []registry.ReplicationAdapter{}
	}
	if model.Limit != nil && len(adapters) > int(*model.Limit) {
		adapters = adapters[:*model.Limit]
	}

	return adapters, nil
}

func outputResult(p *print.Printer, outputFormat, projectLabel string, adapters []registry.ReplicationAdapter) error {
	return p.OutputResult(outputFormat, adapters, func() error {
		if len(adapters) == 0 {
			p.Outputf("No replication adapters found for project %q\n", projectLabel)
			return nil
		}

		table := tables.NewTable()
		table.SetHeader("REGISTRY TYPE", "DISPLAY NAME", "CREDENTIAL TYPE")

		for _, adapter := range adapters {
			table.AddRow(
				adapter.RegistryType,
				adapter.DisplayName,
				adapter.CredentialType,
			)
		}

		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}

		return nil
	})
}
