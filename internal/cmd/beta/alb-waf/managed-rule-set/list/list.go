package list

import (
	"context"
	"fmt"
	"math"

	"github.com/stackitcloud/stackit-cli/internal/pkg/types"

	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/albwaf/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"

	"github.com/spf13/cobra"
	albwaf "github.com/stackitcloud/stackit-sdk-go/services/albwaf/v1api"
)

const (
	limitFlag   = "limit"
	maxPageSize = int64(100)
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	Limit *int64
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Lists all managed rule sets of the ALB WAF",
		Long:  "Lists all managed rule sets (MRS) of the Web Application Firewall (WAF) for application loadbalancers.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`List all managed rule sets`,
				"$ stackit beta alb-waf managed-rule-set list",
			),
			examples.NewExample(
				`List all managed rule sets in JSON format`,
				"$ stackit beta alb-waf managed-rule-set list --output-format json",
			),
			examples.NewExample(
				`List up to 10 managed rule sets`,
				"$ stackit beta alb-waf managed-rule-set list --limit 10",
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

			items, err := fetchManagedRuleSets(ctx, model, apiClient)
			if err != nil {
				return fmt.Errorf("list managed rule sets: %w", err)
			}

			return outputResult(params.Printer, model.OutputFormat, items)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Int64(limitFlag, 0, "Number of managed rule sets to list")
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

func buildRequest(ctx context.Context, model *inputModel, apiClient *albwaf.APIClient, nextPageID string, pageLimit int64) albwaf.ApiListManagedRuleSetsRequest {
	req := apiClient.DefaultAPI.ListManagedRuleSets(ctx, model.ProjectId, model.Region)
	req = req.PageSize(fmt.Sprintf("%d", pageLimit))
	if nextPageID != "" {
		req = req.PageId(nextPageID)
	}
	return req
}

// TODO: Replace with utils function within STACKITSDK-525
func fetchManagedRuleSets(ctx context.Context, model *inputModel, apiClient *albwaf.APIClient) ([]albwaf.GetLimitedManagedRuleSetResponse, error) {
	var nextPageID string
	var items []albwaf.GetLimitedManagedRuleSetResponse
	received := int64(0)
	limit := int64(math.MaxInt64)
	if model.Limit != nil {
		limit = min(limit, *model.Limit)
	}
	for {
		want := min(maxPageSize, limit-received)
		request := buildRequest(ctx, model, apiClient, nextPageID, want)
		response, err := request.Execute()
		if err != nil {
			return nil, err
		}
		items = append(items, response.Items...)
		nextPageID = ""
		if response.NextPageId != nil {
			nextPageID = *response.NextPageId
		}
		received += int64(len(response.Items))
		if nextPageID == "" || received >= limit {
			break
		}
	}
	return items, nil
}

func outputResult(p *print.Printer, outputFormat string, items []albwaf.GetLimitedManagedRuleSetResponse) error {
	return p.OutputResult(outputFormat, items, func() error {
		if len(items) == 0 {
			p.Outputln("No managed rule sets found")
			return nil
		}

		table := tables.NewTable()
		table.SetHeader("NAME", "TYPE", "VERSION", "USAGE COUNT")

		for _, item := range items {
			usageCount := ""
			if item.Usage != nil {
				usageCount = utils.PtrString(item.Usage.Count)
			}
			table.AddRow(
				item.Name,
				string(item.Type),
				item.Version,
				usageCount,
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
