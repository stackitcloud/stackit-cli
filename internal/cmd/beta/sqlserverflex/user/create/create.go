package create

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/goccy/go-yaml"
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
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
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

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Creates a SQLServer Flex user",
		Long: fmt.Sprintf("%s\n\n%s\n%s\n%s\n%s\n\n%s\n%s",
			"Creates a SQLServer Flex user for an instance.",
			"The password is only visible upon creation and cannot be retrieved later.",
			"Alternatively, you can reset the password and access the new one by running:",
			"  $ stackit beta sqlserverflex user reset-password USER_ID --instance-id INSTANCE_ID",
			"Please refer to https://docs.stackit.cloud/stackit/en/creating-logins-and-users-in-sqlserver-flex-instances-210862358.html for additional information.",
			"The allowed user roles for your instance can be obtained by running:",
			"  $ stackit beta sqlserverflex options --user-roles --instance-id INSTANCE_ID",
		),
		Example: examples.Build(
			examples.NewExample(
				`Create a SQLServer Flex user for instance with ID "xxx" and specify the username, role and database`,
				`$ stackit beta sqlserverflex user create --instance-id xxx --username johndoe --roles "##STACKIT_DatabaseManager##"`),
			examples.NewExample(
				`Create a SQLServer Flex user for instance with ID "xxx", specifying multiple roles`,
				`$ stackit beta sqlserverflex user create --instance-id xxx --username johndoe --roles "##STACKIT_LoginManager##,##STACKIT_DatabaseManager##"`),
		),
		Args: args.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := context.Background()
			model, err := parseInput(params.Printer, cmd)
			if err != nil {
				return err
			}

			// Configure API client
			apiClient, err := client.ConfigureClient(params.Printer, params.CliVersion)
			if err != nil {
				return err
			}

			instanceLabel, err := sqlserverflexUtils.GetInstanceName(ctx, apiClient, model.ProjectId, model.InstanceId, model.Region)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get instance name: %v", err)
				instanceLabel = model.InstanceId
			}

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to create a user for instance %q?", instanceLabel)
				err = params.Printer.PromptForConfirmation(prompt)
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

			return outputResult(params.Printer, model, instanceLabel, user)
		},
	}

	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), instanceIdFlag, "ID of the instance")
	cmd.Flags().String(usernameFlag, "", "Username of the user")
	cmd.Flags().StringSlice(rolesFlag, []string{}, "Roles of the user")

	err := flags.MarkFlagsRequired(cmd, instanceIdFlag, usernameFlag, rolesFlag)
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
	req := apiClient.CreateUser(ctx, model.ProjectId, model.InstanceId, model.Region)

	req = req.CreateUserPayload(sqlserverflex.CreateUserPayload{
		Username: model.Username,
		Roles:    model.Roles,
	})
	return req
}

func outputResult(p *print.Printer, model *inputModel, instanceLabel string, user *sqlserverflex.SingleUser) error {
	if user == nil {
		return fmt.Errorf("user response is empty")
	}
	switch model.OutputFormat {
	case print.JSONOutputFormat:
		details, err := json.MarshalIndent(user, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal SQLServer Flex user: %w", err)
		}
		p.Outputln(string(details))

		return nil
	case print.YAMLOutputFormat:
		details, err := yaml.MarshalWithOptions(user, yaml.IndentSequence(true), yaml.UseJSONMarshaler())
		if err != nil {
			return fmt.Errorf("marshal SQLServer Flex user: %w", err)
		}
		p.Outputln(string(details))

		return nil
	default:
		p.Outputf("Created user for instance %q. User ID: %s\n\n", instanceLabel, utils.PtrString(user.Id))
		p.Outputf("Username: %s\n", utils.PtrString(user.Username))
		p.Outputf("Password: %s\n", utils.PtrString(user.Password))
		if user.Roles != nil && len(*user.Roles) != 0 {
			p.Outputf("Roles: [%v]\n", strings.Join(*user.Roles, ", "))
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
