package create

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/types"

	"github.com/spf13/cobra"
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
	"github.com/stackitcloud/stackit-sdk-go/services/dns"
	"github.com/stackitcloud/stackit-sdk-go/services/dns/wait"
)

const (
	zoneIdFlag  = "zone-id"
	commentFlag = "comment"
	nameFlag    = "name"
	recordFlag  = "record"
	ttlFlag     = "ttl"
	typeFlag    = "type"

	defaultType = dns.CREATERECORDSETPAYLOADTYPE_A
	txtType     = dns.CREATERECORDSETPAYLOADTYPE_TXT
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	ZoneId  string
	Comment *string
	Name    *string
	Records []string
	TTL     *int64
	Type    dns.CreateRecordSetPayloadTypes
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Creates a DNS record set",
		Long:  "Creates a DNS record set.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Create a DNS record set with name "my-rr" with records "1.2.3.4" and "5.6.7.8" in zone with ID "xxx"`,
				"$ stackit dns record-set create --zone-id xxx --name my-rr --record 1.2.3.4 --record 5.6.7.8"),
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

			prompt := fmt.Sprintf("Are you sure you want to create a record set for zone %s?", zoneLabel)
			err = params.Printer.PromptForConfirmation(prompt)
			if err != nil {
				return err
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("create DNS record set: %w", err)
			}
			recordSetId := *resp.Rrset.Id

			// Wait for async operation, if async mode not enabled
			if !model.Async {
				s := spinner.New(params.Printer)
				s.Start("Creating record set")
				_, err = wait.CreateRecordSetWaitHandler(ctx, apiClient, model.ProjectId, model.ZoneId, recordSetId).WaitWithContext(ctx)
				if err != nil {
					return fmt.Errorf("wait for DNS record set creation: %w", err)
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
	var typeFlagOptions []string
	for _, val := range dns.AllowedCreateRecordSetPayloadTypesEnumValues {
		typeFlagOptions = append(typeFlagOptions, string(val))
	}

	cmd.Flags().Var(flags.UUIDFlag(), zoneIdFlag, "Zone ID")
	cmd.Flags().String(commentFlag, "", "User comment")
	cmd.Flags().String(nameFlag, "", "Name of the record, should be compliant with RFC1035, Section 2.3.4")
	cmd.Flags().Int64(ttlFlag, 0, "Time to live, if not provided defaults to the zone's default TTL")
	cmd.Flags().StringSlice(recordFlag, []string{}, "Records belonging to the record set")
	cmd.Flags().Var(flags.EnumFlag(false, string(defaultType), typeFlagOptions...), typeFlag, fmt.Sprintf("Record type, one of %q", typeFlagOptions))

	err := flags.MarkFlagsRequired(cmd, zoneIdFlag, nameFlag, recordFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command, _ []string) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	recordType := flags.FlagWithDefaultToStringValue(p, cmd, typeFlag)

	model := inputModel{
		GlobalFlagModel: globalFlags,
		ZoneId:          flags.FlagToStringValue(p, cmd, zoneIdFlag),
		Comment:         flags.FlagToStringPointer(p, cmd, commentFlag),
		Name:            flags.FlagToStringPointer(p, cmd, nameFlag),
		Records:         flags.FlagToStringSliceValue(p, cmd, recordFlag),
		TTL:             flags.FlagToInt64Pointer(p, cmd, ttlFlag),
		Type:            dns.CreateRecordSetPayloadTypes(recordType),
	}

	if model.Type == txtType {
		for idx := range model.Records {
			// Based on RFC 1035 section 2.3.4, TXT Records are limited to 255 Characters
			// Longer strings need to be split into multiple records
			if len(model.Records[idx]) > 255 {
				var err error
				model.Records[idx], err = dnsUtils.FormatTxtRecord(model.Records[idx])
				if err != nil {
					return nil, err
				}
			}
		}
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *dns.APIClient) dns.ApiCreateRecordSetRequest {
	records := make([]dns.RecordPayload, 0)
	for _, r := range model.Records {
		records = append(records, dns.RecordPayload{Content: utils.Ptr(r)})
	}

	req := apiClient.CreateRecordSet(ctx, model.ProjectId, model.ZoneId)
	req = req.CreateRecordSetPayload(dns.CreateRecordSetPayload{
		Comment: model.Comment,
		Name:    model.Name,
		Records: &records,
		Ttl:     model.TTL,
		Type:    &model.Type,
	})
	return req
}

func outputResult(p *print.Printer, model *inputModel, zoneLabel string, resp *dns.RecordSetResponse) error {
	if resp == nil {
		return fmt.Errorf("record set response is empty")
	}
	return p.OutputResult(model.OutputFormat, resp, func() error {
		operationState := "Created"
		if model.Async {
			operationState = "Triggered creation of"
		}
		p.Outputf("%s record set for zone %s. Record set ID: %s\n", operationState, zoneLabel, utils.PtrString(resp.Rrset.Id))
		return nil
	})
}
