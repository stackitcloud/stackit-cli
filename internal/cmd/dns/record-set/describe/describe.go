package describe

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
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

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

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("describe %s", recordSetIdArg),
		Short: "Shows details  of a DNS record set",
		Long:  "Shows details  of a DNS record set.",
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
			apiClient, err := client.ConfigureClient(p)
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

			return outputResult(p, model.OutputFormat, recordSet)
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

func outputResult(p *print.Printer, outputFormat string, recordSet *dns.RecordSet) error {
	switch outputFormat {
	case globalflags.PrettyOutputFormat:
		recordsData := make([]string, 0, len(*recordSet.Records))
		for _, r := range *recordSet.Records {
			recordsData = append(recordsData, *r.Content)
		}
		recordsDataJoin := strings.Join(recordsData, ", ")

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
		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}

		return nil
	default:
		details, err := json.MarshalIndent(recordSet, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal DNS record set: %w", err)
		}
		p.Outputln(string(details))

		return nil
	}
}
