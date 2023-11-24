package describe

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/dns/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/dns"
)

const (
	zoneIdFlag = "zone-id"
)

type flagModel struct {
	*globalflags.GlobalFlagModel
	ZoneId string
}

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "describe",
		Short:   "Get details of a DNS zone",
		Long:    "Get details of a DNS zone",
		Example: `$ stackit dns zone describe --project-id xxx --zone-id xxx`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			model, err := parseFlags(cmd)
			if err != nil {
				return err
			}
			// Configure API client
			apiClient, err := client.ConfigureClient(cmd)
			if err != nil {
				return fmt.Errorf("authentication failed, please run \"stackit auth login\" or \"stackit auth activate-service-account\"")
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("read DNS zone: %w", err)
			}
			zone := resp.Zone

			return outputResult(cmd, model.OutputFormat, zone)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), zoneIdFlag, "Zone ID")

	err := utils.MarkFlagsRequired(cmd, zoneIdFlag)
	cobra.CheckErr(err)
}

func parseFlags(cmd *cobra.Command) (*flagModel, error) {
	globalFlags := globalflags.Parse()
	if globalFlags.ProjectId == "" {
		return nil, fmt.Errorf("project ID not set")
	}

	return &flagModel{
		GlobalFlagModel: globalFlags,
		ZoneId:          utils.FlagToStringValue(cmd, zoneIdFlag),
	}, nil
}

func buildRequest(ctx context.Context, model *flagModel, apiClient *dns.APIClient) dns.ApiGetZoneRequest {
	req := apiClient.GetZone(ctx, model.ProjectId, model.ZoneId)
	return req
}

func outputResult(cmd *cobra.Command, outputFormat string, zone *dns.Zone) error {
	switch outputFormat {
	case globalflags.TableOutputFormat:
		table := tables.NewTable()
		table.SetHeader("ID", "NAME", "DNS_NAME", "STATE")
		table.AddRow(*zone.Id, *zone.Name, *zone.DnsName, *zone.State)
		err := table.Render(cmd)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}

		return nil
	default:
		details, err := json.MarshalIndent(zone, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal DNS zone: %w", err)
		}
		cmd.Println(string(details))

		return nil
	}
}
