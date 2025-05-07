package list

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/goccy/go-yaml"
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

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/iaas"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	Limit *int64
}

const (
	limitFlag = "limit"
)

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Get list of all machine types available in a project",
		Long:  "Get list of all machine types available in a project.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Get list of all machine types`,
				"$ stackit server machine-type list",
			),
			examples.NewExample(
				`Get list of all machine types in JSON format`,
				"$ stackit server machine-type list --output-format json",
			),
			examples.NewExample(
				`List the first 10 machine types`,
				`$ stackit server machine-type list --limit=10`,
			),
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

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("read machine-types: %w", err)
			}

			if resp.Items == nil || len(*resp.Items) == 0 {
				projectLabel, err := projectname.GetProjectName(ctx, params.Printer, params.CliVersion, cmd)
				if err != nil {
					params.Printer.Debug(print.ErrorLevel, "get project name: %v", err)
					projectLabel = model.ProjectId
				}
				params.Printer.Info("No machine-types found for project %q\n", projectLabel)
				return nil
			}

			// limit output
			if model.Limit != nil && len(*resp.Items) > int(*model.Limit) {
				*resp.Items = (*resp.Items)[:*model.Limit]
			}

			return outputResult(params.Printer, model.OutputFormat, *resp)
		},
	}

	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Int64(limitFlag, 0, "Limit the output to the first n elements")
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
		Limit:           flags.FlagToInt64Pointer(p, cmd, limitFlag),
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

func buildRequest(ctx context.Context, model *inputModel, apiClient *iaas.APIClient) iaas.ApiListMachineTypesRequest {
	return apiClient.ListMachineTypes(ctx, model.ProjectId)
}

func outputResult(p *print.Printer, outputFormat string, machineTypes iaas.MachineTypeListResponse) error {
	switch outputFormat {
	case print.JSONOutputFormat:
		details, err := json.MarshalIndent(machineTypes, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal machineTypes: %w", err)
		}
		p.Outputln(string(details))

		return nil
	case print.YAMLOutputFormat:
		details, err := yaml.MarshalWithOptions(machineTypes, yaml.IndentSequence(true), yaml.UseJSONMarshaler())
		if err != nil {
			return fmt.Errorf("marshal machineTypes: %w", err)
		}
		p.Outputln(string(details))

		return nil
	default:
		table := tables.NewTable()
		table.SetTitle("Machine-Types")

		table.SetHeader("NAME", "DESCRIPTION")
		if items := machineTypes.GetItems(); len(items) > 0 {
			for _, machineType := range items {
				table.AddRow(*machineType.Name, utils.PtrString(machineType.Description))
			}
		}

		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}

		return nil
	}
}
