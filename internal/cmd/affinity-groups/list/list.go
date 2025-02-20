package list

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/goccy/go-yaml"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/iaas/client"
	"github.com/stackitcloud/stackit-sdk-go/services/iaas"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	Limit *int64
}

const limitFlag = "limit"

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Lists affinity groups",
		Long:  `Lists affinity groups.`,
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				"Lists all affinity groups",
				"$ stackit affinity-group list",
			),
			examples.NewExample(
				"Lists up to 10 affinity groups",
				"$ stackit affinity-group list --limit=10",
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
			request := buildRequest(ctx, *model, apiClient)
			result, err := request.Execute()
			if err != nil {
				return fmt.Errorf("list affinity groups: %w", err)
			}

			if items := result.Items; items != nil {
				if model.Limit != nil && len(*items) > int(*model.Limit) {
					*items = (*items)[:*model.Limit]
				}
				return outputResult(p, *model, *items)
			}

			p.Outputln("No affinity groups found")
			return nil
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Int64(limitFlag, 0, "Limit the output to the first n elements")
}

func buildRequest(ctx context.Context, model inputModel, apiClient *iaas.APIClient) iaas.ApiListAffinityGroupsRequest {
	return apiClient.ListAffinityGroups(ctx, model.ProjectId)
}

func parseInput(p *print.Printer, cmd *cobra.Command) (*inputModel, error) {
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
		Limit:           limit,
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

func outputResult(p *print.Printer, model inputModel, items []iaas.AffinityGroup) error {
	var outputFormat string
	if model.GlobalFlagModel != nil {
		outputFormat = model.GlobalFlagModel.OutputFormat
	}
	switch outputFormat {
	case print.JSONOutputFormat:
		details, err := json.MarshalIndent(items, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal affinity groups: %w", err)
		}
		p.Outputln(string(details))
	case print.YAMLOutputFormat:
		details, err := yaml.MarshalWithOptions(items, yaml.IndentSequence(true), yaml.UseJSONMarshaler())
		if err != nil {
			return fmt.Errorf("marshal affinity groups: %w", err)
		}
		p.Outputln(string(details))
	default:
		table := tables.NewTable()
		table.SetHeader("ID", "NAME", "POLICY")
		for _, item := range items {
			table.AddRow(
				utils.PtrString(item.Id),
				utils.PtrString(item.Name),
				utils.PtrString(item.Policy),
			)
			table.AddSeparator()
		}

		if err := table.Display(p); err != nil {
			return fmt.Errorf("render table: %w", err)
		}
	}
	return nil
}
