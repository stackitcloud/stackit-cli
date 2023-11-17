package list

import (
	"context"
	"fmt"
	"strings"

	"github.com/stackitcloud/stackit-cli/internal/pkg/config"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/dns/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stackitcloud/stackit-sdk-go/services/dns"
)

const (
	projectIdFlag   = "project-id"
	nameLikeFlag    = "name-like"
	activeFlag      = "is-active"
	orderByNameFlag = "order-by-name"
)

type flagModel struct {
	ProjectId   string
	NameLike    *string
	Active      *bool
	OrderByName *string
}

var Cmd = &cobra.Command{
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

		// Call API
		req := buildRequest(ctx, model, apiClient)
		resp, err := req.Execute()
		if err != nil {
			return fmt.Errorf("get DNS zones: %w", err)
		}
		zones := *resp.Zones
		if len(zones) == 0 {
			fmt.Printf("No zones found for project with ID %s\n", model.ProjectId)
			return nil
		}

		// Show output as table
		table := tables.NewTable()
		table.SetHeader("ID", "NAME", "DNS_NAME", "STATE")
		for _, zone := range zones {
			table.AddRow(*zone.Id, *zone.Name, *zone.DnsName, *zone.State)
		}
		table.Render()

		return nil
	},
}

func init() {
	configureFlags(Cmd)
}

func configureFlags(cmd *cobra.Command) {
	activeFlagOptions := []string{"true", "false"}
	orderByNameFlagOptions := []string{"asc", "desc"}

	cmd.Flags().String(nameLikeFlag, "", "Filter by name")
	cmd.Flags().Var(flags.EnumBoolFlag(), activeFlag, fmt.Sprintf("Filter by active status, one of %q", activeFlagOptions))
	cmd.Flags().Var(flags.EnumFlag(true, orderByNameFlagOptions...), orderByNameFlag, fmt.Sprintf("Order by name, one of %q", orderByNameFlagOptions))
}

func parseFlags(cmd *cobra.Command) (*flagModel, error) {
	projectId := viper.GetString(config.ProjectIdKey)
	if projectId == "" {
		return nil, fmt.Errorf("project ID not set")
	}

	return &flagModel{
		ProjectId:   projectId,
		NameLike:    utils.FlagToStringPointer(cmd, nameLikeFlag),
		Active:      utils.FlagToBoolPointer(cmd, activeFlag),
		OrderByName: utils.FlagToStringPointer(cmd, orderByNameFlag),
	}, nil
}

func buildRequest(ctx context.Context, model *flagModel, apiClient *dns.APIClient) dns.ApiGetZonesRequest {
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
	return req
}
