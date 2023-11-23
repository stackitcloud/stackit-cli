package update

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/commonflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/dns/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/dns"
	"github.com/stackitcloud/stackit-sdk-go/services/dns/wait"
)

const (
	zoneIdFlag      = "zone-id"
	recordSetIdFlag = "record-set-id"
	commentFlag     = "comment"
	nameFlag        = "name"
	recordFlag      = "record"
	ttlFlag         = "ttl"
)

type flagModel struct {
	ProjectId   string
	ZoneId      string
	RecordSetId string
	Comment     *string
	Name        *string
	Records     *[]string
	TTL         *int64
}

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "update",
		Short:   "Updates a DNS record set",
		Long:    "Updates a DNS record set. Performs a partial update; fields not provided are kept unchanged",
		Example: `$ stackit dns record-set update --project-id xxx --zone-id xxx --record-set-id xxx --name my-zone --type A --record 1.2.3.4 --record 5.6.7.8`,
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
			_, err = req.Execute()
			if err != nil {
				return fmt.Errorf("update DNS record set: %w", err)
			}

			// Wait for async operation
			_, err = wait.UpdateRecordSetWaitHandler(ctx, apiClient, model.ProjectId, model.ZoneId, model.RecordSetId).WaitWithContext(ctx)
			if err != nil {
				return fmt.Errorf("wait for DNS record set update: %w", err)
			}

			cmd.Println("Record set updated")
			return nil
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), zoneIdFlag, "Zone ID")
	cmd.Flags().Var(flags.UUIDFlag(), recordSetIdFlag, "Record set ID")
	cmd.Flags().String(commentFlag, "", "User comment")
	cmd.Flags().String(nameFlag, "", "Name of the record, should be compliant with RFC1035, Section 2.3.4")
	cmd.Flags().Int64(ttlFlag, 0, "Time to live, if not provided defaults to the zone's default TTL")
	cmd.Flags().StringSlice(recordFlag, []string{}, "Records belonging to the record set. If this flag is used, records already created that aren't set when running the command will be deleted")

	err := utils.MarkFlagsRequired(cmd, zoneIdFlag, recordSetIdFlag)
	cobra.CheckErr(err)
}

func parseFlags(cmd *cobra.Command) (*flagModel, error) {
	projectId := commonflags.GetString(commonflags.ProjectIdFlag)
	if projectId == "" {
		return nil, fmt.Errorf("project ID not set")
	}

	zoneId := utils.FlagToStringValue(cmd, zoneIdFlag)
	recordSetId := utils.FlagToStringValue(cmd, recordSetIdFlag)
	comment := utils.FlagToStringPointer(cmd, commentFlag)
	name := utils.FlagToStringPointer(cmd, nameFlag)
	records := utils.FlagToStringSlicePointer(cmd, recordFlag)
	ttl := utils.FlagToInt64Pointer(cmd, ttlFlag)

	if comment == nil && name == nil && records == nil && ttl == nil {
		return nil, fmt.Errorf("please specify at least one field to update")
	}

	return &flagModel{
		ProjectId:   projectId,
		ZoneId:      zoneId,
		RecordSetId: recordSetId,
		Comment:     comment,
		Name:        name,
		Records:     records,
		TTL:         ttl,
	}, nil
}

func buildRequest(ctx context.Context, model *flagModel, apiClient *dns.APIClient) dns.ApiUpdateRecordSetRequest {
	var records *[]dns.RecordPayload = nil
	if model.Records != nil {
		records = utils.Ptr(make([]dns.RecordPayload, 0))
		for _, r := range *model.Records {
			records = utils.Ptr(append(*records, dns.RecordPayload{Content: utils.Ptr(r)}))
		}
	}

	req := apiClient.UpdateRecordSet(ctx, model.ProjectId, model.ZoneId, model.RecordSetId)
	req = req.UpdateRecordSetPayload(dns.UpdateRecordSetPayload{
		Comment: model.Comment,
		Name:    model.Name,
		Records: records,
		Ttl:     model.TTL,
	})
	return req
}
