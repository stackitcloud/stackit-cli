package create

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/types"

	"github.com/spf13/cobra"
	postgresflex "github.com/stackitcloud/stackit-sdk-go/services/postgresflex/v3api"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/postgresflex/client"
	postgresflexUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/postgresflex/utils"
)

const (
	instanceIdFlag = "instance-id"
	usernameFlag   = "username"
)

var (
	rolesDefault = []string{"login"}
	roleFlag     = flags.StringEnumSliceFlag(
		"role",
		[]string{"login", "createdb"},
		"Roles of the user,",
		flags.DefaultValues(rolesDefault...),
	)
)

type inputModel struct {
	*globalflags.GlobalFlagModel

	InstanceId string
	Username   string
	Roles      []string
}

func NewCmd(params *types.CmdParams) *cobra.Command {
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
			model, err := parseInput(params.Printer, cmd, args)
			if err != nil {
				return err
			}

			// Configure API client
			apiClient, err := client.ConfigureClient(params.Printer, params.CliVersion)
			if err != nil {
				return err
			}

			instanceLabel, err := postgresflexUtils.GetInstanceName(ctx, apiClient.DefaultAPI, model.ProjectId, model.Region, model.InstanceId)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get instance name: %v", err)
				instanceLabel = model.InstanceId
			}

			prompt := fmt.Sprintf("Are you sure you want to create a user for instance %q?", instanceLabel)
			err = params.Printer.PromptForConfirmation(prompt)
			if err != nil {
				return err
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("create PostgreSQL Flex user: %w", err)
			}

			return outputResult(params.Printer, model.OutputFormat, instanceLabel, resp)
		},
	}

	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), instanceIdFlag, "ID of the instance")
	cmd.Flags().String(usernameFlag, "", "Username of the user")
	roleFlag.Register(cmd)

	err := flags.MarkFlagsRequired(cmd, instanceIdFlag, usernameFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command, _ []string) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		InstanceId:      flags.FlagToStringValue(p, cmd, instanceIdFlag),
		Username:        flags.FlagToStringValue(p, cmd, usernameFlag),
		Roles:           roleFlag.Get(),
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *postgresflex.APIClient) postgresflex.ApiCreateUserRequest {
	req := apiClient.DefaultAPI.CreateUser(ctx, model.ProjectId, model.Region, model.InstanceId)
	req = req.CreateUserPayload(postgresflex.CreateUserPayload{
		Name:  model.Username,
		Roles: model.Roles,
	})
	return req
}

func outputResult(p *print.Printer, outputFormat, instanceLabel string, user *postgresflex.CreateUserResponse) error {
	return p.OutputResult(outputFormat, user, func() error {
		if user == nil {
			return fmt.Errorf("no response passed")
		}

		p.Outputf("Created user for instance %q. User ID: %d\n\n", instanceLabel, user.Id)
		p.Outputf("Username: %s\n", user.Name)
		p.Outputf("Password: %s\n", user.Password)

		return nil
	})
}
