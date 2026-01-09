package resetpassword

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/types"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/mongodbflex/client"
	mongodbflexUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/mongodbflex/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/mongodbflex"
)

const (
	userIdArg = "USER_ID"

	instanceIdFlag = "instance-id"
)

type inputModel struct {
	*globalflags.GlobalFlagModel

	InstanceId string
	UserId     string
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("reset-password %s", userIdArg),
		Short: "Resets the password of a MongoDB Flex user",
		Long: fmt.Sprintf("%s\ns%s",
			"Resets the password of a MongoDB Flex user.",
			"The new password is visible after and cannot be retrieved later.",
		),
		Example: examples.Build(
			examples.NewExample(
				`Reset the password of a MongoDB Flex user with ID "xxx" of instance with ID "yyy"`,
				"$ stackit mongodbflex user reset-password xxx --instance-id yyy"),
		),
		Args: args.SingleArg(userIdArg, utils.ValidateUUID),
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

			instanceLabel, err := mongodbflexUtils.GetInstanceName(ctx, apiClient, model.ProjectId, model.InstanceId, model.Region)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get instance name: %v", err)
				instanceLabel = model.InstanceId
			}

			userLabel, err := mongodbflexUtils.GetUserName(ctx, apiClient, model.ProjectId, model.InstanceId, model.UserId, model.Region)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get user name: %v", err)
				userLabel = model.UserId
			}

			prompt := fmt.Sprintf("Are you sure you want to reset the password of user %q of instance %q? (This cannot be undone)", userLabel, instanceLabel)
			err = params.Printer.PromptForConfirmation(prompt)
			if err != nil {
				return err
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			user, err := req.Execute()
			if err != nil {
				return fmt.Errorf("reset MongoDB Flex user password: %w", err)
			}

			return outputResult(params.Printer, model.OutputFormat, userLabel, instanceLabel, user)
		},
	}

	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), instanceIdFlag, "ID of the instance")

	err := flags.MarkFlagsRequired(cmd, instanceIdFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	userId := inputArgs[0]

	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		InstanceId:      flags.FlagToStringValue(p, cmd, instanceIdFlag),
		UserId:          userId,
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *mongodbflex.APIClient) mongodbflex.ApiResetUserRequest {
	req := apiClient.ResetUser(ctx, model.ProjectId, model.InstanceId, model.UserId, model.Region)
	return req
}

func outputResult(p *print.Printer, outputFormat, userLabel, instanceLabel string, user *mongodbflex.User) error {
	if user == nil {
		return fmt.Errorf("user is nil")
	}

	return p.OutputResult(outputFormat, user, func() error {
		p.Outputf("Reset password for user %q of instance %q\n\n", userLabel, instanceLabel)
		p.Outputf("Username: %s\n", utils.PtrString(user.Username))
		p.Outputf("New password: %s\n", utils.PtrString(user.Password))
		p.Outputf("New URI: %s\n", utils.PtrString(user.Uri))
		return nil
	})
}
