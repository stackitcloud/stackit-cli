package update

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

type flagModel struct {
	ProjectId     string
	ZoneId        string
	Name          *string
	DefaultTTL    *int64
	Primaries     *[]string
	Acl           *string
	RetryTime     *int64
	RefreshTime   *int64
	NegativeCache *int64
	ExpireTime    *int64
	Description   *string
	ContactEmail  *string
}

const (
	projectIdFlag     = "project-id"
	zoneIdFlag        = "zone-id"
	nameFlag          = "name"
	defaultTTLFlag    = "default-ttl"
	primaryFlag       = "primary"
	aclFlag           = "acl"
	retryTimeFlag     = "retry-time"
	refreshTimeFlag   = "refresh-time"
	negativeCacheFlag = "negative-cache"
	expireTimeFlag    = "expire-time"
	descriptionFlag   = "description"
	contactEmailFlag  = "contact-email"
)

var Cmd = &cobra.Command{
	Use:     "update",
	Short:   "Updates a DNS zone",
	Long:    "Updates a DNS zone. Performs a partial update; fields not provided are kept unchanged",
	Example: `$ stackit dns zone update --project-id xxx --zone-id xxx --name my-zone --dns-name my-zone.com`,
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
			return fmt.Errorf("update DNS zone: %w", err)
		}

		// Wait for async operation
		_, err = wait.UpdateZoneWaitHandler(ctx, apiClient, model.ProjectId, model.ZoneId).WaitWithContext(ctx)
		if err != nil {
			return fmt.Errorf("wait for DNS zone update: %w", err)
		}

		cmd.Println("Zone updated")
		return nil
	},
}

func init() {
	configureFlags(Cmd)
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), zoneIdFlag, "Zone ID")
	cmd.Flags().String(nameFlag, "", "User given name of the zone")
	cmd.Flags().Int64(defaultTTLFlag, 1000, "Default time to live")
	cmd.Flags().StringSlice(primaryFlag, []string{}, "Primary name server for secondary zone")
	cmd.Flags().String(aclFlag, "", "Access control list")
	cmd.Flags().Int64(retryTimeFlag, 0, "Retry time")
	cmd.Flags().Int64(refreshTimeFlag, 0, "Refresh time")
	cmd.Flags().Int64(negativeCacheFlag, 0, "Negative cache")
	cmd.Flags().Int64(expireTimeFlag, 0, "Expire time")
	cmd.Flags().String(descriptionFlag, "", "Description of the zone")
	cmd.Flags().String(contactEmailFlag, "", "Contact email for the zone")

	err := utils.MarkFlagsRequired(cmd, zoneIdFlag)
	cobra.CheckErr(err)
}

func parseFlags(cmd *cobra.Command) (*flagModel, error) {
	projectId := viper.GetString(config.ProjectIdKey)
	if projectId == "" {
		return nil, fmt.Errorf("project ID not set")
	}

	return &flagModel{
		ProjectId:     projectId,
		ZoneId:        utils.FlagToStringValue(cmd, zoneIdFlag),
		Name:          utils.FlagToStringPointer(cmd, nameFlag),
		DefaultTTL:    utils.FlagToInt64Pointer(cmd, defaultTTLFlag),
		Primaries:     utils.FlagToStringSlicePointer(cmd, primaryFlag),
		Acl:           utils.FlagToStringPointer(cmd, aclFlag),
		RetryTime:     utils.FlagToInt64Pointer(cmd, retryTimeFlag),
		RefreshTime:   utils.FlagToInt64Pointer(cmd, refreshTimeFlag),
		NegativeCache: utils.FlagToInt64Pointer(cmd, negativeCacheFlag),
		ExpireTime:    utils.FlagToInt64Pointer(cmd, expireTimeFlag),
		Description:   utils.FlagToStringPointer(cmd, descriptionFlag),
		ContactEmail:  utils.FlagToStringPointer(cmd, contactEmailFlag),
	}, nil
}

func buildRequest(ctx context.Context, model *flagModel, apiClient *dns.APIClient) dns.ApiUpdateZoneRequest {
	req := apiClient.UpdateZone(ctx, model.ProjectId, model.ZoneId)
	req = req.UpdateZonePayload(dns.UpdateZonePayload{
		Name:          model.Name,
		DefaultTTL:    model.DefaultTTL,
		Primaries:     model.Primaries,
		Acl:           model.Acl,
		RetryTime:     model.RetryTime,
		RefreshTime:   model.RefreshTime,
		NegativeCache: model.NegativeCache,
		ExpireTime:    model.ExpireTime,
		Description:   model.Description,
		ContactEmail:  model.ContactEmail,
	})
	return req
}
