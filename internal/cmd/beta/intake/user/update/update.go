package update

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/intake"
	"github.com/stackitcloud/stackit-sdk-go/services/intake/wait"

	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	cliErr "github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/intake/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/spinner"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

const (
	userIdArg       = "USER_ID"
	intakeIdFlag    = "intake-id"
	displayNameFlag = "display-name"
	passwordFlag    = "password"
	descriptionFlag = "description"
	typeFlag        = "type"
	labelFlag       = "labels"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	IntakeId    string
	UserId      string
	DisplayName *string
	Password    *string
	Description *string
	Type        *string
	Labels      *map[string]string
}

func NewCmd(p *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("update %s", userIdArg),
		Short: "Updates an Intake User",
		Long:  "Updates an Intake User. Only the specified fields are updated.",
		Args:  args.SingleArg(userIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Update the display name of an Intake User with ID "xxx"`,
				`$ stackit beta intake user update xxx --intake-id yyy --display-name "new-user-name"`),
			examples.NewExample(
				`Update the password of an Intake User`,
				`$ stackit beta intake user update xxx --intake-id yyy --password "new-secret"`),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			model, err := parseInput(p.Printer, cmd, args)
			if err != nil {
				return err
			}

			// Configure API client
			apiClient, err := client.ConfigureClient(p.Printer, p.CliVersion)
			if err != nil {
				return err
			}

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to update Intake User %q?", model.UserId)
				err = p.Printer.PromptForConfirmation(prompt)
				if err != nil {
					return err
				}
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("update Intake User: %w", err)
			}

			// Wait for async operation, if async mode not enabled
			if !model.Async {
				s := spinner.New(p.Printer)
				s.Start("Updating STACKIT Intake User instance")
				_, err = wait.CreateOrUpdateIntakeUserWaitHandler(ctx, apiClient, model.ProjectId, model.Region, model.IntakeId, model.UserId).WaitWithContext(ctx)
				if err != nil {
					return fmt.Errorf("wait for STACKIT Instance creation: %w", err)
				}
				s.Stop()
			}

			return outputResult(p.Printer, model, model.IntakeId, resp)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().String(intakeIdFlag, "", "ID of the Intake")
	cmd.Flags().String(displayNameFlag, "", "Display name")
	cmd.Flags().String(passwordFlag, "", "User password")
	cmd.Flags().String(descriptionFlag, "", "Description")
	cmd.Flags().String(typeFlag, "", "Type of user, 'intake' for writing or 'dead-letter' for reading from the dead-letter queue")
	cmd.Flags().StringToString(labelFlag, nil, `Labels in key=value format, separated by commas. Example: --labels "key1=value1,key2=value2".`)

	err := flags.MarkFlagsRequired(cmd, intakeIdFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	userId := inputArgs[0]

	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &cliErr.ProjectIdError{}
	}

	intakeId := flags.FlagToStringValue(p, cmd, intakeIdFlag)
	if intakeId == "" {
		return nil, &cliErr.FlagValidationError{
			Flag:    intakeIdFlag,
			Details: "can't be empty",
		}
	}

	model := &inputModel{
		GlobalFlagModel: globalFlags,
		IntakeId:        intakeId,
		UserId:          userId,
		DisplayName:     flags.FlagToStringPointer(p, cmd, displayNameFlag),
		Password:        flags.FlagToStringPointer(p, cmd, passwordFlag),
		Description:     flags.FlagToStringPointer(p, cmd, descriptionFlag),
		Type:            flags.FlagToStringPointer(p, cmd, typeFlag),
		Labels:          flags.FlagToStringToStringPointer(p, cmd, labelFlag),
	}

	if model.DisplayName == nil && model.Password == nil && model.Description == nil && model.Type == nil && model.Labels == nil {
		return nil, &cliErr.EmptyUpdateError{}
	}

	p.DebugInputModel(model)
	return model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *intake.APIClient) intake.ApiUpdateIntakeUserRequest {
	req := apiClient.UpdateIntakeUser(ctx, model.ProjectId, model.Region, model.IntakeId, model.UserId)

	payload := intake.UpdateIntakeUserPayload{}
	if model.DisplayName != nil {
		payload.DisplayName = model.DisplayName
	}
	if model.Password != nil {
		payload.Password = model.Password
	}
	if model.Description != nil {
		payload.Description = model.Description
	}
	if model.Type != nil {
		// This line is only reached if model.Type is not nil. Therefore, the conversion is safe.
		payload.Type = (*intake.UserType)(model.Type)
	}
	if model.Labels != nil {
		payload.Labels = model.Labels
	}

	req = req.UpdateIntakeUserPayload(payload)
	return req
}

func outputResult(p *print.Printer, model *inputModel, intakeId string, resp *intake.IntakeUserResponse) error {
	return p.OutputResult(model.OutputFormat, resp, func() error {
		if resp == nil {
			p.Outputf("Updated Intake User for Intake %q, but no intake ID was returned.\n", intakeId)
			return nil
		}

		operationState := "Updated"
		if model.Async {
			operationState = "Triggered update of"
		}
		p.Outputf("%s Intake User for Intake %q. User ID: %s\n", operationState, intakeId, utils.PtrString(resp.Id))
		return nil
	})
}
