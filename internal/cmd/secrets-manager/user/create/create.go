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
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/secrets-manager/client"
	secretsManagerUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/secrets-manager/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/secretsmanager"
)

const (
	instanceIdFlag   = "instance-id"
	descriptionFlag  = "description"
	writeFlag        = "write"
	hidePasswordFlag = "hide-password"
)

type inputModel struct {
	*globalflags.GlobalFlagModel

	InstanceId   string
	Description  *string
	Write        *bool
	HidePassword bool
}

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Creates a Secrets Manager user",
		Long:  "Creates a user for a Secrets Manager instance with generated username and password",
		Example: examples.Build(
			examples.NewExample(
				`Create a Secrets Manager user for instance with ID "xxx"`,
				"$ stackit mongodbflex user create --instance-id xxx"),
			examples.NewExample(
				`Create a Secrets Manager user for instance with ID "xxx" and description "yyy"`,
				"$ stackit mongodbflex user create --instance-id xxx --description yyy"),
			examples.NewExample(
				`Create a Secrets Manager user for instance with ID "xxx" and doesn't display the password`,
				"$ stackit mongodbflex user create --instance-id xxx --hide-password"),
			examples.NewExample(
				`Create a Secrets Manager user for instance with ID "xxx" with write access to the secrets engine`,
				"$ stackit mongodbflex user create --instance-id xxx --write"),
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

			instanceLabel, err := secretsManagerUtils.GetInstanceName(ctx, apiClient, model.ProjectId, model.InstanceId)
			if err != nil {
				instanceLabel = model.InstanceId
			}

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to create a user for instance %q?", instanceLabel)
				err = confirm.PromptForConfirmation(cmd, prompt)
				if err != nil {
					return err
				}
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("create Secrets Manager user: %w", err)
			}

			cmd.Printf("Created user for instance %q. User ID: %s\n\n", instanceLabel, *resp.Id)
			cmd.Printf("Username: %s\n", *resp.Username)
			if model.HidePassword {
				cmd.Printf("Password: <hidden>\n")
			} else {
				cmd.Printf("Password: %s\n", *resp.Password)
			}
			cmd.Printf("Description: %s\n", *resp.Description)
			cmd.Printf("Write Access: %t\n", *resp.Write)

			return nil
		},
	}

	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), instanceIdFlag, "ID of the instance")
	cmd.Flags().String(descriptionFlag, "", "A user chosen description to differentiate between multiple users")
	cmd.Flags().Bool(writeFlag, false, "User write access to the secrets engine. If unset, user is read-only")
	cmd.Flags().Bool(hidePasswordFlag, false, "Hide password in output")

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
		Description:     utils.Ptr(flags.FlagToStringValue(cmd, descriptionFlag)),
		Write:           utils.Ptr(flags.FlagToBoolValue(cmd, writeFlag)),
		HidePassword:    flags.FlagToBoolValue(cmd, hidePasswordFlag),
	}, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *secretsmanager.APIClient) secretsmanager.ApiCreateUserRequest {
	req := apiClient.CreateUser(ctx, model.ProjectId, model.InstanceId)
	req = req.CreateUserPayload(secretsmanager.CreateUserPayload{
		Description: model.Description,
		Write:       model.Write,
	})
	return req
}
