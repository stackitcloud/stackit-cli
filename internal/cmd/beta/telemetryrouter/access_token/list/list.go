package list

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/types"

	"github.com/stackitcloud/stackit-sdk-go/core/experimental/paginate"
	telemetryrouter "github.com/stackitcloud/stackit-sdk-go/services/telemetryrouter/v1api"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	cliErr "github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/telemetryrouter/client"
	telemetryrouterUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/telemetryrouter/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

const (
	limitFlag      = "limit"
	instanceIdFlag = "instance-id"

	// maxPageSize is the maximum number of items the API returns per page.
	maxPageSize = 100
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	InstanceId string
	Limit      *int64
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Lists all TelemetryRouter access tokens of an instance",
		Long:  "Lists all TelemetryRouter access tokens of an instance.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Lists all access tokens of the TelemetryRouter instance "xxx"`,
				"$ stackit beta telemetryrouter access-token list --instance-id xxx",
			),
			examples.NewExample(
				`Lists all access tokens in JSON format`,
				"$ stackit beta telemetryrouter access-token list --instance-id xxx --output-format json",
			),
			examples.NewExample(
				`Lists up to 10 access tokens`,
				"$ stackit beta telemetryrouter access-token list --instance-id xxx --limit 10",
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

			instanceLabel, err := telemetryrouterUtils.GetInstanceName(ctx, apiClient.DefaultAPI, model.ProjectId, model.Region, model.InstanceId)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get instance name: %v", err)
				instanceLabel = model.InstanceId
			}

			items, err := fetchAccessTokens(ctx, model, apiClient)
			if err != nil {
				return err
			}

			return outputResult(params.Printer, model.OutputFormat, items, instanceLabel)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Int64(limitFlag, 0, "Maximum number of entries to list")
	cmd.Flags().Var(flags.UUIDFlag(), instanceIdFlag, "ID of the TelemetryRouter instance")

	err := flags.MarkFlagsRequired(cmd, instanceIdFlag)
	cobra.CheckErr(err)
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
		InstanceId:      flags.FlagToStringValue(p, cmd, instanceIdFlag),
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *telemetryrouter.APIClient) telemetryrouter.ApiListAccessTokensRequest {
	return apiClient.DefaultAPI.ListAccessTokens(ctx, model.ProjectId, model.Region, model.InstanceId)
}

// listAccessTokensResponse adapts telemetryrouter.ListAccessTokensResponse to paginate.Response,
// since the generated type names its list field after the resource (AIP-158) instead of "Items".
type listAccessTokensResponse struct {
	*telemetryrouter.ListAccessTokensResponse
}

func (r listAccessTokensResponse) GetItems() []telemetryrouter.GetAccessTokenResponse {
	return r.GetAccessTokens()
}

// listAccessTokensAdapter adapts telemetryrouter.ApiListAccessTokensRequest to paginate.Request.
type listAccessTokensAdapter struct {
	request telemetryrouter.ApiListAccessTokensRequest
}

func (a *listAccessTokensAdapter) PageSize(pageSize int32) *listAccessTokensAdapter {
	a.request = a.request.PageSize(pageSize)
	return a
}

func (a *listAccessTokensAdapter) PageToken(pageToken string) *listAccessTokensAdapter {
	a.request = a.request.PageToken(pageToken)
	return a
}

func (a *listAccessTokensAdapter) Execute() (listAccessTokensResponse, error) {
	resp, err := a.request.Execute()
	if err != nil {
		return listAccessTokensResponse{}, err
	}
	return listAccessTokensResponse{resp}, nil
}

func fetchAccessTokens(ctx context.Context, model *inputModel, apiClient *telemetryrouter.APIClient) ([]telemetryrouter.GetAccessTokenResponse, error) {
	opts := []paginate.Option{paginate.WithPageSize(maxPageSize)}
	if model.Limit != nil {
		opts = append(opts, paginate.WithLimit(int(*model.Limit))) // nolint:gosec // bounded by prior validation, cannot overflow int
	}

	request := &listAccessTokensAdapter{request: buildRequest(ctx, model, apiClient)}
	items, err := paginate.All(request, opts...)
	if err != nil {
		return nil, fmt.Errorf("list TelemetryRouter access tokens: %w", err)
	}
	if items == nil {
		items = []telemetryrouter.GetAccessTokenResponse{}
	}
	return items, nil
}

func outputResult(p *print.Printer, outputFormat string, tokens []telemetryrouter.GetAccessTokenResponse, instanceLabel string) error {
	return p.OutputResult(outputFormat, tokens, func() error {
		if len(tokens) == 0 {
			p.Outputf("No access tokens found for TelemetryRouter instance %q\n", instanceLabel)
			return nil
		}

		table := tables.NewTable()
		table.SetHeader("ID", "DISPLAY NAME", "STATUS", "CREATOR ID", "EXPIRATION TIME")

		for _, token := range tokens {
			table.AddRow(
				token.Id,
				token.DisplayName,
				token.Status,
				token.CreatorId,
				utils.PtrString(token.ExpirationTime.Get()),
			)
			table.AddSeparator()
		}

		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}
		return nil
	})
}
