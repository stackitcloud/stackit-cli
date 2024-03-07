package update

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

	instanceIdFlag   = "instance-id"
	enableWriteFlag  = "enable-write"
	disableWriteFlag = "disable-write"
)

type inputModel struct {
	*globalflags.GlobalFlagModel

	InstanceId   string
	UserId       string
	EnableWrite  *bool
	DisableWrite *bool
}

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("update %s", userIdArg),
		Short: "Updates the write privileges Secrets Manager user",
		Long:  "Updates the write privileges Secrets Manager user.",
		Example: examples.Build(
			examples.NewExample(
				`Enable write access of a Secrets Manager user with ID "xxx" of instance with ID "yyy"`,
				"$ stackit secrets-manager user update xxx --instance-id yyy --enable-write"),
			examples.NewExample(
				`Disable write access of a Secrets Manager user with ID "xxx" of instance with ID "yyy"`,
				"$ stackit secrets-manager user update xxx --instance-id yyy --disable-write"),
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

			var userLabel string

			userName, userDescription, err := secretsManagerUtils.GetUserDetails(ctx, apiClient, model.ProjectId, model.InstanceId, model.UserId)
			if err != nil {
				userLabel = fmt.Sprintf("%q", model.UserId)
			} else {
				userLabel = fmt.Sprintf("%q (%q)", userName, userDescription)
			}

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to update user %s of instance %q?", userLabel, instanceLabel)
				err = confirm.PromptForConfirmation(cmd, prompt)
				if err != nil {
					return err
				}
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			err = req.Execute()
			if err != nil {
				return fmt.Errorf("update Secrets Manager user: %w", err)
			}

			cmd.Printf("Updated user %s of instance %q\n", userLabel, instanceLabel)
			return nil
		},
	}

	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), instanceIdFlag, "ID of the instance")
	cmd.Flags().Bool(enableWriteFlag, false, "Set the user to have write access.")
	cmd.Flags().Bool(disableWriteFlag, false, "Set the user to have read-only access.")

	err := flags.MarkFlagsRequired(cmd, instanceIdFlag)

	cmd.MarkFlagsMutuallyExclusive(enableWriteFlag, disableWriteFlag)
	cobra.CheckErr(err)
}

func parseInput(cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	userId := inputArgs[0]

	globalFlags := globalflags.Parse(cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	enableWrite := flags.FlagToBoolPointer(cmd, enableWriteFlag)
	disableWrite := flags.FlagToBoolPointer(cmd, disableWriteFlag)

	if enableWrite == nil && disableWrite == nil {
		return nil, &errors.EmptyUpdateError{}
	}

	return &inputModel{
		GlobalFlagModel: globalFlags,
		InstanceId:      flags.FlagToStringValue(cmd, instanceIdFlag),
		EnableWrite:     enableWrite,
		DisableWrite:    disableWrite,
		UserId:          userId,
	}, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *secretsmanager.APIClient) secretsmanager.ApiUpdateUserRequest {
	req := apiClient.UpdateUser(ctx, model.ProjectId, model.InstanceId, model.UserId)

	// model.EnableWrite and model.DisableWrite are mutually exclusive and can't be both nil
	// therefore we can check only one for the value of the write parameter
	write := model.EnableWrite != nil

	req = req.UpdateUserPayload(secretsmanager.UpdateUserPayload{
		Write: utils.Ptr(write),
	})
	return req
}
