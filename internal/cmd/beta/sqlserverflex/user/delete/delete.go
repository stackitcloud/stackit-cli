package delete

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

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("delete %s", userIdArg),
		Short: "Deletes a SQLServer Flex user",
		Long: fmt.Sprintf("%s\n%s",
			"Deletes a SQLServer Flex user by ID. You can get the IDs of users for an instance by running:",
			"  $ stackit beta sqlserverflex user list --instance-id <INSTANCE_ID>",
		),
		Example: examples.Build(
			examples.NewExample(
				`Delete a SQLServer Flex user with ID "xxx" for instance with ID "yyy"`,
				"$ stackit beta sqlserverflex user delete xxx --instance-id yyy"),
		),
		Args: args.SingleArg(userIdArg, nil),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			model, err := parseInput(params.Printer, cmd, args)
			if err != nil {
				return err
			}

			// Configure API client
			apiClient, err := client.ConfigureClient(params.Printer)
			if err != nil {
				return err
			}

			instanceLabel, err := sqlserverflexUtils.GetInstanceName(ctx, apiClient, model.ProjectId, model.InstanceId, model.Region)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get instance name: %v", err)
				instanceLabel = model.InstanceId
			}

			userLabel, err := sqlserverflexUtils.GetUserName(ctx, apiClient, model.ProjectId, model.InstanceId, model.UserId, model.Region)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get user name: %v", err)
				userLabel = model.UserId
			}

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to delete user %q of instance %q? (This cannot be undone)", userLabel, instanceLabel)
				err = params.Printer.PromptForConfirmation(prompt)
				if err != nil {
					return err
				}
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			err = req.Execute()
			if err != nil {
				return fmt.Errorf("delete SQLServer Flex user: %w", err)
			}

			params.Printer.Info("Deleted user %q of instance %q\n", userLabel, instanceLabel)
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

func buildRequest(ctx context.Context, model *inputModel, apiClient *sqlserverflex.APIClient) sqlserverflex.ApiDeleteUserRequest {
	req := apiClient.DeleteUser(ctx, model.ProjectId, model.InstanceId, model.UserId, model.Region)
	return req
}
