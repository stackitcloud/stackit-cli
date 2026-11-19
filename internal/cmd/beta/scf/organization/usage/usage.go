package usage

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
	orgIdArg = "ORGANIZATION_ID"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	OrganizationId string
}

func NewCmd(p *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("usage %s", orgIdArg),
		Short: "Shows the usage summary of a STACKIT Cloud Foundry organization",
		Long:  "Shows the usage summary of a STACKIT Cloud Foundry organization.",
		Args:  args.SingleArg(orgIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Show the usage summary of a STACKIT Cloud Foundry organization with ID "xxx"`,
				`$ stackit beta scf organization usage xxx`),
			examples.NewExample(
				`Show the usage summary of a STACKIT Cloud Foundry organization with ID "xxx" in JSON format`,
				`$ stackit beta scf organization usage xxx --output-format json`),
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
				return fmt.Errorf("get STACKIT Cloud Foundry organization usage summary: %w", err)
			}

			return outputResult(p.Printer, model.OutputFormat, resp)
		},
	}
	return cmd
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	organizationId := inputArgs[0]

	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &cliErr.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		OrganizationId:  organizationId,
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient scf.DefaultAPI) scf.ApiGetOrganizationUsageSummaryRequest {
	return apiClient.GetOrganizationUsageSummary(ctx, model.ProjectId, model.Region, model.OrganizationId)
}

func outputResult(p *print.Printer, outputFormat string, summary *scf.OrganizationUsageSummary) error {
	return p.OutputResult(outputFormat, summary, func() error {
		if summary == nil {
			return fmt.Errorf("received nil response, could not display usage summary")
		}

		usage := summary.GetUsageSummary()

		table := tables.NewTable()
		table.SetHeader("Attribute", "Value")

		table.AddRow("Organization ID", summary.GetOrgId())
		table.AddRow("Platform ID", summary.GetPlatformId())
		table.AddRow("Project ID", summary.GetProjectId())
		table.AddRow("Region", summary.GetRegion())
		table.AddRow("Domains", usage.GetDomains())
		table.AddRow("Memory (MB)", usage.GetMemoryInMb())
		table.AddRow("Per-App Tasks", usage.GetPerAppTasks())
		table.AddRow("Reserved Ports", usage.GetReservedPorts())
		table.AddRow("Routes", usage.GetRoutes())
		table.AddRow("Service Instances", usage.GetServiceInstances())
		table.AddRow("Service Keys", usage.GetServiceKeys())
		table.AddRow("Started Instances", usage.GetStartedInstances())

		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}

		return nil
	})
}
