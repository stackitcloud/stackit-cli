package create

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
	instanceIdFlag  = "instance-id"
	descriptionFlag = "description"
	writeFlag       = "write"
)

type inputModel struct {
	*globalflags.GlobalFlagModel

	InstanceId  string
	Description *string
	Write       *bool
}

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Creates a Secrets Manager user",
		Long: fmt.Sprintf("%s\n%s\n%s",
			"Creates a Secrets Manager user.",
			"The username and password are auto-generated and provided upon creation. The password cannot be retrieved later.",
			"A description can be provided to identify a user.",
		),
		Example: examples.Build(
			examples.NewExample(
				`Create a Secrets Manager user for instance with ID "xxx" and description "yyy"`,
				"$ stackit secrets-manager user create --instance-id xxx --description yyy"),
			examples.NewExample(
				`Create a Secrets Manager user for instance with ID "xxx" with write access to the secrets engine`,
				"$ stackit secrets-manager user create --instance-id xxx --write"),
		),
		Args: args.NoArgs,
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

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to create a user for instance %q?", instanceLabel)
				err = params.Printer.PromptForConfirmation(prompt)
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

			return outputResult(params.Printer, model.OutputFormat, instanceLabel, resp)
		},
	}

	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), instanceIdFlag, "ID of the instance")
	cmd.Flags().String(descriptionFlag, "", "A user chosen description to differentiate between multiple users")
	cmd.Flags().Bool(writeFlag, false, "User write access to the secrets engine. If unset, user is read-only")

	err := flags.MarkFlagsRequired(cmd, instanceIdFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command, _ []string) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		InstanceId:      flags.FlagToStringValue(p, cmd, instanceIdFlag),
		Description:     utils.Ptr(flags.FlagToStringValue(p, cmd, descriptionFlag)),
		Write:           utils.Ptr(flags.FlagToBoolValue(p, cmd, writeFlag)),
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *secretsmanager.APIClient) secretsmanager.ApiCreateUserRequest {
	req := apiClient.CreateUser(ctx, model.ProjectId, model.InstanceId)
	req = req.CreateUserPayload(secretsmanager.CreateUserPayload{
		Description: model.Description,
		Write:       model.Write,
	})
	return req
}

func outputResult(p *print.Printer, outputFormat, instanceLabel string, user *secretsmanager.User) error {
	if user == nil {
		return fmt.Errorf("user is nil")
	}

	return p.OutputResult(outputFormat, user, func() error {
		p.Outputf("Created user for instance %q. User ID: %s\n\n", instanceLabel, utils.PtrString(user.Id))
		p.Outputf("Username: %s\n", utils.PtrString(user.Username))
		p.Outputf("Password: %s\n", utils.PtrString(user.Password))
		p.Outputf("Description: %s\n", utils.PtrString(user.Description))
		p.Outputf("Write Access: %s\n", utils.PtrString(user.Write))

		return nil
	})
}
