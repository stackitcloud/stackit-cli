package create

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
	intakeIdFlag    = "intake-id"
	displayNameFlag = "display-name"
	passwordFlag    = "password"
	descriptionFlag = "description"
	typeFlag        = "type"
	labelsFlag      = "labels"
)

// inputModel struct holds all the input parameters for the command
type inputModel struct {
	*globalflags.GlobalFlagModel

	IntakeId    *string
	DisplayName *string
	Password    *string
	Description *string
	Type        *string
	Labels      *map[string]string
}

func NewCmd(p *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Creates a new Intake User",
		Long:  "Creates a new Intake User, providing secure access credentials for applications to connect to a data stream.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Create a new Intake User with a display name and password for a specific Intake`,
				`$ stackit beta intake user create --intake-id xxx --display-name my-intake-user --password "my-secret-password"`),
			examples.NewExample(
				`Create a new dead-letter queue user with a description and labels`,
				`$ stackit beta intake user create --intake-id xxx --display-name my-dlq-reader --password "another-secret" --type "dead-letter" --description "User for reading undelivered messages" --labels "owner=team-alpha,scope=dlq"`),
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

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to create a new User for Intake %q?", *model.IntakeId)
				err = p.Printer.PromptForConfirmation(prompt)
				if err != nil {
					return err
				}
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
				s.Start("Creating STACKIT Intake User instance")
				_, err = wait.CreateOrUpdateIntakeUserWaitHandler(ctx, apiClient, model.ProjectId, model.Region, *model.IntakeId, resp.GetId()).WaitWithContext(ctx)
				if err != nil {
					return fmt.Errorf("wait for STACKIT Instance creation: %w", err)
				}
				s.Stop()
			}

			return outputResult(p.Printer, model, *model.IntakeId, resp)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), intakeIdFlag, "ID of the Intake to which the user belongs")
	cmd.Flags().String(displayNameFlag, "", "Display name")
	cmd.Flags().String(passwordFlag, "", "User password")
	cmd.Flags().String(descriptionFlag, "", "Description")
	cmd.Flags().String(typeFlag, "", "Type of user, 'intake' for writing to the stream or 'dead-letter' for reading from the dead-letter queue")
	cmd.Flags().StringToString(labelsFlag, nil, "Labels in key=value format, separated by commas. Example: --labels \"key1=value1,key2=value2\"")

	err := flags.MarkFlagsRequired(cmd, intakeIdFlag, displayNameFlag, passwordFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &cliErr.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		IntakeId:        flags.FlagToStringPointer(p, cmd, intakeIdFlag),
		DisplayName:     flags.FlagToStringPointer(p, cmd, displayNameFlag),
		Password:        flags.FlagToStringPointer(p, cmd, passwordFlag),
		Description:     flags.FlagToStringPointer(p, cmd, descriptionFlag),
		Type:            flags.FlagToStringPointer(p, cmd, typeFlag),
		Labels:          flags.FlagToStringToStringPointer(p, cmd, labelsFlag),
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *intake.APIClient) intake.ApiCreateIntakeUserRequest {
	req := apiClient.CreateIntakeUser(ctx, model.ProjectId, model.Region, *model.IntakeId)

	// Build main payload
	payload := intake.CreateIntakeUserPayload{
		DisplayName: model.DisplayName,
		Password:    model.Password,
		Description: model.Description,
		Labels:      model.Labels,
	}

	if model.Type != nil {
		payload.Type = (*intake.UserType)(model.Type)
	}

	req = req.CreateIntakeUserPayload(payload)
	return req
}

func outputResult(p *print.Printer, model *inputModel, intakeId string, resp *intake.IntakeUserResponse) error {
	return p.OutputResult(model.OutputFormat, resp, func() error {
		if resp == nil {
			p.Outputf("Created Intake User for Intake %q, but no intake ID was returned.\n", intakeId)
			return nil
		}

		operationState := "Created"
		if model.Async {
			operationState = "Triggered creation of"
		}
		p.Outputf("%s Intake User for Intake %q. User ID: %s\n", operationState, intakeId, utils.PtrString(resp.Id))
		return nil
	})
}
