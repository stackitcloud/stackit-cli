package describe

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/goccy/go-yaml"
	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
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

func NewCmd(params *params.CmdParams) *cobra.Command {
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
				`Get details of DNS record set with ID "xxx" in zone with ID "yyy" in JSON format`,
				"$ stackit dns record-set describe xxx --zone-id yyy --output-format json"),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			model, err := parseInput(params.Printer, cmd, args)
			if err != nil {
				return err
			}

			// Configure API client
			apiClient, err := client.ConfigureClient(params.Printer, params.CliVersion)
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

			return outputResult(params.Printer, model.OutputFormat, recordSet)
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

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	recordSetId := inputArgs[0]

	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		ZoneId:          flags.FlagToStringValue(p, cmd, zoneIdFlag),
		RecordSetId:     recordSetId,
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

func buildRequest(ctx context.Context, model *inputModel, apiClient *dns.APIClient) dns.ApiGetRecordSetRequest {
	req := apiClient.GetRecordSet(ctx, model.ProjectId, model.ZoneId, model.RecordSetId)
	return req
}

func outputResult(p *print.Printer, outputFormat string, recordSet *dns.RecordSet) error {
	if recordSet == nil {
		return fmt.Errorf("record set response is empty")
	}
	switch outputFormat {
	case print.JSONOutputFormat:
		details, err := json.MarshalIndent(recordSet, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal DNS record set: %w", err)
		}
		p.Outputln(string(details))

		return nil
	case print.YAMLOutputFormat:
		details, err := yaml.MarshalWithOptions(recordSet, yaml.IndentSequence(true), yaml.UseJSONMarshaler())
		if err != nil {
			return fmt.Errorf("marshal DNS record set: %w", err)
		}
		p.Outputln(string(details))

		return nil
	default:
		recordsData := make([]string, 0, len(*recordSet.Records))
		for _, r := range *recordSet.Records {
			recordsData = append(recordsData, *r.Content)
		}
		recordsDataJoin := strings.Join(recordsData, ", ")

		table := tables.NewTable()
		table.AddRow("ID", utils.PtrString(recordSet.Id))
		table.AddSeparator()
		table.AddRow("NAME", utils.PtrString(recordSet.Name))
		table.AddSeparator()
		table.AddRow("STATE", utils.PtrString(recordSet.State))
		table.AddSeparator()
		table.AddRow("TTL", utils.PtrString(recordSet.Ttl))
		table.AddSeparator()
		table.AddRow("TYPE", utils.PtrString(recordSet.Type))
		table.AddSeparator()
		table.AddRow("RECORDS DATA", recordsDataJoin)
		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}

		return nil
	}
}
