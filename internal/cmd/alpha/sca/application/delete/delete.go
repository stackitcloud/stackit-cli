package delete

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	sca "github.com/stackitcloud/stackit-sdk-go/services/sca/v1alphaapi"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/projectname"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/sca/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/sca/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
)

const (
	applicationIDArg  = "APPLICATION_ID"
	environmentIDFlag = "environment-id"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	EnvironmentID string
	ApplicationID string
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a SCA application",
		Long:  "Delete a STACKIT Container Applications (SCA) application.",
		Args:  args.SingleArg(applicationIDArg, nil),
		Example: examples.Build(
			examples.NewExample(
				`Delete a SCA application with ID "xxx" from an environment with ID "yyy"`,
				"$ stackit alpha sca application delete xxx --environment-id yyy"),
		),
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

			globalFlags := globalflags.Parse(params.Printer, cmd)
			if globalFlags.ProjectId == "" {
				return &errors.ProjectIdError{}
			}

			applicationLabel, err := utils.GetApplicationName(ctx, apiClient.DefaultAPI, model.ProjectId, model.EnvironmentID, model.ApplicationID)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get application name: %v", err)
				applicationLabel = model.ApplicationID
			}

			environmentLabel, err := utils.GetEnvironmentName(ctx, apiClient.DefaultAPI, model.ProjectId, model.EnvironmentID)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get environemnt name: %v", err)
				environmentLabel = model.EnvironmentID
			}

			projectLabel, err := projectname.GetProjectName(ctx, params.Printer, params.CliVersion, cmd)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get project name: %v", err)
				projectLabel = model.ProjectId
			}

			prompt := fmt.Sprintf("Are you sure you want to delete the application %q from environment %q and product %q?", applicationLabel, environmentLabel, projectLabel)
			err = params.Printer.PromptForConfirmation(prompt)
			if err != nil {
				return err
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			_, err = req.Execute()
			if err != nil {
				return fmt.Errorf("delete application: %w", err)
			}

			params.Printer.Info("Deleted application %q from environment %q\n", applicationLabel, environmentLabel)

			return nil
		},
	}

	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), environmentIDFlag, "Environment ID (uses default environment if not set)")
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	applicationID := inputArgs[0]

	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	environmentID := flags.FlagToStringValue(p, cmd, environmentIDFlag)
	if environmentID == "" {
		environmentID = globalFlags.ProjectId
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		EnvironmentID:   environmentID,
		ApplicationID:   applicationID,
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *sca.APIClient) sca.ApiDeleteApplicationRequest {
	return apiClient.DefaultAPI.DeleteApplication(ctx, model.ProjectId, model.EnvironmentID, model.ApplicationID)
}
