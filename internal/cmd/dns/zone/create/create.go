package create

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/config"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/dns/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stackitcloud/stackit-sdk-go/services/dns"
	"github.com/stackitcloud/stackit-sdk-go/services/dns/wait"
)

const (
	projectIdFlag     = "project-id"
	nameFlag          = "name"
	dnsNameFlag       = "dns-name"
	defaultTTLFlag    = "default-ttl"
	primaryFlag       = "primary"
	aclFlag           = "acl"
	typeFlag          = "type"
	retryTimeFlag     = "retry-time"
	refreshTimeFlag   = "refresh-time"
	negativeCacheFlag = "negative-cache"
	isReverseZoneFlag = "is-reverse-zone"
	expireTimeFlag    = "expire-time"
	descriptionFlag   = "description"
	contactEmailFlag  = "contact-email"
)

type flagModel struct {
	ProjectId     string
	Name          *string
	DnsName       *string
	DefaultTTL    *int64
	Primaries     *[]string
	Acl           *string
	Type          *string
	RetryTime     *int64
	RefreshTime   *int64
	NegativeCache *int64
	IsReverseZone *bool
	ExpireTime    *int64
	Description   *string
	ContactEmail  *string
}

var Cmd = &cobra.Command{
	Use:     "create",
	Short:   "Creates a DNS zone",
	Long:    "Creates a DNS zone",
	Example: `$ stackit dns zone create --project-id xxx --name my-zone --dns-name my-zone.com`,
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
			return fmt.Errorf("create DNS zone: %w", err)
		}

		// Wait for async operation
		zoneId := *resp.Zone.Id
		_, err = wait.CreateZoneWaitHandler(ctx, apiClient, model.ProjectId, zoneId).WaitWithContext(ctx)
		if err != nil {
			return fmt.Errorf("wait for DNS zone creation: %w", err)
		}

		cmd.Printf("Created zone with ID %s\n", zoneId)
		return nil
	},
}

func init() {
	configureFlags(Cmd)
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().String(nameFlag, "", "User given name of the zone")
	cmd.Flags().String(dnsNameFlag, "", "DNS zone name")
	cmd.Flags().Int64(defaultTTLFlag, 1000, "Default time to live")
	cmd.Flags().StringSlice(primaryFlag, []string{}, "Primary name server for secondary zone")
	cmd.Flags().String(aclFlag, "", "Access control list")
	cmd.Flags().String(typeFlag, "", "Zone type")
	cmd.Flags().Int64(retryTimeFlag, 0, "Retry time")
	cmd.Flags().Int64(refreshTimeFlag, 0, "Refresh time")
	cmd.Flags().Int64(negativeCacheFlag, 0, "Negative cache")
	cmd.Flags().Bool(isReverseZoneFlag, false, "Is reverse zone")
	cmd.Flags().Int64(expireTimeFlag, 0, "Expire time")
	cmd.Flags().String(descriptionFlag, "", "Description of the zone")
	cmd.Flags().String(contactEmailFlag, "", "Contact email for the zone")

	err := utils.MarkFlagsRequired(cmd, nameFlag, dnsNameFlag)
	cobra.CheckErr(err)
}

func parseFlags(cmd *cobra.Command) (*flagModel, error) {
	projectId := viper.GetString(config.ProjectIdKey)
	if projectId == "" {
		return nil, fmt.Errorf("project ID not set")
	}

	return &flagModel{
		ProjectId:     projectId,
		Name:          utils.FlagToStringPointer(cmd, nameFlag),
		DnsName:       utils.FlagToStringPointer(cmd, dnsNameFlag),
		DefaultTTL:    utils.FlagToInt64Pointer(cmd, defaultTTLFlag),
		Primaries:     utils.FlagToStringSlicePointer(cmd, primaryFlag),
		Acl:           utils.FlagToStringPointer(cmd, aclFlag),
		Type:          utils.FlagToStringPointer(cmd, typeFlag),
		RetryTime:     utils.FlagToInt64Pointer(cmd, retryTimeFlag),
		RefreshTime:   utils.FlagToInt64Pointer(cmd, refreshTimeFlag),
		NegativeCache: utils.FlagToInt64Pointer(cmd, negativeCacheFlag),
		IsReverseZone: utils.FlagToBoolPointer(cmd, isReverseZoneFlag),
		ExpireTime:    utils.FlagToInt64Pointer(cmd, expireTimeFlag),
		Description:   utils.FlagToStringPointer(cmd, descriptionFlag),
		ContactEmail:  utils.FlagToStringPointer(cmd, contactEmailFlag),
	}, nil
}

func buildRequest(ctx context.Context, model *flagModel, apiClient *dns.APIClient) dns.ApiCreateZoneRequest {
	req := apiClient.CreateZone(ctx, model.ProjectId)
	req = req.CreateZonePayload(dns.CreateZonePayload{
		Name:          model.Name,
		DnsName:       model.DnsName,
		DefaultTTL:    model.DefaultTTL,
		Primaries:     model.Primaries,
		Acl:           model.Acl,
		Type:          model.Type,
		RetryTime:     model.RetryTime,
		RefreshTime:   model.RefreshTime,
		NegativeCache: model.NegativeCache,
		IsReverseZone: model.IsReverseZone,
		ExpireTime:    model.ExpireTime,
		Description:   model.Description,
		ContactEmail:  model.ContactEmail,
	})
	return req
}
