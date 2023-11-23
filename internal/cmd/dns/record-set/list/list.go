package list

import (
	"context"
	"fmt"
	"strings"

	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/dns/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/dns"
)

const (
	zoneIdFlag      = "zone-id"
	nameLikeFlag    = "name-like"
	activeFlag      = "is-active"
	orderByNameFlag = "order-by-name"
	limitFlag       = "limit"
	pageSizeFlag    = "page-size"
)

type flagModel struct {
	GlobalFlags *globalflags.Model
	ZoneId      string
	NameLike    *string
	Active      *bool
	OrderByName *string
	Limit       *int64
	PageSize    int64
}

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List all DNS record sets",
		Long:    "List all DNS record sets",
		Example: `$ stackit dns record-set list --project-id xxx --zone-id xxx`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			model, err := parseFlags(cmd)
			if err != nil {
				return err
			}

			// Configure API client
			apiClient, err := client.ConfigureClient(cmd)
			if err != nil {
				return fmt.Errorf("authentication failed, please run \"stackit auth login\" or \"stackit auth activate-service-account\"")
			}

			// Fetch record sets
			recordSets, err := fetchRecordSets(ctx, model, apiClient)
			if err != nil {
				return err
			}
			if len(recordSets) == 0 {
				cmd.Printf("No record sets found for zone %s in project with ID %s\n", model.ZoneId, model.GlobalFlags.ProjectId)
				return nil
			}

			// Show output as table
			table := tables.NewTable()
			table.SetHeader("ID", "Name", "Type", "State")
			for i := range recordSets {
				rs := recordSets[i]
				table.AddRow(*rs.Id, *rs.Name, *rs.Type, *rs.State)
			}
			table.Render(cmd)

			return nil
		},
	}

	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	activeFlagOptions := []string{"true", "false"}
	orderByNameFlagOptions := []string{"asc", "desc"}

	cmd.Flags().Var(flags.UUIDFlag(), zoneIdFlag, "Zone ID")
	cmd.Flags().String(nameLikeFlag, "", "Filter by name")
	cmd.Flags().Var(flags.EnumBoolFlag(), activeFlag, fmt.Sprintf("Filter by active status, one of %q", activeFlagOptions))
	cmd.Flags().Var(flags.EnumFlag(true, orderByNameFlagOptions...), orderByNameFlag, fmt.Sprintf("Order by name, one of %q", orderByNameFlagOptions))
	cmd.Flags().Int64(limitFlag, 0, "Maximum number of entries to list")
	cmd.Flags().Int64(pageSizeFlag, 100, "Number of items fetched in each API call. Does not affect the number of items in the command output")

	err := utils.MarkFlagsRequired(cmd, zoneIdFlag)
	cobra.CheckErr(err)
}

func parseFlags(cmd *cobra.Command) (*flagModel, error) {
	globalFlags := globalflags.Parse()
	if globalFlags.ProjectId == "" {
		return nil, fmt.Errorf("project ID not set")
	}

	limit := utils.FlagToInt64Pointer(cmd, limitFlag)
	if limit != nil && *limit < 1 {
		return nil, fmt.Errorf("limit must be greater than 0")
	}

	pageSize, err := utils.FlagWithDefaultToInt64Value(cmd, pageSizeFlag)
	if err != nil {
		return nil, fmt.Errorf("parse %s flag: %w", pageSizeFlag, err)
	}
	if pageSize < 1 {
		return nil, fmt.Errorf("page size must be greater than 0")
	}

	return &flagModel{
		GlobalFlags: globalFlags,
		ZoneId:      utils.FlagToStringValue(cmd, zoneIdFlag),
		NameLike:    utils.FlagToStringPointer(cmd, nameLikeFlag),
		Active:      utils.FlagToBoolPointer(cmd, activeFlag),
		OrderByName: utils.FlagToStringPointer(cmd, orderByNameFlag),
		Limit:       utils.FlagToInt64Pointer(cmd, limitFlag),
		PageSize:    pageSize,
	}, nil
}

func buildRequest(ctx context.Context, model *flagModel, apiClient dnsClient, page int) dns.ApiGetRecordSetsRequest {
	req := apiClient.GetRecordSets(ctx, model.GlobalFlags.ProjectId, model.ZoneId)
	if model.NameLike != nil {
		req = req.NameLike(*model.NameLike)
	}
	if model.Active != nil {
		req = req.ActiveEq(*model.Active)
	}
	if model.OrderByName != nil {
		req = req.OrderByName(strings.ToUpper(*model.OrderByName))
	}
	req = req.PageSize(int32(model.PageSize))
	req = req.Page(int32(page))
	return req
}

type dnsClient interface {
	GetRecordSets(ctx context.Context, projectId, zoneId string) dns.ApiGetRecordSetsRequest
}

func fetchRecordSets(ctx context.Context, model *flagModel, apiClient dnsClient) ([]dns.RecordSet, error) {
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
