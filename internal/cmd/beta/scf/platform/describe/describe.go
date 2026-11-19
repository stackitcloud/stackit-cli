package describe

import (
	"context"
	"fmt"

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
		Use:   fmt.Sprintf("describe %s", platformIdArg),
		Short: "Shows details of a STACKIT Cloud Foundry platform",
		Long:  "Shows details of a STACKIT Cloud Foundry platform.",
		Args:  args.SingleArg(platformIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Get details of a STACKIT Cloud Foundry platform with ID "xxx"`,
				`$ stackit beta scf platform describe xxx`),
			examples.NewExample(
				`Get details of a STACKIT Cloud Foundry platform with ID "xxx" in JSON format`,
				`$ stackit beta scf platform describe xxx --output-format json`),
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
				return fmt.Errorf("get STACKIT Cloud Foundry platform: %w", err)
			}

			return outputResult(p.Printer, model.OutputFormat, resp)
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

func buildRequest(ctx context.Context, model *inputModel, apiClient scf.DefaultAPI) scf.ApiGetPlatformRequest {
	return apiClient.GetPlatform(ctx, model.ProjectId, model.Region, model.PlatformId)
}

func outputResult(p *print.Printer, outputFormat string, platform *scf.Platforms) error {
	return p.OutputResult(outputFormat, platform, func() error {
		if platform == nil {
			return fmt.Errorf("received nil response, could not display details")
		}

		table := tables.NewTable()
		table.SetHeader("Attribute", "Value")

		table.AddRow("ID", platform.GetGuid())
		table.AddRow("Name", platform.GetDisplayName())
		table.AddRow("System ID", platform.GetSystemId())
		table.AddRow("Region", platform.GetRegion())
		table.AddRow("API URL", platform.GetApiUrl())
		table.AddRow("Console URL", utils.PtrString(platform.ConsoleUrl))

		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}

		return nil
	})
}
