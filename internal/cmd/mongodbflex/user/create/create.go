package create

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/confirm"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/mongodbflex/client"
	mongodbflexUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/mongodbflex/utils"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/mongodbflex"
)

const (
	instanceIdFlag = "instance-id"
	usernameFlag   = "username"
	databaseFlag   = "database"
	rolesFlag      = "roles"
)

var (
	rolesDefault = []string{"read"}
)

type inputModel struct {
	*globalflags.GlobalFlagModel

	InstanceId string
	Username   *string
	Database   *string
	Roles      *[]string
}

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a MongoDB Flex user",
		Long: fmt.Sprintf("%s\n%s\n%s\n%s",
			"Create a MongoDB Flex user.",
			"The password is only visible upon creation and cannot be retrieved later.",
			"Alternatively, you can reset the password and access the new one by running:",
			"  $ stackit mongodbflex user reset-password --instance-id <INSTANCE_ID> --user-id <USER_ID>",
		),
		Example: examples.Build(
			examples.NewExample(
				`Create a MongoDB Flex user for instance with ID "xxx" and specify the username`,
				"$ stackit mongodbflex user create --instance-id xxx --username johndoe --roles read --database default"),
			examples.NewExample(
				`Create a MongoDB Flex user for instance with ID "xxx" with an automatically generated username`,
				"$ stackit mongodbflex user create --instance-id xxx --roles read --database default"),
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

			instanceLabel, err := mongodbflexUtils.GetInstanceName(ctx, apiClient, model.ProjectId, model.InstanceId)
			if err != nil {
				instanceLabel = model.InstanceId
			}

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to create a user for instance %s?", instanceLabel)
				err = confirm.PromptForConfirmation(cmd, prompt)
				if err != nil {
					return err
				}
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("create MongoDB Flex user: %w", err)
			}
			user := resp.Item

			cmd.Printf("Created user for instance %s. User ID: %s\n\n", instanceLabel, *user.Id)
			cmd.Printf("Username: %s\n", *user.Username)
			cmd.Printf("Password: %s\n", *user.Password)
			cmd.Printf("Roles: %v\n", *user.Roles)
			cmd.Printf("Database: %s\n", *user.Database)
			cmd.Printf("Host: %s\n", *user.Host)
			cmd.Printf("Port: %d\n", *user.Port)
			cmd.Printf("URI: %s\n", *user.Uri)

			return nil
		},
	}

	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	rolesOptions := []string{"read", "readWrite"}

	cmd.Flags().Var(flags.UUIDFlag(), instanceIdFlag, "ID of the instance")
	cmd.Flags().String(usernameFlag, "", "Username of the user. If not specified, a random username will be assigned")
	cmd.Flags().String(databaseFlag, "", "The database inside the MongoDB instance that the user has access to. If it does not exist, it will be created once the user writes to it")
	cmd.Flags().Var(flags.EnumSliceFlag(false, rolesDefault, rolesOptions...), rolesFlag, fmt.Sprintf("Roles of the user, possible values are %q", rolesOptions))

	err := flags.MarkFlagsRequired(cmd, instanceIdFlag, databaseFlag)
	cobra.CheckErr(err)
}

func parseInput(cmd *cobra.Command) (*inputModel, error) {
	globalFlags := globalflags.Parse(cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	return &inputModel{
		GlobalFlagModel: globalFlags,
		InstanceId:      flags.FlagToStringValue(cmd, instanceIdFlag),
		Username:        flags.FlagToStringPointer(cmd, usernameFlag),
		Database:        flags.FlagToStringPointer(cmd, databaseFlag),
		Roles:           flags.FlagWithDefaultToStringSlicePointer(cmd, rolesFlag),
	}, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *mongodbflex.APIClient) mongodbflex.ApiCreateUserRequest {
	req := apiClient.CreateUser(ctx, model.ProjectId, model.InstanceId)
	req = req.CreateUserPayload(mongodbflex.CreateUserPayload{
		Username: model.Username,
		Database: model.Database,
		Roles:    model.Roles,
	})
	return req
}
