package list

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/goccy/go-yaml"
	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/projectname"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/iaas/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-sdk-go/services/iaas"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	Labels string
}

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "list security groups",
		Long:  "list security groups",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(`list all groups`, `$ stackit beta security-group list`),
			examples.NewExample(`list groups with labels`, `$ stackit beta security-group list --labels label1=value1,label2=value2`),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			return executeList(cmd, p, args)
		},
	}

	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	globalflags.Configure(cmd.Flags())
	cmd.Flags().String("labels", "", "a list of labels in the form <key>=<value>")
}

func executeList(cmd *cobra.Command, p *print.Printer, _ []string) error {
	p.Info("executing list command")
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

	projectLabel, err := projectname.GetProjectName(ctx, p, cmd)
	if err != nil {
		p.Debug(print.ErrorLevel, "get project name: %v", err)
		projectLabel = model.ProjectId
	}

	// Call API
	request := buildRequest(ctx, model, apiClient)
	response, err := request.Execute()
	if err != nil {
		return fmt.Errorf("list security group: %w", err)
	}
	if items := response.GetItems(); items == nil || len(*items) > 0 {
		p.Info("no security groups found for %q", projectLabel)
	} else {
		outputResult(p, model.OutputFormat, *items)
	}

	return nil
}

func parseInput(p *print.Printer, cmd *cobra.Command) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,

		Labels: flags.FlagToStringValue(p, cmd, "labels"),
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

func buildRequest(ctx context.Context, model *inputModel, apiClient *iaas.APIClient) iaas.ApiListSecurityGroupsRequest {
	request := apiClient.ListSecurityGroups(ctx, model.ProjectId)
	request = request.LabelSelector(model.Labels)

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
		details, err := yaml.MarshalWithOptions(items, yaml.IndentSequence(true))
		if err != nil {
			return fmt.Errorf("marshal PostgreSQL Flex instance list: %w", err)
		}
		p.Outputln(string(details))

		return nil
	default:
		table := tables.NewTable()
		table.SetHeader("ID", "NAME", "LABELS", "STATEFUL")
		for _, item := range items {
			table.AddRow(item.Id, item.Name, concatLabels(item.Labels), item.Stateful)
		}
		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}

		return nil
	}
}

func concatLabels(item *map[string]any) string {
	if item == nil {
		return ""
	}
	var builder strings.Builder
	for k, v := range *item {
		builder.WriteString(fmt.Sprintf("%s=%v ", k, v))
	}
	return builder.String()
}
