package list

import (
	"context"
	"fmt"
	"math"
	"strconv"

	"github.com/stackitcloud/stackit-cli/internal/pkg/types"

	"github.com/spf13/cobra"
	albwaf "github.com/stackitcloud/stackit-sdk-go/services/albwaf/v1api"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/projectname"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/albwaf/client"
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
		Short: "Lists all ALB WAF custom rule groups",
		Long:  "Lists all STACKIT Application Load Balancer (ALB) Web Application Firewall (WAF) custom rule groups.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`List all ALB WAF custom rule groups`,
				`$ stackit beta alb-waf custom-rule-group list`,
			),
			examples.NewExample(
				`List the first 10 ALB WAF custom rule groups`,
				`$ stackit beta alb-waf custom-rule-group list --limit=10`,
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
			} else if projectLabel == "" {
				projectLabel = model.ProjectId
			}

			items, err := fetchCustomRuleGroups(ctx, model, apiClient)
			if err != nil {
				return fmt.Errorf("list ALB WAF custom rule groups: %w", err)
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

func buildRequest(ctx context.Context, model *inputModel, apiClient *albwaf.APIClient, pageId string, pageSize int64) albwaf.ApiListCustomRuleGroupRequest {
	req := apiClient.DefaultAPI.ListCustomRuleGroup(ctx, model.ProjectId, model.Region)
	req = req.PageSize(strconv.FormatInt(pageSize, 10))
	if pageId != "" {
		req = req.PageId(pageId)
	}
	return req
}

func fetchCustomRuleGroups(ctx context.Context, model *inputModel, apiClient *albwaf.APIClient) ([]albwaf.GetCustomRuleGroupResponse, error) {
	var pageId string
	var items []albwaf.GetCustomRuleGroupResponse
	received := int64(0)
	limit := int64(math.MaxInt64)
	if model.Limit != nil {
		limit = *model.Limit
	}
	for {
		want := min(int64(maxPageSize), limit-received)
		request := buildRequest(ctx, model, apiClient, pageId, want)
		response, err := request.Execute()
		if err != nil {
			return nil, fmt.Errorf("list custom rule groups: %w", err)
		}
		if response.Items != nil {
			items = append(items, response.Items...)
		}
		pageId = ""
		if response.NextPageId != nil {
			pageId = *response.NextPageId
		}
		received += want
		if pageId == "" || received >= limit {
			break
		}
	}
	return items, nil
}

func outputResult(p *print.Printer, outputFormat, projectLabel string, items []albwaf.GetCustomRuleGroupResponse) error {
	return p.OutputResult(outputFormat, items, func() error {
		if len(items) == 0 {
			p.Outputf("No custom rule groups found for project %q\n", projectLabel)
			return nil
		}

		table := tables.NewTable()
		table.SetHeader("NAME", "RULES", "USED BY")
		for i := range items {
			item := &items[i]

			var usedBy int
			if item.Usage != nil && item.Usage.Count != nil {
				usedBy = int(*item.Usage.Count)
			}

			table.AddRow(
				item.Name,
				len(item.Rules),
				usedBy,
			)
		}
		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}

		return nil
	})
}
