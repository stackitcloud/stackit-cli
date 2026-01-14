package clone

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/types"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/dns/client"
	dnsUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/dns/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/spinner"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/dns"
	"github.com/stackitcloud/stackit-sdk-go/services/dns/wait"
)

const (
	nameFlag          = "name"
	dnsNameFlag       = "dns-name"
	descriptionFlag   = "description"
	adjustRecordsFlag = "adjust-records"
	zoneIdArg         = "ZONE_ID"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	Name          *string
	DnsName       *string
	Description   *string
	AdjustRecords *bool
	ZoneId        string
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("clone %s", zoneIdArg),
		Short: "Clones a DNS zone",
		Long:  "Clones an existing DNS zone with all record sets to a new zone with a different name.",
		Args:  args.SingleArg(zoneIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Clones a DNS zone with ID "xxx" to a new zone with DNS name "www.my-zone.com"`,
				"$ stackit dns zone clone xxx --dns-name www.my-zone.com"),
			examples.NewExample(
				`Clones a DNS zone with ID "xxx" to a new zone with DNS name "www.my-zone.com" and display name "new-zone"`,
				"$ stackit dns zone clone xxx --dns-name www.my-zone.com --name new-zone"),
			examples.NewExample(
				`Clones a DNS zone with ID "xxx" to a new zone with DNS name "www.my-zone.com" and adjust records "true"`,
				"$ stackit dns zone clone xxx --dns-name www.my-zone.com --adjust-records"),
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

			zoneLabel, err := dnsUtils.GetZoneName(ctx, apiClient, model.ProjectId, model.ZoneId)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get zone name: %v", err)
				zoneLabel = model.ZoneId
			}

			prompt := fmt.Sprintf("Are you sure you want to clone the zone %q?", zoneLabel)
			err = params.Printer.PromptForConfirmation(prompt)
			if err != nil {
				return err
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("clone DNS zone: %w", err)
			}
			zoneId := *resp.Zone.Id

			// Wait for async operation, if async mode not enabled
			if !model.Async {
				s := spinner.New(params.Printer)
				s.Start("Cloning zone")
				_, err = wait.CreateZoneWaitHandler(ctx, apiClient, model.ProjectId, zoneId).WaitWithContext(ctx)
				if err != nil {
					return fmt.Errorf("wait for DNS zone cloning: %w", err)
				}
				s.Stop()
			}

			return outputResult(params.Printer, model, zoneLabel, resp)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().String(nameFlag, "", "User given new name for the cloned zone")
	cmd.Flags().String(dnsNameFlag, "", "Fully qualified domain name of the new DNS zone to clone")
	cmd.Flags().String(descriptionFlag, "", "New description for the cloned zone")
	cmd.Flags().Bool(adjustRecordsFlag, false, "Sets content and replaces the DNS name of the original zone with the new DNS name of the cloned zone")

	err := flags.MarkFlagsRequired(cmd, dnsNameFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	zoneId := inputArgs[0]

	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		Name:            flags.FlagToStringPointer(p, cmd, nameFlag),
		DnsName:         flags.FlagToStringPointer(p, cmd, dnsNameFlag),
		Description:     flags.FlagToStringPointer(p, cmd, descriptionFlag),
		AdjustRecords:   flags.FlagToBoolPointer(p, cmd, adjustRecordsFlag),
		ZoneId:          zoneId,
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *dns.APIClient) dns.ApiCloneZoneRequest {
	req := apiClient.CloneZone(ctx, model.ProjectId, model.ZoneId)
	req = req.CloneZonePayload(dns.CloneZonePayload{
		Name:          model.Name,
		DnsName:       model.DnsName,
		Description:   model.Description,
		AdjustRecords: model.AdjustRecords,
	})
	return req
}

func outputResult(p *print.Printer, model *inputModel, projectLabel string, resp *dns.ZoneResponse) error {
	if resp == nil {
		return fmt.Errorf("dns zone response is empty")
	}
	return p.OutputResult(model.OutputFormat, resp, func() error {
		operationState := "Cloned"
		if model.Async {
			operationState = "Triggered cloning of"
		}
		p.Outputf("%s zone for project %q. Zone ID: %s\n", operationState, projectLabel, utils.PtrString(resp.Zone.Id))
		return nil
	})
}
