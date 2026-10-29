package list

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"
	scf "github.com/stackitcloud/stackit-sdk-go/services/scf/v1api"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	cliErr "github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/scf/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

const (
	platformIdArg = "PLATFORM_ID"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	PlatformId string
}

func NewCmd(p *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("list %s", platformIdArg),
		Short: "Lists the quotas of a STACKIT Cloud Foundry platform",
		Long:  "Lists the quotas of a STACKIT Cloud Foundry platform.",
		Args:  args.SingleArg(platformIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`List the quotas of a STACKIT Cloud Foundry platform with ID "xxx"`,
				`$ stackit beta scf platform quota list xxx`),
			examples.NewExample(
				`List the quotas of a STACKIT Cloud Foundry platform with ID "xxx" in JSON format`,
				`$ stackit beta scf platform quota list xxx --output-format json`),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			model, err := parseInput(p.Printer, cmd, args)
			if err != nil {
				return err
			}

			// Configure API client
			apiClient, err := client.ConfigureClient(p.Printer, p.CliVersion)
			if err != nil {
				return err
			}

			// Call API
			req := buildRequest(ctx, model, apiClient.DefaultAPI)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("list STACKIT Cloud Foundry platform quotas: %w", err)
			}

			return outputResult(p.Printer, model, resp)
		},
	}
	return cmd
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	platformId := inputArgs[0]

	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &cliErr.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		PlatformId:      platformId,
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient scf.DefaultAPI) scf.ApiListPlatformQuotasRequest {
	return apiClient.ListPlatformQuotas(ctx, model.ProjectId, model.Region, model.PlatformId)
}

func outputResult(p *print.Printer, model *inputModel, quotas []scf.Quota) error {
	return p.OutputResult(model.OutputFormat, quotas, func() error {
		if len(quotas) == 0 {
			p.Outputf("No quotas found for STACKIT Cloud Foundry platform %s\n", model.PlatformId)
			return nil
		}

		table := tables.NewTable()

		table.SetTitle("STACKIT Cloud Foundry platform quotas")
		table.SetHeader("ID", "NAME", "PLATFORM ID", "CREATED")
		table.AddSeparator()
		for _, quota := range quotas {
			table.AddRow(
				quota.GetGuid(),
				quota.GetName(),
				quota.GetPlatformId(),
				quota.CreatedAt.Format(time.DateTime),
			)
			table.AddSeparator()
		}

		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}
		return nil
	})
}
