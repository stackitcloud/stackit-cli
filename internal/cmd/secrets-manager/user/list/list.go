package list

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/types"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/secrets-manager/client"
	secretsManagerUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/secrets-manager/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/secretsmanager"
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

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Lists all Secrets Manager users",
		Long:  "Lists all Secrets Manager users.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`List all Secrets Manager users of instance with ID "xxx"`,
				"$ stackit secrets-manager user list --instance-id xxx"),
			examples.NewExample(
				`List all Secrets Manager users in JSON format with ID "xxx"`,
				"$ stackit secrets-manager user list --instance-id xxx --output-format json"),
			examples.NewExample(
				`List up to 10 Secrets Manager users with ID "xxx"`,
				"$ stackit secrets-manager user list --instance-id xxx --limit 10"),
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
				return fmt.Errorf("get Secrets Manager users: %w", err)
			}
			if resp.Users == nil || len(*resp.Users) == 0 {
				instanceLabel, err := secretsManagerUtils.GetInstanceName(ctx, apiClient, model.ProjectId, *model.InstanceId)
				if err != nil {
					params.Printer.Debug(print.ErrorLevel, "get instance name: %v", err)
					instanceLabel = *model.InstanceId
				}
				params.Printer.Info("No users found for instance %q\n", instanceLabel)
				return nil
			}
			users := *resp.Users

			// Truncate output
			if model.Limit != nil && len(users) > int(*model.Limit) {
				users = users[:*model.Limit]
			}

			return outputResult(params.Printer, model.OutputFormat, users)
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

func buildRequest(ctx context.Context, model *inputModel, apiClient *secretsmanager.APIClient) secretsmanager.ApiListUsersRequest {
	req := apiClient.ListUsers(ctx, model.ProjectId, *model.InstanceId)
	return req
}

func outputResult(p *print.Printer, outputFormat string, users []secretsmanager.User) error {
	return p.OutputResult(outputFormat, users, func() error {
		table := tables.NewTable()
		table.SetHeader("ID", "USERNAME", "DESCRIPTION", "WRITE ACCESS")
		for i := range users {
			user := users[i]
			table.AddRow(
				utils.PtrString(user.Id),
				utils.PtrString(user.Username),
				utils.PtrString(user.Description),
				utils.PtrString(user.Write),
			)
		}
		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}

		return nil
	})
}
