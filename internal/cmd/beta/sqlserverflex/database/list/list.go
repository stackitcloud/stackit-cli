package list

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/projectname"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/sqlserverflex/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/sqlserverflex"
)

const (
	instanceIdFlag = "instance-id"
	limitFlag      = "limit"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	InstanceId string
	Limit      *int64
}

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Lists all SQLServer Flex databases",
		Long:  "Lists all SQLServer Flex databases.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`List all SQLServer Flex databases of instance with ID "xxx"`,
				"$ stackit beta sqlserverflex database list --instance-id xxx"),
			examples.NewExample(
				`List all SQLServer Flex databases of instance with ID "xxx" in JSON format`,
				"$ stackit beta sqlserverflex database list --instance-id xxx --output-format json"),
			examples.NewExample(
				`List up to 10 SQLServer Flex databases of instance with ID "xxx"`,
				"$ stackit beta sqlserverflex database list --instance-id xxx --limit 10"),
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
				return fmt.Errorf("get SQLServer Flex databases: %w", err)
			}
			if resp.Databases == nil || len(*resp.Databases) == 0 {
				projectLabel, err := projectname.GetProjectName(ctx, params.Printer, params.CliVersion, cmd)
				if err != nil {
					params.Printer.Debug(print.ErrorLevel, "get project name: %v", err)
					projectLabel = model.ProjectId
				}
				params.Printer.Info("No databases found for instance %s on project %s\n", model.InstanceId, projectLabel)
				return nil
			}
			databases := *resp.Databases

			// Truncate output
			if model.Limit != nil && len(databases) > int(*model.Limit) {
				databases = databases[:*model.Limit]
			}

			return outputResult(params.Printer, model.OutputFormat, databases)
		},
	}

	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), instanceIdFlag, "SQLServer Flex instance ID")
	cmd.Flags().Int64(limitFlag, 0, "Maximum number of entries to list")

	err := flags.MarkFlagsRequired(cmd, instanceIdFlag)
	cobra.CheckErr(err)
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
		InstanceId:      flags.FlagToStringValue(p, cmd, instanceIdFlag),
		Limit:           flags.FlagToInt64Pointer(p, cmd, limitFlag),
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *sqlserverflex.APIClient) sqlserverflex.ApiListDatabasesRequest {
	req := apiClient.ListDatabases(ctx, model.ProjectId, model.InstanceId, model.Region)
	return req
}

func outputResult(p *print.Printer, outputFormat string, databases []sqlserverflex.Database) error {
	return p.OutputResult(outputFormat, databases, func() error {
		table := tables.NewTable()
		table.SetHeader("ID", "NAME")
		for i := range databases {
			database := databases[i]
			table.AddRow(utils.PtrString(database.Id), utils.PtrString(database.Name))
		}
		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}

		return nil
	})
}
