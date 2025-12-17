package list

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/types"

	rmClient "github.com/stackitcloud/stackit-cli/internal/pkg/services/resourcemanager/client"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/iaas/client"
	rmUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/resourcemanager/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/iaas"
)

const (
	limitFlag          = "limit"
	organizationIdFlag = "organization-id"
	labelSelectorFlag  = "label-selector"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	Limit          *int64
	OrganizationId *string
	LabelSelector  *string
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Lists all STACKIT Network Areas (SNA) of an organization",
		Long:  "Lists all STACKIT Network Areas (SNA) of an organization.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Lists all network areas of organization "xxx"`,
				"$ stackit network-area list --organization-id xxx",
			),
			examples.NewExample(
				`Lists all network areas of organization "xxx" in JSON format`,
				"$ stackit network-area list --organization-id xxx --output-format json",
			),
			examples.NewExample(
				`Lists up to 10 network areas of organization "xxx"`,
				"$ stackit network-area list --organization-id xxx --limit 10",
			),
			examples.NewExample(
				`Lists all network areas of organization "xxx" which contains the label yyy`,
				"$ stackit network-area list --organization-id xxx --label-selector yyy",
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
				return fmt.Errorf("list network areas: %w", err)
			}

			if resp.Items == nil || len(*resp.Items) == 0 {
				var orgLabel string
				rmApiClient, err := rmClient.ConfigureClient(params.Printer, params.CliVersion)
				if err == nil {
					orgLabel, err = rmUtils.GetOrganizationName(ctx, rmApiClient, *model.OrganizationId)
					if err != nil {
						params.Printer.Debug(print.ErrorLevel, "get organization name: %v", err)
						orgLabel = *model.OrganizationId
					} else if orgLabel == "" {
						orgLabel = *model.OrganizationId
					}
				} else {
					params.Printer.Debug(print.ErrorLevel, "configure resource manager client: %v", err)
				}
				params.Printer.Info("No STACKIT Network Areas found for organization %q\n", orgLabel)
				return nil
			}

			// Truncate output
			items := *resp.Items
			if model.Limit != nil && len(items) > int(*model.Limit) {
				items = items[:*model.Limit]
			}

			return outputResult(params.Printer, model.OutputFormat, items)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Int64(limitFlag, 0, "Maximum number of entries to list")
	cmd.Flags().Var(flags.UUIDFlag(), organizationIdFlag, "Organization ID")
	cmd.Flags().String(labelSelectorFlag, "", "Filter by label")

	err := flags.MarkFlagsRequired(cmd, organizationIdFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command, _ []string) (*inputModel, error) {
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
		LabelSelector:   flags.FlagToStringPointer(p, cmd, labelSelectorFlag),
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *iaas.APIClient) iaas.ApiListNetworkAreasRequest {
	req := apiClient.ListNetworkAreas(ctx, *model.OrganizationId)
	if model.LabelSelector != nil {
		req = req.LabelSelector(*model.LabelSelector)
	}
	return req
}

func outputResult(p *print.Printer, outputFormat string, networkAreas []iaas.NetworkArea) error {
	return p.OutputResult(outputFormat, networkAreas, func() error {
		table := tables.NewTable()
		table.SetHeader("ID", "Name", "# Attached Projects")

		for _, networkArea := range networkAreas {
			table.AddRow(
				utils.PtrString(networkArea.Id),
				utils.PtrString(networkArea.Name),
				utils.PtrString(networkArea.ProjectCount),
			)
			table.AddSeparator()
		}

		p.Outputln(table.Render())
		return nil
	})
}
