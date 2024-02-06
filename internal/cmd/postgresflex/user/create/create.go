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
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/postgresflex/client"
	postgresflexUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/postgresflex/utils"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/postgresflex"
)

const (
	instanceIdFlag = "instance-id"
	usernameFlag   = "username"
	rolesFlag      = "roles"
)

var (
	rolesDefault = []string{"read"}
)

type inputModel struct {
	*globalflags.GlobalFlagModel

	InstanceId string
	Username   *string
	Roles      *[]string
}

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Creates a PostgreSQL Flex user",
		Long: fmt.Sprintf("%s\n%s\n%s\n%s",
			"Creates a PostgreSQL Flex user.",
			"The password is only visible upon creation and cannot be retrieved later.",
			"Alternatively, you can reset the password and access the new one by running:",
			"  $ stackit postgresflex user reset-password USER_ID --instance-id INSTANCE_ID",
		),
		Example: examples.Build(
			examples.NewExample(
				`Create a PostgreSQL Flex user for instance with ID "xxx" and specify the username`,
				"$ stackit postgresflex user create --instance-id xxx --username johndoe --roles read"),
			examples.NewExample(
				`Create a PostgreSQL Flex user for instance with ID "xxx" with an automatically generated username`,
				"$ stackit postgresflex user create --instance-id xxx --roles read"),
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

			instanceLabel, err := postgresflexUtils.GetInstanceName(ctx, apiClient, model.ProjectId, model.InstanceId)
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
				return fmt.Errorf("create PostgreSQL Flex user: %w", err)
			}
			user := resp.Item

			cmd.Printf("Created user for instance %s. User ID: %s\n\n", instanceLabel, *user.Id)
			cmd.Printf("Username: %s\n", *user.Username)
			cmd.Printf("Password: %s\n", *user.Password)
			cmd.Printf("Roles: %v\n", *user.Roles)
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
	cmd.Flags().Var(flags.EnumSliceFlag(false, rolesDefault, rolesOptions...), rolesFlag, fmt.Sprintf("Roles of the user, possible values are %q", rolesOptions))

	err := flags.MarkFlagsRequired(cmd, instanceIdFlag)
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
		Roles:           flags.FlagWithDefaultToStringSlicePointer(cmd, rolesFlag),
	}, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *postgresflex.APIClient) postgresflex.ApiCreateUserRequest {
	req := apiClient.CreateUser(ctx, model.ProjectId, model.InstanceId)
	req = req.CreateUserPayload(postgresflex.CreateUserPayload{
		Username: model.Username,
		Roles:    model.Roles,
	})
	return req
}
