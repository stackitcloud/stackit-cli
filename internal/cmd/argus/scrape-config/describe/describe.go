package describe

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
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/argus/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-sdk-go/services/argus"

	"github.com/spf13/cobra"
)

const (
	jobNameArg = "JOB_NAME"

	instanceIdFlag = "instance-id"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	JobName    string
	InstanceId string
}

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("describe %s", jobNameArg),
		Short: "Shows details of a scrape configuration from an Argus instance",
		Long:  "Shows details of a scrape configuration from an Argus instance.",
		Args:  args.SingleArg(jobNameArg, nil),
		Example: examples.Build(
			examples.NewExample(
				`Get details of a scrape configuration with name "my-config" from Argus instance "xxx"`,
				"$ stackit argus scrape-config describe my-config --instance-id xxx"),
			examples.NewExample(
				`Get details of a scrape configuration with name "my-config" from Argus instance "xxx" in a table format`,
				"$ stackit argus scrape-config describe my-config --output-format pretty"),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			model, err := parseInput(cmd, args)
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
				return fmt.Errorf("read scrape configuration: %w", err)
			}

			return outputResult(p, model.OutputFormat, resp.Data)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), instanceIdFlag, "Instance ID")

	err := flags.MarkFlagsRequired(cmd, instanceIdFlag)
	cobra.CheckErr(err)
}

func parseInput(cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	jobName := inputArgs[0]

	globalFlags := globalflags.Parse(cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	return &inputModel{
		GlobalFlagModel: globalFlags,
		JobName:         jobName,
		InstanceId:      flags.FlagToStringValue(cmd, instanceIdFlag),
	}, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *argus.APIClient) argus.ApiGetScrapeConfigRequest {
	req := apiClient.GetScrapeConfig(ctx, model.InstanceId, model.JobName, model.ProjectId)
	return req
}

func outputResult(p *print.Printer, outputFormat string, config *argus.Job) error {
	switch outputFormat {
	case globalflags.PrettyOutputFormat:

		saml2Enabled := "Enabled"
		if config.Params != nil {
			saml2 := (*config.Params)["saml2"]
			if len(saml2) > 0 && saml2[0] == "disabled" {
				saml2Enabled = "Disabled"
			}
		}

		targets := []string{}
		for _, target := range *config.StaticConfigs {
			targetFmt := ""
			if target.Labels != nil {
				// make map prettier
				for k, v := range *target.Labels {
					if targetFmt != "" {
						targetFmt += " "
					} else {
						targetFmt += "labels: ["
					}
					targetFmt += fmt.Sprintf("%s:%s", k, v)
				}
				if targetFmt != "" {
					targetFmt += "]"
				}
			}
			if target.Targets != nil {
				if targetFmt != "" {
					targetFmt += "; "
				}
				targetFmt += fmt.Sprintf("urls: %v", *target.Targets)
			}
			targets = append(targets, targetFmt)
		}

		table := tables.NewTable()
		table.AddRow("NAME", *config.JobName)
		table.AddSeparator()
		table.AddRow("METRICS PATH", *config.MetricsPath)
		table.AddSeparator()
		table.AddRow("SCHEME", *config.Scheme)
		table.AddSeparator()
		table.AddRow("SCRAPE INTERVAL", *config.ScrapeInterval)
		table.AddSeparator()
		table.AddRow("SCRAPE TIMEOUT", *config.ScrapeTimeout)
		table.AddSeparator()
		table.AddRow("SAML2", saml2Enabled)
		table.AddSeparator()
		if config.BasicAuth == nil {
			table.AddRow("AUTHENTICATION", "None")
		} else {
			table.AddRow("AUTHENTICATION", "Basic Auth")
			table.AddSeparator()
			table.AddRow("USERNAME", *config.BasicAuth.Username)
			table.AddSeparator()
			table.AddRow("PASSWORD", *config.BasicAuth.Password)
		}
		table.AddSeparator()
		for i, target := range targets {
			table.AddRow(fmt.Sprintf("TARGET %d", i+1), target)
			table.AddSeparator()
		}

		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}

		return nil
	default:
		details, err := json.MarshalIndent(config, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal scrape configuration: %w", err)
		}
		p.Outputln(string(details))

		return nil
	}
}
