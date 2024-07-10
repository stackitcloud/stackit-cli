package describe

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/goccy/go-yaml"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/sqlserverflex/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-sdk-go/services/sqlserverflex"

	"github.com/spf13/cobra"
)

const (
	databaseNameArg = "DATABASE_NAME"

	instanceIdFlag = "instance-id"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	DatabaseName string
	InstanceId   string
}

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("describe %s", databaseNameArg),
		Short: "Shows details of an SQLServer Flex database",
		Long:  "Shows details of an SQLServer Flex database.",
		Args:  args.SingleArg(databaseNameArg, nil),
		Example: examples.Build(
			examples.NewExample(
				`Get details of an SQLServer Flex database with name "my-database" of instance with ID "xxx"`,
				"$ stackit beta sqlserverflex database describe my-database --instance-id xxx"),
			examples.NewExample(
				`Get details of an SQLServer Flex database with name "my-database" of instance with ID "xxx" in JSON format`,
				"$ stackit beta sqlserverflex database describe my-database --instance-id xxx --output-format json"),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			model, err := parseInput(p, cmd, args)
			if err != nil {
				return err
			}
			// Configure API client
			apiClient, err := client.ConfigureClient(p)
			if err != nil {
				return err
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("read SQLServer Flex database: %w", err)
			}

			return outputResult(p, model.OutputFormat, resp)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), instanceIdFlag, "SQLServer Flex instance ID")

	err := flags.MarkFlagsRequired(cmd, instanceIdFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	databaseName := inputArgs[0]

	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		DatabaseName:    databaseName,
		InstanceId:      flags.FlagToStringValue(p, cmd, instanceIdFlag),
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

func buildRequest(ctx context.Context, model *inputModel, apiClient *sqlserverflex.APIClient) sqlserverflex.ApiGetDatabaseRequest {
	req := apiClient.GetDatabase(ctx, model.ProjectId, model.InstanceId, model.DatabaseName)
	return req
}

func outputResult(p *print.Printer, outputFormat string, database *sqlserverflex.GetDatabaseResponse) error {
	switch outputFormat {
	case print.JSONOutputFormat:
		details, err := json.MarshalIndent(database, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal SQLServer Flex database: %w", err)
		}
		p.Outputln(string(details))

		return nil
	case print.YAMLOutputFormat:
		details, err := yaml.MarshalWithOptions(database, yaml.IndentSequence(true))
		if err != nil {
			return fmt.Errorf("marshal SQLServer Flex database: %w", err)
		}
		p.Outputln(string(details))

		return nil
	default:
		database := database.Database
		table := tables.NewTable()
		table.AddRow("ID", *database.Id)
		table.AddSeparator()
		table.AddRow("NAME", *database.Name)
		table.AddSeparator()
		if database.CreateDate != nil {
			table.AddRow("CREATE DATE", *database.CreateDate)
			table.AddSeparator()
		}
		if database.Collation != nil {
			table.AddRow("COLLATION", *database.Collation)
			table.AddSeparator()
		}
		if database.Options != nil {
			if database.Options.CompatibilityLevel != nil {
				table.AddRow("COMPATIBILITY LEVEL", *database.Options.CompatibilityLevel)
				table.AddSeparator()
			}
			if database.Options.IsEncrypted != nil {
				table.AddRow("IS ENCRYPTED", *database.Options.IsEncrypted)
				table.AddSeparator()
			}
			if database.Options.Owner != nil {
				table.AddRow("OWNER", *database.Options.Owner)
				table.AddSeparator()
			}
			if database.Options.UserAccess != nil {
				table.AddRow("USER ACCESS", *database.Options.UserAccess)
			}
		}
		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}

		return nil
	}
}
