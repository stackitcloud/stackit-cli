package resetpassword

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
	userIdArg = "USER_ID"

	instanceIdFlag = "instance-id"
)

type inputModel struct {
	*globalflags.GlobalFlagModel

	InstanceId string
	UserId     string
}

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("reset-password %s", userIdArg),
		Short: "Resets the password of a PostgreSQL Flex user",
		Long: fmt.Sprintf("%s\ns%s",
			"Resets the password of a PostgreSQL Flex user.",
			"The new password is visible after and cannot be retrieved later.",
		),
		Example: examples.Build(
			examples.NewExample(
				`Reset the password of a PostgreSQL Flex user with ID "xxx" of instance with ID "yyy"`,
				"$ stackit postgresflex user reset-password xxx --instance-id yyy"),
		),
		Args: args.SingleArg(userIdArg, nil),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			model, err := parseInput(cmd, args)
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

			userLabel, err := postgresflexUtils.GetUserName(ctx, apiClient, model.ProjectId, model.InstanceId, model.UserId)
			if err != nil {
				userLabel = model.UserId
			}

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to reset the password of user %q of instance %q? (This cannot be undone)", userLabel, instanceLabel)
				err = confirm.PromptForConfirmation(cmd, prompt)
				if err != nil {
					return err
				}
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			user, err := req.Execute()
			if err != nil {
				return fmt.Errorf("reset PostgreSQL Flex user password: %w", err)
			}

			cmd.Printf("Reset password for user %q of instance %q\n\n", userLabel, instanceLabel)
			cmd.Printf("Username: %s\n", *user.Item.Username)
			cmd.Printf("New password: %s\n", *user.Item.Password)
			cmd.Printf("New URI: %s\n", *user.Item.Uri)
			return nil
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

func parseInput(cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	userId := inputArgs[0]

	globalFlags := globalflags.Parse(cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	return &inputModel{
		GlobalFlagModel: globalFlags,
		InstanceId:      flags.FlagToStringValue(cmd, instanceIdFlag),
		UserId:          userId,
	}, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *postgresflex.APIClient) postgresflex.ApiResetUserRequest {
	req := apiClient.ResetUser(ctx, model.ProjectId, model.InstanceId, model.UserId)
	return req
}
