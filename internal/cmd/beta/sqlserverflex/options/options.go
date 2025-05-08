package options

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/goccy/go-yaml"
	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/sqlserverflex/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/sqlserverflex"
)

// enforce implementation of interfaces
var (
	_ sqlServerFlexOptionsClient = &sqlserverflex.APIClient{}
)

type sqlServerFlexOptionsClient interface {
	ListFlavorsExecute(ctx context.Context, projectId string, region string) (*sqlserverflex.ListFlavorsResponse, error)
	ListVersionsExecute(ctx context.Context, projectId string, region string) (*sqlserverflex.ListVersionsResponse, error)
	ListStoragesExecute(ctx context.Context, projectId, flavorId string, region string) (*sqlserverflex.ListStoragesResponse, error)
	ListRolesExecute(ctx context.Context, projectId string, instanceId string, region string) (*sqlserverflex.ListRolesResponse, error)
	ListCollationsExecute(ctx context.Context, projectId string, instanceId string, region string) (*sqlserverflex.ListCollationsResponse, error)
	ListCompatibilityExecute(ctx context.Context, projectId string, instanceId string, region string) (*sqlserverflex.ListCompatibilityResponse, error)
}

const (
	flavorsFlag           = "flavors"
	versionsFlag          = "versions"
	storagesFlag          = "storages"
	userRolesFlag         = "user-roles"
	dbCollationsFlag      = "db-collations"
	dbCompatibilitiesFlag = "db-compatibilities"

	flavorIdFlag   = "flavor-id"
	instanceIdFlag = "instance-id"
)

type inputModel struct {
	*globalflags.GlobalFlagModel

	Flavors           bool
	Versions          bool
	Storages          bool
	UserRoles         bool
	DBCollations      bool
	DBCompatibilities bool

	FlavorId   *string
	InstanceId *string
}

type options struct {
	Flavors           *[]sqlserverflex.InstanceFlavorEntry `json:"flavors,omitempty"`
	Versions          *[]string                            `json:"versions,omitempty"`
	Storages          *flavorStorages                      `json:"flavorStorages,omitempty"`
	UserRoles         *instanceUserRoles                   `json:"userRoles,omitempty"`
	DBCollations      *instanceDBCollations                `json:"dbCollations,omitempty"`
	DBCompatibilities *instanceDBCompatibilities           `json:"dbCompatibilities,omitempty"`
}

type flavorStorages struct {
	FlavorId string                              `json:"flavorId"`
	Storages *sqlserverflex.ListStoragesResponse `json:"storages"`
}

type instanceUserRoles struct {
	InstanceId string   `json:"instanceId"`
	UserRoles  []string `json:"userRoles"`
}

type instanceDBCollations struct {
	InstanceId   string                                 `json:"instanceId"`
	DBCollations []sqlserverflex.MssqlDatabaseCollation `json:"dbCollations"`
}

type instanceDBCompatibilities struct {
	InstanceId        string                                     `json:"instanceId"`
	DBCompatibilities []sqlserverflex.MssqlDatabaseCompatibility `json:"dbCompatibilities"`
}

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "options",
		Short: "Lists SQL Server Flex options",
		Long:  "Lists SQL Server Flex options (flavors, versions and storages for a given flavor)\nPass one or more flags to filter what categories are shown.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`List SQL Server Flex flavors options`,
				"$ stackit beta sqlserverflex options --flavors"),
			examples.NewExample(
				`List SQL Server Flex available versions`,
				"$ stackit beta sqlserverflex options --versions"),
			examples.NewExample(
				`List SQL Server Flex storage options for a given flavor. The flavor ID can be retrieved by running "$ stackit beta sqlserverflex options --flavors"`,
				"$ stackit beta sqlserverflex options --storages --flavor-id <FLAVOR_ID>"),
			examples.NewExample(
				`List SQL Server Flex user roles and database compatibilities for a given instance. The IDs of existing instances can be obtained by running "$ stackit beta sqlserverflex instance list"`,
				"$ stackit beta sqlserverflex options --user-roles --db-compatibilities --instance-id <INSTANCE_ID>"),
		),
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := context.Background()
			model, err := parseInput(params.Printer, cmd)
			if err != nil {
				return err
			}

			// Configure API client
			apiClient, err := client.ConfigureClient(params.Printer)
			if err != nil {
				return err
			}

			// Call API
			err = buildAndExecuteRequest(ctx, params.Printer, model, apiClient)
			if err != nil {
				return fmt.Errorf("get SQL Server Flex options: %w", err)
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
	cmd.Flags().Bool(userRolesFlag, false, "Lists supported user roles for a given instance")
	cmd.Flags().Bool(dbCollationsFlag, false, "Lists supported database collations for a given instance")
	cmd.Flags().Bool(dbCompatibilitiesFlag, false, "Lists supported database compatibilities for a given instance")
	cmd.Flags().String(flavorIdFlag, "", `The flavor ID to show storages for. Only relevant when "--storages" is passed`)
	cmd.Flags().String(instanceIdFlag, "", `The instance ID to show user roles, database collations and database compatibilities for. Only relevant when "--user-roles", "--db-collations" or "--db-compatibilities" is passed`)
}

func parseInput(p *print.Printer, cmd *cobra.Command) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)

	flavors := flags.FlagToBoolValue(p, cmd, flavorsFlag)
	versions := flags.FlagToBoolValue(p, cmd, versionsFlag)
	storages := flags.FlagToBoolValue(p, cmd, storagesFlag)
	userRoles := flags.FlagToBoolValue(p, cmd, userRolesFlag)
	dbCollations := flags.FlagToBoolValue(p, cmd, dbCollationsFlag)
	dbCompatibilities := flags.FlagToBoolValue(p, cmd, dbCompatibilitiesFlag)

	flavorId := flags.FlagToStringPointer(p, cmd, flavorIdFlag)
	instanceId := flags.FlagToStringPointer(p, cmd, instanceIdFlag)

	if !flavors && !versions && !storages && !userRoles && !dbCollations && !dbCompatibilities {
		return nil, fmt.Errorf("%s\n\n%s",
			"please specify at least one category for which to list the available options.",
			"Get details on the available flags by re-running your command with the --help flag.")
	}

	if storages && flavorId == nil {
		return nil, fmt.Errorf("%s\n\n%s\n%s",
			`please specify a flavor ID to show storages for by setting the flag "--flavor-id <FLAVOR_ID>".`,
			"You can get the available flavor IDs by running:",
			"  $ stackit beta sqlserverflex options --flavors")
	}

	if (userRoles || dbCollations || dbCompatibilities) && instanceId == nil {
		return nil, fmt.Errorf("%s\n\n%s\n%s",
			`please specify an instance ID to show user roles, database collations or database compatibilities for by setting the flag "--instance-id <INSTANCE_ID>".`,
			"You can get the available instances and their IDs by running:",
			"  $ stackit beta sqlserverflex instance list")
	}

	model := inputModel{
		GlobalFlagModel:   globalFlags,
		Flavors:           flavors,
		Versions:          versions,
		Storages:          storages,
		UserRoles:         userRoles,
		DBCollations:      dbCollations,
		DBCompatibilities: dbCompatibilities,
		FlavorId:          flavorId,
		InstanceId:        instanceId,
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

func buildAndExecuteRequest(ctx context.Context, p *print.Printer, model *inputModel, apiClient sqlServerFlexOptionsClient) error {
	var flavors *sqlserverflex.ListFlavorsResponse
	var versions *sqlserverflex.ListVersionsResponse
	var storages *sqlserverflex.ListStoragesResponse
	var userRoles *sqlserverflex.ListRolesResponse
	var dbCollations *sqlserverflex.ListCollationsResponse
	var dbCompatibilities *sqlserverflex.ListCompatibilityResponse
	var err error

	if model.Flavors {
		flavors, err = apiClient.ListFlavorsExecute(ctx, model.ProjectId, model.Region)
		if err != nil {
			return fmt.Errorf("get SQL Server Flex flavors: %w", err)
		}
	}
	if model.Versions {
		versions, err = apiClient.ListVersionsExecute(ctx, model.ProjectId, model.Region)
		if err != nil {
			return fmt.Errorf("get SQL Server Flex versions: %w", err)
		}
	}
	if model.Storages {
		storages, err = apiClient.ListStoragesExecute(ctx, model.ProjectId, *model.FlavorId, model.Region)
		if err != nil {
			return fmt.Errorf("get SQL Server Flex storages: %w", err)
		}
	}
	if model.UserRoles {
		userRoles, err = apiClient.ListRolesExecute(ctx, model.ProjectId, *model.InstanceId, model.Region)
		if err != nil {
			return fmt.Errorf("get SQL Server Flex user roles: %w", err)
		}
	}
	if model.DBCollations {
		dbCollations, err = apiClient.ListCollationsExecute(ctx, model.ProjectId, *model.InstanceId, model.Region)
		if err != nil {
			return fmt.Errorf("get SQL Server Flex DB collations: %w", err)
		}
	}
	if model.DBCompatibilities {
		dbCompatibilities, err = apiClient.ListCompatibilityExecute(ctx, model.ProjectId, *model.InstanceId, model.Region)
		if err != nil {
			return fmt.Errorf("get SQL Server Flex DB compatibilities: %w", err)
		}
	}

	return outputResult(p, model, flavors, versions, storages, userRoles, dbCollations, dbCompatibilities)
}

func outputResult(p *print.Printer, model *inputModel, flavors *sqlserverflex.ListFlavorsResponse, versions *sqlserverflex.ListVersionsResponse, storages *sqlserverflex.ListStoragesResponse, userRoles *sqlserverflex.ListRolesResponse, dbCollations *sqlserverflex.ListCollationsResponse, dbCompatibilities *sqlserverflex.ListCompatibilityResponse) error {
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
	if userRoles != nil && model.InstanceId != nil {
		options.UserRoles = &instanceUserRoles{
			InstanceId: *model.InstanceId,
			UserRoles:  *userRoles.Roles,
		}
	}
	if dbCollations != nil && model.InstanceId != nil {
		options.DBCollations = &instanceDBCollations{
			InstanceId:   *model.InstanceId,
			DBCollations: *dbCollations.Collations,
		}
	}
	if dbCompatibilities != nil && model.InstanceId != nil {
		options.DBCompatibilities = &instanceDBCompatibilities{
			InstanceId:        *model.InstanceId,
			DBCompatibilities: *dbCompatibilities.Compatibilities,
		}
	}

	switch model.OutputFormat {
	case print.JSONOutputFormat:
		details, err := json.MarshalIndent(options, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal SQL Server Flex options: %w", err)
		}
		p.Outputln(string(details))
		return nil
	case print.YAMLOutputFormat:
		details, err := yaml.MarshalWithOptions(options, yaml.IndentSequence(true), yaml.UseJSONMarshaler())
		if err != nil {
			return fmt.Errorf("marshal SQL Server Flex options: %w", err)
		}
		p.Outputln(string(details))

		return nil
	default:
		return outputResultAsTable(p, model, options)
	}
}

func outputResultAsTable(p *print.Printer, model *inputModel, options *options) error {
	content := []tables.Table{}
	if model.Flavors && len(*options.Flavors) != 0 {
		content = append(content, buildFlavorsTable(*options.Flavors))
	}
	if model.Versions && len(*options.Versions) != 0 {
		content = append(content, buildVersionsTable(*options.Versions))
	}
	if model.Storages && options.Storages.Storages != nil && len(*options.Storages.Storages.StorageClasses) != 0 {
		content = append(content, buildStoragesTable(*options.Storages.Storages))
	}
	if model.UserRoles && len(options.UserRoles.UserRoles) != 0 {
		content = append(content, buildUserRoles(options.UserRoles))
	}
	if model.DBCompatibilities && len(options.DBCompatibilities.DBCompatibilities) != 0 {
		content = append(content, buildDBCompatibilitiesTable(options.DBCompatibilities.DBCompatibilities))
	}
	// Rendered at last because table is very long
	if model.DBCollations && len(options.DBCollations.DBCollations) != 0 {
		content = append(content, buildDBCollationsTable(options.DBCollations.DBCollations))
	}

	err := tables.DisplayTables(p, content)
	if err != nil {
		return fmt.Errorf("display output: %w", err)
	}

	return nil
}

func buildFlavorsTable(flavors []sqlserverflex.InstanceFlavorEntry) tables.Table {
	table := tables.NewTable()
	table.SetTitle("Flavors")
	table.SetHeader("ID", "CPU", "MEMORY", "DESCRIPTION", "VALID INSTANCE TYPES")
	for i := range flavors {
		f := flavors[i]
		table.AddRow(*f.Id, *f.Cpu, *f.Memory, *f.Description, *f.Categories)
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

func buildStoragesTable(storagesResp sqlserverflex.ListStoragesResponse) tables.Table {
	storages := *storagesResp.StorageClasses
	table := tables.NewTable()
	table.SetTitle("Storages")
	table.SetHeader("MINIMUM", "MAXIMUM", "STORAGE CLASS")
	for i := range storages {
		sc := storages[i]
		table.AddRow(*storagesResp.StorageRange.Min, *storagesResp.StorageRange.Max, sc)
	}
	table.EnableAutoMergeOnColumns(1, 2, 3)
	return table
}

func buildUserRoles(roles *instanceUserRoles) tables.Table {
	table := tables.NewTable()
	table.SetTitle("User Roles")
	table.SetHeader("ROLE")
	for i := range roles.UserRoles {
		table.AddRow(roles.UserRoles[i])
	}
	return table
}

func buildDBCollationsTable(dbCollations []sqlserverflex.MssqlDatabaseCollation) tables.Table {
	table := tables.NewTable()
	table.SetTitle("DB Collations")
	table.SetHeader("NAME", "DESCRIPTION")
	for i := range dbCollations {
		table.AddRow(dbCollations[i].CollationName, dbCollations[i].Description)
	}
	return table
}

func buildDBCompatibilitiesTable(dbCompatibilities []sqlserverflex.MssqlDatabaseCompatibility) tables.Table {
	table := tables.NewTable()
	table.SetTitle("DB Compatibilities")
	table.SetHeader("COMPATIBILITY LEVEL", "DESCRIPTION")
	for i := range dbCompatibilities {
		table.AddRow(dbCompatibilities[i].CompatibilityLevel, dbCompatibilities[i].Description)
	}
	return table
}
