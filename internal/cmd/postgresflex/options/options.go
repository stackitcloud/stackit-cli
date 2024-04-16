package options

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/pager"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/postgresflex/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"

	"github.com/spf13/cobra"
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

func NewCmd(p *print.Printer) *cobra.Command {
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
			model, err := parseInput(cmd, p)
			if err != nil {
				return err
			}

			// Configure API client
			apiClient, err := client.ConfigureClient(p)
			if err != nil {
				return err
			}

			// Call API
			err = buildAndExecuteRequest(ctx, p, model, apiClient)
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

func parseInput(cmd *cobra.Command, p *print.Printer) (*inputModel, error) {
	globalFlags := globalflags.Parse(cmd, p)
	flavors := flags.FlagToBoolValue(cmd, flavorsFlag, p)
	versions := flags.FlagToBoolValue(cmd, versionsFlag, p)
	storages := flags.FlagToBoolValue(cmd, storagesFlag, p)
	flavorId := flags.FlagToStringPointer(cmd, flavorIdFlag, p)

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

	return &inputModel{
		GlobalFlagModel: globalFlags,
		Flavors:         flavors,
		Versions:        versions,
		Storages:        storages,
		FlavorId:        flags.FlagToStringPointer(cmd, flavorIdFlag, p),
	}, nil
}

type postgresFlexOptionsClient interface {
	ListFlavorsExecute(ctx context.Context, projectId string) (*postgresflex.ListFlavorsResponse, error)
	ListVersionsExecute(ctx context.Context, projectId string) (*postgresflex.ListVersionsResponse, error)
	ListStoragesExecute(ctx context.Context, projectId, flavorId string) (*postgresflex.ListStoragesResponse, error)
}

func buildAndExecuteRequest(ctx context.Context, p *print.Printer, model *inputModel, apiClient postgresFlexOptionsClient) error {
	var flavors *postgresflex.ListFlavorsResponse
	var versions *postgresflex.ListVersionsResponse
	var storages *postgresflex.ListStoragesResponse
	var err error

	if model.Flavors {
		flavors, err = apiClient.ListFlavorsExecute(ctx, model.ProjectId)
		if err != nil {
			return fmt.Errorf("get PostgreSQL Flex flavors: %w", err)
		}
	}
	if model.Versions {
		versions, err = apiClient.ListVersionsExecute(ctx, model.ProjectId)
		if err != nil {
			return fmt.Errorf("get PostgreSQL Flex versions: %w", err)
		}
	}
	if model.Storages {
		storages, err = apiClient.ListStoragesExecute(ctx, model.ProjectId, *model.FlavorId)
		if err != nil {
			return fmt.Errorf("get PostgreSQL Flex storages: %w", err)
		}
	}

	return outputResult(p, model, flavors, versions, storages)
}

func outputResult(p *print.Printer, model *inputModel, flavors *postgresflex.ListFlavorsResponse, versions *postgresflex.ListVersionsResponse, storages *postgresflex.ListStoragesResponse) error {
	options := &options{}
	if flavors != nil {
		options.Flavors = flavors.Flavors
	}
	if versions != nil {
		options.Versions = versions.Versions
	}
	if storages != nil && model.FlavorId != nil {
		options.Storages = &flavorStorages{
			FlavorId: *model.FlavorId,
			Storages: storages,
		}
	}

	switch model.OutputFormat {
	case globalflags.JSONOutputFormat:
		details, err := json.MarshalIndent(options, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal PostgreSQL Flex options: %w", err)
		}
		p.Outputln(string(details))
		return nil
	default:
		return outputResultAsTable(p, model, options)
	}
}

func outputResultAsTable(p *print.Printer, model *inputModel, options *options) error {
	content := ""
	if model.Flavors {
		content += renderFlavors(*options.Flavors)
	}
	if model.Versions {
		content += renderVersions(*options.Versions)
	}
	if model.Storages {
		content += renderStorages(options.Storages.Storages)
	}

	err := pager.Display(p, content)
	if err != nil {
		return fmt.Errorf("display output: %w", err)
	}

	return nil
}

func renderFlavors(flavors []postgresflex.Flavor) string {
	if len(flavors) == 0 {
		return ""
	}

	table := tables.NewTable()
	table.SetHeader("ID", "CPU", "MEMORY", "DESCRIPTION")
	for i := range flavors {
		f := flavors[i]
		table.AddRow(*f.Id, *f.Cpu, *f.Memory, *f.Description)
	}
	return table.Render()
}

func renderVersions(versions []string) string {
	if len(versions) == 0 {
		return ""
	}

	table := tables.NewTable()
	table.SetHeader("VERSION")
	for i := range versions {
		v := versions[i]
		table.AddRow(v)
	}
	return table.Render()
}

func renderStorages(resp *postgresflex.ListStoragesResponse) string {
	if resp.StorageClasses == nil || len(*resp.StorageClasses) == 0 {
		return ""
	}
	storageClasses := *resp.StorageClasses

	table := tables.NewTable()
	table.SetHeader("MIN STORAGE", "MAX STORAGE", "STORAGE CLASS")
	for i := range storageClasses {
		sc := storageClasses[i]
		table.AddRow(*resp.StorageRange.Min, *resp.StorageRange.Max, sc)
	}
	table.EnableAutoMergeOnColumns(1, 2, 3)
	return table.Render()
}
