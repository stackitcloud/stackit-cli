package list

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/dns/client"
	dnsUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/dns/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/dns"
)

const (
	activeFlag      = "active"
	inactiveFlag    = "inactive"
	zoneIdFlag      = "zone-id"
	deletedFlag     = "deleted"
	nameLikeFlag    = "name-like"
	orderByNameFlag = "order-by-name"
	limitFlag       = "limit"
	pageSizeFlag    = "page-size"

	pageSizeDefault      = 100
	deleteSucceededState = "DELETE_SUCCEEDED"
)

type inputModel struct {
	*globalflags.GlobalFlagModel

	Active      bool
	Inactive    bool
	ZoneId      string
	Deleted     bool
	NameLike    *string
	OrderByName *string
	Limit       *int64
	PageSize    int64
}

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List DNS record sets",
		Long:  `List DNS record sets. Successfully deleted record sets are not listed by default.`,
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`List DNS record-sets for zone with ID "xxx"`,
				"$ stackit dns record-set list --zone-id xxx"),
			examples.NewExample(
				`List DNS record-sets for zone with ID "xxx" in JSON format`,
				"$ stackit dns record-set list --zone-id xxx --output-format json"),
			examples.NewExample(
				`List active DNS record-sets for zone with ID "xxx"`,
				"$ stackit dns record-set list --zone-id xxx --is-active true"),
			examples.NewExample(
				`List up to 10 DNS record-sets for zone with ID "xxx"`,
				"$ stackit dns record-set list --zone-id xxx --limit 10"),
			examples.NewExample(
				`List the deleted DNS record-sets for zone with ID "xxx"`,
				"$ stackit dns record-set list --zone-id xxx --deleted"),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			model, err := parseInput(p, cmd)
			if err != nil {
				return err
			}

			// Configure API client
			apiClient, err := client.ConfigureClient(p)
			if err != nil {
				return err
			}

			// Fetch record sets
			recordSets, err := fetchRecordSets(ctx, model, apiClient)
			if err != nil {
				return err
			}
			if len(recordSets) == 0 {
				zoneLabel, err := dnsUtils.GetZoneName(ctx, apiClient, model.ProjectId, model.ZoneId)
				if err != nil {
					p.Debug(print.ErrorLevel, "get zone name: %v", err)
					zoneLabel = model.ZoneId
				}
				p.Info("No record sets found for zone %s matching the criteria\n", zoneLabel)
				return nil
			}
			return outputResult(p, model.OutputFormat, recordSets)
		},
	}

	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	orderByNameFlagOptions := []string{"asc", "desc"}

	cmd.Flags().Var(flags.UUIDFlag(), zoneIdFlag, "Zone ID")
	cmd.Flags().Bool(activeFlag, false, "Filter for active record sets")
	cmd.Flags().Bool(inactiveFlag, false, "Filter for inactive record sets. Deleted record sets are always inactive and will be included when this flag is set")
	cmd.Flags().Bool(deletedFlag, false, "Filter for deleted record sets")
	cmd.Flags().String(nameLikeFlag, "", "Filter by name")
	cmd.Flags().Var(flags.EnumFlag(true, "", orderByNameFlagOptions...), orderByNameFlag, fmt.Sprintf("Order by name, one of %q", orderByNameFlagOptions))
	cmd.Flags().Int64(limitFlag, 0, "Maximum number of entries to list")
	cmd.Flags().Int64(pageSizeFlag, pageSizeDefault, "Number of items fetched in each API call. Does not affect the number of items in the command output")

	err := flags.MarkFlagsRequired(cmd, zoneIdFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	limit := flags.FlagToInt64Pointer(p, cmd, limitFlag)
	if limit != nil && *limit < 1 {
		return nil, &errors.FlagValidationError{
			Flag:    limitFlag,
			Details: "must be greater than 0",
		}
	}

	pageSize := flags.FlagWithDefaultToInt64Value(p, cmd, pageSizeFlag)
	if pageSize < 1 {
		return nil, &errors.FlagValidationError{
			Flag:    pageSizeFlag,
			Details: "must be greater than 0",
		}
	}

	active := flags.FlagToBoolValue(p, cmd, activeFlag)
	inactive := flags.FlagToBoolValue(p, cmd, inactiveFlag)
	if active && inactive {
		return nil, fmt.Errorf("only one of %s and %s can be set", activeFlag, inactiveFlag)
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		ZoneId:          flags.FlagToStringValue(p, cmd, zoneIdFlag),
		Active:          active,
		Inactive:        inactive,
		Deleted:         flags.FlagToBoolValue(p, cmd, deletedFlag),
		NameLike:        flags.FlagToStringPointer(p, cmd, nameLikeFlag),
		OrderByName:     flags.FlagToStringPointer(p, cmd, orderByNameFlag),
		Limit:           flags.FlagToInt64Pointer(p, cmd, limitFlag),
		PageSize:        pageSize,
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

func buildRequest(ctx context.Context, model *inputModel, apiClient dnsClient, page int) dns.ApiListRecordSetsRequest {
	req := apiClient.ListRecordSets(ctx, model.ProjectId, model.ZoneId)
	if model.Active {
		req = req.ActiveEq(true)
	}
	if model.Inactive {
		req = req.ActiveEq(false)
	}
	if model.Deleted {
		req = req.StateEq(deleteSucceededState)
	} else if !model.Inactive {
		req = req.StateNeq(deleteSucceededState)
	}
	if model.NameLike != nil {
		req = req.NameLike(*model.NameLike)
	}
	if model.OrderByName != nil {
		req = req.OrderByName(strings.ToUpper(*model.OrderByName))
	}
	req = req.PageSize(int32(model.PageSize))
	req = req.Page(int32(page))
	return req
}

type dnsClient interface {
	ListRecordSets(ctx context.Context, projectId, zoneId string) dns.ApiListRecordSetsRequest
}

func fetchRecordSets(ctx context.Context, model *inputModel, apiClient dnsClient) ([]dns.RecordSet, error) {
	if model.Limit != nil && *model.Limit < model.PageSize {
		model.PageSize = *model.Limit
	}
	page := 1
	recordSets := []dns.RecordSet{}
	for {
		// Call API
		req := buildRequest(ctx, model, apiClient, page)
		resp, err := req.Execute()
		if err != nil {
			return nil, fmt.Errorf("get DNS record sets: %w", err)
		}
		respRecordSets := *resp.RrSets
		if len(respRecordSets) == 0 {
			break
		}
		recordSets = append(recordSets, respRecordSets...)
		// Stop if no more pages
		if len(respRecordSets) < int(model.PageSize) {
			break
		}
		// Stop and truncate if limit is reached
		if model.Limit != nil && len(recordSets) >= int(*model.Limit) {
			recordSets = recordSets[:*model.Limit]
			break
		}
		page++
	}
	return recordSets, nil
}

func outputResult(p *print.Printer, outputFormat string, recordSets []dns.RecordSet) error {
	switch outputFormat {
	case print.JSONOutputFormat:
		details, err := json.MarshalIndent(recordSets, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal DNS record set list: %w", err)
		}
		p.Outputln(string(details))

		return nil
	default:
		table := tables.NewTable()
		table.SetHeader("ID", "NAME", "STATUS", "TTL", "TYPE", "RECORD DATA")
		for i := range recordSets {
			rs := recordSets[i]
			recordData := make([]string, 0, len(*rs.Records))
			for _, r := range *rs.Records {
				recordData = append(recordData, *r.Content)
			}
			recordDataJoin := strings.Join(recordData, ", ")
			table.AddRow(*rs.Id, *rs.Name, *rs.State, *rs.Ttl, *rs.Type, recordDataJoin)
		}
		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}

		return nil
	}
}
