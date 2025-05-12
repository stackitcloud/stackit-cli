package update

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
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

func NewCmd(params *params.CmdParams) *cobra.Command {
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
			model, err := parseInput(params.Printer, cmd, args)
			if err != nil {
				return err
			}

			// Configure API client
			apiClient, err := client.ConfigureClient(params.Printer, params.CliVersion)
			if err != nil {
				return err
			}

			instanceLabel, err := secretsManagerUtils.GetInstanceName(ctx, apiClient, model.ProjectId, model.InstanceId)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get instance name: %v", err)
				instanceLabel = model.InstanceId
			}

			userLabel, err := secretsManagerUtils.GetUserLabel(ctx, apiClient, model.ProjectId, model.InstanceId, model.UserId)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get user label: %v", err)
				userLabel = fmt.Sprintf("%q", model.UserId)
			}

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to update user %s of instance %q?", userLabel, instanceLabel)
				err = params.Printer.PromptForConfirmation(prompt)
				if err != nil {
					return err
				}
			}

			// Call API
			req, err := buildRequest(ctx, model, apiClient)
			if err != nil {
				return err
			}

			err = req.Execute()
			if err != nil {
				return fmt.Errorf("update Secrets Manager user: %w", err)
			}

			params.Printer.Info("Updated user %s of instance %q\n", userLabel, instanceLabel)
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

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	userId := inputArgs[0]

	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	enableWrite := flags.FlagToBoolPointer(p, cmd, enableWriteFlag)
	disableWrite := flags.FlagToBoolPointer(p, cmd, disableWriteFlag)

	if enableWrite == nil && disableWrite == nil {
		return nil, &errors.EmptyUpdateError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		InstanceId:      flags.FlagToStringValue(p, cmd, instanceIdFlag),
		EnableWrite:     enableWrite,
		DisableWrite:    disableWrite,
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

func buildRequest(ctx context.Context, model *inputModel, apiClient *secretsmanager.APIClient) (secretsmanager.ApiUpdateUserRequest, error) {
	req := apiClient.UpdateUser(ctx, model.ProjectId, model.InstanceId, model.UserId)

	var write bool

	if model.EnableWrite != nil && model.DisableWrite == nil {
		write = true
	} else if model.DisableWrite != nil && model.EnableWrite == nil {
		write = false
	} else if model.DisableWrite == nil && model.EnableWrite == nil {
		// Should never happen
		return req, fmt.Errorf("one of %s and %s flags needs to be set", enableWriteFlag, disableWriteFlag)
	} else if model.DisableWrite != nil && model.EnableWrite != nil {
		// Should never happen
		return req, fmt.Errorf("%s and %s flags can't be both set", enableWriteFlag, disableWriteFlag)
	}

	req = req.UpdateUserPayload(secretsmanager.UpdateUserPayload{
		Write: utils.Ptr(write),
	})
	return req, nil
}
