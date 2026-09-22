package delete

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/sca/client"
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
		Long:  "Delete a STACKIT Kubernetes Engine (SCA) application.",
		Args:  args.SingleArg(applicationIDArg, nil),
		Example: examples.Build(
			examples.NewExample(
				`Delete a SCA application with ID "xxx" from an environment with ID "yyy"`,
				"$ stackit sca application delete xxx --environment-id yyy"),
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

			// Call API
			req := apiClient.DefaultAPI.DeleteApplication(ctx, model.ProjectId, model.EnvironmentID, model.ApplicationID)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("delete application: %w", err)
			}

			params.Printer.Info("Deleted application %q from environment %q\n", resp.GetId(), resp.GetEnvironmentId())

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
