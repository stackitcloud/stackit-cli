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
	"github.com/stackitcloud/stackit-cli/internal/pkg/projectname"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/dns/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/dns"
)

const (
	activeFlag         = "active"
	inactiveFlag       = "inactive"
	nameLikeFlag       = "name-like"
	orderByNameFlag    = "order-by-name"
	includeDeletedFlag = "include-deleted"
	limitFlag          = "limit"
	pageSizeFlag       = "page-size"

	pageSizeDefault      = 100
	deleteSucceededState = "DELETE_SUCCEEDED"
)

type inputModel struct {
	*globalflags.GlobalFlagModel

	Active         bool
	Inactive       bool
	NameLike       *string
	OrderByName    *string
	IncludeDeleted bool
	Limit          *int64
	PageSize       int64
}

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List DNS zones",
		Long:  `List DNS zones. Successfully deleted zones are not listed by default.`,
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`List DNS zones`,
				"$ stackit dns zone list"),
			examples.NewExample(
				`List DNS zones in JSON format`,
				"$ stackit dns zone list --output-format json"),
			examples.NewExample(
				`List up to 10 DNS zones`,
				"$ stackit dns zone list --limit 10"),
			examples.NewExample(
				`List DNS zones, including deleted`,
				"$ stackit dns zone list --include-deleted"),
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

			// Fetch zones
			zones, err := fetchZones(ctx, model, apiClient)
			if err != nil {
				return err
			}
			if len(zones) == 0 {
				projectLabel, err := projectname.GetProjectName(ctx, cmd, p)
				if err != nil {
					p.Debug(print.ErrorLevel, "get project name: %v", err)
					projectLabel = model.ProjectId
				}
				p.Info("No zones found for project %q matching the criteria\n", projectLabel)
				return nil
			}

			return outputResult(p, model.OutputFormat, zones)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	orderByNameFlagOptions := []string{"asc", "desc"}

	cmd.Flags().Bool(activeFlag, false, "Filter for active zones")
	cmd.Flags().Bool(inactiveFlag, false, "Filter for inactive zones")
	cmd.Flags().String(nameLikeFlag, "", "Filter by name")
	cmd.Flags().Var(flags.EnumFlag(true, "", orderByNameFlagOptions...), orderByNameFlag, fmt.Sprintf("Order by name, one of %q", orderByNameFlagOptions))
	cmd.Flags().Bool(includeDeletedFlag, false, "Includes successfully deleted zones (if unset, these are filtered out)")
	cmd.Flags().Int64(limitFlag, 0, "Maximum number of entries to list")
	cmd.Flags().Int64(pageSizeFlag, pageSizeDefault, "Number of items fetched in each API call. Does not affect the number of items in the command output")
}

func parseInput(cmd *cobra.Command) (*inputModel, error) {
	globalFlags := globalflags.Parse(cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	limit := flags.FlagToInt64Pointer(cmd, limitFlag)
	if limit != nil && *limit < 1 {
		return nil, &errors.FlagValidationError{
			Flag:    limitFlag,
			Details: "must be greater than 0",
		}
	}

	pageSize := flags.FlagWithDefaultToInt64Value(cmd, pageSizeFlag)
	if pageSize < 1 {
		return nil, &errors.FlagValidationError{
			Flag:    pageSizeFlag,
			Details: "must be greater than 0",
		}
	}

	active := flags.FlagToBoolValue(cmd, activeFlag)
	inactive := flags.FlagToBoolValue(cmd, inactiveFlag)
	if active && inactive {
		return nil, fmt.Errorf("only one of %s and %s can be set", activeFlag, inactiveFlag)
	}

	return &inputModel{
		GlobalFlagModel: globalFlags,
		Active:          active,
		Inactive:        inactive,
		IncludeDeleted:  flags.FlagToBoolValue(cmd, includeDeletedFlag),
		NameLike:        flags.FlagToStringPointer(cmd, nameLikeFlag),
		OrderByName:     flags.FlagToStringPointer(cmd, orderByNameFlag),
		Limit:           limit,
		PageSize:        pageSize,
	}, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient dnsClient, page int) dns.ApiListZonesRequest {
	req := apiClient.ListZones(ctx, model.ProjectId)
	if model.Active {
		req = req.ActiveEq(true)
	}
	if model.Inactive {
		req = req.ActiveEq(false)
	}
	if model.NameLike != nil {
		req = req.NameLike(*model.NameLike)
	}
	if model.OrderByName != nil {
		req = req.OrderByName(strings.ToUpper(*model.OrderByName))
	}
	if !model.IncludeDeleted {
		req = req.StateNeq(deleteSucceededState)
	}
	req = req.PageSize(int32(model.PageSize))
	req = req.Page(int32(page))
	return req
}

type dnsClient interface {
	ListZones(ctx context.Context, projectId string) dns.ApiListZonesRequest
}

func fetchZones(ctx context.Context, model *inputModel, apiClient dnsClient) ([]dns.Zone, error) {
	if model.Limit != nil && *model.Limit < model.PageSize {
		model.PageSize = *model.Limit
	}
	page := 1
	zones := []dns.Zone{}
	for {
		// Call API
		req := buildRequest(ctx, model, apiClient, page)
		resp, err := req.Execute()
		if err != nil {
			return nil, fmt.Errorf("get DNS zones: %w", err)
		}
		respZones := *resp.Zones
		if len(respZones) == 0 {
			break
		}
		zones = append(zones, respZones...)
		// Stop if no more pages
		if len(respZones) < int(model.PageSize) {
			break
		}
		// Stop and truncate if limit is reached
		if model.Limit != nil && len(zones) >= int(*model.Limit) {
			zones = zones[:*model.Limit]
			break
		}
		page++
	}
	return zones, nil
}

func outputResult(p *print.Printer, outputFormat string, zones []dns.Zone) error {
	switch outputFormat {
	case globalflags.JSONOutputFormat:
		// Show details
		details, err := json.MarshalIndent(zones, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal DNS zone list: %w", err)
		}
		p.Outputln(string(details))

		return nil
	default:
		table := tables.NewTable()
		table.SetHeader("ID", "NAME", "STATE", "TYPE", "DNS NAME", "RECORD COUNT")
		for i := range zones {
			z := zones[i]
			table.AddRow(*z.Id, *z.Name, *z.State, *z.Type, *z.DnsName, *z.RecordCount)
		}
		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}

		return nil
	}
}
