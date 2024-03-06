package delete

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/confirm"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/secrets-manager/client"
	secretsManagerUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/secrets-manager/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/secretsmanager"
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
		Use:   fmt.Sprintf("delete %s", userIdArg),
		Short: "Deletes a Secrets Manager user",
		Long: fmt.Sprintf("%s\n%s",
			"Deletes a Secrets Manager user by ID. You can get the IDs of users for an instance by running:",
			"  $ stackit secrets-manager user list --instance-id <INSTANCE_ID>",
		),
		Example: examples.Build(
			examples.NewExample(
				`Delete a Secrets Manager user with ID "xxx" for instance with ID "yyy"`,
				"$ stackit secrets-manager user delete xxx --instance-id yyy"),
		),
		Args: args.SingleArg(userIdArg, utils.ValidateUUID),
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

			instanceLabel, err := secretsManagerUtils.GetInstanceName(ctx, apiClient, model.ProjectId, model.InstanceId)
			if err != nil {
				instanceLabel = model.InstanceId
			}

			userLabel, userDescription, err := secretsManagerUtils.GetUserDetails(ctx, apiClient, model.ProjectId, model.InstanceId, model.UserId)
			if err != nil {
				userLabel = model.UserId
			}

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to delete user %q (%s) of instance %q? (This cannot be undone)", userLabel, userDescription, instanceLabel)
				err = confirm.PromptForConfirmation(cmd, prompt)
				if err != nil {
					return err
				}
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			err = req.Execute()
			if err != nil {
				return fmt.Errorf("delete Secrets Manager user: %w", err)
			}

			cmd.Printf("Deleted user %q of instance %q\n", userLabel, instanceLabel)
			return nil
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), instanceIdFlag, "Instance ID")

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

func buildRequest(ctx context.Context, model *inputModel, apiClient *secretsmanager.APIClient) secretsmanager.ApiDeleteUserRequest {
	req := apiClient.DeleteUser(ctx, model.ProjectId, model.InstanceId, model.UserId)
	return req
}
