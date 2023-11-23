package describe

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/commonflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/dns/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/dns"
)

const (
	zoneIdFlag      = "zone-id"
	recordSetIdFlag = "record-set-id"
)

type flagModel struct {
	ProjectId   string
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
			recordSet := *resp.Rrset

			// Show details
			details, err := json.MarshalIndent(recordSet, "", "  ")
			if err != nil {
				return fmt.Errorf("marshal DNS record set: %w", err)
			}
			cmd.Println(string(details))

			return nil
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
	projectId := commonflags.GetString(commonflags.ProjectIdFlag)
	if projectId == "" {
		return nil, fmt.Errorf("project ID not set")
	}

	return &flagModel{
		ProjectId:   projectId,
		ZoneId:      utils.FlagToStringValue(cmd, zoneIdFlag),
		RecordSetId: utils.FlagToStringValue(cmd, recordSetIdFlag),
	}, nil
}

func buildRequest(ctx context.Context, model *flagModel, apiClient *dns.APIClient) dns.ApiGetRecordSetRequest {
	req := apiClient.GetRecordSet(ctx, model.ProjectId, model.ZoneId, model.RecordSetId)
	return req
}
