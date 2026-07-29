package describe

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/types"

	"github.com/spf13/cobra"
	postgresflex "github.com/stackitcloud/stackit-sdk-go/services/postgresflex/v3api"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	cliErr "github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/postgresflex/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
)

const (
	flavorIdArg = "FLAVOR_ID"
)

type inputModel struct {
	*globalflags.GlobalFlagModel

	FlavorId string
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "describe",
		Short: "Show details of a PostgreSQL Flex flavor",
		Long:  "Show details of a PostgreSQL Flex flavor.",
		Args:  args.SingleArg(flavorIdArg, nil),
		Example: examples.Build(
			examples.NewExample(
				`Show PostgreSQL Flex flavor details`,
				"$ stackit postgresflex flavor describe <FLAVOR_ID>"),
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
			flavor, err := buildAndExecuteRequest(ctx, model, apiClient.DefaultAPI)
			if err != nil {
				return fmt.Errorf("get PostgreSQL Flex flavor: %w", err)
			}

			return outputResult(params.Printer, model.OutputFormat, flavor)
		},
	}
	return cmd
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	flavorId := inputArgs[0]

	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &cliErr.ProjectIdError{}
	}
	model := inputModel{
		GlobalFlagModel: globalFlags,
		FlavorId:        flavorId,
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildAndExecuteRequest(ctx context.Context, model *inputModel, apiClient postgresflex.DefaultAPI) (*postgresflex.ListFlavors, error) {
	// Note: There's currently no GET endpoint for a single flavor. So we have to use the list endpoint and
	// do client-site filtering for the correct flavor. A feature request has been opened to add a GET endpoint
	// for a single flavor to the API.

	// TODO: Use GET endpoint instead of list endpoint after https://support.stackit.cloud/servicedesk/customer/portal/3/SSD-23311 is done

	flavors, err := apiClient.ListFlavors(ctx, model.ProjectId, model.Region).Size(100).Execute()
	if err != nil {
		return nil, fmt.Errorf("listing flavors: %w", err)
	}

	for _, flavor := range flavors.Flavors {
		if flavor.Id == model.FlavorId {
			return &flavor, nil
		}
	}

	return nil, fmt.Errorf("flavor \"%s\" not found", model.FlavorId)
}

func outputResult(p *print.Printer, outputFormat string, flavor *postgresflex.ListFlavors) error {
	return p.OutputResult(outputFormat, flavor, func() error {
		if flavor == nil {
			return fmt.Errorf("flavor is nil")
		}

		table := tables.NewTable()
		table.SetTitle("Flavor")
		table.AddRow("ID", flavor.Id)
		table.AddSeparator()
		table.AddRow("CPU", flavor.Cpu)
		table.AddSeparator()
		table.AddRow("MEMORY (GiB)", flavor.Memory)
		table.AddSeparator()
		table.AddRow("NODE TYPE", flavor.NodeType)
		table.AddSeparator()
		table.AddRow("MIN STORAGE (GB)", flavor.MinGB)
		table.AddSeparator()
		table.AddRow("MAX STORAGE (GB)", flavor.MaxGB)
		table.AddSeparator()
		table.AddRow("DESCRIPTION", flavor.Description)

		storageClassesTable := tables.NewTable()
		storageClassesTable.SetTitle("Storages")
		storageClassesTable.SetHeader("STORAGE CLASS", "MAX IO PER SEC", "MAX THROUGH (MB)")

		for _, sc := range flavor.StorageClasses {
			storageClassesTable.AddRow(sc.Class, sc.MaxIoPerSec, sc.MaxThroughInMb)
			storageClassesTable.AddSeparator()
		}

		return tables.DisplayTables(p, []tables.Table{table, storageClassesTable})
	})
}
