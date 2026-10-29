package generatemovecode

import (
	"context"
	"fmt"

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
	zoneIdArg = "ZONE_ID"

	ttlFlag = "ttl"

	ttlMin = 1
	ttlMax = 2592000 // 30 days
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	ZoneId string
	Ttl    *int32
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("generate-move-code %s", zoneIdArg),
		Short: "Generates a move code for a DNS zone",
		Long:  "Generates a move code for a DNS zone, which is required to move the zone to another project.\nThe code is only shown once upon creation, so make sure to save it. By default, the code is valid for 30 days.",
		Args:  args.SingleArg(zoneIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Generate a move code for a DNS zone with ID "xxx"`,
				"$ stackit dns zone generate-move-code xxx"),
			examples.NewExample(
				`Generate a move code for a DNS zone with ID "xxx" that is valid for one hour`,
				"$ stackit dns zone generate-move-code xxx --ttl 3600"),
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

			zoneLabel, err := dnsUtils.GetZoneName(ctx, apiClient.DefaultAPI, model.ProjectId, model.ZoneId)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get zone name: %v", err)
				zoneLabel = model.ZoneId
			}

			prompt := fmt.Sprintf("Are you sure you want to generate a move code for zone %q?", zoneLabel)
			err = params.Printer.PromptForConfirmation(prompt)
			if err != nil {
				return err
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("generate DNS zone move code: %w", err)
			}

			return outputResult(params.Printer, model.OutputFormat, resp)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Int32(ttlFlag, 0, "Validity duration of the move code in seconds (min 1, max 2592000 / 30 days). If not set, the API default is used")
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	zoneId := inputArgs[0]

	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	ttl := flags.FlagToInt32Pointer(p, cmd, ttlFlag)
	if ttl != nil && (*ttl < ttlMin || *ttl > ttlMax) {
		return nil, &errors.FlagValidationError{
			Flag:    ttlFlag,
			Details: fmt.Sprintf("must be between %d and %d", ttlMin, ttlMax),
		}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		ZoneId:          zoneId,
		Ttl:             ttl,
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *dns.APIClient) dns.ApiCreateMoveCodeRequest {
	req := apiClient.DefaultAPI.CreateMoveCode(ctx, model.ProjectId, model.ZoneId)
	if model.Ttl != nil {
		req = req.CreateMoveCodePayload(dns.CreateMoveCodePayload{
			Ttl: model.Ttl,
		})
	}
	return req
}

func outputResult(p *print.Printer, outputFormat string, resp *dns.MoveCodeResponse) error {
	if resp == nil {
		return fmt.Errorf("move code response is empty")
	}

	return p.OutputResult(outputFormat, resp, func() error {
		p.Outputf("Generated move code. It is only shown once, so make sure to save it.\n")
		p.Outputf("Code: %s\n", resp.Code)
		p.Outputf("Expires at: %s\n", resp.ExpiresAt)
		return nil
	})
}
