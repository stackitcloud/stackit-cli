package delete

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/intake/wait"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	cliErr "github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/intake/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/spinner"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/intake"
)

const (
	intakeIdArg = "INTAKE_ID"
)

// inputModel struct holds all the input parameters for the command
type inputModel struct {
	*globalflags.GlobalFlagModel
	IntakeId string
}

// NewCmd creates a new cobra command for deleting an Intake
func NewCmd(p *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("delete %s", intakeIdArg),
		Short: "Deletes an Intake",
		Long:  "Deletes an Intake.",
		Args:  args.SingleArg(intakeIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Delete an Intake with ID "xxx"`,
				`$ stackit beta intake delete xxx`),
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

			prompt := fmt.Sprintf("Are you sure you want to delete an Intake %q?", model.IntakeId)
			err = p.Printer.PromptForConfirmation(prompt)
			if err != nil {
				return err
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			if err = req.Execute(); err != nil {
				return fmt.Errorf("delete Intake: %w", err)
			}

			// Wait for async operation, if async mode not enabled
			if !model.Async {
				s := spinner.New(p.Printer)
				s.Start("Deleting STACKIT Intake instance")
				_, err = wait.DeleteIntakeWaitHandler(ctx, apiClient, model.ProjectId, model.Region, model.IntakeId).WaitWithContext(ctx)
				if err != nil {
					return fmt.Errorf("wait for STACKIT Instance deletion: %w", err)
				}
				s.Stop()
			}

			operationState := "Deleted"
			if model.Async {
				operationState = "Triggered deletion of"
			}
			p.Printer.Outputf("%s stackit Intake instance %s \n", operationState, model.IntakeId)

			return nil
		},
	}
	return cmd
}

// parseInput parses the command arguments and flags into a standardized model
func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	intakeId := inputArgs[0]

	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &cliErr.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		IntakeId:        intakeId,
	}

	p.DebugInputModel(model)
	return &model, nil
}

// buildRequest creates the API request to delete an Intake
func buildRequest(ctx context.Context, model *inputModel, apiClient *intake.APIClient) intake.ApiDeleteIntakeRequest {
	req := apiClient.DeleteIntake(ctx, model.ProjectId, model.Region, model.IntakeId)
	return req
}
