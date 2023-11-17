package create

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/config"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/dns/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stackitcloud/stackit-sdk-go/services/dns"
	"github.com/stackitcloud/stackit-sdk-go/services/dns/wait"
)

const (
	projectIdFlag = "project-id"
	zoneIdFlag    = "zone-id"
	commentFlag   = "comment"
	nameFlag      = "name"
	recordFlag    = "record"
	ttlFlag       = "ttl"
	typeFlag      = "type"
)

type flagModel struct {
	ProjectId string
	ZoneId    string
	Comment   *string
	Name      *string
	Records   []string
	TTL       *int64
	Type      *string
}

var Cmd = &cobra.Command{
	Use:     "create",
	Short:   "Creates a DNS record set",
	Long:    "Creates a DNS record set",
	Example: `$ stackit dns record-set create --project-id xxx --zone-id xxx --name my-zone --type A --record 1.2.3.4 --record 5.6.7.8`,
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
			return fmt.Errorf("create DNS record set: %w", err)
		}

		// Wait for async operation
		recordSetId := *resp.Rrset.Id
		_, err = wait.CreateRecordSetWaitHandler(ctx, apiClient, model.ProjectId, model.ZoneId, recordSetId).WaitWithContext(ctx)
		if err != nil {
			return fmt.Errorf("wait for DNS record set creation: %w", err)
		}

		fmt.Printf("Created record set with ID %s\n", recordSetId)
		return nil
	},
}

func init() {
	configureFlags(Cmd)
}

func configureFlags(cmd *cobra.Command) {
	typeFlagOptions := []string{"A", "AAAA", "SOA", "CNAME", "NS", "MX", "TXT", "SRV", "PTR", "ALIAS", "DNAME", "CAA"}

	cmd.Flags().Var(flags.UUIDFlag(), zoneIdFlag, "Zone ID")
	cmd.Flags().String(commentFlag, "", "User comment")
	cmd.Flags().String(nameFlag, "", "Name of the record, should be compliant with RFC1035, Section 2.3.4")
	cmd.Flags().Int64(ttlFlag, 0, "Time to live, if not provided defaults to the zone's default TTL")
	cmd.Flags().StringSlice(recordFlag, []string{}, "Records belonging to the record set")
	cmd.Flags().Var(flags.EnumFlag(false, typeFlagOptions...), typeFlag, fmt.Sprintf("Zone type, one of %q", typeFlagOptions))

	err := utils.MarkFlagsRequired(cmd, zoneIdFlag, nameFlag, recordFlag, typeFlag)
	cobra.CheckErr(err)
}

func parseFlags(cmd *cobra.Command) (*flagModel, error) {
	projectId := viper.GetString(config.ProjectIdKey)
	if projectId == "" {
		return nil, fmt.Errorf("project ID not set")
	}

	return &flagModel{
		ProjectId: projectId,
		ZoneId:    utils.FlagToStringValue(cmd, zoneIdFlag),
		Comment:   utils.FlagToStringPointer(cmd, commentFlag),
		Name:      utils.FlagToStringPointer(cmd, nameFlag),
		Records:   utils.FlagToStringSliceValue(cmd, recordFlag),
		TTL:       utils.FlagToInt64Pointer(cmd, ttlFlag),
		Type:      utils.FlagToStringPointer(cmd, typeFlag),
	}, nil
}

func buildRequest(ctx context.Context, model *flagModel, apiClient *dns.APIClient) dns.ApiCreateRecordSetRequest {
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
		Type:    model.Type,
	})
	return req
}
