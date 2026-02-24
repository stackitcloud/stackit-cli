package describe

import (
	"context"

	"github.com/stackitcloud/stackit-cli/internal/pkg/types"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/resourcemanager/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/resourcemanager"
)

const (
	organizationIdArg = "ORGANIZATION_ID"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	OrganizationId string
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "describe",
		Short: "Show a organization",
		Long:  "Show a organization.",
		// the arg can be the organization uuid or the container id, which is not a uuid, so no validation needed
		Args: args.SingleArg(organizationIdArg, nil),
		Example: examples.Build(
			examples.NewExample(
				`Describe the organization with the organization uuid "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"`,
				"$ stackit organization describe xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
			),
			examples.NewExample(
				`Describe the organization with the container id "foo-bar-organization"`,
				"$ stackit organization describe foo-bar-organization",
			),
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
				return err
			}

			return outputResult(params.Printer, model.OutputFormat, resp)
		},
	}
	return cmd
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	organizationId := inputArgs[0]
	globalFlags := globalflags.Parse(p, cmd)

	model := inputModel{
		GlobalFlagModel: globalFlags,
		OrganizationId:  organizationId,
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *resourcemanager.APIClient) resourcemanager.ApiGetOrganizationRequest {
	req := apiClient.GetOrganization(ctx, model.OrganizationId)
	return req
}

func outputResult(p *print.Printer, outputFormat string, organization *resourcemanager.OrganizationResponse) error {
	return p.OutputResult(outputFormat, organization, func() error {
		if organization == nil {
			p.Outputln("show organization: empty response")
			return nil
		}

		table := tables.NewTable()

		table.AddRow("ORGANIZATION ID", utils.PtrString(organization.OrganizationId))
		table.AddSeparator()
		table.AddRow("NAME", utils.PtrString(organization.Name))
		table.AddSeparator()
		table.AddRow("CONTAINER ID", utils.PtrString(organization.ContainerId))
		table.AddSeparator()
		table.AddRow("STATUS", utils.PtrString(organization.LifecycleState))
		table.AddSeparator()
		table.AddRow("CREATION TIME", utils.PtrString(organization.CreationTime))
		table.AddSeparator()
		table.AddRow("UPDATE TIME", utils.PtrString(organization.UpdateTime))
		table.AddSeparator()
		table.AddRow("LABELS", utils.JoinStringMap(utils.PtrValue(organization.Labels), ": ", ", "))

		p.Outputln(table.Render())
		return nil
	})
}
