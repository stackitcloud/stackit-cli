package update

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/confirm"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/dns/client"
	dnsUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/dns/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/spinner"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/dns"
	"github.com/stackitcloud/stackit-sdk-go/services/dns/wait"
)

const (
	zoneIdArg = "ZONE_ID"

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

type inputModel struct {
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
		Use:   fmt.Sprintf("update %s", zoneIdArg),
		Short: "Update a DNS zone",
		Long:  "Update a DNS zone. Performs a partial update, i.e., fields not provided are kept unchanged.",
		Args:  args.SingleArg(zoneIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Update the contact email of the DNS zone with ID "xxx"`,
				"$ stackit dns zone update xxx --contact-email someone@domain.com"),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			model, err := parseInput(cmd, args)
			if err != nil {
				return err
			}

			// Configure API client
			apiClient, err := client.ConfigureClient(cmd)
			if err != nil {
				return err
			}

			zoneLabel, err := dnsUtils.GetZoneName(ctx, apiClient, model.ProjectId, model.ZoneId)
			if err != nil {
				zoneLabel = model.ZoneId
			}

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to update zone %s?", zoneLabel)
				err = confirm.PromptForConfirmation(cmd, prompt)
				if err != nil {
					return err
				}
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

			// Wait for async operation, if async mode not enabled
			if !model.Async {
				s := spinner.New(cmd)
				s.Start("Updating zone")
				_, err = wait.PartialUpdateZoneWaitHandler(ctx, apiClient, model.ProjectId, model.ZoneId).WaitWithContext(ctx)
				if err != nil {
					return fmt.Errorf("wait for DNS zone update: %w", err)
				}
				s.Stop()
			}

			operationState := "Updated"
			if model.Async {
				operationState = "Triggered update of"
			}
			cmd.Printf("%s zone %s\n", operationState, zoneLabel)
			return nil
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
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
}

func parseInput(cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	zoneId := inputArgs[0]

	globalFlags := globalflags.Parse(cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	name := flags.FlagToStringPointer(cmd, nameFlag)
	defaultTTL := flags.FlagToInt64Pointer(cmd, defaultTTLFlag)
	primaries := flags.FlagToStringSlicePointer(cmd, primaryFlag)
	acl := flags.FlagToStringPointer(cmd, aclFlag)
	retryTime := flags.FlagToInt64Pointer(cmd, retryTimeFlag)
	refreshTime := flags.FlagToInt64Pointer(cmd, refreshTimeFlag)
	negativeCache := flags.FlagToInt64Pointer(cmd, negativeCacheFlag)
	expireTime := flags.FlagToInt64Pointer(cmd, expireTimeFlag)
	description := flags.FlagToStringPointer(cmd, descriptionFlag)
	contactEmail := flags.FlagToStringPointer(cmd, contactEmailFlag)

	if name == nil && defaultTTL == nil && primaries == nil &&
		acl == nil && retryTime == nil && refreshTime == nil &&
		negativeCache == nil && expireTime == nil && description == nil &&
		contactEmail == nil {
		return nil, &errors.EmptyUpdateError{}
	}

	return &inputModel{
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

func buildRequest(ctx context.Context, model *inputModel, apiClient *dns.APIClient) dns.ApiPartialUpdateZoneRequest {
	req := apiClient.PartialUpdateZone(ctx, model.ProjectId, model.ZoneId)
	req = req.PartialUpdateZonePayload(dns.PartialUpdateZonePayload{
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
