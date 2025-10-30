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
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/sqlserverflex/client"
	sqlserverflexUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/sqlserverflex/utils"
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

	InstanceId *string
	Limit      *int64
}

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Lists all SQLServer Flex users of an instance",
		Long:  "Lists all SQLServer Flex users of an instance.",
		Example: examples.Build(
			examples.NewExample(
				`List all SQLServer Flex users of instance with ID "xxx"`,
				"$ stackit beta sqlserverflex user list --instance-id xxx"),
			examples.NewExample(
				`List all SQLServer Flex users of instance with ID "xxx" in JSON format`,
				"$ stackit beta sqlserverflex user list --instance-id xxx --output-format json"),
			examples.NewExample(
				`List up to 10 SQLServer Flex users of instance with ID "xxx"`,
				"$ stackit beta sqlserverflex user list --instance-id xxx --limit 10"),
		),
		Args: args.NoArgs,
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
				return fmt.Errorf("get SQLServer Flex users: %w", err)
			}
			users := resp.GetItems()

			instanceLabel, err := sqlserverflexUtils.GetInstanceName(ctx, apiClient, model.ProjectId, *model.InstanceId, model.Region)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get instance name: %v", err)
				instanceLabel = *model.InstanceId
			}

			// Truncate output
			if model.Limit != nil && len(users) > int(*model.Limit) {
				users = users[:*model.Limit]
			}

			return outputResult(params.Printer, model.OutputFormat, instanceLabel, users)
		},
	}

	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), instanceIdFlag, "Instance ID")
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
		InstanceId:      flags.FlagToStringPointer(p, cmd, instanceIdFlag),
		Limit:           flags.FlagToInt64Pointer(p, cmd, limitFlag),
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *sqlserverflex.APIClient) sqlserverflex.ApiListUsersRequest {
	req := apiClient.ListUsers(ctx, model.ProjectId, *model.InstanceId, model.Region)
	return req
}

func outputResult(p *print.Printer, outputFormat, instanceLabel string, users []sqlserverflex.InstanceListUser) error {
	return p.OutputResult(outputFormat, users, func() error {
		if len(users) == 0 {
			p.Outputf("No users found for instance %q\n", instanceLabel)
			return nil
		}

		table := tables.NewTable()
		table.SetHeader("ID", "USERNAME")
		for i := range users {
			user := users[i]
			table.AddRow(
				utils.PtrString(user.Id),
				utils.PtrString(user.Username),
			)
		}
		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}

		return nil
	})
}
