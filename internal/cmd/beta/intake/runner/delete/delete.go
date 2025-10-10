package delete

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	cliErr "github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/intake/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/intake"
)

const (
	runnerIdArg = "RUNNER_ID"
)

// inputModel struct holds all the input parameters for the command
type inputModel struct {
	*globalflags.GlobalFlagModel
	RunnerId string
}

// NewDeleteCmd creates a new cobra command for deleting an Intake Runner
func NewDeleteCmd(p *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("delete %s", runnerIdArg),
		Short: "Deletes an Intake Runner",
		Long:  "Deletes an Intake Runner.",
		Args:  args.SingleArg(runnerIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Delete an Intake Runner with ID "xxx"`,
				`$ stackit beta intake runner delete xxx`),
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
				prompt := fmt.Sprintf("Are you sure you want to delete Intake Runner %q?", model.RunnerId)
				err = p.Printer.PromptForConfirmation(prompt)
				if err != nil {
					return err
				}
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			if err = req.Execute(); err != nil {
				return fmt.Errorf("delete Intake Runner: %w", err)
			}
			p.Printer.Info("Deleted runner %s\n", model.RunnerId)

			return nil
		},
	}
	return cmd
}

// parseInput parses the command arguments and flags into a standardized model
func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	runnerId := inputArgs[0]

	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &cliErr.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		RunnerId:        runnerId,
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

// buildRequest creates the API request to delete an Intake Runner
func buildRequest(ctx context.Context, model *inputModel, apiClient *intake.APIClient) intake.ApiDeleteIntakeRunnerRequest {
	req := apiClient.DeleteIntakeRunner(ctx, model.ProjectId, model.Region, model.RunnerId)
	return req
}
