package plans

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
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/rabbitmq/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/rabbitmq"
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
		Use:   "plans",
		Short: "Lists all RabbitMQ service plans",
		Long:  "Lists all RabbitMQ service plans.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`List all RabbitMQ service plans`,
				"$ stackit rabbitmq plans"),
			examples.NewExample(
				`List all RabbitMQ service plans in JSON format`,
				"$ stackit rabbitmq plans --output-format json"),
			examples.NewExample(
				`List up to 10 RabbitMQ service plans`,
				"$ stackit rabbitmq plans --limit 10"),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			model, err := parseInput(cmd)
			if err != nil {
				return err
			}

			// Configure API client
			apiClient, err := client.ConfigureClient(cmd)
			if err != nil {
				return err
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("get RabbitMQ service plans: %w", err)
			}
			plans := *resp.Offerings
			if len(plans) == 0 {
				projectLabel, err := projectname.GetProjectName(ctx, cmd)
				if err != nil {
					projectLabel = model.ProjectId
				}
				p.Info("No plans found for project %q\n", projectLabel)
				return nil
			}

			// Truncate output
			if model.Limit != nil && len(plans) > int(*model.Limit) {
				plans = plans[:*model.Limit]
			}

			return outputResult(p, model.OutputFormat, plans)
		},
	}

	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Int64(limitFlag, 0, "Maximum number of entries to list")
}

func parseInput(cmd *cobra.Command) (*inputModel, error) {
	globalFlags := globalflags.Parse(cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	limit := flags.FlagToInt64Pointer(cmd, limitFlag)
	if limit != nil && *limit < 1 {
		return nil, &errors.FlagValidationError{
			Flag:    limitFlag,
			Details: "must be greater than 0",
		}
	}

	return &inputModel{
		GlobalFlagModel: globalFlags,
		Limit:           limit,
	}, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *rabbitmq.APIClient) rabbitmq.ApiListOfferingsRequest {
	req := apiClient.ListOfferings(ctx, model.ProjectId)
	return req
}

func outputResult(p *print.Printer, outputFormat string, plans []rabbitmq.Offering) error {
	switch outputFormat {
	case globalflags.JSONOutputFormat:
		details, err := json.MarshalIndent(plans, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal RabbitMQ plans: %w", err)
		}
		p.Outputln(string(details))

		return nil
	default:
		table := tables.NewTable()
		table.SetHeader("OFFERING NAME", "VERSION", "ID", "NAME", "DESCRIPTION")
		for i := range plans {
			o := plans[i]
			for j := range *o.Plans {
				plan := (*o.Plans)[j]
				table.AddRow(*o.Name, *o.Version, *plan.Id, *plan.Name, *plan.Description)
			}
			table.AddSeparator()
		}
		table.EnableAutoMergeOnColumns(1, 2)
		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}

		return nil
	}
}
