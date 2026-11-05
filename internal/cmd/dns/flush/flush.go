package flush

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

	"github.com/spf13/cobra"
	dns "github.com/stackitcloud/stackit-sdk-go/services/dns/v1api"
)

const (
	domainArg = "DOMAIN"
)

var typeFlag = flags.StringEnumFlag(
	"type",
	dns.AllowedFlushDomainPayloadTypeEnumValues,
	"Record type to flush. If not set, the cache is flushed for all record types of the domain",
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	Domain string
	Type   *dns.FlushDomainPayloadType
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("flush %s", domainArg),
		Short: "Flushes the cache for a domain on the STACKIT DNS resolvers",
		Long:  "Flushes the cache for a domain on the STACKIT DNS resolvers.",
		Args:  args.SingleArg(domainArg, nil),
		Example: examples.Build(
			examples.NewExample(
				`Flush the cache for domain "www.example.com"`,
				"$ stackit dns flush www.example.com"),
			examples.NewExample(
				`Flush the cache for the A records of domain "www.example.com"`,
				"$ stackit dns flush www.example.com --type A"),
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

			prompt := fmt.Sprintf("Are you sure you want to flush the cache for domain %q?", model.Domain)
			err = params.Printer.PromptForConfirmation(prompt)
			if err != nil {
				return err
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("flush DNS cache: %w", err)
			}

			return outputResult(params.Printer, model, resp)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	typeFlag.Register(cmd.Flags())
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	domain := inputArgs[0]
	globalFlags := globalflags.Parse(p, cmd)

	model := inputModel{
		GlobalFlagModel: globalFlags,
		Domain:          domain,
		Type:            typeFlag.Ptr(),
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *dns.APIClient) dns.ApiFlushDomainRequest {
	req := apiClient.DefaultAPI.FlushDomain(ctx)
	req = req.FlushDomainPayload(dns.FlushDomainPayload{
		Domain: model.Domain,
		Type:   model.Type,
	})
	return req
}

func outputResult(p *print.Printer, model *inputModel, resp *dns.FlushDomainResponse) error {
	return p.OutputResult(model.OutputFormat, resp, func() error {
		if resp == nil {
			return fmt.Errorf("flush DNS cache response is empty")
		}
		if model.Type != nil {
			p.Outputf("Triggered flush of the %s record cache for domain %q\n", *model.Type, model.Domain)
		} else {
			p.Outputf("Triggered flush of the cache for domain %q\n", model.Domain)
		}
		return nil
	})
}
