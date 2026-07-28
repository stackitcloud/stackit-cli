package options

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/types"

	"github.com/spf13/cobra"
	postgresflex "github.com/stackitcloud/stackit-sdk-go/services/postgresflex/v3api"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	cliErr "github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/postgresflex/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
)

const (
	// Deprecated: Will be removed after 2027-01-31.
	flavorsFlag = "flavors"
	// Deprecated: Will be removed after 2027-01-31.
	versionsFlag = "versions"
	// Deprecated: Will be removed after 2027-01-31.
	storagesFlag = "storages"
	// Deprecated: Will be removed after 2027-01-31.
	flavorIdFlag = "flavor-id"
)

// Deprecated: Will be removed after 2027-01-31.
type inputModel struct {
	*globalflags.GlobalFlagModel

	Flavors  bool
	Versions bool
	Storages bool
	FlavorId *string
}

// Deprecated: Will be removed after 2027-01-31.
type options struct {
	Flavors  []postgresflex.ListFlavors `json:"flavors,omitempty"`
	Versions []postgresflex.Version     `json:"versions,omitempty"`
	Storages *flavorStorages            `json:"flavorStorages,omitempty"`
}

// Deprecated: Will be removed after 2027-01-31.
type flavorStorages struct {
	FlavorId string                                          `json:"flavorId"`
	Storages []postgresflex.FlavorStorageClassesStorageClass `json:"storages"`
}

// Deprecated: Will be removed after 2027-01-31.
func NewCmd(params *types.CmdParams) *cobra.Command {
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
		Deprecated: `Command "stackit postgresflex options" is deprecated and will be removed after 2027-01-31. Please use "stackit postgresflex version list", "stackit postgresflex flavors list" and "stackit postgresflex flavor describe" instead.`,
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
			options, err := buildAndExecuteRequest(ctx, model, apiClient.DefaultAPI)
			if err != nil {
				return fmt.Errorf("get PostgreSQL Flex options: %w", err)
			}

			return outputResult(params.Printer, model, options)
		},
	}
	configureFlags(cmd)
	return cmd
}

// Deprecated: Will be removed after 2027-01-31.
func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Bool(flavorsFlag, false, "Lists supported flavors")
	cmd.Flags().Bool(versionsFlag, false, "Lists supported versions")
	cmd.Flags().Bool(storagesFlag, false, "Lists supported storages for a given flavor")
	cmd.Flags().String(flavorIdFlag, "", fmt.Sprintf("The flavor ID to show storages for. Only relevant when \"--%s\" is passed", storagesFlag))
}

// Deprecated: Will be removed after 2027-01-31.
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

// Deprecated: Will be removed after 2027-01-31.
func buildAndExecuteRequest(ctx context.Context, model *inputModel, apiClient postgresflex.DefaultAPI) (*options, error) {
	options := options{}

	if model.Flavors || model.Storages {
		flavors, err := apiClient.ListFlavors(ctx, model.ProjectId, model.Region).Execute()
		if err != nil {
			return nil, fmt.Errorf("get PostgreSQL Flex flavors: %w", err)
		}

		if model.Flavors {
			options.Flavors = flavors.Flavors
		}

		if model.Storages && model.FlavorId != nil {
			for _, f := range flavors.Flavors {
				if f.Id == *model.FlavorId {
					options.Storages = &flavorStorages{
						FlavorId: f.Id,
						Storages: f.StorageClasses,
					}
				}
			}

			if options.Storages == nil {
				return nil, fmt.Errorf("couldn't find flavor with id \"%s\"", *model.FlavorId)
			}
		}
	}

	if model.Versions {
		versions, err := apiClient.ListVersions(ctx, model.ProjectId, model.Region).Execute()
		if err != nil {
			return nil, fmt.Errorf("get PostgreSQL Flex versions: %w", err)
		}

		options.Versions = versions.Versions
	}

	return &options, nil
}

// Deprecated: Will be removed after 2027-01-31.
func outputResult(p *print.Printer, model *inputModel, options *options) error {
	if model == nil || model.GlobalFlagModel == nil {
		return fmt.Errorf("model is nil")
	}

	return p.OutputResult(model.OutputFormat, options, func() error {
		if options == nil {
			return fmt.Errorf("options is nil")
		}

		content := []tables.Table{}
		if model.Flavors && len(options.Flavors) != 0 {
			content = append(content, buildFlavorsTable(options.Flavors))
		}
		if model.Versions && len(options.Versions) != 0 {
			content = append(content, buildVersionsTable(options.Versions))
		}
		if model.Storages && options.Storages != nil && len(options.Storages.Storages) > 0 {
			content = append(content, buildStoragesTable(options.Storages.Storages))
		}

		err := tables.DisplayTables(p, content)
		if err != nil {
			return fmt.Errorf("display output: %w", err)
		}

		return nil
	})
}

// Deprecated: Will be removed after 2027-01-31.
func buildFlavorsTable(flavors []postgresflex.ListFlavors) tables.Table {
	table := tables.NewTable()
	table.SetTitle("Flavors")
	table.SetHeader("ID", "CPU", "MEMORY", "DESCRIPTION")
	for _, f := range flavors {
		table.AddRow(
			f.Id,
			f.Cpu,
			f.Memory,
			f.Description,
		)
		table.AddSeparator()
	}
	return table
}

// Deprecated: Will be removed after 2027-01-31.
func buildVersionsTable(versions []postgresflex.Version) tables.Table {
	table := tables.NewTable()
	table.SetTitle("Versions")
	table.SetHeader("VERSION", "RECOMMEND", "BETA", "DEPRECATED")

	for _, v := range versions {
		table.AddRow(v.Version, v.Recommend, v.Beta, v.Deprecated)
		table.AddSeparator()
	}
	return table
}

// Deprecated: Will be removed after 2027-01-31.
func buildStoragesTable(storageClasses []postgresflex.FlavorStorageClassesStorageClass) tables.Table {
	table := tables.NewTable()
	table.SetTitle("Storages")
	table.SetHeader("STORAGE CLASS", "MAX IO PER SEC", "MAX THROUGH (MB)")

	for _, sc := range storageClasses {
		table.AddRow(sc.Class, sc.MaxIoPerSec, sc.MaxThroughInMb)
		table.AddSeparator()
	}
	return table
}
