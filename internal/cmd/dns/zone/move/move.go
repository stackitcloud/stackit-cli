package move

import (
	"context"
	"fmt"
	"strings"

	"github.com/stackitcloud/stackit-cli/internal/pkg/types"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/dns/client"
	dnsUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/dns/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
	dns "github.com/stackitcloud/stackit-sdk-go/services/dns/v1api"
)

const (
	zoneDomainNameFlag  = "zone-domain-name"
	targetProjectIdFlag = "target-project-id"
	codeFlag            = "code"

	// The generated move code by this command should not be valid for 30 days.
	generatedMoveCodeTtl int32 = 3600 // 1 hour
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	ZoneDomainName  string
	TargetProjectId string
	Code            *string
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "move",
		Short: "Moves a DNS zone to another project",
		Long: "Moves a DNS zone to another project.\n" +
			"The API call this command makes is always performed against the destination project (--target-project-id). Moving requires an authorizing move code for the zone, which can be obtained in one of two ways:\n" +
			"  - Generated automatically by this command: omit --code. The configured project is then used to look up the zone and must be its current project. The generated move code is valid for 1 hour.\n" +
			"  - Generated beforehand, by yourself or by whoever controls the zone's current project, via \"stackit dns zone generate-move-code\" (and optionally checked with \"stackit dns zone validate-move-code\"), then passed explicitly via --code." +
			" The configured project isn't used for the move itself in this case.",
		Args: args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Move a DNS zone with domain name "www.my-zone.com" to project "xxx", automatically generating a move code against the zone's current project`,
				"$ stackit dns zone move --zone-domain-name www.my-zone.com --target-project-id xxx"),
			examples.NewExample(
				`Move a DNS zone with domain name "www.my-zone.com" from project "yyy" to project "xxx", automatically generating a move code. --project-id explicitly sets the source project, i.e. the zone's current project`,
				"$ stackit dns zone move --project-id yyy --zone-domain-name www.my-zone.com --target-project-id xxx"),
			examples.NewExample(
				`Move a DNS zone with domain name "www.my-zone.com" to project "xxx" using a move code generated beforehand (e.g. via "stackit dns zone generate-move-code")`,
				"$ stackit dns zone move --zone-domain-name www.my-zone.com --target-project-id xxx --code @./move-code.txt"),
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

			prompt := fmt.Sprintf("Are you sure you want to move the zone with domain name %q to project %q?", model.ZoneDomainName, model.TargetProjectId)
			err = params.Printer.PromptForConfirmation(prompt)
			if err != nil {
				return err
			}

			code := model.Code
			if code == nil {
				generatedCode, err := generateMoveCode(ctx, params.Printer, apiClient, model)
				if err != nil {
					return err
				}
				code = &generatedCode
			}

			// Call API
			req := buildRequest(ctx, model, *code, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("move DNS zone: %w", err)
			}

			return outputResult(params.Printer, model, resp)
		},
	}
	configureFlags(cmd, params)
	return cmd
}

func generateMoveCode(ctx context.Context, p *print.Printer, apiClient *dns.APIClient, model *inputModel) (string, error) {
	zoneId, err := dnsUtils.GetZoneIdByDnsName(ctx, apiClient.DefaultAPI, model.ProjectId, model.ZoneDomainName)
	if err != nil {
		return "", fmt.Errorf("find DNS zone by domain name: %w", err)
	}

	createResp, err := apiClient.DefaultAPI.CreateMoveCode(ctx, model.ProjectId, zoneId).CreateMoveCodePayload(dns.CreateMoveCodePayload{
		Ttl: utils.Ptr(generatedMoveCodeTtl),
	}).Execute()
	if err != nil {
		return "", fmt.Errorf("generate DNS zone move code: %w", err)
	}

	p.Debug(print.DebugLevel, "generated move code for zone %q", zoneId)
	return createResp.Code, nil
}

func configureFlags(cmd *cobra.Command, params *types.CmdParams) {
	cmd.Flags().String(zoneDomainNameFlag, "", "Domain name of the zone to move")
	cmd.Flags().Var(flags.UUIDFlag(), targetProjectIdFlag, "ID of the project the zone should be moved to")

	code := flags.SecretFlag(codeFlag, params)
	cmd.Flags().Var(code, codeFlag, fmt.Sprintf("%s If not set, a move code is generated automatically for the zone in the configured project, which must then be the zone's current project.", code.Usage()))

	err := flags.MarkFlagsRequired(cmd, zoneDomainNameFlag, targetProjectIdFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command, _ []string) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	// --code is optional: if omitted, ResolveSecretFlag returns nil without
	// prompting, and a code is generated automatically instead (see
	// generateMoveCode).
	code, err := flags.ResolveSecretFlag(p, cmd, codeFlag)
	if err != nil {
		return nil, err
	}
	if code != nil {
		// Move codes never contain surrounding whitespace, but files written
		// by editors usually end with a trailing newline.
		trimmed := strings.TrimSpace(*code)
		code = &trimmed
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		ZoneDomainName:  flags.FlagToStringValue(p, cmd, zoneDomainNameFlag),
		TargetProjectId: flags.FlagToStringValue(p, cmd, targetProjectIdFlag),
		Code:            code,
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, code string, apiClient *dns.APIClient) dns.ApiMoveZoneRequest {
	req := apiClient.DefaultAPI.MoveZone(ctx, model.TargetProjectId)
	req = req.MoveZonePayload(dns.MoveZonePayload{
		Code:        code,
		ZoneDnsName: model.ZoneDomainName,
	})
	return req
}

func outputResult(p *print.Printer, model *inputModel, resp *dns.Message) error {
	if resp == nil {
		return fmt.Errorf("move DNS zone response is empty")
	}

	return p.OutputResult(model.OutputFormat, resp, func() error {
		p.Outputf("Moved zone with domain name %q to project %q\n", model.ZoneDomainName, model.TargetProjectId)
		return nil
	})
}
