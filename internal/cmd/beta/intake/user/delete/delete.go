package delete

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
	userIdArg    = "USER_ID"
	intakeIdFlag = "intake-id"
)

// inputModel struct holds all the input parameters for the command
type inputModel struct {
	*globalflags.GlobalFlagModel
	IntakeId string
	UserId   string
}

// NewCmd creates a new cobra command for deleting an Intake User
func NewCmd(p *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("delete %s", userIdArg),
		Short: "Deletes an Intake User",
		Long:  "Deletes an Intake User.",
		Args:  args.SingleArg(userIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Delete an Intake User with ID "xxx" for Intake "yyy"`,
				`$ stackit beta intake user delete xxx --intake-id yyy`),
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

			prompt := fmt.Sprintf("Are you sure you want to delete Intake User %q?", model.UserId)
			err = p.Printer.PromptForConfirmation(prompt)
			if err != nil {
				return err
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			if err = req.Execute(); err != nil {
				return fmt.Errorf("delete Intake User: %w", err)
			}

			// Wait for async operation, if async mode not enabled
			if !model.Async {
				s := spinner.New(p.Printer)
				s.Start("Deleting STACKIT Intake User")
				_, err = wait.DeleteIntakeUserWaitHandler(ctx, apiClient, model.ProjectId, model.Region, model.IntakeId, model.UserId).WaitWithContext(ctx)
				if err != nil {
					return fmt.Errorf("wait for STACKIT Intake User deletion: %w", err)
				}
				s.Stop()
			}

			operationState := "Deleted"
			if model.Async {
				operationState = "Triggered deletion of"
			}
			p.Printer.Outputf("%s STACKIT Intake User %s\n", operationState, model.UserId)

			return nil
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), intakeIdFlag, "Intake ID")

	err := flags.MarkFlagsRequired(cmd, intakeIdFlag)
	cobra.CheckErr(err)
}

// parseInput parses the command arguments and flags into a standardized model
func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	userId := inputArgs[0]

	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &cliErr.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		IntakeId:        flags.FlagToStringValue(p, cmd, intakeIdFlag),
		UserId:          userId,
	}

	p.DebugInputModel(model)
	return &model, nil
}

// buildRequest creates the API request to delete an Intake User
func buildRequest(ctx context.Context, model *inputModel, apiClient *intake.APIClient) intake.ApiDeleteIntakeUserRequest {
	req := apiClient.DeleteIntakeUser(ctx, model.ProjectId, model.Region, model.IntakeId, model.UserId)
	return req
}
