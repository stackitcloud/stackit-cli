package create

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/goccy/go-yaml"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/postgresflex/client"
	postgresflexUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/postgresflex/utils"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/postgresflex"
)

const (
	instanceIdFlag = "instance-id"
	usernameFlag   = "username"
	roleFlag       = "role"
)

var (
	rolesDefault = []string{"login"}
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
		Short: "Creates a PostgreSQL Flex user",
		Long: fmt.Sprintf("%s\n%s\n%s\n%s",
			"Creates a PostgreSQL Flex user.",
			"The password is only visible upon creation and cannot be retrieved later.",
			"Alternatively, you can reset the password and access the new one by running:",
			"  $ stackit postgresflex user reset-password USER_ID --instance-id INSTANCE_ID",
		),
		Example: examples.Build(
			examples.NewExample(
				`Create a PostgreSQL Flex user for instance with ID "xxx"`,
				"$ stackit postgresflex user create --instance-id xxx --username johndoe"),
			examples.NewExample(
				`Create a PostgreSQL Flex user for instance with ID "xxx" and permission "createdb"`,
				"$ stackit postgresflex user create --instance-id xxx --username johndoe --role createdb"),
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

			instanceLabel, err := postgresflexUtils.GetInstanceName(ctx, apiClient, model.ProjectId, model.InstanceId)
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
				return fmt.Errorf("create PostgreSQL Flex user: %w", err)
			}

			return outputResult(p, model, instanceLabel, resp)
		},
	}

	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	roleOptions := []string{"login", "createdb"}

	cmd.Flags().Var(flags.UUIDFlag(), instanceIdFlag, "ID of the instance")
	cmd.Flags().String(usernameFlag, "", "Username of the user")
	cmd.Flags().Var(flags.EnumSliceFlag(false, rolesDefault, roleOptions...), roleFlag, fmt.Sprintf("Roles of the user, possible values are %q", roleOptions))

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
		Roles:           flags.FlagWithDefaultToStringSlicePointer(p, cmd, roleFlag),
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

func buildRequest(ctx context.Context, model *inputModel, apiClient *postgresflex.APIClient) postgresflex.ApiCreateUserRequest {
	req := apiClient.CreateUser(ctx, model.ProjectId, model.InstanceId)
	req = req.CreateUserPayload(postgresflex.CreateUserPayload{
		Username: model.Username,
		Roles:    model.Roles,
	})
	return req
}

func outputResult(p *print.Printer, model *inputModel, instanceLabel string, resp *postgresflex.CreateUserResponse) error {
	switch model.OutputFormat {
	case print.JSONOutputFormat:
		details, err := json.MarshalIndent(resp, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal PostgresFlex user: %w", err)
		}
		p.Outputln(string(details))

		return nil
	case print.YAMLOutputFormat:
		details, err := yaml.Marshal(resp)
		if err != nil {
			return fmt.Errorf("marshal PostgresFlex user: %w", err)
		}
		p.Outputln(string(details))

		return nil
	default:
		user := resp.Item
		p.Outputf("Created user for instance %q. User ID: %s\n\n", instanceLabel, *user.Id)
		p.Outputf("Username: %s\n", *user.Username)
		p.Outputf("Password: %s\n", *user.Password)
		p.Outputf("Roles: %v\n", *user.Roles)
		p.Outputf("Host: %s\n", *user.Host)
		p.Outputf("Port: %d\n", *user.Port)
		p.Outputf("URI: %s\n", *user.Uri)

		return nil
	}
}
