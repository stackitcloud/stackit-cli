package describe

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/cdn/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	sdkUtils "github.com/stackitcloud/stackit-sdk-go/core/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/cdn"
)

const distributionIDArg = "DISTRIBUTION_ID_ARG"
const flagWithWaf = "with-waf"

type inputModel struct {
	*globalflags.GlobalFlagModel
	DistributionID string
	WithWAF        bool
}

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "describe",
		Short: "Describe a CDN distribution",
		Long:  "Describe a CDN distribution by its ID.",
		Args:  args.SingleArg(distributionIDArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Get details of a CDN distribution with ID "xxx"`,
				`$ stackit beta cdn distribution describe xxx`,
			),
			examples.NewExample(
				`Get details of a CDN, including WAF details, for ID "xxx"`,
				`$ stackit beta cdn distribution describe xxx --with-waf`,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()

			model, err := parseInput(params.Printer, cmd, args)
			if err != nil {
				return err
			}

			apiClient, err := client.ConfigureClient(params.Printer, params.CliVersion)
			if err != nil {
				return err
			}

			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("read distribution: %w", err)
			}

			return outputResult(params.Printer, model.OutputFormat, resp)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Bool(flagWithWaf, false, "Include WAF details in the distribution description")
}

func parseInput(p *print.Printer, cmd *cobra.Command, args []string) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	model := &inputModel{
		GlobalFlagModel: globalFlags,
		DistributionID:  args[0],
		WithWAF:         flags.FlagToBoolValue(p, cmd, flagWithWaf),
	}
	p.DebugInputModel(model)
	return model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *cdn.APIClient) cdn.ApiGetDistributionRequest {
	return apiClient.GetDistribution(ctx, model.ProjectId, model.DistributionID).WithWafStatus(model.WithWAF)
}

func outputResult(p *print.Printer, outputFormat string, distribution *cdn.GetDistributionResponse) error {
	if distribution == nil {
		return fmt.Errorf("distribution response is empty")
	}
	return p.OutputResult(outputFormat, distribution, func() error {
		d := distribution.Distribution
		var content []tables.Table

		content = append(content, buildDistributionTable(d))

		if d.Waf != nil {
			content = append(content, buildWAFTable(d))
		}

		err := tables.DisplayTables(p, content)
		if err != nil {
			return fmt.Errorf("display table: %w", err)
		}
		return nil
	})
}

func buildDistributionTable(d *cdn.Distribution) tables.Table {
	regions := strings.Join(sdkUtils.EnumSliceToStringSlice(*d.Config.Regions), ", ")
	defaultCacheDuration := ""
	if d.Config.DefaultCacheDuration != nil && d.Config.DefaultCacheDuration.IsSet() {
		defaultCacheDuration = *d.Config.DefaultCacheDuration.Get()
	}
	logSinkPushUrl := ""
	if d.Config.LogSink != nil && d.Config.LogSink.LokiLogSink != nil {
		logSinkPushUrl = *d.Config.LogSink.LokiLogSink.PushUrl
	}
	monthlyLimitBytes := ""
	if d.Config.MonthlyLimitBytes != nil {
		monthlyLimitBytes = fmt.Sprintf("%d", *d.Config.MonthlyLimitBytes)
	}
	optimizerEnabled := ""
	if d.Config.Optimizer != nil {
		optimizerEnabled = fmt.Sprintf("%t", *d.Config.Optimizer.Enabled)
	}
	table := tables.NewTable()
	table.SetTitle("Distribution")
	table.AddRow("ID", utils.PtrString(d.Id))
	table.AddSeparator()
	table.AddRow("STATUS", utils.PtrString(d.Status))
	table.AddSeparator()
	table.AddRow("REGIONS", regions)
	table.AddSeparator()
	table.AddRow("CREATED AT", utils.PtrString(d.CreatedAt))
	table.AddSeparator()
	table.AddRow("UPDATED AT", utils.PtrString(d.UpdatedAt))
	table.AddSeparator()
	table.AddRow("PROJECT ID", utils.PtrString(d.ProjectId))
	table.AddSeparator()
	if d.Errors != nil && len(*d.Errors) > 0 {
		var errorDescriptions []string
		for _, err := range *d.Errors {
			errorDescriptions = append(errorDescriptions, *err.En)
		}
		table.AddRow("ERRORS", strings.Join(errorDescriptions, "\n"))
		table.AddSeparator()
	}
	if d.Config.Backend.BucketBackend != nil {
		b := d.Config.Backend.BucketBackend
		table.AddRow("BACKEND TYPE", "BUCKET")
		table.AddSeparator()
		table.AddRow("BUCKET URL", utils.PtrString(b.BucketUrl))
		table.AddSeparator()
		table.AddRow("BUCKET REGION", utils.PtrString(b.Region))
		table.AddSeparator()
	} else if d.Config.Backend.HttpBackend != nil {
		h := d.Config.Backend.HttpBackend
		var geofencing []string
		if h.Geofencing != nil {
			for k, v := range *h.Geofencing {
				geofencing = append(geofencing, fmt.Sprintf("%s: %s", k, strings.Join(v, ", ")))
			}
		}
		table.AddRow("BACKEND TYPE", "HTTP")
		table.AddSeparator()
		table.AddRow("HTTP ORIGIN URL", utils.PtrString(h.OriginUrl))
		table.AddSeparator()
		if h.OriginRequestHeaders != nil {
			table.AddRow("HTTP ORIGIN REQUEST HEADERS", utils.JoinStringMap(*h.OriginRequestHeaders, ": ", ", "))
			table.AddSeparator()
		}
		table.AddRow("HTTP GEOFENCING PROPERTIES", strings.Join(geofencing, "\n"))
		table.AddSeparator()
	}
	table.AddRow("BLOCKED COUNTRIES", strings.Join(*d.Config.BlockedCountries, ", "))
	table.AddSeparator()
	table.AddRow("BLOCKED IPS", strings.Join(*d.Config.BlockedIps, ", "))
	table.AddSeparator()
	table.AddRow("DEFAULT CACHE DURATION", defaultCacheDuration)
	table.AddSeparator()
	table.AddRow("LOG SINK PUSH URL", logSinkPushUrl)
	table.AddSeparator()
	table.AddRow("MONTHLY LIMIT (BYTES)", monthlyLimitBytes)
	table.AddSeparator()
	table.AddRow("OPTIMIZER ENABLED", optimizerEnabled)
	table.AddSeparator()
	// TODO config has yet another WAF block, left it out because the docs say to use the WAF block at the top level to determine enabled rules. There's also mode and type fields here, both left out.
	return table
}

func buildWAFTable(d *cdn.Distribution) tables.Table {
	table := tables.NewTable()
	table.SetTitle("WAF")
	for _, disabled := range *d.Waf.DisabledRules {
		table.AddRow("DISABLED RULE ID", utils.PtrString(disabled.Id))
		table.AddSeparator()
	}
	for _, enabled := range *d.Waf.EnabledRules {
		table.AddRow("ENABLED RULE ID", utils.PtrString(enabled.Id))
		table.AddSeparator()
	}
	for _, logOnly := range *d.Waf.LogOnlyRules {
		table.AddRow("LOG-ONLY RULE ID", utils.PtrString(logOnly.Id))
		table.AddSeparator()
	}
	return table
}
