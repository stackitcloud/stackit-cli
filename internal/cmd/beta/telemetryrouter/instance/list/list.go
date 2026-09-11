package list

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/types"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/core/experimental/paginate"
	telemetryrouter "github.com/stackitcloud/stackit-sdk-go/services/telemetryrouter/v1api"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/projectname"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/telemetryrouter/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	Limit *int64
}

const (
	limitFlag = "limit"

	// maxPageSize is the maximum number of items the API returns per page.
	maxPageSize = 100
)

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Lists TelemetryRouter instances",
		Long:  "Lists TelemetryRouter instances within the project.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`List all TelemetryRouter instances`,
				`$ stackit beta telemetryrouter instance list`,
			),
			examples.NewExample(
				`List the first 10 TelemetryRouter instances`,
				`$ stackit beta telemetryrouter instance list --limit=10`,
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

			projectLabel, err := projectname.GetProjectName(ctx, params.Printer, params.CliVersion, cmd)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get project name: %v", err)
				projectLabel = model.ProjectId
			}
			if projectLabel == "" {
				projectLabel = model.ProjectId
			}

			items, err := fetchInstances(ctx, model, apiClient)
			if err != nil {
				return err
			}

			return outputResult(params.Printer, model.OutputFormat, projectLabel, items)
		},
	}

	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Int64(limitFlag, 0, "Limit the output to the first n elements")
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

func buildRequest(ctx context.Context, model *inputModel, apiClient *telemetryrouter.APIClient) telemetryrouter.ApiListTelemetryRoutersRequest {
	return apiClient.DefaultAPI.ListTelemetryRouters(ctx, model.ProjectId, model.Region)
}

// listTelemetryRoutersResponse adapts telemetryrouter.ListTelemetryRoutersResponse to
// paginate.Response, since the generated type names its list field after the resource
// (AIP-158) instead of "Items".
type listTelemetryRoutersResponse struct {
	*telemetryrouter.ListTelemetryRoutersResponse
}

func (r listTelemetryRoutersResponse) GetItems() []telemetryrouter.TelemetryRouterResponse {
	return r.GetTelemetryRouters()
}

// listTelemetryRoutersAdapter adapts telemetryrouter.ApiListTelemetryRoutersRequest to paginate.Request.
type listTelemetryRoutersAdapter struct {
	request telemetryrouter.ApiListTelemetryRoutersRequest
}

func (a *listTelemetryRoutersAdapter) PageSize(pageSize int32) *listTelemetryRoutersAdapter {
	a.request = a.request.PageSize(pageSize)
	return a
}

func (a *listTelemetryRoutersAdapter) PageToken(pageToken string) *listTelemetryRoutersAdapter {
	a.request = a.request.PageToken(pageToken)
	return a
}

func (a *listTelemetryRoutersAdapter) Execute() (listTelemetryRoutersResponse, error) {
	resp, err := a.request.Execute()
	if err != nil {
		return listTelemetryRoutersResponse{}, err
	}
	return listTelemetryRoutersResponse{resp}, nil
}

func fetchInstances(ctx context.Context, model *inputModel, apiClient *telemetryrouter.APIClient) ([]telemetryrouter.TelemetryRouterResponse, error) {
	opts := []paginate.Option{paginate.WithPageSize(maxPageSize)}
	if model.Limit != nil {
		opts = append(opts, paginate.WithLimit(int(*model.Limit))) // nolint:gosec // bounded by prior validation, cannot overflow int
	}

	request := &listTelemetryRoutersAdapter{request: buildRequest(ctx, model, apiClient)}
	items, err := paginate.All(request, opts...)
	if err != nil {
		return nil, fmt.Errorf("list TelemetryRouter instances: %w", err)
	}
	if items == nil {
		items = []telemetryrouter.TelemetryRouterResponse{}
	}
	return items, nil
}

func outputResult(p *print.Printer, outputFormat, projectLabel string, instances []telemetryrouter.TelemetryRouterResponse) error {
	return p.OutputResult(outputFormat, instances, func() error {
		if len(instances) == 0 {
			p.Outputf("No TelemetryRouter instances found for project %q\n", projectLabel)
			return nil
		}

		table := tables.NewTable()
		table.SetHeader("ID", "DISPLAY NAME", "URI")
		for _, instance := range instances {
			table.AddRow(
				instance.Id,
				instance.DisplayName,
				instance.Uri,
			)
		}
		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}

		return nil
	})
}
