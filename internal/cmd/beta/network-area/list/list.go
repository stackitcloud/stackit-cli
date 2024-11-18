package list

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/goccy/go-yaml"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/iaas/client"
	rmClient "github.com/stackitcloud/stackit-cli/internal/pkg/services/resourcemanager/client"
	rmUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/resourcemanager/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-sdk-go/services/iaas"

	"github.com/spf13/cobra"
)

const (
	limitFlag          = "limit"
	organizationIdFlag = "organization-id"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	Limit          *int64
	OrganizationId *string
}

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Lists all STACKIT Network Areas (SNA) of an organization",
		Long:  "Lists all STACKIT Network Areas (SNA) of an organization.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Lists all network areas of organization "xxx"`,
				"$ stackit beta network-area list --organization-id xxx",
			),
			examples.NewExample(
				`Lists all network areas of organization "xxx" in JSON format`,
				"$ stackit beta network-area list --organization-id xxx --output-format json",
			),
			examples.NewExample(
				`Lists up to 10 network areas of organization "xxx"`,
				"$ stackit beta network-area list --organization-id xxx --limit 10",
			),
		),
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := context.Background()
			model, err := parseInput(p, cmd)
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
				return fmt.Errorf("list network areas: %w", err)
			}

			if resp.Items == nil || len(*resp.Items) == 0 {
				var orgLabel string
				rmApiClient, err := rmClient.ConfigureClient(p)
				if err == nil {
					orgLabel, err = rmUtils.GetOrganizationName(ctx, rmApiClient, *model.OrganizationId)
					if err != nil {
						p.Debug(print.ErrorLevel, "get organization name: %v", err)
						orgLabel = *model.OrganizationId
					}
				} else {
					p.Debug(print.ErrorLevel, "configure resource manager client: %v", err)
				}
				p.Info("No STACKIT Network Areas found for organization %q\n", orgLabel)
				return nil
			}

			// Truncate output
			items := *resp.Items
			if model.Limit != nil && len(items) > int(*model.Limit) {
				items = items[:*model.Limit]
			}

			return outputResult(p, model.OutputFormat, items)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Int64(limitFlag, 0, "Maximum number of entries to list")
	cmd.Flags().Var(flags.UUIDFlag(), organizationIdFlag, "Organization ID")

	err := flags.MarkFlagsRequired(cmd, organizationIdFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	limit := flags.FlagToInt64Pointer(p, cmd, limitFlag)
	if limit != nil && *limit < 1 {
		return nil, &errors.FlagValidationError{
			Flag:    limitFlag,
			Details: "must be greater than 0",
		}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		Limit:           limit,
		OrganizationId:  flags.FlagToStringPointer(p, cmd, organizationIdFlag),
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

func buildRequest(ctx context.Context, model *inputModel, apiClient *iaas.APIClient) iaas.ApiListNetworkAreasRequest {
	return apiClient.ListNetworkAreas(ctx, *model.OrganizationId)
}

func outputResult(p *print.Printer, outputFormat string, networkAreas []iaas.NetworkArea) error {
	switch outputFormat {
	case print.JSONOutputFormat:
		details, err := json.MarshalIndent(networkAreas, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal network area: %w", err)
		}
		p.Outputln(string(details))

		return nil
	case print.YAMLOutputFormat:
		details, err := yaml.MarshalWithOptions(networkAreas, yaml.IndentSequence(true))
		if err != nil {
			return fmt.Errorf("marshal area: %w", err)
		}
		p.Outputln(string(details))

		return nil
	default:
		table := tables.NewTable()
		table.SetHeader("ID", "Name", "Status", "Network Ranges", "# Attached Projects")

		for _, networkArea := range networkAreas {
			table.AddRow(*networkArea.AreaId, *networkArea.Name, *networkArea.State, len(*networkArea.Ipv4.NetworkRanges), *networkArea.ProjectCount)
			table.AddSeparator()
		}

		p.Outputln(table.Render())
		return nil
	}
}
