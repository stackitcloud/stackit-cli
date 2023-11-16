package describe

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"stackit/internal/pkg/args"
	"stackit/internal/pkg/errors"
	"stackit/internal/pkg/examples"
	"stackit/internal/pkg/flags"
	"stackit/internal/pkg/globalflags"
	"stackit/internal/pkg/services/dns/client"
	"stackit/internal/pkg/tables"
	"stackit/internal/pkg/utils"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/dns"
)

const (
	recordSetIdArg = "RECORD_SET_ID"

	zoneIdFlag = "zone-id"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	ZoneId      string
	RecordSetId string
}

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("describe %s", recordSetIdArg),
		Short: "Get details of a DNS record set",
		Long:  "Get details of a DNS record set",
		Args:  args.SingleArg(recordSetIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Get details of DNS record set with ID "xxx" in zone with ID "yyy"`,
				"$ stackit dns record-set describe xxx --zone-id yyy"),
			examples.NewExample(
				`Get details of DNS record set with ID "xxx" in zone with ID "yyy" in a table format`,
				"$ stackit dns record-set describe xxx --zone-id yyy --output-format pretty"),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			model, err := parseInput(cmd, args)
			if err != nil {
				return err
			}

			// Configure API client
			apiClient, err := client.ConfigureClient(cmd)
			if err != nil {
				return err
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("read DNS record set: %w", err)
			}
			recordSet := resp.Rrset

			return outputResult(cmd, model.OutputFormat, recordSet)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), zoneIdFlag, "Zone ID")

	err := flags.MarkFlagsRequired(cmd, zoneIdFlag)
	cobra.CheckErr(err)
}

func parseInput(cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	recordSetId := inputArgs[0]

	globalFlags := globalflags.Parse(cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	return &inputModel{
		GlobalFlagModel: globalFlags,
		ZoneId:          flags.FlagToStringValue(cmd, zoneIdFlag),
		RecordSetId:     recordSetId,
	}, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *dns.APIClient) dns.ApiGetRecordSetRequest {
	req := apiClient.GetRecordSet(ctx, model.ProjectId, model.ZoneId, model.RecordSetId)
	return req
}

func outputResult(cmd *cobra.Command, outputFormat string, recordSet *dns.RecordSet) error {
	switch outputFormat {
	case globalflags.PrettyOutputFormat:
		records := *recordSet.Records
		recordsData := []string{}
		for i := range records {
			recordsData = append(recordsData, *records[i].Content)
		}
		recordsDataJoin := strings.Join(recordsData, ",")

		table := tables.NewTable()
		table.AddRow("ID", *recordSet.Id)
		table.AddSeparator()
		table.AddRow("NAME", *recordSet.Name)
		table.AddSeparator()
		table.AddRow("STATE", *recordSet.State)
		table.AddSeparator()
		table.AddRow("TTL", *recordSet.Ttl)
		table.AddSeparator()
		table.AddRow("TYPE", *recordSet.Type)
		table.AddSeparator()
		table.AddRow("RECORDS DATA", recordsDataJoin)
		err := table.Display(cmd)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}

		return nil
	default:
		details, err := json.MarshalIndent(recordSet, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal DNS record set: %w", err)
		}
		cmd.Println(string(details))

		return nil
	}
}
