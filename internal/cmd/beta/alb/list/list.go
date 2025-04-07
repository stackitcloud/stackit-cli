package list

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/goccy/go-yaml"
	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/projectname"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/alb/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/alb"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	Limit *int64
}

const (
	labelSelectorFlag = "label-selector"
	limitFlag         = "limit"
)

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Lists albs",
		Long:  "Lists application load balancers.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`List all load balancers`,
				`$ stackit beta alb list`,
			),
			examples.NewExample(
				`List the first 10 application load balancers`,
				`$ stackit beta alb list --limit=10`,
			),
		),
		RunE: func(cmd *cobra.Command, _ []string) error {
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

			projectLabel, err := projectname.GetProjectName(ctx, p, cmd)
			if err != nil {
				p.Debug(print.ErrorLevel, "get project name: %v", err)
				projectLabel = model.ProjectId
			} else if projectLabel == "" {
				projectLabel = model.ProjectId
			}

			// Call API
			request := buildRequest(ctx, model, apiClient)

			response, err := request.Execute()
			if err != nil {
				return fmt.Errorf("list load balancerse: %w", err)
			}

			if items := response.LoadBalancers; items != nil && len(*items) == 0 {
				p.Info("No load balancers found for project %q", projectLabel)
			} else {
				if model.Limit != nil && len(*items) > int(*model.Limit) {
					*items = (*items)[:*model.Limit]
				}
				if err := outputResult(p, model.OutputFormat, *items); err != nil {
					return fmt.Errorf("output loadbalancers: %w", err)
				}
			}

			return nil
		},
	}

	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Int64(limitFlag, 0, "Limit the output to the first n elements")
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
		Limit:           limit,
	}

	if p.IsVerbosityDebug() {
		modelStr, err := print.BuildDebugStrFromInputModel(model)
		if err != nil {
			p.Debug(print.ErrorLevel, "convert model to string for debugging: %v", err)
		} else {
			p.Debug(print.DebugLevel, "parsed input values: %s", modelStr)
		}
	}

	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *alb.APIClient) alb.ApiListLoadBalancersRequest {
	request := apiClient.ListLoadBalancers(ctx, model.ProjectId, model.Region)

	return request
}
func outputResult(p *print.Printer, outputFormat string, items []alb.LoadBalancer) error {
	switch outputFormat {
	case print.JSONOutputFormat:
		details, err := json.MarshalIndent(items, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal loadbalancer list: %w", err)
		}
		p.Outputln(string(details))

		return nil
	case print.YAMLOutputFormat:
		details, err := yaml.MarshalWithOptions(items, yaml.IndentSequence(true), yaml.UseJSONMarshaler())
		if err != nil {
			return fmt.Errorf("marshal loadbalancer list: %w", err)
		}
		p.Outputln(string(details))

		return nil
	default:
		table := tables.NewTable()
		table.SetHeader("NAME", "EXTERNAL ADDRESS", "REGION", "STATUS", "VERSION")
		for _, item := range items {
			table.AddRow(utils.PtrString(item.Name),
				utils.PtrString(item.ExternalAddress),
				utils.PtrString(item.Region),
				utils.PtrString(item.Status),
				utils.PtrString(item.Version),
				utils.PtrString(item.ExternalAddress),
			)
		}
		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}

		return nil
	}
}
