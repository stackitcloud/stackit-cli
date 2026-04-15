package list

import (
	"context"
	"fmt"
	"strings"

	"github.com/stackitcloud/stackit-cli/internal/pkg/types"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/iaas"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/projectname"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/iaas/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	LabelSelector *string
	Limit         *int64
}

const (
	labelSelectorFlag = "label-selector"
	limitFlag         = "limit"
)

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Lists security groups",
		Long:  "Lists security groups by its internal ID.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(`Lists all security groups`, `$ stackit security-group list`),
			examples.NewExample(`Lists security groups with labels`, `$ stackit security-group list --label-selector label1=value1,label2=value2`),
			examples.NewExample(
				`Lists all security groups in JSON format`,
				"$ stackit security-group list --output-format json",
			),
			examples.NewExample(
				`Lists up to 10 security groups`,
				"$ stackit security-group list --limit 10",
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
			request := buildRequest(ctx, model, apiClient)

			response, err := request.Execute()
			if err != nil {
				return fmt.Errorf("list security group: %w", err)
			}

			items := response.GetItems()

			projectLabel, err := projectname.GetProjectName(ctx, params.Printer, params.CliVersion, cmd)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get project name: %v", err)
				projectLabel = model.ProjectId
			}

			// Truncate output
			if model.Limit != nil && len(items) > int(*model.Limit) {
				items = items[:*model.Limit]
			}

			return outputResult(params.Printer, model.OutputFormat, projectLabel, items)

		},
	}

	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().String(labelSelectorFlag, "", "Filter by label")
	cmd.Flags().Int64(limitFlag, 0, "Maximum number of entries to list")
}

func parseInput(p *print.Printer, cmd *cobra.Command, _ []string) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	limit := flags.FlagToInt64Pointer(p, cmd, limitFlag)
	if limit != nil && *limit < 1 {
		return nil, &errors.FlagValidationError{
			Flag:    limitFlag,
			Details: "must be greater than 0",
		}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		LabelSelector:   flags.FlagToStringPointer(p, cmd, labelSelectorFlag),
		Limit:           limit,
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *iaas.APIClient) iaas.ApiListSecurityGroupsRequest {
	request := apiClient.ListSecurityGroups(ctx, model.ProjectId, model.Region)
	if model.LabelSelector != nil {
		request = request.LabelSelector(*model.LabelSelector)
	}

	return request
}
func outputResult(p *print.Printer, outputFormat, projectLabel string, items []iaas.SecurityGroup) error {
	return p.OutputResult(outputFormat, items, func() error {
		if len(items) == 0 {
			p.Outputf("No security groups found for project %q\n", projectLabel)
			return nil
		}
		table := tables.NewTable()
		table.SetHeader("ID", "NAME", "STATEFUL", "DESCRIPTION", "LABELS")
		for _, item := range items {
			var labelsString string
			if item.Labels != nil {
				var labels []string
				for key, value := range *item.Labels {
					labels = append(labels, fmt.Sprintf("%s: %s", key, value))
				}
				labelsString = strings.Join(labels, ", ")
			}

			table.AddRow(
				utils.PtrString(item.Id),
				utils.PtrString(item.Name),
				utils.PtrString(item.Stateful),
				utils.PtrString(item.Description),
				labelsString,
			)
		}
		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}

		return nil
	})
}
