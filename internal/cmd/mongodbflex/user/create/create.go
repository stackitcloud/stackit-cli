package create

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/mongodbflex/client"
	mongodbflexUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/mongodbflex/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/mongodbflex"
)

const (
	instanceIdFlag = "instance-id"
	usernameFlag   = "username"
	databaseFlag   = "database"
	roleFlag       = "role"
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

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Creates a MongoDB Flex user",
		Long: fmt.Sprintf("%s\n%s\n%s\n%s",
			"Creates a MongoDB Flex user.",
			"The password is only visible upon creation and cannot be retrieved later.",
			"Alternatively, you can reset the password and access the new one by running:",
			"  $ stackit mongodbflex user reset-password USER_ID --instance-id INSTANCE_ID",
		),
		Example: examples.Build(
			examples.NewExample(
				`Create a MongoDB Flex user for instance with ID "xxx" and specify the username`,
				"$ stackit mongodbflex user create --instance-id xxx --username johndoe --role read --database default"),
			examples.NewExample(
				`Create a MongoDB Flex user for instance with ID "xxx" with an automatically generated username`,
				"$ stackit mongodbflex user create --instance-id xxx --role read --database default"),
		),
		Args: args.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			model, err := parseInput(p, cmd)
			if err != nil {
				return err
			}

			// Configure API client
			apiClient, err := client.ConfigureClient(p)
			if err != nil {
				return err
			}

			instanceLabel, err := mongodbflexUtils.GetInstanceName(ctx, apiClient, model.ProjectId, model.InstanceId)
			if err != nil {
				p.Debug(print.ErrorLevel, "get instance name: %v", err)
				instanceLabel = model.InstanceId
			}

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to create a user for instance %q?", instanceLabel)
				err = p.PromptForConfirmation(prompt)
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

			return outputResult(p, model, instanceLabel, user)
		},
	}

	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	roleOptions := []string{"read", "readWrite"}

	cmd.Flags().Var(flags.UUIDFlag(), instanceIdFlag, "ID of the instance")
	cmd.Flags().String(usernameFlag, "", "Username of the user. If not specified, a random username will be assigned")
	cmd.Flags().String(databaseFlag, "", "The database inside the MongoDB instance that the user has access to. If it does not exist, it will be created once the user writes to it")
	cmd.Flags().Var(flags.EnumSliceFlag(false, rolesDefault, roleOptions...), roleFlag, fmt.Sprintf("Roles of the user, possible values are %q", roleOptions))

	err := flags.MarkFlagsRequired(cmd, instanceIdFlag, databaseFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	return &inputModel{
		GlobalFlagModel: globalFlags,
		InstanceId:      flags.FlagToStringValue(p, cmd, instanceIdFlag),
		Username:        flags.FlagToStringPointer(p, cmd, usernameFlag),
		Database:        flags.FlagToStringPointer(p, cmd, databaseFlag),
		Roles:           flags.FlagWithDefaultToStringSlicePointer(p, cmd, roleFlag),
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

func outputResult(p *print.Printer, model *inputModel, instanceLabel string, user *mongodbflex.User) error {
	switch model.OutputFormat {
	case print.JSONOutputFormat:
		details, err := json.MarshalIndent(user, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal MongoDB Flex user: %w", err)
		}
		p.Outputln(string(details))

		return nil
	default:
		p.Outputf("Created user for instance %q. User ID: %s\n\n", instanceLabel, *user.Id)
		p.Outputf("Username: %s\n", *user.Username)
		p.Outputf("Password: %s\n", *user.Password)
		p.Outputf("Roles: %v\n", *user.Roles)
		p.Outputf("Database: %s\n", *user.Database)
		p.Outputf("Host: %s\n", *user.Host)
		p.Outputf("Port: %d\n", *user.Port)
		p.Outputf("URI: %s\n", *user.Uri)

		return nil
	}
}
