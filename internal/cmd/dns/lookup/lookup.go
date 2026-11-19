package lookup

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/types"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/dns/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
	dns "github.com/stackitcloud/stackit-sdk-go/services/dns/v1api"
)

const (
	domainArg = "DOMAIN"

	allResolversFlag       = "all-resolvers"
	includeRawResponseFlag = "include-raw-response"
)

var recordTypeFlag = flags.StringEnumFlag(
	"record-type",
	dns.AllowedLookupDomainPayloadRecordTypeEnumValues,
	"Record type to look up. If not set, defaults to A",
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	Domain             string
	RecordType         *dns.LookupDomainPayloadRecordType
	AllResolvers       bool
	IncludeRawResponse bool
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("lookup %s", domainArg),
		Short: "Looks up a domain on the STACKIT DNS resolvers",
		Long:  "Looks up a domain on the STACKIT DNS resolvers.",
		Args:  args.SingleArg(domainArg, nil),
		Example: examples.Build(
			examples.NewExample(
				`Look up domain "www.example.com"`,
				"$ stackit dns lookup www.example.com"),
			examples.NewExample(
				`Look up the AAAA records of domain "www.example.com" on all STACKIT DNS resolvers`,
				"$ stackit dns lookup www.example.com --record-type AAAA --all-resolvers"),
			examples.NewExample(
				`Look up domain "www.example.com" including the raw DNS response in JSON format`,
				"$ stackit dns lookup www.example.com --include-raw-response --output-format json"),
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
				return fmt.Errorf("look up DNS domain: %w", err)
			}

			return outputResult(params.Printer, model, resp)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	recordTypeFlag.Register(cmd.Flags())
	cmd.Flags().Bool(allResolversFlag, false, "If set, all STACKIT DNS resolvers are queried instead of a single one")
	cmd.Flags().Bool(includeRawResponseFlag, false, "If set, the raw DNS response is included in the output. Only shown in JSON and YAML output")
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	domain := inputArgs[0]
	globalFlags := globalflags.Parse(p, cmd)

	model := inputModel{
		GlobalFlagModel:    globalFlags,
		Domain:             domain,
		RecordType:         recordTypeFlag.Ptr(),
		AllResolvers:       flags.FlagToBoolValue(p, cmd, allResolversFlag),
		IncludeRawResponse: flags.FlagToBoolValue(p, cmd, includeRawResponseFlag),
	}

	if model.IncludeRawResponse && model.OutputFormat != print.JSONOutputFormat && model.OutputFormat != print.YAMLOutputFormat {
		p.Warn("--%s is only shown in JSON or YAML output, set --output-format json or --output-format yaml to see the raw DNS response\n", includeRawResponseFlag)
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *dns.APIClient) dns.ApiLookupDomainRequest {
	payload := dns.LookupDomainPayload{
		Address:    utils.Ptr(model.Domain),
		RecordType: model.RecordType,
	}
	if model.AllResolvers {
		payload.Type = utils.Ptr(dns.LOOKUPDOMAINPAYLOADTYPE_ALL)
	}
	if model.IncludeRawResponse {
		payload.IncludeRawResponse = utils.Ptr(true)
	}
	if model.Region != "" {
		payload.Region = utils.Ptr(model.Region)
	}

	req := apiClient.DefaultAPI.LookupDomain(ctx)
	req = req.LookupDomainPayload(payload)
	return req
}

func outputResult(p *print.Printer, model *inputModel, resp *dns.LookupDomainResponse) error {
	return p.OutputResult(model.OutputFormat, resp, func() error {
		if resp == nil {
			return fmt.Errorf("look up DNS domain response is empty")
		}
		table := tables.NewTable()
		table.SetHeader("RESOLVER IP", "NAME", "TYPE", "TTL", "DATA")
		rows := 0
		for _, result := range resp.State {
			// When querying all resolvers, also list the ones without records,
			// so it's visible which resolvers answered
			if model.AllResolvers && len(result.Records) == 0 {
				table.AddRow(utils.PtrString(result.ResolverIP), "", "", "", "No records found")
				rows++
				continue
			}
			for _, record := range result.Records {
				table.AddRow(
					utils.PtrString(result.ResolverIP),
					utils.PtrString(record.Name),
					utils.PtrString(record.Type),
					utils.PtrString(record.Ttl),
					utils.PtrString(record.Data),
				)
				rows++
			}
		}
		if rows == 0 {
			p.Outputf("No records found for domain %q\n", model.Domain)
			return nil
		}

		table.EnableAutoMergeOnColumns(1)
		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}
		return nil
	})
}
