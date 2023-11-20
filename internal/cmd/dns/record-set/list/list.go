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
	zoneIdFlag      = "zone-id"
	nameLikeFlag    = "name-like"
	activeFlag      = "is-active"
	orderByNameFlag = "order-by-name"
)

type flagModel struct {
	ProjectId   string
	ZoneId      string
	NameLike    *string
	Active      *bool
	OrderByName *string
}

var Cmd = &cobra.Command{
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

		// Call API
		req := buildRequest(ctx, model, apiClient)
		resp, err := req.Execute()
		if err != nil {
			return fmt.Errorf("get DNS record sets: %w", err)
		}
		recordSets := *resp.RrSets
		if len(recordSets) == 0 {
			cmd.Printf("No record-sets found for zone with ID %s\n", model.ZoneId)
			return nil
		}

		// Show output as table
		table := tables.NewTable()
		table.SetHeader("ID", "Name", "Type", "State")
		for _, recordSet := range recordSets {
			table.AddRow(*recordSet.Id, *recordSet.Name, *recordSet.Type, *recordSet.State)
		}
		table.Render(cmd)

		return nil
	},
}

func init() {
	configureFlags(Cmd)
}

func configureFlags(cmd *cobra.Command) {
	activeFlagOptions := []string{"true", "false"}
	orderByNameFlagOptions := []string{"asc", "desc"}

	cmd.Flags().Var(flags.UUIDFlag(), zoneIdFlag, "Zone ID")
	cmd.Flags().String(nameLikeFlag, "", "Filter by name")
	cmd.Flags().Var(flags.EnumBoolFlag(), activeFlag, fmt.Sprintf("Filter by active status, one of %q", activeFlagOptions))
	cmd.Flags().Var(flags.EnumFlag(true, orderByNameFlagOptions...), orderByNameFlag, fmt.Sprintf("Order by name, one of %q", orderByNameFlagOptions))

	err := utils.MarkFlagsRequired(cmd, zoneIdFlag)
	cobra.CheckErr(err)
}

func parseFlags(cmd *cobra.Command) (*flagModel, error) {
	projectId := viper.GetString(config.ProjectIdKey)
	if projectId == "" {
		return nil, fmt.Errorf("project ID not set")
	}

	return &flagModel{
		ProjectId:   projectId,
		ZoneId:      utils.FlagToStringValue(cmd, zoneIdFlag),
		NameLike:    utils.FlagToStringPointer(cmd, nameLikeFlag),
		Active:      utils.FlagToBoolPointer(cmd, activeFlag),
		OrderByName: utils.FlagToStringPointer(cmd, orderByNameFlag),
	}, nil
}

func buildRequest(ctx context.Context, model *flagModel, apiClient *dns.APIClient) dns.ApiGetRecordSetsRequest {
	req := apiClient.GetRecordSets(ctx, model.ProjectId, model.ZoneId)
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
