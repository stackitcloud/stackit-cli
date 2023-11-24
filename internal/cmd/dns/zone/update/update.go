package update

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/confirm"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/dns/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/dns"
	"github.com/stackitcloud/stackit-sdk-go/services/dns/wait"
)

const (
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

type flagModel struct {
	*globalflags.GlobalFlagModel
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

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
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

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to update zone %s?", model.ZoneId)
				err = confirm.PromptForConfirmation(cmd, prompt)
				if err != nil {
					return err
				}
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
	configureFlags(cmd)
	return cmd
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
	globalFlags := globalflags.Parse(cmd)
	if globalFlags.ProjectId == "" {
		return nil, fmt.Errorf("project ID not set")
	}

	zoneId := utils.FlagToStringValue(cmd, zoneIdFlag)
	name := utils.FlagToStringPointer(cmd, nameFlag)
	defaultTTL := utils.FlagToInt64Pointer(cmd, defaultTTLFlag)
	primaries := utils.FlagToStringSlicePointer(cmd, primaryFlag)
	acl := utils.FlagToStringPointer(cmd, aclFlag)
	retryTime := utils.FlagToInt64Pointer(cmd, retryTimeFlag)
	refreshTime := utils.FlagToInt64Pointer(cmd, refreshTimeFlag)
	negativeCache := utils.FlagToInt64Pointer(cmd, negativeCacheFlag)
	expireTime := utils.FlagToInt64Pointer(cmd, expireTimeFlag)
	description := utils.FlagToStringPointer(cmd, descriptionFlag)
	contactEmail := utils.FlagToStringPointer(cmd, contactEmailFlag)

	if name == nil && defaultTTL == nil && primaries == nil &&
		acl == nil && retryTime == nil && refreshTime == nil &&
		negativeCache == nil && expireTime == nil && description == nil &&
		contactEmail == nil {
		return nil, fmt.Errorf("please specify at least one field to update")
	}

	return &flagModel{
		GlobalFlagModel: globalFlags,
		ZoneId:          zoneId,
		Name:            name,
		DefaultTTL:      defaultTTL,
		Primaries:       primaries,
		Acl:             acl,
		RetryTime:       retryTime,
		RefreshTime:     refreshTime,
		NegativeCache:   negativeCache,
		ExpireTime:      expireTime,
		Description:     description,
		ContactEmail:    contactEmail,
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
