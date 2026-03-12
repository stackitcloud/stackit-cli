package list

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/types"

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
	Limit  *int64
	Filter *string
}

const (
	limitFlag  = "limit"
	filterFlag = "filter"
)

func NewCmd(params *types.CmdParams) *cobra.Command {
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
			examples.NewExample(
				`List machine types with exactly 2 vCPUs`,
				`$ stackit server machine-type list --filter="vcpus==2"`,
			),
			examples.NewExample(
				`List machine types with at least 2 vCPUs and 2048 MB RAM`,
				`$ stackit server machine-type list --filter="vcpus >= 2 && ram >= 2048"`,
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
	cmd.Flags().String(filterFlag, "", "Filter resources by fields. A subset of expr-lang is supported. See https://expr-lang.org/docs/language-definition for usage details")
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
		Limit:           limit,
		Filter:          flags.FlagToStringPointer(p, cmd, filterFlag),
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *iaas.APIClient) iaas.ApiListMachineTypesRequest {
	req := apiClient.ListMachineTypes(ctx, model.ProjectId, model.Region)
	if model.Filter != nil {
		req = req.Filter(*model.Filter)
	}
	return req
}

func outputResult(p *print.Printer, outputFormat string, machineTypes iaas.MachineTypeListResponse) error {
	return p.OutputResult(outputFormat, machineTypes, func() error {
		table := tables.NewTable()
		table.SetTitle("Machine-Types")
		table.SetHeader("NAME", "VCPUS", "RAM (GB)", "DESCRIPTION", "EXTRA SPECS")
		if items := machineTypes.GetItems(); len(items) > 0 {
			for _, machineType := range items {
				extraSpecMap := make(map[string]string)
				if machineType.ExtraSpecs != nil && len(*machineType.ExtraSpecs) > 0 {
					for key, value := range *machineType.ExtraSpecs {
						extraSpecMap[key] = fmt.Sprintf("%v", value)
					}
				}
				ramGB := int64(0)
				if machineType.Ram != nil {
					ramGB = *machineType.Ram / 1024
				}

				table.AddRow(
					utils.PtrString(machineType.Name),
					utils.PtrValue(machineType.Vcpus),
					ramGB,
					utils.PtrString(machineType.Description),
					utils.JoinStringMap(extraSpecMap, ": ", "\n"),
				)
				table.AddSeparator()
			}
		}

		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}

		return nil
	})
}
