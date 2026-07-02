package describe

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/types"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/dns/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
	dns "github.com/stackitcloud/stackit-sdk-go/services/dns/v1api"
)

const (
	zoneIdArg = "ZONE_ID"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	ZoneId string
}

func NewCmd(params *types.CmdParams) *cobra.Command {
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
			model, err := parseInput(params.Printer, cmd, args)
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
				return fmt.Errorf("read DNS zone: %w", err)
			}
			zone := resp.Zone

			return outputResult(params.Printer, model.OutputFormat, &zone)
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

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *dns.APIClient) dns.ApiGetZoneRequest {
	req := apiClient.DefaultAPI.GetZone(ctx, model.ProjectId, model.ZoneId)
	return req
}

func outputResult(p *print.Printer, outputFormat string, zone *dns.Zone) error {
	if zone == nil {
		return fmt.Errorf("zone response is empty")
	}

	return p.OutputResult(outputFormat, zone, func() error {
		table := tables.NewTable()
		table.AddRow("ID", zone.Id)
		table.AddSeparator()
		table.AddRow("NAME", zone.Name)
		table.AddSeparator()
		table.AddRow("DESCRIPTION", utils.PtrString(zone.Description))
		table.AddSeparator()
		table.AddRow("STATE", zone.State)
		table.AddSeparator()
		table.AddRow("TYPE", zone.Type)
		table.AddSeparator()
		table.AddRow("DNS NAME", zone.DnsName)
		table.AddSeparator()
		table.AddRow("REVERSE ZONE", utils.PtrString(zone.IsReverseZone))
		table.AddSeparator()
		table.AddRow("RECORD COUNT", utils.PtrString(zone.RecordCount))
		table.AddSeparator()
		table.AddRow("CONTACT EMAIL", utils.PtrString(zone.ContactEmail))
		table.AddSeparator()
		table.AddRow("DEFAULT TTL", zone.DefaultTTL)
		table.AddSeparator()
		table.AddRow("SERIAL NUMBER", zone.SerialNumber)
		table.AddSeparator()
		table.AddRow("REFRESH TIME", zone.RefreshTime)
		table.AddSeparator()
		table.AddRow("RETRY TIME", zone.RetryTime)
		table.AddSeparator()
		table.AddRow("EXPIRE TIME", zone.ExpireTime)
		table.AddSeparator()
		table.AddRow("NEGATIVE CACHE", zone.NegativeCache)
		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}

		return nil
	})
}
