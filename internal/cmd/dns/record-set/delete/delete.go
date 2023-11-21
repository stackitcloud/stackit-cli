package delete

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/config"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/dns/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stackitcloud/stackit-sdk-go/services/dns"
	"github.com/stackitcloud/stackit-sdk-go/services/dns/wait"
)

const (
	projectIdFlag   = "project-id"
	zoneIdFlag      = "zone-id"
	recordSetIdFlag = "record-set-id"
)

type flagModel struct {
	ProjectId   string
	ZoneId      string
	RecordSetId string
}

var Cmd = &cobra.Command{
	Use:     "delete",
	Short:   "Delete a DNS record set",
	Long:    "Delete a DNS record set",
	Example: `$ stackit dns record-set delete --project-id xxx --zone-id xxx --record-set-id xxx`,
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
		if err != nil {
			return err
		}
		_, err = req.Execute()
		if err != nil {
			return fmt.Errorf("delete DNS record set: %w", err)
		}

		// Wait for async operation
		_, err = wait.DeleteRecordSetWaitHandler(ctx, apiClient, model.ProjectId, model.ZoneId, model.RecordSetId).WaitWithContext(ctx)
		if err != nil {
			return fmt.Errorf("wait for DNS record set deletion: %w", err)
		}

		cmd.Println("Record set deleted")
		return nil
	},
}

func init() {
	configureFlags(Cmd)
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), zoneIdFlag, "Zone ID")
	cmd.Flags().Var(flags.UUIDFlag(), recordSetIdFlag, "Record set ID")

	err := utils.MarkFlagsRequired(cmd, zoneIdFlag, recordSetIdFlag)
	cobra.CheckErr(err)
}

func parseFlags(cmd *cobra.Command) (*flagModel, error) {
	projectId := viper.GetString(config.ProjectIdKey)
	if projectId == "" {
		return nil, fmt.Errorf("project ID not set")
	}

	return &flagModel{
		ProjectId:   projectId,
		ZoneId:      utils.FlagToStringValue(cmd, zoneIdFlag),
		RecordSetId: utils.FlagToStringValue(cmd, recordSetIdFlag),
	}, nil
}

func buildRequest(ctx context.Context, model *flagModel, apiClient *dns.APIClient) dns.ApiDeleteRecordSetRequest {
	req := apiClient.DeleteRecordSet(ctx, model.ProjectId, model.ZoneId, model.RecordSetId)
	return req
}
