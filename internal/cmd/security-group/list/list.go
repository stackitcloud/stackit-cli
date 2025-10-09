package list

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/goccy/go-yaml"
	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
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
	"github.com/stackitcloud/stackit-sdk-go/services/iaas"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	LabelSelector *string
}

const (
	labelSelectorFlag = "label-selector"
)

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Lists security groups",
		Long:  "Lists security groups by its internal ID.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(`List all groups`, `$ stackit security-group list`),
			examples.NewExample(`List groups with labels`, `$ stackit security-group list --label-selector label1=value1,label2=value2`),
		),
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := context.Background()
			model, err := parseInput(params.Printer, cmd)
			if err != nil {
				return err
			}

			// Configure API client
			apiClient, err := client.ConfigureClient(params.Printer, params.CliVersion)
			if err != nil {
				return err
			}

			projectLabel, err := projectname.GetProjectName(ctx, params.Printer, params.CliVersion, cmd)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get project name: %v", err)
				projectLabel = model.ProjectId
			}

			// Call API
			request := buildRequest(ctx, model, apiClient)

			response, err := request.Execute()
			if err != nil {
				return fmt.Errorf("list security group: %w", err)
			}

			if items := response.GetItems(); len(items) == 0 {
				params.Printer.Info("No security groups found for project %q", projectLabel)
			} else {
				if err := outputResult(params.Printer, model.OutputFormat, items); err != nil {
					return fmt.Errorf("output security groups: %w", err)
				}
			}

			return nil
		},
	}

	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().String(labelSelectorFlag, "", "Filter by label")
}

func parseInput(p *print.Printer, cmd *cobra.Command) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		LabelSelector:   flags.FlagToStringPointer(p, cmd, labelSelectorFlag),
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *iaas.APIClient) iaas.ApiListSecurityGroupsRequest {
	request := apiClient.ListSecurityGroups(ctx, model.ProjectId)
	if model.LabelSelector != nil {
		request = request.LabelSelector(*model.LabelSelector)
	}

	return request
}
func outputResult(p *print.Printer, outputFormat string, items []iaas.SecurityGroup) error {
	switch outputFormat {
	case print.JSONOutputFormat:
		details, err := json.MarshalIndent(items, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal PostgreSQL Flex instance list: %w", err)
		}
		p.Outputln(string(details))

		return nil
	case print.YAMLOutputFormat:
		details, err := yaml.MarshalWithOptions(items, yaml.IndentSequence(true), yaml.UseJSONMarshaler())
		if err != nil {
			return fmt.Errorf("marshal PostgreSQL Flex instance list: %w", err)
		}
		p.Outputln(string(details))

		return nil
	default:
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
	}
}
