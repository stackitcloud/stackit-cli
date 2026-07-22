package describe

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/types"

	"github.com/spf13/cobra"
	sqlserverflex "github.com/stackitcloud/stackit-sdk-go/services/sqlserverflex/v3api"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/sqlserverflex/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
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

func NewCmd(params *types.CmdParams) *cobra.Command {
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
				return fmt.Errorf("read SQLServer Flex database: %w", err)
			}

			return outputResult(params.Printer, model.OutputFormat, resp)
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

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *sqlserverflex.APIClient) sqlserverflex.ApiGetDatabaseRequest {
	req := apiClient.DefaultAPI.GetDatabase(ctx, model.ProjectId, model.Region, model.InstanceId, model.DatabaseName)
	return req
}

func outputResult(p *print.Printer, outputFormat string, resp *sqlserverflex.GetDatabaseResponse) error {
	if resp == nil {
		return fmt.Errorf("database response is empty")
	}

	return p.OutputResult(outputFormat, resp, func() error {
		table := tables.NewTable()
		table.AddRow("ID", resp.Id)
		table.AddSeparator()
		table.AddRow("NAME", resp.Name)
		table.AddSeparator()
		table.AddRow("COMPATIBILITY LEVEL", resp.CompatibilityLevel)
		table.AddSeparator()
		table.AddRow("OWNER", resp.Owner)
		table.AddSeparator()
		table.AddRow("COLLATION", resp.CollationName)
		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}

		return nil
	})
}
