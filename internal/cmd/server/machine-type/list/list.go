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
	Limit    *int64
	MinVCPUs *int64
	MinRAM   *int64
}

const (
	limitFlag   = "limit"
	minVcpuFlag = "min-vcpu"
	minRamFlag  = "min-ram"
)

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Get list of all machine types available in a project",
		Long:  "Get list of all machine types available in a project. Supports filtering by minimum vCPU and RAM (GB).",
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
				`Filter for machines with at least 8 vCPUs and 16GB RAM`,
				"$ stackit server machine-type list --min-vcpu 8 --min-ram 16",
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
					projectLabel = model.ProjectId
				}
				params.Printer.Info("No machine-types found for project %q\n", projectLabel)
				return nil
			}

			// Filter the items client-side
			filteredItems := filterMachineTypes(resp.Items, model)

			if len(filteredItems) == 0 {
				params.Printer.Info("No machine-types found matching the criteria\n")
				return nil
			}

			// Apply limit to results
			if model.Limit != nil && len(filteredItems) > int(*model.Limit) {
				filteredItems = filteredItems[:*model.Limit]
			}

			resp.Items = &filteredItems
			return outputResult(params.Printer, model.OutputFormat, *resp)
		},
	}

	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Int64(limitFlag, 0, "Limit the output to the first n elements")
	cmd.Flags().Int64(minVcpuFlag, 0, "Filter by minimum number of vCPUs")
	cmd.Flags().Int64(minRamFlag, 0, "Filter by minimum RAM amount in GB")
}

func parseInput(p *print.Printer, cmd *cobra.Command, _ []string) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	limit := flags.FlagToInt64Pointer(p, cmd, limitFlag)
	if limit != nil && *limit < 1 {
		return nil, &errors.FlagValidationError{Flag: limitFlag, Details: "must be greater than 0"}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		Limit:           limit,
		MinVCPUs:        flags.FlagToInt64Pointer(p, cmd, minVcpuFlag),
		MinRAM:          flags.FlagToInt64Pointer(p, cmd, minRamFlag),
	}

	p.DebugInputModel(model)
	return &model, nil
}

// filterMachineTypes applies logic to filter by resource minimums.
// Discuss: hide deprecated machine-types?
func filterMachineTypes(items *[]iaas.MachineType, model *inputModel) []iaas.MachineType {
	if items == nil {
		return []iaas.MachineType{}
	}

	var filtered []iaas.MachineType
	for _, item := range *items {
		// Minimum vCPU check
		if model.MinVCPUs != nil && *model.MinVCPUs > 0 {
			if item.Vcpus == nil || *item.Vcpus < *model.MinVCPUs {
				continue
			}
		}

		// Minimum RAM check (converting API MB to GB)
		if model.MinRAM != nil && *model.MinRAM > 0 {
			if item.Ram == nil || (*item.Ram/1024) < *model.MinRAM {
				continue
			}
		}

		filtered = append(filtered, item)
	}
	return filtered
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *iaas.APIClient) iaas.ApiListMachineTypesRequest {
	return apiClient.ListMachineTypes(ctx, model.ProjectId, model.Region)
}

func outputResult(p *print.Printer, outputFormat string, machineTypes iaas.MachineTypeListResponse) error {
	return p.OutputResult(outputFormat, machineTypes, func() error {
		table := tables.NewTable()
		table.SetTitle("Machine-Types")
		table.SetHeader("NAME", "VCPUS", "RAM (GB)", "DESCRIPTION", "EXTRA SPECS")

		if items := machineTypes.GetItems(); len(items) > 0 {
			for _, mt := range items {
				extraSpecMap := make(map[string]string)
				if mt.ExtraSpecs != nil && len(*mt.ExtraSpecs) > 0 {
					for key, value := range *mt.ExtraSpecs {
						extraSpecMap[key] = fmt.Sprintf("%v", value)
					}
				}

				ramGB := int64(0)
				if mt.Ram != nil {
					ramGB = *mt.Ram / 1024
				}

				table.AddRow(
					utils.PtrString(mt.Name),
					utils.PtrValue(mt.Vcpus),
					ramGB,
					utils.PtrString(mt.Description),
					utils.JoinStringMap(extraSpecMap, ": ", "\n"),
				)
				table.AddSeparator()
			}
		}
		return table.Display(p)
	})
}
