package list

import (
	"context"
	"encoding/json"
	"fmt"

	"stackit/internal/pkg/args"
	"stackit/internal/pkg/errors"
	"stackit/internal/pkg/examples"
	"stackit/internal/pkg/flags"
	"stackit/internal/pkg/globalflags"
	"stackit/internal/pkg/services/mongodbflex/client"
	mongodbflexUtils "stackit/internal/pkg/services/mongodbflex/utils"
	"stackit/internal/pkg/tables"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/mongodbflex"
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
		Short: "List all MongoDB Flex users of an instance",
		Long:  "List all MongoDB Flex users of an instance.",
		Example: examples.Build(
			examples.NewExample(
				`List all MongoDB Flex users of instance with ID "xxx"`,
				"$ stackit mongodbflex user list --instance-id xxx"),
			examples.NewExample(
				`List all MongoDB Flex users of instance with ID "xxx" in JSON format`,
				"$ stackit mongodbflex user list --instance-id xxx --output-format json"),
			examples.NewExample(
				`List up to 10 MongoDB Flex users of instance with ID "xxx"`,
				"$ stackit mongodbflex user list --instance-id xxx --limit 10"),
		),
		Args: args.NoArgs,
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
				return fmt.Errorf("get MongoDB Flex users: %w", err)
			}
			if resp.Items == nil || len(*resp.Items) == 0 {
				instanceLabel, err := mongodbflexUtils.GetInstanceName(ctx, apiClient, model.ProjectId, *model.InstanceId)
				if err != nil {
					instanceLabel = *model.InstanceId
				}
				cmd.Printf("No users found for instance %s\n", instanceLabel)
				return nil
			}
			users := *resp.Items

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

func buildRequest(ctx context.Context, model *inputModel, apiClient *mongodbflex.APIClient) mongodbflex.ApiListUsersRequest {
	req := apiClient.ListUsers(ctx, model.ProjectId, *model.InstanceId)
	return req
}

func outputResult(cmd *cobra.Command, outputFormat string, users []mongodbflex.ListUser) error {
	switch outputFormat {
	case globalflags.JSONOutputFormat:
		details, err := json.MarshalIndent(users, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal MongoDB Flex user list: %w", err)
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
