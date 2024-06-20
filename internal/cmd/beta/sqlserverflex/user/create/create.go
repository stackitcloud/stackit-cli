package create

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/goccy/go-yaml"
	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/sqlserverflex/client"
	sqlserverflexUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/sqlserverflex/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/sqlserverflex"
)

const (
	instanceIdFlag = "instance-id"
	usernameFlag   = "username"
	rolesFlag      = "roles"
)

type inputModel struct {
	*globalflags.GlobalFlagModel

	InstanceId string
	Username   *string
	Roles      *[]string
}

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Creates a SQLServer Flex user",
		Long: fmt.Sprintf("%s\n%s\n%s\n%s\n%s",
			"Creates a SQLServer Flex user for an instance.",
			"The password is only visible upon creation and cannot be retrieved later.",
			"Alternatively, you can reset the password and access the new one by running:",
			"  $ stackit beta sqlserverflex user reset-password USER_ID --instance-id INSTANCE_ID",
			"Please refer to https://docs.stackit.cloud/stackit/en/creating-logins-and-users-in-sqlserver-flex-instances-210862358.html for additional information.",
		),
		Example: examples.Build(
			examples.NewExample(
				`Create a SQLServer Flex user for instance with ID "xxx" and specify the username, role and database`,
				"$ stackit beta sqlserverflex user create --instance-id xxx --username johndoe --roles my-role --database my-database"),
			examples.NewExample(
				`Create a SQLServer Flex user for instance with ID "xxx", specifying multiple roles`,
				`$ stackit beta sqlserverflex user create --instance-id xxx --username johndoe --roles "my-role-1,my-role-2`),
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

			instanceLabel, err := sqlserverflexUtils.GetInstanceName(ctx, apiClient, model.ProjectId, model.InstanceId)
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
				return fmt.Errorf("create SQLServer Flex user: %w", err)
			}
			user := resp.Item

			return outputResult(p, model, instanceLabel, user)
		},
	}

	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), instanceIdFlag, "ID of the instance")
	cmd.Flags().String(usernameFlag, "", "Username of the user")
	cmd.Flags().StringSlice(rolesFlag, []string{}, "Roles of the user")

	err := flags.MarkFlagsRequired(cmd, instanceIdFlag, usernameFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		InstanceId:      flags.FlagToStringValue(p, cmd, instanceIdFlag),
		Username:        flags.FlagToStringPointer(p, cmd, usernameFlag),
		Roles:           flags.FlagToStringSlicePointer(p, cmd, rolesFlag),
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

func buildRequest(ctx context.Context, model *inputModel, apiClient *sqlserverflex.APIClient) sqlserverflex.ApiCreateUserRequest {
	req := apiClient.CreateUser(ctx, model.ProjectId, model.InstanceId)

	var roles []sqlserverflex.Role
	if model.Roles != nil {
		for _, r := range *model.Roles {
			roles = append(roles, sqlserverflex.Role(r))
		}
	}

	req = req.CreateUserPayload(sqlserverflex.CreateUserPayload{
		Username: model.Username,
		Roles:    &roles,
	})
	return req
}

func outputResult(p *print.Printer, model *inputModel, instanceLabel string, user *sqlserverflex.User) error {
	switch model.OutputFormat {
	case print.JSONOutputFormat:
		details, err := json.MarshalIndent(user, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal SQLServer Flex user: %w", err)
		}
		p.Outputln(string(details))

		return nil
	case print.YAMLOutputFormat:
		details, err := yaml.MarshalWithOptions(user, yaml.IndentSequence(true))
		if err != nil {
			return fmt.Errorf("marshal SQLServer Flex user: %w", err)
		}
		p.Outputln(string(details))

		return nil
	default:
		p.Outputf("Created user for instance %q. User ID: %s\n\n", instanceLabel, *user.Id)
		p.Outputf("Username: %s\n", *user.Username)
		p.Outputf("Password: %s\n", *user.Password)
		if user.Roles != nil && len(*user.Roles) != 0 {
			p.Outputf("Roles: %v\n", *user.Roles)
		}
		if user.Host != nil && *user.Host != "" {
			p.Outputf("Host: %s\n", *user.Host)
		}
		if user.Port != nil {
			p.Outputf("Port: %d\n", *user.Port)
		}
		if user.Uri != nil && *user.Uri != "" {
			p.Outputf("URI: %s\n", *user.Uri)
		}

		return nil
	}
}
