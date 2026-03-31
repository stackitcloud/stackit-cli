package login

import (
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/auth"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"

	"github.com/spf13/cobra"
)

const (
	portFlag          = "port"
	useDeviceFlowFlag = "use-device-flow"
)

type inputModel struct {
	Port          *int
	UseDeviceFlow bool
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "login",
		Short: "Logs in to the STACKIT CLI",
		Long: fmt.Sprintf("%s\n%s",
			"Logs in to the STACKIT CLI using a user account.",
			"By default, the authentication uses a web-based authorization flow and opens a browser window where you can login to your STACKIT account. You can alternatively use the OAuth 2.0 device flow for environments where receiving a local callback is not possible."),
		Args: args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Login to the STACKIT CLI. This command will open a browser window where you can login to your STACKIT account`,
				"$ stackit auth login"),
			examples.NewExample(
				`Login to the STACKIT CLI using OAuth 2.0 device flow`,
				"$ stackit auth login --use-device-flow"),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			model, err := parseInput(params.Printer, cmd, args)
			if err != nil {
				return err
			}

			err = auth.AuthorizeUser(params.Printer, auth.UserAuthConfig{
				IsReauthentication: false,
				Port:               model.Port,
				UseDeviceFlow:      model.UseDeviceFlow,
			})
			if err != nil {
				return fmt.Errorf("authorization failed: %w", err)
			}

			params.Printer.Outputln("Successfully logged into STACKIT CLI.\n")

			return nil
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Int(portFlag, 0,
		"The port on which the callback server will listen to. By default, it tries to bind a port between 8000 and 8020.\n"+
			"When a value is specified, it will only try to use the specified port. Valid values are within the range of 8000 to 8020.",
	)
	cmd.Flags().Bool(useDeviceFlowFlag, false,
		"Use OAuth 2.0 device authorization grant (device flow) instead of the browser callback flow.")
}

func parseInput(p *print.Printer, cmd *cobra.Command, _ []string) (*inputModel, error) {
	port := flags.FlagToIntPointer(p, cmd, portFlag)
	useDeviceFlow := flags.FlagToBoolValue(p, cmd, useDeviceFlowFlag)
	// For the CLI client only callback URLs with localhost:[8000-8020] are valid. Additional callbacks must be enabled in the backend.
	if port != nil && (*port < 8000 || 8020 < *port) {
		return nil, fmt.Errorf("port must be between 8000 and 8020")
	}

	model := inputModel{
		Port:          port,
		UseDeviceFlow: useDeviceFlow,
	}

	p.DebugInputModel(model)
	return &model, nil
}
