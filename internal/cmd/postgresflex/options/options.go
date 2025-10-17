package options

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/cmd/params"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	cliErr "github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/postgresflex/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/postgresflex"
)

const (
	flavorsFlag  = "flavors"
	versionsFlag = "versions"
	storagesFlag = "storages"
	flavorIdFlag = "flavor-id"
)

type inputModel struct {
	*globalflags.GlobalFlagModel

	Flavors  bool
	Versions bool
	Storages bool
	FlavorId *string
}

type options struct {
	Flavors  *[]postgresflex.Flavor `json:"flavors,omitempty"`
	Versions *[]string              `json:"versions,omitempty"`
	Storages *flavorStorages        `json:"flavorStorages,omitempty"`
}

type flavorStorages struct {
	FlavorId string                             `json:"flavorId"`
	Storages *postgresflex.ListStoragesResponse `json:"storages"`
}

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "options",
		Short: "Lists PostgreSQL Flex options",
		Long:  "Lists PostgreSQL Flex options (flavors, versions and storages for a given flavor)\nPass one or more flags to filter what categories are shown.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`List PostgreSQL Flex flavors options`,
				"$ stackit postgresflex options --flavors"),
			examples.NewExample(
				`List PostgreSQL Flex available versions`,
				"$ stackit postgresflex options --versions"),
			examples.NewExample(
				`List PostgreSQL Flex storage options for a given flavor. The flavor ID can be retrieved by running "$ stackit postgresflex options --flavors"`,
				"$ stackit postgresflex options --storages --flavor-id <FLAVOR_ID>"),
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
			err = buildAndExecuteRequest(ctx, params.Printer, model, apiClient)
			if err != nil {
				return fmt.Errorf("get PostgreSQL Flex options: %w", err)
			}

			return nil
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Bool(flavorsFlag, false, "Lists supported flavors")
	cmd.Flags().Bool(versionsFlag, false, "Lists supported versions")
	cmd.Flags().Bool(storagesFlag, false, "Lists supported storages for a given flavor")
	cmd.Flags().String(flavorIdFlag, "", `The flavor ID to show storages for. Only relevant when "--storages" is passed`)
}

func parseInput(p *print.Printer, cmd *cobra.Command, _ []string) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &cliErr.ProjectIdError{}
	}
	flavors := flags.FlagToBoolValue(p, cmd, flavorsFlag)
	versions := flags.FlagToBoolValue(p, cmd, versionsFlag)
	storages := flags.FlagToBoolValue(p, cmd, storagesFlag)
	flavorId := flags.FlagToStringPointer(p, cmd, flavorIdFlag)

	if !flavors && !versions && !storages {
		return nil, fmt.Errorf("%s\n\n%s",
			"please specify at least one category for which to list the available options.",
			"Get details on the available flags by re-running your command with the --help flag.")
	}

	if storages && flavorId == nil {
		return nil, fmt.Errorf("%s\n\n%s\n%s",
			`please specify a flavor ID to show storages for by setting the flag "--flavor-id <FLAVOR_ID>".`,
			"You can get the available flavor IDs by running:",
			"  $ stackit postgresflex options --flavors")
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		Flavors:         flavors,
		Versions:        versions,
		Storages:        storages,
		FlavorId:        flags.FlagToStringPointer(p, cmd, flavorIdFlag),
	}

	p.DebugInputModel(model)
	return &model, nil
}

type postgresFlexOptionsClient interface {
	ListFlavorsExecute(ctx context.Context, projectId, region string) (*postgresflex.ListFlavorsResponse, error)
	ListVersionsExecute(ctx context.Context, projectId, region string) (*postgresflex.ListVersionsResponse, error)
	ListStoragesExecute(ctx context.Context, projectId, region, flavorId string) (*postgresflex.ListStoragesResponse, error)
}

func buildAndExecuteRequest(ctx context.Context, p *print.Printer, model *inputModel, apiClient postgresFlexOptionsClient) error {
	var flavors *postgresflex.ListFlavorsResponse
	var versions *postgresflex.ListVersionsResponse
	var storages *postgresflex.ListStoragesResponse
	var err error

	if model.Flavors {
		flavors, err = apiClient.ListFlavorsExecute(ctx, model.ProjectId, model.Region)
		if err != nil {
			return fmt.Errorf("get PostgreSQL Flex flavors: %w", err)
		}
	}
	if model.Versions {
		versions, err = apiClient.ListVersionsExecute(ctx, model.ProjectId, model.Region)
		if err != nil {
			return fmt.Errorf("get PostgreSQL Flex versions: %w", err)
		}
	}
	if model.Storages {
		storages, err = apiClient.ListStoragesExecute(ctx, model.ProjectId, model.Region, *model.FlavorId)
		if err != nil {
			return fmt.Errorf("get PostgreSQL Flex storages: %w", err)
		}
	}

	return outputResult(p, *model, flavors, versions, storages)
}

func outputResult(p *print.Printer, model inputModel, flavors *postgresflex.ListFlavorsResponse, versions *postgresflex.ListVersionsResponse, storages *postgresflex.ListStoragesResponse) error {
	options := &options{}
	if flavors != nil {
		options.Flavors = flavors.Flavors
	}
	if model.GlobalFlagModel == nil {
		return fmt.Errorf("no global model defined")
	}
	if versions != nil {
		options.Versions = versions.Versions
	}
	if storages != nil && model.FlavorId != nil {
		options.Storages = &flavorStorages{
			FlavorId: utils.PtrString(model.FlavorId),
			Storages: storages,
		}
	}

	return p.OutputResult(model.OutputFormat, options, func() error {
		content := []tables.Table{}
		if model.Flavors && len(*options.Flavors) != 0 {
			content = append(content, buildFlavorsTable(*options.Flavors))
		}
		if model.Versions && len(*options.Versions) != 0 {
			content = append(content, buildVersionsTable(*options.Versions))
		}
		if model.Storages && options.Storages.Storages != nil && len(*options.Storages.Storages.StorageClasses) > 0 {
			content = append(content, buildStoragesTable(*options.Storages.Storages))
		}

		err := tables.DisplayTables(p, content)
		if err != nil {
			return fmt.Errorf("display output: %w", err)
		}

		return nil
	})
}

func buildFlavorsTable(flavors []postgresflex.Flavor) tables.Table {
	table := tables.NewTable()
	table.SetTitle("Flavors")
	table.SetHeader("ID", "CPU", "MEMORY", "DESCRIPTION")
	for i := range flavors {
		f := flavors[i]
		table.AddRow(
			utils.PtrString(f.Id),
			utils.PtrString(f.Cpu),
			utils.PtrString(f.Memory),
			utils.PtrString(f.Description),
		)
	}
	return table
}

func buildVersionsTable(versions []string) tables.Table {
	table := tables.NewTable()
	table.SetTitle("Versions")
	table.SetHeader("VERSION")
	for i := range versions {
		v := versions[i]
		table.AddRow(v)
	}
	return table
}

func buildStoragesTable(storagesResp postgresflex.ListStoragesResponse) tables.Table {
	storages := *storagesResp.StorageClasses
	table := tables.NewTable()
	table.SetTitle("Storages")
	table.SetHeader("MINIMUM", "MAXIMUM", "STORAGE CLASS")
	for i := range storages {
		sc := storages[i]
		table.AddRow(
			utils.PtrString(storagesResp.StorageRange.Min),
			utils.PtrString(storagesResp.StorageRange.Max),
			sc,
		)
	}
	table.EnableAutoMergeOnColumns(1, 2, 3)
	return table
}
