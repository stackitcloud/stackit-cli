package describe

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/goccy/go-yaml"
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

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "describe",
		Short: "Shows details of a network area",
		Long:  "Shows details of a network area in an organization.",
		Args:  args.SingleArg(areaIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Show details of a network area with ID "xxx" in organization with ID "yyy"`,
				"$ stackit beta network-area describe xxx --organization-id yyy",
			),
			examples.NewExample(
				`Show details of a network area with ID "xxx" in organization with ID "yyy" and show attached projects`,
				"$ stackit beta network-area describe xxx --organization-id yyy --show-attached-projects",
			),
			examples.NewExample(
				`Show details of a network area with ID "xxx" in organization with ID "yyy" in JSON format`,
				"$ stackit beta network-area describe xxx --organization-id yyy --output-format json",
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			model, err := parseInput(p, cmd, args)
			if err != nil {
				return err
			}

			// Configure API client
			apiClient, err := client.ConfigureClient(p)
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
				if err != nil {
					return fmt.Errorf("get attached projects: %w", err)
				}
			}

			return outputResult(p, model.OutputFormat, resp, projects)
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

	if p.IsVerbosityDebug() {
		modelStr, err := print.BuildDebugStrFromInputModel(model)
		if err != nil {
			p.Debug(print.ErrorLevel, "convert model to string for debugging: %v", err)
		} else {
			p.Debug(print.DebugLevel, "parsed input values: %s", modelStr)
		}
	}

	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *iaas.APIClient) iaas.ApiGetNetworkAreaRequest {
	return apiClient.GetNetworkArea(ctx, *model.OrganizationId, model.AreaId)
}

func outputResult(p *print.Printer, outputFormat string, networkArea *iaas.NetworkArea, attachedProjects []string) error {
	switch outputFormat {
	case print.JSONOutputFormat:
		details, err := json.MarshalIndent(networkArea, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal network area: %w", err)
		}
		p.Outputln(string(details))

		return nil
	case print.YAMLOutputFormat:
		details, err := yaml.MarshalWithOptions(networkArea, yaml.IndentSequence(true))
		if err != nil {
			return fmt.Errorf("marshal network area: %w", err)
		}
		p.Outputln(string(details))

		return nil
	default:
		var routes []string
		if networkArea.Ipv4.Routes != nil {
			for _, route := range *networkArea.Ipv4.Routes {
				routes = append(routes, fmt.Sprintf("next hop: %s\nprefix: %s", *route.Nexthop, *route.Prefix))
			}
		}

		var networkRanges []string
		if networkArea.Ipv4.NetworkRanges != nil {
			for _, networkRange := range *networkArea.Ipv4.NetworkRanges {
				networkRanges = append(networkRanges, *networkRange.Prefix)
			}
		}

		table := tables.NewTable()
		table.AddRow("ID", *networkArea.AreaId)
		table.AddSeparator()
		table.AddRow("NAME", *networkArea.Name)
		table.AddSeparator()
		table.AddRow("STATE", *networkArea.State)
		table.AddSeparator()
		if len(networkRanges) > 0 {
			table.AddRow("NETWORK RANGES", strings.Join(networkRanges, ","))
		}
		table.AddSeparator()
		for i, route := range routes {
			table.AddRow(fmt.Sprintf("STATIC ROUTE %d", i+1), route)
			table.AddSeparator()
		}
		if networkArea.Ipv4.TransferNetwork != nil {
			table.AddRow("TRANSFER RANGE", *networkArea.Ipv4.TransferNetwork)
			table.AddSeparator()
		}
		if networkArea.Ipv4.DefaultNameservers != nil {
			table.AddRow("DNS NAME SERVERS", strings.Join(*networkArea.Ipv4.DefaultNameservers, ","))
			table.AddSeparator()
		}
		if networkArea.Ipv4.DefaultPrefixLen != nil {
			table.AddRow("DEFAULT PREFIX LENGTH", *networkArea.Ipv4.DefaultPrefixLen)
			table.AddSeparator()
		}
		if networkArea.Ipv4.MaxPrefixLen != nil {
			table.AddRow("MAX PREFIX LENGTH", *networkArea.Ipv4.MaxPrefixLen)
			table.AddSeparator()
		}
		if networkArea.Ipv4.MinPrefixLen != nil {
			table.AddRow("MIN PREFIX LENGTH", *networkArea.Ipv4.MinPrefixLen)
			table.AddSeparator()
		}
		if len(attachedProjects) > 0 {
			table.AddRow("ATTACHED PROJECTS IDS", strings.Join(attachedProjects, "\n"))
			table.AddSeparator()
		} else {
			table.AddRow("# ATTACHED PROJECTS", *networkArea.ProjectCount)
			table.AddSeparator()
		}

		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}
		return nil
	}
}
