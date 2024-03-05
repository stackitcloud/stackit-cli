package list

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/secrets-manager/client"
	secretsManagerUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/secrets-manager/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"

	"github.com/spf13/cobra"
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

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Lists all Secrets Manager users",
		Long:  "Lists all Secrets Manager users.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`List all Secrets Manager users of instance with ID "xxx`,
				"$ stackit secrets-manager user list --instance-id xxx"),
			examples.NewExample(
				`List all Secrets Manager users in JSON format with ID "xxx`,
				"$ stackit secrets-manager user list --instance-id xxx --output-format json"),
			examples.NewExample(
				`List up to 10 Secrets Manager users with ID "xxx"`,
				"$ stackit secrets-manager user list --instance-id xxx --limit 10"),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			model, err := parseInput(cmd)
			if err != nil {
				return err
			}

			// Configure API client
			apiClient, err := client.ConfigureClient(cmd)
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
					instanceLabel = *model.InstanceId
				}
				cmd.Printf("No users found for instance %q\n", instanceLabel)
				return nil
			}
			users := *resp.Users

			// Truncate output
			if model.Limit != nil && len(users) > int(*model.Limit) {
				users = users[:*model.Limit]
			}

			return outputResult(cmd, model.OutputFormat, users)
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

func parseInput(cmd *cobra.Command) (*inputModel, error) {
	globalFlags := globalflags.Parse(cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	limit := flags.FlagToInt64Pointer(cmd, limitFlag)
	if limit != nil && *limit < 1 {
		return nil, &errors.FlagValidationError{
			Flag:    limitFlag,
			Details: "must be greater than 0",
		}
	}

	return &inputModel{
		GlobalFlagModel: globalFlags,
		InstanceId:      flags.FlagToStringPointer(cmd, instanceIdFlag),
		Limit:           flags.FlagToInt64Pointer(cmd, limitFlag),
	}, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *secretsmanager.APIClient) secretsmanager.ApiListUsersRequest {
	req := apiClient.ListUsers(ctx, model.ProjectId, *model.InstanceId)
	return req
}

func outputResult(cmd *cobra.Command, outputFormat string, users []secretsmanager.User) error {
	switch outputFormat {
	case globalflags.JSONOutputFormat:
		details, err := json.MarshalIndent(users, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal Secrets Manager user list: %w", err)
		}
		cmd.Println(string(details))

		return nil
	default:
		table := tables.NewTable()
		table.SetHeader("ID", "USERNAME")
		for i := range users {
			user := users[i]
			table.AddRow(*user.Id, *user.Username)
		}
		err := table.Display(cmd)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}

		return nil
	}
}
