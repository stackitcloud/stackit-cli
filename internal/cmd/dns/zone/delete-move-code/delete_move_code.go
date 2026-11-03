package deletemovecode

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/types"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
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
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	ZoneId string
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("delete-move-code %s", zoneIdArg),
		Short: "Deletes the move code of a DNS zone",
		Long:  "Deletes (invalidates) the move code of a DNS zone, so it can no longer be used to move the zone to another project.",
		Args:  args.SingleArg(zoneIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Delete the move code of a DNS zone with ID "xxx"`,
				"$ stackit dns zone delete-move-code xxx"),
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
			prompt := fmt.Sprintf("Are you sure you want to delete the move code of zone %q? (This cannot be undone)", zoneLabel)
			err = params.Printer.PromptForConfirmation(prompt)
			if err != nil {
				return err
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("delete DNS zone move code: %w", err)
			}

			return outputResult(params.Printer, model, zoneLabel, resp)
		},
	}
	return cmd
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	zoneId := inputArgs[0]
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		ZoneId:          zoneId,
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *dns.APIClient) dns.ApiDeleteMoveCodeRequest {
	req := apiClient.DefaultAPI.DeleteMoveCode(ctx, model.ProjectId, model.ZoneId)
	return req
}

func outputResult(p *print.Printer, model *inputModel, zoneLabel string, resp *dns.Message) error {
	if resp == nil {
		return fmt.Errorf("move code deletion response is empty")
	}

	return p.OutputResult(model.OutputFormat, resp, func() error {
		p.Outputf("Deleted move code of zone %s\n", zoneLabel)
		return nil
	})
}
