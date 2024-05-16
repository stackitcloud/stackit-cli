package list

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/goccy/go-yaml"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/argus/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"

	argusUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/argus/utils"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/argus"
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

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Lists all scrape configurations of an Argus instance",
		Long:  "Lists all scrape configurations of an Argus instance.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`List all scrape configurations of Argus instance "xxx"`,
				"$ stackit argus scrape-config list --instance-id xxx"),
			examples.NewExample(
				`List all scrape configurations of Argus instance "xxx" in JSON format`,
				"$ stackit argus scrape-config list --instance-id xxx --output-format json"),
			examples.NewExample(
				`List up to 10 scrape configurations of Argus instance "xxx"`,
				"$ stackit argus scrape-config list --instance-id xxx --limit 10"),
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
				return fmt.Errorf("get scrape configurations: %w", err)
			}
			configs := *resp.Data
			if len(configs) == 0 {
				instanceLabel, err := argusUtils.GetInstanceName(ctx, apiClient, model.InstanceId, model.ProjectId)
				if err != nil {
					p.Debug(print.ErrorLevel, "get instance name: %v", err)
					instanceLabel = model.InstanceId
				}
				p.Info("No scrape configurations found for instance %q\n", instanceLabel)
				return nil
			}

			// Truncate output
			if model.Limit != nil && len(configs) > int(*model.Limit) {
				configs = configs[:*model.Limit]
			}

			return outputResult(p, model.OutputFormat, configs)
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

func buildRequest(ctx context.Context, model *inputModel, apiClient *argus.APIClient) argus.ApiListScrapeConfigsRequest {
	req := apiClient.ListScrapeConfigs(ctx, model.InstanceId, model.ProjectId)
	return req
}

func outputResult(p *print.Printer, outputFormat string, configs []argus.Job) error {
	switch outputFormat {
	case print.JSONOutputFormat:
		details, err := json.MarshalIndent(configs, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal scrape configurations list: %w", err)
		}
		p.Outputln(string(details))

		return nil
	case print.YAMLOutputFormat:
		details, err := yaml.MarshalWithOptions(configs, yaml.IndentSequence(true))
		if err != nil {
			return fmt.Errorf("marshal scrape configurations list: %w", err)
		}
		p.Outputln(string(details))

		return nil
	default:
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

			table.AddRow(*c.JobName, targets, *c.ScrapeInterval)
		}
		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}

		return nil
	}
}
