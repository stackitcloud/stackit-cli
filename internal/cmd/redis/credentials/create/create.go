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
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/redis/client"
	redisUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/redis/utils"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/redis"
)

const (
	instanceIdFlag   = "instance-id"
	hidePasswordFlag = "hide-password"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	InstanceId   string
	HidePassword bool
}

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Creates credentials for a Redis instance",
		Long:  "Creates credentials (username and password) for a Redis instance.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Create credentials for a Redis instance`,
				"$ stackit redis credentials create --instance-id xxx"),
			examples.NewExample(
				`Create credentials for a Redis instance and hide the password in the output`,
				"$ stackit redis credentials create --instance-id xxx --hide-password"),
		),
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

			instanceLabel, err := redisUtils.GetInstanceName(ctx, apiClient, model.ProjectId, model.InstanceId)
			if err != nil {
				instanceLabel = model.InstanceId
			}

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to create credentials for instance %s?", instanceLabel)
				err = confirm.PromptForConfirmation(cmd, prompt)
				if err != nil {
					return err
				}
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("create Redis credentials: %w", err)
			}

			cmd.Printf("Created credentials for instance %s. Credentials ID: %s\n\n", instanceLabel, *resp.Id)
			cmd.Printf("Username: %s\n", *resp.Raw.Credentials.Username)
			if model.HidePassword {
				cmd.Printf("Password: <hidden>\n")
			} else {
				cmd.Printf("Password: %s\n", *resp.Raw.Credentials.Password)
			}
			cmd.Printf("Host: %s\n", *resp.Raw.Credentials.Host)
			cmd.Printf("Port: %d\n", *resp.Raw.Credentials.Port)
			cmd.Printf("URI: %s\n", *resp.Uri)
			return nil
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), instanceIdFlag, "Instance ID")
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
		HidePassword:    flags.FlagToBoolValue(cmd, hidePasswordFlag),
	}, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *redis.APIClient) redis.ApiCreateCredentialsRequest {
	req := apiClient.CreateCredentials(ctx, model.ProjectId, model.InstanceId)
	return req
}
