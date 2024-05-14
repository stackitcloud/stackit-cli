package describe

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ghodss/yaml"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/dns/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/dns"
)

const (
	zoneIdArg = "ZONE_ID"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	ZoneId string
}

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("describe %s", zoneIdArg),
		Short: "Shows details of a DNS zone",
		Long:  "Shows details of a DNS zone.",
		Args:  args.SingleArg(zoneIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Get details of a DNS zone with ID "xxx"`,
				"$ stackit dns zone describe xxx"),
			examples.NewExample(
				`Get details of a DNS zone with ID "xxx" in JSON format`,
				"$ stackit dns zone describe xxx --output-format json"),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			model, err := parseInput(p, cmd, args)
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
				return fmt.Errorf("read DNS zone: %w", err)
			}
			zone := resp.Zone

			return outputResult(p, model.OutputFormat, zone)
		},
	}
	return cmd
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	zoneId := inputArgs[0]

	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		ZoneId:          zoneId,
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

func buildRequest(ctx context.Context, model *inputModel, apiClient *dns.APIClient) dns.ApiGetZoneRequest {
	req := apiClient.GetZone(ctx, model.ProjectId, model.ZoneId)
	return req
}

func outputResult(p *print.Printer, outputFormat string, zone *dns.Zone) error {
	switch outputFormat {
	case print.JSONOutputFormat:
		details, err := json.MarshalIndent(zone, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal DNS zone: %w", err)
		}
		p.Outputln(string(details))

		return nil
	case print.YAMLOutputFormat:
		details, err := yaml.Marshal(zone)
		if err != nil {
			return fmt.Errorf("marshal DNS zone: %w", err)
		}
		p.Outputln(string(details))

		return nil
	default:
		table := tables.NewTable()
		table.AddRow("ID", *zone.Id)
		table.AddSeparator()
		table.AddRow("NAME", *zone.Name)
		table.AddSeparator()
		table.AddRow("DESCRIPTION", *zone.Description)
		table.AddSeparator()
		table.AddRow("STATE", *zone.State)
		table.AddSeparator()
		table.AddRow("TYPE", *zone.Type)
		table.AddSeparator()
		table.AddRow("DNS NAME", *zone.DnsName)
		table.AddSeparator()
		table.AddRow("REVERSE ZONE", *zone.IsReverseZone)
		table.AddSeparator()
		table.AddRow("RECORD COUNT", *zone.RecordCount)
		table.AddSeparator()
		table.AddRow("CONTACT EMAIL", *zone.ContactEmail)
		table.AddSeparator()
		table.AddRow("DEFAULT TTL", *zone.DefaultTTL)
		table.AddSeparator()
		table.AddRow("SERIAL NUMBER", *zone.SerialNumber)
		table.AddSeparator()
		table.AddRow("REFRESH TIME", *zone.RefreshTime)
		table.AddSeparator()
		table.AddRow("RETRY TIME", *zone.RetryTime)
		table.AddSeparator()
		table.AddRow("EXPIRE TIME", *zone.ExpireTime)
		table.AddSeparator()
		table.AddRow("NEGATIVE CACHE", *zone.NegativeCache)
		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}

		return nil
	}
}
