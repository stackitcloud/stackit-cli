package create

import (
	"context"
	"fmt"

	sdkUtils "github.com/stackitcloud/stackit-sdk-go/core/utils"

	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/projectname"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/dns/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/spinner"

	"github.com/spf13/cobra"
	dns "github.com/stackitcloud/stackit-sdk-go/services/dns/v1api"
	"github.com/stackitcloud/stackit-sdk-go/services/dns/v1api/wait"
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
	Name          string
	DnsName       string
	DefaultTTL    *int32
	Primaries     []string
	Acl           *string
	Type          *dns.CreateZonePayloadType
	RetryTime     *int32
	RefreshTime   *int32
	NegativeCache *int32
	IsReverseZone *bool
	ExpireTime    *int32
	Description   *string
	ContactEmail  *string
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Creates a DNS zone",
		Long:  "Creates a DNS zone.",
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
			model, err := parseInput(params.Printer, cmd, args)
			if err != nil {
				return err
			}

			// Configure API client
			apiClient, err := client.ConfigureClient(params.Printer, params.CliVersion)
			if err != nil {
				return err
			}

			projectLabel, err := projectname.GetProjectName(ctx, params.Printer, params.CliVersion, cmd)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get project name: %v", err)
				projectLabel = model.ProjectId
			}

			prompt := fmt.Sprintf("Are you sure you want to create a zone for project %q?", projectLabel)
			err = params.Printer.PromptForConfirmation(prompt)
			if err != nil {
				return err
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("create DNS zone: %w", err)
			}
			zoneId := resp.Zone.Id

			// Wait for async operation, if async mode not enabled
			if !model.Async {
				err := spinner.Run(params.Printer, "Creating zone", func() error {
					_, err = wait.CreateZoneWaitHandler(ctx, apiClient.DefaultAPI, model.ProjectId, zoneId).WaitWithContext(ctx)
					return err
				})
				if err != nil {
					return fmt.Errorf("wait for DNS zone creation: %w", err)
				}
			}

			return outputResult(params.Printer, model, projectLabel, resp)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().String(nameFlag, "", "User given name of the zone")
	cmd.Flags().String(dnsNameFlag, "", "Fully qualified domain name of the DNS zone")
	cmd.Flags().Int32(defaultTTLFlag, 1000, "Default time to live")
	cmd.Flags().StringSlice(primaryFlag, []string{}, "Primary name server for secondary zone")
	cmd.Flags().String(aclFlag, "", "Access control list")
	cmd.Flags().Var(flags.EnumFlag(false, "", append(sdkUtils.EnumSliceToStringSlice(dns.AllowedCreateZonePayloadTypeEnumValues), "")...), typeFlag, fmt.Sprintf("Zone type, one of: %q", utils.FormatPossibleValues(sdkUtils.EnumSliceToStringSlice(dns.AllowedCreateZonePayloadTypeEnumValues)...)))
	cmd.Flags().Int32(retryTimeFlag, 0, "Retry time")
	cmd.Flags().Int32(refreshTimeFlag, 0, "Refresh time")
	cmd.Flags().Int32(negativeCacheFlag, 0, "Negative cache")
	cmd.Flags().Bool(isReverseZoneFlag, false, "Is reverse zone")
	cmd.Flags().Int32(expireTimeFlag, 0, "Expire time")
	cmd.Flags().String(descriptionFlag, "", "Description of the zone")
	cmd.Flags().String(contactEmailFlag, "", "Contact email for the zone")

	err := flags.MarkFlagsRequired(cmd, nameFlag, dnsNameFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command, _ []string) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	var zoneType *dns.CreateZonePayloadType
	if zoneTypeString := flags.FlagToStringPointer(p, cmd, typeFlag); zoneTypeString != nil && *zoneTypeString != "" {
		zoneType = dns.CreateZonePayloadType(*zoneTypeString).Ptr()
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		Name:            flags.FlagToStringValue(p, cmd, nameFlag),
		DnsName:         flags.FlagToStringValue(p, cmd, dnsNameFlag),
		DefaultTTL:      flags.FlagToInt32Pointer(p, cmd, defaultTTLFlag),
		Primaries:       flags.FlagToStringSliceValue(p, cmd, primaryFlag),
		Acl:             flags.FlagToStringPointer(p, cmd, aclFlag),
		Type:            zoneType,
		RetryTime:       flags.FlagToInt32Pointer(p, cmd, retryTimeFlag),
		RefreshTime:     flags.FlagToInt32Pointer(p, cmd, refreshTimeFlag),
		NegativeCache:   flags.FlagToInt32Pointer(p, cmd, negativeCacheFlag),
		IsReverseZone:   flags.FlagToBoolPointer(p, cmd, isReverseZoneFlag),
		ExpireTime:      flags.FlagToInt32Pointer(p, cmd, expireTimeFlag),
		Description:     flags.FlagToStringPointer(p, cmd, descriptionFlag),
		ContactEmail:    flags.FlagToStringPointer(p, cmd, contactEmailFlag),
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *dns.APIClient) dns.ApiCreateZoneRequest {
	req := apiClient.DefaultAPI.CreateZone(ctx, model.ProjectId)
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

func outputResult(p *print.Printer, model *inputModel, projectLabel string, resp *dns.ZoneResponse) error {
	if resp == nil {
		return fmt.Errorf("dns zone response is empty")
	}
	return p.OutputResult(model.OutputFormat, resp, func() error {
		operationState := "Created"
		if model.Async {
			operationState = "Triggered creation of"
		}
		p.Outputf("%s zone for project %q. Zone ID: %s\n", operationState, projectLabel, resp.Zone.Id)
		return nil
	})
}
