package create

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
	"github.com/stackitcloud/stackit-cli/internal/pkg/projectname"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/intake/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/spinner"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/intake"
	"github.com/stackitcloud/stackit-sdk-go/services/intake/wait"
)

const (
	displayNameFlag = "display-name"
	intakeIdFlag    = "intake-id"
	passwordFlag    = "password"
	userTypeFlag    = "type"
	descriptionFlag = "description"
	labelsFlag      = "labels"
)

// inputModel struct holds all the input parameters for the command
type inputModel struct {
	*globalflags.GlobalFlagModel
	DisplayName *string
	IntakeId    *string
	Password    *string
	UserType    *string
	Description *string
	Labels      *map[string]string
}

func NewCmd(p *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Creates a new Intake User",
		Long:  "Creates a new Intake User for a specific Intake.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Create a new Intake User with required parameters`,
				`$ stackit beta intake user create --display-name intake-user --intake-id xxx --password "SuperSafepass123\!"`),
			examples.NewExample(
				`Create a new Intake User for the dead-letter queue with labels`,
				`$ stackit beta intake user create --display-name dlq-user --intake-id xxx --password "SuperSafepass123\!" --type dead-letter --labels "env=prod"`),
		),
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := context.Background()
			model, err := parseInput(p.Printer, cmd)
			if err != nil {
				return err
			}

			// Configure API client
			apiClient, err := client.ConfigureClient(p.Printer, p.CliVersion)
			if err != nil {
				return err
			}

			projectLabel, err := projectname.GetProjectName(ctx, p.Printer, p.CliVersion, cmd)
			if err != nil {
				p.Printer.Debug(print.ErrorLevel, "get project name: %v", err)
				projectLabel = model.ProjectId
			}

			prompt := fmt.Sprintf("Are you sure you want to create an Intake User for project %q?", projectLabel)
			err = p.Printer.PromptForConfirmation(prompt)
			if err != nil {
				return err
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("create Intake User: %w", err)
			}

			// Wait for async operation, if async mode not enabled
			if !model.Async {
				s := spinner.New(p.Printer)
				s.Start("Creating STACKIT Intake User")
				_, err = wait.CreateOrUpdateIntakeUserWaitHandler(ctx, apiClient, model.ProjectId, model.Region, *model.IntakeId, resp.GetId()).WaitWithContext(ctx)
				if err != nil {
					return fmt.Errorf("wait for STACKIT Intake User creation: %w", err)
				}
				s.Stop()
			}

			return outputResult(p.Printer, model, projectLabel, resp)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().String(displayNameFlag, "", "Display name")
	cmd.Flags().Var(flags.UUIDFlag(), intakeIdFlag, "The UUID of the Intake to associate the user with")
	cmd.Flags().String(passwordFlag, "", "Password for the user. Must contain lower, upper, number, and special characters (min 12 chars)")
	cmd.Flags().String(userTypeFlag, string(intake.USERTYPE_INTAKE), "Type of user. One of 'intake' (default) or 'dead-letter'")
	cmd.Flags().String(descriptionFlag, "", "Description")
	cmd.Flags().StringToString(labelsFlag, nil, "Labels in key=value format, separated by commas")

	err := flags.MarkFlagsRequired(cmd, displayNameFlag, intakeIdFlag, passwordFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &cliErr.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		DisplayName:     flags.FlagToStringPointer(p, cmd, displayNameFlag),
		IntakeId:        flags.FlagToStringPointer(p, cmd, intakeIdFlag),
		Password:        flags.FlagToStringPointer(p, cmd, passwordFlag),
		UserType:        flags.FlagToStringPointer(p, cmd, userTypeFlag),
		Description:     flags.FlagToStringPointer(p, cmd, descriptionFlag),
		Labels:          flags.FlagToStringToStringPointer(p, cmd, labelsFlag),
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *intake.APIClient) intake.ApiCreateIntakeUserRequest {
	req := apiClient.CreateIntakeUser(ctx, model.ProjectId, model.Region, *model.IntakeId)

	var userType *intake.UserType
	if model.UserType != nil {
		userType = utils.Ptr(intake.UserType(*model.UserType))
	}

	payload := intake.CreateIntakeUserPayload{
		DisplayName: model.DisplayName,
		Password:    model.Password,
		Type:        userType,
		Description: model.Description,
		Labels:      model.Labels,
	}

	req = req.CreateIntakeUserPayload(payload)
	return req
}

func outputResult(p *print.Printer, model *inputModel, projectLabel string, resp *intake.IntakeUserResponse) error {
	return p.OutputResult(model.OutputFormat, resp, func() error {
		if resp == nil {
			p.Outputf("Triggered creation of Intake User for project %q, but no user ID was returned.\n", projectLabel)
			return nil
		}

		operationState := "Created"
		if model.Async {
			operationState = "Triggered creation of"
		}
		p.Outputf("%s Intake User for project %q. User ID: %s\n", operationState, projectLabel, utils.PtrString(resp.Id))
		return nil
	})
}
