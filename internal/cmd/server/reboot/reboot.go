package reboot

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/types"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/iaas/client"
	iaasUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/iaas/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/iaas"

	"github.com/spf13/cobra"
)

const (
	serverIdArg = "SERVER_ID"

	hardRebootFlag    = "hard"
	defaultHardReboot = false
	hardRebootAction  = "hard"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	ServerId   string
	HardReboot bool
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("reboot %s", serverIdArg),
		Short: "Reboots a server",
		Long:  "Reboots a server.",
		Args:  args.SingleArg(serverIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Perform a soft reboot of a server with ID "xxx"`,
				"$ stackit server reboot xxx",
			),
			examples.NewExample(
				`Perform a hard reboot of a server with ID "xxx"`,
				"$ stackit server reboot xxx --hard",
			),
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

			serverLabel, err := iaasUtils.GetServerName(ctx, apiClient, model.ProjectId, model.Region, model.ServerId)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get server name: %v", err)
				serverLabel = model.ServerId
			} else if serverLabel == "" {
				serverLabel = model.ServerId
			}
			prompt := fmt.Sprintf("Are you sure you want to reboot server %q?", serverLabel)
			err = params.Printer.PromptForConfirmation(prompt)
			if err != nil {
				return err
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			err = req.Execute()
			if err != nil {
				return fmt.Errorf("server reboot: %w", err)
			}

			params.Printer.Info("Server %q rebooted\n", serverLabel)

			return nil
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().BoolP(hardRebootFlag, "b", defaultHardReboot, "Performs a hard reboot. (default false)")
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	serverId := inputArgs[0]

	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		ServerId:        serverId,
		HardReboot:      flags.FlagToBoolValue(p, cmd, hardRebootFlag),
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *iaas.APIClient) iaas.ApiRebootServerRequest {
	req := apiClient.RebootServer(ctx, model.ProjectId, model.Region, model.ServerId)
	// if hard reboot is set the action must be set (soft is default)
	if model.HardReboot {
		req = req.Action(hardRebootAction)
	}
	return req
}
