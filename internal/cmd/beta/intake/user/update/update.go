package update

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	cliErr "github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/intake/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/spinner"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/intake"
	"github.com/stackitcloud/stackit-sdk-go/services/intake/wait"
)

const (
	userIdArg = "USER_ID"

	intakeIdFlag    = "intake-id"
	displayNameFlag = "display-name"
	descriptionFlag = "description"
	passwordFlag    = "password"
	userTypeFlag    = "type"
	labelsFlag      = "labels"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	IntakeId    string
	UserId      string
	DisplayName *string
	Description *string
	Password    *string
	UserType    *string
	Labels      *map[string]string
}

func NewCmd(p *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("update %s", userIdArg),
		Short: "Updates an Intake User",
		Long:  "Updates an Intake User. Only the specified fields are updated.",
		Args:  args.SingleArg(userIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Update the display name of an Intake User`,
				`$ stackit beta intake user update xxx --intake-id yyy --display-name "new-user-name"`),
			examples.NewExample(
				`Update the password and description for an Intake User`,
				`$ stackit beta intake user update xxx --intake-id yyy --password "NewSecret123\!" --description "Updated description"`),
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

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("update Intake User: %w", err)
			}

			// Wait for async operation, if async mode not enabled
			if !model.Async {
				s := spinner.New(p.Printer)
				s.Start("Updating STACKIT Intake User")
				_, err = wait.CreateOrUpdateIntakeUserWaitHandler(ctx, apiClient, model.ProjectId, model.Region, model.IntakeId, model.UserId).WaitWithContext(ctx)
				if err != nil {
					return fmt.Errorf("wait for STACKIT Intake User update: %w", err)
				}
				s.Stop()
			}

			return outputResult(p.Printer, model, resp)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), intakeIdFlag, "Intake ID")
	cmd.Flags().String(displayNameFlag, "", "Display name")
	cmd.Flags().String(descriptionFlag, "", "Description")
	cmd.Flags().String(passwordFlag, "", "Password for the user. Must contain lower, upper, number, and special characters (min 12 chars)")
	cmd.Flags().String(userTypeFlag, "", "Type of user. One of 'intake' or 'dead-letter'")
	cmd.Flags().StringToString(labelsFlag, nil, `Labels in key=value format, separated by commas. Example: --labels "key1=value1,key2=value2".`)

	err := flags.MarkFlagsRequired(cmd, intakeIdFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	userId := inputArgs[0]

	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &cliErr.ProjectIdError{}
	}

	model := &inputModel{
		GlobalFlagModel: globalFlags,
		IntakeId:        flags.FlagToStringValue(p, cmd, intakeIdFlag),
		UserId:          userId,
		DisplayName:     flags.FlagToStringPointer(p, cmd, displayNameFlag),
		Description:     flags.FlagToStringPointer(p, cmd, descriptionFlag),
		Password:        flags.FlagToStringPointer(p, cmd, passwordFlag),
		UserType:        flags.FlagToStringPointer(p, cmd, userTypeFlag),
		Labels:          flags.FlagToStringToStringPointer(p, cmd, labelsFlag),
	}

	if model.DisplayName == nil && model.Description == nil && model.Password == nil && model.UserType == nil && model.Labels == nil {
		return nil, &cliErr.EmptyUpdateError{}
	}

	p.DebugInputModel(model)
	return model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *intake.APIClient) intake.ApiUpdateIntakeUserRequest {
	req := apiClient.UpdateIntakeUser(ctx, model.ProjectId, model.Region, model.IntakeId, model.UserId)

	payload := intake.UpdateIntakeUserPayload{
		DisplayName: model.DisplayName,
		Description: model.Description,
		Password:    model.Password,
		Labels:      model.Labels,
	}

	if model.UserType != nil {
		userType := intake.UserType(*model.UserType)
		payload.Type = &userType
	}

	req = req.UpdateIntakeUserPayload(payload)
	return req
}

func outputResult(p *print.Printer, model *inputModel, resp *intake.IntakeUserResponse) error {
	return p.OutputResult(model.OutputFormat, resp, func() error {
		if resp == nil {
			p.Outputf("Triggered update of Intake User for intake %q, but no user ID was returned.\n", model.IntakeId)
			return nil
		}

		operationState := "Updated"
		if model.Async {
			operationState = "Triggered update of"
		}
		p.Outputf("%s Intake User for intake %q. User ID: %s\n", operationState, model.IntakeId, utils.PtrString(resp.Id))
		return nil
	})
}
