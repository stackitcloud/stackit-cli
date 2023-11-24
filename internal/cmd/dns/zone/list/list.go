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
	nameLikeFlag    = "name-like"
	activeFlag      = "is-active"
	orderByNameFlag = "order-by-name"
	limitFlag       = "limit"
	pageSizeFlag    = "page-size"
)

type flagModel struct {
	ProjectId   string
	NameLike    *string
	Active      *bool
	OrderByName *string
	Limit       *int64
	PageSize    int64
}

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List all DNS zones",
		Long:    "List all DNS zones",
		Example: `$ stackit dns zone list --project-id xxx`,
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

			// Fetch zones
			zones, err := fetchZones(ctx, model, apiClient)
			if err != nil {
				return err
			}
			if len(zones) == 0 {
				cmd.Printf("No zones found for project with ID %s\n", model.ProjectId)
				return nil
			}

			// Show output as table
			table := tables.NewTable()
			table.SetHeader("ID", "NAME", "DNS_NAME", "STATE")
			for i := range zones {
				z := zones[i]
				table.AddRow(*z.Id, *z.Name, *z.DnsName, *z.State)
			}
			err = table.Render()
			if err != nil {
				return fmt.Errorf("render table: %w", err)
			}

			return nil
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	activeFlagOptions := []string{"true", "false"}
	orderByNameFlagOptions := []string{"asc", "desc"}

	cmd.Flags().String(nameLikeFlag, "", "Filter by name")
	cmd.Flags().Var(flags.EnumBoolFlag(), activeFlag, fmt.Sprintf("Filter by active status, one of %q", activeFlagOptions))
	cmd.Flags().Var(flags.EnumFlag(true, orderByNameFlagOptions...), orderByNameFlag, fmt.Sprintf("Order by name, one of %q", orderByNameFlagOptions))
	cmd.Flags().Int64(limitFlag, 0, "Maximum number of entries to list")
	cmd.Flags().Int64(pageSizeFlag, 100, "Number of items fetched in each API call. Does not affect the number of items in the command output")
}

func parseFlags(cmd *cobra.Command) (*flagModel, error) {
	projectId := globalflags.GetString(globalflags.ProjectIdFlag)
	if projectId == "" {
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
		ProjectId:   projectId,
		NameLike:    utils.FlagToStringPointer(cmd, nameLikeFlag),
		Active:      utils.FlagToBoolPointer(cmd, activeFlag),
		OrderByName: utils.FlagToStringPointer(cmd, orderByNameFlag),
		Limit:       limit,
		PageSize:    pageSize,
	}, nil
}

func buildRequest(ctx context.Context, model *flagModel, apiClient dnsClient, page int) dns.ApiGetZonesRequest {
	req := apiClient.GetZones(ctx, model.ProjectId)
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
	GetZones(ctx context.Context, projectId string) dns.ApiGetZonesRequest
}

func fetchZones(ctx context.Context, model *flagModel, apiClient dnsClient) ([]dns.Zone, error) {
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
