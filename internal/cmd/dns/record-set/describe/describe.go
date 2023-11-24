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
	zoneIdFlag      = "zone-id"
	recordSetIdFlag = "record-set-id"
)

type flagModel struct {
	*globalflags.GlobalFlagModel
	ZoneId      string
	RecordSetId string
}

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "describe",
		Short:   "Get details of a DNS record set",
		Long:    "Get details of a DNS record set",
		Example: `$ stackit dns record-set describe --project-id xxx --zone-id xxx --record-set-id xxx`,
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
				return fmt.Errorf("read DNS record set: %w", err)
			}
			recordSet := resp.Rrset

			return outputResult(cmd, model.OutputFormat, recordSet)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), zoneIdFlag, "Zone ID")
	cmd.Flags().Var(flags.UUIDFlag(), recordSetIdFlag, "Record Set ID")

	err := utils.MarkFlagsRequired(cmd, zoneIdFlag, recordSetIdFlag)
	cobra.CheckErr(err)
}

func parseFlags(cmd *cobra.Command) (*flagModel, error) {
	globalFlags := globalflags.Parse(cmd)
	if globalFlags.ProjectId == "" {
		return nil, fmt.Errorf("project ID not set")
	}

	return &flagModel{
		GlobalFlagModel: globalFlags,
		ZoneId:          utils.FlagToStringValue(cmd, zoneIdFlag),
		RecordSetId:     utils.FlagToStringValue(cmd, recordSetIdFlag),
	}, nil
}

func buildRequest(ctx context.Context, model *flagModel, apiClient *dns.APIClient) dns.ApiGetRecordSetRequest {
	req := apiClient.GetRecordSet(ctx, model.ProjectId, model.ZoneId, model.RecordSetId)
	return req
}

func outputResult(cmd *cobra.Command, outputFormat string, recordSet *dns.RecordSet) error {
	switch outputFormat {
	case globalflags.TableOutputFormat:
		table := tables.NewTable()
		table.SetHeader("ID", "Name", "Type", "State")
		table.AddRow(*recordSet.Id, *recordSet.Name, *recordSet.Type, *recordSet.State)
		err := table.Render(cmd)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}

		return nil
	default:
		details, err := json.MarshalIndent(recordSet, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal DNS record set: %w", err)
		}
		cmd.Println(string(details))

		return nil
	}
}
