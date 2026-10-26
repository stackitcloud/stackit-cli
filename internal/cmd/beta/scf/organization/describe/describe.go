package describe

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
	orgIdArg = "ORGANIZATION_ID"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	OrganizationId string
}

func NewCmd(p *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("describe %s", orgIdArg),
		Short: "Shows details of a STACKIT Cloud Foundry organization",
		Long:  "Shows details of a STACKIT Cloud Foundry organization.",
		Args:  args.SingleArg(orgIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Get details of a STACKIT Cloud Foundry organization with ID "xxx"`,
				`$ stackit beta scf organization describe xxx`),
			examples.NewExample(
				`Get details of a STACKIT Cloud Foundry organization with ID "xxx" in JSON format`,
				`$ stackit beta scf organization describe xxx --output-format json`),
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
				return fmt.Errorf("get STACKIT Cloud Foundry organization: %w", err)
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

func buildRequest(ctx context.Context, model *inputModel, apiClient scf.DefaultAPI) scf.ApiGetOrganizationRequest {
	return apiClient.GetOrganization(ctx, model.ProjectId, model.Region, model.OrganizationId)
}

func outputResult(p *print.Printer, outputFormat string, organization *scf.Organization) error {
	return p.OutputResult(outputFormat, organization, func() error {
		if organization == nil {
			return fmt.Errorf("received nil response, could not display details")
		}

		table := tables.NewTable()
		table.SetHeader("Attribute", "Value")

		table.AddRow("ID", organization.GetGuid())
		table.AddRow("Name", organization.GetName())
		table.AddRow("Platform ID", organization.GetPlatformId())
		table.AddRow("Project ID", organization.GetProjectId())
		table.AddRow("Quota ID", organization.GetQuotaId())
		table.AddRow("Region", organization.GetRegion())
		table.AddRow("Status", organization.GetStatus())
		table.AddRow("Suspended", organization.GetSuspended())
		table.AddRow("Created", organization.CreatedAt.Format(time.DateTime))
		table.AddRow("Updated", organization.UpdatedAt.Format(time.DateTime))

		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}

		return nil
	})
}
