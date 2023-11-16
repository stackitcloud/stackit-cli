package update

import (
	"context"
	"fmt"

	"stackit/internal/pkg/args"
	"stackit/internal/pkg/confirm"
	"stackit/internal/pkg/errors"
	"stackit/internal/pkg/examples"
	"stackit/internal/pkg/flags"
	"stackit/internal/pkg/globalflags"
	"stackit/internal/pkg/services/mongodbflex/client"
	mongodbflexUtils "stackit/internal/pkg/services/mongodbflex/utils"
	"stackit/internal/pkg/utils"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/mongodbflex"
)

const (
	userIdArg = "USER_ID"

	instanceIdFlag = "instance-id"
	databaseFlag   = "database"
	rolesFlag      = "roles"
)

type inputModel struct {
	*globalflags.GlobalFlagModel

	InstanceId string
	UserId     string
	Database   *string
	Roles      *[]string
}

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("update %s", userIdArg),
		Short: "Update a MongoDB Flex user",
		Long:  "Update a MongoDB Flex user.",
		Example: examples.Build(
			examples.NewExample(
				`Update the roles of a MongoDB Flex user with ID "xxx" of instance with ID "yyy"`,
				"$ stackit mongodbflex user update xxx --instance-id yyy --roles read"),
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

			instanceLabel, err := mongodbflexUtils.GetInstanceName(ctx, apiClient, model.ProjectId, model.InstanceId)
			if err != nil {
				instanceLabel = model.InstanceId
			}

			userLabel, err := mongodbflexUtils.GetUserName(ctx, apiClient, model.ProjectId, model.InstanceId, model.UserId)
			if err != nil {
				userLabel = model.UserId
			}

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to update user %s of instance %s?", userLabel, instanceLabel)
				err = confirm.PromptForConfirmation(cmd, prompt)
				if err != nil {
					return err
				}
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			err = req.Execute()
			if err != nil {
				return fmt.Errorf("update MongoDB Flex user: %w", err)
			}

			cmd.Printf("Updated user %s of instance %s\n", userLabel, instanceLabel)
			return nil
		},
	}

	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	rolesOptions := []string{"read", "readWrite"}

	cmd.Flags().Var(flags.UUIDFlag(), instanceIdFlag, "ID of the instance")
	cmd.Flags().String(databaseFlag, "", "The database inside the MongoDB instance that the user has access to. If it does not exist, it will be created once the user writes to it")
	cmd.Flags().Var(flags.EnumSliceFlag(false, nil, rolesOptions...), rolesFlag, fmt.Sprintf("Roles of the user, possible values are %q", rolesOptions))

	err := flags.MarkFlagsRequired(cmd, instanceIdFlag)
	cobra.CheckErr(err)
}

func parseInput(cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	userId := inputArgs[0]

	globalFlags := globalflags.Parse(cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	database := flags.FlagToStringPointer(cmd, databaseFlag)
	roles := flags.FlagToStringSlicePointer(cmd, rolesFlag)

	if database == nil && roles == nil {
		return nil, &errors.EmptyUpdateError{}
	}

	return &inputModel{
		GlobalFlagModel: globalFlags,
		InstanceId:      flags.FlagToStringValue(cmd, instanceIdFlag),
		UserId:          userId,
		Database:        database,
		Roles:           roles,
	}, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *mongodbflex.APIClient) mongodbflex.ApiPartialUpdateUserRequest {
	req := apiClient.PartialUpdateUser(ctx, model.ProjectId, model.InstanceId, model.UserId)
	req = req.PartialUpdateUserPayload(mongodbflex.PartialUpdateUserPayload{
		Database: model.Database,
		Roles:    model.Roles,
	})
	return req
}
