package create

import (
	"context"
	"fmt"

	"stackit/internal/pkg/args"
	"stackit/internal/pkg/confirm"
	"stackit/internal/pkg/errors"
	"stackit/internal/pkg/examples"
	"stackit/internal/pkg/flags"
	"stackit/internal/pkg/globalflags"
	"stackit/internal/pkg/projectname"
	"stackit/internal/pkg/services/dns/client"
	"stackit/internal/pkg/spinner"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/dns"
	"github.com/stackitcloud/stackit-sdk-go/services/dns/wait"
)

const (
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

type inputModel struct {
	*globalflags.GlobalFlagModel
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

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Creates a DNS zone",
		Long:  "Creates a DNS zone",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Create a DNS zone with name "my-zone" and DNS name "www.my-zone.com"`,
				"$ stackit dns zone create --name my-zone --dns-name www.my-zone.com"),
			examples.NewExample(
				`Create a DNS zone with name "my-zone", DNS name "www.my-zone.com" and default time to live of 1000ms`,
				"$ stackit dns zone create --name my-zone --dns-name www.my-zone.com --default-ttl 1000"),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			model, err := parseInput(cmd)
			if err != nil {
				return err
			}

			// Configure API client
			apiClient, err := client.ConfigureClient(cmd)
			if err != nil {
				return err
			}

			projectLabel, err := projectname.GetProjectName(ctx, cmd)
			if err != nil {
				projectLabel = model.ProjectId
			}

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to create a zone for project %s?", projectLabel)
				err = confirm.PromptForConfirmation(cmd, prompt)
				if err != nil {
					return err
				}
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("create DNS zone: %w", err)
			}
			zoneId := *resp.Zone.Id

			// Wait for async operation, if async mode not enabled
			if !model.Async {
				s := spinner.New(cmd)
				s.Start("Creating zone")
				_, err = wait.CreateZoneWaitHandler(ctx, apiClient, model.ProjectId, zoneId).WaitWithContext(ctx)
				if err != nil {
					return fmt.Errorf("wait for DNS zone creation: %w", err)
				}
				s.Stop()
			}

			operationState := "Created"
			if model.Async {
				operationState = "Triggered creation of"
			}
			cmd.Printf("%s zone for project %s. Zone ID: %s\n", operationState, projectLabel, zoneId)
			return nil
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().String(nameFlag, "", "User given name of the zone")
	cmd.Flags().String(dnsNameFlag, "", "Fully qualified domain name of the DNS zone")
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

	err := flags.MarkFlagsRequired(cmd, nameFlag, dnsNameFlag)
	cobra.CheckErr(err)
}

func parseInput(cmd *cobra.Command) (*inputModel, error) {
	globalFlags := globalflags.Parse(cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	return &inputModel{
		GlobalFlagModel: globalFlags,
		Name:            flags.FlagToStringPointer(cmd, nameFlag),
		DnsName:         flags.FlagToStringPointer(cmd, dnsNameFlag),
		DefaultTTL:      flags.FlagToInt64Pointer(cmd, defaultTTLFlag),
		Primaries:       flags.FlagToStringSlicePointer(cmd, primaryFlag),
		Acl:             flags.FlagToStringPointer(cmd, aclFlag),
		Type:            flags.FlagToStringPointer(cmd, typeFlag),
		RetryTime:       flags.FlagToInt64Pointer(cmd, retryTimeFlag),
		RefreshTime:     flags.FlagToInt64Pointer(cmd, refreshTimeFlag),
		NegativeCache:   flags.FlagToInt64Pointer(cmd, negativeCacheFlag),
		IsReverseZone:   flags.FlagToBoolPointer(cmd, isReverseZoneFlag),
		ExpireTime:      flags.FlagToInt64Pointer(cmd, expireTimeFlag),
		Description:     flags.FlagToStringPointer(cmd, descriptionFlag),
		ContactEmail:    flags.FlagToStringPointer(cmd, contactEmailFlag),
	}, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *dns.APIClient) dns.ApiCreateZoneRequest {
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
