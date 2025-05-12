package update

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
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
	recordSetIdArg = "RECORD_SET_ID"

	zoneIdFlag  = "zone-id"
	commentFlag = "comment"
	nameFlag    = "name"
	recordFlag  = "record"
	ttlFlag     = "ttl"
	txtType     = "TXT"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	ZoneId      string
	RecordSetId string
	Comment     *string
	Name        *string
	Records     *[]string
	TTL         *int64
	Type        *string
}

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("update %s", recordSetIdArg),
		Short: "Updates a DNS record set",
		Long:  "Updates a DNS record set.",
		Args:  args.SingleArg(recordSetIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Update the time to live of the record-set with ID "xxx" for zone with ID "yyy"`,
				"$ stackit dns record-set update xxx --zone-id yyy --ttl 100"),
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

			recordSetLabel, err := dnsUtils.GetRecordSetName(ctx, apiClient, model.ProjectId, model.ZoneId, model.RecordSetId)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get record set name: %v", err)
				recordSetLabel = model.RecordSetId
			}

			typeLabel, err := dnsUtils.GetRecordSetType(ctx, apiClient, model.ProjectId, model.ZoneId, model.RecordSetId)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get record set type: %v", err)
			}
			model.Type = typeLabel

			if utils.PtrString(model.Type) == txtType {
				err = parseTxtRecord(model.Records)
				if err != nil {
					return err
				}
			}

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to update record set %s of zone %s?", recordSetLabel, zoneLabel)
				err = params.Printer.PromptForConfirmation(prompt)
				if err != nil {
					return err
				}
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			_, err = req.Execute()
			if err != nil {
				return fmt.Errorf("update DNS record set: %w", err)
			}

			// Wait for async operation, if async mode not enabled
			if !model.Async {
				s := spinner.New(params.Printer)
				s.Start("Updating record set")
				_, err = wait.PartialUpdateRecordSetWaitHandler(ctx, apiClient, model.ProjectId, model.ZoneId, model.RecordSetId).WaitWithContext(ctx)
				if err != nil {
					return fmt.Errorf("wait for DNS record set update: %w", err)
				}
				s.Stop()
			}

			operationState := "Updated"
			if model.Async {
				operationState = "Triggered update of"
			}
			params.Printer.Info("%s record set %s of zone %s\n", operationState, recordSetLabel, zoneLabel)
			return nil
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), zoneIdFlag, "Zone ID")
	cmd.Flags().String(commentFlag, "", "User comment")
	cmd.Flags().String(nameFlag, "", "Name of the record, should be compliant with RFC1035, Section 2.3.4")
	cmd.Flags().Int64(ttlFlag, 0, "Time to live, if not provided defaults to the zone's default TTL")
	cmd.Flags().StringSlice(recordFlag, []string{}, "Records belonging to the record set. If this flag is used, records already created that aren't set when running the command will be deleted")

	err := flags.MarkFlagsRequired(cmd, zoneIdFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	recordSetId := inputArgs[0]

	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	zoneId := flags.FlagToStringValue(p, cmd, zoneIdFlag)
	comment := flags.FlagToStringPointer(p, cmd, commentFlag)
	name := flags.FlagToStringPointer(p, cmd, nameFlag)
	records := flags.FlagToStringSlicePointer(p, cmd, recordFlag)
	ttl := flags.FlagToInt64Pointer(p, cmd, ttlFlag)

	if comment == nil && name == nil && records == nil && ttl == nil {
		return nil, &errors.EmptyUpdateError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		ZoneId:          zoneId,
		RecordSetId:     recordSetId,
		Comment:         comment,
		Name:            name,
		Records:         records,
		TTL:             ttl,
	}

	if p.IsVerbosityDebug() {
		modelStr, err := print.BuildDebugStrFromInputModel(model)
		if err != nil {
			p.Debug(print.ErrorLevel, "convert model to string for debugging: %v", err)
		} else {
			p.Debug(print.DebugLevel, "parsed input values: %s", modelStr)
		}
	}

	return &model, nil
}

func parseTxtRecord(records *[]string) error {
	if records == nil {
		return nil
	}
	if len(*records) == 0 {
		return nil
	}

	for idx := range *records {
		var err error
		// Based on RFC 1035 section 2.3.4, TXT Records are limited to 255 Characters.
		// Longer strings need to be split into multiple records
		(*records)[idx], err = dnsUtils.FormatTxtRecord((*records)[idx])
		if err != nil {
			return err
		}
	}

	return nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *dns.APIClient) dns.ApiPartialUpdateRecordSetRequest {
	var records *[]dns.RecordPayload = nil
	if model.Records != nil {
		records = utils.Ptr(make([]dns.RecordPayload, 0))
		for _, r := range *model.Records {
			records = utils.Ptr(append(*records, dns.RecordPayload{Content: utils.Ptr(r)}))
		}
	}

	req := apiClient.PartialUpdateRecordSet(ctx, model.ProjectId, model.ZoneId, model.RecordSetId)
	req = req.PartialUpdateRecordSetPayload(dns.PartialUpdateRecordSetPayload{
		Comment: model.Comment,
		Name:    model.Name,
		Records: records,
		Ttl:     model.TTL,
	})
	return req
}
