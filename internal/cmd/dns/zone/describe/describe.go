package describe

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/config"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/dns/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stackitcloud/stackit-sdk-go/services/dns"
)

const (
	projectIdFlag = "project-id"
	zoneIdFlag    = "zone-id"
)

type flagModel struct {
	ProjectId string
	ZoneId    string
}

var Cmd = &cobra.Command{
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
		zone := *resp.Zone

		// Show details
		details, err := json.MarshalIndent(zone, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal DNS zone: %w", err)
		}
		fmt.Println(string(details))

		return nil
	},
}

func init() {
	configureFlags(Cmd)
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), zoneIdFlag, "Zone ID")

	err := utils.MarkFlagsRequired(cmd, zoneIdFlag)
	cobra.CheckErr(err)
}

func parseFlags(cmd *cobra.Command) (*flagModel, error) {
	projectId := viper.GetString(config.ProjectIdKey)
	if projectId == "" {
		return nil, fmt.Errorf("project ID not set")
	}

	return &flagModel{
		ProjectId: projectId,
		ZoneId:    utils.FlagToStringValue(cmd, zoneIdFlag),
	}, nil
}

func buildRequest(ctx context.Context, model *flagModel, apiClient *dns.APIClient) dns.ApiGetZoneRequest {
	req := apiClient.GetZone(ctx, model.ProjectId, model.ZoneId)
	return req
}
