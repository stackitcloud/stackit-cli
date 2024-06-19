package resetpassword

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
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/sqlserverflex/client"
	sqlserverflexUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/sqlserverflex/utils"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/sqlserverflex"
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

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("reset-password %s", userIdArg),
		Short: "Resets the password of an SQLServer Flex user",
		Long: fmt.Sprintf("%s\ns%s",
			"Resets the password of an SQLServer Flex user.",
			"The new password is visible after and cannot be retrieved later.",
		),
		Example: examples.Build(
			examples.NewExample(
				`Reset the password of an SQLServer Flex user with ID "xxx" of instance with ID "yyy"`,
				"$ stackit sqlserverflex user reset-password xxx --instance-id yyy"),
		),
		Args: args.SingleArg(userIdArg, nil),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			model, err := parseInput(p, cmd, args)
			if err != nil {
				return err
			}

			// Configure API client
			apiClient, err := client.ConfigureClient(p)
			if err != nil {
				return err
			}

			instanceLabel, err := sqlserverflexUtils.GetInstanceName(ctx, apiClient, model.ProjectId, model.InstanceId)
			if err != nil {
				p.Debug(print.ErrorLevel, "get instance name: %v", err)
				instanceLabel = model.InstanceId
			}

			userLabel, err := sqlserverflexUtils.GetUserName(ctx, apiClient, model.ProjectId, model.InstanceId, model.UserId)
			if err != nil {
				p.Debug(print.ErrorLevel, "get user name: %v", err)
				userLabel = model.UserId
			}

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to reset the password of user %q of instance %q? (This cannot be undone)", userLabel, instanceLabel)
				err = p.PromptForConfirmation(prompt)
				if err != nil {
					return err
				}
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			user, err := req.Execute()
			if err != nil {
				return fmt.Errorf("reset SQLServer Flex user password: %w", err)
			}

			return outputResult(p, model, userLabel, instanceLabel, user.Item)
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

func buildRequest(ctx context.Context, model *inputModel, apiClient *sqlserverflex.APIClient) sqlserverflex.ApiResetUserRequest {
	req := apiClient.ResetUser(ctx, model.ProjectId, model.InstanceId, model.UserId)
	return req
}

func outputResult(p *print.Printer, model *inputModel, userLabel, instanceLabel string, user *sqlserverflex.User) error {
	switch model.OutputFormat {
	case print.JSONOutputFormat:
		details, err := json.MarshalIndent(user, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal SQLServer Flex reset password: %w", err)
		}
		p.Outputln(string(details))

		return nil
	case print.YAMLOutputFormat:
		details, err := yaml.MarshalWithOptions(user, yaml.IndentSequence(true))
		if err != nil {
			return fmt.Errorf("marshal SQLServer Flex reset password: %w", err)
		}
		p.Outputln(string(details))

		return nil
	default:
		p.Outputf("Reset password for user %q of instance %q\n\n", userLabel, instanceLabel)
		p.Outputf("Username: %s\n", *user.Username)
		p.Outputf("New password: %s\n", *user.Password)
		if user.Uri != nil && *user.Uri != "" {
			p.Outputf("New URI: %s\n", *user.Uri)
		}
		return nil
	}
}
