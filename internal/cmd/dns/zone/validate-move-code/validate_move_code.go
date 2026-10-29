package validatemovecode

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
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
	dns "github.com/stackitcloud/stackit-sdk-go/services/dns/v1api"
)

const (
	zoneIdArg = "ZONE_ID"

	codeFlag = "code"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	ZoneId string
	Code   string
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("validate-move-code %s", zoneIdArg),
		Short: "Validates the move code of a DNS zone",
		Long:  "Validates the move code of a DNS zone. Checks that the provided code is valid, not expired, and belongs to the given zone.",
		Args:  args.SingleArg(zoneIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Validate the move code stored in "move-code.txt" of a DNS zone with ID "xxx"`,
				"$ stackit dns zone validate-move-code xxx --code @./move-code.txt"),
			examples.NewExample(
				`Validate a move code read from stdin for a DNS zone with ID "xxx"`,
				"$ echo yyy | stackit dns zone validate-move-code xxx --code -"),
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
				return fmt.Errorf("validate DNS zone move code: %w", err)
			}

			return outputResult(params.Printer, model.OutputFormat, resp)
		},
	}
	configureFlags(cmd, params)
	return cmd
}

func configureFlags(cmd *cobra.Command, params *types.CmdParams) {
	code := flags.SecretFlag(codeFlag, params)
	cmd.Flags().Var(code, codeFlag, code.Usage())

	err := flags.MarkFlagsRequired(cmd, codeFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	zoneId := inputArgs[0]

	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	code, err := flags.ResolveSecretFlag(p, cmd, codeFlag)
	if err != nil {
		return nil, err
	}
	if code == nil {
		return nil, fmt.Errorf("flag --%s is required", codeFlag)
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		ZoneId:          zoneId,
		// Move codes never contain surrounding whitespace, but files written
		// by editors usually end with a trailing newline.
		Code: strings.TrimSpace(*code),
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *dns.APIClient) dns.ApiValidateMoveCodeRequest {
	req := apiClient.DefaultAPI.ValidateMoveCode(ctx, model.ProjectId, model.ZoneId)
	req = req.ValidateMoveCodePayload(dns.ValidateMoveCodePayload{
		Code: model.Code,
	})
	return req
}

func outputResult(p *print.Printer, outputFormat string, resp *dns.Message) error {
	if resp == nil {
		return fmt.Errorf("move code validation response is empty")
	}

	return p.OutputResult(outputFormat, resp, func() error {
		p.Outputf("The move code is valid.\n")
		return nil
	})
}
