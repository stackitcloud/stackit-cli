package create

import (
	"context"
	"fmt"

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
	zoneIdFlag  = "zone-id"
	commentFlag = "comment"
	nameFlag    = "name"
	recordFlag  = "record"
	ttlFlag     = "ttl"
	typeFlag    = "type"

	defaultType = "A"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	ZoneId  string
	Comment *string
	Name    *string
	Records []string
	TTL     *int64
	Type    string
}

func NewCmd(p *print.Printer) *cobra.Command {
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
			model, err := parseInput(cmd)
			if err != nil {
				return err
			}

			// Configure API client
			apiClient, err := client.ConfigureClient(p)
			if err != nil {
				return err
			}

			zoneLabel, err := dnsUtils.GetZoneName(ctx, apiClient, model.ProjectId, model.ZoneId)
			if err != nil {
				p.Debug(print.ErrorLevel, "get zone name: %v", err)
				zoneLabel = model.ZoneId
			}

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to create a record set for zone %s?", zoneLabel)
				err = p.PromptForConfirmation(prompt)
				if err != nil {
					return err
				}
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
				s := spinner.New(p)
				s.Start("Creating record set")
				_, err = wait.CreateRecordSetWaitHandler(ctx, apiClient, model.ProjectId, model.ZoneId, recordSetId).WaitWithContext(ctx)
				if err != nil {
					return fmt.Errorf("wait for DNS record set creation: %w", err)
				}
				s.Stop()
			}

			operationState := "Created"
			if model.Async {
				operationState = "Triggered creation of"
			}
			p.Outputf("%s record set for zone %s. Record set ID: %s\n", operationState, zoneLabel, recordSetId)
			return nil
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	typeFlagOptions := []string{"A", "AAAA", "SOA", "CNAME", "NS", "MX", "TXT", "SRV", "PTR", "ALIAS", "DNAME", "CAA"}

	cmd.Flags().Var(flags.UUIDFlag(), zoneIdFlag, "Zone ID")
	cmd.Flags().String(commentFlag, "", "User comment")
	cmd.Flags().String(nameFlag, "", "Name of the record, should be compliant with RFC1035, Section 2.3.4")
	cmd.Flags().Int64(ttlFlag, 0, "Time to live, if not provided defaults to the zone's default TTL")
	cmd.Flags().StringSlice(recordFlag, []string{}, "Records belonging to the record set")
	cmd.Flags().Var(flags.EnumFlag(false, defaultType, typeFlagOptions...), typeFlag, fmt.Sprintf("Record type, one of %q", typeFlagOptions))

	err := flags.MarkFlagsRequired(cmd, zoneIdFlag, nameFlag, recordFlag)
	cobra.CheckErr(err)
}

func parseInput(cmd *cobra.Command) (*inputModel, error) {
	globalFlags := globalflags.Parse(cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	return &inputModel{
		GlobalFlagModel: globalFlags,
		ZoneId:          flags.FlagToStringValue(cmd, zoneIdFlag),
		Comment:         flags.FlagToStringPointer(cmd, commentFlag),
		Name:            flags.FlagToStringPointer(cmd, nameFlag),
		Records:         flags.FlagToStringSliceValue(cmd, recordFlag),
		TTL:             flags.FlagToInt64Pointer(cmd, ttlFlag),
		Type:            flags.FlagWithDefaultToStringValue(cmd, typeFlag),
	}, nil
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
