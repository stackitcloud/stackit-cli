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
		Use:   "list",
		Short: "Lists all RabbitMQ instances",
		Long:  "Lists all RabbitMQ instances.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`List all RabbitMQ instances`,
				"$ stackit rabbitmq instance list"),
			examples.NewExample(
				`List all RabbitMQ instances in JSON format`,
				"$ stackit rabbitmq instance list --output-format json"),
			examples.NewExample(
				`List up to 10 RabbitMQ instances`,
				"$ stackit rabbitmq instance list --limit 10"),
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

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("get RabbitMQ instances: %w", err)
			}
			instances := *resp.Instances
			if len(instances) == 0 {
				projectLabel, err := projectname.GetProjectName(ctx, p, cmd)
				if err != nil {
					p.Debug(print.ErrorLevel, "get project name: %v", err)
					projectLabel = model.ProjectId
				}
				p.Info("No instances found for project %q\n", projectLabel)
				return nil
			}

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

func buildRequest(ctx context.Context, model *inputModel, apiClient *rabbitmq.APIClient) rabbitmq.ApiListInstancesRequest {
	req := apiClient.ListInstances(ctx, model.ProjectId)
	return req
}

func outputResult(p *print.Printer, outputFormat string, instances []rabbitmq.Instance) error {
	switch outputFormat {
	case print.JSONOutputFormat:
		details, err := json.MarshalIndent(instances, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal RabbitMQ instance list: %w", err)
		}
		p.Outputln(string(details))

		return nil
	default:
		table := tables.NewTable()
		table.SetHeader("ID", "NAME", "LAST OPERATION TYPE", "LAST OPERATION STATE")
		for i := range instances {
			instance := instances[i]
			table.AddRow(*instance.InstanceId, *instance.Name, *instance.LastOperation.Type, *instance.LastOperation.State)
		}
		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}

		return nil
	}
}
