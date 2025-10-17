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
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/observability/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	observabilityUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/observability/utils"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/observability"
)

const (
	limitFlag      = "limit"
	instanceIdFlag = "instance-id"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	Limit      *int64
	InstanceId string
}

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Lists all scrape configurations of an Observability instance",
		Long:  "Lists all scrape configurations of an Observability instance.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`List all scrape configurations of Observability instance "xxx"`,
				"$ stackit observability scrape-config list --instance-id xxx"),
			examples.NewExample(
				`List all scrape configurations of Observability instance "xxx" in JSON format`,
				"$ stackit observability scrape-config list --instance-id xxx --output-format json"),
			examples.NewExample(
				`List up to 10 scrape configurations of Observability instance "xxx"`,
				"$ stackit observability scrape-config list --instance-id xxx --limit 10"),
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

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("get scrape configurations: %w", err)
			}
			configs := *resp.Data
			if len(configs) == 0 {
				instanceLabel, err := observabilityUtils.GetInstanceName(ctx, apiClient, model.InstanceId, model.ProjectId)
				if err != nil {
					params.Printer.Debug(print.ErrorLevel, "get instance name: %v", err)
					instanceLabel = model.InstanceId
				}
				params.Printer.Info("No scrape configurations found for instance %q\n", instanceLabel)
				return nil
			}

			// Truncate output
			if model.Limit != nil && len(configs) > int(*model.Limit) {
				configs = configs[:*model.Limit]
			}

			return outputResult(params.Printer, model.OutputFormat, configs)
		},
	}

	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Int64(limitFlag, 0, "Maximum number of entries to list")
	cmd.Flags().Var(flags.UUIDFlag(), instanceIdFlag, "Instance ID")

	err := flags.MarkFlagsRequired(cmd, instanceIdFlag)
	cobra.CheckErr(err)
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
		InstanceId:      flags.FlagToStringValue(p, cmd, instanceIdFlag),
	}, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *observability.APIClient) observability.ApiListScrapeConfigsRequest {
	req := apiClient.ListScrapeConfigs(ctx, model.InstanceId, model.ProjectId)
	return req
}

func outputResult(p *print.Printer, outputFormat string, configs []observability.Job) error {
	return p.OutputResult(outputFormat, configs, func() error {
		table := tables.NewTable()
		table.SetHeader("NAME", "TARGETS", "SCRAPE INTERVAL")
		for i := range configs {
			c := configs[i]

			targets := 0
			if c.StaticConfigs != nil {
				for _, sc := range *c.StaticConfigs {
					if sc.Targets == nil {
						continue
					}
					targets += len(*sc.Targets)
				}
			}

			table.AddRow(
				utils.PtrString(c.JobName),
				targets,
				utils.PtrString(c.ScrapeInterval),
			)
		}
		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}

		return nil
	})
}
