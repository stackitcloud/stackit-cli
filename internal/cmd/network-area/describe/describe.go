package describe

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/stackitcloud/stackit-cli/internal/pkg/types"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/iaas/client"
	iaasUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/iaas/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/iaas"

	"github.com/spf13/cobra"
)

const (
	areaIdArg                = "AREA_ID"
	organizationIdFlag       = "organization-id"
	showAttachedProjectsFlag = "show-attached-projects"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	OrganizationId       *string
	AreaId               string
	ShowAttachedProjects bool
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("describe %s", areaIdArg),
		Short: "Shows details of a STACKIT Network Area",
		Long:  "Shows details of a STACKIT Network Area in an organization.",
		Args:  args.SingleArg(areaIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Show details of a network area with ID "xxx" in organization with ID "yyy"`,
				"$ stackit network-area describe xxx --organization-id yyy",
			),
			examples.NewExample(
				`Show details of a network area with ID "xxx" in organization with ID "yyy" and show attached projects`,
				"$ stackit network-area describe xxx --organization-id yyy --show-attached-projects",
			),
			examples.NewExample(
				`Show details of a network area with ID "xxx" in organization with ID "yyy" in JSON format`,
				"$ stackit network-area describe xxx --organization-id yyy --output-format json",
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
				return fmt.Errorf("read network area: %w", err)
			}

			var projects []string

			if model.ShowAttachedProjects {
				projects, err = iaasUtils.ListAttachedProjects(ctx, apiClient, *model.OrganizationId, model.AreaId)
				if err != nil && errors.Is(err, iaasUtils.ErrItemsNil) {
					projects = []string{}
				} else if err != nil {
					return fmt.Errorf("get attached projects: %w", err)
				}
			}

			return outputResult(params.Printer, model.OutputFormat, resp, projects)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), organizationIdFlag, "Organization ID")
	cmd.Flags().Bool(showAttachedProjectsFlag, false, "Whether to show attached projects. If a network area has several attached projects, their retrieval may take some time and the output may be extensive.")

	err := flags.MarkFlagsRequired(cmd, organizationIdFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	areaId := inputArgs[0]

	globalFlags := globalflags.Parse(p, cmd)

	model := inputModel{
		GlobalFlagModel:      globalFlags,
		OrganizationId:       flags.FlagToStringPointer(p, cmd, organizationIdFlag),
		AreaId:               areaId,
		ShowAttachedProjects: flags.FlagToBoolValue(p, cmd, showAttachedProjectsFlag),
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *iaas.APIClient) iaas.ApiGetNetworkAreaRequest {
	return apiClient.GetNetworkArea(ctx, *model.OrganizationId, model.AreaId)
}

func outputResult(p *print.Printer, outputFormat string, networkArea *iaas.NetworkArea, attachedProjects []string) error {
	if networkArea == nil {
		return fmt.Errorf("network area is nil")
	}

	return p.OutputResult(outputFormat, networkArea, func() error {
		table := tables.NewTable()
		table.AddRow("ID", utils.PtrString(networkArea.Id))
		table.AddSeparator()
		table.AddRow("NAME", utils.PtrString(networkArea.Name))
		table.AddSeparator()
		if networkArea.Labels != nil && len(*networkArea.Labels) > 0 {
			var labels []string
			for key, value := range *networkArea.Labels {
				labels = append(labels, fmt.Sprintf("%s: %s", key, value))
			}
			table.AddRow("LABELS", strings.Join(labels, "\n"))
			table.AddSeparator()
		}
		if len(attachedProjects) > 0 {
			table.AddRow("ATTACHED PROJECTS IDS", strings.Join(attachedProjects, "\n"))
			table.AddSeparator()
		} else {
			table.AddRow("# ATTACHED PROJECTS", utils.PtrString(networkArea.ProjectCount))
			table.AddSeparator()
		}
		table.AddRow("CREATED AT", utils.PtrString(networkArea.CreatedAt))
		table.AddSeparator()
		table.AddRow("UPDATED AT", utils.PtrString(networkArea.UpdatedAt))
		table.AddSeparator()

		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}
		return nil
	})
}
